"use client";

import { useAuth } from "@/hooks/useAuth";
import Box from "@mui/material/Box";
import { useTheme } from "@mui/material/styles";
import { useState } from "react";
import Header from "./Header";
import Sidebar from "./Sidebar";

export default function AuthenticatedLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const theme = useTheme();
  useAuth(); // Check authentication and redirect if not authenticated

  const toggleSidebar = () => {
    setSidebarOpen(!sidebarOpen);
  };

  return (
    <Box sx={{ display: "flex", minHeight: "100vh", flexDirection: "column" }}>
      <Header onMenuClick={toggleSidebar} isSidebarOpen={sidebarOpen} />
      <Sidebar open={sidebarOpen} />
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          p: 3,
          width: {
            xs: `calc(100% - ${sidebarOpen ? "240px" : theme.spacing(9)})`,
          },
          ml: sidebarOpen ? "240px" : theme.spacing(9),
          mt: "64px", // Height of the header
          transition: theme.transitions.create(["margin", "width"], {
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
