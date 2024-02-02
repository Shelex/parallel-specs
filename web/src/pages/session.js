import { useCallback, useEffect, useState } from "react";
import { Link, useParams, useHistory } from "react-router-dom";
import useFetch from "use-http";
import useWebSocket from "react-use-websocket";
import Loading from "../components/atoms/Loading";
import { ShowAlert } from "../components/atoms/Alert";
import { DeleteButton } from "../components/atoms/DeleteButton";
import { Timer } from "../components/atoms/ProgressTimer";
import { displayTimestamp, secondsToDuration } from "../format/displayDate";
import { defineSpecColor } from "../format/specStatus";
import { endpoints, url } from "../api";
import { auth } from "../services/auth.service";

export const Session = () => {
  const { id } = useParams();
  const [session, setSession] = useState();

  const history = useHistory();

  const { get, del, response, loading, error } = useFetch(endpoints.session);

  const loadSession = useCallback(async () => {
    const response = await get(id);
    if (response && !response.errors) {
      setSession(response);
    }
  }, [get, id]);

  useEffect(() => {
    loadSession();
  }, [loadSession]);

  const user = auth.info();

  useWebSocket(url.ws, {
    share: true,
    queryParams: {
      projectId: session?.projectId,
      sessionId: id,
      userId: user?.id,
    },
    onMessage: (event) => {
      if (!session?.id && !session?.projectId) {
        return;
      }

      let message = event.data;
      try {
        message = JSON.parse(message);
      } catch (e) {
        console.error(e);
      }

      if (message.event.id !== id && message.sessionId !== id) {
        return;
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
        !isStarted && message.status && (item.status = message.status);

        if (!isStarted) {
          item.duration = item.finishedAt - item.startedAt;
        }
        session.executions[index] = item;
        return session;
      };

      const update = {
        session: updateSession,
        execution: updateExecution,
        default: (session) => session,
      };

      const fn = update[message.event.topic] || update.default;
      const updatedSession = fn(session, message);

      setSession(updatedSession);
    },
  });

  const onDelete = useCallback(
    async (e) => {
      e.preventDefault();
      const projectId = session?.projectId;
      del(id).then(() => {
        history.push(projectId ? `/project/${projectId}` : `/projects`);
      });
    },
    [del, id, session, history]
  );

  if (loading) {
    return <Loading />;
  }

  if (error) {
    return ShowAlert(response);
  }

  const statsByMachines = calculateStatByMachine(session);

  return (
    <div className="max-w-6xl px-4 mx-auto mt-8">
      {session?.projectId && (
        <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-1 px-2 rounded">
          <Link to={`/project/${session?.projectId}`}>Back to project</Link>
        </button>
      )}
      <div className="text-2xl">Session id "{session?.id}"</div>
      <div className="text-1xl">
        {displayTimestamp(session?.startedAt)} :{" "}
        {displayTimestamp(session?.finishedAt)}
      </div>
      {session && Specs(session, statsByMachines)}
      <div className="mt-10">{ByMachine(statsByMachines)}</div>
      {session?.startedAt && !session?.finishedAt && (
        <button className="bg-yellow-500 hover:bg-yellow-700 text-white font-bold py-1 px-2 rounded mt-5">
          <Link to={`/emulate/${id}`}>Emulate manually</Link>
        </button>
      )}
      <DeleteButton
        title="Delete session"
        onClick={onDelete}
        loading={loading}
      />
    </div>
  );
};

const Specs = (session, statsByMachines) => {
  const { executions } = session;
  return (
    executions?.length && (
      <div>
        <div>{executions?.length} files</div>
        {executions.length > 0 && (
          <table className="table-auto border-collapse border border-blue-400">
            <thead className="space-x-1">
              <tr className="bg-blue-600 px-auto py-auto">
                <th className="w-1/2 border border-blue-400">
                  <span className="text-gray-100 font-semibold">FileName</span>
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
                <th className="w-1/6 border border-blue-400">
                  <span className="text-gray-100 font-semibold">Machine</span>
                </th>
                <th className="w-1/6 border border-blue-400">
                  <span className="text-gray-100 font-semibold">
                    Duration %
                  </span>
                </th>
              </tr>
            </thead>
            <tbody className="bg-gray-200">
              {executions?.length &&
                [...executions]
                  .sort((a, b) => b.estimatedDuration - a.estimatedDuration)
                  .map((execution) => Spec(execution, statsByMachines))}
            </tbody>
          </table>
        )}
      </div>
    )
  );
};

const Spec = (execution, statsByMachines) => {
  const machineStat = statsByMachines.stats.find(
    (stat) => stat.machine === execution.machineId
  );

  if (!machineStat) {
    return;
  }

  const color = defineSpecColor(execution);
  return (
    <tr key={execution.specName} className="bg-white">
      <td className="font-semibold border border-blue-400">
        <Link to={`/spec/${execution.specId}`}>{execution.specName}</Link>
      </td>
      <td className="border border-blue-400">
        {secondsToDuration(execution.estimatedDuration)}
      </td>
      <td className="border border-blue-400">
        {execution.startedAt > 0 && execution.finishedAt === 0 ? (
          <Timer
            initialTime={execution.startedAt}
            estimated={execution.estimatedDuration}
          />
        ) : (
          secondsToDuration(execution.duration)
        )}
      </td>
      <td className={`border border-blue-400 bg-${color}`}>
        {execution.status}
      </td>
      <td className="border border-blue-400">{execution.machineId}</td>
      <td className="border border-blue-400">
        {machineStat?.duration > 0
          ? ((execution.duration / machineStat?.duration) * 100).toFixed(2)
          : "n/a"}{" "}
        %
      </td>
    </tr>
  );
};

const calculateStatByMachine = (session) => {
  if (!session || !session?.executions) {
    return {};
  }

  const machines = Array.from(
    new Set(session?.executions.map((item) => item?.machineId).filter((x) => x))
  );

  const stats = machines.map((machine) => {
    return {
      machine: machine,
      duration: session?.executions
        .filter((execution) => execution.machineId === machine)
        .map((execution) => execution.duration)
        .reduce((a, b) => a + b, 0),
    };
  });

  return {
    machines,
    stats,
  };
};

const ByMachine = (perMachine) => {
  if (!perMachine?.machines?.length) {
    return;
  }

  const { machines, stats } = perMachine;

  return (
    <div>
      <div>{machines.length} machines</div>
      {machines.length > 0 && (
        <table className="table-auto border-collapse border border-blue-400 w-full">
          <thead className="space-x-1">
            <tr className="bg-blue-600 px-auto py-auto">
              <th className="w-1/2 border border-blue-400">
                <span className="text-gray-100 font-semibold">MachineID</span>
              </th>
              <th className="w-1/6 border border-blue-400">
                <span className="text-gray-100 font-semibold">Duration</span>
              </th>
            </tr>
          </thead>
          <tbody className="bg-gray-200">
            {stats
              .sort((a, b) => a.machine.localeCompare(b.machine))
              .map((stat) => (
                <tr key={stat.machine} className="bg-white">
                  <td className="font-semibold border border-blue-400">
                    {stat.machine}
                  </td>
                  <td className="border border-blue-400">
                    {secondsToDuration(stat.duration)}
                  </td>
                </tr>
              ))}
          </tbody>
        </table>
      )}
    </div>
  );
};
