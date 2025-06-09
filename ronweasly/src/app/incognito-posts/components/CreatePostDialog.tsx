"use client";

import { useCreateIncognitoPost } from "@/hooks/useIncognitoPosts";
import { useTranslation } from "@/hooks/useTranslation";
import {
  Box,
  Button,
  Chip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
  Typography,
} from "@mui/material";
import { AddIncognitoPostRequest, VTag, VTagID } from "@vetchium/typespec";
import { useEffect, useState } from "react";
import TagSelector from "./TagSelector";

interface CreatePostDialogProps {
  open: boolean;
  onClose: () => void;
  onPostCreated: () => void;
  onError: (error: string) => void;
}

export default function CreatePostDialog({
  open,
  onClose,
  onPostCreated,
  onError,
}: CreatePostDialogProps) {
  const { t } = useTranslation();
  const { isCreating, error, createPost } = useCreateIncognitoPost();

  const [content, setContent] = useState("");
  const [selectedTags, setSelectedTags] = useState<VTagID[]>([]);
  const [availableTags, setAvailableTags] = useState<VTag[]>([]);
  const [contentError, setContentError] = useState("");
  const [tagsError, setTagsError] = useState("");

  useEffect(() => {
    if (error) {
      onError(error.message || t("incognitoPosts.errors.createFailed"));
    }
  }, [error, onError, t]);

  const handleClose = () => {
    if (!isCreating) {
      setContent("");
      setSelectedTags([]);
      setContentError("");
      setTagsError("");
      onClose();
    }
  };

  const validateForm = () => {
    let isValid = true;

    // Validate content
    if (!content.trim()) {
      setContentError(t("incognitoPosts.compose.contentRequired"));
      isValid = false;
    } else if (content.length > 1024) {
      setContentError(t("incognitoPosts.compose.contentTooLong"));
      isValid = false;
    } else {
      setContentError("");
    }

    // Validate tags
    if (selectedTags.length === 0) {
      setTagsError(t("incognitoPosts.compose.requiredTags"));
      isValid = false;
    } else if (selectedTags.length > 3) {
      setTagsError(t("incognitoPosts.compose.maxTags"));
      isValid = false;
    } else {
      setTagsError("");
    }

    return isValid;
  };

  const handleSubmit = async () => {
    if (!validateForm()) {
      return;
    }

    const request = new AddIncognitoPostRequest();
    request.content = content.trim();
    request.tag_ids = selectedTags;

    const postId = await createPost(request);
    if (postId) {
      onPostCreated();
      handleClose();
    }
  };

  const handleTagToggle = (tagId: VTagID) => {
    setSelectedTags((prev) => {
      if (prev.includes(tagId)) {
        return prev.filter((id) => id !== tagId);
      } else if (prev.length < 3) {
        return [...prev, tagId];
      }
      return prev;
    });
  };

  const handleTagsLoaded = (tags: VTag[]) => {
    setAvailableTags(tags);
  };

  const getSelectedTagNames = () => {
    return selectedTags
      .map((tagId) => availableTags.find((tag) => tag.id === tagId)?.name)
      .filter(Boolean);
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="md" fullWidth>
      <DialogTitle>{t("incognitoPosts.compose.title")}</DialogTitle>
      <DialogContent>
        <Box sx={{ display: "flex", flexDirection: "column", gap: 3, mt: 1 }}>
          {/* Content Input */}
          <TextField
            label={t("incognitoPosts.compose.contentLabel")}
            placeholder={t("incognitoPosts.compose.contentPlaceholder")}
            multiline
            rows={6}
            value={content}
            onChange={(e) => setContent(e.target.value)}
            error={!!contentError}
            helperText={
              contentError ||
              `${content.length}/1024 ${t("posts.charactersLimit")}`
            }
            fullWidth
          />

          {/* Tag Selection */}
          <Box>
            <Typography variant="subtitle2" gutterBottom>
              {t("incognitoPosts.compose.tagsLabel")}
            </Typography>
            <Typography variant="body2" color="text.secondary" gutterBottom>
              {t("incognitoPosts.compose.selectTags")} (1-3)
            </Typography>

            {/* Available Tags */}
            <Box sx={{ display: "flex", flexWrap: "wrap", gap: 1, mb: 2 }}>
              {availableTags.map((tag) => (
                <Chip
                  key={tag.id}
                  label={tag.name}
                  onClick={() => handleTagToggle(tag.id)}
                  color={selectedTags.includes(tag.id) ? "primary" : "default"}
                  variant={
                    selectedTags.includes(tag.id) ? "filled" : "outlined"
                  }
                  clickable
                />
              ))}
            </Box>

            {/* Selected Tags Display */}
            {selectedTags.length > 0 && (
              <Box sx={{ mb: 1 }}>
                <Typography variant="body2" color="text.secondary" gutterBottom>
                  Selected: {getSelectedTagNames().join(", ")}
                </Typography>
              </Box>
            )}

            {tagsError && (
              <Typography variant="body2" color="error">
                {tagsError}
              </Typography>
            )}
          </Box>

          {/* Hidden TagSelector to load tags */}
          <Box sx={{ display: "none" }}>
            <TagSelector
              selectedTag=""
              onTagSelect={() => {}}
              onError={onError}
            />
          </Box>
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose} disabled={isCreating}>
          {t("incognitoPosts.compose.cancel")}
        </Button>
        <Button
          onClick={handleSubmit}
          variant="contained"
          disabled={isCreating || !content.trim() || selectedTags.length === 0}
        >
          {isCreating
            ? t("incognitoPosts.compose.posting")
            : t("incognitoPosts.compose.submit")}
        </Button>
      </DialogActions>
    </Dialog>
  );
}
