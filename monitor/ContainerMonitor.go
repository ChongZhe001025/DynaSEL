package monitor

import (
	"DynaSEL-latest/policy"
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
)

func MonitorMobyDir(strConfigDirPath string) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}

	eventCh, errCh := cli.Events(context.Background(), events.ListOptions{})

	fmt.Println("Listening for Docker events...")

	for {
		select {
		case event := <-eventCh:
			if event.Type == "container" {
				switch event.Action {
				case "create":
					fmt.Printf("Container created: ID=%s Name=%s\n", event.ID, event.Actor.Attributes["name"])
					policy.CreateSElinuxPolicyCil(strConfigDirPath, event.ID)
				case "destroy":
					fmt.Printf("Container destroyed: ID=%s Name=%s\n", event.ID, event.Actor.Attributes["name"])
				}
			}

		case err := <-errCh:
			if err != nil {
				log.Fatalf("Error while listening for events: %v", err)
			}
		}
	}
}
