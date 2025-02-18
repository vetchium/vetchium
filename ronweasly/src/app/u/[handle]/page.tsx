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
import Bio from "@/components/Bio";
import ProfilePicture from "@/components/ProfilePicture";
import { useProfile } from "@/hooks/useProfile";
import { config } from "@/config";
import Alert from "@mui/material/Alert";

export default function ProfilePage() {
  const params = useParams();
  const router = useRouter();
  const userHandle = params.handle as string;
  const { myHandle, isLoading: isLoadingHandle } = useMyHandle();
  const { t } = useTranslation();
  const isOwnProfile = myHandle === userHandle;
  const { bio, isLoading: isLoadingBio, error } = useProfile(userHandle);

  if (isLoadingHandle || isLoadingBio) {
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
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error.message}
          </Alert>
        )}

        <Box
          sx={{
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            mb: 4,
          }}
        >
          <ProfilePicture
            imageUrl={`${config.API_SERVER_PREFIX}/hub/profile-picture/${userHandle}`}
            size={150}
          />

          {isOwnProfile && (
            <Button
              variant="contained"
              color="primary"
              onClick={() => router.push("/my-profile")}
              sx={{ mt: 2 }}
            >
              {t("profile.editMyProfile")}
            </Button>
          )}
        </Box>

        {bio && (
          <Box sx={{ mb: 4 }}>
            <Bio bio={bio} isLoading={false} />
          </Box>
        )}

        <Box sx={{ mt: 4 }}>
          <WorkHistory userHandle={userHandle} canEdit={false} />
        </Box>
      </Box>
    </AuthenticatedLayout>
  );
}
