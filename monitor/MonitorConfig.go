package monitor

import (
	"DynaSEL-latest/policy"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

var bolMonitorIsPaused = false

func MonitorConfigJson(strArrConfigParentDirPath []string) {

	fileChangeWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Error creating watcher: %v", err)
	}
	defer fileChangeWatcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-fileChangeWatcher.Events:
				if !ok {
					return
				}
				if !bolMonitorIsPaused {
					strContainerIdPath := event.Name

					intLastSlashIndex := strings.LastIndex(strContainerIdPath, "/")
					if intLastSlashIndex != -1 {
						strContainerIdPath = strContainerIdPath[:intLastSlashIndex]
					}

					bolMonitorIsPaused = true

					err := fileChangeWatcher.Remove(strContainerIdPath)
					if err != nil {
						log.Printf("Error removing watcher: %v", err)
					}

					createCilFile(strContainerIdPath)

					fmt.Println("Resuming monitoring...")
					err = fileChangeWatcher.Add(strContainerIdPath)
					if err != nil {
						log.Printf("Error adding watcher: %v", err)
					}
					bolMonitorIsPaused = false
				}

			case err, ok := <-fileChangeWatcher.Errors:
				if !ok {
					return
				}
				log.Printf("Error: %v\n", err)
			}
		}
	}()

	for _, strConfigParentDirPath := range strArrConfigParentDirPath {
		if _, err := os.Stat(strConfigParentDirPath); os.IsNotExist(err) {
			log.Fatalf("Directory does not exist: %s", strConfigParentDirPath)
		}

		err := fileChangeWatcher.Add(strConfigParentDirPath)
		if err != nil {
			log.Printf("Failed to watch file %s: %v", strConfigParentDirPath, err)
		} else {
			fmt.Printf("Started watching: %s\n", strConfigParentDirPath)
		}
	}
	select {}
}

func createCilFile(strContainerIdPath string) {
	lastSlashIndex := strings.LastIndex(strContainerIdPath, "/")
	strConfigDirPath := strContainerIdPath[:lastSlashIndex]
	strContainerID := strContainerIdPath[lastSlashIndex:]

	policy.CreateCilFile(strConfigDirPath, strContainerID)
	time.Sleep(5 * time.Second)
	fmt.Println("Task completed!")
}
