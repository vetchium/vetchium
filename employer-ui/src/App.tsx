import React from "react";
import { Routes, Route, Navigate } from "react-router-dom";
import { Layout } from "antd";
import MainLayout from "./layouts/MainLayout";
import SignIn from "./pages/SignIn";
import Dashboard from "./pages/Dashboard";
import Locations from "./pages/Locations";
import Departments from "./pages/Departments";
import Openings from "./pages/Openings";
import CreateOpening from "./pages/CreateOpening";
import Applications from "./pages/Applications";
import OrgUsers from "./pages/OrgUsers";
import Settings from "./pages/Settings";
import { useAuth } from "./hooks/useAuth";

const App: React.FC = () => {
  const { isAuthenticated } = useAuth();

  return (
    <Routes>
      <Route
        path="/signin"
        element={!isAuthenticated ? <SignIn /> : <Navigate to="/" />}
      />

      <Route element={<MainLayout />}>
        <Route
          path="/"
          element={isAuthenticated ? <Dashboard /> : <Navigate to="/signin" />}
        />
        <Route
          path="/locations"
          element={isAuthenticated ? <Locations /> : <Navigate to="/signin" />}
        />
        <Route
          path="/departments"
          element={
            isAuthenticated ? <Departments /> : <Navigate to="/signin" />
          }
        />
        <Route
          path="/openings"
          element={isAuthenticated ? <Openings /> : <Navigate to="/signin" />}
        />
        <Route
          path="/openings/create"
          element={
            isAuthenticated ? <CreateOpening /> : <Navigate to="/signin" />
          }
        />
        <Route
          path="/applications"
          element={
            isAuthenticated ? <Applications /> : <Navigate to="/signin" />
          }
        />
        <Route
          path="/org-users"
          element={isAuthenticated ? <OrgUsers /> : <Navigate to="/signin" />}
        />
        <Route
          path="/settings"
          element={isAuthenticated ? <Settings /> : <Navigate to="/signin" />}
        />
      </Route>

      <Route path="*" element={<Navigate to="/" />} />
    </Routes>
  );
};

export default App;
