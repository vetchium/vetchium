import React, { useState, useEffect } from "react";
import Cookies from "js-cookie";
import Home from "./pages/Home";
import Login from "./pages/Login";

const App: React.FC = () => {
  const [signedIn, setSignedIn] = useState(false);

  useEffect(() => {
    const cookie = Cookies.get("signedIn");
    setSignedIn(cookie === "true");
  }, []);

  return signedIn ? (
    <Home
      onSignOut={() => {
        setSignedIn(false);
        Cookies.remove("signedIn");
        window.location.href = "/";
      }}
    />
  ) : (
    <Login
      onSignIn={() => {
        setSignedIn(true);
        Cookies.set("signedIn", "true");
      }}
    />
  );
};

export default App;
