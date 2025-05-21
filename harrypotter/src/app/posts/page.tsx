"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import ComposeSection from "@/components/posts/ComposeSection";
import TabPanel from "@/components/posts/TabPanel";
import TimelineTab from "@/components/posts/TimelineTab";
import { useTranslation } from "@/hooks/useTranslation";
import CloseIcon from "@mui/icons-material/Close";
import {
  Box,
  Button,
  CircularProgress,
  Paper,
  Snackbar,
  Tab,
  Tabs,
  Typography,
} from "@mui/material";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { Suspense, useEffect, useState } from "react";

function PostsContent() {
  const { t } = useTranslation();
  const router = useRouter();
  const [tabValue, setTabValue] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  // Check authentication on component mount
  useEffect(() => {
    const token = Cookies.get("session_token");
    if (!token) {
      router.push("/signin");
    }
  }, [router]);

  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  // Handler for when a new post is created
  const handlePostCreated = () => {
    // Increment the trigger to cause the timeline to refresh
    setRefreshTrigger((prev) => prev + 1);
    setSuccess(t("posts.createSuccess"));
  };

  return (
    <Box sx={{ maxWidth: 800, mx: "auto", mt: 4 }}>
      <Typography variant="h4" gutterBottom align="center">
        {t("posts.title")}
      </Typography>

      {/* Compose Section */}
      <ComposeSection
        onPostCreated={handlePostCreated}
        onError={setError}
        onSuccess={setSuccess}
      />

      {/* Tabs Section */}
      <Paper sx={{ width: "100%" }}>
        <Box sx={{ borderBottom: 1, borderColor: "divider" }}>
          <Tabs
            value={tabValue}
            onChange={handleTabChange}
            aria-label="posts tabs"
            variant="fullWidth"
          >
            <Tab label={t("posts.timeline")} id="posts-tab-0" />
          </Tabs>
        </Box>
        <TabPanel value={tabValue} index={0}>
          <TimelineTab refreshTrigger={refreshTrigger} onError={setError} />
        </TabPanel>
      </Paper>

      {/* Notifications */}
      <Snackbar
        open={!!error}
        autoHideDuration={6000}
        onClose={() => setError(null)}
        message={error}
        action={
          <Button color="inherit" size="small" onClick={() => setError(null)}>
            <CloseIcon fontSize="small" />
          </Button>
        }
      />

      <Snackbar
        open={!!success}
        autoHideDuration={6000}
        onClose={() => setSuccess(null)}
        message={success}
        action={
          <Button color="inherit" size="small" onClick={() => setSuccess(null)}>
            <CloseIcon fontSize="small" />
          </Button>
        }
      />
    </Box>
  );
}

export default function PostsPage() {
  return (
    <AuthenticatedLayout>
      <Suspense
        fallback={
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
        }
      >
        <PostsContent />
      </Suspense>
    </AuthenticatedLayout>
  );
}
