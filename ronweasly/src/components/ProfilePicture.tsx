import { useRef, useState, useEffect } from "react";
import Box from "@mui/material/Box";
import Avatar from "@mui/material/Avatar";
import IconButton from "@mui/material/IconButton";
import CameraAltIcon from "@mui/icons-material/CameraAlt";
import DeleteIcon from "@mui/icons-material/Delete";
import Dialog from "@mui/material/Dialog";
import DialogTitle from "@mui/material/DialogTitle";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogActions from "@mui/material/DialogActions";
import Button from "@mui/material/Button";
import { useTranslation } from "@/hooks/useTranslation";
import Cookies from "js-cookie";
import { config } from "@/config";

interface ProfilePictureProps {
  imageUrl?: string;
  size?: number;
  onImageSelect?: (file: File) => Promise<void>;
  onRemove?: () => void;
  isLoading?: boolean;
}

export default function ProfilePicture({
  imageUrl,
  size = 150,
  onImageSelect,
  onRemove,
  isLoading = false,
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

  const handleImageClick = () => {
    if (isLoading) return;

    if (imageData) {
      setZoomOpen(true);
    } else if (onImageSelect) {
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

  return (
    <>
      <Box
        sx={{
          position: "relative",
          width: size,
          height: size,
          mx: "auto",
          mb: 4,
        }}
      >
        <Avatar
          src={imageData}
          sx={{
            width: size,
            height: size,
            cursor: isLoading ? "default" : "pointer",
            opacity: isLoading ? 0.7 : 1,
            transition: "opacity 0.2s",
          }}
          onClick={handleImageClick}
        />
        {imageData && onImageSelect && (
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
            disabled={isLoading}
            onClick={() => setConfirmOpen(true)}
            aria-label={t("profile.picture.remove")}
            size="small"
          >
            <DeleteIcon fontSize="small" />
          </IconButton>
        )}
        {onImageSelect && (
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
            disabled={isLoading}
            onClick={() => fileInputRef.current?.click()}
            aria-label={t("profile.picture.change")}
          >
            <CameraAltIcon />
          </IconButton>
        )}
        {onImageSelect && (
          <input
            type="file"
            ref={fileInputRef}
            onChange={handleFileChange}
            accept="image/*"
            style={{ display: "none" }}
            aria-label={t("profile.picture.upload")}
          />
        )}
      </Box>

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
