"use client";

import { config } from "@/config";
import { useMyDetails } from "@/hooks/useMyDetails";
import { useTranslation } from "@/hooks/useTranslation";
import DeleteIcon from "@mui/icons-material/Delete";
import {
  Avatar,
  Box,
  Button,
  CircularProgress,
  Collapse,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  IconButton,
  TextField,
  Typography,
  useTheme,
} from "@mui/material";
import {
  AddPostCommentRequest,
  AddPostCommentResponse,
  DeleteMyCommentRequest,
  DeletePostCommentRequest,
  GetPostCommentsRequest,
  PostComment,
} from "@vetchium/typespec";
import { formatDistanceToNow } from "date-fns";
import Cookies from "js-cookie";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

interface CommentsProps {
  postId: string;
  commentsCount: number;
  canComment: boolean;
  amIAuthor: boolean;
  onCommentsCountChange: (newCount: number) => void;
}

export default function Comments({
  postId,
  commentsCount,
  canComment,
  amIAuthor,
  onCommentsCountChange,
}: CommentsProps) {
  const { t } = useTranslation();
  const theme = useTheme();
  const router = useRouter();
  const { details } = useMyDetails();

  const [isExpanded, setIsExpanded] = useState(false);
  const [comments, setComments] = useState<PostComment[]>([]);
  const [loading, setLoading] = useState(false);
  const [newComment, setNewComment] = useState("");
  const [addingComment, setAddingComment] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [commentToDelete, setCommentToDelete] = useState<PostComment | null>(
    null
  );
  const [deletingComment, setDeletingComment] = useState(false);
  const [showCommentInput, setShowCommentInput] = useState(false);

  const maxCommentLength = 4096;

  // Load comments when expanded
  useEffect(() => {
    if (isExpanded && comments.length === 0 && commentsCount > 0) {
      loadComments();
    }
  }, [isExpanded, commentsCount]);

  // Handle toggle comments event from PostCard
  useEffect(() => {
    const handleToggleComments = () => {
      if (commentsCount > 0) {
        setIsExpanded(!isExpanded);
      }
      if (canComment) {
        setShowCommentInput(!showCommentInput);
      }
    };

    const handleRefreshComments = () => {
      // Clear comments and reset state when comments are deleted
      setComments([]);
      setIsExpanded(false);
      setShowCommentInput(false);
      onCommentsCountChange(0);
    };

    const element = document.getElementById(`comments-${postId}`);
    if (element) {
      element.addEventListener("toggleComments", handleToggleComments);
      element.addEventListener("refreshComments", handleRefreshComments);
      return () => {
        element.removeEventListener("toggleComments", handleToggleComments);
        element.removeEventListener("refreshComments", handleRefreshComments);
      };
    }
  }, [
    postId,
    isExpanded,
    showCommentInput,
    commentsCount,
    canComment,
    onCommentsCountChange,
  ]);

  const loadComments = async () => {
    setLoading(true);
    const token = Cookies.get("session_token");
    if (!token) {
      router.push("/login");
      return;
    }

    try {
      const request = new GetPostCommentsRequest();
      request.post_id = postId;
      request.limit = 5; // Show top 5 comments

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/get-post-comments`,
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
        if (response.status === 401) {
          Cookies.remove("session_token", { path: "/" });
          router.push("/login");
          return;
        }
        throw new Error(`Failed to load comments: ${response.statusText}`);
      }

      const data: PostComment[] = await response.json();
      setComments(data);
    } catch (error) {
      console.error("Error loading comments:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleAddComment = async () => {
    if (!newComment.trim()) {
      return;
    }

    if (newComment.length > maxCommentLength) {
      return;
    }

    setAddingComment(true);
    const token = Cookies.get("session_token");
    if (!token) {
      router.push("/login");
      return;
    }

    try {
      const request = new AddPostCommentRequest();
      request.post_id = postId;
      request.content = newComment.trim();

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/add-post-comment`,
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
        if (response.status === 401) {
          Cookies.remove("session_token", { path: "/" });
          router.push("/login");
          return;
        }
        if (response.status === 403) {
          // Comments disabled
          return;
        }
        throw new Error(`Failed to add comment: ${response.statusText}`);
      }

      const data: AddPostCommentResponse = await response.json();

      // Create new comment object for immediate display
      const newCommentObj: PostComment = {
        id: data.comment_id,
        content: newComment.trim(),
        author_name: details?.full_name || "",
        author_handle: details?.handle || "",
        created_at: new Date(),
      };

      // Add to comments list and update count
      setComments([newCommentObj, ...comments]);
      setNewComment("");
      onCommentsCountChange(commentsCount + 1);

      // Expand comments if not already expanded
      if (!isExpanded) {
        setIsExpanded(true);
      }
    } catch (error) {
      console.error("Error adding comment:", error);
    } finally {
      setAddingComment(false);
    }
  };

  const handleDeleteComment = async (comment: PostComment) => {
    setDeletingComment(true);
    const token = Cookies.get("session_token");
    if (!token) {
      router.push("/login");
      return;
    }

    try {
      const isMyComment = comment.author_handle === details?.handle;
      const endpoint = isMyComment
        ? `${config.API_SERVER_PREFIX}/hub/delete-my-comment`
        : `${config.API_SERVER_PREFIX}/hub/delete-post-comment`;

      const request = isMyComment
        ? ({
            post_id: postId,
            comment_id: comment.id,
          } as DeleteMyCommentRequest)
        : ({
            post_id: postId,
            comment_id: comment.id,
          } as DeletePostCommentRequest);

      const response = await fetch(endpoint, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(request),
      });

      if (!response.ok) {
        if (response.status === 401) {
          Cookies.remove("session_token", { path: "/" });
          router.push("/login");
          return;
        }
        throw new Error(`Failed to delete comment: ${response.statusText}`);
      }

      // Remove comment from list and update count
      setComments(comments.filter((c) => c.id !== comment.id));
      onCommentsCountChange(Math.max(0, commentsCount - 1));
    } catch (error) {
      console.error("Error deleting comment:", error);
    } finally {
      setDeletingComment(false);
      setDeleteDialogOpen(false);
      setCommentToDelete(null);
    }
  };

  const openDeleteDialog = (comment: PostComment) => {
    setCommentToDelete(comment);
    setDeleteDialogOpen(true);
  };

  const canDeleteComment = (comment: PostComment) => {
    return amIAuthor || comment.author_handle === details?.handle;
  };

  // Don't render anything if no comments and can't comment
  if (commentsCount === 0 && !canComment) {
    return null;
  }

  return (
    <Box id={`comments-${postId}`} sx={{ mt: 0.5 }}>
      {/* Compact comment interface */}
      {commentsCount > 0 && (
        <Collapse in={isExpanded}>
          <Box sx={{ mt: 1, mb: 1 }}>
            {loading ? (
              <Box sx={{ display: "flex", justifyContent: "center", p: 1 }}>
                <CircularProgress size={20} />
              </Box>
            ) : comments.length > 0 ? (
              <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
                {comments.slice(0, 3).map((comment) => (
                  <Box
                    key={comment.id}
                    sx={{
                      display: "flex",
                      gap: 1,
                      p: 1,
                      backgroundColor: theme.palette.background.default,
                      borderRadius: 1,
                      border: `1px solid ${theme.palette.divider}`,
                    }}
                  >
                    <Link
                      href={`/u/${comment.author_handle}`}
                      target="_blank"
                      rel="noopener noreferrer"
                      style={{ textDecoration: "none" }}
                    >
                      <Avatar sx={{ width: 24, height: 24 }}>
                        {comment.author_name?.charAt(0) ||
                          comment.author_handle.charAt(0)}
                      </Avatar>
                    </Link>
                    <Box sx={{ flex: 1, minWidth: 0 }}>
                      <Box
                        sx={{
                          display: "flex",
                          alignItems: "center",
                          gap: 0.5,
                          mb: 0.25,
                        }}
                      >
                        <Link
                          href={`/u/${comment.author_handle}`}
                          target="_blank"
                          rel="noopener noreferrer"
                          style={{ textDecoration: "none", color: "inherit" }}
                        >
                          <Typography
                            variant="caption"
                            sx={{
                              fontWeight: 500,
                              color: theme.palette.text.primary,
                            }}
                          >
                            {comment.author_name || comment.author_handle}
                          </Typography>
                        </Link>
                        <Typography
                          variant="caption"
                          sx={{ color: theme.palette.text.secondary }}
                        >
                          @{comment.author_handle}
                        </Typography>
                        <Typography
                          variant="caption"
                          sx={{ color: theme.palette.text.secondary }}
                        >
                          ·
                        </Typography>
                        <Typography
                          variant="caption"
                          sx={{ color: theme.palette.text.secondary }}
                        >
                          {formatDistanceToNow(new Date(comment.created_at), {
                            addSuffix: true,
                          })}
                        </Typography>
                        {canDeleteComment(comment) && (
                          <IconButton
                            size="small"
                            onClick={() => openDeleteDialog(comment)}
                            sx={{
                              ml: "auto",
                              color: theme.palette.text.secondary,
                              "&:hover": {
                                color: theme.palette.error.main,
                              },
                              p: 0.25,
                            }}
                          >
                            <DeleteIcon sx={{ fontSize: "0.8rem" }} />
                          </IconButton>
                        )}
                      </Box>
                      <Typography
                        variant="body2"
                        sx={{
                          color: theme.palette.text.primary,
                          whiteSpace: "pre-wrap",
                          wordBreak: "break-word",
                          fontSize: "0.8rem",
                          lineHeight: 1.3,
                        }}
                      >
                        {comment.content}
                      </Typography>
                    </Box>
                  </Box>
                ))}
                {comments.length > 3 && (
                  <Typography
                    variant="caption"
                    sx={{
                      color: theme.palette.text.secondary,
                      textAlign: "center",
                      cursor: "pointer",
                      "&:hover": {
                        color: theme.palette.primary.main,
                      },
                    }}
                    onClick={() => {
                      // Navigate to full post page to see all comments
                      router.push(`/posts/${postId}`);
                    }}
                  >
                    View {comments.length - 3} more comments
                  </Typography>
                )}
              </Box>
            ) : (
              <Typography
                variant="caption"
                sx={{
                  color: theme.palette.text.secondary,
                  textAlign: "center",
                  p: 1,
                }}
              >
                {t("comments.noComments")}
              </Typography>
            )}
          </Box>
        </Collapse>
      )}

      {/* Compact comment input */}
      {canComment && (
        <Collapse in={showCommentInput}>
          <Box sx={{ mt: 1, mb: 1 }}>
            <Box sx={{ display: "flex", gap: 1, alignItems: "flex-start" }}>
              <Avatar sx={{ width: 24, height: 24 }}>
                {details?.full_name?.charAt(0) || details?.handle?.charAt(0)}
              </Avatar>
              <Box sx={{ flex: 1 }}>
                <TextField
                  fullWidth
                  multiline
                  rows={2}
                  value={newComment}
                  onChange={(e) => setNewComment(e.target.value)}
                  placeholder={t("comments.addPlaceholder")}
                  variant="outlined"
                  size="small"
                  sx={{
                    "& .MuiOutlinedInput-root": {
                      fontSize: "0.8rem",
                    },
                  }}
                />
                <Box
                  sx={{
                    display: "flex",
                    justifyContent: "space-between",
                    alignItems: "center",
                    mt: 0.5,
                  }}
                >
                  <Typography
                    variant="caption"
                    sx={{
                      color:
                        newComment.length > maxCommentLength
                          ? theme.palette.error.main
                          : theme.palette.text.secondary,
                    }}
                  >
                    {newComment.length}/{maxCommentLength}
                  </Typography>
                  <Box sx={{ display: "flex", gap: 1 }}>
                    <Button
                      size="small"
                      onClick={() => {
                        setShowCommentInput(false);
                        setNewComment("");
                      }}
                      sx={{ fontSize: "0.7rem" }}
                    >
                      Cancel
                    </Button>
                    <Button
                      variant="contained"
                      size="small"
                      onClick={handleAddComment}
                      disabled={
                        addingComment ||
                        !newComment.trim() ||
                        newComment.length > maxCommentLength
                      }
                      sx={{ fontSize: "0.7rem" }}
                    >
                      {addingComment ? (
                        <CircularProgress size={12} color="inherit" />
                      ) : (
                        t("comments.add")
                      )}
                    </Button>
                  </Box>
                </Box>
              </Box>
            </Box>
          </Box>
        </Collapse>
      )}

      {/* Delete confirmation dialog */}
      <Dialog
        open={deleteDialogOpen}
        onClose={() => setDeleteDialogOpen(false)}
      >
        <DialogTitle>
          {commentToDelete?.author_handle === details?.handle
            ? t("comments.deleteMyComment")
            : t("comments.deleteComment")}
        </DialogTitle>
        <DialogContent>
          <DialogContentText>
            {commentToDelete?.author_handle === details?.handle
              ? t("comments.deleteMyCommentConfirm")
              : t("comments.deleteCommentConfirm")}
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)}>
            {t("common.cancel")}
          </Button>
          <Button
            onClick={() => {
              if (commentToDelete) {
                handleDeleteComment(commentToDelete);
              }
            }}
            color="error"
            disabled={deletingComment}
          >
            {deletingComment ? (
              <CircularProgress size={16} color="inherit" />
            ) : (
              t("common.delete")
            )}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
