import { Routes, Route } from "react-router-dom";
import Dashboard from "../pages/Dashboard";
import MyProfile from "../pages/MyProfile";
function Router() {
  return (
    <Routes>
      <Route path="/home" element={<Dashboard />} />
      <Route path="/my-profile" element={<MyProfile />} />
    </Routes>
  );
}

export default Router;
