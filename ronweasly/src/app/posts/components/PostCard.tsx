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
    <Card sx={{ mb: 2, width: "100%" }}>
      <CardHeader
        avatar={
          <Avatar aria-label="user avatar" src={avatarUrl ?? undefined}>
            {/* Fallback: Initials */}
            {post.author_name?.charAt(0) || post.author_handle.charAt(0)}
          </Avatar>
        }
        title={post.author_name || post.author_handle}
        subheader={`@${post.author_handle} · ${timeAgo}`}
        action={
          !hideOpenInNewTab ? (
            <Tooltip title={t("common.externalLink.message")}>
              <IconButton
                component={Link}
                href={`/posts/${post.id}`}
                target="_blank"
                rel="noopener noreferrer"
                aria-label={t("common.externalLink.message")}
              >
                <OpenInNewIcon />
              </IconButton>
            </Tooltip>
          ) : null
        }
      />
      <CardContent>
        <Typography variant="body1" component="p" whiteSpace="pre-wrap">
          {post.content}
        </Typography>
        {post.tags && Array.isArray(post.tags) && post.tags.length > 0 && (
          <Box sx={{ mt: 2, display: "flex", flexWrap: "wrap", gap: 1 }}>
            {post.tags.map((tag) => (
              <Chip
                key={tag}
                label={tag}
                size="small"
                icon={<LocalOfferIcon />}
                variant="outlined"
              />
            ))}
          </Box>
        )}
      </CardContent>
    </Card>
  );
}
