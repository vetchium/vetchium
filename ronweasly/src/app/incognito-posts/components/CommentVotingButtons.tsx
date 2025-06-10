"use client";

import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import {
  ThumbDown,
  ThumbDownOutlined,
  ThumbUp,
  ThumbUpOutlined,
} from "@mui/icons-material";
import { Box, IconButton, Tooltip, Typography } from "@mui/material";
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
        if (response.status === 422) {
          throw new Error(t("incognitoPosts.voting.cannotVoteOwn"));
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
    <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
      {/* Upvote */}
      <Tooltip title={getUpvoteTooltip()}>
        <span>
          <IconButton
            size="small"
            onClick={handleUpvote}
            disabled={!canUpvote || isVoting}
            color={meUpvoted ? "primary" : "default"}
          >
            {meUpvoted ? (
              <ThumbUp fontSize="small" />
            ) : (
              <ThumbUpOutlined fontSize="small" />
            )}
          </IconButton>
        </span>
      </Tooltip>

      <Typography variant="body2" color="text.secondary" sx={{ minWidth: 20 }}>
        {upvotesCount}
      </Typography>

      {/* Downvote */}
      <Tooltip title={getDownvoteTooltip()}>
        <span>
          <IconButton
            size="small"
            onClick={handleDownvote}
            disabled={!canDownvote || isVoting}
            color={meDownvoted ? "error" : "default"}
          >
            {meDownvoted ? (
              <ThumbDown fontSize="small" />
            ) : (
              <ThumbDownOutlined fontSize="small" />
            )}
          </IconButton>
        </span>
      </Tooltip>

      <Typography variant="body2" color="text.secondary" sx={{ minWidth: 20 }}>
        {downvotesCount}
      </Typography>
    </Box>
  );
}
