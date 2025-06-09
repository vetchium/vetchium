"use client";

import { useTranslation } from "@/hooks/useTranslation";
import { Comment, Delete } from "@mui/icons-material";
import {
  Box,
  Button,
  Card,
  CardContent,
  Chip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  Typography,
} from "@mui/material";
import { IncognitoPostSummary } from "@vetchium/typespec";
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
  const router = useRouter();
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

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

  const formatScore = (score: number) => {
    if (score === 0) return "0";
    return score > 0 ? `+${score}` : `${score}`;
  };

  const formatCount = (count: number, singular: string, plural: string) => {
    return count === 1 ? `${count} ${singular}` : `${count} ${plural}`;
  };

  if (post.is_deleted) {
    return (
      <Card sx={{ opacity: 0.6 }}>
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
      <Card sx={{ mb: 2 }}>
        <CardContent>
          {/* Header */}
          <Box
            sx={{
              display: "flex",
              justifyContent: "space-between",
              alignItems: "flex-start",
              mb: 2,
            }}
          >
            <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
              <Typography variant="body2" color="text.secondary">
                {post.is_created_by_me
                  ? t("incognitoPosts.post.createdBy")
                  : t("incognitoPosts.post.anonymous")}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                •
              </Typography>
              <Typography variant="body2" color="text.secondary">
                {new Date(post.created_at).toLocaleDateString()}
              </Typography>
            </Box>

            {post.is_created_by_me && (
              <IconButton
                size="small"
                onClick={handleDeleteClick}
                color="error"
                title={t("incognitoPosts.post.deletePost")}
              >
                <Delete fontSize="small" />
              </IconButton>
            )}
          </Box>

          {/* Content */}
          <Typography variant="body1" sx={{ mb: 2, whiteSpace: "pre-wrap" }}>
            {post.content}
          </Typography>

          {/* Tags */}
          <Box sx={{ display: "flex", flexWrap: "wrap", gap: 1, mb: 2 }}>
            {post.tags.map((tag) => (
              <Chip
                key={tag.id}
                label={tag.name}
                size="small"
                variant="outlined"
              />
            ))}
          </Box>

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

            {/* Comments and View Details */}
            <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
              <Box sx={{ display: "flex", alignItems: "center", gap: 0.5 }}>
                <Comment fontSize="small" color="action" />
                <Typography variant="body2" color="text.secondary">
                  {formatCount(
                    post.comments_count,
                    t("incognitoPosts.post.comments"),
                    t("incognitoPosts.post.commentsPlural")
                  )}
                </Typography>
              </Box>

              <Button size="small" onClick={handleViewDetails}>
                {t("incognitoPosts.post.viewDetails")}
              </Button>
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
