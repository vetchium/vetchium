"use client";

import { useEffect } from "react";
import Box from "@mui/material/Box";
import Paper from "@mui/material/Paper";
import Typography from "@mui/material/Typography";
import CircularProgress from "@mui/material/CircularProgress";
import Alert from "@mui/material/Alert";
import Stack from "@mui/material/Stack";
import Button from "@mui/material/Button";
import Link from "next/link";
import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import { useTranslation } from "@/hooks/useTranslation";
import { useColleagueApprovals } from "@/hooks/useColleagueApprovals";
import { useColleagues } from "@/hooks/useColleagues";
import type { HubUserShort } from "@psankar/vetchi-typespec";
import { config } from "@/config";
import ProfilePicture from "@/components/ProfilePicture";

export default function MyApprovalsPage() {
  const { t } = useTranslation();
  const { approvals, isLoading, error, fetchApprovals } =
    useColleagueApprovals();
  const { approveColleague, rejectColleague, isApproving, isRejecting } =
    useColleagues();

  useEffect(() => {
    fetchApprovals();
  }, []);

  const handleApprove = async (handle: string) => {
    try {
      await approveColleague(handle);
      await fetchApprovals();
    } catch (err) {
      // Error handling is done in the hook
    }
  };

  const handleReject = async (handle: string) => {
    try {
      await rejectColleague(handle);
      await fetchApprovals();
    } catch (err) {
      // Error handling is done in the hook
    }
  };

  return (
    <AuthenticatedLayout>
      <Box sx={{ maxWidth: 800, mx: "auto", mt: 4, px: 2 }}>
        <Typography variant="h4" gutterBottom>
          {t("approvals.title")}
        </Typography>

        <Paper sx={{ p: 3, mb: 4 }}>
          <Typography variant="h6" gutterBottom>
            {t("approvals.colleagueApprovals")}
          </Typography>

          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error.message}
            </Alert>
          )}

          {isLoading ? (
            <Box sx={{ display: "flex", justifyContent: "center", p: 3 }}>
              <CircularProgress />
            </Box>
          ) : approvals?.approvals.length === 0 ? (
            <Typography
              color="text.secondary"
              sx={{ textAlign: "center", p: 3 }}
            >
              {t("approvals.noApprovals")}
            </Typography>
          ) : (
            <Stack spacing={2}>
              {approvals?.approvals.map((user: HubUserShort) => (
                <Paper key={user.handle} variant="outlined" sx={{ p: 2 }}>
                  <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
                    <ProfilePicture
                      imageUrl={`${config.API_SERVER_PREFIX}/hub/profile-picture/${user.handle}`}
                      size={40}
                    />
                    <Box sx={{ flex: 1 }}>
                      <Link
                        href={`/u/${user.handle}`}
                        style={{ textDecoration: "none", color: "inherit" }}
                      >
                        <Typography variant="subtitle1" component="div">
                          {user.name}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          @{user.handle}
                        </Typography>
                      </Link>
                      <Typography variant="body2" sx={{ mt: 1 }}>
                        {user.short_bio}
                      </Typography>
                    </Box>
                    <Stack direction="row" spacing={1}>
                      <Button
                        variant="contained"
                        color="primary"
                        onClick={() => handleApprove(user.handle)}
                        disabled={isApproving || isRejecting}
                      >
                        {t("common.approve")}
                      </Button>
                      <Button
                        variant="outlined"
                        color="error"
                        onClick={() => handleReject(user.handle)}
                        disabled={isApproving || isRejecting}
                      >
                        {t("common.reject")}
                      </Button>
                    </Stack>
                  </Box>
                </Paper>
              ))}
            </Stack>
          )}
        </Paper>
      </Box>
    </AuthenticatedLayout>
  );
}
