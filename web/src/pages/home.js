import { useEffect } from "react";
import { Redirect } from "react-router-dom";
import SignInForm from "../components/form/SignIn";
import { auth } from "../services/auth.service";

export const Home = () => {
  const token = auth.get();

  useEffect(() => {
    document.title = "Split Specs";
  }, []);

  return (
    <section>{token ? <Redirect to="/projects" /> : <SignInForm />}</section>
  );
};
