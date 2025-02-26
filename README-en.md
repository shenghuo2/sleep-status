# Sleep Status Frontend

A simple frontend example to show sleep status.

## Features

- 🎨 Modern, responsive UI
  - Dark mode support (follows system preference)
  - Smooth transitions and hover effects
  - Clean and minimalist design
- 🌙 Offline Mode Effects
  - Grayscale and dimming filter
  - Floating candles with flame animations
  - DVD screensaver-style bouncing candles
  - Dark mode auto-activation
- 🔄 Real-time status updates (every minute)
- 📊 Sleep Timeline Display
  - Shows recent sleep records
  - Updates every 5 minutes
  - Time scale markers (0, 6, 12, 18, 24 hours)
  - Optimized duration format
- 📈 Sleep Statistics
  - Average sleep and wake times over the past 7 days
  - Average sleep duration calculation
  - Automatic date range updates
- ⚡️ Built with pure HTML, CSS (Tailwind) and JavaScript
- 💪 Loading states and error handling
- 🧪 Test Mode
  - Toggle button to test online/offline effects
  - Instant visual feedback

## Usage

Simply open `index.html` in your browser. The page will automatically fetch and display the sleep status from the API endpoint.

The test button in the top-right corner allows you to preview the online/offline state transitions and effects.

## API Integration

The frontend connects to the following endpoints:
- `GET https://sleep-status.shenghuo2.top/status` - Get current sleep status
- `GET https://sleep-status.shenghuo2.top/records` - Get sleep timeline records

## Development

To modify the frontend:
1. Edit `index.html` to change the structure or Tailwind classes
2. The page uses Tailwind CSS via CDN for styling
3. JavaScript is embedded in the HTML file for simplicity

## Visual Effects

When the status is offline:
- The page enters dark mode
- A grayscale and dimming filter is applied
- Four animated candles appear and bounce around the screen
- Each candle has a flickering flame effect
- The overall brightness is reduced

## License

MIT

---
[简体中文](README.md)
