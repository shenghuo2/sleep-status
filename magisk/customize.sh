#!/sbin/sh

SKIPMOUNT=true
PROPFILE=false
POSTFSDATA=false
LATESTARTSERVICE=true

MODPATH=${0%/*}

ui_print "********************************"
ui_print "Sleep Status Reporter Installer"
ui_print "********************************"

set_perm_recursive $MODPATH 0 0 0755 0755

mkdir -p $MODPATH/logs

if [ ! -f /data/adb/modules/sleep_status_reporter/config.conf ]; then
    cp $MODPATH/config.conf.example $MODPATH/config.conf
fi

ui_print "- 配置文件位置："
ui_print "- /data/adb/modules/sleep_status_reporter/config.conf"
ui_print ""
ui_print "- 请先修改配置文件后重启手机"
