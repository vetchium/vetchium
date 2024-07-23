import { useState } from "react";
import Home from "./pages/Home";
import SignIn from "./pages/SignIn";

function App() {
  const [signedIn, setSignedIn] = useState(false);

  return signedIn ? (
    <Home
      onSignOut={() => {
        setSignedIn(false);
        window.location.href = "/vetchi";
      }}
    />
  ) : (
    <SignIn
      onSignIn={() => {
        setSignedIn(true);
      }}
    />
  );
}

export default App;
