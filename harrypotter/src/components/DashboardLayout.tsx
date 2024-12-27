"use client";

import { Box } from "@mui/material";
import Header from "./Header";
import Sidebar from "./Sidebar";
import Footer from "./Footer";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <Box sx={{ display: "flex", minHeight: "100vh", flexDirection: "column" }}>
      <Header />
      <Box sx={{ display: "flex", flex: 1 }}>
        <Sidebar />
        <Box
          component="main"
          sx={{
            flexGrow: 1,
            p: 3,
            mt: 8,
            display: "flex",
            flexDirection: "column",
          }}
        >
          {children}
          <Footer />
        </Box>
      </Box>
    </Box>
  );
}
