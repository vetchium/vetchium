"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import { config } from "@/config";
import { useAuth } from "@/hooks/useAuth";
import { useTranslation } from "@/hooks/useTranslation";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import CloseIcon from "@mui/icons-material/Close";
import {
  Alert,
  Box,
  Button,
  Card,
  CardContent,
  Chip,
  CircularProgress,
  Container,
  Divider,
  IconButton,
  Paper,
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
import CommentVotingButtons from "../components/CommentVotingButtons";
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

      if (replyToComment) {
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

  const handlePostVote = async (action: "upvote" | "downvote" | "unvote") => {
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        throw new Error("User not authenticated");
      }

      const request = {
        incognito_post_id: postId,
      };

      const endpoint =
        action === "upvote"
          ? "/hub/upvote-incognito-post"
          : action === "downvote"
          ? "/hub/downvote-incognito-post"
          : "/hub/unvote-incognito-post";

      const response = await fetch(`${config.API_SERVER_PREFIX}${endpoint}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(request),
      });

      if (!response.ok) {
        throw new Error(`Failed to ${action} post: ${response.statusText}`);
      }

      handleVoteUpdated();
    } catch (error) {
      handleError(
        error instanceof Error
          ? error.message
          : t("incognitoPosts.errors.voteFailed")
      );
    }
  };

  const renderComment = (comment: IncognitoPostComment) => (
    <Card
      key={comment.comment_id}
      variant="outlined"
      sx={{ ml: comment.depth * 2 }}
    >
      <CardContent>
        <Box
          sx={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "flex-start",
            mb: 1,
          }}
        >
          <Typography variant="caption" color="text.secondary">
            {new Date(comment.created_at).toLocaleString()}
            {comment.depth > 0 && (
              <>
                {" • "}
                {t("incognitoPosts.post.inReplyTo")}
              </>
            )}
            {comment.is_created_by_me && (
              <>
                {" • "}
                {t("incognitoPosts.post.createdByYou")}
              </>
            )}
          </Typography>
        </Box>

        {comment.is_deleted ? (
          <Typography variant="body2" color="text.secondary" fontStyle="italic">
            {t("incognitoPosts.comments.deleted")}
          </Typography>
        ) : (
          <Typography variant="body1" sx={{ mb: 2, whiteSpace: "pre-wrap" }}>
            {comment.content}
          </Typography>
        )}

        <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
          <CommentVotingButtons
            postId={postId}
            commentId={comment.comment_id}
            upvotesCount={comment.upvotes_count}
            downvotesCount={comment.downvotes_count}
            score={comment.score}
            meUpvoted={comment.me_upvoted}
            meDownvoted={comment.me_downvoted}
            canUpvote={comment.can_upvote}
            canDownvote={comment.can_downvote}
            onVoteUpdated={handleCommentVoteUpdated}
            onError={handleError}
          />

          {!comment.is_deleted && comment.depth < 3 && (
            <Button
              size="small"
              onClick={() => setReplyToComment(comment.comment_id)}
            >
              {t("incognitoPosts.post.reply")}
            </Button>
          )}
        </Box>
      </CardContent>
    </Card>
  );

  return (
    <Container maxWidth="lg" sx={{ py: 3 }}>
      {/* Header */}
      <Box sx={{ mb: 3, display: "flex", alignItems: "center", gap: 2 }}>
        <IconButton onClick={() => router.back()}>
          <ArrowBackIcon />
        </IconButton>
        <Typography variant="h5">
          {t("incognitoPosts.post.viewDetails")}
        </Typography>
      </Box>

      {/* Post Content */}
      {isLoadingPost ? (
        <Box sx={{ display: "flex", justifyContent: "center", my: 4 }}>
          <CircularProgress />
        </Box>
      ) : post ? (
        <Paper sx={{ mb: 3 }}>
          <CardContent>
            {/* Post Header */}
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
                  {new Date(post.created_at).toLocaleString()}
                  {post.is_created_by_me && (
                    <>
                      {" • "}
                      {t("incognitoPosts.post.createdByYou")}
                    </>
                  )}
                </Typography>
              </Box>
            </Box>

            {/* Post Content */}
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

            {/* Tags */}
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

            {/* Voting and Stats */}
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
            </Box>
          </CardContent>
        </Paper>
      ) : (
        <Alert severity="error">Post not found</Alert>
      )}

      {/* Comments Section */}
      <Paper sx={{ p: 3 }}>
        <Typography variant="h6" gutterBottom>
          Comments
        </Typography>

        {/* Add Comment */}
        <Box sx={{ mb: 3 }}>
          {replyToComment && (
            <Alert severity="info" sx={{ mb: 2 }}>
              Replying to comment...{" "}
              <Button size="small" onClick={() => setReplyToComment(null)}>
                Cancel
              </Button>
            </Alert>
          )}

          <TextField
            fullWidth
            multiline
            rows={3}
            placeholder={t("incognitoPosts.post.commentPlaceholder")}
            value={newComment}
            onChange={(e) => setNewComment(e.target.value)}
            disabled={isSubmittingComment}
            sx={{ mb: 2 }}
          />

          <Box sx={{ display: "flex", gap: 2 }}>
            <Button
              variant="contained"
              onClick={handleAddComment}
              disabled={!newComment.trim() || isSubmittingComment}
            >
              {isSubmittingComment
                ? "Posting..."
                : t("incognitoPosts.post.postComment")}
            </Button>

            {replyToComment && (
              <Button onClick={() => setReplyToComment(null)}>
                {t("incognitoPosts.post.cancelComment")}
              </Button>
            )}
          </Box>
        </Box>

        <Divider sx={{ mb: 2 }} />

        {/* Comments List */}
        {isLoadingComments && comments.length === 0 ? (
          <Box sx={{ display: "flex", justifyContent: "center", my: 4 }}>
            <CircularProgress />
          </Box>
        ) : comments.length === 0 ? (
          <Typography variant="body2" color="text.secondary" align="center">
            No comments yet. Be the first to comment!
          </Typography>
        ) : (
          <Box sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
            {comments.map(renderComment)}

            {hasMoreComments && (
              <Box sx={{ display: "flex", justifyContent: "center", mt: 2 }}>
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
