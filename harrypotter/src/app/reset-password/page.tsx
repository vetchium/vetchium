"use client";

import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import Visibility from "@mui/icons-material/Visibility";
import VisibilityOff from "@mui/icons-material/VisibilityOff";
import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Container from "@mui/material/Container";
import FormHelperText from "@mui/material/FormHelperText";
import IconButton from "@mui/material/IconButton";
import Paper from "@mui/material/Paper";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import { validatePassword } from "@vetchium/typespec";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { Suspense, useEffect, useState } from "react";

interface ResetPasswordRequest {
  token: string;
  password: string;
}

function ResetPasswordContent() {
  const { t } = useTranslation();
  const searchParams = useSearchParams();
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const [loading, setLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [token, setToken] = useState<string | null>(null);

  useEffect(() => {
    const tokenParam = searchParams?.get("token");
    if (!tokenParam) {
      setError(t("auth.resetPassword.error.invalidToken"));
    } else {
      setToken(tokenParam);
    }
  }, [searchParams, t]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setSuccess(false);

    if (!token) {
      setError(t("auth.resetPassword.error.invalidToken"));
      return;
    }

    // Validate password
    if (!password.trim()) {
      setError(t("auth.resetPassword.error.passwordRequired"));
      return;
    }

    const passwordValidation = validatePassword(password);
    if (!passwordValidation.isValid) {
      setError(
        passwordValidation.error || t("auth.resetPassword.error.weakPassword")
      );
      return;
    }

    // Check if passwords match
    if (password !== confirmPassword) {
      setError(t("auth.resetPassword.error.passwordsDoNotMatch"));
      return;
    }

    setLoading(true);

    try {
      const request: ResetPasswordRequest = {
        token: token,
        password: password,
      };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/reset-password`,
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
          case 400:
          case 401:
          case 404:
            throw new Error(t("auth.resetPassword.error.invalidToken"));
          case 500:
          case 501:
          case 502:
          case 503:
          case 504:
            throw new Error(t("common.serverError"));
          default:
            throw new Error(t("auth.resetPassword.error.resetFailed"));
        }
      }

      // Success
      setSuccess(true);
      setPassword("");
      setConfirmPassword("");
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : t("auth.resetPassword.error.resetFailed")
      );
    } finally {
      setLoading(false);
    }
  };

  const handleClickShowPassword = () => setShowPassword((show) => !show);
  const handleClickShowConfirmPassword = () =>
    setShowConfirmPassword((show) => !show);

  if (!token && !error) {
    return (
      <Container component="main" maxWidth="xs">
        <Box sx={{ marginTop: 8, display: "flex", justifyContent: "center" }}>
          <Typography>{t("common.loading")}</Typography>
        </Box>
      </Container>
    );
  }

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
            {t("auth.resetPassword.title")}
          </Typography>
          <Typography
            variant="body2"
            color="text.secondary"
            sx={{ mt: 1, mb: 3, textAlign: "center" }}
          >
            {t("auth.resetPassword.description")}
          </Typography>

          {success ? (
            <Box sx={{ width: "100%" }}>
              <Alert severity="success" sx={{ mb: 2 }}>
                {t("auth.resetPassword.success")}
              </Alert>
              <Button
                component={Link}
                href="/signin"
                fullWidth
                variant="contained"
                sx={{ mt: 2 }}
              >
                {t("auth.resetPassword.backToSignin")}
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
                name="password"
                label={t("auth.resetPassword.passwordLabel")}
                type={showPassword ? "text" : "password"}
                id="password"
                autoComplete="new-password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                disabled={loading}
                InputProps={{
                  endAdornment: (
                    <IconButton
                      aria-label="toggle password visibility"
                      onClick={handleClickShowPassword}
                      edge="end"
                    >
                      {showPassword ? <VisibilityOff /> : <Visibility />}
                    </IconButton>
                  ),
                }}
              />

              <TextField
                margin="normal"
                required
                fullWidth
                name="confirmPassword"
                label={t("auth.resetPassword.confirmPasswordLabel")}
                type={showConfirmPassword ? "text" : "password"}
                id="confirmPassword"
                autoComplete="new-password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                disabled={loading}
                InputProps={{
                  endAdornment: (
                    <IconButton
                      aria-label="toggle confirm password visibility"
                      onClick={handleClickShowConfirmPassword}
                      edge="end"
                    >
                      {showConfirmPassword ? <VisibilityOff /> : <Visibility />}
                    </IconButton>
                  ),
                }}
              />

              <FormHelperText sx={{ mt: 1, mb: 2 }}>
                {t("auth.resetPassword.passwordHelp")}
              </FormHelperText>

              <Button
                type="submit"
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
                disabled={loading || !token}
              >
                {loading
                  ? t("common.loading")
                  : t("auth.resetPassword.resetButton")}
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

export default function ResetPasswordPage() {
  return (
    <Suspense
      fallback={
        <Container component="main" maxWidth="xs">
          <Box sx={{ marginTop: 8, display: "flex", justifyContent: "center" }}>
            <Typography>Loading...</Typography>
          </Box>
        </Container>
      }
    >
      <ResetPasswordContent />
    </Suspense>
  );
}
