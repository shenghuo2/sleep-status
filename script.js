// =================== 配置 ===================
const MY_NAME = "shenghuo2";
const API_URL = "https://sleep-status.shenghuo2.top";
const MY_GITHUB_URL = "https://github.com/shenghuo2";
const REPO_URL = "https://github.com/shenghuo2/sleep-status";
const MODULE_URL = "https://github.com/shenghuo2/sleep-status/tree/feature/magisk-module";
const FRONTEND_URL = "https://github.com/shenghuo2/sleep-status/tree/feature/example-frontend";
const INFO_FROM = "Xiaomi 14 设备在线状态";
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
      'text-xs text-gray-600 dark:text-gray-300 transition-all duration-300 z-20';
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

function testTimeConversion() {
  console.log('===== 时间转换测试 =====');
  
  const testTimes = [
    "2025-03-26T16:25:59Z",      // UTC时间，应转为UTC+8的2025-03-27 00:25:59
    "2025-03-27T00:05:55Z",      // UTC时间，应转为UTC+8的2025-03-27 08:05:55
    "2025-03-26T08:08:52+08:00", // 已经是UTC+8，应保持不变
    "2025-03-26T02:08:52+02:00"  // UTC+2时间，应转为UTC+8的2025-03-26 08:08:52
  ];
  
  testTimes.forEach(time => {
    const converted = convertToUTC8(time);
    console.log(`${time} => ${converted.toISOString()}`);
    console.log(`UTC+8时间: ${formatUTC8Time(converted)}`);
  });
  
  console.log('===== 测试结束 =====');
}

function calculateSleepSegments(records) {
  // 按日期分组记录
  const recordsByDate = {};
  
  // 首先将所有记录转换为UTC+8时间
  const utc8Records = records.map(record => {
    const utc8Time = convertToUTC8(record.time);
    // 提取日期部分作为分组键
    const dateStr = utc8Time.toISOString().split('T')[0];
    return {
      ...record,
      utc8Time,
      dateStr
    };
  });
  
  // 将所有记录按时间排序（从早到晚）
  const sortedRecords = [...utc8Records].sort((a, b) => a.utc8Time - b.utc8Time);
  
  // 找出所有的睡眠-唤醒对
  const sleepWakePairs = [];
  
  for (let i = 0; i < sortedRecords.length - 1; i++) {
    const current = sortedRecords[i];
    const next = sortedRecords[i + 1];
    
    // 如果当前记录是睡眠，下一个是唤醒，则形成一对
    if (current.action === 'sleep' && next.action === 'wake') {
      const sleepTime = current.utc8Time;
      const wakeTime = next.utc8Time;
      
      // 计算在24小时制中的位置（用于时间线显示）
      const sleepHour = sleepTime.getUTCHours() + sleepTime.getUTCMinutes() / 60;
      const wakeHour = wakeTime.getUTCHours() + wakeTime.getUTCMinutes() / 60;
      
      // 获取睡眠和唤醒的日期
      const sleepDate = sleepTime.toISOString().split('T')[0];
      const wakeDate = wakeTime.toISOString().split('T')[0];
      
      // 检查是否是跨日睡眠
      const isCrossDaySleep = sleepDate !== wakeDate;
      
      sleepWakePairs.push({
        sleepTime,
        wakeTime,
        sleepHour,
        wakeHour,
        sleepDate,
        wakeDate,
        isCrossDaySleep
      });
      
      // 跳过已处理的唤醒记录
      i++;
    }
  }
  
  // 初始化睡眠片段对象
  const sleepSegments = {};
  
  // 处理每一对睡眠-唤醒记录
  sleepWakePairs.forEach(pair => {
    const { sleepTime, wakeTime, sleepHour, wakeHour, sleepDate, wakeDate, isCrossDaySleep } = pair;
    
    // 确保睡眠日期的键存在
    if (!sleepSegments[sleepDate]) {
      sleepSegments[sleepDate] = [];
    }
    
    // 将睡眠片段添加到睡眠日期
    if (isCrossDaySleep) {
      // 跨日睡眠：在睡眠日期添加从睡眠时间到24:00的片段
      sleepSegments[sleepDate].push({
        start: sleepHour,
        end: 24,
        startTime: sleepTime,
        endTime: wakeTime,
        isCrossDaySleep: true
      });
      
      // 确保唤醒日期的键存在
      if (!sleepSegments[wakeDate]) {
        sleepSegments[wakeDate] = [];
      }
      
      // 在唤醒日期添加从00:00到唤醒时间的片段
      sleepSegments[wakeDate].push({
        start: 0,
        end: wakeHour,
        startTime: sleepTime,
        endTime: wakeTime,
        isCrossDaySleep: true,
        isWakePart: true
      });
    } else {
      // 非跨日睡眠：正常添加
      sleepSegments[sleepDate].push({
        start: sleepHour,
        end: wakeHour,
        startTime: sleepTime,
        endTime: wakeTime,
        isCrossDaySleep: false
      });
    }
  });
  
  // 检查是否有未配对的睡眠记录（最后一个记录是睡眠）
  const lastRecord = sortedRecords[sortedRecords.length - 1];
  if (lastRecord && lastRecord.action === 'sleep') {
    const sleepTime = lastRecord.utc8Time;
    const sleepHour = sleepTime.getUTCHours() + sleepTime.getUTCMinutes() / 60;
    const sleepDate = sleepTime.toISOString().split('T')[0];
    
    // 使用当前时间作为结束时间
    const now = convertToUTC8(new Date().toISOString());
    const today = now.toISOString().split('T')[0];
    
    // 确保睡眠日期的键存在
    if (!sleepSegments[sleepDate]) {
      sleepSegments[sleepDate] = [];
    }
    
    if (sleepDate === today) {
      // 如果是今天的记录，使用当前时间
      const endHour = now.getUTCHours() + now.getUTCMinutes() / 60;
      
      sleepSegments[sleepDate].push({
        start: sleepHour,
        end: endHour,
        startTime: sleepTime,
        endTime: now,
        isCrossDaySleep: false
      });
    } else {
      // 如果是之前的记录，假设跨日睡眠
      sleepSegments[sleepDate].push({
        start: sleepHour,
        end: 24,
        startTime: sleepTime,
        endTime: now,
        isCrossDaySleep: true
      });
      
      // 确保今天的键存在
      if (!sleepSegments[today]) {
        sleepSegments[today] = [];
      }
      
      // 添加从00:00到当前时间的片段
      const endHour = now.getUTCHours() + now.getUTCMinutes() / 60;
      sleepSegments[today].push({
        start: 0,
        end: endHour,
        startTime: sleepTime,
        endTime: now,
        isCrossDaySleep: true,
        isWakePart: true
      });
    }
  }
  
  return sleepSegments;
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

function updateWeeklyStats(records, recentDates) {
  // 如果没有记录，显示无数据
  if (!records || records.length < 2 || !recentDates || recentDates.length === 0) {
    document.getElementById('avgSleepTime').textContent = '--:--';
    document.getElementById('avgWakeTime').textContent = '--:--';
    document.getElementById('avgSleepDuration').textContent = '-小时-分钟';
    document.getElementById('statsDateRange').textContent = '无数据';
    document.getElementById('statsTitle').textContent = '睡眠统计:';
    return;
  }
  
  // 首先将所有记录转换为UTC+8时间并按时间排序（从早到晚）
  const utc8Records = [...records].map(record => ({
    ...record,
    utc8Time: convertToUTC8(record.time)
  })).sort((a, b) => a.utc8Time - b.utc8Time);
  
  // 将recentDates转换为UTC+8时间的日期对象
  const utc8RecentDates = recentDates.map(date => {
    const dateWithTime = date + 'T12:00:00Z';
    return convertToUTC8(dateWithTime);
  });
  
  // 使用提供的recentDates参数范围内的数据
  const earliestAllowedDate = new Date(utc8RecentDates[utc8RecentDates.length - 1]);
  earliestAllowedDate.setHours(0, 0, 0, 0);
  
  const latestAllowedDate = new Date(utc8RecentDates[0]);
  latestAllowedDate.setHours(23, 59, 59, 999);
  
  // 找出所有的睡眠-唤醒对
  const sleepWakePairs = [];
  
  for (let i = 0; i < utc8Records.length - 1; i++) {
    const current = utc8Records[i];
    const next = utc8Records[i + 1];
    
    if (current.action === 'sleep' && next.action === 'wake') {
      const sleepTime = current.utc8Time;
      const wakeTime = next.utc8Time;
      
      if (sleepTime >= earliestAllowedDate && wakeTime <= latestAllowedDate) {
        const durationMinutes = Math.abs(wakeTime - sleepTime) / (60 * 1000);
        
        const isPrimary = durationMinutes >= 180;
        
        const sleepDate = sleepTime.toISOString().split('T')[0];
        const wakeDate = wakeTime.toISOString().split('T')[0];
        const isCrossDaySleep = sleepDate !== wakeDate;
        
        sleepWakePairs.push({
          sleepTime,
          wakeTime,
          duration: durationMinutes,
          sleepDate,
          wakeDate,
          isCrossDaySleep,
          isPrimarySleep: isPrimary
        });
        
        i++;
      }
    }
  }
  
  if (sleepWakePairs.length === 0) {
    document.getElementById('avgSleepTime').textContent = '--:--';
    document.getElementById('avgWakeTime').textContent = '--:--';
    document.getElementById('avgSleepDuration').textContent = '-小时-分钟';
    document.getElementById('statsDateRange').textContent = '无数据';
    document.getElementById('statsTitle').textContent = '睡眠统计:';
    return;
  }
  
  const dailySleepCycles = {};
  
  utc8RecentDates.forEach(date => {
    const dateStr = date.toISOString().split('T')[0];
    if (!dailySleepCycles[dateStr]) {
      dailySleepCycles[dateStr] = [];
    }
  });
  
  sleepWakePairs.forEach(pair => {
    if (!dailySleepCycles[pair.sleepDate]) {
      dailySleepCycles[pair.sleepDate] = [];
    }
    dailySleepCycles[pair.sleepDate].push(pair);
  });
  
  Object.keys(dailySleepCycles).forEach(date => {
    dailySleepCycles[date].sort((a, b) => a.sleepTime - b.sleepTime);
    
    let foundPrimary = false;
    dailySleepCycles[date].forEach(cycle => {
      if (!foundPrimary && cycle.duration >= 180) {
        cycle.isPrimarySleep = true;
        foundPrimary = true;
      } else {
        cycle.isPrimarySleep = false;
      }
    });
  });
  
  let totalAdjustedSleepHours = 0;
  let totalWakeHours = 0;
  let primarySleepCount = 0;
  
  const dailyTotalDurations = {};
  
  console.log("处理睡眠-唤醒对，计算平均时间：");
  
  sleepWakePairs.forEach(pair => {
    // 只有主要睡眠（长于3小时的睡眠）才计入每日总睡眠时长
    if (pair.isPrimarySleep) {
      const date = pair.sleepDate;
      if (!dailyTotalDurations[date]) {
        dailyTotalDurations[date] = 0;
      }
      dailyTotalDurations[date] += pair.duration;
      
      // 只有主要睡眠才计入入睡和清醒时间
      // 对于睡眠时间，直接使用24小时制中的小时数和分钟数
      let sleepHour = pair.sleepTime.getUTCHours();
      const sleepMinute = pair.sleepTime.getUTCMinutes();
      
      // 将睡眠时间转换为小时（带小数）
      let sleepTimeInHours = sleepHour + sleepMinute / 60;
      
      // 如果睡眠时间在18:00以后，则将其视为"前一天的-小时"
      if (sleepHour >= 18) {
        sleepTimeInHours = sleepTimeInHours - 24;
      }
      
      // 对于唤醒时间，保持正常计算
      const wakeHour = pair.wakeTime.getUTCHours();
      const wakeMinute = pair.wakeTime.getUTCMinutes();
      const wakeTimeInHours = wakeHour + wakeMinute / 60;
      
      console.log(`睡眠记录: ${formatUTC8Time(pair.sleepTime)} - ${formatUTC8Time(pair.wakeTime)}`);
      console.log(`  原始睡眠小时: ${sleepHour}:${sleepMinute} (${sleepHour + sleepMinute/60}小时)`);
      console.log(`  调整后睡眠小时: ${sleepTimeInHours < 0 ? '-' : ''}${Math.floor(Math.abs(sleepTimeInHours))}:${Math.round(Math.abs(sleepTimeInHours - Math.floor(sleepTimeInHours)) * 60)}`);
      console.log(`  唤醒小时: ${wakeHour}:${wakeMinute} (${wakeTimeInHours}小时)`);
      console.log(`  持续时间: ${pair.duration}分钟`);
      
      totalAdjustedSleepHours += sleepTimeInHours;
      totalWakeHours += wakeTimeInHours;
      primarySleepCount++;
      
      console.log(`  累计睡眠小时: ${totalAdjustedSleepHours}, 累计唤醒小时: ${totalWakeHours}, 计数: ${primarySleepCount}`);
    }
  });
  
  // 计算每日平均睡眠时长
  let totalDailyDuration = 0;
  const daysWithSleep = Object.keys(dailyTotalDurations).length;
  
  Object.values(dailyTotalDurations).forEach(duration => {
    totalDailyDuration += duration;
  });
  
  // 计算平均值
  const avgSleepHours = primarySleepCount > 0 ? totalAdjustedSleepHours / primarySleepCount : 0;
  console.log(`平均睡眠小时(调整后): ${avgSleepHours}`);
  
  // 处理负数小时的情况（晚上的时间）
  let displayAvgSleepHours = avgSleepHours;
  if (displayAvgSleepHours < 0) {
    // 将负数小时转换为24小时制的晚上时间
    displayAvgSleepHours += 24;
    console.log(`显示用的平均睡眠小时(调整后): ${displayAvgSleepHours}`);
  }
  
  // 将小时转换为小时和分钟
  let avgSleepHour = Math.floor(displayAvgSleepHours);
  let avgSleepMinute = Math.round((displayAvgSleepHours - avgSleepHour) * 60);
  
  // 处理分钟为60的情况
  if (avgSleepMinute === 60) {
    avgSleepMinute = 0;
    avgSleepHour += 1;
    if (avgSleepHour === 24) {
      avgSleepHour = 0;
    }
  }
  
  console.log(`最终平均睡眠时间: ${avgSleepHour}:${avgSleepMinute}`);

  const avgWakeHours = primarySleepCount > 0 ? totalWakeHours / primarySleepCount : 0;
  console.log(`平均唤醒小时: ${avgWakeHours}`);
  
  // 将小时转换为小时和分钟
  let avgWakeHour = Math.floor(avgWakeHours);
  let avgWakeMinute = Math.round((avgWakeHours - avgWakeHour) * 60);
  
  // 处理分钟为60的情况
  if (avgWakeMinute === 60) {
    avgWakeMinute = 0;
    avgWakeHour += 1;
    if (avgWakeHour === 24) {
      avgWakeHour = 0;
    }
  }
  
  // 处理可能的溢出
  if (avgWakeHour >= 24) {
    avgWakeHour %= 24;
  }
  
  document.getElementById('avgSleepTime').textContent = 
    primarySleepCount > 0 ? `${String(avgSleepHour).padStart(2, '0')}:${String(avgSleepMinute).padStart(2, '0')}` : '--:--';
  
  document.getElementById('avgWakeTime').textContent = 
    primarySleepCount > 0 ? `${String(avgWakeHour).padStart(2, '0')}:${String(avgWakeMinute).padStart(2, '0')}` : '--:--';
  
  const avgDailyDurationMinutes = daysWithSleep > 0 ? Math.round(totalDailyDuration / daysWithSleep) : 0;
  const avgDailyDurationHours = Math.floor(avgDailyDurationMinutes / 60);
  const avgDailyDurationMins = avgDailyDurationMinutes % 60;

  document.getElementById('avgSleepDuration').textContent = 
    daysWithSleep > 0 ? (avgDailyDurationHours === 0 ? `${avgDailyDurationMins}m` : `${avgDailyDurationHours}h${avgDailyDurationMins > 0 ? avgDailyDurationMins + 'm' : ''}`) : '-小时-分钟';
  
  const dateFormat = date => `${date.getFullYear()}/${String(date.getMonth() + 1).padStart(2, '0')}/${String(date.getDate()).padStart(2, '0')}`;
  
  const statsTitle = document.getElementById('statsTitle');
  if (statsTitle) {
    if (recentDates.length === 1) {
      statsTitle.textContent = '单日统计:';
    } else {
      statsTitle.textContent = `近${recentDates.length}天统计:`;
    }
  }
  
  if (recentDates.length === 1) {
    const utc8Date = convertToUTC8(recentDates[0] + 'T12:00:00Z');
    document.getElementById('statsDateRange').textContent = dateFormat(utc8Date);
  } else {
    const oldestDate = convertToUTC8(recentDates[recentDates.length - 1] + 'T12:00:00Z');
    const newestDate = convertToUTC8(recentDates[0] + 'T12:00:00Z');
    document.getElementById('statsDateRange').textContent = 
      `${dateFormat(oldestDate)} ~ ${dateFormat(newestDate)}`;
  }
}

function updateSleepTimelines() {
  fetch(`${API_URL}/records`)
    .then(response => {
      if (!response.ok) {
        throw new Error('网络响应不正常');
      }
      return response.json();
    })
    .then(data => {
      if (!data.success || !Array.isArray(data.records)) {
        throw new Error('无效的数据格式');
      }

      const container = document.getElementById('sleepTimelines');
      container.innerHTML = '';
      
      const sleepSegments = calculateSleepSegments(data.records);
      
      const allDates = Object.keys(sleepSegments).sort().reverse();
      
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
          
          if (recentContinuousDates.length >= 7) {
            break;
          }
        }
      }
      
      recentContinuousDates.forEach(date => {
        const segments = sleepSegments[date].map(segment => {
          const startFormatted = formatUTC8Time(segment.startTime);
          const endFormatted = formatUTC8Time(segment.endTime);
          
          const durationHours = (segment.endTime - segment.startTime) / (1000 * 60 * 60);
          
          return {
            start: segment.start * 100 / 24,
            end: segment.end * 100 / 24,
            text: `${startFormatted}-${endFormatted} (${formatDuration(durationHours)})`
          };
        });
        
        const timeline = createTimelineElement(date, segments);
        container.appendChild(timeline);
      });
      
      updateWeeklyStats(data.records, recentContinuousDates);
      
      const timelineTitle = document.getElementById('timelineTitle');
      if (timelineTitle) {
        if (recentContinuousDates.length === 0) {
          timelineTitle.textContent = '离线记录:';
        } else if (recentContinuousDates.length === 1) {
          timelineTitle.textContent = '单日离线记录:';
        } else {
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
      console.error('获取记录失败:', error);
      document.getElementById('timelineDescription').textContent = 
        '获取记录失败: ' + error.toString();
    });
}

document.addEventListener("DOMContentLoaded", () => {
  initConfig();
  
  testTimeConversion();
  
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
