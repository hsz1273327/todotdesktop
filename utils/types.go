package utils

type DesktopEntry struct {
	Version        string
	Name           string
	Encoding       string
	Comment        string
	Exec           string
	Icon           string
	Terminal       bool
	Type           string
	Categories     string
	StartupWMClass string
	StartupNotify  bool
}

// Options 包含了创建快捷方式所需的所有用户输入
type Options struct {
	ExecPath     string
	AppName      string
	Comment      string
	IconFilePath string
	Version      string
	WMClass      string
	ExecType     string
	Terminal     string // 'true', 'false', or 'auto'
}
