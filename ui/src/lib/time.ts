export const convertZulu2Beijing = (zuluTime: string) => {
  const utcDate = new Date(zuluTime);

  // Format the Beijing time
  const formattedBeijingTime = new Intl.DateTimeFormat("zh-CN", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false,
    timeZone: "Asia/Shanghai",
  }).format(utcDate);

  return formattedBeijingTime;
};

export const getDate = (zuluTime: string) => {
  const time = convertZulu2Beijing(zuluTime);
  return time.split(" ")[0];
};

export const getLeftDays = (zuluTime: string) => {
  const time = convertZulu2Beijing(zuluTime);
  const date = time.split(" ")[0];
  const now = new Date();
  const target = new Date(date);
  const diff = target.getTime() - now.getTime();
  const days = Math.floor(diff / (1000 * 60 * 60 * 24));
  return days;
};

export function getTimeBefore(days: number): string {
  // 获取当前时间
  const currentDate = new Date();

  // 减去指定的天数
  currentDate.setUTCDate(currentDate.getUTCDate() - days);

  // 格式化日期为 yyyy-mm-dd
  const year = currentDate.getUTCFullYear();
  const month = String(currentDate.getUTCMonth() + 1).padStart(2, "0"); // 月份从 0 开始
  const day = String(currentDate.getUTCDate()).padStart(2, "0");

  // 格式化时间为 hh:ii:ss
  const hours = String(currentDate.getUTCHours()).padStart(2, "0");
  const minutes = String(currentDate.getUTCMinutes()).padStart(2, "0");
  const seconds = String(currentDate.getUTCSeconds()).padStart(2, "0");

  // 组合成 yyyy-mm-dd hh:ii:ss 格式
  const formattedDate = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;

  return formattedDate;
}

export function getTimeAfter(days: number): string {
  // 获取当前时间
  const currentDate = new Date();

  // 加上指定的天数
  currentDate.setUTCDate(currentDate.getUTCDate() + days);

  // 格式化日期为 yyyy-mm-dd
  const year = currentDate.getUTCFullYear();
  const month = String(currentDate.getUTCMonth() + 1).padStart(2, "0"); // 月份从 0 开始
  const day = String(currentDate.getUTCDate()).padStart(2, "0");

  // 格式化时间为 hh:ii:ss
  const hours = String(currentDate.getUTCHours()).padStart(2, "0");
  const minutes = String(currentDate.getUTCMinutes()).padStart(2, "0");
  const seconds = String(currentDate.getUTCSeconds()).padStart(2, "0");

  // 组合成 yyyy-mm-dd hh:ii:ss 格式
  const formattedDate = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;

  return formattedDate;
}

