import { useCallback, useState } from "react";
import Spinner from "../atoms/Spinner";
import useFetch from "use-http";
import { endpoints } from "../../api";
import { auth } from "../../services/auth.service";

export const SignIn = ({ history }) => {
  const { post, error, response, loading } = useFetch(endpoints.login);

  const [values, setValues] = useState();

  const onSubmit = useCallback(
    async (e) => {
      e.preventDefault();
      const { email, password } = values;
      const response = await post({ email, password });
      if (response && !response.errors) {
        auth.set(response?.token);
        history.push("/");
      }
    },
    [post, values, history]
  );

  const onChange = useCallback((e) => {
    setValues((prev) => ({
      ...prev,
      [e.target.name]: e.target.value,
    }));
  }, []);

  return (
    <form onSubmit={onSubmit}>
      <div className="max-w-lg mx-auto mb-2">
        <div>
          <input
            className="form-input"
            type="email"
            name="email"
            placeholder="Please enter your email"
            autoComplete="off"
            required
            onChange={onChange}
          />
        </div>

        <div>
          <input
            className="form-input"
            type="password"
            name="password"
            placeholder="Please enter your password"
            autoComplete="off"
            required
            onChange={onChange}
          />
        </div>
      </div>

      <div className="text-xs font-semibold text-red-500">
        {error && (response?.data?.errors?.join("; ") || response?.data)}
      </div>

      <div className="mt-8">
        <button
          type="submit"
          className="bg-blue-800 w-full py-3 rounded-md text-white hover:bg-blue-900 focus:outline-none"
        >
          {loading ? <Spinner /> : <p>Sign in</p>}
        </button>
      </div>
    </form>
  );
};
