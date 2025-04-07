"use client";

import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import { Visibility, VisibilityOff } from "@mui/icons-material";
import {
  Alert,
  Box,
  Button,
  Container,
  IconButton,
  TextField,
  Typography,
} from "@mui/material";
import { SetOnboardPasswordRequest } from "@vetchium/typespec";
import { useParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";

export default function SignupOrgUser() {
  const params = useParams();
  const token = params.token as string;
  const [domain, setDomain] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const [redirectCounter, setRedirectCounter] = useState(3);
  const [isLoading, setIsLoading] = useState(false);
  const { t } = useTranslation();
  const router = useRouter();

  const handleClickShowPassword = () => setShowPassword((show) => !show);
  const handleClickShowConfirmPassword = () =>
    setShowConfirmPassword((show) => !show);

  const validateInputs = () => {
    if (!domain.trim()) {
      setError(t("validation.name.required"));
      return false;
    }

    if (password !== confirmPassword) {
      setError(t("auth.passwordsDontMatch"));
      return false;
    }

    return true;
  };

  useEffect(() => {
    // Redirect to sign in page after showing success message
    if (success) {
      const timer = setTimeout(() => {
        router.push("/signin");
      }, 3000);

      return () => clearTimeout(timer);
    }
  }, [success, router]);

  useEffect(() => {
    // Count down timer for redirect
    if (success && redirectCounter > 0) {
      const interval = setInterval(() => {
        setRedirectCounter((prev) => prev - 1);
      }, 1000);

      return () => clearInterval(interval);
    }
  }, [success, redirectCounter]);

  const handleSubmit = async () => {
    if (!validateInputs()) {
      return;
    }

    setIsLoading(true);
    setError("");
    setSuccess(false);

    try {
      const requestBody: SetOnboardPasswordRequest = {
        client_id: domain,
        password,
        token,
      };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/set-onboard-password`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(requestBody),
        }
      );

      if (response.status === 200) {
        // Success - show success message then redirect
        setSuccess(true);
        setRedirectCounter(3);
      } else if (response.status === 422) {
        // Expired or invalid token
        setError(t("auth.expiredOrInvalidToken"));
      } else {
        // Other errors
        setError(t("common.serverError"));
      }
    } catch (err) {
      console.error(err);
      setError(t("common.serverError"));
    } finally {
      setIsLoading(false);
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
        <Typography component="h1" variant="h5">
          {t("auth.signupOrgUser")}
        </Typography>
        <Box sx={{ mt: 1, width: "100%" }}>
          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}
          {success && (
            <Alert severity="success" sx={{ mb: 2 }}>
              {t("auth.passwordSetupSuccess")}{" "}
              {t("common.redirecting", { seconds: redirectCounter })}
            </Alert>
          )}
          <TextField
            margin="normal"
            required
            fullWidth
            id="domain"
            label={t("auth.domain")}
            name="domain"
            autoComplete="off"
            autoFocus
            value={domain}
            onChange={(e) => setDomain(e.target.value)}
            disabled={success || isLoading}
          />
          <TextField
            margin="normal"
            required
            fullWidth
            name="password"
            label={t("auth.password")}
            type={showPassword ? "text" : "password"}
            id="password"
            autoComplete="new-password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            disabled={success || isLoading}
            InputProps={{
              endAdornment: (
                <IconButton
                  aria-label="toggle password visibility"
                  onClick={handleClickShowPassword}
                  edge="end"
                  disabled={success || isLoading}
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
            label={t("auth.confirmPassword")}
            type={showConfirmPassword ? "text" : "password"}
            id="confirmPassword"
            autoComplete="new-password"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            disabled={success || isLoading}
            InputProps={{
              endAdornment: (
                <IconButton
                  aria-label="toggle confirm password visibility"
                  onClick={handleClickShowConfirmPassword}
                  edge="end"
                  disabled={success || isLoading}
                >
                  {showConfirmPassword ? <VisibilityOff /> : <Visibility />}
                </IconButton>
              ),
            }}
          />
          <Button
            fullWidth
            variant="contained"
            sx={{ mt: 3, mb: 2 }}
            onClick={handleSubmit}
            disabled={isLoading || success}
          >
            {t("auth.completeSignup")}
          </Button>
        </Box>
      </Box>
    </Container>
  );
}
