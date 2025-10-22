package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// FileExists 检查指定路径的文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// // SanitizeFileName 清理文件名以符合规范
// func SanitizeFileName(name string) string {
// 	name = strings.ToLower(name)
// 	name = strings.ReplaceAll(name, " ", "-")
// 	name = strings.ReplaceAll(name, "_", "-")
// 	// 移除所有非字母、数字、连字符的字符
// 	re := regexp.MustCompile(`[^a-z0-9-]`)
// 	res := re.ReplaceAllString(name, "")
// 	if res == "" {
// 		return "application"
// 	}
// 	return res
// }

func dotDesktopFilePath(appname string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dotdesktopdir := filepath.Join(homeDir, ".local", "share", "applications")
	fileName := appname + ".desktop"
	outPath := filepath.Join(dotdesktopdir, fileName)
	return outPath, nil
}

func iconFilePath(appname string, suffix string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	icondir := filepath.Join(homeDir, ".local", "share", "icons")
	to := filepath.Join(icondir, appname+suffix)
	return to, nil
}

// DotDesktopExist 检查applications目录下的 .desktop 文件是否存在
func DotDesktopExist(appname string) (bool, error) {
	outPath, err := dotDesktopFilePath(appname)
	if err != nil {
		return false, err
	}
	return FileExists(outPath), nil
}

// IconExist 检查icon目录下的 .desktop 文件是否存在
func IconExist(appname string, suffix string) (bool, error) {
	if suffix != ".svg" {
		suffix = ".png"
	}
	outPath, err := iconFilePath(appname, suffix)
	if err != nil {
		return false, err
	}

	return FileExists(outPath), nil
}

// CopyIcon 复制图标文件到用户图标目录
func CopyIcon(appname string, from string) (string, error) {
	suffix := strings.ToLower(filepath.Ext(from))
	if suffix != ".svg" && suffix != ".png" {
		return "", errors.New("不支持的图标格式，仅支持 .png 和 .svg")
	}
	outPath, err := iconFilePath(appname, suffix)
	if err != nil {
		return "", err
	}

	input, err := os.ReadFile(from)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(outPath, input, 0644)
	if err != nil {
		return "", err
	}
	return appname + suffix, nil
}

// DeleteIcon 删除用户图标目录下的图标文件
func DeleteIcon(appname string, suffix string) error {
	if suffix != ".svg" && suffix != ".png" {
		suffix = ".png"
	}
	outPath, err := iconFilePath(appname, suffix)
	if err != nil {
		return err
	}
	err = os.Remove(outPath)
	if err != nil {
		return err
	}
	return nil
}

func ReadDotDesktopFile(appname string) (string, error) {
	outPath, err := dotDesktopFilePath(appname)
	if err != nil {
		return "", err
	}
	content, err := os.ReadFile(outPath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// WriteDotDesktopFile 将 .desktop 文件内容写入指定目录
func WriteDotDesktopFile(appname string, content string) error {
	outPath, err := dotDesktopFilePath(appname)
	if err != nil {
		return err
	}
	err = os.WriteFile(outPath, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}

// SetExecutablePermission 给指定路径的文件添加可执行权限
func SetExecutablePermission(path string) error {
	return os.Chmod(path, 0755)
}
