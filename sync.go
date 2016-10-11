package main

import "github.com/docker/go-dockercloud/dockercloud"
//import "github.com/aws/aws-sdk-go/service/ec2"
import "log"
import "fmt"
import "os"
//import "strings"

func main() {
	dockercloud.User = os.Getenv("DOCKER_CLOUD_USER")
	dockercloud.ApiKey = os.Getenv("DOCKER_CLOUD_KEY")
	dockercloud.Namespace = os.Getenv("DOCKER_CLOUD_NAMESPACE")

	nodeList, err := dockercloud.ListNodes()

	if err != nil {
  		log.Println(err)
	}

	log.Println("Nodes:")
	fmt.Printf("%+v", nodeList)

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
