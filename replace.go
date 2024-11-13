package kakaxi

import (
	"path"
	"regexp"
	"strings"
)

func ReplaceHTML(body []byte, currentPath string, currentHost string) (bbody []byte, err error) {
	var (
		prefixHttps = "https://" + currentHost
		prefixHttp  = "http://" + currentHost
	)

	mapping := func(target string) string {
		var (
			isHttps bool
		)

		target = strings.ReplaceAll(target, "&#x2F;", "/")

		isHttps = strings.HasPrefix(target, "http:") || strings.HasPrefix(target, "https:")
		if isHttps {
			if strings.HasPrefix(target, prefixHttps) {
				target = strings.TrimPrefix(target, prefixHttps)
			} else if strings.HasPrefix(target, prefixHttp) {
				target = strings.TrimPrefix(target, prefixHttp)
			} else {
				//otherDomain
				return target
			}
		}
		return getRelativePath(currentPath, target)
	}

	// 执行替换
	newHTML := ExpandURL(Bytes2string(body), mapping)
	return String2Bytes(newHTML), nil
}

func getRelativePath(currentPath, targetPath string) string {
	// 标准化路径，使用 filepath.Clean 来处理所有路径
	currentPath = path.Clean(currentPath)
	targetPath = path.Clean(targetPath)

	// 检查目标路径是否是相对路径
	if !path.IsAbs(targetPath) {
		// 如果是相对路径，则基于当前路径构建完整路径
		return path.Join(currentPath, targetPath)
	}

	// 特殊处理根目录
	if currentPath == "/" {
		// 如果当前路径是根目录，直接去掉目标路径的根斜杠
		return strings.TrimPrefix(targetPath, "/")
	}

	// 分割路径成切片，过滤掉空字符串
	currentParts := splitAndFilter(currentPath)
	targetParts := splitAndFilter(targetPath)

	// 找到共同前缀的长度
	i := 0
	for i < len(currentParts) && i < len(targetParts) && currentParts[i] == targetParts[i] {
		i++
	}

	// 计算需要返回上层目录的次数
	backCount := len(currentParts) - i

	// 构建结果
	var result []string

	// 添加返回上层的 ..
	for j := 0; j < backCount; j++ {
		result = append(result, "..")
	}

	// 添加目标路径的剩余部分
	result = append(result, targetParts[i:]...)

	// 如果结果为空，返回"."表示当前目录
	if len(result) == 0 {
		return "."
	}

	// 使用路径分隔符连接
	return path.Join(result...)
}

// splitAndFilter 分割路径并过滤空字符串
func splitAndFilter(path string) []string {
	parts := strings.Split(path, "/")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

type URLInfo struct {
	Tag       string // HTML标签
	Attr      string // 属性名
	URL       string // URL值
	Position  int    // 在文本中的位置
	TagIndex  [2]int
	AttrIndex [2]int
	URLIndex  [2]int
}

func ExtractURLInfo(html string) []URLInfo {
	// 更详细的正则表达式，捕获标签和属性名
	detailPattern := `(?i)<((?:a|area|base|link|img|script|iframe|embed|source|track|input|frame|form|object|video|audio|meta))\s+[^>]*?(href|src|action|data|poster|srcset|content)\s*=\s*(["'])(.*?)(?:["'])`

	re := regexp.MustCompile(detailPattern)
	matches := re.FindAllStringSubmatchIndex(html, -1)

	results := make([]URLInfo, 0, len(matches))
	for _, match := range matches {
		if len(match) >= 10 { // 确保有足够的捕获组
			info := URLInfo{
				Tag:       html[match[2]:match[3]],
				TagIndex:  [2]int{match[2], match[3]},
				Attr:      html[match[4]:match[5]],
				AttrIndex: [2]int{match[4], match[5]},
				URL:       html[match[8]:match[9]],
				URLIndex:  [2]int{match[8], match[9]},
				Position:  match[0],
			}
			results = append(results, info)
		}
	}
	return results
}

func ExpandURL(html string, mapping func(string) string) (result string) {
	infos := ExtractURLInfo(html)
	indexOffset := 0
	result = html
	for _, info := range infos {
		old := info.URL
		newURL := mapping(info.URL)
		result = result[:info.URLIndex[0]+indexOffset] + newURL + result[info.URLIndex[1]+indexOffset:]
		// 更新偏移量
		indexOffset += len(newURL) - len(old)
	}
	return result
}
