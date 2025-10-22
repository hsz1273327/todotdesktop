package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hsz1273327/todotdesktop/utils"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] -filename <path> <appname>\n", filepath.Base(os.Args[0]))
	flag.PrintDefaults()
}

// createDesktopWorkflow 协调创建 .desktop 文件的整个流程
func createDesktopWorkflow(opts utils.Options) error {
	log.Println("开始创建快捷方式...")

	// 1. 给文件赋可执行权限
	log.Printf("设置可执行权限: %s", opts.ExecPath)
	err := utils.SetExecutablePermission(opts.ExecPath)
	if err != nil {
		return fmt.Errorf("设置可执行权限失败: %w", err)
	}
	// 2. 准备 .desktop 文件的内容(包含对icon的处理)
	log.Println("准备 .desktop 文件数据...")
	entryData, err := utils.GenerateDotDesktopEntry(opts)
	if err != nil {
		return fmt.Errorf("生成 .desktop 文件的EntryData数据失败: %w", err)
	}

	// 3. 将 .desktop 文件写入目标位置
	log.Println("写入 .desktop 文件...")
	content, err := utils.GenerateDotDesktopContent(*entryData)
	if err != nil {
		return fmt.Errorf("生成 .desktop 文件内容失败: %w", err)
	}
	utils.WriteDotDesktopFile(opts.AppName, content)
	log.Println("快捷方式创建成功!")
	// 4. 处理 WMClass 的特殊逻辑
	// 提醒用户更新 WMClass
	if opts.WMClass == "" && utils.IsX11() {
		log.Println("如果你希望快捷方式正确关联到应用窗口,请手动设置合适的 WMClass 值.")
		log.Println("在x11环境中,你可以使用 'xprop' 命令来获取应用窗口的 WM_CLASS 属性,步骤如下:")
		log.Println("1. 启动应用.")
		log.Println("2. 成功启动应用后,在终端中运行: xprop WM_CLASS")
		log.Println("3. 鼠标指针会变成一个十字形,点击应用窗口.")
		log.Println("然后点击应用窗口,你将看到类似下面的输出:")
		log.Println("WM_CLASS(STRING) = <your-wmclass>, <YourAppName>")
		log.Println("记下第一个值 (your-wmclass),并使用它来更新快捷方式.")
		log.Println("在获取到WMClass值后使用更新命令更新快捷方式,例如:")
		log.Printf("  %s -exec %s -wmclass <WMClass值> %s\n", filepath.Base(os.Args[0]), opts.ExecPath, opts.AppName)
	}

	return nil
}
func updateDesktopWorkflow(opts utils.Options) error {
	log.Println("开始更新快捷方式...")
	// 1. 给文件赋可执行权限
	log.Printf("设置可执行权限: %s", opts.ExecPath)
	err := utils.SetExecutablePermission(opts.ExecPath)
	if err != nil {
		return fmt.Errorf("设置可执行权限失败: %w", err)
	}

	// 2. 准备 .desktop 文件的内容(包含对icon的处理)
	log.Println("准备 .desktop 文件数据...")
	entryData, err := utils.GenerateDotDesktopEntry(opts)
	if err != nil {
		return fmt.Errorf("生成 .desktop 文件的EntryData数据失败: %w", err)
	}

	// 3. 读取现有的 .desktop 文件内容构造 DesktopEntry 结构体
	log.Println("读取现有的 .desktop 文件内容...")
	existingContent, err := utils.ReadDotDesktopFile(opts.AppName)
	if err != nil {
		return fmt.Errorf("读取现有的 .desktop 文件内容失败: %w", err)
	}

	// 4. 合并新旧数据，优先使用新数据
	log.Println("合并新旧数据...")
	newcontent, err := utils.UpdateDotDesktopContent(existingContent, *entryData)
	if err != nil {
		return fmt.Errorf("合并新旧数据失败: %w", err)
	}

	// 5. 将更新后的 .desktop 文件写入目标位置
	log.Println("写入更新后的 .desktop 文件...")
	err = utils.WriteDotDesktopFile(opts.AppName, newcontent)
	if err != nil {
		return fmt.Errorf("写入更新后的 .desktop 文件失败: %w", err)
	}
	log.Println("快捷方式更新成功!")
	return nil
}

func main() {
	// 定义命令行标志
	execPath := flag.String("exec", "", "应用的执行文件路径 (必选)")
	comment := flag.String("comment", "", "描述 (可选)")
	iconFilePath := flag.String("icon", "", "图标路径 (可选)")
	version := flag.String("version", "", "应用版本 (可选, 否则会尝试从文件名解析)")
	execType := flag.String("type", "", "应用类型 (e.g., 'script', 'electron')")
	terminal := flag.String("terminal", "auto", "是否在终端中运行 ('true', 'false', 'auto')")
	wmclass := flag.String("wmclass", "", "窗口的 WM_CLASS 值 (可选)")

	flag.Usage = usage
	flag.Parse()

	// 校验参数
	if *execPath == "" {
		fmt.Fprintln(os.Stderr, "错误: -exec 参数是必选的")
		usage()
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "错误: 需要一个应用名")
		usage()
		os.Exit(1)
	}
	appName := args[0]

	// 填充选项结构体
	opts := utils.Options{
		ExecPath:     *execPath,
		AppName:      appName,
		Comment:      *comment,
		IconFilePath: *iconFilePath,
		Version:      *version,
		WMClass:      *wmclass,
		ExecType:     *execType,
		Terminal:     *terminal,
	}

	DotDesktopExistok, err := utils.DotDesktopExist(appName)
	if err != nil {
		log.Fatalf("检查 .desktop 文件存在性时出错: %v", err)
	}
	if DotDesktopExistok {
		log.Printf("[警告] 名为 '%s' 的 .desktop 文件已存在,操作将更新该文件", appName)
		// 执行更新工作流
		err := updateDesktopWorkflow(opts)
		if err != nil {
			log.Fatalf("错误: %v", err)
		}
	} else {
		log.Printf("[警告] 名为 '%s' 的 .desktop 文件不存在,将创建新文件", appName)
		// 执行创建工作流
		err := createDesktopWorkflow(opts)
		if err != nil {
			log.Fatalf("错误: %v", err)
		}
	}

}
