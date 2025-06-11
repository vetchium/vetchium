"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import { config } from "@/config";
import { useAuth } from "@/hooks/useAuth";
import { useTranslation } from "@/hooks/useTranslation";
import AddCircleOutlineIcon from "@mui/icons-material/AddCircleOutline";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import CloseIcon from "@mui/icons-material/Close";
import RemoveCircleOutlineIcon from "@mui/icons-material/RemoveCircleOutline";
import {
  Alert,
  Box,
  Button,
  Card,
  CardContent,
  Chip,
  CircularProgress,
  Container,
  IconButton,
  Snackbar,
  TextField,
  Typography,
} from "@mui/material";
import {
  AddIncognitoPostCommentResponse,
  GetIncognitoPostCommentsRequest,
  GetIncognitoPostCommentsResponse,
  GetIncognitoPostRequest,
  IncognitoPost,
  IncognitoPostComment,
} from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useParams, useRouter } from "next/navigation";
import { Suspense, useEffect, useState } from "react";
import VotingButtons from "../components/VotingButtons";

// Define the interface inline since it's missing from typespec
interface AddIncognitoPostCommentRequest {
  incognito_post_id: string;
  content: string;
  in_reply_to?: string;
}

function IncognitoPostDetailsContent() {
  const { t } = useTranslation();
  const router = useRouter();
  const params = useParams();
  useAuth();

  const postId = params?.postId as string;

  const [post, setPost] = useState<IncognitoPost | null>(null);
  const [comments, setComments] = useState<IncognitoPostComment[]>([]);
  const [isLoadingPost, setIsLoadingPost] = useState(true);
  const [isLoadingComments, setIsLoadingComments] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [newComment, setNewComment] = useState("");
  const [isSubmittingComment, setIsSubmittingComment] = useState(false);
  const [replyToComment, setReplyToComment] = useState<string | null>(null);
  const [commentsPaginationKey, setCommentsPaginationKey] = useState<
    string | undefined
  >();
  const [hasMoreComments, setHasMoreComments] = useState(false);
  const [collapsedComments, setCollapsedComments] = useState<Set<string>>(
    new Set()
  );

  // Handle case where postId is not available
  if (!postId) {
    return (
      <Container maxWidth="lg" sx={{ py: 3 }}>
        <Alert severity="error">Invalid post ID</Alert>
      </Container>
    );
  }

  useEffect(() => {
    loadPost();
    loadComments(true);
  }, [postId]);

  const loadPost = async () => {
    setIsLoadingPost(true);
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        throw new Error("User not authenticated");
      }

      const request = new GetIncognitoPostRequest();
      request.incognito_post_id = postId;

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/get-incognito-post`,
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
        if (response.status === 404) {
          throw new Error("Post not found");
        }
        throw new Error(`Failed to fetch post: ${response.statusText}`);
      }

      const data: IncognitoPost = await response.json();
      setPost(data);
    } catch (error) {
      setError(
        error instanceof Error
          ? error.message
          : t("incognitoPosts.errors.loadFailed")
      );
    } finally {
      setIsLoadingPost(false);
    }
  };

  const loadComments = async (refresh = false) => {
    setIsLoadingComments(true);
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        throw new Error("User not authenticated");
      }

      const request = new GetIncognitoPostCommentsRequest();
      request.incognito_post_id = postId;
      request.limit = 25;
      request.include_nested_depth = 3;

      if (!refresh && commentsPaginationKey) {
        request.pagination_key = commentsPaginationKey;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/get-incognito-post-comments`,
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

      const data: GetIncognitoPostCommentsResponse = await response.json();

      if (refresh) {
        setComments(data.comments);
      } else {
        setComments((prev) => [...prev, ...data.comments]);
      }

      setCommentsPaginationKey(data.pagination_key);
      setHasMoreComments(!!data.pagination_key);
    } catch (error) {
      setError(
        error instanceof Error
          ? error.message
          : t("incognitoPosts.errors.commentLoadFailed")
      );
    } finally {
      setIsLoadingComments(false);
    }
  };

  const handleAddComment = async () => {
    if (!newComment.trim()) return;

    setIsSubmittingComment(true);
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        throw new Error("User not authenticated");
      }

      const request: AddIncognitoPostCommentRequest = {
        incognito_post_id: postId,
        content: newComment.trim(),
      };

      // Only add in_reply_to if it's a real comment ID (not "new")
      if (replyToComment && replyToComment !== "new") {
        request.in_reply_to = replyToComment;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/add-incognito-post-comment`,
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
        throw new Error(`Failed to add comment: ${response.statusText}`);
      }

      const data: AddIncognitoPostCommentResponse = await response.json();
      setSuccess(t("incognitoPosts.success.commentPosted"));
      setNewComment("");
      setReplyToComment(null);

      // Refresh comments to show the new one
      loadComments(true);
    } catch (error) {
      setError(
        error instanceof Error
          ? error.message
          : t("incognitoPosts.errors.commentFailed")
      );
    } finally {
      setIsSubmittingComment(false);
    }
  };

  const handleError = (errorMessage: string) => {
    setError(errorMessage);
  };

  const handleVoteUpdated = () => {
    // Refresh the post to get updated vote counts
    loadPost();
  };

  const handleCommentVoteUpdated = () => {
    // Refresh comments to get updated vote counts
    loadComments(true);
  };

  const toggleCollapseComment = (commentId: string) => {
    setCollapsedComments((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(commentId)) {
        newSet.delete(commentId);
      } else {
        newSet.add(commentId);
      }
      return newSet;
    });
  };

  const isCommentCollapsed = (commentId: string) => {
    return collapsedComments.has(commentId);
  };

  const buildCommentTree = (comments: IncognitoPostComment[]) => {
    type CommentNode = IncognitoPostComment & { children: CommentNode[] };
    const commentMap = new Map<string, CommentNode>();
    const rootComments: CommentNode[] = [];

    const sortedComments = [...comments].sort(
      (a, b) =>
        new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
    );

    for (const comment of sortedComments) {
      const commentNode: CommentNode = { ...comment, children: [] };
      commentMap.set(comment.comment_id, commentNode);

      if (comment.in_reply_to && commentMap.has(comment.in_reply_to)) {
        const parent = commentMap.get(comment.in_reply_to);
        parent?.children.push(commentNode);
      } else {
        rootComments.push(commentNode);
      }
    }

    return rootComments;
  };

  const commentTree = buildCommentTree(comments);

  const getDepthPattern = (depth: number) => {
    const patterns = [
      { borderStyle: "solid", color: "#2196f3" },
      { borderStyle: "dashed", color: "#f44336" },
      { borderStyle: "dotted", color: "#4caf50" },
      { borderStyle: "solid", color: "#ff9800" },
      { borderStyle: "dashed", color: "#9c27b0" },
    ];
    return patterns[depth % patterns.length];
  };

  const getDepthBackground = (depth: number) => {
    const backgrounds = [
      "rgba(33, 150, 243, 0.08)",
      "rgba(244, 67, 54, 0.08)",
      "rgba(76, 175, 80, 0.08)",
      "rgba(255, 152, 0, 0.08)",
      "rgba(156, 39, 176, 0.08)",
    ];
    return backgrounds[depth % backgrounds.length];
  };

  const CommentNode = ({
    comment,
  }: {
    comment: IncognitoPostComment & { children: any[] };
  }) => {
    const isCollapsed = isCommentCollapsed(comment.comment_id);
    const hasReplies = comment.children.length > 0;
    const showCollapseButton = hasReplies;

    const depthPattern = getDepthPattern(comment.depth);
    const depthBackground = getDepthBackground(comment.depth);

    return (
      <Box sx={{ position: "relative" }}>
        <Box
          sx={{
            backgroundColor: depthBackground,
            borderRadius: 1,
            p: 1.5,
            mt: 2,
            position: "relative",
            zIndex: 1,
          }}
        >
          <Box
            sx={{
              display: "flex",
              gap: 1,
              alignItems: isCollapsed ? "center" : "flex-start",
            }}
          >
            {showCollapseButton && (
              <IconButton
                size="small"
                onClick={() => toggleCollapseComment(comment.comment_id)}
                sx={{
                  p: 0,
                  color: depthPattern.color,
                }}
              >
                {isCollapsed ? (
                  <AddCircleOutlineIcon fontSize="inherit" />
                ) : (
                  <RemoveCircleOutlineIcon fontSize="inherit" />
                )}
              </IconButton>
            )}
            <Box sx={{ flex: 1 }}>
              {isCollapsed ? (
                <Typography variant="body2" color="text.secondary" noWrap>
                  {comment.is_deleted
                    ? "[deleted]"
                    : comment.content.length > 60
                    ? `${comment.content.substring(0, 60)}...`
                    : comment.content}
                </Typography>
              ) : (
                <>
                  {comment.is_deleted ? (
                    <Typography
                      variant="body2"
                      color="text.secondary"
                      fontStyle="italic"
                      sx={{ mb: 1 }}
                    >
                      [deleted]
                    </Typography>
                  ) : (
                    <Typography
                      variant="body1"
                      sx={{
                        mb: 1,
                        whiteSpace: "pre-wrap",
                        fontSize: "0.9rem",
                        lineHeight: 1.5,
                        color: "text.primary",
                      }}
                    >
                      {comment.content}
                    </Typography>
                  )}
                  <Box
                    sx={{
                      display: "flex",
                      alignItems: "center",
                      gap: 1,
                      fontSize: "0.75rem",
                    }}
                  >
                    <Box
                      sx={{ display: "flex", alignItems: "center", gap: 0.25 }}
                    >
                      <IconButton
                        size="small"
                        onClick={async () => {
                          try {
                            const token = Cookies.get("session_token");
                            if (!token)
                              throw new Error("User not authenticated");
                            const endpoint = comment.me_upvoted
                              ? "/hub/unvote-incognito-post-comment"
                              : "/hub/upvote-incognito-post-comment";
                            const res = await fetch(
                              `${config.API_SERVER_PREFIX}${endpoint}`,
                              {
                                method: "POST",
                                headers: {
                                  "Content-Type": "application/json",
                                  Authorization: `Bearer ${token}`,
                                },
                                body: JSON.stringify({
                                  incognito_post_id: postId,
                                  comment_id: comment.comment_id,
                                }),
                              }
                            );
                            if (!res.ok) throw new Error("Vote failed");
                            handleCommentVoteUpdated();
                          } catch (e) {
                            handleError(
                              e instanceof Error ? e.message : "Vote failed"
                            );
                          }
                        }}
                        disabled={!comment.can_upvote}
                        sx={{
                          color: comment.me_upvoted
                            ? "primary.main"
                            : "text.secondary",
                        }}
                      >
                        ▲
                      </IconButton>
                      <Typography variant="caption" sx={{ minWidth: "16px" }}>
                        {comment.upvotes_count}
                      </Typography>
                      <IconButton
                        size="small"
                        onClick={async () => {
                          try {
                            const token = Cookies.get("session_token");
                            if (!token)
                              throw new Error("User not authenticated");
                            const endpoint = comment.me_downvoted
                              ? "/hub/unvote-incognito-post-comment"
                              : "/hub/downvote-incognito-post-comment";
                            const res = await fetch(
                              `${config.API_SERVER_PREFIX}${endpoint}`,
                              {
                                method: "POST",
                                headers: {
                                  "Content-Type": "application/json",
                                  Authorization: `Bearer ${token}`,
                                },
                                body: JSON.stringify({
                                  incognito_post_id: postId,
                                  comment_id: comment.comment_id,
                                }),
                              }
                            );
                            if (!res.ok) throw new Error("Vote failed");
                            handleCommentVoteUpdated();
                          } catch (e) {
                            handleError(
                              e instanceof Error ? e.message : "Vote failed"
                            );
                          }
                        }}
                        disabled={!comment.can_downvote}
                        sx={{
                          color: comment.me_downvoted
                            ? "error.main"
                            : "text.secondary",
                        }}
                      >
                        ▼
                      </IconButton>
                      <Typography variant="caption" sx={{ minWidth: "16px" }}>
                        {comment.downvotes_count}
                      </Typography>
                    </Box>
                    {!comment.is_deleted && comment.depth < 3 && (
                      <Button
                        size="small"
                        sx={{
                          textTransform: "none",
                          color: "text.secondary",
                        }}
                        onClick={() => setReplyToComment(comment.comment_id)}
                      >
                        Reply
                      </Button>
                    )}
                    <Typography
                      variant="caption"
                      color="text.secondary"
                      sx={{ fontSize: "0.8rem", ml: "auto" }}
                    >
                      {new Intl.DateTimeFormat(undefined, {
                        month: "short",
                        day: "numeric",
                        hour: "2-digit",
                        minute: "2-digit",
                      }).format(new Date(comment.created_at))}
                      {comment.is_created_by_me && " • (you)"}
                    </Typography>
                  </Box>
                </>
              )}
            </Box>
          </Box>
        </Box>

        {hasReplies && !isCollapsed && (
          <Box
            sx={{
              pl: 3,
              ml: 3,
              borderLeft: `2px ${depthPattern.borderStyle} ${depthPattern.color}`,
            }}
          >
            {comment.children.map((child: any) => (
              <CommentNode key={child.comment_id} comment={child} />
            ))}
          </Box>
        )}
      </Box>
    );
  };

  return (
    <Container maxWidth="lg" sx={{ py: 3 }}>
      <Box sx={{ mb: 3, display: "flex", alignItems: "center", gap: 2 }}>
        <IconButton
          onClick={() => {
            if (window.history.length > 1) {
              router.back();
            } else {
              router.push("/incognito-posts");
            }
          }}
        >
          <ArrowBackIcon />
        </IconButton>
      </Box>

      {isLoadingPost ? (
        <Box sx={{ display: "flex", justifyContent: "center", my: 4 }}>
          <CircularProgress />
        </Box>
      ) : post ? (
        <Card variant="outlined" sx={{ mb: 3 }}>
          <CardContent>
            <Box
              sx={{
                display: "flex",
                justifyContent: "space-between",
                alignItems: "flex-start",
                mb: 2,
              }}
            >
              <Box>
                <Typography variant="caption" color="text.secondary">
                  {new Intl.DateTimeFormat(undefined, {
                    year: "numeric",
                    month: "short",
                    day: "numeric",
                    hour: "2-digit",
                    minute: "2-digit",
                  }).format(new Date(post.created_at))}
                  {post.is_created_by_me && (
                    <>
                      {" • "}
                      {t("incognitoPosts.post.createdByYou")}
                    </>
                  )}
                </Typography>
              </Box>
            </Box>

            {post.is_deleted ? (
              <Typography
                variant="body1"
                color="text.secondary"
                fontStyle="italic"
                sx={{ mb: 2 }}
              >
                {t("incognitoPosts.post.deleted")}
              </Typography>
            ) : (
              <Typography
                variant="body1"
                sx={{ mb: 2, whiteSpace: "pre-wrap" }}
              >
                {post.content}
              </Typography>
            )}

            <Box sx={{ display: "flex", flexWrap: "wrap", gap: 1, mb: 2 }}>
              {post.tags.map((tag) => (
                <Chip
                  key={tag.id}
                  label={tag.name}
                  size="small"
                  clickable
                  onClick={() => router.push(`/incognito-posts/tag/${tag.id}`)}
                />
              ))}
            </Box>

            <Box sx={{ display: "flex", alignItems: "center", gap: 3 }}>
              <VotingButtons
                postId={postId}
                upvotesCount={post.upvotes_count}
                downvotesCount={post.downvotes_count}
                score={post.score}
                meUpvoted={post.me_upvoted}
                meDownvoted={post.me_downvoted}
                canUpvote={post.can_upvote}
                canDownvote={post.can_downvote}
                onVoteUpdated={handleVoteUpdated}
                onError={handleError}
              />
              <Button
                variant="outlined"
                size="small"
                onClick={() => setReplyToComment("new")}
                sx={{
                  textTransform: "none",
                  fontSize: "0.875rem",
                }}
              >
                Add Comment
              </Button>
            </Box>
          </CardContent>
        </Card>
      ) : (
        <Alert severity="error">Post not found</Alert>
      )}

      <Box sx={{ mt: 3 }}>
        <Typography variant="h6" gutterBottom sx={{ mb: 2 }}>
          Comments
        </Typography>

        {replyToComment && (
          <Box
            sx={{
              mb: 3,
              p: 2,
              backgroundColor: "background.paper",
              borderRadius: 1,
              border: "1px solid",
              borderColor: "divider",
            }}
          >
            <Alert severity="info" sx={{ mb: 2 }}>
              {replyToComment === "new"
                ? "Adding new comment..."
                : "Replying to comment..."}{" "}
              <Button size="small" onClick={() => setReplyToComment(null)}>
                Cancel
              </Button>
            </Alert>
            <TextField
              fullWidth
              multiline
              rows={3}
              placeholder={
                replyToComment === "new"
                  ? t("incognitoPosts.post.commentPlaceholder")
                  : "Write your reply..."
              }
              value={newComment}
              onChange={(e) => setNewComment(e.target.value)}
              disabled={isSubmittingComment}
              sx={{ mb: 2 }}
              autoFocus
            />
            <Box sx={{ display: "flex", gap: 2 }}>
              <Button
                variant="contained"
                onClick={handleAddComment}
                disabled={!newComment.trim() || isSubmittingComment}
              >
                {isSubmittingComment
                  ? "Posting..."
                  : replyToComment === "new"
                  ? "Post Comment"
                  : "Post Reply"}
              </Button>
              <Button onClick={() => setReplyToComment(null)}>Cancel</Button>
            </Box>
          </Box>
        )}

        {isLoadingComments && comments.length === 0 ? (
          <Box sx={{ display: "flex", justifyContent: "center", my: 4 }}>
            <CircularProgress />
          </Box>
        ) : comments.length === 0 ? (
          <Box sx={{ textAlign: "center", py: 4 }}>
            <Typography variant="body2" color="text.secondary">
              No comments yet. Be the first to comment!
            </Typography>
          </Box>
        ) : (
          <Box sx={{ display: "flex", flexDirection: "column" }}>
            {commentTree.map((comment) => (
              <CommentNode key={comment.comment_id} comment={comment} />
            ))}
            {hasMoreComments && (
              <Box sx={{ display: "flex", justifyContent: "center", mt: 3 }}>
                <Button
                  variant="outlined"
                  onClick={() => loadComments(false)}
                  disabled={isLoadingComments}
                >
                  {isLoadingComments ? "Loading..." : "Load More Comments"}
                </Button>
              </Box>
            )}
          </Box>
        )}
      </Box>

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

export default function IncognitoPostDetailsPage() {
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
        <IncognitoPostDetailsContent />
      </Suspense>
    </AuthenticatedLayout>
  );
}
