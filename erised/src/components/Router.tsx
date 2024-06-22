import { Route, Routes } from "react-router-dom";
import Openings from "../pages/Openings";
import CreateOpening from "../pages/CreateOpening";

function Router() {
  return (
    <Routes>
      <Route path="/openings" element={<Openings />} />
      <Route path="/create-opening" element={<CreateOpening />} />
      <Route path="/org-settings" element={<div>Org Settings</div>} />
      <Route path="/account-settings" element={<div>Account Settings</div>} />
    </Routes>
  );
}

export default Router;
