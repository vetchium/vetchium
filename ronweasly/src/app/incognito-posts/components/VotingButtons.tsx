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
import {
  DownvoteIncognitoPostRequest,
  UnvoteIncognitoPostRequest,
  UpvoteIncognitoPostRequest,
} from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useState } from "react";

interface VotingButtonsProps {
  postId: string;
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

export default function VotingButtons({
  postId,
  upvotesCount,
  downvotesCount,
  score,
  meUpvoted,
  meDownvoted,
  canUpvote,
  canDownvote,
  onVoteUpdated,
  onError,
}: VotingButtonsProps) {
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
      let request: any = {};

      switch (action) {
        case "upvote":
          endpoint = "/hub/upvote-incognito-post";
          request = new UpvoteIncognitoPostRequest();
          request.incognito_post_id = postId;
          break;
        case "downvote":
          endpoint = "/hub/downvote-incognito-post";
          request = new DownvoteIncognitoPostRequest();
          request.incognito_post_id = postId;
          break;
        case "unvote":
          endpoint = "/hub/unvote-incognito-post";
          request = new UnvoteIncognitoPostRequest();
          request.incognito_post_id = postId;
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
          throw new Error(
            "This post has been deleted and can no longer be voted on."
          );
        } else if (response.status === 422) {
          if (action === "upvote" || action === "downvote") {
            throw new Error(
              "Cannot vote on this post. You may have already voted in the opposite direction or this is your own post."
            );
          } else {
            throw new Error(
              "Cannot remove vote from this post. This may be your own post."
            );
          }
        } else if (response.status === 401) {
          throw new Error("You must be logged in to vote.");
        } else {
          throw new Error(`Failed to ${action}. Please try again.`);
        }
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
    if (isVoting) return;
    handleVote(meUpvoted ? "unvote" : "upvote");
  };

  const handleDownvote = () => {
    if (isVoting) return;
    handleVote(meDownvoted ? "unvote" : "downvote");
  };

  const getUpvoteTooltip = () => {
    if (meDownvoted) return "Remove your downvote to upvote";
    if (!canUpvote && !meUpvoted) return "You cannot vote on your own post";
    if (meUpvoted) return t("incognitoPosts.voting.unvote");
    return t("incognitoPosts.voting.upvote");
  };

  const getDownvoteTooltip = () => {
    if (meUpvoted) return "Remove your upvote to downvote";
    if (!canDownvote && !meDownvoted) return "You cannot vote on your own post";
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
            disabled={isVoting || meDownvoted || (!canUpvote && !meUpvoted)}
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
            disabled={isVoting || meUpvoted || (!canDownvote && !meDownvoted)}
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
