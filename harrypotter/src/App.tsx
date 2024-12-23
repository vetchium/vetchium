import { useState, useEffect } from "react";
import Cookies from "js-cookie";
import Home from "./pages/Home";
import SignIn from "./pages/SignIn";
import {
  OrgUserRoles,
  OrgUserRole,
  LoginRequest,
} from "@psankar/vetchi-typespec";

function App() {
  const [signedIn, setSignedIn] = useState(false);
  const [userRole, setUserRole] = useState<OrgUserRole>(
    OrgUserRoles.OPENINGS_VIEWER
  );

  useEffect(() => {
    const cookie = Cookies.get("signedIn");
    setSignedIn(cookie === "true");
  }, []);

  return signedIn ? (
    <Home
      onSignOut={() => {
        setSignedIn(false);
        setUserRole(OrgUserRoles.OPENINGS_VIEWER);
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
