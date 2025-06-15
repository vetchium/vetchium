"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import { config } from "@/config";
import { useAuth } from "@/hooks/useAuth";
import { useTranslation } from "@/hooks/useTranslation";
import {
  ThumbDown,
  ThumbDownOutlined,
  ThumbUp,
  ThumbUpOutlined,
} from "@mui/icons-material";
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
  AddIncognitoPostCommentRequest,
  DeleteIncognitoPostCommentRequest,
  DownvoteIncognitoPostCommentRequest,
  GetIncognitoPostCommentsRequest,
  GetIncognitoPostCommentsResponse,
  GetIncognitoPostRequest,
  IncognitoPost,
  IncognitoPostComment,
  IncognitoPostCommentSortBy,
  UnvoteIncognitoPostCommentRequest,
  UpvoteIncognitoPostCommentRequest,
} from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useParams, useRouter } from "next/navigation";
import { Suspense, useEffect, useState } from "react";
import VotingButtons from "../components/VotingButtons";

function IncognitoPostDetailsContent() {
  const { t } = useTranslation();
  const router = useRouter();
  const params = useParams();
  useAuth();

  const postId = params?.postId as string;

  const [post, setPost] = useState<IncognitoPost | null>(null);
  const [comments, setComments] = useState<IncognitoPostComment[]>([]);
  const [sortedCommentTree, setSortedCommentTree] = useState<any[]>([]);
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
  const [totalCommentsCount, setTotalCommentsCount] = useState(0);
  const [sortBy, setSortBy] = useState<IncognitoPostCommentSortBy>(
    IncognitoPostCommentSortBy.Top
  );
  const [collapsedComments, setCollapsedComments] = useState<Set<string>>(
    new Set()
  );
  const [loadingReplies, setLoadingReplies] = useState<Set<string>>(new Set());

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
  }, [postId, sortBy]);

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
      request.direct_replies_per_comment = 3;
      request.sort_by = sortBy;

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

      let newComments: IncognitoPostComment[];
      if (refresh) {
        newComments = data.comments;
        setComments(data.comments);
      } else {
        newComments = [...comments, ...data.comments];
        setComments((prev) => [...prev, ...data.comments]);
      }

      // The API now loads exactly direct_replies_per_comment direct replies per top-level comment
      // No need for additional preview loading

      // Build and sort the comment tree only when loading from server
      const newSortedTree = buildCommentTree(newComments);
      setSortedCommentTree(newSortedTree);

      setCommentsPaginationKey(data.pagination_key);
      setHasMoreComments(!!data.pagination_key);
      setTotalCommentsCount(data.total_comments_count);
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

      const request = new AddIncognitoPostCommentRequest();
      request.incognito_post_id = postId;
      request.content = newComment.trim();

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

      await response.json();
      setSuccess(t("incognitoPosts.success.commentPosted"));
      setNewComment("");
      setReplyToComment(null);
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
    loadPost();
  };

  const handleDeleteComment = async (commentId: string) => {
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        throw new Error("User not authenticated");
      }

      const request = new DeleteIncognitoPostCommentRequest();
      request.incognito_post_id = postId;
      request.comment_id = commentId;

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/delete-incognito-post-comment`,
        {
          method: "DELETE",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(request),
        }
      );

      if (!response.ok) {
        throw new Error(`Failed to delete comment: ${response.statusText}`);
      }
      loadComments(true);
    } catch (error) {
      handleError(
        error instanceof Error ? error.message : "Failed to delete comment"
      );
    }
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

  const loadMoreReplies = async (
    commentId: string,
    incognitoPostId: string,
    loadCount?: number
  ) => {
    setLoadingReplies((prev) => new Set(prev).add(commentId));
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        throw new Error("User not authenticated");
      }

      const request = {
        incognito_post_id: incognitoPostId,
        parent_comment_id: commentId,
        limit: loadCount || 50,
        direct_only: true,
        max_depth: 2,
      };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/get-comment-replies`,
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
        throw new Error(`Failed to load replies: ${response.statusText}`);
      }

      const data = await response.json();

      // Update the sorted tree with the new replies
      setSortedCommentTree((prevTree) => {
        return updateCommentWithReplies(prevTree, commentId, data.replies);
      });

      // Also update the comments array
      setComments((prevComments) => {
        return [...prevComments, ...data.replies];
      });

      setSuccess("Loaded more replies");
    } catch (error) {
      setError(
        error instanceof Error ? error.message : "Failed to load more replies"
      );
    } finally {
      setLoadingReplies((prev) => {
        const newSet = new Set(prev);
        newSet.delete(commentId);
        return newSet;
      });
    }
  };

  const buildCommentTree = (comments: IncognitoPostComment[]) => {
    type CommentNode = IncognitoPostComment & { children: CommentNode[] };

    const commentMap = new Map<string, CommentNode>();
    const rootComments: CommentNode[] = [];

    // First pass: create all comment nodes
    comments.forEach((comment) => {
      commentMap.set(comment.comment_id, {
        ...comment,
        children: [],
      });
    });

    // Second pass: build the tree structure
    comments.forEach((comment) => {
      const commentNode = commentMap.get(comment.comment_id)!;

      if (comment.in_reply_to) {
        const parentNode = commentMap.get(comment.in_reply_to);
        if (parentNode) {
          parentNode.children.push(commentNode);
        } else {
          // Parent not found in current page, treat as root
          rootComments.push(commentNode);
        }
      } else {
        rootComments.push(commentNode);
      }
    });

    // Sort root comments and children according to sort preference
    const sortComments = (comments: CommentNode[]) => {
      return comments.sort((a, b) => {
        switch (sortBy) {
          case IncognitoPostCommentSortBy.Top:
            if (a.score !== b.score) return b.score - a.score;
            return (
              new Date(b.created_at).getTime() -
              new Date(a.created_at).getTime()
            );
          case IncognitoPostCommentSortBy.New:
            return (
              new Date(b.created_at).getTime() -
              new Date(a.created_at).getTime()
            );
          case IncognitoPostCommentSortBy.Old:
            return (
              new Date(a.created_at).getTime() -
              new Date(b.created_at).getTime()
            );
          default:
            return b.score - a.score;
        }
      });
    };

    // Sort each level
    const sortedRootComments = sortComments(rootComments);
    sortedRootComments.forEach((comment) => {
      comment.children = sortComments(comment.children);
    });

    return sortedRootComments;
  };

  // Function to update comment data in the sorted tree without re-sorting
  const updateCommentInTree = (
    tree: any[],
    targetCommentId: string,
    updatedComment: IncognitoPostComment
  ): any[] => {
    return tree.map((comment) => {
      if (comment.comment_id === targetCommentId) {
        return { ...updatedComment, children: comment.children };
      }
      if (comment.children && comment.children.length > 0) {
        return {
          ...comment,
          children: updateCommentInTree(
            comment.children,
            targetCommentId,
            updatedComment
          ),
        };
      }
      return comment;
    });
  };

  // Function to update comment with new direct replies (predictable loading)
  const updateCommentWithReplies = (
    tree: any[],
    targetCommentId: string,
    newReplies: IncognitoPostComment[]
  ): any[] => {
    return tree.map((comment) => {
      if (comment.comment_id === targetCommentId) {
        // For predictable loading, we only add direct replies (immediate children)
        // Filter to only direct children of the target comment
        const directReplies = newReplies.filter(
          (reply) => reply.in_reply_to === targetCommentId
        );

        // Combine existing children with new direct replies
        const existingChildren = comment.children || [];
        const existingChildIds = new Set(
          existingChildren.map((c: any) => c.comment_id)
        );

        // Only add new replies that aren't already present
        const newDirectReplies = directReplies.filter(
          (reply) => !existingChildIds.has(reply.comment_id)
        );

        // Convert new replies to the tree format
        const newChildren = newDirectReplies.map((reply) => ({
          ...reply,
          children: [], // Direct replies start with no children
        }));

        return {
          ...comment,
          children: [...existingChildren, ...newChildren],
        };
      }
      if (comment.children && comment.children.length > 0) {
        return {
          ...comment,
          children: updateCommentWithReplies(
            comment.children,
            targetCommentId,
            newReplies
          ),
        };
      }
      return comment;
    });
  };

  // Helper function to flatten a comment tree back to a flat array
  const flattenCommentTree = (tree: any[]): IncognitoPostComment[] => {
    const result: IncognitoPostComment[] = [];

    const flatten = (comments: any[]) => {
      comments.forEach((comment) => {
        // Add the comment itself (without children to avoid circular references)
        const { children, ...commentData } = comment;
        result.push(commentData);

        // Recursively flatten children
        if (children && children.length > 0) {
          flatten(children);
        }
      });
    };

    flatten(tree);
    return result;
  };

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
    onDeleteComment,
    renderDepth = 0,
  }: {
    comment: IncognitoPostComment & { children: any[] };
    onDeleteComment: (commentId: string) => void;
    renderDepth?: number;
  }) => {
    const [isVoting, setIsVoting] = useState(false);
    const isCollapsed = isCommentCollapsed(comment.comment_id);
    const hasReplies = comment.children.length > 0;
    const hasMoreReplies = comment.replies_count > comment.children.length;
    const showCollapseButton = hasReplies || hasMoreReplies; // Show if has children OR more to load

    // Use renderDepth for visual styling instead of database depth
    const depthPattern = getDepthPattern(renderDepth);
    const depthBackground = getDepthBackground(renderDepth);

    const handleVote = async (action: "upvote" | "downvote" | "unvote") => {
      if (isVoting) return;
      setIsVoting(true);
      try {
        const token = Cookies.get("session_token");
        if (!token) throw new Error("User not authenticated");

        let endpoint = "";
        let request: any;

        // If user clicks on their existing vote, unvote instead
        if (action === "upvote" && comment.me_upvoted) {
          action = "unvote";
        } else if (action === "downvote" && comment.me_downvoted) {
          action = "unvote";
        }

        switch (action) {
          case "upvote":
            endpoint = "/hub/upvote-incognito-post-comment";
            request = new UpvoteIncognitoPostCommentRequest();
            break;
          case "downvote":
            endpoint = "/hub/downvote-incognito-post-comment";
            request = new DownvoteIncognitoPostCommentRequest();
            break;
          case "unvote":
            endpoint = "/hub/unvote-incognito-post-comment";
            request = new UnvoteIncognitoPostCommentRequest();
            break;
        }
        request.incognito_post_id = postId;
        request.comment_id = comment.comment_id;

        const res = await fetch(`${config.API_SERVER_PREFIX}${endpoint}`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(request),
        });

        if (!res.ok) {
          throw new Error("Vote failed");
        }

        // Update the comment's vote counts and states locally
        const updatedComment = { ...comment };
        if (action === "unvote") {
          if (comment.me_upvoted) {
            updatedComment.upvotes_count--;
            updatedComment.me_upvoted = false;
          }
          if (comment.me_downvoted) {
            updatedComment.downvotes_count--;
            updatedComment.me_downvoted = false;
          }
          updatedComment.score =
            updatedComment.upvotes_count - updatedComment.downvotes_count;
          updatedComment.can_upvote = true;
          updatedComment.can_downvote = true;
        } else if (action === "upvote") {
          if (comment.me_downvoted) {
            updatedComment.downvotes_count--;
            updatedComment.me_downvoted = false;
          }
          updatedComment.upvotes_count++;
          updatedComment.me_upvoted = true;
          updatedComment.score =
            updatedComment.upvotes_count - updatedComment.downvotes_count;
          updatedComment.can_upvote = false;
          updatedComment.can_downvote = false;
        } else if (action === "downvote") {
          if (comment.me_upvoted) {
            updatedComment.upvotes_count--;
            updatedComment.me_upvoted = false;
          }
          updatedComment.downvotes_count++;
          updatedComment.me_downvoted = true;
          updatedComment.score =
            updatedComment.upvotes_count - updatedComment.downvotes_count;
          updatedComment.can_upvote = false;
          updatedComment.can_downvote = false;
        }

        // Update the comment in the comments array AND the sorted tree
        setComments((prevComments) => {
          return prevComments.map((c) => {
            if (c.comment_id === comment.comment_id) {
              return updatedComment;
            }
            return c;
          });
        });

        // Update the sorted tree without re-sorting to preserve positions
        setSortedCommentTree((prevTree) => {
          return updateCommentInTree(
            prevTree,
            comment.comment_id,
            updatedComment
          );
        });
      } catch (e) {
        handleError(e instanceof Error ? e.message : "Vote failed");
      } finally {
        setIsVoting(false);
      }
    };

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
                        onClick={() =>
                          handleVote(comment.me_upvoted ? "unvote" : "upvote")
                        }
                        disabled={
                          (!comment.can_upvote && !comment.me_upvoted) ||
                          isVoting
                        }
                        sx={{
                          color: comment.me_upvoted
                            ? "primary.main"
                            : "text.secondary",
                        }}
                      >
                        {comment.me_upvoted ? (
                          <ThumbUp fontSize="inherit" />
                        ) : (
                          <ThumbUpOutlined fontSize="inherit" />
                        )}
                      </IconButton>
                      <Typography variant="caption" sx={{ minWidth: "16px" }}>
                        {comment.upvotes_count}
                      </Typography>
                      <IconButton
                        size="small"
                        onClick={() =>
                          handleVote(
                            comment.me_downvoted ? "unvote" : "downvote"
                          )
                        }
                        disabled={
                          (!comment.can_downvote && !comment.me_downvoted) ||
                          isVoting
                        }
                        sx={{
                          color: comment.me_downvoted
                            ? "error.main"
                            : "text.secondary",
                        }}
                      >
                        {comment.me_downvoted ? (
                          <ThumbDown fontSize="inherit" />
                        ) : (
                          <ThumbDownOutlined fontSize="inherit" />
                        )}
                      </IconButton>
                      <Typography variant="caption" sx={{ minWidth: "16px" }}>
                        {comment.downvotes_count}
                      </Typography>
                    </Box>
                    {!comment.is_deleted &&
                      comment.depth < (config.MAX_COMMENT_DEPTH || 4) && (
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
                    {comment.is_created_by_me && !comment.is_deleted && (
                      <Button
                        size="small"
                        sx={{
                          textTransform: "none",
                          color: "text.secondary",
                        }}
                        onClick={() => onDeleteComment(comment.comment_id)}
                      >
                        Delete
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

        {(hasReplies || hasMoreReplies) && !isCollapsed && (
          <Box
            sx={{
              pl: 3,
              ml: 3,
              borderLeft: `2px ${depthPattern.borderStyle} ${depthPattern.color}`,
            }}
          >
            {renderDepth < 4 ? (
              // Normal nested rendering for shallow depths
              <>
                {comment.children.map((child: any) => (
                  <CommentNode
                    key={child.comment_id}
                    comment={child}
                    onDeleteComment={onDeleteComment}
                    renderDepth={renderDepth + 1}
                  />
                ))}
                {hasMoreReplies && (
                  <Box sx={{ mt: 1, pl: 1 }}>
                    <Button
                      size="small"
                      variant="text"
                      sx={{
                        textTransform: "none",
                        color: "primary.main",
                        fontSize: "0.8rem",
                      }}
                      disabled={loadingReplies.has(comment.comment_id)}
                      onClick={() => {
                        const remainingReplies =
                          comment.replies_count - comment.children.length;
                        const loadCount = Math.min(remainingReplies, 10); // Load max 10 at a time
                        loadMoreReplies(comment.comment_id, postId, loadCount);
                      }}
                    >
                      {loadingReplies.has(comment.comment_id)
                        ? "Loading..."
                        : `Load ${Math.min(
                            comment.replies_count - comment.children.length,
                            10
                          )} more replies`}
                    </Button>
                  </Box>
                )}
              </>
            ) : (
              // Show "continue this thread" for deep nesting
              <Box sx={{ mt: 1, pl: 1 }}>
                <Button
                  size="small"
                  variant="text"
                  sx={{
                    textTransform: "none",
                    color: "primary.main",
                    fontSize: "0.8rem",
                    fontStyle: "italic",
                  }}
                  onClick={() => {
                    // For now, load replies normally but could open a focused view
                    const remainingReplies =
                      comment.replies_count - comment.children.length;
                    const loadCount = Math.min(remainingReplies, 10);
                    loadMoreReplies(comment.comment_id, postId, loadCount);
                  }}
                >
                  Continue this thread →
                  {comment.replies_count > 0 &&
                    ` (${comment.replies_count} replies)`}
                </Button>
              </Box>
            )}
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
        <Typography variant="h5">Incognito Post</Typography>
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
        <Box sx={{ display: "flex", alignItems: "center", gap: 2, mb: 2 }}>
          <Typography variant="h6" gutterBottom sx={{ mb: 0 }}>
            Comments ({totalCommentsCount})
          </Typography>
          <Box sx={{ display: "flex", gap: 1 }}>
            {Object.values(IncognitoPostCommentSortBy).map((sort) => (
              <Button
                key={sort}
                size="small"
                variant={sortBy === sort ? "contained" : "outlined"}
                onClick={() => setSortBy(sort)}
                sx={{ textTransform: "capitalize" }}
              >
                {sort}
              </Button>
            ))}
          </Box>
        </Box>

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
            {sortedCommentTree.map((comment) => (
              <CommentNode
                key={comment.comment_id}
                comment={comment}
                onDeleteComment={handleDeleteComment}
                renderDepth={0}
              />
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
