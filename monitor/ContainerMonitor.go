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

	// 創建一個可取消的 Context
	ctx, cancel := context.WithCancel(context.Background())

	for {
		// 啟動監控事件的 Goroutine
		go func() {
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

							fmt.Println("Resuming monitoring...")
							ctx, cancel = context.WithCancel(context.Background())
							return
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
		}()

		<-ctx.Done()
	}
}
