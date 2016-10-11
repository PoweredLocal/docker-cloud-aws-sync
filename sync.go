package main

import "github.com/docker/go-dockercloud/dockercloud"
import "github.com/aws/aws-sdk-go/service/ec2"
import "github.com/aws/aws-sdk-go/aws/session"
import "log"
import "fmt"
import "os"
//import "strings"

func getEnv(name string) string {
	value := os.Getenv(name)

	if len(value) == 0 {
		log.Println("Please set " + name + " variable")
		os.Exit(1)
	}

	return value
}

func main() {
	dockercloud.User = getEnv("DOCKER_CLOUD_USER")
	dockercloud.ApiKey = getEnv("DOCKER_CLOUD_KEY")
	dockercloud.Namespace = getEnv("DOCKER_CLOUD_NAMESPACE")
	sg := getEnv("AWS_SG_ID")	

	nodeList, err := dockercloud.ListNodes()

	if err != nil {
  		log.Println(err)
	}

	nodeIps := make([]string, 0)

	log.Println("Nodes:")

	for i := 0; i < len(nodeList.Objects); i++ {
		nodeIps = append(nodeIps, nodeList.Objects[i].Public_ip)
    }

    fmt.Printf("IPs are:\n%+v\n", nodeIps)
    log.Println(sg)

	//c := make(chan dockercloud.Event)
	//e := make(chan error)

	svc := ec2.New(session.New())

	// Call the DescribeInstances Operation
	resp, err := svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}

	// resp has all of the response data, pull out instance IDs:
	fmt.Println("> Number of reservation sets: ", len(resp.Reservations))
	for idx, res := range resp.Reservations {
		fmt.Println("  > Number of instances: ", len(res.Instances))
		for _, inst := range resp.Reservations[idx].Instances {
			fmt.Println("    - Instance ID: ", *inst.InstanceId)
		}
	}

/*	go dockercloud.Events(c, e)

	for {
    	select {
        	case event := <-c:
            	log.Println(event)
	        case err := <-e:
    	        log.Println(err)
    	}
	} */
}
