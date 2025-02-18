import { useState } from "react";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import { useTranslation } from "@/hooks/useTranslation";
import IconButton from "@mui/material/IconButton";
import EditIcon from "@mui/icons-material/Edit";
import SaveIcon from "@mui/icons-material/Save";
import CancelIcon from "@mui/icons-material/Cancel";

interface BioData {
  handle: string;
  full_name: string;
  short_bio: string;
  long_bio: string;
}

interface BioProps {
  bio: BioData;
  onSave: (bio: BioData) => Promise<void>;
  isLoading?: boolean;
}

export default function Bio({ bio, onSave, isLoading = false }: BioProps) {
  const { t } = useTranslation();
  const [isEditing, setIsEditing] = useState(false);
  const [editedBio, setEditedBio] = useState<BioData>(bio);

  const handleSave = async () => {
    await onSave(editedBio);
    setIsEditing(false);
  };

  const handleCancel = () => {
    setEditedBio(bio);
    setIsEditing(false);
  };

  if (isEditing) {
    return (
      <Box sx={{ width: "100%" }}>
        <Typography variant="h6" sx={{ mb: 2 }}>
          {t("profile.bio.title")}
        </Typography>
        <Box
          component="form"
          sx={{ display: "flex", flexDirection: "column", gap: 2 }}
        >
          <TextField
            label={t("profile.bio.fullName")}
            value={editedBio.full_name}
            onChange={(e) =>
              setEditedBio({ ...editedBio, full_name: e.target.value })
            }
            fullWidth
            required
          />
          <TextField
            label={t("profile.bio.handle")}
            value={editedBio.handle}
            onChange={(e) =>
              setEditedBio({ ...editedBio, handle: e.target.value })
            }
            fullWidth
            required
            InputProps={{
              startAdornment: "@",
            }}
          />
          <TextField
            label={t("profile.bio.shortBio")}
            value={editedBio.short_bio}
            onChange={(e) =>
              setEditedBio({ ...editedBio, short_bio: e.target.value })
            }
            fullWidth
            required
            multiline
            rows={2}
          />
          <TextField
            label={t("profile.bio.longBio")}
            value={editedBio.long_bio}
            onChange={(e) =>
              setEditedBio({ ...editedBio, long_bio: e.target.value })
            }
            fullWidth
            multiline
            rows={4}
          />
          <Box sx={{ display: "flex", gap: 2, justifyContent: "flex-end" }}>
            <Button
              variant="outlined"
              onClick={handleCancel}
              startIcon={<CancelIcon />}
              disabled={isLoading}
            >
              {t("profile.bio.cancel")}
            </Button>
            <Button
              variant="contained"
              onClick={handleSave}
              startIcon={<SaveIcon />}
              disabled={isLoading}
            >
              {t("profile.bio.save")}
            </Button>
          </Box>
        </Box>
      </Box>
    );
  }

  return (
    <Box sx={{ width: "100%", position: "relative" }}>
      <IconButton
        sx={{ position: "absolute", right: 0, top: 0 }}
        onClick={() => setIsEditing(true)}
        aria-label={t("profile.bio.title")}
      >
        <EditIcon />
      </IconButton>
      <Typography variant="h3" component="h1" sx={{ mb: 1 }}>
        {bio.full_name}
      </Typography>
      <Typography variant="subtitle1" color="text.secondary" sx={{ mb: 1 }}>
        @{bio.handle} • {bio.short_bio}
      </Typography>
      <Typography variant="body1" sx={{ mt: 2, whiteSpace: "pre-wrap" }}>
        {bio.long_bio}
      </Typography>
    </Box>
  );
}
