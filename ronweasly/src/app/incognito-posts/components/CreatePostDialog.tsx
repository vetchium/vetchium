"use client";

import { config } from "@/config";
import { useCreateIncognitoPost } from "@/hooks/useIncognitoPosts";
import { useTranslation } from "@/hooks/useTranslation";
import {
  Alert,
  Autocomplete,
  Box,
  Button,
  Chip,
  CircularProgress,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
  Typography,
} from "@mui/material";
import { AddIncognitoPostRequest, VTag } from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useEffect, useState } from "react";

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
  const [selectedTags, setSelectedTags] = useState<VTag[]>([]);
  const [availableTags, setAvailableTags] = useState<VTag[]>([]);
  const [contentError, setContentError] = useState("");
  const [tagsError, setTagsError] = useState("");
  const [isLoadingTags, setIsLoadingTags] = useState(false);
  const [tagInputValue, setTagInputValue] = useState("");

  // Load tags when dialog opens
  useEffect(() => {
    if (open && availableTags.length === 0) {
      loadTags();
    }
  }, [open]);

  const loadTags = async () => {
    setIsLoadingTags(true);
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        throw new Error("User not authenticated");
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/filter-vtags`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({ prefix: "" }),
        }
      );

      if (!response.ok) {
        throw new Error(`Failed to fetch tags: ${response.statusText}`);
      }

      const data: VTag[] = await response.json();
      setAvailableTags(data);
    } catch (error) {
      onError(
        error instanceof Error
          ? error.message
          : t("incognitoPosts.errors.tagLoadFailed")
      );
    } finally {
      setIsLoadingTags(false);
    }
  };

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
      setTagInputValue("");
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
    request.tag_ids = selectedTags.map((tag) => tag.id);

    const postId = await createPost(request);
    if (postId) {
      onPostCreated();
      handleClose();
    }
  };

  const handleTagsChange = (_: any, newValue: VTag[]) => {
    if (newValue.length <= 3) {
      setSelectedTags(newValue);
      // Clear tag error when user selects a tag
      if (tagsError && newValue.length > 0) {
        setTagsError("");
      }
    }
  };

  const isFormValid =
    content.trim() && selectedTags.length > 0 && selectedTags.length <= 3;

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="md" fullWidth>
      <DialogTitle>{t("incognitoPosts.compose.title")}</DialogTitle>
      <DialogContent>
        <Box sx={{ display: "flex", flexDirection: "column", gap: 3, mt: 1 }}>
          {/* Warning messages at the top of the dialog */}
          <Box sx={{ mb: 2 }}>
            <Alert severity="warning" sx={{ mb: 2 }}>
              {t("incognitoPosts.compose.warnings.dataRetention")}
            </Alert>
            <Alert severity="info">
              {t("incognitoPosts.compose.warnings.respectfulPosting")}
            </Alert>
          </Box>

          {/* Content Input */}
          <TextField
            label={t("incognitoPosts.compose.contentLabel")}
            placeholder={t("incognitoPosts.compose.contentPlaceholder")}
            multiline
            rows={6}
            value={content}
            onChange={(e) => {
              setContent(e.target.value);
              // Clear content error when user types
              if (contentError) {
                setContentError("");
              }
            }}
            error={!!contentError}
            helperText={contentError || `${content.length}/1024 characters`}
            fullWidth
          />

          {/* Tag Selection with Autocomplete */}
          <Box>
            <Typography variant="subtitle2" gutterBottom>
              {t("incognitoPosts.compose.tagsLabel")}
            </Typography>

            <Autocomplete
              multiple
              options={availableTags}
              getOptionLabel={(option) => option.name}
              value={selectedTags}
              onChange={handleTagsChange}
              inputValue={tagInputValue}
              onInputChange={(_, newInputValue) => {
                setTagInputValue(newInputValue);
              }}
              loading={isLoadingTags}
              renderTags={(value, getTagProps) =>
                value.map((option, index) => (
                  <Chip
                    variant="outlined"
                    label={option.name}
                    {...getTagProps({ index })}
                    key={option.id}
                  />
                ))
              }
              renderInput={(params) => (
                <TextField
                  {...params}
                  placeholder={t("incognitoPosts.compose.tagPlaceholder")}
                  error={!!tagsError}
                  helperText={
                    tagsError ||
                    `Select 1-3 tags (${selectedTags.length}/3 selected)`
                  }
                  InputProps={{
                    ...params.InputProps,
                    endAdornment: (
                      <>
                        {isLoadingTags ? (
                          <CircularProgress color="inherit" size={20} />
                        ) : null}
                        {params.InputProps.endAdornment}
                      </>
                    ),
                  }}
                />
              )}
              fullWidth
              limitTags={3}
              disableCloseOnSelect
              isOptionEqualToValue={(option, value) => option.id === value.id}
              getOptionDisabled={() => selectedTags.length >= 3}
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
          disabled={isCreating || !isFormValid}
        >
          {isCreating
            ? t("incognitoPosts.compose.posting")
            : t("incognitoPosts.compose.submit")}
        </Button>
      </DialogActions>
    </Dialog>
  );
}
