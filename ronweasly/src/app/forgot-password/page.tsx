"use client";

import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Container from "@mui/material/Container";
import Paper from "@mui/material/Paper";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import {
  ForgotPasswordRequest,
  validateEmailAddress,
} from "@vetchium/typespec";
import Link from "next/link";
import { useState } from "react";

export default function ForgotPasswordPage() {
  const { t } = useTranslation();
  const [email, setEmail] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setSuccess(false);

    // Validate email
    if (!email.trim()) {
      setError(t("auth.forgotPassword.error.invalidEmail"));
      return;
    }

    if (!validateEmailAddress(email)) {
      setError(t("auth.forgotPassword.error.invalidEmail"));
      return;
    }

    setLoading(true);

    try {
      const request: ForgotPasswordRequest = { email: email.trim() };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/forgot-password`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(request),
        }
      );

      if (!response.ok) {
        switch (response.status) {
          case 404:
            throw new Error(t("auth.forgotPassword.error.userNotFound"));
          case 400:
            throw new Error(t("auth.forgotPassword.error.invalidEmail"));
          case 500:
          case 501:
          case 502:
          case 503:
          case 504:
            throw new Error(t("auth.errors.serverError"));
          default:
            throw new Error(t("auth.forgotPassword.error.sendFailed"));
        }
      }

      // Success
      setSuccess(true);
      setEmail("");
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : t("auth.forgotPassword.error.sendFailed")
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <Container component="main" maxWidth="xs">
      <Box
        sx={{
          marginTop: 8,
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
        }}
      >
        <Paper
          elevation={3}
          sx={{
            p: 4,
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            width: "100%",
          }}
        >
          {/* Logo */}
          <img
            src="/logo.webp"
            alt="Vetchium Logo"
            width={60}
            height={60}
            style={{ marginBottom: "16px" }}
          />
          <Typography component="h1" variant="h5">
            {t("auth.forgotPassword.title")}
          </Typography>
          <Typography
            variant="body2"
            color="text.secondary"
            sx={{ mt: 1, mb: 3, textAlign: "center" }}
          >
            {t("auth.forgotPassword.description")}
          </Typography>

          {success ? (
            <Box sx={{ width: "100%" }}>
              <Alert severity="success" sx={{ mb: 2 }}>
                {t("auth.forgotPassword.success")}
              </Alert>
              <Button
                component={Link}
                href="/login"
                fullWidth
                variant="contained"
                sx={{ mt: 2 }}
              >
                {t("auth.forgotPassword.backToLogin")}
              </Button>
            </Box>
          ) : (
            <Box
              component="form"
              onSubmit={handleSubmit}
              noValidate
              sx={{ mt: 1, width: "100%" }}
            >
              {error && (
                <Alert severity="error" sx={{ mb: 2 }}>
                  {error}
                </Alert>
              )}
              <TextField
                margin="normal"
                required
                fullWidth
                id="email"
                label={t("auth.forgotPassword.emailLabel")}
                name="email"
                autoComplete="email"
                autoFocus
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                disabled={loading}
              />
              <Button
                type="submit"
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
                disabled={loading}
              >
                {loading
                  ? t("common.loading")
                  : t("auth.forgotPassword.sendLinkButton")}
              </Button>
              <Box sx={{ textAlign: "center" }}>
                <Link href="/login" style={{ textDecoration: "none" }}>
                  <Typography variant="body2" color="primary">
                    {t("auth.forgotPassword.backToLogin")}
                  </Typography>
                </Link>
              </Box>
            </Box>
          )}
        </Paper>
      </Box>
    </Container>
  );
}
