"use client";

import { useMyHandle } from "@/hooks/useMyHandle";
import { WorkHistory } from "../u/[handle]/WorkHistory";
import { useTranslation } from "@/hooks/useTranslation";
import CircularProgress from "@mui/material/CircularProgress";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import OfficialEmails from "@/components/OfficialEmails";
import Divider from "@mui/material/Divider";

export default function MyProfilePage() {
  const { myHandle, isLoading: isLoadingHandle } = useMyHandle();
  const { t } = useTranslation();

  if (isLoadingHandle) {
    return (
      <AuthenticatedLayout>
        <Box sx={{ display: "flex", justifyContent: "center", mt: 4 }}>
          <CircularProgress />
        </Box>
      </AuthenticatedLayout>
    );
  }

  if (!myHandle) {
    return null; // or handle unauthorized state
  }

  return (
    <AuthenticatedLayout>
      <Box sx={{ maxWidth: 800, mx: "auto", mt: 4, px: 2 }}>
        <Box sx={{ mb: 4 }}>
          <Typography variant="h4" sx={{ mb: 1 }}>
            {myHandle}
          </Typography>
          <Typography variant="subtitle1" color="text.secondary">
            @{myHandle}
          </Typography>
        </Box>

        <Box sx={{ mt: 4, mb: 4 }}>
          <OfficialEmails />
        </Box>

        <Divider sx={{ my: 6 }} />

        <Box sx={{ mt: 4 }}>
          <WorkHistory userHandle={myHandle} canEdit={true} />
        </Box>
      </Box>
    </AuthenticatedLayout>
  );
}
