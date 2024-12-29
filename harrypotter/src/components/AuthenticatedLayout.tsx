"use client";

import { useState } from "react";
import Box from "@mui/material/Box";
import { useTheme } from "@mui/material/styles";
import Header from "./Header";
import Sidebar from "./Sidebar";

export default function AuthenticatedLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const theme = useTheme();

  const toggleSidebar = () => {
    setSidebarOpen(!sidebarOpen);
  };

  return (
    <Box sx={{ display: "flex", minHeight: "100vh", flexDirection: "column" }}>
      <Header onMenuClick={toggleSidebar} />
      <Sidebar open={sidebarOpen} />
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          p: 3,
          width: "100%",
          ml: sidebarOpen ? "240px" : theme.spacing(9),
          mt: "64px", // Height of the header
          transition: theme.transitions.create(["margin"], {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.enteringScreen,
          }),
        }}
      >
        {children}
      </Box>
    </Box>
  );
}
