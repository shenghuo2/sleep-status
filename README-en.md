# Sleep Status Frontend

A simple frontend example to show sleep status.

## Features

- Modern, responsive UI
  - Dark mode support (follows system preference)
  - Smooth transitions and hover effects
  - Clean and minimalist design
- Offline Mode Effects
  - Grayscale and dimming filter
  - Floating candles with flame animations
  - DVD screensaver-style bouncing candles
  - Dark mode auto-activation
- Real-time status updates (every minute)
- Sleep Timeline Display
  - Shows recent sleep records
  - Updates every 5 minutes
  - Hourly tick marks with 6-hour major ticks
  - Optimized duration format
  - **New: Smart display of recent continuous sleep records, automatically skipping older data**
  - **New: Dynamic title showing the actual number of days with records**
- Sleep Statistics
  - **Updated: Shows statistics for the same number of days as offline records**
  - **Updated: "Daily Average Sleep Duration" replaces "Average Sleep Duration"**
  - **Updated: "Average Wake-up Time" replaces "Average Wake Time"**
  - **New: Improved algorithm - sleep segments shorter than 3 hours and secondary sleep segments don't affect average sleep/wake times**
  - **New: All sleep segments contribute to daily total sleep duration for comprehensive data**
  - **New: Timezone indicator (UTC+8), all times shown in China Standard Time**
  - Average sleep and wake-up times calculation
  - Automatic date range updates
- Built with pure HTML, CSS (Tailwind) and JavaScript
- Loading states and error handling
- Test Mode (commented out)
  - Toggle button to test online/offline effects
  - Instant visual feedback
- **New: Project Information Links**
  - **GitHub profile link**
  - **Data collection Magisk module link**
  - **Frontend source code link**

## Usage

Simply open `index.html` in your browser. The page will automatically fetch and display the sleep status from the API endpoint.

## API Integration

The frontend connects to the following endpoints:
- `GET https://sleep-status.shenghuo2.top/status` - Get current sleep status
- `GET https://sleep-status.shenghuo2.top/sleep-stats` - Get sleep statistics and records

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
