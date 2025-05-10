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
import Link from "next/link";
import { useState } from "react";

export default function SignupRequestPage() {
  const { t } = useTranslation();
  const [email, setEmail] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setSuccess(false);

    try {
      const response = await fetch(`${config.API_SERVER_PREFIX}/hub/signup`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email }),
      });

      if (!response.ok) {
        switch (response.status) {
          case 460:
            throw new Error(t("signup.errors.domainNotSupported"));
          case 461:
            throw new Error(t("signup.errors.alreadyMemberOrInvited"));
          case 400:
            throw new Error(t("signup.errors.invalidEmail"));
          case 500:
          case 501:
          case 502:
          case 503:
          case 504:
            throw new Error(t("common.error.serverError"));
          default:
            throw new Error(t("signup.errors.signupFailed"));
        }
      }

      // Signup successful
      setSuccess(true);
      setEmail("");
    } catch (err) {
      setError(err instanceof Error ? err.message : t("signup.errors.signupFailed"));
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
            {t("signup.title")}
          </Typography>
          <Typography variant="body2" sx={{ mt: 1, mb: 3, textAlign: "center" }}>
            {t("signup.description")}
          </Typography>

          {success ? (
            <Box sx={{ width: "100%" }}>
              <Alert severity="success" sx={{ mb: 2 }}>
                {t("signup.success")}
              </Alert>
              <Button
                component={Link}
                href="/login"
                fullWidth
                variant="contained"
                sx={{ mt: 2 }}
              >
                {t("auth.backToLogin")}
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
                label={t("common.email")}
                name="email"
                autoComplete="email"
                autoFocus
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
              <Button
                type="submit"
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
              >
                {t("signup.submitButton")}
              </Button>
              <Box sx={{ textAlign: "center" }}>
                <Link href="/login" style={{ textDecoration: "none" }}>
                  <Typography variant="body2" color="primary">
                    {t("auth.backToLogin")}
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
