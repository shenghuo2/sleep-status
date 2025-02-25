# Sleep Status - Magisk Module

![version](https://img.shields.io/github/v/release/shenghuo2/sleep-status?include_prereleases&label=version)

[English](./README.md) | [简体中文](./README-zh.md)

This is the Magisk module for the Sleep Status project. It sends heartbeat signals to the server and automatically reports sleep status when network connection is lost.

## Features

- Automatically detect network connectivity through heartbeat signals
- Report sleep status when network is unavailable
- Support heartbeat monitoring
- Runs as a system service via Magisk
- Supports Android ARM64 devices

## Installation

1. Download the latest `sleep-status-magisk.zip` from [Releases](https://github.com/shenghuo2/sleep-status/releases)
2. Install through Magisk Manager
3. Edit `/data/adb/modules/sleep-status/config.conf` to set your server URL and API key
4. Reboot your device for changes to take effect

## Configuration

The configuration file (`config.conf`) contains:
```
SERVER_URL="http://your-server:8000"  # Your server URL
API_KEY="your-key"                    # API key for authentication
```

After modifying the configuration, you must reboot your device for changes to take effect.

## Logs

You can check the module's operation in Magisk logs:
```bash
tail -f /cache/magisk.log | grep sleep-status
```

## Troubleshooting

1. If the module is not showing up in Magisk Manager:
   - Make sure you're using the latest Magisk version
   - Check if the module is properly installed in `/data/adb/modules/`

2. If status is not being reported:
   - Check if the configuration file exists and is properly configured
   - Verify server URL is accessible from your device
   - Check Magisk logs for any errors
   - Make sure you've rebooted after configuration changes

## Building from Source

The module can be built using the provided build script:
```bash
./build.sh
```

This will create a flashable zip file that can be installed through Magisk Manager.

## Contributing

Feel free to submit issues and pull requests to improve the module.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
