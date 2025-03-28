"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import Bio from "@/components/Bio";
import { Education } from "@/components/Education";
import ProfilePicture from "@/components/ProfilePicture";
import { config } from "@/config";
import { useColleagues } from "@/hooks/useColleagues";
import { useMyHandle } from "@/hooks/useMyHandle";
import { useProfile } from "@/hooks/useProfile";
import { useTranslation } from "@/hooks/useTranslation";
import BlockIcon from "@mui/icons-material/Block";
import CancelIcon from "@mui/icons-material/Cancel";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import LinkOffIcon from "@mui/icons-material/LinkOff";
import MoreVertIcon from "@mui/icons-material/MoreVert";
import PendingIcon from "@mui/icons-material/Pending";
import PersonAddIcon from "@mui/icons-material/PersonAdd";
import VerifiedIcon from "@mui/icons-material/Verified";
import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import CircularProgress from "@mui/material/CircularProgress";
import Container from "@mui/material/Container";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogTitle from "@mui/material/DialogTitle";
import Divider from "@mui/material/Divider";
import IconButton from "@mui/material/IconButton";
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";
import Paper from "@mui/material/Paper";
import Stack from "@mui/material/Stack";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import { useParams, useRouter } from "next/navigation";
import { useState } from "react";
import { WorkHistory } from "./WorkHistory";

export default function UserProfilePage() {
  const { t } = useTranslation();
  const params = useParams();
  const router = useRouter();

  if (!params?.handle) {
    return (
      <Box sx={{ p: 3 }}>
        <Typography color="error">{t("common.error.invalidParams")}</Typography>
        <Button
          variant="contained"
          onClick={() => router.back()}
          sx={{ mt: 2 }}
        >
          {t("common.back")}
        </Button>
      </Box>
    );
  }

  const userHandle = params.handle as string;
  const { myHandle, isLoading: isLoadingHandle } = useMyHandle();
  const isOwnProfile = myHandle === userHandle;
  const {
    bio,
    isLoading: isLoadingBio,
    error,
    refetch,
  } = useProfile(userHandle);
  const {
    connectColleague,
    approveColleague,
    rejectColleague,
    unlinkColleague,
    isConnecting,
    isApproving,
    isRejecting,
    isUnlinking,
  } = useColleagues();
  const [connectionError, setConnectionError] = useState<string | null>(null);
  const [menuAnchorEl, setMenuAnchorEl] = useState<null | HTMLElement>(null);
  const [unlinkDialogOpen, setUnlinkDialogOpen] = useState(false);
  const [unlinkConfirmHandle, setUnlinkConfirmHandle] = useState("");

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

  const handleApproveRequest = async () => {
    if (!bio) return;

    setConnectionError(null);

    try {
      await approveColleague(bio.handle);
      // Refetch the profile to get updated connection state
      await refetch();
    } catch (error) {
      setConnectionError(
        error instanceof Error
          ? t(error.message)
          : t("profile.error.approvalFailed")
      );
    }
  };

  const handleDeclineRequest = async () => {
    if (!bio) return;

    setConnectionError(null);

    try {
      await rejectColleague(bio.handle);
      // Refetch the profile to get updated connection state
      await refetch();
    } catch (error) {
      setConnectionError(
        error instanceof Error
          ? t(error.message)
          : t("profile.error.rejectFailed")
      );
    }
  };

  const handleMenuClick = (event: React.MouseEvent<HTMLElement>) => {
    setMenuAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setMenuAnchorEl(null);
  };

  const handleUnlinkClick = () => {
    handleMenuClose();
    setUnlinkDialogOpen(true);
  };

  const handleUnlinkDialogClose = () => {
    setUnlinkDialogOpen(false);
    setUnlinkConfirmHandle("");
    setConnectionError(null);
  };

  const handleUnlinkConfirm = async () => {
    if (!bio) return;

    if (unlinkConfirmHandle !== bio.handle) {
      setConnectionError(t("profile.error.handleMismatch"));
      return;
    }

    try {
      await unlinkColleague(bio.handle);
      handleUnlinkDialogClose();
      await refetch();
    } catch (error) {
      setConnectionError(
        error instanceof Error
          ? t(error.message)
          : t("profile.error.unlinkFailed")
      );
    }
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
              startIcon={
                isApproving ? (
                  <CircularProgress size={20} />
                ) : (
                  <CheckCircleIcon />
                )
              }
              onClick={handleApproveRequest}
              disabled={isApproving || isRejecting}
              fullWidth
            >
              {isApproving ? t("common.loading") : t("profile.approveRequest")}
            </Button>
            <Button
              variant="outlined"
              color="error"
              startIcon={
                isRejecting ? <CircularProgress size={20} /> : <CancelIcon />
              }
              onClick={handleDeclineRequest}
              disabled={isApproving || isRejecting}
              fullWidth
            >
              {isRejecting ? t("common.loading") : t("profile.declineRequest")}
            </Button>
            {connectionError && (
              <Alert severity="error">{connectionError}</Alert>
            )}
          </Stack>
        );

      case "CONNECTED":
        return (
          <Stack spacing={2}>
            <Alert severity="success" icon={<VerifiedIcon />} sx={{ mb: 2 }}>
              {t("profile.mutuallyVerifiedColleague")}
            </Alert>
            <Box sx={{ display: "flex", justifyContent: "flex-end" }}>
              <IconButton onClick={handleMenuClick} size="small">
                <MoreVertIcon />
              </IconButton>
            </Box>
            <Menu
              anchorEl={menuAnchorEl}
              open={Boolean(menuAnchorEl)}
              onClose={handleMenuClose}
            >
              <MenuItem
                onClick={handleUnlinkClick}
                sx={{ color: "error.main" }}
              >
                <LinkOffIcon sx={{ mr: 1 }} />
                {t("profile.unlinkConnection")}
              </MenuItem>
            </Menu>
            <Dialog
              open={unlinkDialogOpen}
              onClose={handleUnlinkDialogClose}
              maxWidth="sm"
              fullWidth
            >
              <DialogTitle sx={{ color: "error.main" }}>
                {t("profile.unlinkConfirmTitle")}
              </DialogTitle>
              <DialogContent>
                <Typography sx={{ mb: 2 }}>
                  {t("profile.unlinkConfirmMessage", { handle: bio.handle })}
                </Typography>
                <TextField
                  fullWidth
                  label={t("profile.unlinkConfirmHandleLabel")}
                  value={unlinkConfirmHandle}
                  onChange={(e) => setUnlinkConfirmHandle(e.target.value)}
                  error={Boolean(connectionError)}
                  helperText={connectionError}
                  sx={{ mt: 1 }}
                />
              </DialogContent>
              <DialogActions>
                <Button onClick={handleUnlinkDialogClose}>
                  {t("common.cancel")}
                </Button>
                <Button
                  onClick={handleUnlinkConfirm}
                  color="error"
                  variant="contained"
                  disabled={isUnlinking || unlinkConfirmHandle !== bio.handle}
                  startIcon={
                    isUnlinking ? (
                      <CircularProgress size={20} />
                    ) : (
                      <LinkOffIcon />
                    )
                  }
                >
                  {isUnlinking
                    ? t("common.loading")
                    : t("profile.unlinkConfirm")}
                </Button>
              </DialogActions>
            </Dialog>
            {connectionError && (
              <Alert severity="error">{connectionError}</Alert>
            )}
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
                    </Box>
                  </Box>
                )}
              </Box>

              <Divider sx={{ mb: 4 }} />

              {/* Work History section */}
              <Box>
                <WorkHistory userHandle={userHandle} canEdit={false} />
              </Box>

              <Divider sx={{ my: 4 }} />

              {/* Education section */}
              <Box>
                <Education userHandle={userHandle} canEdit={false} />
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
