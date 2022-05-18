import dayjs from "dayjs";
const duration = require("dayjs/plugin/duration");
dayjs.extend(duration);

export const timestampToDate = (timestamp) => {
  return dayjs(timestamp).format("DD-MM-YYYY HH:mm:ss:SSS");
};

export const displayTimestamp = (timestamp) => {
  return timestamp > 0 ? timestampToDate(timestamp) : "_";
};

export const secondsToDuration = (unix) => {
  return dayjs.duration(unix).format("HH:mm:ss:SSS");
};
