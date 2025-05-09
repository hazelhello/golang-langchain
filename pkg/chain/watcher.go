package chain

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
	"time"
)

// StartWatching 启动模板目录监控（新增文件）
func (c *Chain) StartWatching(templateDir string) error {
	watcher, err := fsnotify.NewWatcher() // 创建文件监控器
	if err != nil {
		return fmt.Errorf("创建监控器失败: %w", err)
	}

	go func() {
		defer watcher.Close()
		var debounceTimer *time.Timer // 防抖计时器
		
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// 只处理写入事件
				if event.Op&fsnotify.Write == fsnotify.Write {
					// 防抖处理：500ms内多次变更只触发一次
					if debounceTimer != nil {
						debounceTimer.Stop()
					}
					debounceTimer = time.AfterFunc(500*time.Millisecond, func() {
						content, err := os.ReadFile(event.Name)
						if err != nil {
							log.Printf("⚠️ 读取模板文件失败: %v", err)
							return
						}
						tplName := filepath.Base(event.Name)
						if err := c.AddTemplate(tplName, string(content)); err != nil {
							log.Printf("⚠️ 模板更新失败: %v", err)
						} else {
							log.Printf("🔄 模板已热重载: %s", tplName)
						}
					})
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("监控错误: %v", err)
			}
		}
	}()

	// 添加监控目录
	return watcher.Add(templateDir)
}