import { Routes, Route } from "react-router-dom";
import Dashboard from "../pages/Dashboard";
import MyProfile from "../pages/MyProfile";
import AddWorkHistory from "../pages/AddWorkHistory";

function Router() {
  return (
    <Routes>
      <Route path="/home" element={<Dashboard />} />
      <Route path="/my-profile" element={<MyProfile />} />
      <Route path="/add-work-history" element={<AddWorkHistory />} />
    </Routes>
  );
}

export default Router;
