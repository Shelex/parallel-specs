import { memo, useCallback, useState } from "react";
import Spinner from "../atoms/Spinner";
import useFetch from "use-http";
import { endpoints } from "../../api";
import { auth } from "../../services/auth.service";

const SignUpForm = ({ history }) => {
  const { post, error, response, loading } = useFetch(endpoints.register);

  const [values, setValues] = useState();

  const validate = (values) =>
    !values ||
    !values?.email ||
    !values?.password ||
    !values?.passwordConfirm ||
    values?.password?.length < 3 ||
    values?.password !== values?.passwordConfirm;

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
    [history, post, values]
  );

  const onChange = useCallback((e) => {
    setValues((prev) => ({
      ...prev,
      [e.target.name]: e.target.value,
    }));
  }, []);

  return (
    <form onSubmit={onSubmit}>
      <div className="max-w-lg mx-auto bg-white p-4">
        <div className="mb-4">
          <label className="form-label" htmlFor="email">
            Email
          </label>
          <input
            className="form-input"
            type="email"
            name="email"
            id="email"
            placeholder="Please enter your email"
            autoComplete="off"
            required
            onChange={onChange}
          />
        </div>

        <div className="mb-4">
          <label className="form-label" htmlFor="password">
            Password
          </label>
          <input
            className="form-input"
            type="password"
            name="password"
            id="password"
            placeholder="Please enter your password, min 4 chars"
            autoComplete="off"
            required
            onChange={onChange}
          />
        </div>

        <div className="mb-4">
          <label className="form-label" htmlFor="password-check">
            Confirm Password
          </label>
          <input
            type="password"
            name="passwordConfirm"
            id="password-check"
            className="form-input"
            placeholder="Please enter your password again"
            autoComplete="off"
            required
            onChange={onChange}
          />
        </div>

        <div className="text-xs font-semibold text-red-500">
          {error && (response?.data?.errors?.join("; ") || response?.data)}
        </div>

        <div className="mt-12">
          <button
            disabled={validate(values)}
            type="submit"
            className="bg-blue-800 w-full py-3 rounded-md text-white hover:bg-blue-900 focus:outline-none disabled:opacity-50"
          >
            {loading ? <Spinner /> : `Sign up`}
          </button>
        </div>
      </div>
    </form>
  );
};

export default memo(SignUpForm);
