import { useCallback, useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { secondsToDuration, displayTimestamp } from "../format/displayDate";
import { defineSpecColor } from "../format/specStatus";
import Loading from "../components/atoms/Loading";
import { ShowAlert } from "../components/atoms/Alert";
import useFetch from "use-http";
import { endpoints } from "../api";

export const Spec = () => {
  const { id } = useParams();

  const [spec, setSpec] = useState();

  const { get, response, loading, error } = useFetch(endpoints.spec);

  const loadSpec = useCallback(async () => {
    const response = await get(id);
    if (response && !error) {
      setSpec(response);
    }
  }, [get, id, error]);

  useEffect(() => {
    loadSpec();
  }, [loadSpec]);

  if (error) {
    return ShowAlert(response);
  }

  if (loading) {
    return <Loading />;
  }

  return (
    <div className="max-w-6xl px-4 mx-auto mt-8">
      <div className="text-2xl">Project "{spec?.project?.name}"</div>
      <div className="text-2xl">File "{spec?.name}"</div>
      <div>{spec?.executions.length} latest executions</div>
      {spec?.executions.length > 0 && (
        <table className="table-auto border-collapse border border-blue-400 w-full">
          <thead className="space-x-1">
            <tr className="bg-blue-600 px-auto py-auto">
              <th className="w-1/2 border border-blue-400">
                <span className="text-gray-100 font-semibold">Session</span>
              </th>
              <th className="w-1/8 border border-blue-400">
                <span className="text-gray-100 font-semibold">Start</span>
              </th>
              <th className="w-1/8 border border-blue-400">
                <span className="text-gray-100 font-semibold">End</span>
              </th>
              <th className="w-1/8 border border-blue-400">
                <span className="text-gray-100 font-semibold">Estimated</span>
              </th>
              <th className="w-1/8 border border-blue-400">
                <span className="text-gray-100 font-semibold">Duration</span>
              </th>
              <th className="w-1/6 border border-blue-400">
                <span className="text-gray-100 font-semibold">Status</span>
              </th>
            </tr>
          </thead>
          <tbody className="bg-gray-200">
            {spec?.executions
              .sort((a, b) => b.startedAt - a.startedAt)
              .map((stat) => {
                const color = defineSpecColor(stat);
                return (
                  <tr key={stat.sessionId} className="bg-white">
                    <td className="font-semibold border border-blue-400">
                      <Link to={`/session/${stat.sessionId}`}>
                        {stat.sessionId}
                      </Link>
                    </td>
                    <td className="border border-blue-400">
                      {displayTimestamp(stat.startedAt)}
                    </td>
                    <td className="border border-blue-400">
                      {displayTimestamp(stat.finishedAt)}
                    </td>
                    <td className="border border-blue-400">
                      {secondsToDuration(stat.estimatedDuration)}
                    </td>
                    <td className="border border-blue-400">
                      {secondsToDuration(stat.duration)}
                    </td>
                    <td className={`border border-blue-400 bg-${color}`}>
                      {stat.status}
                    </td>
                  </tr>
                );
              })}
          </tbody>
        </table>
      )}
    </div>
  );
};
