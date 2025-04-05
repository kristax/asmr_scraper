package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/go-resty/resty/v2"
)

// Translate 将输入字符串翻译成中文，若输入已是中文则直接返回
func Translate(s string) (string, error) {
	if isChinese(s) {
		return s, nil
	}

	apiKey := os.Getenv("GOOGLE_TRANSLATE_API_KEY")
	if apiKey == "" {
		return "", errors.New("GOOGLE_TRANSLATE_API_KEY 环境变量未设置")
	}

	client := resty.New()
	resp, err := client.R().
		SetQueryParam("key", apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"q":      []string{s},
			"target": "zh-CN",
		}).
		Post("https://translation.googleapis.com/language/translate/v2")

	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("API 错误: 状态码 %d, 响应: %s", resp.StatusCode(), resp.String())
	}

	var result struct {
		Data struct {
			Translations []struct {
				TranslatedText string `json:"translatedText"`
			} `json:"translations"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if len(result.Data.Translations) == 0 {
		return "", errors.New("未找到翻译结果")
	}

	return result.Data.Translations[0].TranslatedText, nil
}

// isChinese 检测字符串是否为中文字符串（包含中文标点和数字）
func isChinese(s string) bool {
	// 允许中文字符、中文标点、数字和空格
	allowedPattern := `^[\p{Han}\d\s，。！？、；：“”‘’（）《》【】…·—]*$`
	if matched, _ := regexp.MatchString(allowedPattern, s); !matched {
		return false
	}

	// 必须包含至少一个中文字符
	hanPattern := `\p{Han}`
	matched, _ := regexp.MatchString(hanPattern, s)
	return matched
}
