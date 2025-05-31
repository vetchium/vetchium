"use client";

import { config } from "@/config";
import { useAuth } from "@/contexts/AuthContext";
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
import {
  EmployerSignInRequest,
  EmployerSignInResponse,
  GetOnboardStatusRequest,
  GetOnboardStatusResponse,
  OnboardStatuses,
} from "@vetchium/typespec";
import Cookies from "js-cookie";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";

export default function SignIn() {
  const [domain, setDomain] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [showCredentials, setShowCredentials] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const router = useRouter();
  const { t } = useTranslation();
  const { setUserEmail } = useAuth();

  const handleClickShowPassword = () => setShowPassword((show) => !show);

  const handleDomainVerification = async () => {
    try {
      const requestBody: GetOnboardStatusRequest = {
        client_id: domain,
      };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/get-onboard-status`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(requestBody),
        }
      );

      const data: GetOnboardStatusResponse = await response.json();

      if (data.status === OnboardStatuses.DOMAIN_NOT_VERIFIED) {
        setError(t("auth.domainNotVerifiedDetail"));
      } else if (
        data.status === OnboardStatuses.DOMAIN_VERIFIED_ONBOARD_PENDING
      ) {
        setError(t("auth.domainVerifyPendingDetail"));
      } else if (data.status === OnboardStatuses.DOMAIN_ONBOARDED) {
        setShowCredentials(true);
        setError("");
      }
    } catch {
      setError(t("auth.serverError"));
    }
  };

  const handleSignIn = async () => {
    try {
      const requestBody: EmployerSignInRequest = {
        client_id: domain,
        email,
        password,
      };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/signin`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(requestBody),
        }
      );

      if (response.status === 200) {
        const data: EmployerSignInResponse = await response.json();
        Cookies.set("tfa_token", data.token, { path: "/" });
        setUserEmail(email);
        router.push("/tfa");
      } else if (response.status === 422) {
        setError(t("auth.accountDisabled"));
      } else if (response.status >= 500) {
        setError(t("auth.serverError"));
      } else {
        setError(t("auth.invalidCredentials"));
      }
    } catch {
      setError(t("auth.serverError"));
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
          {t("auth.signin")}
        </Typography>
        <Box sx={{ mt: 1, width: "100%" }}>
          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
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
          />
          {!showCredentials && (
            <Alert severity="info" sx={{ mt: 1 }}>
              {t("auth.domainHelperText")}
            </Alert>
          )}
          {!showCredentials ? (
            <Button
              fullWidth
              variant="contained"
              sx={{ mt: 3, mb: 2 }}
              onClick={handleDomainVerification}
            >
              {t("auth.verify")}
            </Button>
          ) : (
            <>
              <TextField
                margin="normal"
                required
                fullWidth
                id="email"
                label={t("auth.email")}
                name="email"
                autoComplete="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
              <TextField
                margin="normal"
                required
                fullWidth
                name="password"
                label={t("auth.password")}
                type={showPassword ? "text" : "password"}
                id="password"
                autoComplete="current-password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
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
              <Button
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
                onClick={handleSignIn}
              >
                {t("auth.submit")}
              </Button>
              <Box sx={{ textAlign: "center" }}>
                <Link
                  href="/forgot-password"
                  style={{ textDecoration: "none" }}
                >
                  <Typography variant="body2" color="primary">
                    {t("auth.forgotPasswordLink")}
                  </Typography>
                </Link>
              </Box>
            </>
          )}
        </Box>
      </Box>
    </Container>
  );
}
