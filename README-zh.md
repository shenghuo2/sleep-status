# Sleep Status - Magisk 模块

![version](https://img.shields.io/github/v/release/shenghuo2/sleep-status?include_prereleases&label=version)

[English](./README.md) | [简体中文](./README-zh.md)

这是 Sleep Status 项目的 Magisk 模块。它通过向服务器发送心跳信号，在网络连接丢失时自动报告睡眠状态。

## 功能特点

- 通过心跳信号自动检测网络连接状态
- 在网络不可用时报告睡眠状态
- 支持心跳监控
- 通过 Magisk 作为系统服务运行
- 支持 Android ARM64 设备

## 安装方法

1. 从 [Releases](https://github.com/shenghuo2/sleep-status/releases) 下载最新的 `sleep-status-magisk.zip`
2. 通过 Magisk Manager 安装
3. 编辑 `/data/adb/modules/sleep-status/config.conf` 设置你的服务器 URL 和 API 密钥
4. 重启设备以使更改生效

## 配置说明

配置文件（`config.conf`）包含：
```
SERVER_URL="http://your-server:8000"  # 你的服务器 URL
API_KEY="your-key"                    # 用于认证的 API 密钥
```

修改配置后，必须重启设备才能使更改生效。

## 日志查看

你可以在 Magisk 日志中查看模块的运行情况：
```bash
tail -f /cache/magisk.log | grep sleep-status
```

## 故障排除

1. 如果模块没有在 Magisk Manager 中显示：
   - 确保你使用的是最新版本的 Magisk
   - 检查模块是否正确安装在 `/data/adb/modules/` 目录下

2. 如果状态没有被报告：
   - 检查配置文件是否存在且配置正确
   - 验证服务器 URL 是否可以从设备访问
   - 检查 Magisk 日志中是否有错误信息
   - 确保在修改配置后已经重启设备

## 从源码构建

可以使用提供的构建脚本构建模块：
```bash
./build.sh
```

这将创建一个可以通过 Magisk Manager 安装的刷机包。

## 贡献代码

欢迎提交 issues 和 pull requests 来改进模块。

## 开源协议

本项目采用 MIT 协议 - 详见 LICENSE 文件
