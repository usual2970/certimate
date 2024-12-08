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

export const getLeftDays = (zuluTime: string) => {
  const time = convertZulu2Beijing(zuluTime);
  const date = time.split(" ")[0];
  const now = new Date();
  const target = new Date(date);
  const diff = target.getTime() - now.getTime();
  const days = Math.floor(diff / (1000 * 60 * 60 * 24));
  return days;
};

export const diffDays = (date1: string, date2: string) => {
  const target1 = new Date(date1);
  const target2 = new Date(date2);
  const diff = target1.getTime() - target2.getTime();
  const days = Math.ceil(diff / (1000 * 60 * 60 * 24));
  return days;
};
