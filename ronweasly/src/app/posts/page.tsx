"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import { useAuth } from "@/hooks/useAuth";
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
} from "@mui/material";
import { useRouter } from "next/navigation";
import { Suspense, useState } from "react";
import ComposeSection from "./components/ComposeSection";
import TabPanel from "./components/TabPanel";
import TimelineTab from "./components/TimelineTab";

function PostsContent() {
  const { t } = useTranslation();
  const router = useRouter();
  useAuth(); // Check authentication and redirect if not authenticated
  const [tabValue, setTabValue] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  // Handler for when a new post is created
  const handlePostCreated = () => {
    // Increment the trigger to cause the timeline to refresh
    setRefreshTrigger((prev) => prev + 1);
  };

  return (
    <Box sx={{ maxWidth: 800, mx: "auto", mt: 4 }}>
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
            <Tab label={t("posts.following")} id="posts-tab-0" />
            <Tab label={t("posts.trending")} id="posts-tab-1" />
          </Tabs>
        </Box>
        <TabPanel value={tabValue} index={0}>
          <TimelineTab type="following" refreshTrigger={refreshTrigger} />
        </TabPanel>
        <TabPanel value={tabValue} index={1}>
          <TimelineTab type="trending" refreshTrigger={refreshTrigger} />
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
