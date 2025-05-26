"use client";

import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import OpenInNewIcon from "@mui/icons-material/OpenInNew";
import ThumbDownIcon from "@mui/icons-material/ThumbDown";
import ThumbDownOutlinedIcon from "@mui/icons-material/ThumbDownOutlined";
import ThumbUpIcon from "@mui/icons-material/ThumbUp";
import ThumbUpOutlinedIcon from "@mui/icons-material/ThumbUpOutlined";
import {
  Avatar,
  Box,
  Card,
  CardContent,
  CardHeader,
  Chip,
  IconButton,
  Tooltip,
  Typography,
  useTheme,
} from "@mui/material";
import { Post } from "@vetchium/typespec";
import { formatDistanceToNow } from "date-fns";
import Cookies from "js-cookie";
import Link from "next/link";
import { useEffect, useState } from "react";
import CommentSettings from "./CommentSettings";
import Comments from "./Comments";

// Cache for profile pictures with timestamps to enable expiration
interface ProfilePictureCacheEntry {
  url: string | null;
  timestamp: number;
}

const profilePictureCache: Record<string, ProfilePictureCacheEntry> = {};
const CACHE_EXPIRATION_MS = 20 * 60 * 1000; // 20 minutes cache expiration

interface PostCardProps {
  post: Post;
  hideOpenInNewTab?: boolean;
}

export default function PostCard({
  post,
  hideOpenInNewTab = false,
}: PostCardProps) {
  const { t } = useTranslation();
  const theme = useTheme();
  const timeAgo = formatDistanceToNow(new Date(post.created_at), {
    addSuffix: true,
  });
  const [avatarUrl, setAvatarUrl] = useState<string | null>(null);
  const [isUpvoted, setIsUpvoted] = useState(post.me_upvoted);
  const [isDownvoted, setIsDownvoted] = useState(post.me_downvoted);
  const [upvotesCount, setUpvotesCount] = useState(post.upvotes_count);
  const [downvotesCount, setDownvotesCount] = useState(post.downvotes_count);
  const [canUpvote, setCanUpvote] = useState(post.can_upvote);
  const [canDownvote, setCanDownvote] = useState(post.can_downvote);
  const [canComment, setCanComment] = useState(post.can_comment);
  const [commentsCount, setCommentsCount] = useState(post.comments_count);

  useEffect(() => {
    let objectUrl: string | null = null;
    const authorHandle = post.author_handle;

    const fetchAvatar = async () => {
      const token = Cookies.get("session_token");
      // No need to fetch if token is missing, fallback will be used
      if (!token || !authorHandle) {
        setAvatarUrl(null);
        return;
      }

      // Check if this author's profile picture is already in cache and not expired
      const now = Date.now();
      const cachedEntry = profilePictureCache[authorHandle];
      if (cachedEntry && now - cachedEntry.timestamp < CACHE_EXPIRATION_MS) {
        setAvatarUrl(cachedEntry.url);
        return;
      }

      const imageUrl = `${config.API_SERVER_PREFIX}/hub/profile-picture/${authorHandle}`;

      try {
        const response = await fetch(imageUrl, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.ok) {
          const blob = await response.blob();
          // Check if blob has size, empty blob indicates no image
          if (blob.size > 0) {
            objectUrl = URL.createObjectURL(blob);
            setAvatarUrl(objectUrl);
            // Store in cache with timestamp
            profilePictureCache[authorHandle] = {
              url: objectUrl,
              timestamp: now,
            };
          } else {
            setAvatarUrl(null);
            profilePictureCache[authorHandle] = {
              url: null,
              timestamp: now,
            };
          }
        } else {
          // Handle errors like 401, 404, etc. - use fallback
          setAvatarUrl(null);
          profilePictureCache[authorHandle] = {
            url: null,
            timestamp: now,
          };
        }
      } catch (error) {
        console.error(`Error fetching avatar for ${authorHandle}:`, error);
        setAvatarUrl(null);
        profilePictureCache[authorHandle] = {
          url: null,
          timestamp: now,
        };
      }
    };

    fetchAvatar();

    // Cleanup function
    return () => {
      // Only revoke URLs that aren't used in the cache or if they're expired
      const now = Date.now();
      const cachedEntry = profilePictureCache[authorHandle];
      const isExpired =
        cachedEntry && now - cachedEntry.timestamp >= CACHE_EXPIRATION_MS;

      if (
        objectUrl &&
        (!cachedEntry || isExpired || cachedEntry.url !== objectUrl)
      ) {
        URL.revokeObjectURL(objectUrl);
      }
    };
  }, [post.author_handle]);

  return (
    <Card
      sx={{
        mb: 2.5,
        width: "100%",
        border: "none",
        boxShadow: "0 1px 2px rgba(0,0,0,0.06)",
        borderRadius: "8px",
        backgroundColor: theme.palette.background.paper,
      }}
    >
      <CardHeader
        avatar={
          <Link
            href={`/u/${post.author_handle}`}
            target="_blank"
            rel="noopener noreferrer"
            style={{ textDecoration: "none" }}
          >
            <Avatar
              aria-label="user avatar"
              src={avatarUrl ?? undefined}
              sx={{
                width: 48,
                height: 48,
                border: `1px solid ${theme.palette.divider}`,
              }}
            >
              {/* Fallback: Initials */}
              {post.author_name?.charAt(0) || post.author_handle.charAt(0)}
            </Avatar>
          </Link>
        }
        title={
          <Box sx={{ mb: 0.25 }}>
            <Link
              href={`/u/${post.author_handle}`}
              target="_blank"
              rel="noopener noreferrer"
              style={{ textDecoration: "none", color: "inherit" }}
            >
              <Typography
                variant="subtitle1"
                component="span"
                sx={{
                  fontWeight: 500,
                  color: theme.palette.text.primary,
                  lineHeight: 1.3,
                }}
              >
                {post.author_name || post.author_handle}
              </Typography>
            </Link>
          </Box>
        }
        subheader={
          <Box>
            <Typography
              variant="body2"
              component="span"
              sx={{
                color: theme.palette.text.secondary,
                fontSize: "0.8rem",
                lineHeight: 1.2,
              }}
            >
              <Link
                href={`/u/${post.author_handle}`}
                target="_blank"
                rel="noopener noreferrer"
                style={{ textDecoration: "none", color: "inherit" }}
              >
                @{post.author_handle}
              </Link>
              {` · ${timeAgo}`}
            </Typography>
          </Box>
        }
        action={
          <Box sx={{ display: "flex", alignItems: "center", gap: 0.5 }}>
            {post.am_i_author && (
              <CommentSettings
                postId={post.id}
                canComment={canComment}
                onCommentSettingsChange={setCanComment}
              />
            )}
            {!hideOpenInNewTab && (
              <Tooltip title={t("common.externalLink.message")}>
                <IconButton
                  component={Link}
                  href={`/posts/${post.id}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  aria-label={t("common.externalLink.message")}
                  sx={{
                    color: theme.palette.text.secondary,
                    mt: -0.5,
                  }}
                >
                  <OpenInNewIcon sx={{ fontSize: "1.125rem" }} />
                </IconButton>
              </Tooltip>
            )}
          </Box>
        }
        sx={{
          alignItems: "flex-start",
          p: 2,
          "& .MuiCardHeader-content": {
            overflow: "hidden",
          },
        }}
      />
      <CardContent sx={{ pt: 0.5, pb: "16px !important" }}>
        <Box sx={{ display: "flex", alignItems: "center", mb: 1.5, gap: 2 }}>
          <Box sx={{ display: "flex", alignItems: "center", gap: 0.5 }}>
            <IconButton
              onClick={async () => {
                const token = Cookies.get("session_token");
                if (!token) return;

                try {
                  const endpoint = isUpvoted
                    ? `${config.API_SERVER_PREFIX}/hub/unvote-user-post`
                    : `${config.API_SERVER_PREFIX}/hub/upvote-user-post`;

                  const response = await fetch(endpoint, {
                    method: "POST",
                    headers: {
                      Authorization: `Bearer ${token}`,
                      "Content-Type": "application/json",
                    },
                    body: JSON.stringify({ post_id: post.id }),
                  });

                  if (response.ok) {
                    if (isUpvoted) {
                      setUpvotesCount((prev) => prev - 1);
                      setIsUpvoted(false);
                      setCanUpvote(true);
                      setCanDownvote(true);
                    } else if (isDownvoted) {
                      setDownvotesCount((prev) => prev - 1);
                      setIsDownvoted(false);
                      setUpvotesCount((prev) => prev + 1);
                      setIsUpvoted(true);
                      setCanUpvote(false);
                      setCanDownvote(false);
                      setCanUpvote(false);
                      setCanDownvote(false);
                    } else {
                      setUpvotesCount((prev) => prev + 1);
                      setIsUpvoted(true);
                      setCanUpvote(false);
                      setCanDownvote(false);
                    }
                  }
                } catch (error) {
                  console.error("Error voting:", error);
                }
              }}
              disabled={!canUpvote && !isUpvoted && !post.am_i_author}
              size="small"
              sx={{
                color: isUpvoted
                  ? theme.palette.primary.main
                  : theme.palette.text.secondary,
                "&:hover": {
                  color: theme.palette.primary.main,
                },
              }}
            >
              {isUpvoted ? (
                <ThumbUpIcon fontSize="small" />
              ) : (
                <ThumbUpOutlinedIcon fontSize="small" />
              )}
            </IconButton>
            <Typography
              variant="body2"
              sx={{ color: theme.palette.text.secondary }}
            >
              {upvotesCount}
            </Typography>
          </Box>

          <Box sx={{ display: "flex", alignItems: "center", gap: 0.5 }}>
            <IconButton
              onClick={async () => {
                const token = Cookies.get("session_token");
                if (!token) return;

                try {
                  const endpoint = isDownvoted
                    ? `${config.API_SERVER_PREFIX}/hub/unvote-user-post`
                    : `${config.API_SERVER_PREFIX}/hub/downvote-user-post`;

                  const response = await fetch(endpoint, {
                    method: "POST",
                    headers: {
                      Authorization: `Bearer ${token}`,
                      "Content-Type": "application/json",
                    },
                    body: JSON.stringify({ post_id: post.id }),
                  });

                  if (response.ok) {
                    if (isDownvoted) {
                      setDownvotesCount((prev) => prev - 1);
                      setIsDownvoted(false);
                      setCanUpvote(true);
                      setCanDownvote(true);
                    } else if (isUpvoted) {
                      setUpvotesCount((prev) => prev - 1);
                      setIsUpvoted(false);
                      setDownvotesCount((prev) => prev + 1);
                      setIsDownvoted(true);
                      setCanUpvote(false);
                      setCanDownvote(false);
                      setCanUpvote(false);
                      setCanDownvote(false);
                    } else {
                      setDownvotesCount((prev) => prev + 1);
                      setIsDownvoted(true);
                      setCanUpvote(false);
                      setCanDownvote(false);
                    }
                  }
                } catch (error) {
                  console.error("Error voting:", error);
                }
              }}
              disabled={!canDownvote && !isDownvoted && !post.am_i_author}
              size="small"
              sx={{
                color: isDownvoted
                  ? theme.palette.error.main
                  : theme.palette.text.secondary,
                "&:hover": {
                  color: theme.palette.error.main,
                },
              }}
            >
              {isDownvoted ? (
                <ThumbDownIcon fontSize="small" />
              ) : (
                <ThumbDownOutlinedIcon fontSize="small" />
              )}
            </IconButton>
            <Typography
              variant="body2"
              sx={{ color: theme.palette.text.secondary }}
            >
              {downvotesCount}
            </Typography>
          </Box>
        </Box>
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
        {post.tags && Array.isArray(post.tags) && post.tags.length > 0 && (
          <Box sx={{ mt: 1.5, display: "flex", flexWrap: "wrap", gap: 0.75 }}>
            {post.tags.map((tag) => (
              <Chip
                key={tag}
                label={`#${tag}`}
                size="small"
                variant="outlined"
                clickable
                sx={{
                  borderRadius: "4px",
                  backgroundColor: "transparent",
                  borderColor: "transparent",
                  color: theme.palette.text.secondary,
                  fontSize: "0.75rem",
                  p: "0 4px",
                  height: "auto",
                  "& .MuiChip-label": {
                    padding: "0",
                  },
                  "&:hover": {
                    backgroundColor: theme.palette.action.hover,
                  },
                }}
              />
            ))}
          </Box>
        )}

        {/* Comments section */}
        <Comments
          postId={post.id}
          commentsCount={commentsCount}
          canComment={canComment}
          amIAuthor={post.am_i_author}
          onCommentsCountChange={setCommentsCount}
        />
      </CardContent>
    </Card>
  );
}
