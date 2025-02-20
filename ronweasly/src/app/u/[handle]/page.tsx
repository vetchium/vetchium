"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import Bio from "@/components/Bio";
import ProfilePicture from "@/components/ProfilePicture";
import { config } from "@/config";
import { useMyHandle } from "@/hooks/useMyHandle";
import { useProfile } from "@/hooks/useProfile";
import { useTranslation } from "@/hooks/useTranslation";
import PersonAddIcon from "@mui/icons-material/PersonAdd";
import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import CircularProgress from "@mui/material/CircularProgress";
import Container from "@mui/material/Container";
import Divider from "@mui/material/Divider";
import Paper from "@mui/material/Paper";
import Typography from "@mui/material/Typography";
import { useParams, useRouter } from "next/navigation";
import { WorkHistory } from "./WorkHistory";

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

  const handleAddColleague = () => {
    // TODO: Implement the add colleague functionality
    console.log("Add colleague clicked");
  };

  return (
    <AuthenticatedLayout>
      <Container maxWidth="md">
        <Box sx={{ py: 4 }}>
          {error && (
            <Alert severity="error" sx={{ mb: 3 }}>
              {error.message}
            </Alert>
          )}

          <Box sx={{ display: "flex", gap: 3 }}>
            {/* Main content */}
            <Paper
              elevation={0}
              sx={{
                p: 4,
                borderRadius: 2,
                bgcolor: "background.default",
                flex: 1,
              }}
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

            {/* Actions sidebar */}
            {bio?.colleaguable && (
              <Box sx={{ width: 280, flexShrink: 0 }}>
                <Card elevation={0}>
                  <CardContent>
                    <Typography variant="h6" sx={{ mb: 2 }}>
                      {t("profile.actions")}
                    </Typography>
                    <Button
                      variant="outlined"
                      startIcon={<PersonAddIcon />}
                      onClick={handleAddColleague}
                      fullWidth
                    >
                      {t("profile.addAsColleague")}
                    </Button>
                  </CardContent>
                </Card>
              </Box>
            )}
          </Box>
        </Box>
      </Container>
    </AuthenticatedLayout>
  );
}
