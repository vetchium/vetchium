import { Route, Routes } from "react-router-dom";
import Openings from "../pages/Openings";

function Router() {
  return (
    <Routes>
      <Route path="/openings" element={<Openings />} />
      <Route path="/org-settings" element={<div>Org Settings</div>} />
      <Route path="/account-settings" element={<div>Account Settings</div>} />
    </Routes>
  );
}

export default Router;
