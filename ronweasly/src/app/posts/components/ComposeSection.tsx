"use client";

import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import AddIcon from "@mui/icons-material/Add";
import LocalOfferIcon from "@mui/icons-material/LocalOffer";
import {
  Autocomplete,
  Box,
  Button,
  Chip,
  CircularProgress,
  Paper,
  TextField,
  Typography,
} from "@mui/material";
import { createFilterOptions } from "@mui/material/Autocomplete";
import { VTag } from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

// Interface for tags including free input
interface TagOption extends VTag {
  inputValue?: string;
}

interface ComposeProps {
  onPostCreated: () => void;
  onError: (errorMessage: string) => void;
  onSuccess: (successMessage: string) => void;
}

interface AddPostResponse {
  post_id: string;
}

export default function ComposeSection({
  onPostCreated,
  onError,
  onSuccess,
}: ComposeProps) {
  const { t } = useTranslation();
  const router = useRouter();
  const [postContent, setPostContent] = useState("");
  const [selectedTags, setSelectedTags] = useState<VTag[]>([]);
  const [tagSuggestions, setTagSuggestions] = useState<VTag[]>([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [loading, setLoading] = useState(false);

  // Filter configuration for Autocomplete
  const filter = createFilterOptions<TagOption>();

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

    const token = Cookies.get("session_token");
    if (!token) {
      router.push("/login");
      return;
    }

    setLoading(true);

    try {
      const response = await fetch(`${config.API_SERVER_PREFIX}/hub/add-post`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          content: postContent,
          tag_ids: selectedTags.map((tag) => tag.id).filter(Boolean),
          new_tags: selectedTags
            .filter((tag) => !tag.id)
            .map((tag) => tag.name),
        }),
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

  return (
    <Paper
      sx={{
        p: 3,
        mb: 4,
        display: "flex",
        flexDirection: "column",
        gap: 2,
      }}
    >
      <Typography variant="h6">{t("posts.compose")}</Typography>
      <TextField
        fullWidth
        multiline
        rows={4}
        value={postContent}
        onChange={(e) => setPostContent(e.target.value)}
        placeholder={t("posts.placeholder")}
        variant="outlined"
      />

      <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
        <Box
          sx={{
            display: "flex",
            flexWrap: "wrap",
            gap: 1,
            alignItems: "center",
          }}
        >
          {selectedTags.map((tag) => (
            <Chip
              key={tag.id || tag.name}
              label={tag.name}
              onDelete={() => {
                setSelectedTags(selectedTags.filter((t) => t !== tag));
              }}
              color="primary"
              icon={<LocalOfferIcon />}
              size="small"
            />
          ))}
        </Box>

        {selectedTags.length < 3 && (
          <Autocomplete
            key={`tag-input-${selectedTags.length}`}
            id="tags-search"
            options={tagSuggestions}
            freeSolo
            value={null}
            inputValue={searchQuery}
            onInputChange={(_, newInputValue) => {
              setSearchQuery(newInputValue);
            }}
            onChange={(_, newValue) => {
              if (!newValue) return;

              // Don't add if already at max
              if (selectedTags.length >= 3) return;

              // Create the new tag
              let newTag: VTag;

              if (typeof newValue === "string") {
                // String input - create tag from the string
                newTag = { name: newValue, id: "" };
              } else if (
                typeof (newValue as TagOption).inputValue === "string"
              ) {
                // Create tag from inputValue
                newTag = {
                  name: (newValue as TagOption).inputValue as string,
                  id: "",
                };
              } else {
                // Existing tag
                newTag = newValue as VTag;
              }

              // Don't add if already selected
              if (selectedTags.some((tag) => tag.name === newTag.name)) {
                return;
              }

              setSelectedTags([...selectedTags, newTag]);
              setSearchQuery("");
            }}
            filterOptions={(options, params) => {
              const filtered = filter(options as TagOption[], params);
              const { inputValue } = params;

              // Only suggest creating a new tag if it's not already in suggestions
              // and not already selected and not empty
              const isExisting = options.some(
                (option) => option.name === inputValue
              );
              const isSelected = selectedTags.some(
                (tag) => tag.name === inputValue
              );

              if (inputValue !== "" && !isExisting && !isSelected) {
                filtered.push({
                  inputValue,
                  name: inputValue,
                  id: "",
                } as TagOption);
              }

              return filtered;
            }}
            getOptionLabel={(option) => {
              if (typeof option === "string") {
                return option;
              }
              return option.name;
            }}
            renderOption={(props, option) => {
              const { key, ...otherProps } = props;
              return (
                <li key={key || option.id || option.name} {...otherProps}>
                  {!option.id ? (
                    <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                      <AddIcon fontSize="small" />
                      <span>{t("posts.newTag", { name: option.name })}</span>
                    </Box>
                  ) : (
                    <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                      <LocalOfferIcon fontSize="small" />
                      <span>{option.name}</span>
                    </Box>
                  )}
                </li>
              );
            }}
            renderInput={(params) => (
              <TextField
                {...params}
                size="small"
                placeholder={
                  selectedTags.length >= 3
                    ? t("posts.maxTags")
                    : t("posts.searchTags")
                }
                disabled={selectedTags.length >= 3}
                InputProps={{
                  ...params.InputProps,
                  startAdornment: <LocalOfferIcon sx={{ mr: 1 }} />,
                }}
              />
            )}
          />
        )}
      </Box>

      <Box sx={{ display: "flex", justifyContent: "flex-end" }}>
        <Button
          variant="contained"
          onClick={handlePublish}
          disabled={loading || !postContent.trim()}
        >
          {loading ? (
            <CircularProgress size={24} color="inherit" />
          ) : (
            t("posts.publish")
          )}
        </Button>
      </Box>
    </Paper>
  );
}
