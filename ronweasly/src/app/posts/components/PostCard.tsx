"use client";

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
import Link from "next/link";

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

  return (
    <Card sx={{ mb: 2, width: "100%" }}>
      <CardHeader
        avatar={
          <Avatar aria-label="user avatar">
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
