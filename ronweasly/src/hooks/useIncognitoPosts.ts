"use client";

import { config } from "@/config";
import {
  AddIncognitoPostRequest,
  AddIncognitoPostResponse,
  GetIncognitoPostsRequest,
  GetIncognitoPostsResponse,
  GetMyIncognitoPostsRequest,
  GetMyIncognitoPostsResponse,
  IncognitoPostTimeFilter,
  VTagID,
} from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useCallback, useState } from "react";

interface UseIncognitoPostsResult {
  posts: GetIncognitoPostsResponse["posts"];
  isLoading: boolean;
  error: Error | null;
  hasMorePages: boolean;
  loadPosts: (
    tagId: VTagID,
    timeFilter?: IncognitoPostTimeFilter,
    refresh?: boolean
  ) => Promise<void>;
  loadMorePosts: () => Promise<void>;
}

interface UseMyIncognitoPostsResult {
  posts: GetMyIncognitoPostsResponse["posts"];
  isLoading: boolean;
  error: Error | null;
  hasMorePages: boolean;
  loadMyPosts: (refresh?: boolean) => Promise<void>;
  loadMoreMyPosts: () => Promise<void>;
}

interface UseCreateIncognitoPostResult {
  isCreating: boolean;
  error: Error | null;
  createPost: (request: AddIncognitoPostRequest) => Promise<string | null>;
}

export const useIncognitoPosts = (): UseIncognitoPostsResult => {
  const [posts, setPosts] = useState<GetIncognitoPostsResponse["posts"]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const [paginationKey, setPaginationKey] = useState<string | undefined>();
  const [currentTagId, setCurrentTagId] = useState<VTagID>("");
  const [currentTimeFilter, setCurrentTimeFilter] =
    useState<IncognitoPostTimeFilter>(IncognitoPostTimeFilter.Past24Hours);
  const [hasMorePages, setHasMorePages] = useState(false);

  const loadPosts = useCallback(
    async (
      tagId: VTagID,
      timeFilter: IncognitoPostTimeFilter = IncognitoPostTimeFilter.Past24Hours,
      refresh = false
    ) => {
      setIsLoading(true);
      setError(null);

      try {
        const token = Cookies.get("session_token");
        if (!token) {
          throw new Error("User not authenticated");
        }

        const request = new GetIncognitoPostsRequest();
        request.tag_id = tagId;
        request.time_filter = timeFilter;
        request.limit = 25;

        if (!refresh && tagId === currentTagId && paginationKey) {
          request.pagination_key = paginationKey;
        }

        const response = await fetch(
          `${config.API_SERVER_PREFIX}/hub/get-incognito-posts`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify(request),
          }
        );

        if (!response.ok) {
          throw new Error(`Failed to fetch posts: ${response.statusText}`);
        }

        const data: GetIncognitoPostsResponse = await response.json();

        if (refresh || tagId !== currentTagId) {
          setPosts(data.posts);
        } else {
          setPosts((prev) => [...prev, ...data.posts]);
        }

        setPaginationKey(data.pagination_key);
        setCurrentTagId(tagId);
        setCurrentTimeFilter(timeFilter);
        setHasMorePages(!!data.pagination_key);
      } catch (err) {
        setError(
          err instanceof Error ? err : new Error("Unknown error occurred")
        );
      } finally {
        setIsLoading(false);
      }
    },
    [currentTagId, paginationKey]
  );

  const loadMorePosts = useCallback(async () => {
    if (!hasMorePages || isLoading || !currentTagId) return;
    await loadPosts(currentTagId, currentTimeFilter, false);
  }, [hasMorePages, isLoading, currentTagId, currentTimeFilter, loadPosts]);

  return {
    posts,
    isLoading,
    error,
    hasMorePages,
    loadPosts,
    loadMorePosts,
  };
};

export const useMyIncognitoPosts = (): UseMyIncognitoPostsResult => {
  const [posts, setPosts] = useState<GetMyIncognitoPostsResponse["posts"]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const [paginationKey, setPaginationKey] = useState<string | undefined>();
  const [hasMorePages, setHasMorePages] = useState(false);

  const loadMyPosts = useCallback(
    async (refresh = false) => {
      setIsLoading(true);
      setError(null);

      try {
        const token = Cookies.get("session_token");
        if (!token) {
          throw new Error("User not authenticated");
        }

        const request = new GetMyIncognitoPostsRequest();
        request.limit = 25;

        if (!refresh && paginationKey) {
          request.pagination_key = paginationKey;
        }

        const response = await fetch(
          `${config.API_SERVER_PREFIX}/hub/get-my-incognito-posts`,
          {
            method: "GET",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify(request),
          }
        );

        if (!response.ok) {
          throw new Error(`Failed to fetch my posts: ${response.statusText}`);
        }

        const data: GetMyIncognitoPostsResponse = await response.json();

        if (refresh) {
          setPosts(data.posts);
        } else {
          setPosts((prev) => [...prev, ...data.posts]);
        }

        setPaginationKey(data.pagination_key);
        setHasMorePages(!!data.pagination_key);
      } catch (err) {
        setError(
          err instanceof Error ? err : new Error("Unknown error occurred")
        );
      } finally {
        setIsLoading(false);
      }
    },
    [paginationKey]
  );

  const loadMoreMyPosts = useCallback(async () => {
    if (!hasMorePages || isLoading) return;
    await loadMyPosts(false);
  }, [hasMorePages, isLoading, loadMyPosts]);

  return {
    posts,
    isLoading,
    error,
    hasMorePages,
    loadMyPosts,
    loadMoreMyPosts,
  };
};

export const useCreateIncognitoPost = (): UseCreateIncognitoPostResult => {
  const [isCreating, setIsCreating] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const createPost = useCallback(
    async (request: AddIncognitoPostRequest): Promise<string | null> => {
      setIsCreating(true);
      setError(null);

      try {
        const token = Cookies.get("session_token");
        if (!token) {
          throw new Error("User not authenticated");
        }

        if (!request.IsValid()) {
          throw new Error("Invalid post data");
        }

        const response = await fetch(
          `${config.API_SERVER_PREFIX}/hub/add-incognito-post`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify(request),
          }
        );

        if (!response.ok) {
          if (response.status === 400) {
            throw new Error("Invalid post content or tags");
          }
          throw new Error(`Failed to create post: ${response.statusText}`);
        }

        const data: AddIncognitoPostResponse = await response.json();
        return data.incognito_post_id;
      } catch (err) {
        const error =
          err instanceof Error ? err : new Error("Unknown error occurred");
        setError(error);
        return null;
      } finally {
        setIsCreating(false);
      }
    },
    []
  );

  return {
    isCreating,
    error,
    createPost,
  };
};
