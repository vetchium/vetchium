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

    const fetchAvatar = async () => {
      const token = Cookies.get("session_token");
      // No need to fetch if token is missing, fallback will be used
      if (!token || !post.author_handle) {
        setAvatarUrl(null); // Ensure avatarUrl is null if no token or handle
        return;
      }

      const imageUrl = `${config.API_SERVER_PREFIX}/hub/profile-picture/${post.author_handle}`;

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
          } else {
            setAvatarUrl(null); // No image available, use fallback
          }
        } else {
          // Handle errors like 401, 404, etc. - use fallback
          console.error(
            `Failed to fetch avatar for ${post.author_handle}: ${response.status}`
          );
          setAvatarUrl(null);
        }
      } catch (error) {
        console.error(
          `Error fetching avatar for ${post.author_handle}:`,
          error
        );
        setAvatarUrl(null); // Network or other errors, use fallback
      }
    };

    fetchAvatar();

    // Cleanup function
    return () => {
      if (objectUrl) {
        URL.revokeObjectURL(objectUrl);
        setAvatarUrl(null); // Clear state on cleanup as well
      }
    };
  }, [post.author_handle]); // Re-run effect if author handle changes

  return (
    <Card sx={{ mb: 2, width: "100%" }}>
      <CardHeader
        avatar={
          <Avatar
            aria-label="user avatar"
            src={avatarUrl ?? undefined} // Use object URL if available, otherwise undefined lets fallback render
          >
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
