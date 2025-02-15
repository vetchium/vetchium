"use client";

import { useMyHandle } from "@/hooks/useMyHandle";
import { useParams, useRouter } from "next/navigation";
import { WorkHistory } from "./WorkHistory";
import { useTranslation } from "@/hooks/useTranslation";
import CircularProgress from "@mui/material/CircularProgress";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import AuthenticatedLayout from "@/components/AuthenticatedLayout";

export default function ProfilePage() {
  const params = useParams();
  const router = useRouter();
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
        <Box
          sx={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "flex-start",
            mb: 4,
          }}
        >
          <Box>
            <Typography variant="h4" sx={{ mb: 1 }}>
              {userHandle}
            </Typography>
            <Typography variant="subtitle1" color="text.secondary">
              @{userHandle}
            </Typography>
          </Box>
          {isOwnProfile && (
            <Button
              variant="contained"
              color="primary"
              onClick={() => router.push("/my-profile")}
            >
              {t("profile.editMyProfile")}
            </Button>
          )}
        </Box>

        <Box sx={{ mt: 4 }}>
          <WorkHistory userHandle={userHandle} canEdit={false} />
        </Box>
      </Box>
    </AuthenticatedLayout>
  );
}
