import { Route, Routes } from "react-router-dom";

function Router() {
  return (
    <Routes>
      <Route path="/one" element={<div>One</div>} />
      <Route path="/two" element={<div>Two</div>} />
      <Route path="/three" element={<div>Three</div>} />
    </Routes>
  );
}

export default Router;
