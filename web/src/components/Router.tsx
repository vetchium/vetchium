import { Routes, Route } from "react-router-dom";
import Dashboard from "../pages/Dashboard";

function Router() {
  return (
    <Routes>
      <Route path="/home" element={<Dashboard />} />
    </Routes>
  );
}

export default Router;
