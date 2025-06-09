"use client";

import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import {
  Alert,
  Box,
  Button,
  Card,
  CardContent,
  Chip,
  CircularProgress,
  Typography,
} from "@mui/material";
import {
  GetMyIncognitoPostCommentsRequest,
  GetMyIncognitoPostCommentsResponse,
  MyIncognitoPostComment,
} from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";

interface MyCommentsTabProps {
  onError: (error: string) => void;
}

export default function MyCommentsTab({ onError }: MyCommentsTabProps) {
  const { t } = useTranslation();
  const router = useRouter();
  const [comments, setComments] = useState<MyIncognitoPostComment[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [paginationKey, setPaginationKey] = useState<string | undefined>();
  const [hasMorePages, setHasMorePages] = useState(false);

  const loadComments = useCallback(
    async (refresh = false) => {
      setIsLoading(true);
      try {
        const token = Cookies.get("session_token");
        if (!token) {
          throw new Error("User not authenticated");
        }

        const request = new GetMyIncognitoPostCommentsRequest();
        request.limit = 25;

        if (!refresh && paginationKey) {
          request.pagination_key = paginationKey;
        }

        const response = await fetch(
          `${config.API_SERVER_PREFIX}/hub/get-my-incognito-post-comments`,
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
          throw new Error(`Failed to fetch comments: ${response.statusText}`);
        }

        const data: GetMyIncognitoPostCommentsResponse = await response.json();

        if (refresh) {
          setComments(data.comments);
        } else {
          setComments((prev) => [...prev, ...data.comments]);
        }

        setPaginationKey(data.pagination_key);
        setHasMorePages(!!data.pagination_key);
      } catch (error) {
        onError(
          error instanceof Error
            ? error.message
            : t("incognitoPosts.errors.commentLoadFailed")
        );
      } finally {
        setIsLoading(false);
      }
    },
    [paginationKey, onError, t]
  );

  useEffect(() => {
    loadComments(true);
  }, [loadComments]);

  const handleLoadMore = () => {
    loadComments(false);
  };

  const handleViewPost = (postId: string) => {
    router.push(`/incognito-posts/${postId}`);
  };

  const formatCount = (count: number, singular: string, plural: string) => {
    return count === 1 ? `${count} ${singular}` : `${count} ${plural}`;
  };

  return (
    <Box>
      {/* Header */}
      <Box sx={{ mb: 3 }}>
        <Typography variant="h6" gutterBottom>
          {t("incognitoPosts.myComments.title")}
        </Typography>
        {comments.length > 0 && (
          <Typography variant="body2" color="text.secondary">
            {t("incognitoPosts.myComments.totalComments")}: {comments.length}
          </Typography>
        )}
      </Box>

      {/* Content */}
      {comments.length === 0 && !isLoading ? (
        <Alert severity="info">
          <Typography variant="body1">
            {t("incognitoPosts.myComments.noComments")}
          </Typography>
        </Alert>
      ) : (
        <Box>
          {/* Comments List */}
          <Box sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
            {comments.map((comment) => (
              <Card key={comment.comment_id} variant="outlined">
                <CardContent>
                  {/* Comment Header */}
                  <Box
                    sx={{
                      display: "flex",
                      justifyContent: "space-between",
                      alignItems: "flex-start",
                      mb: 2,
                    }}
                  >
                    <Box>
                      <Typography variant="body2" color="text.secondary">
                        {new Date(comment.created_at).toLocaleDateString()}
                        {comment.depth > 0 && (
                          <>
                            {" • "}
                            {t("incognitoPosts.comments.replyTo")} (depth:{" "}
                            {comment.depth})
                          </>
                        )}
                      </Typography>
                    </Box>
                    <Button
                      size="small"
                      onClick={() => handleViewPost(comment.incognito_post_id)}
                    >
                      {t("incognitoPosts.myComments.viewPost")}
                    </Button>
                  </Box>

                  {/* Comment Content */}
                  {comment.is_deleted ? (
                    <Typography
                      variant="body2"
                      color="text.secondary"
                      fontStyle="italic"
                      sx={{ mb: 2 }}
                    >
                      {t("incognitoPosts.comments.deleted")}
                    </Typography>
                  ) : (
                    <Typography
                      variant="body1"
                      sx={{ mb: 2, whiteSpace: "pre-wrap" }}
                    >
                      {comment.content}
                    </Typography>
                  )}

                  {/* Vote Stats */}
                  <Box
                    sx={{
                      display: "flex",
                      alignItems: "center",
                      gap: 2,
                      mb: 2,
                    }}
                  >
                    <Typography variant="body2" color="text.secondary">
                      {formatCount(
                        comment.upvotes_count,
                        t("incognitoPosts.post.upvotes"),
                        t("incognitoPosts.post.upvotesPlural")
                      )}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {formatCount(
                        comment.downvotes_count,
                        t("incognitoPosts.post.downvotes"),
                        t("incognitoPosts.post.downvotesPlural")
                      )}
                    </Typography>
                    <Typography
                      variant="body2"
                      sx={{
                        fontWeight: "bold",
                        color:
                          comment.score > 0
                            ? "success.main"
                            : comment.score < 0
                            ? "error.main"
                            : "text.secondary",
                      }}
                    >
                      {t("incognitoPosts.post.score")}: {comment.score}
                    </Typography>
                  </Box>

                  {/* Post Preview */}
                  <Box
                    sx={{ mt: 2, p: 2, bgcolor: "grey.50", borderRadius: 1 }}
                  >
                    <Typography
                      variant="body2"
                      color="text.secondary"
                      gutterBottom
                    >
                      {t("incognitoPosts.myComments.postPreview")}:
                    </Typography>
                    <Typography variant="body2" sx={{ mb: 1 }}>
                      {comment.post_content_preview}
                    </Typography>
                    <Box sx={{ display: "flex", flexWrap: "wrap", gap: 0.5 }}>
                      {comment.post_tags.map((tag) => (
                        <Chip
                          key={tag.id}
                          label={tag.name}
                          size="small"
                          variant="outlined"
                        />
                      ))}
                    </Box>
                  </Box>
                </CardContent>
              </Card>
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

          {/* No More Comments Message */}
          {!hasMorePages && comments.length > 0 && !isLoading && (
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
