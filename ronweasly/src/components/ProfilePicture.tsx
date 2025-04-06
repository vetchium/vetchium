import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import CameraAltIcon from "@mui/icons-material/CameraAlt";
import DeleteIcon from "@mui/icons-material/Delete";
import Avatar from "@mui/material/Avatar";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import CircularProgress from "@mui/material/CircularProgress";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";
import IconButton from "@mui/material/IconButton";
import Link from "@mui/material/Link";
import Typography from "@mui/material/Typography";
import { HubUserTier, HubUserTiers } from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useEffect, useRef, useState } from "react";

interface ProfilePictureProps {
  imageUrl?: string;
  size?: number;
  onImageSelect?: (file: File) => Promise<void>;
  onRemove?: () => void;
  isLoading?: boolean;
  userTier?: HubUserTier | null;
  isTierLoading?: boolean;
}

export default function ProfilePicture({
  imageUrl,
  size = 150,
  onImageSelect,
  onRemove,
  isLoading = false,
  userTier,
  isTierLoading = false,
}: ProfilePictureProps) {
  const { t } = useTranslation();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const deleteButtonRef = useRef<HTMLButtonElement>(null);
  const cancelButtonRef = useRef<HTMLButtonElement>(null);
  const token = Cookies.get("session_token");
  const [imageData, setImageData] = useState<string | undefined>(undefined);
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [zoomOpen, setZoomOpen] = useState(false);
  const [timestamp, setTimestamp] = useState(Date.now());

  useEffect(() => {
    if (!imageUrl || !token) return;

    fetch(`${imageUrl}?t=${timestamp}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then((response) => {
        if (response.status === 404) {
          // Not an error, just means no profile picture
          setImageData(undefined);
          return null;
        }
        if (!response.ok) {
          throw new Error("Failed to load image");
        }
        return response.blob();
      })
      .then((blob) => {
        if (!blob) return; // Handle the 404 case
        const url = URL.createObjectURL(blob);
        setImageData(url);
        return () => URL.revokeObjectURL(url);
      })
      .catch((error) => {
        console.error("Failed to load profile picture:", error);
        setImageData(undefined);
      });
  }, [imageUrl, token, timestamp]);

  const isPaidUser = userTier === HubUserTiers.PaidHubUserTier;

  const handleImageClick = () => {
    if (isLoading || isTierLoading) return;

    if (imageData) {
      setZoomOpen(true);
    } else if (onImageSelect && isPaidUser) {
      fileInputRef.current?.click();
    }
  };

  const handleFileChange = async (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    const file = event.target.files?.[0];
    if (file && onImageSelect) {
      try {
        await onImageSelect(file);
        // Clear the input so the same file can be selected again
        event.target.value = "";
        // Refresh the image
        setTimestamp(Date.now());
      } catch (error) {
        console.error("Failed to upload profile picture:", error);
      }
    }
  };

  const handleRemove = async () => {
    if (!token || isLoading) return;

    try {
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/remove-profile-picture`,
        {
          method: "POST",
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (!response.ok) {
        throw new Error(t("profile.picture.removeFailed"));
      }

      setImageData(undefined);
      setConfirmOpen(false);
      setTimestamp(Date.now());
      onRemove?.();
    } catch (error) {
      console.error("Failed to remove profile picture:", error);
    }
  };

  const handleClose = () => {
    setConfirmOpen(false);
  };

  const showUpgradeMessage = onImageSelect && !isPaidUser && !isTierLoading;

  return (
    <>
      <Box
        sx={{
          position: "relative",
          width: size,
          height: size,
          mx: "auto",
          mb: 2,
        }}
      >
        <Avatar
          src={imageData}
          sx={{
            width: size,
            height: size,
            cursor: isLoading || isTierLoading ? "default" : "pointer",
            opacity: isLoading || isTierLoading ? 0.7 : 1,
            transition: "opacity 0.2s",
          }}
          onClick={handleImageClick}
        />
        {isTierLoading && (
          <CircularProgress
            size={24}
            sx={{
              position: "absolute",
              top: "50%",
              left: "50%",
              marginTop: "-12px",
              marginLeft: "-12px",
            }}
          />
        )}
        {imageData && onImageSelect && isPaidUser && (
          <IconButton
            ref={deleteButtonRef}
            sx={{
              position: "absolute",
              top: 0,
              right: 0,
              backgroundColor: "error.main",
              color: "white",
              "&:hover": {
                backgroundColor: "error.dark",
              },
              width: 32,
              height: 32,
            }}
            disabled={isLoading || isTierLoading}
            onClick={() => setConfirmOpen(true)}
            aria-label={t("profile.picture.remove")}
            size="small"
          >
            <DeleteIcon fontSize="small" />
          </IconButton>
        )}
        {onImageSelect && isPaidUser && (
          <IconButton
            sx={{
              position: "absolute",
              bottom: 0,
              right: 0,
              backgroundColor: "background.paper",
              "&:hover": {
                backgroundColor: "action.hover",
              },
            }}
            disabled={isLoading || isTierLoading}
            onClick={handleImageClick}
            aria-label={t("profile.picture.change")}
          >
            <CameraAltIcon />
          </IconButton>
        )}
        {onImageSelect && isPaidUser && (
          <input
            type="file"
            ref={fileInputRef}
            onChange={handleFileChange}
            accept="image/*"
            style={{ display: "none" }}
          />
        )}
      </Box>
      {showUpgradeMessage && (
        <Typography
          variant="caption"
          display="block"
          textAlign="center"
          sx={{ mb: 4 }}
        >
          {t("profile.picture.upgradePrompt")}{" "}
          <Link href="/upgrade" underline="hover">
            {t("profile.picture.upgradeLink")}
          </Link>
        </Typography>
      )}

      <Dialog
        open={confirmOpen}
        onClose={handleClose}
        aria-labelledby="remove-profile-picture-dialog"
        aria-describedby="remove-profile-picture-description"
      >
        <DialogTitle id="remove-profile-picture-dialog">
          {t("profile.picture.removeConfirmTitle")}
        </DialogTitle>
        <DialogContent>
          <DialogContentText id="remove-profile-picture-description">
            {t("profile.picture.removeConfirmMessage")}
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose} autoFocus>
            {t("common.cancel")}
          </Button>
          <Button onClick={handleRemove} color="error" variant="contained">
            {t("profile.picture.removeConfirm")}
          </Button>
        </DialogActions>
      </Dialog>

      <Dialog
        open={zoomOpen}
        onClose={() => setZoomOpen(false)}
        maxWidth="lg"
        aria-labelledby="zoom-profile-picture-dialog"
      >
        <DialogContent sx={{ p: 0 }}>
          {imageData && (
            <img
              src={imageData}
              alt={t("profile.picture.fullSize")}
              style={{
                maxWidth: "100%",
                maxHeight: "90vh",
                display: "block",
                margin: "0 auto",
              }}
            />
          )}
        </DialogContent>
      </Dialog>
    </>
  );
}
