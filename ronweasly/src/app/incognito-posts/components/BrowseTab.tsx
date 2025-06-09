"use client";

import { useIncognitoPosts } from "@/hooks/useIncognitoPosts";
import { useTranslation } from "@/hooks/useTranslation";
import {
  Alert,
  Box,
  Button,
  CircularProgress,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  Typography,
} from "@mui/material";
import { IncognitoPostTimeFilter, VTagID } from "@vetchium/typespec";
import { useEffect, useState } from "react";
import IncognitoPostCard from "./IncognitoPostCard";
import TagSelector from "./TagSelector";

interface BrowseTabProps {
  refreshTrigger: number;
  onError: (error: string) => void;
  onSuccess: (message: string) => void;
}

export default function BrowseTab({
  refreshTrigger,
  onError,
  onSuccess,
}: BrowseTabProps) {
  const { t } = useTranslation();
  const { posts, isLoading, error, hasMorePages, loadPosts, loadMorePosts } =
    useIncognitoPosts();

  const [selectedTag, setSelectedTag] = useState<VTagID>("");
  const [timeFilter, setTimeFilter] = useState<IncognitoPostTimeFilter>(
    IncognitoPostTimeFilter.Past24Hours
  );

  useEffect(() => {
    if (error) {
      onError(error.message || t("incognitoPosts.errors.loadFailed"));
    }
  }, [error, onError, t]);

  useEffect(() => {
    if (selectedTag && refreshTrigger > 0) {
      loadPosts(selectedTag, timeFilter, true);
    }
  }, [refreshTrigger, selectedTag, timeFilter, loadPosts]);

  const handleTagSelect = (tagId: VTagID) => {
    setSelectedTag(tagId);
    if (tagId) {
      loadPosts(tagId, timeFilter, true);
    }
  };

  const handleTimeFilterChange = (newFilter: IncognitoPostTimeFilter) => {
    setTimeFilter(newFilter);
    if (selectedTag) {
      loadPosts(selectedTag, newFilter, true);
    }
  };

  const handleLoadMore = () => {
    loadMorePosts();
  };

  const handlePostDeleted = () => {
    if (selectedTag) {
      loadPosts(selectedTag, timeFilter, true);
    }
    onSuccess(t("incognitoPosts.success.postDeleted"));
  };

  const handleVoteUpdated = () => {
    onSuccess(t("incognitoPosts.success.voteUpdated"));
  };

  return (
    <Box>
      {/* Filters */}
      <Box sx={{ mb: 3, display: "flex", gap: 2, flexWrap: "wrap" }}>
        <Box sx={{ minWidth: 200 }}>
          <TagSelector
            selectedTag={selectedTag}
            onTagSelect={handleTagSelect}
            onError={onError}
          />
        </Box>

        <FormControl sx={{ minWidth: 150 }}>
          <InputLabel>{t("incognitoPosts.filters.timeFilter")}</InputLabel>
          <Select
            value={timeFilter}
            label={t("incognitoPosts.filters.timeFilter")}
            onChange={(e) =>
              handleTimeFilterChange(e.target.value as IncognitoPostTimeFilter)
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

      {/* Content */}
      {!selectedTag ? (
        <Alert severity="info">{t("incognitoPosts.feed.selectTagFirst")}</Alert>
      ) : (
        <Box>
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
                  onVoteUpdated={handleVoteUpdated}
                  onError={onError}
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
      )}
    </Box>
  );
}
