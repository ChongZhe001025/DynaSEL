package monitor

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

var isPaused = false

func MonitorConfigJson(strConfigJsonPath string) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Error creating watcher: %v", err)
	}
	defer watcher.Close()

	// 要監控的路徑
	if _, err := os.Stat(strConfigJsonPath); os.IsNotExist(err) {
		log.Fatalf("Directory does not exist: %s", strConfigJsonPath)
	}

	// 啟動 Goroutine 處理事件
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// 當偵測到事件且監控未暫停時
				if !isPaused {
					fmt.Printf("Event detected: %s\n", event)

					// 暫停監控
					isPaused = true

					// 關閉 watcher
					err := watcher.Remove(strConfigJsonPath)
					if err != nil {
						log.Printf("Error removing watcher: %v", err)
					}

					// 執行其他程式
					fmt.Println("Executing task...")
					executeTask() // 執行自訂任務

					// 恢復監控
					fmt.Println("Resuming monitoring...")
					err = watcher.Add(strConfigJsonPath)
					if err != nil {
						log.Printf("Error adding watcher: %v", err)
					}
					isPaused = false
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				// 處理錯誤
				log.Printf("Error: %v\n", err)
			}
		}
	}()

	// 添加監控目錄
	err = watcher.Add(strConfigJsonPath)
	if err != nil {
		log.Fatalf("Error adding path to watcher: %v", err)
	}

	fmt.Printf("Monitoring changes in: %s\n", strConfigJsonPath)

	// 保持程式運行
	select {}

}

// executeTask 模擬執行其他程式
func executeTask() {
	// 模擬執行耗時任務
	fmt.Println("Task in progress...")
	time.Sleep(5 * time.Second) // 模擬5秒任務
	fmt.Println("Task completed!")
}

// func createWatcher(path string) error {
// 	// 創建文件系統監控器
// 	watcher, err := fsnotify.NewWatcher()
// 	if err != nil {
// 		return fmt.Errorf("failed to create watcher: %w", err)
// 	}
// 	defer watcher.Close()

// 	bolChangeFile := false

// 	bolChangeFile = startWatcher(path, watcher, bolChangeFile)

// 	if bolChangeFile {
// 		fmt.Println("detect change!")
// 	}

// 	return nil
// }

// func startWatcher(path string, watcher *fsnotify.Watcher, bolChangeFile bool) bool {
// 	err := watcher.Add(path)
// 	if err != nil {
// 		return false
// 	}
// 	done := make(chan struct{})

// 	go func() {
// 		defer close(done) // 確保 Goroutine 退出時關閉通道
// 		for {
// 			select {
// 			case _, ok := <-watcher.Events:
// 				if !ok {
// 					log.Println("Events channel closed")
// 					return
// 				}
// 				bolChangeFile = true
// 				watcher.Close() // 關閉監控器
// 				return          // 退出 Goroutine
// 			case err, ok := <-watcher.Errors:
// 				if !ok {
// 					log.Println("Errors channel closed")
// 					return
// 				}
// 				log.Printf("Watcher error: %v", err)
// 			}
// 		}
// 	}()
// 	close(done)
// 	return bolChangeFile
// }

// // if _, err := os.Stat(strConfigJsonPath); os.IsNotExist(err) {
// // 	log.Fatalf("File not found: %s", strConfigJsonPath)
// // }

// // // 初始化監控器
// // watcher, err := fsnotify.NewWatcher()
// // if err != nil {
// // 	log.Fatalf("Error initializing watcher: %v", err)
// // }
// // defer watcher.Close()

// // // 啟動 Goroutine 處理事件
// // go func() {
// // 	for {
// // 		select {
// // 		case event, ok := <-watcher.Events:
// // 			if !ok {
// // 				fmt.Println("fail")
// // 				return
// // 			}
// // 			fmt.Println("event.Name:" + event.Name)
// // 			fmt.Println("event.Name:" + event.Op.String())

// // 			// 檢查變更類型
// // 			if event.Op&fsnotify.Write == fsnotify.Write {
// // 				fmt.Printf("File modified: %s\n", event.Name)
// // 				// 這裡可以新增重新讀取檔案或其他操作
// // 			}
// // 		case err, ok := <-watcher.Errors:
// // 			fmt.Println(err)
// // 			if !ok {
// // 				return
// // 			}
// // 			log.Printf("Error: %v", err)
// // 		}
// // 	}
// // }()

// // // 將檔案加入監控清單
// // err = watcher.Add(strConfigJsonPath)
// // if err != nil {
// // 	log.Fatalf("Error adding file to watcher: %v", err)
// // }

// // // 阻塞主執行緒
// // select {}
