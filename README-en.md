# Sleep-Status

![version](https://img.shields.io/github/v/release/shenghuo2/sleep-status?include_prereleases&label=version)

[English](./README-en.md)  [简体中文](./README.md)

## Project Introduction

Sleep-Status is a simple backend service written in Go. The service reads and modifies the `sleep` status by reading the configuration file `config.json`.

Inspired by [this BiliBili video](https://www.bilibili.com/video/BV1fE421A7PE/).

Go was chosen because of its cross-platform capabilities, making it easy to run the service on multiple operating systems.

## Features

- Provides `/status` route to get the current `sleep` status.
- Provides `/change` route to modify the `sleep` status.
- Supports access log recording, logging IP addresses of requests to `/status` and `/change` routes in the `access.log` file.
- Automatically creates the `config.json` file with default values if it does not exist.
- Supports one-click deployment with Render.
- Supports one-click deployment using Docker and Docker Compose.
- Supports running as a binary without dependencies for quick execution.

# Deployment

## Quick Deployment

### Render

Click the button below to deploy with [Render](https://render.com/) in one click

[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy?repo=https://github.com/shenghuo2/sleep-status)

However, the free tier of Render will automatically spin down when there are no requests.

*Free instances spin down after periods of inactivity. They do not support SSH access, scaling, one-off jobs, or persistent disks. Select any paid instance type to enable these features.*

### Docker

Image name: `shenghuo2/sleep-status:latest`

Run command:

```shell
docker run -d -p 8000:8000 --name sleep-status shenghuo2/sleep-status:latest
```

The service will start on port 8000. Use `docker logs sleep-status` to view the automatically generated random password.

### Docker Compose

```yaml
version: '3.8'
services:
  sleep-status:
    image: shenghuo2/sleep-status:latest
    ports:
      - "8000:8000"
```

The service will start on port 8000. Use `docker logs sleep-status` to view the automatically generated random password.

## Self Deployment

### Compile to Binary

1. Ensure Go is installed (recommended version Go 1.16 or above).
2. Download or clone this project to your local machine:

   ```sh
   git clone https://github.com/shenghuo2/sleep-status.git
   cd sleep-status
   ```

3. Compile the project:

   ```sh
   go build .
   ```

# Usage

1. Run the compiled executable:

   ```sh
   ./sleep-status [--port=] [--host=]
   ```

   > Parameters in `[ ]` are optional.

   By default, the service will listen on `0.0.0.0` port `8000`.

2. API Description:

    - Get the current `sleep` status:

      ```sh
      curl http://<host>:<port>/status
      ```

      Example response:

      ```json
      {
        "sleep": false
      }
      ```

    - Modify the `sleep` status:

      ```sh
      curl "http://<host>:<port>/change?key=<your_key>&status=<status>"
      ```

      Parameter description:
        - `key`: The `key` value in the configuration file, used to validate the request.
        - `status`: The desired `sleep` status, 1 for `true`, 0 for `false`.

      Successful response example:

      ```json
      {
        "success": true,
        "result": "status changed! now sleep is 1"
      }
      ```

      Failed response example:

      ```json
      {
        "success": false,
        "result": "sleep already 1"
      }
      ```

   Failure status code is 401

## Configuration File

The default `config.json` file format is as follows:

```json
{
  "sleep": false,
  "key": "default_key"
}
```

When the program runs for the first time, if the configuration file does not exist, it will automatically create the file and generate a random 16-character key.

## Access Logs

Requests to `/status` and `/change` routes will be logged in the `access.log` file in the following format:

```
2024-08-03T13:18:53+08:00 - [::1]:19469 - /status

2024-08-03T13:19:01+08:00 - [::1]:19469 - /change
```

## Supporting Project

### Status Reporting (Android Implementation)

https://github.com/shenghuo2/sleep-status-sender

This app allows users to select either "Awake" or "Asleep" status and send the selected status to a designated server.

There is also a settings page where users can configure the server's BASE_URL and test the connection to the server.

# Others

## Changelog

- v0.0.9
    - Automatically create a random strong password on first startup
    - Log all routes
    - Output the current key at startup
    - Support `Render` one-click deployment blueprint
    - Support `Docker` and `Docker Compose` quick deployment

- v0.0.3: Support CORS requests
- v0.0.2: Support recording sleep and wake times
- v0.0.1: Initial version, providing basic functionality.

## TODO

- [x] Support access log recording for all routes.
- [ ] Implement reading **host** and **port** values from `config.json`
- [x] Implement recording **sleep and wake times**
- [x] Support `CORS`
- [x] Automatically create a random `key` on first startup
- [x] Output the current `key` at startup
- [ ] Add more status and operation APIs.
- [ ] Provide more detailed error handling and response information.
- [x] Support one-click deployment for `Paas`
    - Implemented using `Render`
- [x] Automatically build and push new images on each `push`
- [x] Support `Docker` and `Docker Compose`