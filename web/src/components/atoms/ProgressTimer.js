import { useState, useEffect } from "react";
import { secondsToDuration } from "../../format/displayDate";

export const Timer = ({ tickFrequency = 500, initialTime, estimated }) => {
  const [progress, setProgress] = useState(0);

  useEffect(() => {
    const progress_timer = setTimeout(() => {
      setProgress(Date.now() - initialTime);
    }, tickFrequency);

    return () => {
      clearInterval(progress_timer);
    };
  }, [progress, tickFrequency, initialTime]);

  useEffect(() => {
    if (progress >= estimated) {
      return;
    }
  }, [estimated, progress]);

  const calculatePercentage = (progress) => {
    return Math.floor((progress / estimated) * 100);
  };

  const percentage = calculatePercentage(progress);

  const percentageDisplay = percentage >= 99 ? 99 : percentage;

  return (
    <div>
      <div className="w-full bg-slate-300 rounded-full">
        <div
          className="bg-blue-600 text-xs font-medium text-blue-100 text-center p-0.5 leading-none rounded-full"
          style={{ width: `${percentageDisplay}%` }}
        ></div>
      </div>
      {secondsToDuration(progress)}
    </div>
  );
};
