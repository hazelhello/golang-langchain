package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	
	"github.com/hazel/golang-langchain/pkg/chain" // 替换为实际模块路径
)

func main() {
	// 1. 初始化链式处理器
	myChain := chain.NewChain()
	
	// 2. 加载初始模板
	initialTpl := `你是一个{{.Expert}}，请用{{.Level}}级别解释：{{.Question}}`
	if err := myChain.AddTemplate("explain", initialTpl); err != nil {
		log.Fatal("初始模板加载失败: ", err)
	}
	
	// 3. 启动模板热加载（新增功能调用）
	tplDir := "./templates"
	if err := os.MkdirAll(tplDir, 0755); err != nil {
		log.Printf("⚠️ 创建模板目录失败: %v", err)
	} else if err := myChain.StartWatching(tplDir); err != nil {
		log.Printf("⚠️ 模板监控启动失败: %v", err)
	} else {
		log.Printf("👀 正在监控模板目录: %s", tplDir)
	}
	
	// 4. 示例调用
	inputs := map[string]interface{}{
		"Expert":   "量子物理学家",
		"Level":    "科普级",
		"Question": "量子隧穿效应",
	}
	
	// 执行链式调用
	rawOutput, err := myChain.Run("explain", inputs)
	if err != nil {
		log.Fatal("执行失败: ", err)
	}
	
	// 尝试结构化解析（新增处理逻辑）
	if parsed, err := chain.ParseJSONOutput(rawOutput); err == nil {
		fmt.Printf("【结构化输出】\n回答: %s\n关键词: %v\n", 
			parsed.Answer, parsed.Keywords)
	} else {
		fmt.Printf("【原始输出】\n%s\n", rawOutput)
	}
}