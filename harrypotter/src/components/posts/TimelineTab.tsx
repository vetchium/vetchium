import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import DeleteIcon from "@mui/icons-material/Delete";
import MoreVertIcon from "@mui/icons-material/MoreVert";
import {
  Avatar,
  Box,
  Button,
  Card,
  CardContent,
  CardHeader,
  Chip,
  CircularProgress,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  IconButton,
  Menu,
  MenuItem,
  Typography,
} from "@mui/material";
import { EmployerPost } from "@vetchium/typespec/common/posts";
import {
  DeleteEmployerPostRequest,
  ListEmployerPostsResponse,
} from "@vetchium/typespec/employer/posts";
import Cookies from "js-cookie";
import { useCallback, useEffect, useState } from "react";

interface TimelineTabProps {
  refreshTrigger: number;
  onError: (error: string) => void;
}

export default function TimelineTab({
  refreshTrigger,
  onError,
}: TimelineTabProps) {
  const { t } = useTranslation();
  const [posts, setPosts] = useState<EmployerPost[]>([]);
  const [loading, setLoading] = useState(false);
  const [paginationKey, setPaginationKey] = useState<string | undefined>(
    undefined
  );
  const [hasMore, setHasMore] = useState(false);
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [selectedPost, setSelectedPost] = useState<string | null>(null);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  // Use useCallback to memoize the fetchPosts function
  const fetchPosts = useCallback(
    async (loadMore = false) => {
      setLoading(true);
      try {
        const token = Cookies.get("session_token");

        // Create request body with just the properties needed by the API
        const requestBody = {
          pagination_key: loadMore ? paginationKey : undefined,
          limit: 10,
        };

        const response = await fetch(
          `${config.API_SERVER_PREFIX}/employer/list-posts`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify(requestBody),
          }
        );

        if (!response.ok) {
          throw new Error(t("posts.fetchError"));
        }

        const data: ListEmployerPostsResponse = await response.json();
        // Make sure posts is an array and has expected structure
        const postsArray = Array.isArray(data.posts) ? data.posts : [];
        setPosts((prev) => (loadMore ? [...prev, ...postsArray] : postsArray));
        setPaginationKey(data.pagination_key || undefined);
        setHasMore(postsArray.length > 0 && !!data.pagination_key);
      } catch (err) {
        console.debug("fetch posts error", err);
        onError(t("posts.fetchError"));
      } finally {
        setLoading(false);
      }
    },
    [paginationKey, t, onError]
  );

  // Fetch posts when the component mounts or refreshTrigger changes
  useEffect(() => {
    fetchPosts();
  }, [refreshTrigger, fetchPosts]);

  const handleLoadMore = () => {
    if (paginationKey) {
      fetchPosts(true);
    }
  };

  const handleMenuOpen = (
    event: React.MouseEvent<HTMLElement>,
    postId: string
  ) => {
    setAnchorEl(event.currentTarget);
    setSelectedPost(postId);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleDeleteClick = () => {
    handleMenuClose();
    setDeleteDialogOpen(true);
  };

  const handleDeleteCancel = () => {
    setDeleteDialogOpen(false);
    setSelectedPost(null);
  };

  const handleDeleteConfirm = async () => {
    if (!selectedPost) return;

    try {
      const token = Cookies.get("session_token");

      const requestBody: DeleteEmployerPostRequest = {
        post_id: selectedPost,
      };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/delete-post`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(requestBody),
        }
      );

      if (!response.ok) {
        throw new Error(t("posts.deleteError"));
      }

      // Remove the deleted post from the state
      setPosts((prev) => prev.filter((post) => post.id !== selectedPost));
    } catch (err) {
      console.debug("delete post error", err);
      onError(t("posts.deleteError"));
    } finally {
      setDeleteDialogOpen(false);
      setSelectedPost(null);
    }
  };

  const formatDate = (dateString: string | Date) => {
    const date =
      typeof dateString === "string" ? new Date(dateString) : dateString;
    return new Intl.DateTimeFormat("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "numeric",
      minute: "numeric",
    }).format(date);
  };

  return (
    <Box>
      {posts.length === 0 && !loading ? (
        <Typography align="center" color="textSecondary" sx={{ py: 4 }}>
          {t("posts.noPostsFound")}
        </Typography>
      ) : (
        posts.map((post) => (
          <Card key={post.id} sx={{ mb: 2 }}>
            <CardHeader
              avatar={
                <Avatar>
                  {post.employer_name
                    ? post.employer_name.charAt(0).toUpperCase()
                    : "V"}
                </Avatar>
              }
              action={
                <IconButton
                  aria-label="post-menu"
                  onClick={(e) => handleMenuOpen(e, post.id)}
                >
                  <MoreVertIcon />
                </IconButton>
              }
              title={post.employer_name || "Unknown Organization"}
              subheader={formatDate(post.updated_at)}
            />
            <CardContent>
              <Typography variant="body1" sx={{ whiteSpace: "pre-wrap" }}>
                {post.content}
              </Typography>
              {post.tags && post.tags.length > 0 && (
                <Box sx={{ mt: 2, display: "flex", flexWrap: "wrap", gap: 1 }}>
                  {post.tags.map((tag) => (
                    <Chip
                      key={tag}
                      label={tag}
                      size="small"
                      variant="outlined"
                    />
                  ))}
                </Box>
              )}
            </CardContent>
          </Card>
        ))
      )}

      {loading && (
        <Box sx={{ display: "flex", justifyContent: "center", py: 3 }}>
          <CircularProgress />
        </Box>
      )}

      {hasMore && !loading && (
        <Box sx={{ display: "flex", justifyContent: "center", mt: 2 }}>
          <Button variant="outlined" onClick={handleLoadMore}>
            {t("posts.loadMore")}
          </Button>
        </Box>
      )}

      {/* Post actions menu */}
      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleMenuClose}
      >
        <MenuItem onClick={handleDeleteClick}>
          <DeleteIcon fontSize="small" sx={{ mr: 1 }} />
          {t("common.delete")}
        </MenuItem>
      </Menu>

      {/* Delete confirmation dialog */}
      <Dialog open={deleteDialogOpen} onClose={handleDeleteCancel}>
        <DialogTitle>{t("common.warning")}</DialogTitle>
        <DialogContent>
          <DialogContentText>{t("posts.deleteConfirm")}</DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleDeleteCancel}>{t("common.cancel")}</Button>
          <Button onClick={handleDeleteConfirm} color="error" autoFocus>
            {t("common.delete")}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
