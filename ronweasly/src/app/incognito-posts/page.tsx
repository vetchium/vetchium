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
  Typography,
} from "@mui/material";
import { Suspense, useState } from "react";
import BrowseTab from "./components/BrowseTab";
import CreatePostDialog from "./components/CreatePostDialog";
import MyCommentsTab from "./components/MyCommentsTab";
import MyPostsTab from "./components/MyPostsTab";

function IncognitoPostsContent() {
  const { t } = useTranslation();
  useAuth();
  const [tabValue, setTabValue] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [refreshTrigger, setRefreshTrigger] = useState(0);
  const [createDialogOpen, setCreateDialogOpen] = useState(false);

  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  const handlePostCreated = () => {
    setRefreshTrigger((prev) => prev + 1);
    setSuccess(t("incognitoPosts.success.postCreated"));
  };

  const handleError = (errorMessage: string) => {
    setError(errorMessage);
  };

  const handleSuccess = (successMessage: string) => {
    setSuccess(successMessage);
  };

  return (
    <Box sx={{ maxWidth: 1000, mx: "auto", mt: 4, px: 2 }}>
      {/* Header */}
      <Box sx={{ mb: 4, textAlign: "center" }}>
        <Typography variant="h4" component="h1" gutterBottom>
          {t("incognitoPosts.title")}
        </Typography>
        <Typography variant="subtitle1" color="text.secondary" gutterBottom>
          {t("incognitoPosts.subtitle")}
        </Typography>
        <Button
          variant="contained"
          color="primary"
          onClick={() => setCreateDialogOpen(true)}
          sx={{ mt: 2 }}
        >
          {t("incognitoPosts.createPost")}
        </Button>
      </Box>

      {/* Tabs Section */}
      <Paper sx={{ width: "100%" }}>
        <Box sx={{ borderBottom: 1, borderColor: "divider" }}>
          <Tabs
            value={tabValue}
            onChange={handleTabChange}
            aria-label="incognito posts tabs"
            variant="fullWidth"
          >
            <Tab
              label={t("incognitoPosts.browseByTags")}
              id="incognito-tab-0"
            />
            <Tab
              label={t("incognitoPosts.myPosts.title")}
              id="incognito-tab-1"
            />
            <Tab
              label={t("incognitoPosts.myComments.title")}
              id="incognito-tab-2"
            />
          </Tabs>
        </Box>

        {/* Tab Panels */}
        <Box sx={{ p: 3 }}>
          {tabValue === 0 && (
            <BrowseTab
              refreshTrigger={refreshTrigger}
              onError={handleError}
              onSuccess={handleSuccess}
            />
          )}
          {tabValue === 1 && (
            <MyPostsTab
              refreshTrigger={refreshTrigger}
              onError={handleError}
              onSuccess={handleSuccess}
            />
          )}
          {tabValue === 2 && (
            <MyCommentsTab
              refreshTrigger={refreshTrigger}
              onError={handleError}
              onSuccess={handleSuccess}
            />
          )}
        </Box>
      </Paper>

      {/* Create Post Dialog */}
      <CreatePostDialog
        open={createDialogOpen}
        onClose={() => setCreateDialogOpen(false)}
        onPostCreated={handlePostCreated}
        onError={handleError}
      />

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

export default function IncognitoPostsPage() {
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
        <IncognitoPostsContent />
      </Suspense>
    </AuthenticatedLayout>
  );
}
