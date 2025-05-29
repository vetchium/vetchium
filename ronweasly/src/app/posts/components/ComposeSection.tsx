"use client";

import { config } from "@/config";
import { useMyDetails } from "@/hooks/useMyDetails";
import { useTranslation } from "@/hooks/useTranslation";
import LocalOfferIcon from "@mui/icons-material/LocalOffer";
import {
  Alert,
  Autocomplete,
  Avatar,
  Box,
  Button,
  Chip,
  CircularProgress,
  Collapse,
  Paper,
  TextField,
  Typography,
  useTheme,
} from "@mui/material";
import {
  AddFTPostRequest,
  AddPostRequest,
  AddPostResponse,
  HubUserTiers,
  VTag,
} from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

interface ComposeProps {
  onPostCreated: () => void;
  onError: (errorMessage: string) => void;
  onSuccess: (successMessage: string) => void;
}

export default function ComposeSection({
  onPostCreated,
  onError,
  onSuccess,
}: ComposeProps) {
  const { t } = useTranslation();
  const theme = useTheme();
  const router = useRouter();
  const {
    details,
    isLoading: isLoadingDetails,
    error: detailsError,
  } = useMyDetails();
  const [postContent, setPostContent] = useState("");
  const [selectedTags, setSelectedTags] = useState<VTag[]>([]);
  const [tagSuggestions, setTagSuggestions] = useState<VTag[]>([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [loading, setLoading] = useState(false);
  const [isExpanded, setIsExpanded] = useState(false);

  // Determine if user is on free tier
  const isFreeTier = details?.tier === HubUserTiers.FreeHubUserTier;
  const maxContentLength = isFreeTier ? 255 : 4096;

  // Fetch tag suggestions when user types
  useEffect(() => {
    const fetchTags = async () => {
      if (searchQuery.length >= 2) {
        const token = Cookies.get("session_token");
        if (!token) {
          router.push("/login");
          return;
        }

        try {
          const response = await fetch(
            `${config.API_SERVER_PREFIX}/hub/filter-vtags`,
            {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${token}`,
              },
              body: JSON.stringify({
                prefix: searchQuery,
              }),
            }
          );

          if (!response.ok) {
            if (response.status === 401) {
              Cookies.remove("session_token", { path: "/" });
              router.push("/login");
              return;
            }
            throw new Error(`Failed to fetch tags: ${response.statusText}`);
          }

          const data = await response.json();
          setTagSuggestions(Array.isArray(data) ? data : []);
        } catch (error) {
          console.error("Error fetching tags:", error);
          onError(t("posts.error.tagsFailed"));
          setTagSuggestions([]);
        }
      } else {
        setTagSuggestions([]);
      }
    };

    const debounceTimer = setTimeout(fetchTags, 300);
    return () => clearTimeout(debounceTimer);
  }, [searchQuery, t, router, onError]);

  const handlePublish = async () => {
    if (!postContent.trim()) {
      onError(t("posts.error.contentRequired"));
      return;
    }

    // Check content length based on tier
    if (postContent.length > maxContentLength) {
      onError(t("posts.error.contentTooLong"));
      return;
    }

    const token = Cookies.get("session_token");
    if (!token) {
      router.push("/login");
      return;
    }

    setLoading(true);

    try {
      // Use different endpoints based on user tier
      const endpoint = isFreeTier
        ? `${config.API_SERVER_PREFIX}/hub/add-ft-post`
        : `${config.API_SERVER_PREFIX}/hub/add-post`;

      // Prepare request body based on tier - both use same structure now
      let requestBody: AddPostRequest | AddFTPostRequest;
      if (isFreeTier) {
        requestBody = new AddFTPostRequest();
      } else {
        requestBody = new AddPostRequest();
      }

      requestBody.content = postContent;
      requestBody.tag_ids = selectedTags.map((tag) => tag.id).filter(Boolean);

      const response = await fetch(endpoint, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(requestBody),
      });

      if (!response.ok) {
        if (response.status === 401) {
          Cookies.remove("session_token", { path: "/" });
          router.push("/login");
          return;
        }
        throw new Error(`Failed to create post: ${response.statusText}`);
      }

      // Get post ID from response
      const data: AddPostResponse = await response.json();
      const postId = data.post_id;

      // Reset form
      setPostContent("");
      setSelectedTags([]);
      setIsExpanded(false);
      onSuccess(t("posts.success"));

      // Notify parent about successful post creation
      onPostCreated();

      // Redirect to the post details page
      router.push(`/posts/${postId}`);
    } catch (error) {
      console.error("Error creating post:", error);
      onError(t("posts.error.createFailed"));
    } finally {
      setLoading(false);
    }
  };

  // Show loading state while fetching user details
  if (isLoadingDetails) {
    return (
      <Paper
        sx={{
          p: 2,
          mb: 3,
          display: "flex",
          justifyContent: "center",
        }}
      >
        <CircularProgress size={24} />
      </Paper>
    );
  }

  if (detailsError) {
    return (
      <Paper sx={{ p: 2, mb: 3 }}>
        <Alert severity="error">{t("posts.error.loadingDetails")}</Alert>
      </Paper>
    );
  }

  return (
    <>
      {/* Compact composer */}
      <Paper
        sx={{
          p: 2,
          mb: 3,
          border: `1px solid ${theme.palette.divider}`,
          borderRadius: 2,
        }}
      >
        <Box sx={{ display: "flex", gap: 2, alignItems: "flex-start" }}>
          <Avatar sx={{ width: 40, height: 40 }}>
            {details?.full_name?.charAt(0) || details?.handle?.charAt(0)}
          </Avatar>
          <Box sx={{ flex: 1 }}>
            <TextField
              fullWidth
              placeholder={t("posts.compose.placeholder")}
              variant="outlined"
              size="small"
              value={postContent}
              onChange={(e) => setPostContent(e.target.value)}
              onClick={() => setIsExpanded(true)}
              multiline={isExpanded}
              rows={isExpanded ? 3 : 1}
              sx={{
                "& .MuiOutlinedInput-root": {
                  borderRadius: "20px",
                  backgroundColor: theme.palette.background.default,
                },
              }}
            />

            <Collapse in={isExpanded}>
              <Box sx={{ mt: 2 }}>
                {/* Free tier character limit notice */}
                {isFreeTier && (
                  <Alert severity="info" sx={{ mb: 2, py: 1 }}>
                    <Typography
                      variant="caption"
                      sx={{ display: "block", mb: 0.5 }}
                    >
                      <strong>{t("posts.freeTierLimits.title")}</strong>
                    </Typography>
                    <Typography variant="caption" sx={{ display: "block" }}>
                      {t("posts.freeTierLimits.characterLimit")}
                    </Typography>
                  </Alert>
                )}

                {/* Tags section */}
                <Box sx={{ mb: 2 }}>
                  {/* Show selected tags */}
                  {selectedTags.length > 0 && (
                    <Box
                      sx={{
                        display: "flex",
                        flexWrap: "wrap",
                        gap: 0.5,
                        mb: 1,
                      }}
                    >
                      {selectedTags.map((tag) => (
                        <Chip
                          key={tag.id}
                          label={tag.name}
                          onDelete={() => {
                            setSelectedTags(
                              selectedTags.filter((t) => t !== tag)
                            );
                          }}
                          color="primary"
                          size="small"
                          sx={{ fontSize: "0.7rem", height: "24px" }}
                        />
                      ))}
                    </Box>
                  )}

                  {/* Tag input - only show if less than 3 tags */}
                  {selectedTags.length < 3 && (
                    <Autocomplete
                      size="small"
                      options={tagSuggestions}
                      value={null}
                      inputValue={searchQuery}
                      onInputChange={(_, newInputValue) => {
                        setSearchQuery(newInputValue);
                      }}
                      onChange={(_, newValue) => {
                        if (!newValue) return;

                        // Don't add if already at max
                        if (selectedTags.length >= 3) return;

                        // Only work with existing tags
                        const newTag = newValue as VTag;

                        // Don't add if already selected
                        if (selectedTags.some((tag) => tag.id === newTag.id)) {
                          return;
                        }

                        setSelectedTags([...selectedTags, newTag]);
                        setSearchQuery("");
                      }}
                      getOptionLabel={(option) => option.name}
                      renderOption={(props, option) => {
                        const { key, ...otherProps } = props;
                        return (
                          <li key={key || option.id} {...otherProps}>
                            <Box
                              sx={{
                                display: "flex",
                                alignItems: "center",
                                gap: 1,
                              }}
                            >
                              <LocalOfferIcon fontSize="small" />
                              <span>{option.name}</span>
                            </Box>
                          </li>
                        );
                      }}
                      renderInput={(params) => (
                        <TextField
                          {...params}
                          placeholder={
                            selectedTags.length >= 3
                              ? t("posts.maxTags")
                              : t("posts.searchTagsOnly")
                          }
                          disabled={selectedTags.length >= 3}
                          InputProps={{
                            ...params.InputProps,
                            startAdornment: (
                              <LocalOfferIcon
                                sx={{ mr: 1, fontSize: "1rem" }}
                              />
                            ),
                          }}
                          sx={{
                            "& .MuiOutlinedInput-root": {
                              fontSize: "0.8rem",
                            },
                          }}
                        />
                      )}
                    />
                  )}
                </Box>

                {/* Character count */}
                <Box
                  sx={{
                    display: "flex",
                    justifyContent: "flex-start",
                    alignItems: "center",
                    mb: 1,
                  }}
                >
                  <Typography
                    variant="caption"
                    sx={{
                      color:
                        postContent.length > maxContentLength
                          ? theme.palette.error.main
                          : theme.palette.text.secondary,
                    }}
                  >
                    {postContent.length}/{maxContentLength}
                  </Typography>
                </Box>

                {/* Action buttons */}
                <Box
                  sx={{
                    display: "flex",
                    justifyContent: "space-between",
                    alignItems: "center",
                  }}
                >
                  <Button
                    size="small"
                    onClick={() => {
                      setIsExpanded(false);
                      setPostContent("");
                      setSelectedTags([]);
                    }}
                  >
                    Cancel
                  </Button>
                  <Button
                    variant="contained"
                    size="small"
                    onClick={handlePublish}
                    disabled={
                      loading ||
                      !postContent.trim() ||
                      postContent.length > maxContentLength
                    }
                    sx={{ borderRadius: "20px" }}
                  >
                    {loading ? (
                      <CircularProgress size={16} color="inherit" />
                    ) : (
                      t("posts.publish")
                    )}
                  </Button>
                </Box>
              </Box>
            </Collapse>
          </Box>
        </Box>
      </Paper>
    </>
  );
}
