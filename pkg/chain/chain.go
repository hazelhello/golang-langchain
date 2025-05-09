package chain

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

// CacheProvider 定义缓存接口（需补充）
type CacheProvider interface {
	Get(key string) (string, bool)
	Set(key, value string)
}

// LLMClient 定义LLM客户端接口（需补充）
type LLMClient interface {
	Call(prompt string) (string, error)
}

// Chain 核心链式处理结构体
type Chain struct {
	templates    map[string]*template.Template // 存储模板的映射
	cache        CacheProvider                 // 缓存接口
	llmClient    LLMClient                     // LLM客户端接口
}

// NewChain 初始化链式处理器
func NewChain() *Chain {
	return &Chain{
		templates: make(map[string]*template.Template), // 初始化模板映射
		cache:     NewRedisCache("localhost:6379"),     // 默认使用Redis缓存
		llmClient: NewOpenAIClient(os.Getenv("OPENAI_KEY")), // 从环境变量获取OpenAI Key
	}
}

// AddTemplate 添加命名模板
func (c *Chain) AddTemplate(name, tplText string) error {
	tpl, err := template.New(name).Parse(tplText) // 解析模板文本
	if err != nil {
		return fmt.Errorf("模板解析失败: %w", err)
	}
	c.templates[name] = tpl // 存入映射表
	return nil
}

// Run 执行链式调用流程
func (c *Chain) Run(templateName string, inputs map[string]interface{}) (string, error) {
	// 1. 模板渲染
	var prompt strings.Builder
	if tpl, exists := c.templates[templateName]; exists {
		if err := tpl.Execute(&prompt, inputs); err != nil { // 注入变量生成最终prompt
			return "", fmt.Errorf("模板渲染失败: %w", err)
		}
	} else {
		return "", fmt.Errorf("模板未找到: %s", templateName)
	}

	// 2. 缓存检查
	if cached, hit := c.cache.Get(prompt.String()); hit {
		return cached, nil // 命中缓存直接返回
	}

	// 3. 调用LLM
	response, err := c.llmClient.Call(prompt.String())
	if err != nil {
		return "", fmt.Errorf("LLM调用失败: %w", err)
	}

	// 4. 缓存结果
	c.cache.Set(prompt.String(), response)
	return response, nil
}