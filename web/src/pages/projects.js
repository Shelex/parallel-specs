import { useCallback, useEffect, useState } from "react";
import { Link } from "react-router-dom";
import Loading from "../components/atoms/Loading";
import { ShowAlert } from "../components/atoms/Alert";
import useFetch from "use-http";
import useWebSocket from "react-use-websocket";
import { endpoints, url } from "../api";
import { auth } from "../services/auth.service";

export const Projects = () => {
  const [projects, setProjects] = useState([]);
  const { get, response, loading, error } = useFetch(endpoints.projects);

  const [connect, setConnect] = useState(false);

  const loadProjects = useCallback(async () => {
    const response = await get();
    if (response) {
      setProjects(response.projects);
      setConnect(true);
    }
  }, [get]);

  useEffect(() => {
    loadProjects();
  }, [loadProjects]);

  const user = auth.info();
  useWebSocket(
    url.ws,
    {
      queryParams: { userId: user?.id },
      onMessage: (event) => {
        let message = event.data;
        try {
          message = JSON.parse(message);
        } catch (e) {
          console.error(e);
        }
        if (message?.event?.topic !== "project") {
          return;
        }
        loadProjects();
      },
    },
    connect
  );

  if (loading) {
    return <Loading />;
  }

  return (
    <div className="max-w-6xl px-4 mx-auto mt-8">
      {error ? (
        ShowAlert(response?.data?.errors?.join("; ") || response?.data)
      ) : projects.length ? (
        <div>
          <div className="text-2xl">Projects:</div>
          <div className="grid gap-3 grid-cols-3 mt-10">
            {projects.map((project) => ProjectItem(project))}
          </div>
        </div>
      ) : (
        <ProjectsEmpty />
      )}
    </div>
  );
};

const ProjectsEmpty = () => {
  return (
    <div className="max-w-6xl px-4 mx-auto mt-8">
      No projects available. You can integrate with:
      <li key="playground">
        api docs, schema and playground are available at
        <a
          className="text-blue-600 mx-2"
          href="https://split-specs.shelex.dev/swagger"
        >
          Swagger page,
        </a>
      </li>
      <li key="client">
        <a
          className="text-blue-600 mx-2"
          href="https://github.com/Shelex/split-specs-client"
        >
          js client library
        </a>
      </li>
      <li key="emulation">
        or just check
        <Link className="text-blue-600 mx-2" to="/emulate">
          session emulation
        </Link>
        to find out how it is working :)
      </li>
    </div>
  );
};

const ProjectItem = ({ name, id }) => {
  return (
    <Link to={`/project/${id}`} key={id}>
      <div className="rounded-md py-3 px-6 inline-block border-2 border-blue-600 items-center">
        <p className="align-middle break-all">{name}</p>
      </div>
    </Link>
  );
};
