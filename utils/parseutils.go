package utils

import (
	"bytes"
	"errors"
	"regexp"
	"text/template"
)

// ParseVersionFromName 从文件名中提取版本号
func ParseVersionFromName(filename string) string {
	re := regexp.MustCompile(`[vV]?(\d+\.\d+(\.\d+)?(-\w+)?)`)
	matches := re.FindStringSubmatch(filename)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// GenerateDotDesktopContent 生成 .desktop 文件的内容
func GenerateDotDesktopContent(entry DesktopEntry) (string, error) {

	tmpl := `[Desktop Entry]
{{if .Version}}Version={{.Version}}
{{end}}Name={{.Name}}
Encoding={{.Encoding}}
{{if .Comment}}Comment={{.Comment}}
{{end}}Exec={{.Exec}}
{{if .Icon}}Icon={{.Icon}}
{{end}}Terminal={{if .Terminal}}true{{else}}false{{end}}
Type={{.Type}}
Categories={{.Categories}}
{{if .StartupWMClass}}StartupWMClass={{.StartupWMClass}}
{{end}}StartupNotify={{if .StartupNotify}}true{{else}}false{{end}}
`
	t, err := template.New("desktop").Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, entry)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ParseDotDesktopContent 解析 .desktop 文件内容为 DesktopEntry 结构体
func ParseDotDesktopContent(content string) (DesktopEntry, error) {
	lines := bytes.Split([]byte(content), []byte("\n"))
	contentMap := make(map[string]string)
	for _, line := range lines {
		parts := bytes.SplitN(line, []byte("="), 2)
		if len(parts) == 2 {
			contentMap[string(parts[0])] = string(parts[1])
		}
	}
	// 构造 DesktopEntry 结构体
	entry := DesktopEntry{
		Version:        contentMap["Version"],
		Name:           contentMap["Name"],
		Comment:        contentMap["Comment"],
		Exec:           contentMap["Exec"],
		Icon:           contentMap["Icon"],
		Terminal:       contentMap["Terminal"] == "true",
		Type:           contentMap["Type"],
		Categories:     contentMap["Categories"],
		StartupWMClass: contentMap["StartupWMClass"],
		StartupNotify:  contentMap["StartupNotify"] == "true",
	}
	return entry, nil
}

// updateEntry 根据 updatedEntry 更新 originalEntry 的字段
func updateEntry(originalEntry DesktopEntry, updatedEntry DesktopEntry) (DesktopEntry, error) {
	result := DesktopEntry{}
	if updatedEntry.Name != "" && updatedEntry.Name != originalEntry.Name {
		return DesktopEntry{}, errors.New("原始应用名与指定名不一致,请手动检查")
	}
	result.Name = originalEntry.Name
	// 更新版本
	if updatedEntry.Version != "" && updatedEntry.Version != originalEntry.Version {
		result.Version = updatedEntry.Version
	} else {
		result.Version = originalEntry.Version
	}
	// 更新编码
	if updatedEntry.Encoding != "" && updatedEntry.Encoding != originalEntry.Encoding {
		result.Encoding = updatedEntry.Encoding
	} else {
		result.Encoding = originalEntry.Encoding
	}
	// 更新注释
	if updatedEntry.Comment != "" && updatedEntry.Comment != originalEntry.Comment {
		result.Comment = updatedEntry.Comment
	} else {
		result.Comment = originalEntry.Comment
	}
	// 更新执行命令
	if updatedEntry.Exec != "" && updatedEntry.Exec != originalEntry.Exec {
		result.Exec = updatedEntry.Exec
	} else {
		result.Exec = originalEntry.Exec
	}
	// 更新图标
	if updatedEntry.Icon != "" && updatedEntry.Icon != originalEntry.Icon {
		result.Icon = updatedEntry.Icon
	} else {
		result.Icon = originalEntry.Icon
	}
	// 更新终端选项
	if updatedEntry.Terminal != originalEntry.Terminal {
		result.Terminal = updatedEntry.Terminal
	} else {
		result.Terminal = originalEntry.Terminal
	}
	// 更新类型
	if updatedEntry.Type != "" && updatedEntry.Type != originalEntry.Type {
		result.Type = updatedEntry.Type
	} else {
		result.Type = originalEntry.Type
	}
	// 更新分类
	if updatedEntry.Categories != "" && updatedEntry.Categories != originalEntry.Categories {
		result.Categories = updatedEntry.Categories
	} else {
		result.Categories = originalEntry.Categories
	}
	// 更新启动WM类
	if updatedEntry.StartupWMClass != "" && updatedEntry.StartupWMClass != originalEntry.StartupWMClass {
		result.StartupWMClass = updatedEntry.StartupWMClass
	} else {
		result.StartupWMClass = originalEntry.StartupWMClass
	}
	// 更新启动通知
	if updatedEntry.StartupNotify != originalEntry.StartupNotify {
		result.StartupNotify = updatedEntry.StartupNotify
	} else {
		result.StartupNotify = originalEntry.StartupNotify
	}
	return result, nil
}

// UpdateDotDesktopContent 更新已有 .desktop 文件的内容
func UpdateDotDesktopContent(originalContent string, updatedEntry DesktopEntry) (string, error) {
	originalEntry, err := ParseDotDesktopContent(originalContent)
	if err != nil {
		return "", err
	}
	mergedEntry, err := updateEntry(originalEntry, updatedEntry)
	if err != nil {
		return "", err
	}
	return GenerateDotDesktopContent(mergedEntry)
}
