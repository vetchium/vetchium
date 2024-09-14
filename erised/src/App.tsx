import { useState, useEffect } from "react";
import Cookies from "js-cookie";
import Home from "./pages/Home";
import SignIn from "./pages/SignIn";

function App() {
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
    <SignIn
      onSignIn={() => {
        setSignedIn(true);
        Cookies.set("signedIn", "true");
      }}
    />
  );
}

export default App;
