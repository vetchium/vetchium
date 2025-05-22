"use client";

import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import {
  Box,
  Button,
  CircularProgress,
  Paper,
  Typography,
} from "@mui/material";
import { EmployerPost, Post } from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useRef, useState } from "react";
import EmployerPostCard from "./EmployerPostCard";
import PostCard from "./PostCard";

interface TimelineTabProps {
  refreshTrigger?: number;
  type: "following" | "trending";
}

export default function TimelineTab({
  refreshTrigger = 0,
  type,
}: TimelineTabProps) {
  const { t } = useTranslation();
  const router = useRouter();
  const [posts, setPosts] = useState<Post[]>([]);
  const [employerPosts, setEmployerPosts] = useState<EmployerPost[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [paginationKey, setPaginationKey] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState(true);

  // Keep track of active fetches to prevent duplicate calls
  const isFetchingRef = useRef(false);
  const tRef = useRef(t);
  const urlRef = useRef(`${config.API_SERVER_PREFIX}/hub/get-my-home-timeline`);

  // Track if initial load has happened to prevent multiple calls
  const initialLoadDoneRef = useRef(false);

  // Update refs when dependencies change
  useEffect(() => {
    tRef.current = t;
  }, [t]);

  // Simple fetch posts function
  const fetchPosts = useCallback(
    async (refresh = false) => {
      // Currently only supporting following timeline
      if (type !== "following") {
        return;
      }

      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        return;
      }

      // Don't fetch if already loading
      if (loading || isFetchingRef.current) {
        return;
      }

      // Don't fetch if we've reached the end and not refreshing
      if (!refresh && paginationKey === "") {
        return;
      }

      isFetchingRef.current = true;
      setLoading(true);
      setError(null);

      try {
        // Create the request payload
        const requestPayload = {
          pagination_key: refresh ? undefined : paginationKey,
          limit: 10,
        };

        const response = await fetch(urlRef.current, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(requestPayload),
        });

        if (!response.ok) {
          if (response.status === 401) {
            Cookies.remove("session_token", { path: "/" });
            router.push("/login");
            return;
          }
          if (response.status === 422) {
            setHasMore(false);
            setLoading(false);
            isFetchingRef.current = false;
            return;
          }
          throw new Error(`Failed to fetch timeline: ${response.statusText}`);
        }

        const data = await response.json();

        // Ensure tags are always arrays in each post
        const safePosts = (data.posts || []).map((post: Post) => ({
          ...post,
          tags: Array.isArray(post.tags) ? post.tags : [],
        }));

        // Ensure tags are always arrays in each employer post
        const safeEmployerPosts = (data.employer_posts || []).map(
          (post: EmployerPost) => ({
            ...post,
            tags: Array.isArray(post.tags) ? post.tags : [],
          })
        );

        setPosts((prevPosts) => {
          if (refresh) {
            return safePosts;
          } else {
            // Filter out any posts that we already have to prevent duplicates
            const existingIds = new Set(prevPosts.map((post: Post) => post.id));
            const newPosts = safePosts.filter(
              (post: Post) => !existingIds.has(post.id)
            );

            // If we received posts but all were duplicates, we've likely reached the end
            if (safePosts.length > 0 && newPosts.length === 0) {
              setHasMore(false);
              return prevPosts;
            }

            return [...prevPosts, ...newPosts];
          }
        });

        setEmployerPosts((prevPosts) => {
          if (refresh) {
            return safeEmployerPosts;
          } else {
            // Filter out any posts that we already have to prevent duplicates
            const existingIds = new Set(
              prevPosts.map((post: EmployerPost) => post.id)
            );
            const newPosts = safeEmployerPosts.filter(
              (post: EmployerPost) => !existingIds.has(post.id)
            );

            return [...prevPosts, ...newPosts];
          }
        });

        // Handle empty pagination key as end of data according to API spec
        if (data.pagination_key === "") {
          setPaginationKey("");
          setHasMore(false);
        } else {
          setPaginationKey(data.pagination_key);
          setHasMore(
            (data.posts && data.posts.length > 0) ||
              (data.employer_posts && data.employer_posts.length > 0)
          );
        }
      } catch (error) {
        console.error("Error fetching timeline:", error);
        setError(tRef.current("posts.error.fetchFailed"));
      } finally {
        setLoading(false);
        isFetchingRef.current = false;
      }
    },
    // Remove loading from the dependency array to break the infinite loop
    [type, router, paginationKey]
  );

  // Initial load effect (runs only once)
  useEffect(() => {
    if (!initialLoadDoneRef.current) {
      initialLoadDoneRef.current = true;
      fetchPosts(true);
    }
  }, []); // Empty dependency array ensures it only runs once on mount

  // Handle refresh trigger separately
  useEffect(() => {
    if (refreshTrigger > 0 && initialLoadDoneRef.current) {
      fetchPosts(true);
    }
  }, [refreshTrigger, fetchPosts]);

  // Handle load more button click
  const loadMore = () => {
    if (!loading && hasMore && paginationKey !== "") {
      fetchPosts(false);
    }
  };

  if (error) {
    return (
      <Paper
        elevation={0}
        sx={{
          p: 3,
          textAlign: "center",
          borderRadius: "12px",
          backgroundColor: "#fff0f0",
          border: "1px solid #fcc",
        }}
      >
        <Typography color="error" variant="body1" sx={{ fontWeight: 500 }}>
          {error}
        </Typography>
        <Button
          sx={{
            mt: 2,
            fontWeight: 600,
            textTransform: "none",
            padding: "6px 16px",
          }}
          variant="contained"
          color="primary"
          onClick={() => fetchPosts(true)}
        >
          {t("common.retry")}
        </Button>
      </Paper>
    );
  }

  if (type === "trending") {
    return (
      <Paper
        elevation={0}
        sx={{
          minHeight: 200,
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          borderRadius: "12px",
          backgroundColor: "#f8f9fa",
          border: "1px solid #e9ecef",
        }}
      >
        <Typography
          variant="body1"
          color="text.secondary"
          sx={{ fontWeight: 500 }}
        >
          {t("posts.trendingComingSoon")}
        </Typography>
      </Paper>
    );
  }

  // Check if there are no posts to show
  const hasNoPosts =
    !error && posts.length === 0 && employerPosts.length === 0 && !loading;

  if (hasNoPosts) {
    return (
      <Paper
        elevation={0}
        sx={{
          minHeight: 200,
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          flexDirection: "column",
          gap: 2,
          p: 3,
          borderRadius: "12px",
          backgroundColor: "#f8f9fa",
          border: "1px solid #e9ecef",
        }}
      >
        <Typography
          variant="body1"
          color="text.secondary"
          sx={{ fontWeight: 500 }}
        >
          {t("posts.noTimelinePosts")}
        </Typography>
      </Paper>
    );
  }

  // Helper function to get timestamp for sorting
  const getTimestamp = (item: Post | EmployerPost): number => {
    if ("updated_at" in item && item.updated_at) {
      return new Date(item.updated_at).getTime();
    }
    return new Date(item.created_at).getTime();
  };

  // Create a combined timeline of user posts and employer posts
  // Sort by timestamp for proper chronological display
  const combinedTimeline = [...posts, ...employerPosts].sort((a, b) => {
    return getTimestamp(b) - getTimestamp(a);
  });

  // Helper function to determine if an item is an EmployerPost
  const isEmployerPost = (item: Post | EmployerPost): item is EmployerPost => {
    return "employer_name" in item && "employer_domain_name" in item;
  };

  return (
    <Box
      sx={{
        p: { xs: 1, sm: 2 },
        backgroundColor: "#f5f7f9",
        borderRadius: "8px",
      }}
    >
      <Box sx={{ maxWidth: "100%", mx: "auto" }}>
        {/* Render combined timeline */}
        {combinedTimeline.map((item, index) => (
          <Box
            key={`timeline-${item.id}-${index}`}
            sx={{
              mb: 2.5,
              "&:last-child": {
                mb: 0,
              },
            }}
          >
            {isEmployerPost(item) ? (
              <EmployerPostCard post={item} />
            ) : (
              <PostCard post={item} />
            )}
          </Box>
        ))}
      </Box>

      {loading && (
        <Box sx={{ textAlign: "center", p: 4 }}>
          <CircularProgress size={36} />
        </Box>
      )}

      {!loading && hasMore && paginationKey !== "" && (
        <Box sx={{ textAlign: "center", p: 2 }}>
          <Button
            variant="outlined"
            onClick={loadMore}
            sx={{
              fontWeight: 600,
              textTransform: "none",
              "&:hover": {
                backgroundColor: "primary.lighter",
              },
              padding: "8px 20px",
            }}
            color="primary"
          >
            {t("posts.loadMore")}
          </Button>
        </Box>
      )}

      {!loading && !hasMore && (
        <Box sx={{ textAlign: "center", p: 2 }}>
          <Typography
            variant="body2"
            color="text.secondary"
            sx={{
              fontSize: "0.875rem",
              fontWeight: 500,
            }}
          >
            {t("posts.noMorePosts")}
          </Typography>
        </Box>
      )}
    </Box>
  );
}
