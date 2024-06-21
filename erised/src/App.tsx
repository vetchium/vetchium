import { useState } from "react";
import Home from "./pages/Home";
import SignIn from "./pages/SignIn";

function App() {
  const [signedIn, setSignedIn] = useState(false);

  if (signedIn) {
    return <Home onSignOut={() => setSignedIn(false)} />;
  }

  return <SignIn onSignIn={() => setSignedIn(true)} />;
}

export default App;
