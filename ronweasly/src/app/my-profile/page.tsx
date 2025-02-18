"use client";

import { useMyHandle } from "@/hooks/useMyHandle";
import { useProfile } from "@/hooks/useProfile";
import { WorkHistory } from "../u/[handle]/WorkHistory";
import { useTranslation } from "@/hooks/useTranslation";
import CircularProgress from "@mui/material/CircularProgress";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import OfficialEmails from "@/components/OfficialEmails";
import Divider from "@mui/material/Divider";
import Bio from "@/components/Bio";
import ProfilePicture from "@/components/ProfilePicture";
import Alert from "@mui/material/Alert";
import { config } from "@/config";

export default function MyProfilePage() {
  const { myHandle, isLoading: isLoadingHandle } = useMyHandle();
  const { t } = useTranslation();
  const {
    bio,
    isLoading: isLoadingBio,
    error,
    isSaving,
    updateBio,
    uploadProfilePicture,
  } = useProfile(myHandle ?? "");

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

  const isLoading = isLoadingBio || isLoadingHandle;

  return (
    <AuthenticatedLayout>
      <Box sx={{ maxWidth: 800, mx: "auto", mt: 4, px: 2 }}>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error.message}
          </Alert>
        )}

        <ProfilePicture
          imageUrl={`${config.API_SERVER_PREFIX}/hub/profile-picture/${myHandle}`}
          onImageSelect={uploadProfilePicture}
          isLoading={isSaving}
        />

        {bio && (
          <Box sx={{ mb: 4 }}>
            <Bio bio={bio} onSave={updateBio} isLoading={isSaving} />
          </Box>
        )}

        <Divider sx={{ my: 6 }} />

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
