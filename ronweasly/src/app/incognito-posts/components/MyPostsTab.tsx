"use client";

import { useMyIncognitoPosts } from "@/hooks/useIncognitoPosts";
import { useTranslation } from "@/hooks/useTranslation";
import {
  Alert,
  Box,
  Button,
  CircularProgress,
  Typography,
} from "@mui/material";
import { useEffect } from "react";
import IncognitoPostCard from "./IncognitoPostCard";

interface MyPostsTabProps {
  refreshTrigger: number;
  onError: (error: string) => void;
  onSuccess: (message: string) => void;
}

export default function MyPostsTab({
  refreshTrigger,
  onError,
  onSuccess,
}: MyPostsTabProps) {
  const { t } = useTranslation();
  const {
    posts,
    isLoading,
    error,
    hasMorePages,
    loadMyPosts,
    loadMoreMyPosts,
  } = useMyIncognitoPosts();

  useEffect(() => {
    if (error) {
      onError(error.message || t("incognitoPosts.errors.loadFailed"));
    }
  }, [error, onError, t]);

  useEffect(() => {
    loadMyPosts(true);
  }, [refreshTrigger, loadMyPosts]);

  const handleLoadMore = () => {
    loadMoreMyPosts();
  };

  const handlePostDeleted = () => {
    loadMyPosts(true);
    onSuccess(t("incognitoPosts.success.postDeleted"));
  };

  const handleVoteUpdated = () => {
    onSuccess(t("incognitoPosts.success.voteUpdated"));
  };

  return (
    <Box>
      {/* Header */}
      <Box sx={{ mb: 3 }}>
        <Typography variant="h6" gutterBottom>
          {t("incognitoPosts.myPosts.title")}
        </Typography>
        {posts.length > 0 && (
          <Typography variant="body2" color="text.secondary">
            {t("incognitoPosts.myPosts.totalPosts")}: {posts.length}
          </Typography>
        )}
      </Box>

      {/* Content */}
      {posts.length === 0 && !isLoading ? (
        <Alert severity="info">
          <Typography variant="body1" gutterBottom>
            {t("incognitoPosts.myPosts.noPosts")}
          </Typography>
          <Typography variant="body2">
            {t("incognitoPosts.myPosts.createFirst")}
          </Typography>
        </Alert>
      ) : (
        <Box>
          {/* Posts List */}
          <Box sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
            {posts.map((post) => (
              <IncognitoPostCard
                key={post.incognito_post_id}
                post={post}
                onDeleted={handlePostDeleted}
                onVoteUpdated={handleVoteUpdated}
                onError={onError}
              />
            ))}
          </Box>

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
      )}
    </Box>
  );
}
