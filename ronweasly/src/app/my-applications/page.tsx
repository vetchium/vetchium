"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import { config } from "@/config";
import { useAuth } from "@/hooks/useAuth";
import { useTranslation } from "@/hooks/useTranslation";
import OpenInNewIcon from "@mui/icons-material/OpenInNew";
import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Chip from "@mui/material/Chip";
import CircularProgress from "@mui/material/CircularProgress";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";
import Paper from "@mui/material/Paper";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import { ApplicationState, HubApplication } from "@vetchium/typespec";
import Cookies from "js-cookie";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

const formatApplicationState = (
  state: ApplicationState,
  t: (key: string) => string
) => {
  switch (state) {
    case "APPLIED":
      return t("myApplications.applicationState.applied");
    case "REJECTED":
      return t("myApplications.applicationState.rejected");
    case "SHORTLISTED":
      return t("myApplications.applicationState.shortlisted");
    case "WITHDRAWN":
      return t("myApplications.applicationState.withdrawn");
    case "EXPIRED":
      return t("myApplications.applicationState.expired");
    default:
      return state;
  }
};

const getChipColor = (state: ApplicationState) => {
  switch (state) {
    case "APPLIED":
      return "primary";
    case "REJECTED":
      return "error";
    case "SHORTLISTED":
      return "success";
    case "WITHDRAWN":
      return "default";
    case "EXPIRED":
      return "warning";
    default:
      return "default";
  }
};

export default function MyApplicationsPage() {
  const { t } = useTranslation();
  const router = useRouter();
  useAuth(); // Check authentication and redirect if not authenticated
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [applications, setApplications] = useState<HubApplication[]>([]);
  const [withdrawDialogOpen, setWithdrawDialogOpen] = useState(false);
  const [selectedApplication, setSelectedApplication] = useState<string | null>(
    null
  );
  const [withdrawing, setWithdrawing] = useState(false);

  useEffect(() => {
    const fetchApplications = async () => {
      const token = Cookies.get("session_token");
      if (!token) {
        setError(t("common.error.notAuthenticated"));
        return;
      }

      try {
        const response = await fetch(
          `${config.API_SERVER_PREFIX}/hub/my-applications`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify({
              limit: 40,
            }),
          }
        );

        if (!response.ok) {
          if (response.status === 401) {
            setError(t("common.error.sessionExpired"));
            Cookies.remove("session_token", { path: "/" });
            return;
          }
          throw new Error(
            `Failed to fetch applications: ${response.statusText}`
          );
        }

        const data = await response.json();
        setApplications(data ?? []);
      } catch (error) {
        console.error("Error fetching applications:", error);
        setError(t("myApplications.error.loadFailed"));
      } finally {
        setLoading(false);
      }
    };

    fetchApplications();
  }, []);

  const handleWithdraw = async () => {
    if (!selectedApplication) return;

    const token = Cookies.get("session_token");
    if (!token) {
      setError(t("common.error.notAuthenticated"));
      return;
    }

    setWithdrawing(true);
    try {
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/withdraw-application`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            application_id: selectedApplication,
          }),
        }
      );

      if (!response.ok) {
        if (response.status === 401) {
          setError(t("common.error.sessionExpired"));
          Cookies.remove("session_token", { path: "/" });
          return;
        }
        throw new Error(t("myApplications.withdrawError"));
      }

      // Update the application state locally
      setApplications((prevApplications) =>
        prevApplications.map((app) =>
          app.application_id === selectedApplication
            ? { ...app, state: "WITHDRAWN" as ApplicationState }
            : app
        )
      );

      setWithdrawDialogOpen(false);
    } catch (error) {
      console.error("Error withdrawing application:", error);
      setError(t("myApplications.withdrawError"));
    } finally {
      setWithdrawing(false);
      setSelectedApplication(null);
    }
  };

  if (loading) {
    return (
      <AuthenticatedLayout>
        <Box sx={{ display: "flex", justifyContent: "center", mt: 4 }}>
          <CircularProgress />
        </Box>
      </AuthenticatedLayout>
    );
  }

  return (
    <AuthenticatedLayout>
      <Box sx={{ maxWidth: 800, mx: "auto", mt: 4, px: 2 }}>
        <Typography variant="h4" gutterBottom>
          {t("myApplications.title")}
        </Typography>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        {(applications ?? []).length === 0 ? (
          <Paper sx={{ p: 4, textAlign: "center" }}>
            <Typography color="text.secondary">
              {t("myApplications.noApplications")}
            </Typography>
          </Paper>
        ) : (
          <Stack spacing={2}>
            {applications.map((application) => (
              <Paper key={application.application_id} sx={{ p: 3 }}>
                <Box
                  sx={{
                    display: "flex",
                    justifyContent: "space-between",
                    alignItems: "flex-start",
                    mb: 2,
                  }}
                >
                  <Box>
                    <Typography variant="h6" gutterBottom>
                      {application.opening_title}
                    </Typography>
                    <Typography
                      variant="subtitle1"
                      color="text.secondary"
                      gutterBottom
                    >
                      {application.employer_name}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {t("myApplications.appliedOn", {
                        date: new Date(application.created_at)
                          .toLocaleString("en-US", {
                            year: "numeric",
                            month: "short",
                            day: "2-digit",
                            hour: "numeric",
                            minute: "2-digit",
                            hour12: true,
                          })
                          .replace(/,/g, ""),
                      })}
                    </Typography>
                  </Box>
                  <Chip
                    label={formatApplicationState(application.state, t)}
                    color={getChipColor(application.state) as any}
                  />
                </Box>

                <Box sx={{ display: "flex", gap: 1 }}>
                  <Link
                    href={`/org/${application.employer_domain}/opening/${application.opening_id}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    style={{ textDecoration: "none" }}
                  >
                    <Button variant="outlined" endIcon={<OpenInNewIcon />}>
                      {t("myApplications.viewOpening")}
                    </Button>
                  </Link>
                  {application.state === "APPLIED" && (
                    <Button
                      variant="outlined"
                      color="error"
                      onClick={() => {
                        setSelectedApplication(application.application_id);
                        setWithdrawDialogOpen(true);
                      }}
                    >
                      {t("myApplications.withdrawApplication")}
                    </Button>
                  )}
                </Box>
              </Paper>
            ))}
          </Stack>
        )}

        <Dialog
          open={withdrawDialogOpen}
          onClose={() => setWithdrawDialogOpen(false)}
        >
          <DialogTitle>{t("myApplications.withdrawApplication")}</DialogTitle>
          <DialogContent>
            <DialogContentText>
              {t("myApplications.withdrawConfirmation")}
            </DialogContentText>
          </DialogContent>
          <DialogActions>
            <Button
              onClick={() => setWithdrawDialogOpen(false)}
              disabled={withdrawing}
            >
              Cancel
            </Button>
            <Button
              onClick={handleWithdraw}
              color="error"
              disabled={withdrawing}
              autoFocus
            >
              {withdrawing ? (
                <CircularProgress size={24} color="inherit" />
              ) : (
                t("myApplications.withdrawApplication")
              )}
            </Button>
          </DialogActions>
        </Dialog>
      </Box>
    </AuthenticatedLayout>
  );
}
