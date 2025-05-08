package chain

import (
	"strings"
	"text/template"
)

type Chain struct {
	templates    map[string]*template.Template
	cache        CacheProvider
	llmClient    LLMClient
}

// 初始化链
func NewChain() *Chain {
	return &Chain{
		templates: make(map[string]*template.Template),
		cache:     NewRedisCache("localhost:6379"),
		llmClient: NewOpenAIClient(os.Getenv("OPENAI_KEY")),
	}
}

// 添加模板
func (c *Chain) AddTemplate(name, tplText string) error {
	tpl, err := template.New(name).Parse(tplText)
	if err != nil {
		return err
	}
	c.templates[name] = tpl
	return nil
}

// 执行链式调用
func (c *Chain) Run(templateName string, inputs map[string]interface{}) (string, error) {
	// 1. 渲染Prompt
	var prompt strings.Builder
	if tpl, exists := c.templates[templateName]; exists {
		if err := tpl.Execute(&prompt, inputs); err != nil {
			return "", err
		}
	} else {
		return "", fmt.Errorf("template not found")
	}

	// 2. 检查缓存
	if cached, hit := c.cache.Get(prompt.String()); hit {
		return cached, nil
	}

	// 3. 调用LLM
	response, err := c.llmClient.Call(prompt.String())
	if err != nil {
		return "", err
	}

	// 4. 缓存结果
	c.cache.Set(prompt.String(), response)
	return response, nil
}
