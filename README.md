
# Sleep-Status

## 项目介绍

Sleep-Status 是一个使用 Go 语言编写的简单服务。该服务通过读取配置文件 `config.json` 来获取和修改 `sleep` 状态，并通过 RESTful API 提供相关接口。灵感来源于 [这个Bilibili视频](https://www.bilibili.com/video/BV1fE421A7PE/)。使用 Go 语言是因为其跨平台交叉编译特性，使得服务可以在多个操作系统上运行。

## 功能特点

- 提供 `/status` 路由来获取当前的 `sleep` 状态。
- 提供 `/change` 路由来修改 `sleep` 状态。
- 支持访问日志记录，将 `/status` 路由请求的 IP 地址记录到 `access.log` 文件中。
- 配置文件 `config.json` 不存在时会自动创建并填入默认值。

## 安装

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

## 使用

1. 运行编译后的可执行文件：

   ```sh
   ./sleep-status [--port=] [--host=]
   ```

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

## 配置文件

默认的 `config.json` 文件格式如下：

```json
{
  "sleep": false,
  "key": "default_key"
}
```

首次运行程序时，如果配置文件不存在，将自动创建该文件并填入默认值。请修改 `key` 值以确保 API 安全。

## 访问日志

所有对 `/status` 路由的访问请求都会记录到 `access.log` 文件中，记录格式如下：

```
2023-07-24T15:04:05Z - 127.0.0.1:12345
```

## 更新日志

- 版本 1.0.0: 初始版本，提供基本功能。

## 待实现功能

- [ ] 支持对 `/change` 路由的访问日志记录。
- [ ] 增加更多状态和操作的 API。
- [ ] 提供更加详细的错误处理和响应信息。
