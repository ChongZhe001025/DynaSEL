package automation

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func ApplyPolicyToContainer() {
	// 設定變數
	containerName := "my_container"
	newContainerName := "my_new_container"
	exportFile := "container_backup.tar"
	imageName := "my_new_image"

	// 1. 載入 SELinux .pp 模組
	if err := loadSELinuxPolicy("docker_fix.pp"); err != nil {
		log.Fatalf("SELinux policy loading failed: %v", err)
	}

	// 2. 停止容器
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf("Docker client initialization failed: %v", err)
	}
	if err := stopContainer(cli, containerName); err != nil {
		log.Fatalf("Failed to stop container: %v", err)
	}

	// 3. 導出容器數據
	if err := exportContainer(cli, containerName, exportFile); err != nil {
		log.Fatalf("Failed to export container: %v", err)
	}

	// 4. 重新導入容器數據並創建新的映像
	if err := importContainer(cli, exportFile, imageName); err != nil {
		log.Fatalf("Failed to import container: %v", err)
	}

	// 5. 創建並啟動新容器
	if err := createAndStartContainer(cli, imageName, newContainerName); err != nil {
		log.Fatalf("Failed to create and start new container: %v", err)
	}

	log.Println("New container created and started successfully!")
}

// 載入 SELinux Policy
func loadSELinuxPolicy(policyFile string) error {
	cmd := exec.Command("semodule", "-i", policyFile)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command failed: %s, %v", stderr.String(), err)
	}
	return nil
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
func importContainer(cli *client.Client, importFile, imageName string) error {
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
	log.Printf("Container data from %s imported as image %s", importFile, imageName)
	return nil
}

// 創建並啟動新容器
func createAndStartContainer(cli *client.Client, imageName, containerName string) error {
	ctx := context.Background()
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, &container.HostConfig{
		SecurityOpt: []string{
			"label:type:container_t",
			"label:level:s0:c123,c456",
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
