"use client";

import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import CircularProgress from "@mui/material/CircularProgress";
import Paper from "@mui/material/Paper";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import { HubUserInviteRequest } from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation"; // Needed for redirect on 401
import { useState } from "react";

export default function UserInvite() {
  const { t } = useTranslation();
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const [emailError, setEmailError] = useState<string | null>(null);

  const validateEmail = (email: string): boolean => {
    const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    return emailRegex.test(email);
  };

  const handleEmailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newEmail = e.target.value;
    setEmail(newEmail);
    if (newEmail && !validateEmail(newEmail)) {
      setEmailError(t("settings.inviteUser.error.invalidEmail"));
    } else {
      setEmailError(null);
      setError(null); // Clear general error when fixing email
    }
    setSuccess(false); // Clear success message on new input
  };

  const handleInviteUser = async () => {
    if (!email || !validateEmail(email)) {
      setEmailError(t("settings.inviteUser.error.invalidEmail"));
      // setError(t("settings.inviteUser.error.invalidEmail")); // Set specific email error instead
      return;
    }

    const token = Cookies.get("session_token");
    if (!token) {
      // Handle case where component is rendered but user logs out somehow
      router.push("/login");
      return;
    }

    setLoading(true);
    setError(null);
    setSuccess(false);
    setEmailError(null); // Clear specific email error

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
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(request),
        }
      );

      if (!response.ok) {
        if (response.status === 401) {
          Cookies.remove("session_token");
          router.push("/login");
          return;
        }
        // Consider parsing error response body if available
        // const errorData = await response.json().catch(() => null);
        throw new Error(t("settings.inviteUser.error.failed"));
      }

      setSuccess(true);
      setEmail(""); // Clear email field on success
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
    <Paper sx={{ p: 3, mt: 3 }}>
      <Typography variant="h6" gutterBottom>
        {t("settings.inviteUser.title")}
      </Typography>
      <Typography variant="body1" color="text.secondary" paragraph>
        {t("settings.inviteUser.description")}
      </Typography>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}
      {success && (
        <Alert severity="success" sx={{ mb: 2 }}>
          {t("settings.inviteUser.success")}
        </Alert>
      )}

      <Box
        component="form"
        noValidate
        sx={{ mt: 2 }}
        onSubmit={(e) => {
          e.preventDefault();
          handleInviteUser();
        }}
      >
        <TextField
          fullWidth
          label={t("common.email")}
          placeholder={t("settings.inviteUser.emailPlaceholder")}
          value={email}
          onChange={handleEmailChange}
          disabled={loading}
          type="email"
          required
          error={!!emailError}
          helperText={emailError}
          inputProps={{
            pattern: "[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
          }}
          margin="normal"
        />

        <Box sx={{ mt: 2, display: "flex", justifyContent: "flex-end" }}>
          <Button
            variant="contained"
            type="submit"
            disabled={loading || !!emailError}
            startIcon={loading ? <CircularProgress size={20} /> : undefined}
          >
            {t("settings.inviteUser.inviteButton")}
          </Button>
        </Box>
      </Box>
    </Paper>
  );
}
