"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import ChangeHandle from "@/components/ChangeHandle";
import UserInvite from "@/components/UserInvite";
import { useMyHandle } from "@/hooks/useMyHandle";
import { useMyTier } from "@/hooks/useMyTier";
import { useTranslation } from "@/hooks/useTranslation";
import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import CircularProgress from "@mui/material/CircularProgress";
import Container from "@mui/material/Container";
import Paper from "@mui/material/Paper";
import Typography from "@mui/material/Typography";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function Settings() {
  const { t } = useTranslation();
  const router = useRouter();

  // State for fetching initial page data
  const {
    myHandle,
    isLoading: isLoadingHandle,
    error: handleErorr,
  } = useMyHandle();
  const { tier, isLoading: isLoadingTier, error: tierError } = useMyTier();

  const isPageLoading = isLoadingHandle || isLoadingTier;
  const initialLoadingError = handleErorr || tierError;

  // Auth check
  useEffect(() => {
    const token = Cookies.get("session_token");
    if (!token) {
      router.push("/login");
    }
  }, [router]);

  return (
    <AuthenticatedLayout>
      <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          {t("settings.title")}
        </Typography>

        {/* General loading error */}
        {initialLoadingError && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {initialLoadingError.message}
          </Alert>
        )}

        {/* Change Handle Section */}
        <Paper sx={{ p: 3, mt: 3 }}>
          {
            isPageLoading ? (
              <Box sx={{ display: "flex", justifyContent: "center" }}>
                <CircularProgress />
              </Box>
            ) : myHandle && tier ? (
              <ChangeHandle currentHandle={myHandle} userTier={tier} />
            ) : !initialLoadingError ? ( // Only show if not already showing loading error
              <Typography color="text.secondary">
                {t("common.error.serverError")}{" "}
                {/* Fallback if handle/tier missing */}
              </Typography>
            ) : null /* Error handled above */
          }
        </Paper>

        {/* Invite User Section */}
        <UserInvite />
      </Container>
    </AuthenticatedLayout>
  );
}
