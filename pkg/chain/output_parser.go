package chain

import (
	"encoding/json"
	"strings"
)

// QAOutput 结构化输出格式（新增结构体）
type QAOutput struct {
	Answer   string   `json:"answer"`   // 回答内容
	Keywords []string `json:"keywords"` // 关键词列表
	Confidence float64 `json:"confidence,omitempty"` // 置信度（可选）
}

// ParseJSONOutput 解析LLM的JSON输出（新增函数）
func ParseJSONOutput(output string) (*QAOutput, error) {
	var result QAOutput
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w", err)
	}
	return &result, nil
}

// ExtractKeywords 从文本提取关键词（备用方案）
func ExtractKeywords(text string) []string {
	// 简单实现：按非字母字符分割并去重
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !(r >= 'a' && r <= 'z') && !(r >= 'A' && r <= 'Z')
	})
	
	unique := make(map[string]struct{})
	for _, w := range words {
		if len(w) > 3 { // 过滤短词
			unique[strings.ToLower(w)] = struct{}{}
		}
	}
	
	keywords := make([]string, 0, len(unique))
	for k := range unique {
		keywords = append(keywords, k)
	}
	return keywords
}