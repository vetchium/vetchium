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
import { EmployerPost, GetPostDetailsRequest, Post } from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import EmployerPostCard from "../components/EmployerPostCard";
import PostCard from "../components/PostCard";

type PostType = Post | EmployerPost;

export default function PostDetailPage() {
  const { t } = useTranslation();
  const router = useRouter();
  const params = useParams();
  const postId = (params?.postId as string) || "";
  const [loading, setLoading] = useState(true);
  const [post, setPost] = useState<PostType | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isNotFound, setIsNotFound] = useState(false);

  // Function to determine if a post is an EmployerPost
  const isEmployerPost = (post: any): post is EmployerPost => {
    return post && "employer_name" in post && "employer_domain_name" in post;
  };

  // Fetch post details on component mount
  useEffect(() => {
    const token = Cookies.get("session_token");
    if (!token) {
      router.push("/login");
      return;
    }

    const fetchPostDetails = async () => {
      try {
        // Try fetching user post first
        const requestBody: GetPostDetailsRequest = {
          post_id: postId,
        };

        let response = await fetch(
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

        // If not found as a user post, try as an employer post
        if (response.status === 404) {
          try {
            response = await fetch(
              `${config.API_SERVER_PREFIX}/employer/get-post`,
              {
                method: "POST",
                headers: {
                  "Content-Type": "application/json",
                  Authorization: `Bearer ${token}`,
                },
                body: JSON.stringify({ post_id: postId }),
              }
            );
          } catch (err) {
            console.error("Error fetching employer post:", err);
          }
        }

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

        // Ensure tags is always an array
        const safePost = {
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
    // Check if there's history to go back to
    if (window.history.length > 1) {
      router.back();
    } else {
      // No history, navigate to home page
      router.push("/");
    }
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
              {isEmployerPost(post) ? (
                <EmployerPostCard post={post} hideOpenInNewTab={true} />
              ) : (
                <PostCard post={post as Post} hideOpenInNewTab={true} />
              )}
            </Box>
          ) : (
            <Typography>{t("posts.notFound")}</Typography>
          )}
        </Paper>
      </Box>
    </AuthenticatedLayout>
  );
}
