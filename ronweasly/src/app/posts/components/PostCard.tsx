"use client";

import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import LocalOfferIcon from "@mui/icons-material/LocalOffer";
import OpenInNewIcon from "@mui/icons-material/OpenInNew";
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
} from "@mui/material";
import { Post } from "@vetchium/typespec";
import { formatDistanceToNow } from "date-fns";
import Cookies from "js-cookie";
import Link from "next/link";
import { useEffect, useState } from "react";

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
  const timeAgo = formatDistanceToNow(new Date(post.created_at), {
    addSuffix: true,
  });
  const [avatarUrl, setAvatarUrl] = useState<string | null>(null);

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
        mb: 2,
        width: "100%",
        border: "1px solid #e0e0e0",
        boxShadow: "0 1px 3px rgba(0,0,0,0.08)",
        borderRadius: "8px",
        backgroundColor: "#ffffff",
      }}
    >
      <CardHeader
        avatar={
          <Avatar
            aria-label="user avatar"
            src={avatarUrl ?? undefined}
            sx={{
              width: 48,
              height: 48,
              border: "1px solid #e0e0e0",
            }}
          >
            {/* Fallback: Initials */}
            {post.author_name?.charAt(0) || post.author_handle.charAt(0)}
          </Avatar>
        }
        title={
          <Typography
            variant="subtitle1"
            component="span"
            sx={{
              fontWeight: 600,
              color: "#000000de",
            }}
          >
            {post.author_name || post.author_handle}
          </Typography>
        }
        subheader={
          <Typography
            variant="body2"
            component="span"
            sx={{
              color: "#00000099",
              fontSize: "0.875rem",
            }}
          >
            @{post.author_handle} · {timeAgo}
          </Typography>
        }
        action={
          !hideOpenInNewTab ? (
            <Tooltip title={t("common.externalLink.message")}>
              <IconButton
                component={Link}
                href={`/posts/${post.id}`}
                target="_blank"
                rel="noopener noreferrer"
                aria-label={t("common.externalLink.message")}
                sx={{ color: "#0000008a" }}
              >
                <OpenInNewIcon />
              </IconButton>
            </Tooltip>
          ) : null
        }
        sx={{
          p: 2,
          pb: 1,
          borderBottom: "1px solid #f5f5f5",
        }}
      />
      <CardContent sx={{ p: 2, pt: 1.5 }}>
        <Typography
          variant="body1"
          component="p"
          whiteSpace="pre-wrap"
          sx={{
            color: "#000000de",
            lineHeight: 1.5,
            fontSize: "0.9375rem",
            mb: 1,
          }}
        >
          {post.content}
        </Typography>
        {post.tags && Array.isArray(post.tags) && post.tags.length > 0 && (
          <Box sx={{ mt: 2, display: "flex", flexWrap: "wrap", gap: 1 }}>
            {post.tags.map((tag) => (
              <Chip
                key={tag}
                label={tag}
                size="small"
                icon={<LocalOfferIcon sx={{ fontSize: "0.875rem" }} />}
                variant="outlined"
                sx={{
                  borderRadius: "16px",
                  backgroundColor: "#f5f5f5",
                  borderColor: "#e0e0e0",
                  "& .MuiChip-label": {
                    fontSize: "0.75rem",
                    color: "#000000de",
                  },
                }}
              />
            ))}
          </Box>
        )}
      </CardContent>
    </Card>
  );
}
