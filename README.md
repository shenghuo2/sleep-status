# Sleep-Status

![version](https://img.shields.io/github/v/release/shenghuo2/sleep-status?include_prereleases&label=version)

[English](./README-en.md)  [简体中文](./README.md)

## 项目介绍

Sleep-Status 是一个使用 Go 语言编写的简单后端服务。该服务通过读取配置文件 `config.json` 来获取和修改 `sleep` 状态。

灵感来源于 [这个BiliBili视频](https://www.bilibili.com/video/BV1fE421A7PE/)。

使用 Go 语言是因为其跨平台交叉编译特性，使得服务可以便捷的在多个操作系统上运行。

## 功能特点

- 提供 `/status` 路由来获取当前的 `sleep` 状态。
- 提供 `/change` 路由来修改 `sleep` 状态。
- 提供 `/heartbeat` 路由接收心跳信号，自动检测设备状态。
- 支持访问日志记录，将所有路由请求的 IP 地址记录到 `access.log` 文件中。
- 配置文件 `config.json` 不存在时会自动创建并填入默认值。
- 支持配置文件版本管理，自动迁移旧版本配置。
- 支持心跳超时自动设置睡眠状态。
- 支持`Render`一键部署
- 支持使用`Docker`和`Docker Compose`一键部署
- 支持使用二进制文件，无需依赖快速运行

## 配置说明

配置文件 `config.json` 包含以下字段：

```json
{
  "version": 2,           // 配置文件版本
  "sleep": false,         // 睡眠状态
  "key": "your-key",      // API密钥
  "heartbeat_enabled": false,  // 是否启用心跳检测
  "heartbeat_timeout": 60      // 心跳超时时间（秒）
}
```

### 心跳检测

当启用心跳检测时（`heartbeat_enabled=true`），服务器会：
1. 监听来自客户端的心跳信号（`/heartbeat` 路由）
2. 如果超过 `heartbeat_timeout` 秒没有收到心跳，自动将状态设置为睡眠
3. 收到心跳信号时，如果状态为睡眠，自动设置为醒来

客户端需要定期发送心跳请求：
```bash
curl "http://your-server:8000/heartbeat?key=your-key"
```

## 部署

## 快速部署

### Render

点击下方按钮，使用[Render](https://render.com/)一键部署


[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy?repo=https://github.com/shenghuo2/sleep-status)

不过Render的免费套餐，在无请求的时候会自动 down

*Free instances spin down after periods of inactivity. They do not support SSH access, scaling, one-off jobs, or persistent disks. Select any paid instance type to enable these features.*

### Docker

镜像名 `shenghuo2/sleep-status:latest`

运行命令

```shell
docker run -d -p 8000:8000 --name sleep-status shenghuo2/sleep-status:latest
```

将在本地的8000端口启动，使用`docker logs sleep-status`查看自动生成的随机密码


### Docker Compose

```yaml
# version: '3.8'
services:
  sleep-status:
    image: shenghuo2/sleep-status:latest
    ports:
      - "8000:8000"
```

将在本地的8000端口启动，使用`docker logs sleep-status`查看自动生成的随机密码

## 自行部署

### 编译成二进制文件

1. 确保已安装 Go 语言环境（推荐使用 Go 1.16 及以上版本）。
2. 下载或克隆此项目到本地：

   ```sh
   git clone https://github.com/shenghuo2/sleep-status.git
   cd sleep-status
   ```

3. 编译项目：

   ```sh
   go build .
   ```

# 使用

1. 运行编译后的可执行文件：

   ```sh
   ./sleep-status [--port=] [--host=]
   ```

   >  `[ ]`内为可选参数

   默认情况下，服务将监听在 `0.0.0.0` 的 `8000` 端口。

2. API 接口说明：

   - 获取当前的 `sleep` 状态：

     ```sh
     curl http://<host>:<port>/status
     ```

     响应示例：

     ```json
     {
       "sleep": false
     }
     ```

   - 修改 `sleep` 状态：

     ```sh
     curl "http://<host>:<port>/change?key=<your_key>&status=<status>"
     ```

     参数说明：
      - `key`：配置文件中的 `key` 值，用于校验请求合法性。
      - `status`：期望的 `sleep` 状态，1 表示 `true`，0 表示 `false`。

     成功响应示例：

     ```json
     {
       "success": true,
       "result": "status changed! now sleep is 1"
     }
     ```

     失败响应示例：

     ```json
     {
       "success": false,
       "result": "sleep already 1"
     }
     ```

    失败状态码为401

## 配置文件

默认的 `config.json` 文件格式如下：

```json
{
  "sleep": false,
  "key": "default_key"
}
```

首次运行程序时，如果配置文件不存在，将自动创建该文件并随机生成16位密钥。

## 访问日志

对 `/status` `/change` 路由的访问请求都会记录到 `access.log` 文件中，记录格式如下：

```
2024-08-03T13:18:53+08:00 - [::1]:19469 - /status

2024-08-03T13:19:01+08:00 - [::1]:19469 - /change
```

## 配套项目

### 状态上报（安卓实现）

https://github.com/shenghuo2/sleep-status-sender

可以通过该应用选择“睡醒”或“睡着”状态，将状态发送到指定的服务器。

以及一个设置页面，用户可以配置服务器的 BASE_URL，并测试与服务器的连接。

### Magisk 模块 ([magisk-module 分支](https://github.com/shenghuo2/sleep-status/tree/magisk-module))

提供 Magisk 模块实现，可以在 Root 设备上实现更深度的系统集成。

### 示例前端 ([frontend-example 分支](https://github.com/shenghuo2/sleep-status/tree/frontend-example))

提供一个基于 Web 的示例前端实现，展示如何与服务端进行交互。

[在线示例](https://blog.shenghuo2.top/test)

## 更新日志

### v0.1.0 (2025-02-25)
- 新增心跳检测功能
  - 支持自动检测设备状态
  - 可配置心跳超时时间
  - 超时自动设置睡眠状态
- 新增配置文件版本管理
  - 自动迁移旧版本配置
  - 支持配置文件版本号
- 新增部署方式
  - 支持 Render 一键部署
  - 支持 Docker 和 Docker Compose 部署
  - 提供二进制文件快速部署

### v0.0.9
- 初次启动时自动创建随机强密码 
- 对所有路由进行日志记录 
- 启动时输出当前key
- 支持 `Render` 蓝图一键部署
- 支持 `Docker`和`Docker Compose` 快速部署

### v0.0.3
- 支持 CORS 请求

### v0.0.2
- 支持入睡和醒来时间的记录

### v0.0.1
- 初始版本，提供基本功能

## TODO

- [x] 支持对所有路由的访问日志记录。
- [ ] 实现在 `config.json` 读取 **host** 和 **port** 的值
- [x] 实现对于**入睡和醒来时间** 的记录 
- [x] 支持 `CORS`
- [x] 实现第一次启动时，自动创建随机`key`
- [x] 每次启动时输出当前`key`
- [ ] 增加更多状态和操作的 API。
- [ ] 提供更加详细的错误处理和响应信息。
- [x] 支持`Paas`类方法一键部署
  - 已使用 `Render` 实现
- [x] 每次`push`自动build新镜像并推送
- [x] 支持`Docker`和`Docker Compose` 
