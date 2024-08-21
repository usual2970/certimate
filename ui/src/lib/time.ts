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
