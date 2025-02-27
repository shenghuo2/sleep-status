# Sleep-Status

![version](https://img.shields.io/github/v/release/shenghuo2/sleep-status?include_prereleases&label=version)

[English](./README-en.md)  [简体中文](./README.md)

<details>
<summary>Table of Contents</summary>

- [Project Introduction](#project-introduction)
- [Features](#features)
- [Configuration](#configuration)
  - [Heartbeat Detection](#heartbeat-detection)
- [Deployment](#deployment)
  - [Quick Deployment](#quick-deployment)
    - [Render](#render)
    - [Docker](#docker)
    - [Docker Compose](#docker-compose)
  - [Self Deployment](#self-deployment)
    - [Compile to Binary](#compile-to-binary)
- [Usage](#usage)
  - [API Description](#api-description)
  - [Configuration File](#configuration)
  - [Access Logs](#access-logs)
- [Supporting Project](#supporting-project)
  - [Status Reporting (Android Implementation)](#status-reporting-android-implementation)
- [Changelog](#changelog)
  - [v0.1.2 (2025-02-27)](#v012-2025-02-27)
  - [v0.1.1 (2025-02-26)](#v011-2025-02-26)
  - [v0.1.0 (2025-02-25)](#v010-2025-02-25)
  - [v0.0.9](#v009)
  - [v0.0.3](#v003)
  - [v0.0.2](#v002)
  - [v0.0.1](#v001)
- [TODO](#todo)

</details>

## Project Introduction

Sleep-Status is a simple backend service written in Go. The service reads and modifies the `sleep` status by reading the configuration file `config.json`.

Inspired by [this BiliBili video](https://www.bilibili.com/video/BV1fE421A7PE/).

Go was chosen because of its cross-platform capabilities, making it easy to run the service on multiple operating systems.

## Features

- Provides `/status` route to get the current `sleep` status.
- Provides `/change` route to modify the `sleep` status.
- Provides `/heartbeat` route to receive heartbeat signals for automatic device status detection.
- Supports access log recording, logging IP addresses of all route requests in the `access.log` file.
- Automatically creates the `config.json` file with default values if it does not exist.
- Supports configuration file versioning with automatic migration of old configurations.
- Supports automatic sleep status setting on heartbeat timeout.
- Supports one-click deployment with Render.
- Supports one-click deployment using Docker and Docker Compose.
- Supports running as a binary without dependencies for quick execution.

## Deployment

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

## Configuration

The `config.json` file contains the following fields:

```json
{
  "version": 2,           // Configuration file version
  "sleep": false,         // Sleep status
  "key": "your-key",      // API key
  "heartbeat_enabled": false,  // Whether to enable heartbeat detection
  "heartbeat_timeout": 60      // Heartbeat timeout in seconds
}
```

### Heartbeat Detection

When heartbeat detection is enabled (`heartbeat_enabled=true`), the server will:
1. Listen for heartbeat signals from clients (via the `/heartbeat` route)
2. Automatically set status to sleep if no heartbeat is received for `heartbeat_timeout` seconds
3. Automatically set status to awake when receiving a heartbeat while in sleep status

Clients need to send periodic heartbeat requests:
```bash
curl "http://your-server:8000/heartbeat?key=your-key"
```

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

### v0.1.2 (2025-02-27)
- Bug Fixes
  - Fixed deadlock issue in `/records` API during concurrent access
  - Optimized heartbeat detection mechanism to prevent lock contention with other operations
  - Improved stability of sleep records migration functionality
- Performance Optimization
  - Implemented asynchronous logging to avoid blocking API responses
  - Implemented fine-grained locking mechanism to improve concurrent performance
  - Added non-locking version of record saving functions to reduce lock contention
- Code Refactoring
  - Extracted status change logic to dedicated functions for better maintainability
  - Optimized heartbeat checker implementation to avoid deadlocks caused by nested locks

### v0.1.1 (2025-02-26)
- New Features
  - Added `/records` API to view recent sleep records
  - Support for customizing the number of returned records (default 30)
- Improvements
  - Optimized sleep record storage format, switched to JSON array
  - Added automatic migration functionality for old record formats
  - Improved error handling with more friendly error messages

### v0.1.0 (2025-02-25)
- Added heartbeat detection functionality
  - Automatically detect sleep status based on heartbeat signals
  - Configurable heartbeat timeout
  - Can be enabled/disabled via configuration

### v0.0.9
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