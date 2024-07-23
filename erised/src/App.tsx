import { useState, useEffect } from "react";
import Cookies from "js-cookie";
import Home from "./pages/Home";
import SignIn from "./pages/SignIn";

function App() {
  const [signedIn, setSignedIn] = useState(false);

  useEffect(() => {
    // Check for a valid cookie on app load
    const cookie = Cookies.get("signedIn");
    setSignedIn(cookie === "true");
  }, []);

  return signedIn ? (
    <Home
      onSignOut={() => {
        setSignedIn(false);
        Cookies.remove("signedIn"); // Clear the cookie on sign out
        window.location.href = "/vetchi";
      }}
    />
  ) : (
    <SignIn
      onSignIn={() => {
        setSignedIn(true);
        Cookies.set("signedIn", "true"); // Set the cookie on sign in
      }}
    />
  );
}

export default App;
