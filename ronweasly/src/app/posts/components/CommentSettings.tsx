"use client";

import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import SettingsIcon from "@mui/icons-material/Settings";
import {
  Box,
  Button,
  Checkbox,
  CircularProgress,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  FormControlLabel,
  IconButton,
  Menu,
  MenuItem,
  Typography,
} from "@mui/material";
import {
  DisablePostCommentsRequest,
  EnablePostCommentsRequest,
} from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useState } from "react";

interface CommentSettingsProps {
  postId: string;
  canComment: boolean;
  onCommentSettingsChange: (canComment: boolean) => void;
}

export default function CommentSettings({
  postId,
  canComment,
  onCommentSettingsChange,
}: CommentSettingsProps) {
  const { t } = useTranslation();
  const router = useRouter();

  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [deleteExistingComments, setDeleteExistingComments] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleEnableComments = async () => {
    setLoading(true);
    const token = Cookies.get("session_token");
    if (!token) {
      router.push("/login");
      return;
    }

    try {
      const request = new EnablePostCommentsRequest();
      request.post_id = postId;

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/enable-post-comments`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(request),
        }
      );

      if (!response.ok) {
        if (response.status === 401) {
          Cookies.remove("session_token", { path: "/" });
          router.push("/login");
          return;
        }
        throw new Error(`Failed to enable comments: ${response.statusText}`);
      }

      onCommentSettingsChange(true);
      handleMenuClose();
    } catch (error) {
      console.error("Error enabling comments:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleDisableComments = async () => {
    setLoading(true);
    const token = Cookies.get("session_token");
    if (!token) {
      router.push("/login");
      return;
    }

    try {
      const request = new DisablePostCommentsRequest();
      request.post_id = postId;
      request.delete_existing_comments = deleteExistingComments;

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/disable-post-comments`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(request),
        }
      );

      if (!response.ok) {
        if (response.status === 401) {
          Cookies.remove("session_token", { path: "/" });
          router.push("/login");
          return;
        }
        throw new Error(`Failed to disable comments: ${response.statusText}`);
      }

      onCommentSettingsChange(false);
      setDialogOpen(false);
      handleMenuClose();
      setDeleteExistingComments(false);
    } catch (error) {
      console.error("Error disabling comments:", error);
    } finally {
      setLoading(false);
    }
  };

  const openDisableDialog = () => {
    setDialogOpen(true);
    handleMenuClose();
  };

  return (
    <>
      <IconButton
        size="small"
        onClick={handleMenuOpen}
        sx={{
          color: "text.secondary",
          "&:hover": {
            color: "text.primary",
          },
        }}
      >
        <SettingsIcon fontSize="small" />
      </IconButton>

      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleMenuClose}
        anchorOrigin={{
          vertical: "bottom",
          horizontal: "right",
        }}
        transformOrigin={{
          vertical: "top",
          horizontal: "right",
        }}
      >
        <MenuItem
          onClick={canComment ? openDisableDialog : handleEnableComments}
        >
          {canComment
            ? t("comments.disableComments")
            : t("comments.enableComments")}
        </MenuItem>
      </Menu>

      <Dialog
        open={dialogOpen}
        onClose={() => setDialogOpen(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>{t("comments.disableComments")}</DialogTitle>
        <DialogContent>
          <DialogContentText sx={{ mb: 2 }}>
            {t("comments.confirmDisableComments")}
          </DialogContentText>
          <FormControlLabel
            control={
              <Checkbox
                checked={deleteExistingComments}
                onChange={(e) => setDeleteExistingComments(e.target.checked)}
              />
            }
            label={
              <Box>
                <Typography variant="body2">
                  {t("comments.deleteExistingComments")}
                </Typography>
                <Typography variant="caption" color="text.secondary">
                  {t("comments.deleteExistingCommentsHelp")}
                </Typography>
              </Box>
            }
          />
          {deleteExistingComments && (
            <Typography
              variant="body2"
              color="error"
              sx={{ mt: 1, fontWeight: 500 }}
            >
              {t("comments.confirmDisableCommentsWithDelete")}
            </Typography>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDialogOpen(false)}>
            {t("common.cancel")}
          </Button>
          <Button
            onClick={handleDisableComments}
            color="error"
            disabled={loading}
          >
            {loading ? (
              <CircularProgress size={16} color="inherit" />
            ) : (
              t("comments.disableComments")
            )}
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
}
