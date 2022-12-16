package main

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	options := types.ContainerLogsOptions{ShowStdout: true}
	for _, container := range containers {
		out, err := cli.ContainerLogs(ctx, container.ID, options)
		if err != nil {
			panic(err)
		}

		fmt.Printf("The \"" + container.Image + "\" container, with the ID \"" + container.ID + "\" logged: ")
		fmt.Println()

		buf := new(strings.Builder)
		io.Copy(buf, out)
		fmt.Println(buf.String())
	}
	time.Sleep(time.Minute * 3)

}
