"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import AddIcon from "@mui/icons-material/Add";
import CloseIcon from "@mui/icons-material/Close";
import LocalOfferIcon from "@mui/icons-material/LocalOffer";
import Autocomplete, { createFilterOptions } from "@mui/material/Autocomplete";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Chip from "@mui/material/Chip";
import CircularProgress from "@mui/material/CircularProgress";
import Paper from "@mui/material/Paper";
import Snackbar from "@mui/material/Snackbar";
import Tab from "@mui/material/Tab";
import Tabs from "@mui/material/Tabs";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import { VTag } from "@vetchium/typespec";
import Cookies from "js-cookie";
import { Suspense, useEffect, useState } from "react";

// Interface for tags including free input
interface TagOption extends VTag {
  inputValue?: string;
}

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

// Filter configuration for Autocomplete
const filter = createFilterOptions<TagOption>();

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`posts-tabpanel-${index}`}
      aria-labelledby={`posts-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
    </div>
  );
}

function PostsContent() {
  const { t } = useTranslation();
  const [postContent, setPostContent] = useState("");
  const [selectedTags, setSelectedTags] = useState<VTag[]>([]);
  const [tagSuggestions, setTagSuggestions] = useState<VTag[]>([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [tabValue, setTabValue] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  // Fetch tag suggestions when user types
  useEffect(() => {
    const fetchTags = async () => {
      if (searchQuery.length >= 2) {
        const token = Cookies.get("session_token");
        if (!token) return;

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
              setError(t("common.error.sessionExpired"));
              Cookies.remove("session_token", { path: "/" });
              return;
            }
            throw new Error(`Failed to fetch tags: ${response.statusText}`);
          }

          const data = await response.json();
          setTagSuggestions(Array.isArray(data) ? data : []);
        } catch (error) {
          console.error("Error fetching tags:", error);
          setError(t("posts.error.tagsFailed"));
          setTagSuggestions([]);
        }
      } else {
        setTagSuggestions([]);
      }
    };

    const debounceTimer = setTimeout(fetchTags, 300);
    return () => clearTimeout(debounceTimer);
  }, [searchQuery, t]);

  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  const handlePublish = async () => {
    if (!postContent.trim()) {
      setError(t("posts.error.contentRequired"));
      return;
    }

    const token = Cookies.get("session_token");
    if (!token) {
      setError(t("common.error.sessionExpired"));
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await fetch(`${config.API_SERVER_PREFIX}/hub/add-post`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          content: postContent,
          vtag_ids: selectedTags.map((tag) => tag.id).filter(Boolean),
        }),
      });

      if (!response.ok) {
        if (response.status === 401) {
          setError(t("common.error.sessionExpired"));
          Cookies.remove("session_token", { path: "/" });
          return;
        }
        throw new Error(`Failed to create post: ${response.statusText}`);
      }

      // Reset form
      setPostContent("");
      setSelectedTags([]);
      setSuccess(t("posts.success"));
    } catch (error) {
      console.error("Error creating post:", error);
      setError(t("posts.error.createFailed"));
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box sx={{ maxWidth: 800, mx: "auto", mt: 4 }}>
      <Typography variant="h4" gutterBottom align="center">
        {t("posts.title")}
      </Typography>

      {/* Compose Section */}
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
              id="tags-search"
              options={tagSuggestions}
              freeSolo
              value={null}
              onChange={(_, newValue) => {
                if (!newValue) return;

                // Don't add if already at max
                if (selectedTags.length >= 3) return;

                // Don't add if already selected
                if (
                  typeof newValue === "string" &&
                  selectedTags.some((tag) => tag.name === newValue)
                )
                  return;
                if (
                  typeof newValue !== "string" &&
                  selectedTags.some(
                    (tag) =>
                      tag.id === newValue.id || tag.name === newValue.name
                  )
                )
                  return;

                const newTag =
                  typeof newValue === "string"
                    ? { name: newValue, id: "" }
                    : newValue;

                setSelectedTags([...selectedTags, newTag]);
                setSearchQuery("");
              }}
              filterOptions={(options, params) => {
                const filtered = filter(options, params);
                const { inputValue } = params;

                // Only suggest creating a new tag if it's not already in suggestions
                // and not already selected
                const isExisting = options.some(
                  (option) => option.name === inputValue
                );
                const isSelected = selectedTags.some(
                  (tag) => tag.name === inputValue
                );

                if (inputValue !== "" && !isExisting && !isSelected) {
                  filtered.push({
                    name: inputValue,
                    id: "", // Empty ID indicates it's a new tag
                  });
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
                      <Box
                        sx={{ display: "flex", alignItems: "center", gap: 1 }}
                      >
                        <AddIcon fontSize="small" />
                        <span>{t("posts.newTag", { name: option.name })}</span>
                      </Box>
                    ) : (
                      <Box
                        sx={{ display: "flex", alignItems: "center", gap: 1 }}
                      >
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
                  onChange={(e) => setSearchQuery(e.target.value)}
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

      {/* Tabs Section */}
      <Paper sx={{ width: "100%" }}>
        <Box sx={{ borderBottom: 1, borderColor: "divider" }}>
          <Tabs
            value={tabValue}
            onChange={handleTabChange}
            aria-label="posts tabs"
            variant="fullWidth"
          >
            <Tab label={t("posts.following")} id="posts-tab-0" />
            <Tab label={t("posts.trending")} id="posts-tab-1" />
          </Tabs>
        </Box>
        <TabPanel value={tabValue} index={0}>
          <Box
            sx={{
              minHeight: 200,
              display: "flex",
              justifyContent: "center",
              alignItems: "center",
            }}
          >
            <Typography variant="body1" color="text.secondary">
              Following posts will appear here
            </Typography>
          </Box>
        </TabPanel>
        <TabPanel value={tabValue} index={1}>
          <Box
            sx={{
              minHeight: 200,
              display: "flex",
              justifyContent: "center",
              alignItems: "center",
            }}
          >
            <Typography variant="body1" color="text.secondary">
              Trending posts will appear here
            </Typography>
          </Box>
        </TabPanel>
      </Paper>

      {/* Notifications */}
      <Snackbar
        open={!!error}
        autoHideDuration={6000}
        onClose={() => setError(null)}
        message={error}
        action={
          <Button color="inherit" size="small" onClick={() => setError(null)}>
            <CloseIcon fontSize="small" />
          </Button>
        }
      />

      <Snackbar
        open={!!success}
        autoHideDuration={6000}
        onClose={() => setSuccess(null)}
        message={success}
        action={
          <Button color="inherit" size="small" onClick={() => setSuccess(null)}>
            <CloseIcon fontSize="small" />
          </Button>
        }
      />
    </Box>
  );
}

export default function PostsPage() {
  return (
    <AuthenticatedLayout>
      <Suspense
        fallback={
          <Box
            sx={{
              display: "flex",
              justifyContent: "center",
              alignItems: "center",
              minHeight: "50vh",
            }}
          >
            <CircularProgress />
          </Box>
        }
      >
        <PostsContent />
      </Suspense>
    </AuthenticatedLayout>
  );
}
