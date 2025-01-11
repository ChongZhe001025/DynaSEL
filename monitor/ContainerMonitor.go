package monitor

import (
	"DynaSEL-latest/policy"
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
)

func MonitorContainers(strConfigDirPath string) {
	// 建立 Docker 客戶端
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}

	for {
		// 創建新的 Context 和事件通道
		ctx, cancel := context.WithCancel(context.Background())

		eventCh, errCh := cli.Events(ctx, events.ListOptions{})
		fmt.Println("Listening for Docker events...")

		for {
			select {
			case event := <-eventCh:
				if event.Type == "container" {
					switch event.Action {
					case "create":
						fmt.Printf("Container created: ID=%s Name=%s\n", event.ID, event.Actor.Attributes["name"])

						cancel()
						fmt.Println("Monitoring paused...")
						policy.CreateSELinuxPolicyCil(strConfigDirPath, event.ID)

						goto RestartMonitor
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
	RestartMonitor:
		fmt.Println("Resuming monitoring...")
	}
}
