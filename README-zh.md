# 睡了吗？

一个简单的前端页面示例，用来展示我是不是睡着了。

## 功能特点

- 🎨 现代化界面
  - 自动适配深色模式（跟随系统设置）
  - 平滑动画和交互效果
  - 简约设计风格
- 🔄 每分钟自动更新状态
- ⚡️ 极简技术栈：HTML + CSS（Tailwind）+ JavaScript
- 💪 完整的加载状态和错误提示

## 使用方法

用浏览器打开 `index.html` 就可以了，页面会自动获取并显示最新状态。

## 接口说明

页面会调用这个接口：
- `GET https://sleep-status.shenghuo2.top/status` - 获取当前状态

## 开发相关

想要修改页面：
1. 直接改 `index.html` 就行，包括页面结构和样式（用 Tailwind 的类）
2. 样式用的是 Tailwind CSS 的 CDN 版本
3. JavaScript 代码就直接写在 HTML 里面了，简单粗暴

## 开源协议

MIT

---
[English](README.md)
