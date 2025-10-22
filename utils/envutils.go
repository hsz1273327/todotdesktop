package utils

import (
	"os"
	"strings"
)

func IsX11() bool {
	// 检查 XDG_SESSION_TYPE 环境变量
	sessionType := os.Getenv("XDG_SESSION_TYPE")
	if strings.ToLower(sessionType) == "x11" {
		return true
	}

	// 作为一个备用检查，检查 DISPLAY 变量。
	// 在 Wayland 下运行 Xwayland 时，DISPLAY 也可能被设置，但 XDG_SESSION_TYPE 提供了更准确的信息。
	if sessionType == "" && os.Getenv("DISPLAY") != "" {
		// 如果没有设置 XDG_SESSION_TYPE，但设置了 DISPLAY，我们倾向于认为是 X11
		// （但在现代 Linux 系统中，这可能不完全可靠，因为 Wayland 也可以设置 DISPLAY 来支持 Xwayland）。
		// 最好依赖 XDG_SESSION_TYPE。
		// 在这里，我们只在 XDG_SESSION_TYPE 未设置时使用 DISPLAY 辅助判断。
		return true
	}

	return false
}
