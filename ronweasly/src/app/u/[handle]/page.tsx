"use client";

import { useMyHandle } from "@/hooks/useMyHandle";
import { useParams } from "next/navigation";
import { WorkHistory } from "./WorkHistory";
import { useTranslation } from "@/hooks/useTranslation";
import CircularProgress from "@mui/material/CircularProgress";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import AuthenticatedLayout from "@/components/AuthenticatedLayout";

export default function ProfilePage() {
  const params = useParams();
  const userHandle = params.handle as string;
  const { myHandle, isLoading: isLoadingHandle } = useMyHandle();
  const { t } = useTranslation();
  const isOwnProfile = myHandle === userHandle;

  if (isLoadingHandle) {
    return (
      <AuthenticatedLayout>
        <Box sx={{ display: "flex", justifyContent: "center", mt: 4 }}>
          <CircularProgress />
        </Box>
      </AuthenticatedLayout>
    );
  }

  return (
    <AuthenticatedLayout>
      <Box sx={{ maxWidth: 800, mx: "auto", mt: 4, px: 2 }}>
        <Typography variant="h4" gutterBottom>
          {isOwnProfile ? t("profile.myProfile") : `${userHandle}'s Profile`}
        </Typography>

        <Box sx={{ mt: 4 }}>
          <WorkHistory userHandle={userHandle} canEdit={isOwnProfile} />
        </Box>
      </Box>
    </AuthenticatedLayout>
  );
}
