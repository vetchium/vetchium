"use client";

import { useState } from "react";
import Box from "@mui/material/Box";
import Header from "./Header";
import Sidebar from "./Sidebar";
import Footer from "./Footer";

export default function AuthenticatedLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const [sidebarOpen, setSidebarOpen] = useState(true);

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
          width: { xs: `calc(100% - ${sidebarOpen ? "240px" : "72px"})` },
          ml: sidebarOpen ? "240px" : "72px",
          mt: "64px", // Height of the header
          transition: (theme) =>
            theme.transitions.create(["margin", "width"], {
              easing: theme.transitions.easing.sharp,
              duration: theme.transitions.duration.enteringScreen,
            }),
        }}
      >
        {children}
      </Box>
      <Footer />
    </Box>
  );
}
