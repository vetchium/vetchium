"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import Bio from "@/components/Bio";
import ProfilePicture from "@/components/ProfilePicture";
import { config } from "@/config";
import { useMyHandle } from "@/hooks/useMyHandle";
import { useProfile } from "@/hooks/useProfile";
import { useTranslation } from "@/hooks/useTranslation";
import { useColleagues } from "@/hooks/useColleagues";
import PersonAddIcon from "@mui/icons-material/PersonAdd";
import VerifiedIcon from "@mui/icons-material/Verified";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import CancelIcon from "@mui/icons-material/Cancel";
import BlockIcon from "@mui/icons-material/Block";
import LinkOffIcon from "@mui/icons-material/LinkOff";
import PendingIcon from "@mui/icons-material/Pending";
import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Chip from "@mui/material/Chip";
import CircularProgress from "@mui/material/CircularProgress";
import Container from "@mui/material/Container";
import Divider from "@mui/material/Divider";
import Paper from "@mui/material/Paper";
import Typography from "@mui/material/Typography";
import Stack from "@mui/material/Stack";
import { useParams, useRouter } from "next/navigation";
import { WorkHistory } from "./WorkHistory";
import { useState } from "react";

export default function ProfilePage() {
  const params = useParams();
  const router = useRouter();
  const userHandle = params.handle as string;
  const { myHandle, isLoading: isLoadingHandle } = useMyHandle();
  const { t } = useTranslation();
  const isOwnProfile = myHandle === userHandle;
  const {
    bio,
    isLoading: isLoadingBio,
    error,
    refetch,
  } = useProfile(userHandle);
  const { connectColleague, isConnecting } = useColleagues();
  const [connectionError, setConnectionError] = useState<string | null>(null);

  if (isLoadingHandle || isLoadingBio) {
    return (
      <AuthenticatedLayout>
        <Box sx={{ display: "flex", justifyContent: "center", mt: 4 }}>
          <CircularProgress />
        </Box>
      </AuthenticatedLayout>
    );
  }

  const handleAddColleague = async () => {
    if (!bio) return;

    setConnectionError(null);

    try {
      await connectColleague(bio.handle);
      // Refetch the profile to get updated connection state
      await refetch();
    } catch (error) {
      setConnectionError(
        error instanceof Error
          ? t(error.message)
          : t("profile.error.connectionFailed")
      );
    }
  };

  const handleApproveRequest = () => {
    console.log("Approve request clicked");
  };

  const handleDeclineRequest = () => {
    console.log("Decline request clicked");
  };

  const handleUnlinkConnection = () => {
    console.log("Unlink connection clicked");
  };

  const renderConnectionActions = () => {
    if (isOwnProfile || !bio) return null;

    const state = bio.colleague_connection_state;

    switch (state) {
      case "CAN_SEND_REQUEST":
        return (
          <Stack spacing={2}>
            <Button
              variant="outlined"
              startIcon={
                isConnecting ? (
                  <CircularProgress size={20} />
                ) : (
                  <PersonAddIcon />
                )
              }
              onClick={handleAddColleague}
              disabled={isConnecting}
              fullWidth
            >
              {isConnecting ? t("common.loading") : t("profile.addAsColleague")}
            </Button>
            {connectionError && (
              <Alert severity="error">{connectionError}</Alert>
            )}
          </Stack>
        );

      case "CANNOT_SEND_REQUEST":
        return (
          <Alert severity="info" sx={{ mb: 2 }}>
            {t("profile.cannotAddAsColleague")}
          </Alert>
        );

      case "REQUEST_SENT_PENDING":
        return (
          <Alert severity="info" icon={<PendingIcon />} sx={{ mb: 2 }}>
            {t("profile.requestPending")}
          </Alert>
        );

      case "REQUEST_RECEIVED_PENDING":
        return (
          <Stack spacing={2}>
            <Alert severity="info" sx={{ mb: 2 }}>
              {t("profile.receivedColleagueRequest")}
            </Alert>
            <Button
              variant="contained"
              color="success"
              startIcon={<CheckCircleIcon />}
              onClick={handleApproveRequest}
              fullWidth
            >
              {t("profile.approveRequest")}
            </Button>
            <Button
              variant="outlined"
              color="error"
              startIcon={<CancelIcon />}
              onClick={handleDeclineRequest}
              fullWidth
            >
              {t("profile.declineRequest")}
            </Button>
          </Stack>
        );

      case "CONNECTED":
        return (
          <Stack spacing={2}>
            <Alert severity="success" icon={<VerifiedIcon />} sx={{ mb: 2 }}>
              {t("profile.connectedAsColleagues")}
            </Alert>
            <Button
              variant="outlined"
              color="warning"
              startIcon={<LinkOffIcon />}
              onClick={handleUnlinkConnection}
              fullWidth
            >
              {t("profile.unlinkConnection")}
            </Button>
          </Stack>
        );

      case "REJECTED_BY_ME":
        return (
          <Alert severity="warning" icon={<BlockIcon />} sx={{ mb: 2 }}>
            {t("profile.youRejectedTheirRequest")}
          </Alert>
        );

      case "REJECTED_BY_THEM":
        return (
          <Alert severity="warning" icon={<BlockIcon />} sx={{ mb: 2 }}>
            {t("profile.theyRejectedYourRequest")}
          </Alert>
        );

      case "UNLINKED_BY_ME":
        return (
          <Alert severity="info" icon={<LinkOffIcon />} sx={{ mb: 2 }}>
            {t("profile.youUnlinkedConnection")}
          </Alert>
        );

      case "UNLINKED_BY_THEM":
        return (
          <Alert severity="info" icon={<LinkOffIcon />} sx={{ mb: 2 }}>
            {t("profile.theyUnlinkedConnection")}
          </Alert>
        );

      default:
        return null;
    }
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
                    <Box
                      sx={{
                        display: "flex",
                        alignItems: "center",
                        gap: 2,
                        mb: 2,
                      }}
                    >
                      <Bio bio={bio} isLoading={false} />
                      {bio.colleague_connection_state === "CONNECTED" && (
                        <Chip
                          icon={<VerifiedIcon />}
                          label={t("profile.verifiedColleague")}
                          color="primary"
                          variant="outlined"
                          size="small"
                        />
                      )}
                    </Box>
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
            {!isOwnProfile && (
              <Box sx={{ width: 280, flexShrink: 0 }}>
                <Card elevation={0}>
                  <CardContent>
                    <Typography variant="h6" sx={{ mb: 2 }}>
                      {t("profile.actions")}
                    </Typography>
                    {renderConnectionActions()}
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
