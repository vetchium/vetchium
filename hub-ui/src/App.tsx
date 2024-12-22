import React from "react";
import { Routes, Route, Navigate } from "react-router-dom";
import { Layout } from "antd";
import MainLayout from "./layouts/MainLayout";
import SignIn from "./pages/SignIn";
import Dashboard from "./pages/Dashboard";
import Openings from "./pages/Openings";
import Applications from "./pages/Applications";
import Candidacies from "./pages/Candidacies";
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
          path="/openings"
          element={isAuthenticated ? <Openings /> : <Navigate to="/signin" />}
        />
        <Route
          path="/applications"
          element={
            isAuthenticated ? <Applications /> : <Navigate to="/signin" />
          }
        />
        <Route
          path="/candidacies"
          element={
            isAuthenticated ? <Candidacies /> : <Navigate to="/signin" />
          }
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
