import React, { useState, useEffect } from "react";
import Cookies from "js-cookie";
import Home from "./pages/Home";
import Login from "./pages/Login";
import { ConfigProvider } from "antd";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";

const App: React.FC = () => {
  const [loggedIn, setLoggedIn] = useState(false);

  useEffect(() => {
    const cookie = Cookies.get("signedIn");
    setLoggedIn(cookie === "true");
  }, []);

  return (
    <ConfigProvider
      theme={{
        token: {
          colorPrimary: "#00B96B", // A shade of green
        },
      }}
    >
      {loggedIn ? (
        <Home
          onLogOut={() => {
            setLoggedIn(false);
            Cookies.remove("loggedIn");
            window.location.href = "/";
          }}
        />
      ) : (
        <Login
          onLogIn={() => {
            setLoggedIn(true);
            Cookies.set("loggedIn", "true");
          }}
        />
      )}
    </ConfigProvider>
  );
};

export default App;
