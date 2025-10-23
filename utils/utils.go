package utils

import (
	"errors"
	"path/filepath"
)

func GenerateDotDesktopEntry(opts Options) (*DesktopEntry, error) {
	// 处理 WMClass 的默认值
	if opts.WMClass == "" && !IsX11() {
		opts.WMClass = opts.AppName
	}
	// 处理版本号的默认值
	version := opts.Version
	if version == "" {
		appfilename := filepath.Base(opts.ExecPath)
		_version := ParseVersionFromName(appfilename)
		if _version != "" {
			version = _version
		}
	}
	//处理描述信息的默认值
	comment := opts.Comment
	if comment == "" {
		if opts.ExecType == "script" {
			comment = opts.AppName + " 脚本快捷方式"
		} else {
			comment = opts.AppName + " 应用快捷方式"
		}
	}
	// 处理图标
	icon := ""
	if opts.IconFilePath != "" {
		iconsuffix := filepath.Ext(opts.IconFilePath)
		if iconsuffix != ".svg" && iconsuffix != ".png" {
			return nil, errors.New("不支持的图标格式，仅支持 .png 和 .svg")
		}
		iconExitOk, err := IconExist(opts.AppName, iconsuffix)
		if err != nil {
			return nil, err
		}
		if iconExitOk {
			err = DeleteIcon(opts.AppName, iconsuffix)
			if err != nil {
				return nil, err
			}
			icon, err = CopyIcon(opts.AppName, opts.IconFilePath)
			if err != nil {
				return nil, err
			}
		} else {
			icon, err = CopyIcon(opts.AppName, opts.IconFilePath)
			if err != nil {
				return nil, err
			}
		}
	} else {
		pngIconExitOk, err := IconExist(opts.AppName, ".png")
		if err != nil {
			return nil, err
		}
		if pngIconExitOk {
			icon = opts.AppName + ".png"
		} else {
			svgIconExitOk, err := IconExist(opts.AppName, ".svg")
			if err != nil {
				return nil, err
			}
			if svgIconExitOk {
				icon = opts.AppName + ".svg"
			}
		}
	}
	switch opts.ExecType {
	case "script":
		{
			terminal := true
			if opts.Terminal == "false" {
				terminal = false
			}
			return &DesktopEntry{
				Version:        version,
				Name:           opts.AppName,
				Encoding:       "UTF-8",
				Comment:        comment,
				Exec:           opts.ExecPath,
				Icon:           icon,
				Terminal:       terminal,
				Type:           "Application",
				Categories:     "Utility;",
				StartupWMClass: opts.WMClass,
			}, nil
		}
	case "chromium":
		{
			terminal := false
			if opts.Terminal == "true" {
				terminal = true
			}
			return &DesktopEntry{
				Version:        opts.Version,
				Name:           opts.AppName,
				Encoding:       "UTF-8",
				Comment:        comment,
				Exec:           opts.ExecPath + " --no-sandbox" + " %U",
				Icon:           icon,
				Terminal:       terminal,
				Type:           "Application",
				Categories:     "Application;",
				StartupWMClass: opts.WMClass,
			}, nil
		}
	default:
		{
			terminal := false
			if opts.Terminal == "true" {
				terminal = true
			}
			return &DesktopEntry{
				Version:        version,
				Name:           opts.AppName,
				Encoding:       "UTF-8",
				Comment:        comment,
				Exec:           opts.ExecPath + " %U",
				Icon:           icon,
				Terminal:       terminal,
				Type:           "Application",
				Categories:     "Application;",
				StartupWMClass: opts.WMClass,
			}, nil
		}
	}
}
