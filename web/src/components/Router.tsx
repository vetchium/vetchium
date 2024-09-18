import { Routes, Route } from "react-router-dom";
import Dashboard from "../pages/Dashboard";
import MyProfile from "../pages/MyProfile";
import AddWorkHistory from "../pages/AddWorkHistory";
import UserProfile from "../pages/UserProfile";

function Router() {
  return (
    <Routes>
      <Route path="/home" element={<Dashboard />} />
      <Route path="/my-profile" element={<MyProfile />} />
      <Route path="/add-work-history" element={<AddWorkHistory />} />
      <Route path="/u/:id" element={<UserProfile />} />
    </Routes>
  );
}

export default Router;
