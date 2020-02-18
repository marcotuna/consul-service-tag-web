package main

import (
	"consul-service-tag-web/service"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func main() {

	// Read Service Tags from Environment
	configServiceTags := strings.Split(os.Getenv("SERVICE_TAGS"), ",")

	consulClient, err := service.NewConsulClient()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	serviceTags, err := consulClient.GetConsulServices(configServiceTags)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	serviceTagsJSON, err := json.Marshal(serviceTags)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(serviceTagsJSON))
}
