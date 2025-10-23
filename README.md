# todotdesktop

用于给linux环境下的脚本或appimage创建快捷方式

## 简介

对于没有安装程序的应用程序（例如AppImage）或脚本文件，您可以使用此工具为它们创建桌面快捷方式. 已经存在的同名快捷方式会被覆盖.不存在的则会被创建.

## 安装

本项目是一个go语言编写的命令行工具,你可以通过以下两种方式安装它:

### 预编译版本

请在本项目的[发布页面](https://github.com/hsz1273327/todotdesktop/releases)下载最新的预编译二进制文件.下载后,请将其解压,赋予可执行权限并移动到您的系统PATH中的某个目录,以`/usr/local/bin`为例:

```bash
chmod +x todotdesktop
mv todotdesktop /usr/local/bin/
```

### 编译安装

要编译安装此工具，您需要确保已安装Go语言环境.然后,您可以使用以下命令进行安装:

```bash
go install github.com/LinuxTools/todotdesktop@latest
```

## 使用方法

使用`todotdesktop`命令需要指定最后一位为应用的名字.除此之外必须指定`-exec`来指定应用可执行文件的路径.其他参数均为可选.

```bash
todotdesktop -exec /path/to/your/appimage_or_script AppName
```

### 补充说明

> 关于快捷方式存放位置

本工具会将创建的快捷方式存放在`~/.local/share/applications/`目录下,这是大多数linux发行版默认的用户级应用快捷方式存放位置.
如果你想让所有用户都能使用这个快捷方式,你可以手动将生成的`.desktop`文件复制到`/usr/share/applications/`目录下,并且需要指定对应的可执行文件给所有用户都有可执行权限.当然这个操作我认为是安全风险较高的,请谨慎操作,因此本工具并不支持直接创建到该目录下.

> 关于执行类型

本工具设计的目的是为了支持appimage和脚本两种类型的应用.而appimage中又根据是否基于chromium内核分为chromium内核和非chromium内核两种类型.因此,本工具支持三种执行类型: `default`, `chromium`和`script`,使用`-type`指定.这三种给的快捷方式配置会有些不同,默认执行类型为`default`也就是非chromium内核的appimage.

```bash
todotdesktop -type chromium -exec /path/to/your/appimage_or_script AppName
```

在执行本工具前建议先确认应用的类型,以便指定正确的执行类型.脚本好说,appimage的话建议先直接运行一下,如果是基于chromium内核的appimage,一般会无法运行.

先执行一下先确保可以运行是个好习惯,毕竟有些appimage需要额外的设置或者一些其他依赖才能运行.

> 关于版本号

鉴于很多appimage本身带版本号信息,对于这类应用我们可以不用指定`-version`参数,但如果是脚本文件,建议指定版本号.

```bash
todotdesktop -version 0.0.1 -exec /path/to/your/script AppName
```

> 关于图标

本工具仅支持`png`和`svg`,appimage下载下来一般不带图标,获取图标的方法一般有两种:

+ 从appimage中提取图标文件

    ```bash
    chmod +x YourAppImageFile.AppImage
    ./YourAppImageFile.AppImage --appimage-extract
    ```

    你会得到一个`squashfs-root`目录,你可以去里面找符合要求的图标文件.等到创建好桌面快捷方式后,可以删除这个目录.

+ 去网上下载图标文件

在获取到图标文件后,可以通过`-icon`参数指定图标文件路径,工具会将其复制到`~/.local/share/icons/`目录下并以你指定的应用名称命名.

```bash
todotdesktop -icon /path/to/icon -exec /path/to/your/script AppName
```

另外可以顺便观察下应用能不能在dock中正确显示图标,如果不能我们就还需要指定`-wmclass`参数. 这个参数在基于chromium内核的appimage中比较常见需要额外指定.我们可以通过`xprop`命令来获取应用的WM_CLASS值.

通常在x11环境下可以通过如下步骤获取WM_CLASS值:

1. 启动应用.
2. 成功启动应用后,打开一个终端,在终端中运行: `xprop WM_CLASS`
3. 鼠标指针会变成一个十字形,点击应用窗口.
4. 后点击应用窗口,你将看到类似下面的输出:

    ```bash
    M_CLASS(STRING) = <your-wmclass>, <YourAppName>
    ```

    第一个值`<your-wmclass>`就是我们要的

然后我们可以通过如下命令创建快捷方式:

```bash
todotdesktop -type chromium -wmclass <your-wmclass> -exec /path/to/your/appimage AppName
```