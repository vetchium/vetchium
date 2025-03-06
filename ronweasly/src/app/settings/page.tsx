"use client";

import { useState } from "react";
import { useTranslation } from "@/hooks/useTranslation";
import Container from "@mui/material/Container";
import Typography from "@mui/material/Typography";
import Paper from "@mui/material/Paper";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import Box from "@mui/material/Box";
import CircularProgress from "@mui/material/CircularProgress";
import Alert from "@mui/material/Alert";
import Snackbar from "@mui/material/Snackbar";
import { HubUserInviteRequest } from "@/types/hub/hubusers";
import { config } from "@/config";

export default function Settings() {
  const { t } = useTranslation();
  const [email, setEmail] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  const handleInviteUser = async () => {
    if (!email) {
      setError(t("settings.inviteUser.error.invalidEmail"));
      return;
    }

    setLoading(true);
    setError(null);
    setSuccess(false);

    try {
      const request: HubUserInviteRequest = {
        email,
      };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/invite-hub-user`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(request),
        }
      );

      if (!response.ok) {
        throw new Error(t("settings.inviteUser.error.failed"));
      }

      setSuccess(true);
      setEmail("");
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : t("settings.inviteUser.error.failed")
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Typography variant="h4" component="h1" gutterBottom>
        {t("settings.title")}
      </Typography>

      <Paper sx={{ p: 3, mt: 3 }}>
        <Typography variant="h6" gutterBottom>
          {t("settings.inviteUser.title")}
        </Typography>
        <Typography variant="body1" color="text.secondary" paragraph>
          {t("settings.inviteUser.description")}
        </Typography>

        <Box component="form" noValidate sx={{ mt: 2 }}>
          <TextField
            fullWidth
            label={t("common.email")}
            placeholder={t("settings.inviteUser.emailPlaceholder")}
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            disabled={loading}
            type="email"
            margin="normal"
          />

          <Box sx={{ mt: 2, display: "flex", justifyContent: "flex-end" }}>
            <Button
              variant="contained"
              onClick={handleInviteUser}
              disabled={loading}
              startIcon={loading ? <CircularProgress size={20} /> : undefined}
            >
              {t("settings.inviteUser.inviteButton")}
            </Button>
          </Box>
        </Box>
      </Paper>

      <Snackbar
        open={error !== null}
        autoHideDuration={6000}
        onClose={() => setError(null)}
      >
        <Alert severity="error" onClose={() => setError(null)}>
          {error}
        </Alert>
      </Snackbar>

      <Snackbar
        open={success}
        autoHideDuration={6000}
        onClose={() => setSuccess(false)}
      >
        <Alert severity="success" onClose={() => setSuccess(false)}>
          {t("settings.inviteUser.success")}
        </Alert>
      </Snackbar>
    </Container>
  );
}
