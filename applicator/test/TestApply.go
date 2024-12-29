package test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func TestApplyPolicyToContainer(strContainerID string) {

	// 創建 Docker 客戶端
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}

	// 呼叫 ContainerInspect 取得容器資訊
	ctx := context.Background()
	jsonContainerInspect, err := cli.ContainerInspect(ctx, strContainerID)
	if err != nil {
		log.Fatalf("Failed to inspect container: %v", err)
	}

	strContainerName := jsonContainerInspect.Name
	if len(strContainerName) > 0 && strContainerName[0] == '/' {
		strContainerName = strContainerName[1:]
	}

	strContainerImage := jsonContainerInspect.Config.Image
	strContainerExportedPathName := "SysFiles/ExportedTarFiles/" + strContainerID + ".tar"

	// 1. 停止容器
	if err := stopContainer(cli, strContainerName); err != nil {
		log.Fatalf("Failed to stop container: %v", err)
	}

	// 2. 導出容器數據
	if err := exportContainer(cli, strContainerName, strContainerExportedPathName); err != nil {
		log.Fatalf("Failed to export container: %v", err)
	}

	// 3. 重新導入容器數據並創建新的映像
	if err := importContainer(cli, strContainerExportedPathName, strContainerImage); err != nil {
		log.Fatalf("Failed to import container: %v", err)
	}

	// 4. 刪除容器
	// if err := removeContainer(cli, strContainerID); err != nil {
	// 	log.Fatalf("Failed to remove container: %v", err)
	// }

	// 5. 創建並啟動新容器
	strLabelType := ("container_" + strContainerID)
	if err := createAndStartContainer(cli, strContainerImage, "new_"+strContainerName, strLabelType); err != nil {
		log.Fatalf("Failed to create and start new container: %v", err)
	}

	log.Println("New container created and started successfully!")
}

// 停止容器
func stopContainer(cli *client.Client, containerName string) error {
	ctx := context.Background()
	if err := cli.ContainerStop(ctx, containerName, container.StopOptions{}); err != nil {
		return fmt.Errorf("unable to stop container: %v", err)
	}
	log.Printf("Container %s stopped", containerName)
	return nil
}

// 導出容器
func exportContainer(cli *client.Client, containerName, outputPath string) error {
	ctx := context.Background()
	reader, err := cli.ContainerExport(ctx, containerName)
	if err != nil {
		return fmt.Errorf("failed to export container: %v", err)
	}
	defer reader.Close()

	// 使用 os.Create 替代 ioutil.WriteFile
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create export file: %v", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return fmt.Errorf("failed to write container data: %v", err)
	}
	log.Printf("Container %s exported to %s", containerName, outputPath)
	return nil
}

// 導入容器數據
func importContainer(cli *client.Client, importFile, containerImage string) error {
	ctx := context.Background()
	file, err := os.Open(importFile)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	_, err = cli.ImageLoad(ctx, file, true)
	if err != nil {
		return fmt.Errorf("failed to import container: %v", err)
	}
	log.Printf("Container data from %s imported as image %s", importFile, containerImage)
	return nil
}

// 刪除容器
// func removeContainer(cli *client.Client, containerID string) error {
// 	ctx := context.Background()

// 	err := cli.ContainerRemove(ctx, containerID, container.RemoveOptions{
// 		RemoveVolumes: true, // 刪除容器的掛載卷
// 		Force:         true, // 強制刪除運行中的容器
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to remove container %s: %v", containerID, err)
// 	}

// 	log.Printf("Container %s removed successfully", containerID)
// 	return nil
// }

// 創建並啟動新容器
func createAndStartContainer(cli *client.Client, containerImage, containerName string, strLabelType string) error {
	ctx := context.Background()
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: containerImage,
	}, &container.HostConfig{
		SecurityOpt: []string{
			fmt.Sprintf("label:type:%s", strLabelType),
		},
	}, &network.NetworkingConfig{}, nil, containerName)
	if err != nil {
		return fmt.Errorf("failed to create container: %v", err)
	}

	// 啟動容器，移除 ContainerStartOptions
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %v", err)
	}
	log.Printf("Container %s started", containerName)
	return nil
}
