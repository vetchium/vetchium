import { Route, Routes } from "react-router-dom";
import Applications from "../pages/Applications";
import Candidacy from "../pages/Candidacy";
import Candidates from "../pages/Candidates";
import CreateInterview from "../pages/CreateInterview";
import CreateOpening from "../pages/CreateOpening";
import Departments from "../pages/Departments";
import Interview from "../pages/Interview";
import Interviews from "../pages/Interviews";
import LocationSelector from "../pages/Locations";
import Openings from "../pages/Openings";
import Users from "../pages/Users";

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
      <Route
        path="/create-interview/:candidacy_id"
        element={<CreateInterview />}
      />
      <Route path="/interview/:interview_id" element={<Interview />} />
      <Route path="/interviews" element={<Interviews />} />
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
