import Header from "./Header";
import Footer from "./atoms/Footer";

export const Layout = ({ children }) => {
  return (
    <div className="h-screen w-screen flex flex-col">
      <Header title="Split specs" />

      <main className="flex-1">{children}</main>

      <Footer title="&copy; Shevtsov Oleksandr 2022" />
    </div>
  );
};
