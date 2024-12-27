"use client";

import { Box, Typography } from "@mui/material";
import { useTranslation } from "@/hooks/useTranslation";
import DashboardLayout from "@/components/DashboardLayout";

export default function Dashboard() {
  const { t } = useTranslation();

  return (
    <DashboardLayout>
      <Box sx={{ width: "100%" }}>
        <Typography variant="h4" gutterBottom>
          {t("dashboard.welcome")}
        </Typography>
      </Box>
    </DashboardLayout>
  );
}
