package main

import "github.com/docker/go-dockercloud/dockercloud"
import "github.com/aws/aws-sdk-go/service/ec2"
import "github.com/aws/aws-sdk-go/aws/session"
import "log"
import "os"

/*
 * Return an environment variable. If it's not set - crash
 */
func getEnv(name string) string {
	value := os.Getenv(name)

	if len(value) == 0 {
		log.Println("Please set " + name + " variable")
		os.Exit(1)
	}

	return value
}

/*
 * Returns an array of public IP addresses  
 */
func getNodeIps() []string {
	nodeList, err := dockercloud.ListNodes()

	if err != nil {
  		log.Println(err)
	} else {
		log.Println("Received public IP list from Docker Cloud")
	}

	nodeIps := make([]string, 0)

	if len(nodeList.Objects) == 0 {
		log.Println("There are no nodes in your Docker Cloud account")
		os.Exit(1)
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
	var inboundRules ec2.AuthorizeSecurityGroupIngressInput
	var flushRules ec2.RevokeSecurityGroupIngressInput
    var allProtocol string = "-1"

    log.Println("Flushing security group... ")

	svc := ec2.New(session.New())
    inboundRules.GroupId = &groupId
    flushRules.GroupId = &groupId

	params := &ec2.DescribeSecurityGroupsInput{ GroupIds: []*string{ &groupId }}
	resp, err := svc.DescribeSecurityGroups(params)

	for i := 0; i < len(resp.SecurityGroups[0].IpPermissions); i++ {
		existing := resp.SecurityGroups[0].IpPermissions[i]
		flushRules.IpPermissions = append(flushRules.IpPermissions, existing)
    }
	
	_, err = svc.RevokeSecurityGroupIngress(&flushRules)

	if err == nil {
		log.Println("done")
	}

	log.Println("Adding current node IPs to the group... ")

	// Add AWS internal network
	ips = append(ips, "10.0.0.0/8")

	for i := 0; i < len(ips); i++ {
	    entry := new(ec2.IpPermission)
    	entry.IpProtocol = &allProtocol
	    entry.IpRanges = []*ec2.IpRange{{CidrIp: &ips[i]}}
    	inboundRules.IpPermissions = append(inboundRules.IpPermissions, entry)
	}

    _, err = svc.AuthorizeSecurityGroupIngress(&inboundRules)
    
    if err != nil {
		panic(err)
	} else {
		log.Println("done")
	}
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

func main() {
	initDockerCloud()

	modifySecurityGroup(getEnv("AWS_SG_ID"), getNodeIps())

	listenToEvents()
}
