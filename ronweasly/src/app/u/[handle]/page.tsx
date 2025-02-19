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
import Paper from "@mui/material/Paper";
import Container from "@mui/material/Container";
import Divider from "@mui/material/Divider";

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
      <Container maxWidth="md">
        <Box sx={{ py: 4 }}>
          {error && (
            <Alert severity="error" sx={{ mb: 3 }}>
              {error.message}
            </Alert>
          )}

          <Paper
            elevation={0}
            sx={{ p: 4, borderRadius: 2, bgcolor: "background.default" }}
          >
            <Box
              sx={{
                display: "flex",
                gap: 4,
                mb: 6,
              }}
            >
              {/* Left column - Profile Picture */}
              <Box
                sx={{
                  display: "flex",
                  flexDirection: "column",
                  alignItems: "center",
                  flexShrink: 0,
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

              {/* Right column - Bio */}
              {bio && (
                <Box sx={{ flex: 1 }}>
                  <Bio bio={bio} isLoading={false} />
                </Box>
              )}
            </Box>

            <Divider sx={{ mb: 4 }} />

            {/* Work History section */}
            <Box>
              <WorkHistory userHandle={userHandle} canEdit={false} />
            </Box>
          </Paper>
        </Box>
      </Container>
    </AuthenticatedLayout>
  );
}
