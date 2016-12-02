package main

import "github.com/docker/go-dockercloud/dockercloud"
import "github.com/aws/aws-sdk-go/service/ec2"
import "github.com/aws/aws-sdk-go/aws/session"
import "log"
import "os"
import "strings"

/*
 * Return an environment variable. If it's not set - crash
 */
func getEnv(name string) string {
	value := os.Getenv(name)

	if len(value) == 0 {
		panic("Please set " + name + " variable")
	}

	return value
}

/*
 * Returns an array of public IP addresses  
 */
func getNodeIps() []string {
	nodeList, err := dockercloud.ListNodes()

	if err != nil {
  		panic(err)
	}

	log.Println("Received public IP list from Docker Cloud")
	
	nodeIps := make([]string, 0)

	if len(nodeList.Objects) == 0 {
		log.Println("There are no nodes in your Docker Cloud account yet")
	}

	for i := 0; i < len(nodeList.Objects); i++ {
		if len(nodeList.Objects[i].Public_ip) > 0 {
			nodeIps = append(nodeIps, nodeList.Objects[i].Public_ip + "/32")
		}
    }

    return nodeIps
}

/*
 * Infinite loop - listening to Docker Cloud events
 */
func listenToEvents() {
	log.Println("Listening to Docker Cloud events")

	c := make(chan dockercloud.Event)
	e := make(chan error)

	go dockercloud.Events(c, e)

	for {
    	select {
        	case event := <-c:
            	log.Println(event)
	        case err := <-e:
    	        log.Println(err)
    	}
	} 
}

/*
 * Rewrite inbound rules for the security group
 */
func modifySecurityGroup(groupId string, ips []string) {
	var newRules ec2.AuthorizeSecurityGroupIngressInput
	var oldRules ec2.RevokeSecurityGroupIngressInput
    var allProtocol string = "-1"

    log.Println("Flushing security group... ")

	svc := ec2.New(session.New())
    newRules.GroupId = &groupId
    oldRules.GroupId = &groupId

	params := &ec2.DescribeSecurityGroupsInput{ GroupIds: []*string{ &groupId }}
	resp, err := svc.DescribeSecurityGroups(params)

	for i := 0; i < len(resp.SecurityGroups[0].IpPermissions); i++ {
		existing := resp.SecurityGroups[0].IpPermissions[i]
		oldRules.IpPermissions = append(oldRules.IpPermissions, existing)
    }
	
	_, err = svc.RevokeSecurityGroupIngress(&oldRules)
	if err == nil {
		log.Println("success")
	}

	log.Println("Adding current node IPs to the group... ")

	// Add AWS internal network
	ips = append(ips, "10.0.0.0/8")

	for i := 0; i < len(ips); i++ {
	    entry := new(ec2.IpPermission)
    	entry.IpProtocol = &allProtocol
	    entry.IpRanges = []*ec2.IpRange{{CidrIp: &ips[i]}}
    	newRules.IpPermissions = append(newRules.IpPermissions, entry)
	}

    _, err = svc.AuthorizeSecurityGroupIngress(&newRules)
    if err != nil {
		panic(err)
	}

	log.Println("success")
}

/*
 * Initialize Docker Cloud SDK
 */
func initDockerCloud() {
	dockercloud.User = getEnv("DOCKER_CLOUD_USER")
	dockercloud.ApiKey = getEnv("DOCKER_CLOUD_KEY")

	if len(os.Getenv("DOCKER_CLOUD_NAMESPACE")) > 0 {
		dockercloud.Namespace = os.Getenv("DOCKER_CLOUD_NAMESPACE")
	}
}

/*
 * Main block
 */
func main() {
	initDockerCloud()

	groups := strings.Split(getEnv("AWS_SG_ID"), ',')

	for i := 0; i < len(groups); i++ {
		modifySecurityGroup(groups[i], getNodeIps())
	}

	listenToEvents()
}
