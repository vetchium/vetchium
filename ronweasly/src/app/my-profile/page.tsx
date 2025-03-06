"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import Bio from "@/components/Bio";
import OfficialEmails from "@/components/OfficialEmails";
import ProfilePicture from "@/components/ProfilePicture";
import { config } from "@/config";
import { useMyHandle } from "@/hooks/useMyHandle";
import { useProfile } from "@/hooks/useProfile";
import { useTranslation } from "@/hooks/useTranslation";
import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import CircularProgress from "@mui/material/CircularProgress";
import Divider from "@mui/material/Divider";
import { useState } from "react";
import { WorkHistory } from "../u/[handle]/WorkHistory";

export default function MyProfilePage() {
  const { myHandle, isLoading: isLoadingHandle } = useMyHandle();
  const { t } = useTranslation();
  const [timestamp, setTimestamp] = useState(Date.now());
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
          imageUrl={`${config.API_SERVER_PREFIX}/hub/profile-picture/${myHandle}?t=${timestamp}`}
          onImageSelect={uploadProfilePicture}
          onRemove={() => setTimestamp(Date.now())}
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
