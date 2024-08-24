import { Route, Routes } from "react-router-dom";
import Applications from "../pages/Applications";
import CreateOpening from "../pages/CreateOpening";
import Departments from "../pages/Departments";
import LocationSelector from "../pages/Locations";
import Openings from "../pages/Openings";
import Users from "../pages/Users";
import Candidates from "../pages/Candidates";
import Candidacy from "../pages/Candidacy";

function Router() {
  return (
    <Routes>
      <Route
        path="/openings/:opening_id/applications"
        element={<Applications />}
      />
      <Route path="/openings" element={<Openings />} />
      <Route path="/candidacy/:candidacy_id" element={<Candidacy />} />
      <Route path="/candidates" element={<Candidates />} />
      <Route path="/create-opening" element={<CreateOpening />} />
      <Route path="/org-settings" element={<div>Org Settings</div>} />
      <Route path="/org-settings/locations" element={<LocationSelector />} />
      <Route path="/org-settings/departments" element={<Departments />} />
      <Route path="/org-settings/users" element={<Users />} />
      <Route path="/account-settings" element={<div>Account Settings</div>} />
    </Routes>
  );
}

export default Router;
