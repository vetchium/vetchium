"use client";

import { useTranslation } from "@/hooks/useTranslation";
import LocalOfferIcon from "@mui/icons-material/LocalOffer";
import {
  Avatar,
  Box,
  Card,
  CardContent,
  CardHeader,
  Chip,
  Typography,
} from "@mui/material";
import { Post } from "@vetchium/typespec";
import { formatDistanceToNow } from "date-fns";

interface PostCardProps {
  post: Post;
}

export default function PostCard({ post }: PostCardProps) {
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
