import { BrowserRouter as Router, Route, useHistory } from "react-router-dom";
import { Layout } from "../components/Layout";
import { Home } from "./home";
import { Projects } from "./projects";
import { Project } from "./project";
import { Session } from "./session";
import { Spec } from "./spec";
import { Emulate } from "./emulate";
import { ApiKeys } from "./apiKeys";
import { CreateApiKey } from "./createApiKey";
import { auth } from "../services/auth.service";
import { Provider } from "use-http";
import { url } from "../api";

export const Pages = () => {
  const history = useHistory()
  const options = {
    headers: {
      Accept: "application/json",
      "Content-type": "application/json",
    },
    cachePolicy: "no-cache",
    onError: ({ error }) => {
      if (error?.name) {
        const code = parseInt(error.name);
        if (!isNaN(code) && [401, 403].includes(code)) {
          auth.logout();
          history.push("/");
        }
      }
    },
    interceptors: {
      request: async ({ options, url, path, route }) => {
        if (!options?.headers?.Authorization) {
          const token = auth.get();
          if (token) {
            options.headers.Authorization = `Bearer ${token}`;
          }
        }
        return options;
      },
    },
  };

  const token = auth.get();
  if (token) {
    options.headers.Authorization = `Bearer ${token}`;
  }

  return (
    <Router>
      <Layout>
        <Provider url={url.api} options={options}>
          <Route exact path="/" component={Home} />
          <div>
            <Route path="/projects" component={Projects} />
            <Route path="/project/:id" component={Project} />
            <Route path="/spec/:id" component={Spec} />
            <Route path="/session/:id" component={Session} />
            <Route path="/emulate/:sessionId?" component={Emulate} />
            <Route path="/apiKeys" component={ApiKeys} />
            <Route path="/createApiKey" component={CreateApiKey} />
          </div>
        </Provider>
      </Layout>
    </Router>
  );
};
