"use client";

import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import { Box, Button, Typography } from "@mui/material";
import Cookies from "js-cookie";
import { useState } from "react";

interface CommentVotingButtonsProps {
  postId: string;
  commentId: string;
  upvotesCount: number;
  downvotesCount: number;
  score: number;
  meUpvoted: boolean;
  meDownvoted: boolean;
  canUpvote: boolean;
  canDownvote: boolean;
  onVoteUpdated: () => void;
  onError: (error: string) => void;
}

export default function CommentVotingButtons({
  postId,
  commentId,
  upvotesCount,
  downvotesCount,
  score,
  meUpvoted,
  meDownvoted,
  canUpvote,
  canDownvote,
  onVoteUpdated,
  onError,
}: CommentVotingButtonsProps) {
  const { t } = useTranslation();
  const [isVoting, setIsVoting] = useState(false);

  const handleVote = async (action: "upvote" | "downvote" | "unvote") => {
    if (isVoting) return;

    setIsVoting(true);
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        throw new Error("User not authenticated");
      }

      let endpoint = "";
      const request = {
        incognito_post_id: postId,
        comment_id: commentId,
      };

      switch (action) {
        case "upvote":
          endpoint = "/hub/upvote-incognito-post-comment";
          break;
        case "downvote":
          endpoint = "/hub/downvote-incognito-post-comment";
          break;
        case "unvote":
          endpoint = "/hub/unvote-incognito-post-comment";
          break;
      }

      const response = await fetch(`${config.API_SERVER_PREFIX}${endpoint}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(request),
      });

      if (!response.ok) {
        if (response.status === 404) {
          throw new Error(t("incognitoPosts.errors.commentNotFound"));
        } else if (response.status === 422) {
          throw new Error(t("incognitoPosts.voting.cannotVoteOwn"));
        } else if (response.status === 401) {
          throw new Error(t("incognitoPosts.errors.mustBeLoggedIn"));
        }
        throw new Error(`Failed to ${action}: ${response.statusText}`);
      }

      onVoteUpdated();
    } catch (error) {
      onError(
        error instanceof Error
          ? error.message
          : t("incognitoPosts.voting.votingError")
      );
    } finally {
      setIsVoting(false);
    }
  };

  const handleUpvote = () => {
    if (meUpvoted) {
      handleVote("unvote");
    } else {
      handleVote("upvote");
    }
  };

  const handleDownvote = () => {
    if (meDownvoted) {
      handleVote("unvote");
    } else {
      handleVote("downvote");
    }
  };

  const getUpvoteTooltip = () => {
    if (!canUpvote) return t("incognitoPosts.voting.cannotVoteOwn");
    if (meUpvoted) return t("incognitoPosts.voting.unvote");
    return t("incognitoPosts.voting.upvote");
  };

  const getDownvoteTooltip = () => {
    if (!canDownvote) return t("incognitoPosts.voting.cannotVoteOwn");
    if (meDownvoted) return t("incognitoPosts.voting.unvote");
    return t("incognitoPosts.voting.downvote");
  };

  return (
    <Box sx={{ display: "flex", alignItems: "center", gap: 0.5 }}>
      {/* Upvote Button */}
      <Button
        size="small"
        onClick={handleUpvote}
        disabled={!canUpvote || isVoting}
        sx={{
          minWidth: "auto",
          p: 0.25,
          fontSize: "0.75rem",
          color: meUpvoted ? "primary.main" : "text.secondary",
          textTransform: "none",
          "&:hover": {
            backgroundColor: "action.hover",
          },
        }}
        aria-label={t("incognitoPosts.voting.upvote")}
      >
        ▲
      </Button>

      {/* Score */}
      <Typography
        variant="caption"
        sx={{
          mx: 0.5,
          fontSize: "0.75rem",
          fontWeight: meUpvoted || meDownvoted ? "bold" : "normal",
          color: meUpvoted
            ? "primary.main"
            : meDownvoted
            ? "error.main"
            : "text.secondary",
          minWidth: "20px",
          textAlign: "center",
        }}
      >
        {score}
      </Typography>

      {/* Downvote Button */}
      <Button
        size="small"
        onClick={handleDownvote}
        disabled={!canDownvote || isVoting}
        sx={{
          minWidth: "auto",
          p: 0.25,
          fontSize: "0.75rem",
          color: meDownvoted ? "error.main" : "text.secondary",
          textTransform: "none",
          "&:hover": {
            backgroundColor: "action.hover",
          },
        }}
        aria-label={t("incognitoPosts.voting.downvote")}
      >
        ▼
      </Button>
    </Box>
  );
}
