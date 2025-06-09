"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import { useAuth } from "@/hooks/useAuth";
import { useIncognitoPosts } from "@/hooks/useIncognitoPosts";
import { useTranslation } from "@/hooks/useTranslation";
import CloseIcon from "@mui/icons-material/Close";
import {
  Alert,
  Box,
  Button,
  CircularProgress,
  Container,
  FormControl,
  InputLabel,
  MenuItem,
  Paper,
  Select,
  Snackbar,
  Typography,
} from "@mui/material";
import { IncognitoPostTimeFilter, VTagID } from "@vetchium/typespec";
import { useParams, useRouter } from "next/navigation";
import { Suspense, useEffect, useState } from "react";
import CreatePostDialog from "../../components/CreatePostDialog";
import IncognitoPostCard from "../../components/IncognitoPostCard";

function IncognitoPostsByTagContent() {
  const { t } = useTranslation();
  const router = useRouter();
  const params = useParams();
  useAuth();

  const tagId = params?.tagId as VTagID;

  // Handle case where tagId is not available
  if (!tagId) {
    return (
      <Container maxWidth="lg" sx={{ py: 3 }}>
        <Alert severity="error">Invalid tag ID</Alert>
      </Container>
    );
  }

  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [timeFilter, setTimeFilter] = useState<IncognitoPostTimeFilter>(
    IncognitoPostTimeFilter.Past24Hours
  );

  const {
    posts,
    isLoading,
    error: postsError,
    hasMorePages,
    loadPosts,
    loadMorePosts,
  } = useIncognitoPosts();

  useEffect(() => {
    if (postsError) {
      setError(postsError.message || t("incognitoPosts.errors.loadFailed"));
    }
  }, [postsError, t]);

  useEffect(() => {
    if (tagId) {
      loadPosts(tagId, timeFilter, true);
    }
  }, [tagId, timeFilter, loadPosts]);

  const handleTimeFilterChange = (newFilter: IncognitoPostTimeFilter) => {
    setTimeFilter(newFilter);
    if (tagId) {
      loadPosts(tagId, newFilter, true);
    }
  };

  const handleLoadMore = () => {
    loadMorePosts();
  };

  const handlePostCreated = () => {
    setSuccess(t("incognitoPosts.success.postCreated"));
    // Refresh posts after creation
    if (tagId) {
      loadPosts(tagId, timeFilter, true);
    }
  };

  const handleError = (errorMessage: string) => {
    setError(errorMessage);
  };

  const handlePostDeleted = () => {
    if (tagId) {
      loadPosts(tagId, timeFilter, true);
    }
  };

  return (
    <Container maxWidth="lg" sx={{ py: 3 }}>
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom align="center">
          {t("incognitoPosts.title")}
        </Typography>
        <Typography
          variant="body1"
          color="text.secondary"
          align="center"
          sx={{ mb: 3 }}
        >
          {t("incognitoPosts.description")}
        </Typography>

        {/* Create Post Button */}
        <Box sx={{ display: "flex", justifyContent: "center", mb: 3 }}>
          <Button
            variant="contained"
            color="primary"
            onClick={() => setCreateDialogOpen(true)}
            sx={{ minWidth: 200 }}
          >
            {t("incognitoPosts.browsing.createPost")}
          </Button>
        </Box>
      </Box>

      <Paper sx={{ width: "100%" }}>
        <Box sx={{ p: 3 }}>
          {/* Header with tag info and filters */}
          <Box
            sx={{
              mb: 3,
              display: "flex",
              gap: 2,
              flexWrap: "wrap",
              alignItems: "center",
            }}
          >
            <Typography variant="h6">Posts tagged: {tagId}</Typography>

            <Button
              variant="outlined"
              onClick={() => router.push("/incognito-posts")}
              sx={{ ml: "auto" }}
            >
              ← Browse All Tags
            </Button>
          </Box>

          {/* Time Filter */}
          <Box sx={{ mb: 3 }}>
            <FormControl sx={{ minWidth: 200 }}>
              <InputLabel>{t("incognitoPosts.filters.timeFilter")}</InputLabel>
              <Select
                value={timeFilter}
                label={t("incognitoPosts.filters.timeFilter")}
                onChange={(e) =>
                  handleTimeFilterChange(
                    e.target.value as IncognitoPostTimeFilter
                  )
                }
              >
                <MenuItem value={IncognitoPostTimeFilter.Past24Hours}>
                  {t("incognitoPosts.timeFilters.past24Hours")}
                </MenuItem>
                <MenuItem value={IncognitoPostTimeFilter.PastWeek}>
                  {t("incognitoPosts.timeFilters.pastWeek")}
                </MenuItem>
                <MenuItem value={IncognitoPostTimeFilter.PastMonth}>
                  {t("incognitoPosts.timeFilters.pastMonth")}
                </MenuItem>
                <MenuItem value={IncognitoPostTimeFilter.PastYear}>
                  {t("incognitoPosts.timeFilters.pastYear")}
                </MenuItem>
              </Select>
            </FormControl>
          </Box>

          {/* Posts List */}
          {posts.length === 0 && !isLoading ? (
            <Alert severity="info">
              {t("incognitoPosts.feed.noPostsForTag")}
            </Alert>
          ) : (
            <Box sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
              {posts.map((post) => (
                <IncognitoPostCard
                  key={post.incognito_post_id}
                  post={post}
                  onDeleted={handlePostDeleted}
                  onVoteUpdated={() => {}}
                  onError={handleError}
                />
              ))}
            </Box>
          )}

          {/* Loading indicator */}
          {isLoading && (
            <Box sx={{ display: "flex", justifyContent: "center", mt: 2 }}>
              <CircularProgress />
            </Box>
          )}

          {/* Load More Button */}
          {hasMorePages && !isLoading && (
            <Box sx={{ display: "flex", justifyContent: "center", mt: 3 }}>
              <Button variant="outlined" onClick={handleLoadMore}>
                {t("incognitoPosts.feed.loadMore")}
              </Button>
            </Box>
          )}

          {/* No More Posts Message */}
          {!hasMorePages && posts.length > 0 && !isLoading && (
            <Typography
              variant="body2"
              color="text.secondary"
              align="center"
              sx={{ mt: 3 }}
            >
              {t("incognitoPosts.feed.noMorePosts")}
            </Typography>
          )}
        </Box>
      </Paper>

      {/* Create Post Dialog */}
      <CreatePostDialog
        open={createDialogOpen}
        onClose={() => setCreateDialogOpen(false)}
        onPostCreated={handlePostCreated}
        onError={handleError}
      />

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
    </Container>
  );
}

export default function IncognitoPostsByTagPage() {
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
        <IncognitoPostsByTagContent />
      </Suspense>
    </AuthenticatedLayout>
  );
}
