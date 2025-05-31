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
import { validateEmailAddress } from "@vetchium/typespec";
import Link from "next/link";
import { useState } from "react";

interface ForgotPasswordRequest {
  email: string;
}

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
        `${config.API_SERVER_PREFIX}/employer/forgot-password`,
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
            throw new Error(t("common.serverError"));
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
          <Box
            sx={{
              width: 60,
              height: 60,
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              backgroundColor: "primary.main",
              borderRadius: "50%",
              mb: 2,
            }}
          >
            <Typography variant="h4" color="white" fontWeight="bold">
              V
            </Typography>
          </Box>
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
                href="/signin"
                fullWidth
                variant="contained"
                sx={{ mt: 2 }}
              >
                {t("auth.forgotPassword.backToSignin")}
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
                <Link href="/signin" style={{ textDecoration: "none" }}>
                  <Typography variant="body2" color="primary">
                    {t("auth.forgotPassword.backToSignin")}
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
