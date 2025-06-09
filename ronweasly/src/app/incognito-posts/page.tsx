"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import { useAuth } from "@/hooks/useAuth";
import { useTranslation } from "@/hooks/useTranslation";
import CloseIcon from "@mui/icons-material/Close";
import {
  Box,
  Button,
  CircularProgress,
  Container,
  Paper,
  Snackbar,
  Tab,
  Tabs,
  Typography,
} from "@mui/material";
import React, { Suspense, useState } from "react";
import BrowseTab from "./components/BrowseTab";
import CreatePostDialog from "./components/CreatePostDialog";
import MyCommentsTab from "./components/MyCommentsTab";
import MyPostsTab from "./components/MyPostsTab";

function IncognitoPostsContent() {
  const { t } = useTranslation();
  useAuth(); // Check authentication and redirect if not authenticated
  const [currentTab, setCurrentTab] = useState(0);
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const handleTabChange = (_: React.SyntheticEvent, newValue: number) => {
    setCurrentTab(newValue);
  };

  const handlePostCreated = () => {
    setSuccess(t("incognitoPosts.success.postCreated"));
    // Refresh the browse tab if it's active
    if (currentTab === 0) {
      window.location.reload(); // Simple refresh for now
    }
  };

  const handleError = (errorMessage: string) => {
    setError(errorMessage);
  };

  return (
    <Container maxWidth="lg" sx={{ py: 3 }}>
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom align="center">
          {t("incognitoPosts.title")}
        </Typography>
        <Typography
          variant="body1"
          color="text.secondary"
          align="center"
          sx={{ mb: 3 }}
        >
          {t("incognitoPosts.description")}
        </Typography>

        {/* Create Post Button */}
        <Box sx={{ display: "flex", justifyContent: "center", mb: 3 }}>
          <Button
            variant="contained"
            color="primary"
            onClick={() => setCreateDialogOpen(true)}
            sx={{ minWidth: 200 }}
          >
            {t("incognitoPosts.browsing.createPost")}
          </Button>
        </Box>
      </Box>

      <Paper sx={{ width: "100%" }}>
        <Tabs
          value={currentTab}
          onChange={handleTabChange}
          indicatorColor="primary"
          textColor="primary"
          centered
        >
          <Tab label={t("incognitoPosts.browsing.title")} />
          <Tab label={t("incognitoPosts.browsing.myPosts")} />
          <Tab label={t("incognitoPosts.browsing.myComments")} />
        </Tabs>

        <Box sx={{ p: 3 }}>
          {currentTab === 0 && <BrowseTab onError={handleError} />}
          {currentTab === 1 && <MyPostsTab onError={handleError} />}
          {currentTab === 2 && <MyCommentsTab onError={handleError} />}
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
    </Container>
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
