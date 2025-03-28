# Sleep-Status

![version](https://img.shields.io/github/v/release/shenghuo2/sleep-status?include_prereleases&label=version)

[English](./README-en.md)  [简体中文](./README.md)

<details>
<summary>目录</summary>

- [项目介绍](#项目介绍)
- [功能特点](#功能特点)
- [配置说明](#配置说明)
  - [心跳检测](#心跳检测)
- [部署](#部署)
  - [快速部署](#快速部署)
    - [Render](#render)
    - [Docker](#docker)
    - [Docker Compose](#docker-compose)
  - [自行部署](#自行部署)
    - [编译成二进制文件](#编译成二进制文件)
- [使用](#使用)
  - [API接口说明](#api接口说明)
  - [配置文件](#配置文件)
  - [访问日志](#访问日志)
- [配套项目](#配套项目)
  - [状态上报（安卓实现）](#状态上报安卓实现)
  - [Magisk 模块](#magisk-模块-magisk-module-分支)
  - [示例前端](#示例前端-frontend-example-分支)
- [更新日志](#更新日志)
  - [v0.1.2 (2025-02-27)](#v012-2025-02-27)
  - [v0.1.1 (2025-02-26)](#v011-2025-02-26)
  - [v0.1.0 (2025-02-25)](#v010-2025-02-25)
  - [v0.0.9](#v009)
  - [v0.0.3](#v003)
  - [v0.0.2](#v002)
  - [v0.0.1](#v001)
- [TODO](#todo)

</details>

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

   - 获取最近的睡眠记录：

     ```sh
     curl http://<host>:<port>/records[?limit=30]
     ```

     参数说明：
     - `limit`：可选，返回的记录数量，默认为30条

     响应示例：

     ```json
     {
       "success": true,
       "records": [
         {
           "action": "sleep",
           "time": "2025-02-25T15:30:00Z"
         },
         {
           "action": "wake",
           "time": "2025-02-25T23:00:00Z"
         }
       ]
     }
     ```

   - 获取睡眠统计数据：

     ```sh
     curl http://<host>:<port>/sleep-stats[?days=7&show_time_str=0]
     ```

     参数说明：
     - `days`：可选，统计的天数范围，默认为7天
     - `show_time_str`：可选，是否显示原始时间字符串，1表示显示，0表示不显示，默认为0

     响应示例：

     ```json
     {
       "success": true,
       "stats": {
         "avg_sleep_time": "23:30",
         "avg_wake_time": "07:15",
         "avg_duration_minutes": 465,
         "periods": [
           {
             "sleep_time": 1711382400,
             "wake_time": 1711407600,
             "duration_minutes": 470
           },
           {
             "sleep_time": 1711468800,
             "wake_time": 1711493100,
             "duration_minutes": 405,
             "is_short": true
           }
         ]
       },
       "days": 2,
       "request_days": 7
     }
     ```

     字段说明：
     - `avg_sleep_time`：平均入睡时间（HH:MM格式）
     - `avg_wake_time`：平均醒来时间（HH:MM格式）
     - `avg_duration_minutes`：平均睡眠时长（分钟）
     - `periods`：睡眠周期列表
       - `sleep_time`：入睡时间（Unix时间戳，秒）
       - `wake_time`：醒来时间（Unix时间戳，秒）
       - `duration_minutes`：睡眠时长（分钟）
       - `is_short`：是否为短睡眠（小于3小时），仅当为短睡眠时才显示此字段
     - `days`：实际统计的天数
     - `request_days`：请求的天数，仅当与实际天数不同时才显示

     注意：
     - 短睡眠段（小于3小时）不计入平均入睡和醒来时间的计算，但会计入平均睡眠时长
     - 当数据不足时，`days`字段会显示实际的天数范围，而不是请求的天数

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

### v0.1.3 (2025-03-28)
- 新功能
  - 添加 `/sleep-stats` API 用于获取睡眠统计数据
  - 支持计算平均入睡时间、醒来时间和睡眠时长
  - 支持展示成对的睡眠周期（使用 Unix 时间戳）
- 改进
  - 支持短睡眠时间段（小于3小时）的处理
  - 当数据不足时，自动调整实际统计的天数
  - 可选是否显示原始时间字符串，默认不显示
- 代码优化
  - 改进了时区处理，正确处理 UTC 和本地时间转换
  - 优化了睡眠数据的计算逻辑

### v0.1.2 (2025-02-27)
- 修复问题
  - 修复 `/records` API 在并发访问时可能导致的死锁问题
  - 优化心跳检测机制，防止与其他操作产生锁竞争
  - 改进睡眠记录迁移功能的稳定性
- 性能优化
  - 使用异步日志记录，避免阻塞API响应
  - 实现细粒度锁机制，提高并发性能
  - 添加无锁版本的记录保存函数，减少锁竞争
- 代码重构
  - 将状态变更逻辑提取到独立函数，提高代码可维护性
  - 优化心跳检测器的实现，避免嵌套锁导致的死锁

### v0.1.1 (2025-02-26)
- 新增功能
  - 添加 `/records` API，支持查看最近的睡眠记录
  - 支持自定义返回记录数量（默认30条）
- 优化改进
  - 优化睡眠记录的存储格式，改用JSON数组
  - 添加旧版本记录格式自动迁移功能
  - 改进错误处理，提供更友好的错误信息

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
