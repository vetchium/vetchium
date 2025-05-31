"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import Bio from "@/components/Bio";
import { Certifications } from "@/components/Certifications";
import ChangeHandle from "@/components/ChangeHandle";
import { Education } from "@/components/Education";
import OfficialEmails from "@/components/OfficialEmails";
import { Patents } from "@/components/Patents";
import ProfilePicture from "@/components/ProfilePicture";
import { Publications } from "@/components/Publications";
import { config } from "@/config";
import { useAuth } from "@/hooks/useAuth";
import { useMyHandle } from "@/hooks/useMyHandle";
import { useMyTier } from "@/hooks/useMyTier";
import { useProfile } from "@/hooks/useProfile";
import { useTranslation } from "@/hooks/useTranslation";
import EditIcon from "@mui/icons-material/Edit";
import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import CircularProgress from "@mui/material/CircularProgress";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";
import Divider from "@mui/material/Divider";
import IconButton from "@mui/material/IconButton";
import Link from "@mui/material/Link";
import Typography from "@mui/material/Typography";
import { HubUserTiers } from "@vetchium/typespec";
import { useState } from "react";
import { WorkHistory } from "../u/[handle]/WorkHistory";

export default function MyProfilePage() {
  const { myHandle, isLoading: isLoadingHandle } = useMyHandle();
  const { t } = useTranslation();
  useAuth(); // Check authentication and redirect if not authenticated
  const [timestamp, setTimestamp] = useState(Date.now());
  const [showChangeHandle, setShowChangeHandle] = useState(false);
  const [showUpgradeDialog, setShowUpgradeDialog] = useState(false);
  const {
    bio,
    isLoading: isLoadingBio,
    error: bioError,
    isSaving,
    updateBio,
    uploadProfilePicture,
  } = useProfile(myHandle ?? "");
  const { tier, isLoading: isLoadingTier, error: tierError } = useMyTier();

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

  const isLoading = isLoadingBio || isLoadingHandle || isLoadingTier;
  const error = bioError || tierError;

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
          userTier={tier}
          isTierLoading={isLoadingTier}
        />

        {/* Handle Display and Edit Section */}
        <Box sx={{ display: "flex", alignItems: "center", mt: 2 }}>
          <Typography variant="h6" component="h2">
            {"@"}
            <strong>{myHandle}</strong>
          </Typography>
          <IconButton
            aria-label="edit handle"
            size="small"
            sx={{ ml: 1 }}
            onClick={() => {
              if (tier === HubUserTiers.PaidHubUserTier) {
                setShowChangeHandle(true);
              } else {
                setShowUpgradeDialog(true);
              }
            }}
          >
            <EditIcon fontSize="small" />
          </IconButton>
        </Box>

        {isLoadingBio ? (
          <Box sx={{ display: "flex", justifyContent: "center", my: 4 }}>
            <CircularProgress />
          </Box>
        ) : bio ? (
          <Box sx={{ mb: 4 }}>
            <Bio bio={bio} onSave={updateBio} isLoading={isSaving} />
          </Box>
        ) : null}

        <Divider sx={{ my: 6 }} />

        <Box sx={{ mt: 4, mb: 4 }}>
          <OfficialEmails />
        </Box>

        <Divider sx={{ my: 6 }} />

        <Box sx={{ mt: 4, mb: 4 }}>
          <WorkHistory userHandle={myHandle} canEdit={true} />
        </Box>

        <Divider sx={{ my: 6 }} />

        <Box sx={{ mt: 4, mb: 4 }}>
          <Education userHandle={myHandle} canEdit={true} />
        </Box>

        <Divider sx={{ my: 6 }} />

        <Box sx={{ mt: 4, mb: 4 }}>
          <Patents userHandle={myHandle} canEdit={true} />
        </Box>

        <Divider sx={{ my: 6 }} />

        <Box sx={{ mt: 4, mb: 4 }}>
          <Publications userHandle={myHandle} canEdit={true} />
        </Box>

        <Divider sx={{ my: 6 }} />

        <Box sx={{ mt: 4 }}>
          <Certifications userHandle={myHandle} canEdit={true} />
        </Box>
      </Box>

      {/* Change Handle Dialog */}
      {showChangeHandle && tier && (
        <Dialog
          open={showChangeHandle}
          onClose={() => setShowChangeHandle(false)}
          maxWidth="sm"
          fullWidth
        >
          <DialogTitle>{t("profile.changeHandle.title")}</DialogTitle>
          <DialogContent>
            <ChangeHandle
              currentHandle={myHandle}
              userTier={tier}
              onSuccess={() => setShowChangeHandle(false)}
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setShowChangeHandle(false)}>
              {t("common.close")}
            </Button>
          </DialogActions>
        </Dialog>
      )}

      {/* Upgrade Required Dialog */}
      <Dialog
        open={showUpgradeDialog}
        onClose={() => setShowUpgradeDialog(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>{t("profile.upgradeRequired.title")}</DialogTitle>
        <DialogContent>
          <DialogContentText>
            {t("profile.upgradeRequired.message")}
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowUpgradeDialog(false)}>
            {t("common.cancel")}
          </Button>
          <Button
            component={Link}
            href="/upgrade"
            variant="contained"
            color="primary"
            onClick={() => setShowUpgradeDialog(false)}
          >
            {t("profile.upgradeRequired.upgradeButton")}
          </Button>
        </DialogActions>
      </Dialog>
    </AuthenticatedLayout>
  );
}
