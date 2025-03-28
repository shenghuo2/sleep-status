// =================== 配置 ===================
const MY_NAME = "shenghuo2";
const API_URL = "https://sleep-status.shenghuo2.top";
const MY_GITHUB_URL = "https://github.com/shenghuo2";
const REPO_URL = "https://github.com/shenghuo2/sleep-status";
const MODULE_URL = "https://github.com/shenghuo2/sleep-status/tree/feature/magisk-module";
const FRONTEND_URL = "https://github.com/shenghuo2/sleep-status/tree/feature/example-frontend";
const INFO_FROM = "Xiaomi 14 设备在线状态";
// 统计周期（天数）
// 注意：设置过大的值可能会影响性能和显示效果
const STATS_DAYS = 7;
// =================== 配置 ===================

function formatTime(hours) {
  const h = Math.floor(hours);
  const m = Math.round((hours - h) * 60);
  return `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}`;
}

function formatDuration(hours) {
  const h = Math.floor(hours);
  const m = Math.round((hours - h) * 60);
  
  if (h === 0) {
    return `${m}m`;
  } else {
    return `${h}h${m > 0 ? m + 'm' : ''}`;
  }
}

function updateTime() {
  const now = new Date();
  // 确保使用UTC+8时间
  const utc8Now = convertToUTC8(now.toISOString());
  return utc8Now.getFullYear() + '-' + 
         String(utc8Now.getMonth() + 1).padStart(2, '0') + '-' + 
         String(utc8Now.getDate()).padStart(2, '0') + ' ' + 
         String(utc8Now.getHours()).padStart(2, '0') + ':' + 
         String(utc8Now.getMinutes()).padStart(2, '0') + ':' + 
         String(utc8Now.getSeconds()).padStart(2, '0');
}

function createTimelineElement(date, segments) {
  const container = document.createElement('div');
  container.className = 'mb-8 last:mb-0';

  const dateLabel = document.createElement('div');
  dateLabel.className = 'text-sm text-gray-600 dark:text-gray-400 mb-2';
  dateLabel.textContent = date;
  container.appendChild(dateLabel);

  const timelineWrapper = document.createElement('div');
  timelineWrapper.className = 'relative';
  container.appendChild(timelineWrapper);

  // 时间线和刻度线
  const timelineContainer = document.createElement('div');
  timelineContainer.className = 'relative h-6 bg-gray-100 dark:bg-gray-700 rounded-md overflow-hidden';
  timelineWrapper.appendChild(timelineContainer);

  // 添加刻度线
  for (let hour = 0; hour <= 24; hour += 3) {
    const tick = document.createElement('div');
    tick.className = 'absolute top-0 bottom-0 w-px bg-gray-300 dark:bg-gray-600 z-10';
    tick.style.left = `${(hour / 24) * 100}%`;
    
    // 主要刻度线加粗
    if (hour % 6 === 0) {
      tick.classList.add('bg-gray-400', 'dark:bg-gray-500');
    }
    
    timelineContainer.appendChild(tick);
  }

  // 添加时间刻度
  const timeMarkers = document.createElement('div');
  timeMarkers.className = 'flex justify-between text-xs text-gray-400 dark:text-gray-500 mt-1';
  
  // 添加6小时间隔的刻度
  for (let hour = 0; hour <= 24; hour += 6) {
    const marker = document.createElement('div');
    marker.className = 'relative';
    marker.style.width = '1px'; // 仅用于定位
    
    // 刻度文本
    const text = document.createElement('div');
    text.className = 'absolute transform -translate-x-1/2';
    text.textContent = hour === 0 || hour === 24 ? hour : hour;
    marker.appendChild(text);
    
    timeMarkers.appendChild(marker);
  }
  
  timelineWrapper.appendChild(timeMarkers);

  segments.forEach(segment => {
    const segmentElement = document.createElement('div');
    segmentElement.className = 
      'absolute h-full bg-blue-200 dark:bg-blue-800 flex items-center justify-center ' +
      'text-xs text-gray-600 dark:text-gray-300 transition-all duration-300 z-20 ' +
      (segment.cssClass || '');
    segmentElement.style.left = segment.start + '%';
    segmentElement.style.width = (segment.end - segment.start) + '%';
    
    const timeText = document.createElement('div');
    timeText.className = 'px-2 whitespace-nowrap';
    timeText.textContent = segment.text;
    segmentElement.appendChild(timeText);
    
    timelineContainer.appendChild(segmentElement);
  });

  return container;
}

// 将任何时间格式转换为UTC+8时区的时间
function convertToUTC8(timeString) {
  // 解析时间字符串为Date对象
  const date = new Date(timeString);
  
  // 获取UTC时间的各个部分
  const utcYear = date.getUTCFullYear();
  const utcMonth = date.getUTCMonth();
  const utcDay = date.getUTCDate();
  const utcHours = date.getUTCHours();
  const utcMinutes = date.getUTCMinutes();
  const utcSeconds = date.getUTCSeconds();
  
  // 创建UTC+8时间（直接在UTC时间上加8小时）
  const utc8Hours = (utcHours + 8) % 24;
  const nextDay = (utcHours + 8) >= 24;
  
  // 如果加8小时后超过24小时，日期需要加1
  const utc8Date = new Date(Date.UTC(
    utcYear, 
    utcMonth, 
    utcDay + (nextDay ? 1 : 0), 
    utc8Hours,
    utcMinutes, 
    utcSeconds
  ));
  
  return utc8Date;
}

// 格式化UTC+8时间为时:分格式
function formatUTC8Time(date) {
  const hours = date.getUTCHours();
  const minutes = date.getUTCMinutes();
  return `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}`;
}





function createFloatingElements(sleep) {
  const candleContainer = document.getElementById('candleContainer');
  let candleAnimationFrame;
  
  // 清除现有的元素和动画
  candleContainer.innerHTML = '';
  if (candleAnimationFrame) {
    cancelAnimationFrame(candleAnimationFrame);
  }
  
  if (sleep) {
    // 创建蜡烛
    const candles = [];
    const candleCount = 4;
    
    for (let i = 0; i < candleCount; i++) {
      const candleDiv = document.createElement('div');
      candleDiv.className = 'floating-candle';
      
      // 添加火焰
      const flame = document.createElement('div');
      flame.className = 'candle-flame';
      flame.textContent = '🔥';
      candleDiv.appendChild(flame);
      
      // 添加蜡烛
      const candle = document.createElement('div');
      candle.textContent = '🕯️';
      candleDiv.appendChild(candle);
      
      candleContainer.appendChild(candleDiv);
      
      // 初始位置和速度
      const speed = 2;
      candles.push({
        element: candleDiv,
        x: Math.random() * (window.innerWidth - 50),
        y: Math.random() * (window.innerHeight - 50),
        dx: (Math.random() - 0.5) * speed,
        dy: (Math.random() - 0.5) * speed,
        width: 50,
        height: 80
      });
    }

    // 动画函数
    function animateCandles() {
      candles.forEach(candle => {
        candle.x += candle.dx;
        candle.y += candle.dy;

        if (candle.x <= 0 || candle.x >= window.innerWidth - candle.width) {
          candle.dx *= -1;
          candle.dy += (Math.random() - 0.5) * 0.5;
        }
        if (candle.y <= 0 || candle.y >= window.innerHeight - candle.height) {
          candle.dy *= -1;
          candle.dx += (Math.random() - 0.5) * 0.5;
        }

        const maxSpeed = 3;
        const speed = Math.sqrt(candle.dx * candle.dx + candle.dy * candle.dy);
        if (speed > maxSpeed) {
          candle.dx = (candle.dx / speed) * maxSpeed;
          candle.dy = (candle.dy / speed) * maxSpeed;
        }

        candle.element.style.transform = `translate(${candle.x}px, ${candle.y}px)`;
      });

      candleAnimationFrame = requestAnimationFrame(animateCandles);
    }

    animateCandles();

    window.addEventListener('resize', () => {
      candles.forEach(candle => {
        candle.x = Math.min(candle.x, window.innerWidth - candle.width);
        candle.y = Math.min(candle.y, window.innerHeight - candle.height);
      });
    });
  }
}

function updateStatus() {
  fetch(`${API_URL}/status`)
    .then(response => {
      if (!response.ok) {
        throw new Error('网络响应不正常: ' + response.statusText);
      }
      return response.json();
    })
    .then(data => {
      const statusDiv = document.getElementById('status');
      const descDiv = document.getElementById('description');
      const mainBody = document.getElementById('mainBody');
      console.log('响应数据:', data);
      
      if (data.sleep !== undefined) {
        updateUIState(data.sleep);
        
        // 如果有时间戳，转换为UTC+8时间
        if (data.timestamp) {
          const utc8Time = convertToUTC8(data.timestamp);
          document.getElementById('updateTime').textContent = 
            utc8Time.getFullYear() + '-' + 
            String(utc8Time.getMonth() + 1).padStart(2, '0') + '-' + 
            String(utc8Time.getDate()).padStart(2, '0') + ' ' + 
            String(utc8Time.getHours()).padStart(2, '0') + ':' + 
            String(utc8Time.getMinutes()).padStart(2, '0') + ':' + 
            String(utc8Time.getSeconds()).padStart(2, '0');
        } else {
          document.getElementById('updateTime').textContent = updateTime();
        }
      } else {
        throw new Error('无效的数据格式');
      }
    })
    .catch(error => {
      console.error('获取状态失败:', error);
      const statusDiv = document.getElementById('status');
      statusDiv.textContent = '获取状态失败 ❌';
      statusDiv.className = 'text-2xl my-3 status-error';
      document.getElementById('description').textContent = error.toString();
      document.getElementById('updateTime').textContent = updateTime();
    });
}

/**
 * 清空统计信息
 */
function clearStatistics() {
  document.getElementById('avgSleepTime').textContent = '--:--';
  document.getElementById('avgWakeTime').textContent = '--:--';
  document.getElementById('avgSleepDuration').textContent = '-小时-分钟';
  document.getElementById('statsDateRange').textContent = '无数据';
  document.getElementById('statsTitle').textContent = '睡眠统计:';
}

/**
 * 从 API 响应中更新睡眠统计信息
 * @param {Object} data - 从 /sleep-stats 获取的数据
 */
function updateWeeklyStatsFromAPI(data) {
  // 如果没有数据，显示无数据
  if (!data.success || !data.stats) {
    clearStatistics();
    return;
  }
  
  // 更新平均入睡和醒来时间
  document.getElementById('avgSleepTime').textContent = data.stats.avg_sleep_time;
  document.getElementById('avgWakeTime').textContent = data.stats.avg_wake_time;
  
  // 格式化平均睡眠时长
  const avgMinutes = data.stats.avg_duration_minutes;
  const avgHours = Math.floor(avgMinutes / 60);
  const avgMins = avgMinutes % 60;
  document.getElementById('avgSleepDuration').textContent = 
    avgHours === 0 ? `${avgMins}分钟` : `${avgHours}小时${avgMins > 0 ? avgMins + '分钟' : ''}`;
  
  // 更新统计标题
  const statsTitle = document.getElementById('statsTitle');
  if (statsTitle) {
    if (data.days === 1) {
      statsTitle.textContent = '单日统计:';
    } else {
      statsTitle.textContent = `近${data.days}天统计:`;
    }
  }
  
  // 更新日期范围
  document.getElementById('statsDateRange').textContent = `最近${data.days}天`;
}

function updateUIState(sleep) {
  const statusDiv = document.getElementById('status');
  const descDiv = document.getElementById('description');
  const mainBody = document.getElementById('mainBody');

  if (sleep) {
    statusDiv.textContent = '似了 😴';
    statusDiv.className = 'text-2xl my-3 status-sleeping';
    descDiv.textContent = '目前断网了，除了电话大概都联系不上';
    mainBody.classList.add('offline-mode');
    document.documentElement.classList.add('dark');
    createFloatingElements(true);
  } else {
    statusDiv.textContent = '活着 🌞';
    statusDiv.className = 'text-2xl my-3 status-awake';
    descDiv.textContent = '目前在线，可以通过任何可用的联系方式联系本人。';
    mainBody.classList.remove('offline-mode');
    document.documentElement.classList.remove('dark');
    createFloatingElements(false);
  }
}



function updateSleepTimelines() {
  // 使用 /sleep-stats 接口获取睡眠统计数据
  fetch(`${API_URL}/sleep-stats?days=${STATS_DAYS}&show_time_str=0&show_sleep=1`)
    .then(response => {
      if (!response.ok) {
        throw new Error('网络响应不正常');
      }
      return response.json();
    })
    .then(data => {
      if (!data.success || !data.stats || !data.stats.periods) {
        throw new Error('无效的数据格式');
      }

      const container = document.getElementById('sleepTimelines');
      container.innerHTML = '';
      
      // 按日期分组睡眠周期
      const sleepSegmentsByDate = {};
      
      // 处理每个睡眠周期
      data.stats.periods.forEach(period => {
        // 将时间戳转换为 UTC+8 时间
        const sleepTime = convertToUTC8(new Date(period.sleep_time * 1000));
        const wakeTime = convertToUTC8(new Date(period.wake_time * 1000));
        
        // 获取日期字符串（YYYY-MM-DD 格式）
        const sleepDate = sleepTime.toISOString().split('T')[0];
        const wakeDate = wakeTime.toISOString().split('T')[0];
        
        // 检查是否是跨天睡眠
        const isCrossDaySleep = sleepDate !== wakeDate;
        
        // 计算 24 小时制中的开始和结束时间（用于时间线显示）
        const sleepHour = sleepTime.getUTCHours() + sleepTime.getUTCMinutes() / 60;
        const wakeHour = wakeTime.getUTCHours() + wakeTime.getUTCMinutes() / 60;
        
        // 计算睡眠时长（小时）
        const durationHours = period.duration_minutes / 60;
        
        if (isCrossDaySleep) {
          // 跨天睡眠需要分别处理两天的片段
          
          // 第一天的睡眠时长（从入睡到午夜）
          const firstDayDuration = (24 - sleepHour);
          
          // 第二天的睡眠时长（从午夜到醒来）
          const secondDayDuration = wakeHour;
          
          // 确定哪一天的睡眠时间更长（用于决定显示文字的位置）
          const firstDayIsLonger = firstDayDuration > secondDayDuration;
          
          // 确保第一天的日期键存在
          if (!sleepSegmentsByDate[sleepDate]) {
            sleepSegmentsByDate[sleepDate] = [];
          }
          
          // 添加第一天的睡眠段（从睡眠时间到午夜）
          sleepSegmentsByDate[sleepDate].push({
            sleepTime,
            wakeTime,
            sleepHour,
            wakeHour: 24,  // 到午夜
            duration: firstDayDuration,
            isShort: period.is_short || false,
            isCrossDaySleep: true,
            isFirstDay: true,
            showText: firstDayIsLonger  // 如果第一天时间更长，则在这里显示文字
          });
          
          // 确保第二天的日期键存在
          if (!sleepSegmentsByDate[wakeDate]) {
            sleepSegmentsByDate[wakeDate] = [];
          }
          
          // 添加第二天的睡眠段（从午夜到醒来时间）
          sleepSegmentsByDate[wakeDate].push({
            sleepTime,
            wakeTime,
            sleepHour: 0,  // 从午夜开始
            wakeHour,
            duration: secondDayDuration,
            isShort: period.is_short || false,
            isCrossDaySleep: true,
            isFirstDay: false,
            showText: !firstDayIsLonger  // 如果第二天时间更长，则在这里显示文字
          });
        } else {
          // 非跨天睡眠，正常处理
          // 如果该日期还没有记录，创建一个空数组
          if (!sleepSegmentsByDate[sleepDate]) {
            sleepSegmentsByDate[sleepDate] = [];
          }
          
          // 添加睡眠段
          sleepSegmentsByDate[sleepDate].push({
            sleepTime,
            wakeTime,
            sleepHour,
            wakeHour,
            duration: durationHours,
            isShort: period.is_short || false,
            isCrossDaySleep: false,
            showText: true  // 非跨天睡眠总是显示文字
          });
        }
      });
      
      // 如果当前正在睡眠，添加当前睡眠段
      if (data.sleep && data.current_sleep_at) {
        const currentSleepTime = convertToUTC8(new Date(data.current_sleep_at * 1000));
        const now = convertToUTC8(new Date());
        
        const sleepDate = currentSleepTime.toISOString().split('T')[0];
        const currentDate = now.toISOString().split('T')[0];
        
        // 计算睡眠时长（小时）
        const durationHours = (now - currentSleepTime) / (1000 * 60 * 60);
        
        // 计算 24 小时制中的开始时间
        const sleepHour = currentSleepTime.getUTCHours() + currentSleepTime.getUTCMinutes() / 60;
        const nowHour = now.getUTCHours() + now.getUTCMinutes() / 60;
        
        // 如果该日期还没有记录，创建一个空数组
        if (!sleepSegmentsByDate[sleepDate]) {
          sleepSegmentsByDate[sleepDate] = [];
        }
        
        // 添加当前睡眠段
        if (sleepDate === currentDate) {
          // 同一天的睡眠
          sleepSegmentsByDate[sleepDate].push({
            sleepTime: currentSleepTime,
            wakeTime: now,
            sleepHour,
            wakeHour: nowHour,
            duration: durationHours,
            isShort: durationHours < 3,
            isCurrent: true,
            isCrossDaySleep: false,
            showText: true  // 非跨天睡眠总是显示文字
          });
        } else {
          // 跨天睡眠
          
          // 第一天的睡眠时长（从入睡到午夜）
          const firstDayDuration = (24 - sleepHour);
          
          // 第二天的睡眠时长（从午夜到现在）
          const secondDayDuration = nowHour;
          
          // 确定哪一天的睡眠时间更长（用于决定显示文字的位置）
          const firstDayIsLonger = firstDayDuration > secondDayDuration;
          
          // 添加第一天的睡眠段（从睡眠时间到午夜）
          sleepSegmentsByDate[sleepDate].push({
            sleepTime: currentSleepTime,
            wakeTime: now,
            sleepHour,
            wakeHour: 24,
            duration: firstDayDuration,
            isShort: false,
            isCurrent: true,
            isCrossDaySleep: true,
            isFirstDay: true,
            showText: firstDayIsLonger  // 如果第一天时间更长，则在这里显示文字
          });
          
          // 如果当前日期还没有记录，创建一个空数组
          if (!sleepSegmentsByDate[currentDate]) {
            sleepSegmentsByDate[currentDate] = [];
          }
          
          // 添加当前日期的睡眠段（从午夜到现在）
          sleepSegmentsByDate[currentDate].push({
            sleepTime: currentSleepTime,
            wakeTime: now,
            sleepHour: 0,
            wakeHour: nowHour,
            duration: secondDayDuration,
            isShort: false,
            isCurrent: true,
            isCrossDaySleep: true,
            isFirstDay: false,
            showText: !firstDayIsLonger  // 如果第二天时间更长，则在这里显示文字
          });
        }
      }
      
      // 获取所有日期并按降序排序（最近的日期在前）
      const allDates = Object.keys(sleepSegmentsByDate).sort().reverse();
      
      const recentContinuousDates = [];
      const maxDayGap = 3;
      
      if (allDates.length > 0) {
        recentContinuousDates.push(allDates[0]);
        
        for (let i = 1; i < allDates.length; i++) {
          const currentDate = new Date(allDates[i]);
          const previousDate = new Date(allDates[i-1]);
          
          const daysBetween = Math.round(
            (previousDate - currentDate) / (24 * 60 * 60 * 1000)
          );
          
          if (daysBetween > maxDayGap) {
            break;
          }
          
          recentContinuousDates.push(allDates[i]);
          
          if (recentContinuousDates.length >= STATS_DAYS) {
            break;
          }
        }
      }
      
      recentContinuousDates.forEach(date => {
        const segments = sleepSegmentsByDate[date].map(segment => {
          const startFormatted = formatUTC8Time(segment.sleepTime);
          const endFormatted = segment.isCurrent ? '现在' : formatUTC8Time(segment.wakeTime);
          
          // 处理跨天的情况
          let start = segment.sleepHour;
          let end = segment.wakeHour;
          
          // 如果结束时间小于开始时间，说明跨天了
          if (end < start) {
            end += 24; // 将结束时间调整为第二天
          }
          
          // 决定是否显示文字（对于跨天睡眠，只在一个片段上显示文字）
          const displayText = segment.showText ? `${startFormatted}-${endFormatted} (${formatDuration(segment.duration)})` : '';
          
          // 对于跨天睡眠，根据是第一天还是第二天来调整显示文字
          const segmentClass = segment.isCrossDaySleep ? 
            (segment.isFirstDay ? 'first-day-segment' : 'second-day-segment') : 
            'normal-segment';
          
          return {
            start: start * 100 / 24,
            end: end * 100 / 24,
            text: displayText,
            isCurrent: segment.isCurrent,
            cssClass: segmentClass
          };
        });
        
        const timeline = createTimelineElement(date, segments);
        container.appendChild(timeline);
      });
      
      // 更新统计信息
      updateWeeklyStatsFromAPI(data);
      
      const timelineTitle = document.getElementById('timelineTitle');
      if (timelineTitle) {
        if (recentContinuousDates.length === 0) {
          timelineTitle.textContent = '离线记录:';
        } else if (recentContinuousDates.length === 1) {
          timelineTitle.textContent = '单日离线记录:';
        } else {
          // 使用实际显示的天数，而不是配置的天数
          timelineTitle.textContent = `近${recentContinuousDates.length}天离线记录:`;
        }
      }
      
      if (recentContinuousDates.length === 0) {
        document.getElementById('timelineDescription').textContent = '暂无离线记录';
      } else if (recentContinuousDates.length === 1) {
        document.getElementById('timelineDescription').textContent = `显示 ${recentContinuousDates[0]} 的离线记录`;
      } else {
        const oldestDate = recentContinuousDates[recentContinuousDates.length - 1];
        const newestDate = recentContinuousDates[0];
        document.getElementById('timelineDescription').textContent = 
          `显示 ${oldestDate} 至 ${newestDate} 的离线记录`;
      }
    })
    .catch(error => {
      console.error('获取睡眠数据失败:', error);
      document.getElementById('timelineDescription').textContent = 
        '获取睡眠数据失败: ' + error.toString();
      
      // 清空统计信息
      clearStatistics();
    });
}

document.addEventListener("DOMContentLoaded", () => {
  initConfig();
  
  
  updateStatus();
  updateSleepTimelines();
  
  const testButton = document.getElementById('testButton');
  let testMode = false;
  testButton.addEventListener('click', () => {
    testMode = !testMode;
    updateUIState(testMode);
  });
});

function initConfig() {
  document.title = MY_NAME + '现在似了吗?';
  
  if (document.getElementById('my-name')) {
    document.getElementById('my-name').textContent = MY_NAME;
  }
  if (document.getElementById('my-name1')) {
    document.getElementById('my-name1').textContent = MY_NAME;
  }
  if (document.getElementById('my-name2')) {
    document.getElementById('my-name2').textContent = MY_NAME;
  }
  
  if (document.getElementById('my-status')) {
    document.getElementById('my-status').textContent = MY_NAME + "'s Status:";
  }
  
  if (document.getElementById('my-github')) {
    document.getElementById('my-github').href = MY_GITHUB_URL;
  }
  if (document.getElementById('repo-url')) {
    document.getElementById('repo-url').href = REPO_URL;
  }
  if (document.getElementById('module-url')) {
    document.getElementById('module-url').href = MODULE_URL;
  }
  if (document.getElementById('frontend-url')) {
    document.getElementById('frontend-url').href = FRONTEND_URL;
  }
  
  if (document.getElementById('info-from')) {
    document.getElementById('info-from').textContent = INFO_FROM;
  }
}

setInterval(updateStatus, 60000);
setInterval(updateSleepTimelines, 300000);
