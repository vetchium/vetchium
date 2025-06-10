"use client";

import { useTranslation } from "@/hooks/useTranslation";
import BusinessIcon from "@mui/icons-material/Business";
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
  useTheme,
} from "@mui/material";
import { EmployerPost } from "@vetchium/typespec";
import { formatDistanceToNow } from "date-fns";
import Link from "next/link";

interface EmployerPostCardProps {
  post: EmployerPost;
  hideOpenInNewTab?: boolean;
}

export default function EmployerPostCard({
  post,
  hideOpenInNewTab = false,
}: EmployerPostCardProps) {
  const { t } = useTranslation();
  const theme = useTheme();
  const timeAgo = formatDistanceToNow(new Date(post.created_at), {
    addSuffix: true,
  });

  return (
    <Card
      sx={{
        mb: 2.5,
        width: "100%",
        border: "none",
        boxShadow: "0 1px 2px rgba(0,0,0,0.06)",
        borderRadius: "8px",
        backgroundColor: theme.palette.background.paper,
        borderTop: `2px solid ${theme.palette.primary.main}`,
        borderLeft: `4px solid ${theme.palette.primary.main}`,
      }}
    >
      <CardHeader
        avatar={
          <Link
            href={`https://${post.employer_domain_name}`}
            target="_blank"
            rel="noopener noreferrer"
            style={{ textDecoration: "none" }}
          >
            <Avatar
              aria-label="employer avatar"
              sx={{
                width: 48,
                height: 48,
                border: `1px solid ${theme.palette.divider}`,
                bgcolor: theme.palette.primary.light,
              }}
            >
              <BusinessIcon />
            </Avatar>
          </Link>
        }
        title={
          <Box sx={{ mb: 0.25 }}>
            <Link
              href={`https://${post.employer_domain_name}`}
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
                {post.employer_name}
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
                href={`https://${post.employer_domain_name}`}
                target="_blank"
                rel="noopener noreferrer"
                style={{ textDecoration: "none", color: "inherit" }}
              >
                {post.employer_domain_name}
              </Link>
              {` · ${timeAgo}`}
            </Typography>
          </Box>
        }
        action={
          !hideOpenInNewTab ? (
            <Tooltip title={t("common.externalLink.message")}>
              <IconButton
                component={Link}
                href={`/employer-posts/${post.id}`}
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
          ) : null
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
          <Box sx={{ mt: 1.5, display: "flex", flexWrap: "wrap", gap: 0.5 }}>
            {post.tags.map((tag) => (
              <Chip
                key={tag}
                label={tag}
                size="small"
                variant="filled"
                clickable
                sx={{
                  borderRadius: "16px",
                  backgroundColor: theme.palette.primary.main + "15",
                  color: theme.palette.primary.main,
                  fontSize: "0.75rem",
                  height: "24px",
                  border: `1px solid ${theme.palette.primary.main}30`,
                  "& .MuiChip-label": {
                    padding: "0 8px",
                    fontWeight: 500,
                  },
                  "&:hover": {
                    backgroundColor: theme.palette.primary.main + "25",
                    borderColor: theme.palette.primary.main + "50",
                  },
                  "&:focus": {
                    backgroundColor: theme.palette.primary.main + "25",
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
