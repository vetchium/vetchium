"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import { useAuth } from "@/hooks/useAuth";
import { useTranslation } from "@/hooks/useTranslation";
import "@fontsource/roboto/300.css";
import "@fontsource/roboto/400.css";
import "@fontsource/roboto/500.css";
import "@fontsource/roboto/700.css";
import Box from "@mui/material/Box";
import Grid from "@mui/material/Grid";
import Paper from "@mui/material/Paper";
import Typography from "@mui/material/Typography";
import Link from "next/link";

export default function DashboardPage() {
  const { t } = useTranslation();
  useAuth(); // Check authentication and redirect if not authenticated

  return (
    <AuthenticatedLayout>
      <Box sx={{ flexGrow: 1 }}>
        <Typography variant="h4" gutterBottom>
          {t("dashboard.title")}
        </Typography>
        <Grid container spacing={3}>
          <Grid item xs={12} md={6} lg={4}>
            <Link href="/my-applications" style={{ textDecoration: "none" }}>
              <Paper
                sx={{
                  p: 3,
                  display: "flex",
                  flexDirection: "column",
                  height: 240,
                  cursor: "pointer",
                  "&:hover": {
                    bgcolor: "action.hover",
                  },
                }}
              >
                <Typography variant="h6" gutterBottom>
                  {t("navigation.myApplications")}
                </Typography>
                {/* Add content here */}
              </Paper>
            </Link>
          </Grid>
          <Grid item xs={12} md={6} lg={4}>
            <Paper
              sx={{
                p: 3,
                display: "flex",
                flexDirection: "column",
                height: 240,
              }}
            >
              <Typography variant="h6" gutterBottom>
                {t("dashboard.activeOpenings")}
              </Typography>
              {/* Add content here */}
            </Paper>
          </Grid>
          <Grid item xs={12} md={6} lg={4}>
            <Paper
              sx={{
                p: 3,
                display: "flex",
                flexDirection: "column",
                height: 240,
              }}
            >
              <Typography variant="h6" gutterBottom>
                {t("dashboard.upcomingInterviews")}
              </Typography>
              {/* Add content here */}
            </Paper>
          </Grid>
        </Grid>
      </Box>
    </AuthenticatedLayout>
  );
}
