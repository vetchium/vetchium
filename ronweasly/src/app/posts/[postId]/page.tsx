"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import InfoOutlinedIcon from "@mui/icons-material/InfoOutlined";
import {
  Alert,
  Box,
  Button,
  CircularProgress,
  Paper,
  Typography,
} from "@mui/material";
import { GetPostDetailsRequest, Post } from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import PostCard from "../components/PostCard";

export default function PostDetailPage() {
  const { t } = useTranslation();
  const router = useRouter();
  const params = useParams();
  const postId = (params?.postId as string) || "";
  const [loading, setLoading] = useState(true);
  const [post, setPost] = useState<Post | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isNotFound, setIsNotFound] = useState(false);

  // Fetch post details on component mount
  useEffect(() => {
    const token = Cookies.get("session_token");
    if (!token) {
      router.push("/login");
      return;
    }

    const fetchPostDetails = async () => {
      try {
        // Using the GetPostDetailsRequest structure from typespec
        const requestBody: GetPostDetailsRequest = {
          post_id: postId,
        };

        const response = await fetch(
          `${config.API_SERVER_PREFIX}/hub/get-post-details`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify(requestBody),
          }
        );

        if (response.status === 404) {
          setIsNotFound(true);
          setError(t("posts.notFoundError") || "Post not found");
          setLoading(false);
          return;
        }

        if (!response.ok) {
          throw new Error(`Failed to fetch post: ${response.status}`);
        }

        const data = await response.json();

        // Check if the response contains the post directly or as a nested object
        const postData = data.post || data;

        // Validate that we have a valid post object with required fields
        if (!postData || !postData.id || !postData.content) {
          throw new Error("Invalid post data received");
        }

        // Ensure post.tags is always an array
        const safePost: Post = {
          ...postData,
          tags: Array.isArray(postData.tags) ? postData.tags : [],
        };

        setPost(safePost);
        setLoading(false);
      } catch (err) {
        console.error("Error fetching post details:", err);
        setError(err instanceof Error ? err.message : "Failed to load post");
        setLoading(false);
      }
    };

    fetchPostDetails();
  }, [postId, router, t]);

  const handleBack = () => {
    router.back();
  };

  if (loading) {
    return (
      <AuthenticatedLayout>
        <Box
          sx={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            minHeight: "50vh",
          }}
        >
          <CircularProgress />
        </Box>
      </AuthenticatedLayout>
    );
  }

  if (isNotFound) {
    return (
      <AuthenticatedLayout>
        <Box sx={{ maxWidth: 800, mx: "auto", mt: 4 }}>
          <Button
            startIcon={<ArrowBackIcon />}
            onClick={handleBack}
            sx={{ mb: 2 }}
          >
            {t("common.back")}
          </Button>
          <Paper sx={{ p: 3, mb: 4 }}>
            <Alert severity="info" icon={<InfoOutlinedIcon />} sx={{ mb: 2 }}>
              {t("posts.notFoundError") || "Post not found"}
            </Alert>
            <Typography variant="body1">
              {t("posts.notFoundDescription") ||
                "The post you're looking for could not be found. It may have been deleted or you may not have permission to view it."}
            </Typography>
          </Paper>
        </Box>
      </AuthenticatedLayout>
    );
  }

  if (error && !isNotFound) {
    return (
      <AuthenticatedLayout>
        <Box sx={{ maxWidth: 800, mx: "auto", mt: 4 }}>
          <Button
            startIcon={<ArrowBackIcon />}
            onClick={handleBack}
            sx={{ mb: 2 }}
          >
            {t("common.back")}
          </Button>
          <Paper sx={{ p: 3, mb: 4 }}>
            <Typography color="error">{error}</Typography>
          </Paper>
        </Box>
      </AuthenticatedLayout>
    );
  }

  return (
    <AuthenticatedLayout>
      <Box sx={{ maxWidth: 800, mx: "auto", mt: 4 }}>
        <Button
          startIcon={<ArrowBackIcon />}
          onClick={handleBack}
          sx={{ mb: 2 }}
        >
          {t("common.back")}
        </Button>

        <Paper sx={{ p: 3, mb: 4 }}>
          <Typography variant="h4" gutterBottom>
            {t("posts.viewPost")}
          </Typography>

          {post ? (
            <Box sx={{ mt: 3 }}>
              <PostCard post={post} />
            </Box>
          ) : (
            <Typography>{t("posts.notFound")}</Typography>
          )}
        </Paper>
      </Box>
    </AuthenticatedLayout>
  );
}
