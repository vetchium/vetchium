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
import {
  EmployerPost,
  GetEmployerPostDetailsRequest,
} from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import EmployerPostCard from "../../posts/components/EmployerPostCard";

export default function EmployerPostDetailPage() {
  const { t } = useTranslation();
  const router = useRouter();
  const params = useParams();
  const postId = (params?.id as string) || "";
  const [loading, setLoading] = useState(true);
  const [post, setPost] = useState<EmployerPost | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isNotFound, setIsNotFound] = useState(false);

  // Fetch employer post details on component mount
  useEffect(() => {
    const token = Cookies.get("session_token");
    if (!token) {
      router.push("/login");
      return;
    }

    const fetchEmployerPostDetails = async () => {
      try {
        const requestBody: GetEmployerPostDetailsRequest = {
          employer_post_id: postId,
        };

        const response = await fetch(
          `${config.API_SERVER_PREFIX}/hub/get-employer-post-details`,
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

        if (response.status === 401) {
          Cookies.remove("session_token", { path: "/" });
          router.push("/login");
          return;
        }

        if (!response.ok) {
          throw new Error(`Failed to fetch employer post: ${response.status}`);
        }

        const data: EmployerPost = await response.json();

        // Validate that we have a valid employer post object with required fields
        if (
          !data ||
          !data.id ||
          !data.content ||
          !data.employer_name ||
          !data.employer_domain_name
        ) {
          throw new Error("Invalid employer post data received");
        }

        // Ensure tags is always an array
        const safePost: EmployerPost = {
          ...data,
          tags: Array.isArray(data.tags) ? data.tags : [],
        };

        setPost(safePost);
        setLoading(false);
      } catch (err) {
        console.error("Error fetching employer post details:", err);
        setError(err instanceof Error ? err.message : "Failed to load post");
        setLoading(false);
      }
    };

    fetchEmployerPostDetails();
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
            {t("posts.viewPost") || "View Post"}
          </Typography>

          {post ? (
            <Box sx={{ mt: 3 }}>
              <EmployerPostCard post={post} hideOpenInNewTab={true} />
            </Box>
          ) : (
            <Typography>{t("posts.notFound") || "Post not found"}</Typography>
          )}
        </Paper>
      </Box>
    </AuthenticatedLayout>
  );
}
