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
  const observerRef = useRef<IntersectionObserver | null>(null);
  const lastPostElementRef = useRef<HTMLDivElement | null>(null);

  // Keep track of active fetches to prevent duplicate calls
  const isFetchingRef = useRef(false);
  const tRef = useRef(t);
  const urlRef = useRef(`${config.API_SERVER_PREFIX}/hub/get-my-home-timeline`);

  // Update refs when dependencies change
  useEffect(() => {
    tRef.current = t;
  }, [t]);

  // Fetch the right timeline data based on type
  // eslint-disable-next-line react-hooks/exhaustive-deps
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

      // Check if we should proceed with the fetch
      if (loading || isFetchingRef.current || (!hasMore && !refresh)) {
        return;
      }

      // If the pagination key is an empty string and not a refresh, don't fetch
      if (!refresh && paginationKey === "") {
        return;
      }

      isFetchingRef.current = true;
      setLoading(true);
      setError(null);

      try {
        const response = await fetch(urlRef.current, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            pagination_key: refresh ? undefined : paginationKey,
            limit: 10,
          }),
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

        setPosts((prevPosts) =>
          refresh ? safePosts : [...prevPosts, ...safePosts]
        );

        // Handle empty pagination key as end of data
        if (data.pagination_key === "") {
          setPaginationKey("");
          setHasMore(false);
        } else {
          setPaginationKey(data.pagination_key || null);
          // Only set hasMore to true if we have both posts and a non-empty pagination key
          setHasMore(
            safePosts.length > 0 &&
              !!data.pagination_key &&
              data.pagination_key !== ""
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
    // Remove all reactive dependencies that aren't necessary
    [type, router]
  );

  // Refresh timeline when refreshTrigger changes or on initial load
  useEffect(() => {
    fetchPosts(true);
  }, [refreshTrigger, fetchPosts]);

  // Setup intersection observer for infinite scrolling
  useEffect(() => {
    // Don't setup observer if already loading, no posts, or no more data
    if (loading || posts.length === 0 || !hasMore || paginationKey === "") {
      return;
    }

    // Cleanup previous observer
    if (observerRef.current) {
      observerRef.current.disconnect();
    }

    const callback = (entries: IntersectionObserverEntry[]) => {
      // Only fetch if entry is intersecting, we have more data, and not currently loading
      if (
        entries[0]?.isIntersecting &&
        hasMore &&
        !loading &&
        !isFetchingRef.current &&
        paginationKey !== ""
      ) {
        fetchPosts(false);
      }
    };

    observerRef.current = new IntersectionObserver(callback, {
      rootMargin: "100px",
      threshold: 0.1,
    });

    // Only observe if we have a last post element
    if (lastPostElementRef.current) {
      observerRef.current.observe(lastPostElementRef.current);
    }

    return () => {
      if (observerRef.current) {
        observerRef.current.disconnect();
      }
    };
  }, [loading, hasMore, posts.length, fetchPosts, paginationKey]);

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
        <div
          key={post.id}
          ref={index === posts.length - 1 ? lastPostElementRef : null}
        >
          <PostCard post={post} />
        </div>
      ))}

      {loading && (
        <Box sx={{ textAlign: "center", p: 3 }}>
          <CircularProgress size={40} />
        </Box>
      )}
    </Box>
  );
}
