"use client";

import { useTranslation } from "@/hooks/useTranslation";
import { Comment, Delete } from "@mui/icons-material";
import {
  Box,
  Button,
  Card,
  CardContent,
  CardHeader,
  Chip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  Typography,
  useTheme,
} from "@mui/material";
import { IncognitoPostSummary } from "@vetchium/typespec";
import { formatDistanceToNow } from "date-fns";
import { useRouter } from "next/navigation";
import { useState } from "react";
import VotingButtons from "./VotingButtons";

interface IncognitoPostCardProps {
  post: IncognitoPostSummary;
  onDeleted: () => void;
  onVoteUpdated: () => void;
  onError: (error: string) => void;
}

export default function IncognitoPostCard({
  post,
  onDeleted,
  onVoteUpdated,
  onError,
}: IncognitoPostCardProps) {
  const { t } = useTranslation();
  const theme = useTheme();
  const router = useRouter();
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  const timeAgo = formatDistanceToNow(new Date(post.created_at), {
    addSuffix: true,
  });

  const fullDateTime = new Intl.DateTimeFormat(undefined, {
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(new Date(post.created_at));

  const handleViewDetails = () => {
    router.push(`/incognito-posts/${post.incognito_post_id}`);
  };

  const handleDeleteClick = () => {
    setDeleteDialogOpen(true);
  };

  const handleDeleteConfirm = async () => {
    try {
      // TODO: Implement delete API call
      setDeleteDialogOpen(false);
      onDeleted();
    } catch (error) {
      onError(t("incognitoPosts.errors.deleteFailed"));
    }
  };

  const formatCount = (count: number, singular: string, plural: string) => {
    return count === 1 ? `${count} ${singular}` : `${count} ${plural}`;
  };

  if (post.is_deleted) {
    return (
      <Card
        sx={{
          mb: 2.5,
          width: "100%",
          border: "none",
          boxShadow: "0 1px 2px rgba(0,0,0,0.06)",
          borderRadius: "8px",
          backgroundColor: theme.palette.background.paper,
          borderTop: `2px solid ${theme.palette.warning.main}`,
          borderLeft: `4px solid ${theme.palette.warning.main}`,
          opacity: 0.6,
        }}
      >
        <CardContent>
          <Typography variant="body2" color="text.secondary" fontStyle="italic">
            {t("incognitoPosts.post.deleted")}
          </Typography>
        </CardContent>
      </Card>
    );
  }

  return (
    <>
      <Card
        sx={{
          mb: 2.5,
          width: "100%",
          border: "none",
          boxShadow: "0 1px 2px rgba(0,0,0,0.06)",
          borderRadius: "8px",
          backgroundColor: theme.palette.background.paper,
          borderTop: `2px solid #7c3aed`,
          borderLeft: `4px solid #7c3aed`,
        }}
      >
        <CardHeader
          title={
            <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
              {post.is_created_by_me && (
                <Typography
                  variant="body2"
                  sx={{
                    color: "#6d28d9",
                    fontWeight: 500,
                    fontSize: "0.8rem",
                  }}
                >
                  {t("incognitoPosts.post.createdByYou")}
                </Typography>
              )}
              <Typography
                variant="body2"
                sx={{
                  color: theme.palette.text.secondary,
                  fontSize: "0.8rem",
                  lineHeight: 1.2,
                }}
                title={fullDateTime}
              >
                {timeAgo}
              </Typography>
            </Box>
          }
          action={
            post.is_created_by_me ? (
              <IconButton
                size="small"
                onClick={handleDeleteClick}
                color="error"
                title={t("incognitoPosts.post.deletePost")}
                sx={{ mt: -0.5 }}
              >
                <Delete fontSize="small" />
              </IconButton>
            ) : null
          }
          sx={{
            alignItems: "flex-start",
            p: 2,
            pb: 1,
            "& .MuiCardHeader-content": {
              overflow: "hidden",
            },
          }}
        />

        <CardContent sx={{ pt: 0.5, pb: "16px !important" }}>
          {/* Content */}
          <Typography
            variant="body1"
            component="p"
            whiteSpace="pre-wrap"
            sx={{
              color: theme.palette.text.primary,
              lineHeight: 1.5,
              fontSize: "0.9rem",
              mb: 1.5,
            }}
          >
            {post.content}
          </Typography>

          {/* Tags */}
          {post.tags && post.tags.length > 0 && (
            <Box
              sx={{
                mt: 1.5,
                mb: 2,
                display: "flex",
                flexWrap: "wrap",
                gap: 0.5,
              }}
            >
              {post.tags.map((tag) => (
                <Chip
                  key={tag.id}
                  label={tag.name}
                  size="small"
                  variant="filled"
                  clickable
                  sx={{
                    borderRadius: "16px",
                    backgroundColor: "#f3e8ff",
                    color: "#7c3aed",
                    fontSize: "0.75rem",
                    height: "24px",
                    border: "1px solid #c4b5fd",
                    "& .MuiChip-label": {
                      padding: "0 8px",
                      fontWeight: 500,
                    },
                    "&:hover": {
                      backgroundColor: "#ede9fe",
                      borderColor: "#a78bfa",
                    },
                    "&:focus": {
                      backgroundColor: "#ede9fe",
                    },
                  }}
                />
              ))}
            </Box>
          )}

          {/* Actions */}
          <Box
            sx={{
              display: "flex",
              justifyContent: "space-between",
              alignItems: "center",
            }}
          >
            {/* Voting */}
            <VotingButtons
              postId={post.incognito_post_id}
              upvotesCount={post.upvotes_count}
              downvotesCount={post.downvotes_count}
              score={post.score}
              meUpvoted={post.me_upvoted}
              meDownvoted={post.me_downvoted}
              canUpvote={post.can_upvote}
              canDownvote={post.can_downvote}
              onVoteUpdated={onVoteUpdated}
              onError={onError}
            />

            {/* Comments */}
            <Box
              sx={{
                display: "flex",
                alignItems: "center",
                gap: 0.5,
                cursor: "pointer",
                borderRadius: 1,
                px: 1,
                py: 0.5,
                "&:hover": {
                  backgroundColor: theme.palette.action.hover,
                },
              }}
              onClick={handleViewDetails}
            >
              <Comment fontSize="small" color="action" />
              <Typography variant="body2" color="text.secondary">
                {formatCount(
                  post.comments_count,
                  t("incognitoPosts.post.comments"),
                  t("incognitoPosts.post.commentsPlural")
                )}
              </Typography>
            </Box>
          </Box>
        </CardContent>
      </Card>

      {/* Delete Confirmation Dialog */}
      <Dialog
        open={deleteDialogOpen}
        onClose={() => setDeleteDialogOpen(false)}
      >
        <DialogTitle>{t("incognitoPosts.post.deletePost")}</DialogTitle>
        <DialogContent>
          <Typography>{t("incognitoPosts.post.confirmDelete")}</Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)}>
            {t("incognitoPosts.compose.cancel")}
          </Button>
          <Button
            onClick={handleDeleteConfirm}
            color="error"
            variant="contained"
          >
            {t("incognitoPosts.post.deletePost")}
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
}
