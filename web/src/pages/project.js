import { useCallback, useState, useEffect } from "react";
import { Link, useParams, useHistory } from "react-router-dom";
import ReactPaginate from "react-paginate";
import Loading from "../components/atoms/Loading";
import { ShowAlert } from "../components/atoms/Alert";
import useWebSocket from "react-use-websocket";
import { displayTimestamp, secondsToDuration } from "../format/displayDate";
import { pluralize } from "../format/text";

import { DeleteButton } from "../components/atoms/DeleteButton";
import useFetch from "use-http";
import { endpoints, url } from "../api";
import { auth } from "../services/auth.service";
import { Timer } from "../components/atoms/ProgressTimer";

const itemsPerPage = 15;

export const Project = () => {
  const { id } = useParams();
  const [currentPage, setCurrentPage] = useState(0);
  const [pageCount, setPageCount] = useState(0);
  const [itemOffset, setItemOffset] = useState(0);
  const [itemCount, setItemCount] = useState(0);
  const [project, setProject] = useState();
  const [connect, setConnect] = useState(false);

  const history = useHistory();

  const { get, del, response, loading, error } = useFetch(endpoints.projects);

  const loadProjectSessions = useCallback(async () => {
    const project = await get(
      endpoints.sessions(id, {
        limit: itemsPerPage,
        offset: itemOffset || 0,
      })
    );
    if (project && !project.errors) {
      setProject(project);
      const itemCount = project?.total;
      setItemCount(itemCount);
      setPageCount(Math.ceil(itemCount / itemsPerPage));
      setItemOffset((currentPage * itemsPerPage) % itemCount);
      setConnect(true);
    }
  }, [get, id, currentPage, itemOffset]);

  useEffect(() => {
    loadProjectSessions();
  }, [loadProjectSessions]);
  const user = auth.info();

  useWebSocket(
    url.ws,
    {
      queryParams: { projectId: id, userId: user?.id },
      onMessage: (event) => {
        let message = event.data;
        try {
          message = JSON.parse(message);
        } catch (e) {
          console.error(e);
        }

        if (message.event.kind === "created") {
          loadProjectSessions();
          return;
        }

        setProject((project) => {
          if (!project?.sessions) {
            return project;
          }
          const index = project?.sessions?.findIndex(
            (session) => session.id === message.event.id
          );
          if (index < 0) {
            return;
          }
          const session = project?.sessions[index];

          const prop = `${message.event.kind}At`;
          session[prop] = message.time;
          project.sessions[index] = session;
          return project;
        });
      },
    },
    connect
  );

  const handlePageClick = (event) => {
    setCurrentPage(event?.selected);
    const newOffset = (event.selected * itemsPerPage) % itemCount;
    setItemOffset(newOffset);
  };

  const onDelete = useCallback(
    async (e) => {
      e.preventDefault();
      del(id).then(() => {
        history.push("/projects");
      });
    },
    [del, id, history]
  );

  if (loading) {
    return <Loading />;
  }

  return error ? (
    ShowAlert(response?.data?.errors?.join("; ") || response?.data)
  ) : (
    <div className="max-w-6xl px-4 mx-auto mt-8">
      <div className="text-2xl">{project?.name}</div>
      <div>
        {project?.sessions &&
          Sessions(project, {
            itemCount,
            pageCount,
            handlePageClick,
            currentPage,
          })}
        {project && (
          <DeleteButton
            title="Delete project"
            onClick={onDelete}
            loading={loading}
          />
        )}
      </div>
    </div>
  );
};

const Sessions = (
  project,
  { itemCount, pageCount, currentPage, handlePageClick }
) => {
  const orderedSessions = [...project.sessions].sort(
    (a, b) => b?.createdAt - a?.createdAt
  );

  return (
    <div>
      <p>
        {itemCount}
        {pluralize(" session", itemCount)}
      </p>

      {orderedSessions?.length > 0 && (
        <div>
          <div>
            <table className="table-auto border-collapse border border-blue-400">
              <thead className="space-x-1">
                <tr className="bg-blue-600 px-auto py-auto">
                  <th className="w-1/5 border border-blue-400">
                    <span className="text-gray-100 font-semibold">ID</span>
                  </th>
                  <th className="w-1/8 border border-blue-400">
                    <span className="text-gray-100 font-semibold">
                      Estimated Duration
                    </span>
                  </th>
                  <th className="w-1/8 border border-blue-400">
                    <span className="text-gray-100 font-semibold">
                      Duration
                    </span>
                  </th>
                  <th className="w-1/5 border border-blue-400">
                    <span className="text-gray-100 font-semibold">Start</span>
                  </th>

                  <th className="w-1/5 border border-blue-400">
                    <span className="text-gray-100 font-semibold">End</span>
                  </th>
                </tr>
              </thead>
              <tbody className="bg-gray-200">
                {orderedSessions?.length &&
                  orderedSessions.map((session) => Session(session))}
              </tbody>
            </table>
          </div>
          <div id="container" className="flex flex-row justify-center">
            <ReactPaginate
              nextLabel="next >"
              onPageChange={handlePageClick}
              pageRangeDisplayed={3}
              marginPagesDisplayed={2}
              forcePage={currentPage}
              pageCount={pageCount}
              previousLabel="< previous"
              pageClassName="page-item"
              pageLinkClassName="page-link"
              previousClassName="page-item"
              previousLinkClassName="page-link"
              nextClassName="page-item"
              nextLinkClassName="page-link"
              breakLabel="..."
              breakClassName="page-item"
              breakLinkClassName="page-link"
              containerClassName="pagination"
              activeClassName="page-active"
              disabledClassName="page-disabled"
              renderOnZeroPageCount={null}
            />
          </div>
        </div>
      )}
    </div>
  );
};

const Session = (session) => {
  if (!session?.id) {
    return;
  }
  const displayStart = displayTimestamp(session?.startedAt);
  const displayEnd = displayTimestamp(session?.finishedAt);

  const sessionEstimatedDuration = session?.executions
    ?.map((exec) => exec?.estimatedDuration)
    .reduce((a, b) => (a += b), 0);
  const estimatedDurationText = secondsToDuration(sessionEstimatedDuration);
  const durationText = secondsToDuration(
    session.finishedAt - session.startedAt
  );

  const isStarted = session.startedAt > 0;
  const completed = session.finishedAt > 0;
  return (
    <tr key={session?.id} className="bg-white">
      <td className="font-semibold border border-blue-400">
        <Link to={`/session/${session.id}`}>{session?.id}</Link>
      </td>
      <td className="border border-blue-400">{estimatedDurationText}</td>
      <td className="border border-blue-400">
        {isStarted && !completed ? (
          <Timer
            initialTime={session.startedAt}
            estimated={sessionEstimatedDuration}
          />
        ) : (
          durationText
        )}
      </td>
      <td className="border border-blue-400">{displayStart}</td>
      <td className="border border-blue-400">{displayEnd}</td>
    </tr>
  );
};
