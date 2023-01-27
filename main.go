package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var containerList = make(map[string]string)

func main() {
	ctx := context.Background()
	monitoringDatabaseExists := false

	for {
		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			log.Fatal(err)
		}

		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
		if err != nil {
			log.Fatal(err)
		}

		//Check if Database Exisits
		for _, container := range containers {
			if container.ID == "MonitoringDatabase" {
				monitoringDatabaseExists = true
			}
		}

		if !monitoringDatabaseExists {
			sendToWebsite()
		}
		time.Sleep(time.Second * 3)
	}
}

func sendToDatabase() {
	//Send to Database
	//Insert Database Code
}

func sendToWebsite() {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	//list of streams
	//running loop = Iterate through streams (also check if new containers exist or old ones are off (than close stream))
	options := types.ContainerLogsOptions{ShowStdout: true, Timestamps: true}

	var lastTimestampMap = make(map[string]time.Time) // should probably be a map, mapping ID to Scanner for the container

	for k := range lastTimestampMap {
		delete(lastTimestampMap, k)
	}

	for _, container := range containers {
		out, err := cli.ContainerLogs(ctx, container.ID, options)

		if err != nil {
			log.Fatal(err)
		}

		if container.Image != "docker-monitoring-sw" {
			var scanner = bufio.NewScanner(out)
			//scannerMap[container.ID] = scanner

			// first time get all logs right away
			for scanner.Scan() {
				// do something with the old logs

			}
		}
	}

	/*
		for containerId, scanner := range scannerMap {
			//For each "checkup" run, (main loop run) -> retrieve new container list, for each container in the new list, check if their ID is already a Key in scannerMap
			//Depending on that handle creation/removal (create a scanner or close a scanner)

			fmt.Println("Log for : " + containerId)
			for scanner.Scan() {
				fmt.Printf(scanner.Text())
			}
	*/

	for _, container := range containers {
		out, err := cli.ContainerLogs(ctx, container.ID, options)
		if err != nil {
			log.Fatal(err)
		}
		if container.Image != "docker-monitoring-sw" {
			fmt.Printf("The \"" + container.Image + "\" container, with the ID \"" + container.ID[:10] + "\" logged: ")
			fmt.Println()
			buf := new(strings.Builder)
			io.Copy(buf, out)

			scanner := bufio.NewScanner(strings.NewReader(buf.String()))
			for scanner.Scan() {
				inputArray := strings.Split(scanner.Text(), " ")
				var timeStamp string = inputArray[0]
				for i := 1; i < len(inputArray); i++ {
					fmt.Print(inputArray[i] + " ")
				}
				if val, ok := containerList[container.ID]; ok {
					//Check if already exisiting TimeStamp is before new one
					timeStampNew, _ := time.Parse(time.RFC3339, timeStamp)
					timeStampOld, _ := time.Parse(time.RFC3339, val)

					if timeStampOld.Unix() < timeStampNew.Unix() {
						containerList[container.ID] = timeStamp
						//Send Currend Log
						for i := 1; i < len(inputArray); i++ {
							fmt.Print(inputArray[i] + " ")
						}
						fmt.Println()
					}
				} else {
					containerList[container.ID] = timeStamp
				}
			}
		}
	}
	fmt.Println()
	for key, value := range containerList {
		fmt.Println("key: ", key, " Value: ", value)
	}
	fmt.Println()
}
