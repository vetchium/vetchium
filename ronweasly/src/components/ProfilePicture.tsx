import { useRef, useState, useEffect } from "react";
import Box from "@mui/material/Box";
import Avatar from "@mui/material/Avatar";
import IconButton from "@mui/material/IconButton";
import CameraAltIcon from "@mui/icons-material/CameraAlt";
import { useTranslation } from "@/hooks/useTranslation";
import Cookies from "js-cookie";

interface ProfilePictureProps {
  imageUrl?: string;
  size?: number;
  onImageSelect: (file: File) => Promise<void>;
  isLoading?: boolean;
}

export default function ProfilePicture({
  imageUrl,
  size = 150,
  onImageSelect,
  isLoading = false,
}: ProfilePictureProps) {
  const { t } = useTranslation();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const token = Cookies.get("session_token");
  const [imageData, setImageData] = useState<string | undefined>(undefined);

  useEffect(() => {
    if (!imageUrl || !token) return;

    fetch(imageUrl, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then((response) => response.blob())
      .then((blob) => {
        const url = URL.createObjectURL(blob);
        setImageData(url);
        return () => URL.revokeObjectURL(url);
      })
      .catch((error) => {
        console.error("Failed to load profile picture:", error);
      });
  }, [imageUrl, token]);

  const handleImageClick = () => {
    if (!isLoading) {
      fileInputRef.current?.click();
    }
  };

  const handleFileChange = async (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    const file = event.target.files?.[0];
    if (file) {
      await onImageSelect(file);
      // Clear the input so the same file can be selected again
      event.target.value = "";
    }
  };

  return (
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
        onClick={handleImageClick}
        aria-label={t("profile.picture.change")}
      >
        <CameraAltIcon />
      </IconButton>
      <input
        type="file"
        ref={fileInputRef}
        onChange={handleFileChange}
        accept="image/*"
        style={{ display: "none" }}
        aria-label={t("profile.picture.upload")}
      />
    </Box>
  );
}
