"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import { useTranslation } from "@/hooks/useTranslation";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import {
  Box,
  Button,
  Card,
  CardContent,
  CircularProgress,
  Paper,
  Typography,
} from "@mui/material";
import Cookies from "js-cookie";
import { useParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";

export default function PostDetailPage() {
  const { t } = useTranslation();
  const router = useRouter();
  const params = useParams();
  const postId = (params?.postId as string) || "";
  const [loading, setLoading] = useState(true);

  // Check authentication on component mount
  useEffect(() => {
    const token = Cookies.get("session_token");
    if (!token) {
      router.push("/login");
      return;
    }

    // In a real implementation, we would fetch the post data here
    // For now, just simulate a loading state
    const timer = setTimeout(() => {
      setLoading(false);
    }, 500);

    return () => clearTimeout(timer);
  }, [router]);

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

          <Card sx={{ mt: 3 }}>
            <CardContent>
              <Typography variant="body2" color="text.secondary" gutterBottom>
                {t("posts.postId")}:
              </Typography>
              <Typography variant="body1" component="div" sx={{ mb: 2 }}>
                {postId}
              </Typography>

              <Typography variant="body2" color="text.secondary" gutterBottom>
                {t("posts.content")}:
              </Typography>
              <Typography variant="body1">
                {t("posts.contentPlaceholder")}
              </Typography>
            </CardContent>
          </Card>

          <Box sx={{ mt: 4 }}>
            <Typography variant="body2" color="text.secondary">
              {t("posts.detailsComingSoon")}
            </Typography>
          </Box>
        </Paper>
      </Box>
    </AuthenticatedLayout>
  );
}
