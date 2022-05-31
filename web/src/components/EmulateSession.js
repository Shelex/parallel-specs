import { useCallback, useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { timestampToDate, secondsToDuration } from "../format/displayDate";

import Spinner from "./atoms/Spinner";
import useFetch from "use-http";
import useWebSocket from "react-use-websocket";
import { endpoints, url } from "../api";
import { auth } from "../services/auth.service";
import { Timer } from "./atoms/ProgressTimer";

export const EmulateSession = ({ sessionResponse }) => {
  const { sessionId, projectName, projectId } = sessionResponse;
  const [session, setSession] = useState({});

  const { get, response, error, loading } = useFetch(endpoints.session);

  const fetchSession = useCallback(async () => {
    const response = await get(sessionId);
    if (response && !response.errors) {
      setSession(response);
    }
  }, [get, sessionId]);

  useEffect(() => {
    fetchSession();
  }, [fetchSession]);

  const user = auth.info();
  useWebSocket(url.ws, {
    share: true,
    queryParams: { projectId, sessionId, userId: user?.id },
    onMessage: (event) => {
      if (!session && !session?.id) {
        return session;
      }

      let message = event.data;
      try {
        message = JSON.parse(message);
      } catch (e) {
        console.error(e);
      }

      const updateProperty = `${message.event.kind}At`;

      const updateSession = (session, message) => {
        session[updateProperty] = message.time;
        return session;
      };

      const updateExecution = (session, message) => {
        const isStarted = message.event.kind === "started";

        const index = session?.executions.findIndex((item) =>
          isStarted
            ? item.specId === message.event.id
            : item.machineId === message.machineId &&
              item.startedAt > 0 &&
              item.finishedAt === 0
        );
        if (index === -1) {
          return session;
        }

        const item = session.executions[index];

        item[updateProperty] = message?.time;
        item.machineId = message.machineId;

        if (!isStarted) {
          item.duration = item.finishedAt - item.startedAt;
        }
        session.executions[index] = item;
        return session;
      };

      const update = {
        session: updateSession,
        execution: updateExecution,
      };

      const fn = update[message.event.topic] || update.session;
      const updatedSession = fn(session, message);

      setSession(updatedSession);
    },
  });

  const [machineId, setMachineId] = useState("default");

  const onChange = useCallback((e) => {
    setMachineId(e.target.value);
  }, []);

  const onNextSpec = (machineId) => async (e) => {
    e.preventDefault();
    await get(
      endpoints.next(sessionId, {
        machineId: machineId,
      })
    );
  };

  return (
    <div className="max-w-6xl px-4 mx-auto">
      <p>Session created</p>
      <p>project: {projectName}</p>
      <p>id: {sessionId}</p>
      <p className="w-max">
        <Link
          to={`/session/${sessionId}`}
          location={sessionId}
          target="_blank"
          rel="noopener noreferrer"
        >
          <button className="bg-green-500 w-full px-2  py-3 rounded-md text-white hover:bg-green-700 focus:outline-none disabled:opacity-50">
            open session
          </button>
        </Link>
      </p>
      <div>
        <SpecsTable specs={session?.executions} />
      </div>
      <div className="mt-5">
        <input
          className="form-input"
          type="text"
          name="machineId"
          defaultValue={machineId || "default"}
          placeholder="Please enter name of current machine"
          autoComplete="on"
          required
          onChange={onChange}
        />
        <div className="text-xs font-semibold text-red-500">
          {error && (response?.data?.errors?.join("; ") || response?.data)}
        </div>
        <button
          className={`bg-green-500 hover:bg-green-700 text-white font-bold py-3 px-2 mt-5 rounded w-full`}
          onClick={onNextSpec(machineId)}
        >
          {loading ? <Spinner /> : <p>Request next spec for {machineId}</p>}
        </button>
      </div>
    </div>
  );
};

const SpecsTable = ({ specs }) => {
  return specs ? (
    <div className="mt-5">
      <table className="table-auto border-collapse border border-blue-400">
        <thead className="space-x-1">
          <tr className="bg-blue-600 px-auto py-auto">
            <th className="w-1/3">
              <span className="text-gray-100 font-semibold">Name</span>
            </th>
            <th className="w-1/6">
              <span className="text-gray-100 font-semibold">Estimated</span>
            </th>

            <th className="w-1/6">
              <span className="text-gray-100 font-semibold">Duration</span>
            </th>

            <th className="w-1/6">
              <span className="text-gray-100 font-semibold">Start</span>
            </th>

            <th className="w-1/6">
              <span className="text-gray-100 font-semibold">End</span>
            </th>

            <th className="w-1/6">
              <span className="text-gray-100 font-semibold">Machine</span>
            </th>
          </tr>
        </thead>
        <tbody className="bg-gray-200">
          {specs.map((spec) => (
            <tr key={spec.id} className="bg-white">
              <td className="font-semibold border border-blue-400">
                <Link to={`/spec/${spec.specId}`}>{spec.specName}</Link>
              </td>
              <td className="border border-blue-400">
                {secondsToDuration(spec.estimatedDuration)}
              </td>

              <td className="border border-blue-400">
                {spec.startedAt > 0 && spec.finishedAt === 0 ? (
                  <Timer
                    initialTime={spec.startedAt}
                    estimated={spec.estimatedDuration}
                  />
                ) : (
                  secondsToDuration(spec.duration)
                )}
              </td>

              <td className="border border-blue-400">
                {spec.startedAt > 0 ? timestampToDate(spec.startedAt) : ""}
              </td>
              <td className="border border-blue-400">
                {spec.finishedAt > 0 ? timestampToDate(spec.finishedAt) : ""}
              </td>

              <td className="border border-blue-400">
                {spec.machineId || "none"}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  ) : (
    <p>no specs received</p>
  );
};
