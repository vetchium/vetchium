"use client";

import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import { Box, Button, CircularProgress, Typography } from "@mui/material";
import { Post } from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useRef, useState } from "react";
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

        // Handle empty pagination key as end of data according to API spec
        if (data.pagination_key === "") {
          setPaginationKey("");
          setHasMore(false);
        } else {
          setPaginationKey(data.pagination_key);
          setHasMore(data.posts && data.posts.length > 0);
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
      <Box sx={{ p: 2, textAlign: "center" }}>
        <Typography color="error">{error}</Typography>
        <Button
          sx={{ mt: 1 }}
          variant="outlined"
          onClick={() => fetchPosts(true)}
        >
          {t("common.retry")}
        </Button>
      </Box>
    );
  }

  if (type === "trending") {
    return (
      <Box
        sx={{
          minHeight: 200,
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
        }}
      >
        <Typography variant="body1" color="text.secondary">
          {t("posts.trendingComingSoon")}
        </Typography>
      </Box>
    );
  }

  if (!error && posts.length === 0 && !loading) {
    return (
      <Box
        sx={{
          minHeight: 200,
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          flexDirection: "column",
          gap: 2,
          p: 3,
        }}
      >
        <Typography variant="body1" color="text.secondary">
          {t("posts.noTimelinePosts")}
        </Typography>
      </Box>
    );
  }

  return (
    <Box sx={{ p: 2 }}>
      {posts.map((post, index) => (
        <div key={`${post.id}-${index}`}>
          <PostCard post={post} />
        </div>
      ))}

      {loading && (
        <Box sx={{ textAlign: "center", p: 3 }}>
          <CircularProgress size={40} />
        </Box>
      )}

      {!loading && hasMore && paginationKey !== "" && (
        <Box sx={{ textAlign: "center", p: 2 }}>
          <Button variant="outlined" onClick={loadMore}>
            {t("posts.loadMore")}
          </Button>
        </Box>
      )}

      {!loading && !hasMore && (
        <Box sx={{ textAlign: "center", p: 2 }}>
          <Typography variant="body2" color="text.secondary">
            {t("posts.noMorePosts")}
          </Typography>
        </Box>
      )}
    </Box>
  );
}
