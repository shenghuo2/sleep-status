#!/bin/bash

# Set module name
MODULE_NAME="sleep_status_reporter"

# Create a clean build directory
rm -rf build
mkdir -p build

# Copy all necessary files
cp -r META-INF build/
cp module.prop build/
cp service.sh build/
cp config.conf.example build/config.conf
cp customize.sh build/

# Create empty logs directory
mkdir -p build/logs

# Set correct permissions
find build -type d -exec chmod 755 {} \;
find build -type f -exec chmod 644 {} \;
chmod 755 build/service.sh
chmod 755 build/customize.sh
chmod 755 build/META-INF/com/google/android/update-binary
chmod 755 build/META-INF/com/google/android/updater-script

# Create the zip file
cd build
zip -r ../$MODULE_NAME.zip ./*
cd ..

echo "Module has been built: $MODULE_NAME.zip"
echo "Please install this zip file through Magisk Manager"
