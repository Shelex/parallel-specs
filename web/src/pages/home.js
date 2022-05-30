import { useEffect, useCallback, useState } from "react";
import { useHistory } from "react-router-dom";
import { SignIn } from "../components/form/SignIn";
import { SignUp } from "../components/form/SignUp";
import { auth } from "../services/auth.service";

export const Home = () => {
  const [currentModal, setModal] = useState("login");

  const token = auth.get();

  useEffect(() => {
    document.title = "Split Specs";
  }, []);

  const onNavClick = useCallback((e) => {
    setModal(e.target.id);
  }, []);

  const showModal = (modal) => {
    const style = "text-blue-600 font-semibold";

    const [id, text] =
      modal === "login"
        ? ["register", "I am new here"]
        : ["login", "Back to login"];

    return (
      <button id={id} className={style} onClick={onNavClick}>
        {text}
      </button>
    );
  };

  const history = useHistory();

  if (token) {
    console.log(`HOME HAS TOKEN`)
    history.push("/projects");
  }

  return (
    <section>
      <div className="fixed inset-0 bg-gray-900 bg-opacity-60">
        <div className="min-w-full min-h-full flex items-center justify-center px-4">
          <div className="max-w-md w-full bg-white rounded-md p-6 shadow-2xl">
            <div className="mb-6">
              <h2 className="text-center text-3xl font-extrabold text-gray-700 mt-4">
                Split specs
              </h2>
              <p className="text-center text-xs text-gray-600 mt-1">
                {showModal(currentModal)}
              </p>
            </div>
            {currentModal === "login" ? (
              <SignIn history={history} />
            ) : (
              <SignUp history={history} />
            )}
          </div>
        </div>
      </div>
    </section>
  );
};
