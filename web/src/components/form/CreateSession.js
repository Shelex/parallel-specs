import { useCallback, useState } from "react";
import Spinner from "../atoms/Spinner";
import Loading from "../atoms/Loading";
import { EmulateSession } from "../EmulateSession";
import useFetch from "use-http";
import { endpoints } from "../../api";

export const CreateSessionForm = () => {
  const [session, setSession] = useState();
  const {
    post: createSession,
    response,
    error,
    loading,
  } = useFetch(endpoints.session);

  const [values, setValues] = useState();

  const validate = (values) => !values || !values?.projectName;

  const defaultSpecs = "a,b,c";

  const onSubmit = useCallback(
    async (e) => {
      e.preventDefault();
      const { projectName, files } = values;

      const specFiles = (files || defaultSpecs)
        .split(",")
        .filter((x) => x)
        .map((fileName) => fileName.trim());

      const res = await createSession({ projectName, specFiles });
      if (res && !res.errors) {
        setSession(res);
      }
    },
    [createSession, values]
  );

  const onChange = useCallback((e) => {
    setValues((prev) => ({
      ...prev,
      [e.target.name]: e.target.value,
    }));
  }, []);

  if (loading) {
    return <Loading />;
  }

  return session ? (
    <EmulateSession sessionResponse={session} />
  ) : (
    <div className="min-w-full flex items-center justify-center px-4">
      <div className="max-w-md w-full bg-white rounded-md p-6 shadow-2xl">
        <div className="mb-6">
          <form onSubmit={onSubmit}>
            <p>Emulate new session</p>
            <div className="mx-auto bg-white mt-4">
              <div className="mb-6">
                <label className="form-label" htmlFor="projectName">
                  Please enter project name
                </label>
                <input
                  className="form-input"
                  type="text"
                  name="projectName"
                  placeholder="Please enter name of project"
                  autoComplete="on"
                  required
                  onChange={onChange}
                />
              </div>

              <div className="mb-6">
                <label className="form-label" htmlFor="files">
                  Please enter comma-separated spec files
                </label>
                <input
                  className="form-input"
                  type="text"
                  name="files"
                  placeholder="Please enter spec files"
                  autoComplete="off"
                  defaultValue={defaultSpecs}
                  onChange={onChange}
                />
              </div>

              <div className="text-xs font-semibold text-red-500">
                {error &&
                  (response?.data?.errors?.join("; ") || response?.data)}
              </div>

              <div className="mt-12">
                <button
                  disabled={validate(values)}
                  type="submit"
                  className="bg-blue-800 w-full py-3 rounded-md text-white hover:bg-blue-900 focus:outline-none disabled:opacity-50"
                >
                  {loading ? <Spinner /> : `Create Session`}
                </button>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};
