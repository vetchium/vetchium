"use client";

import {
  Box,
  Button,
  Container,
  TextField,
  Typography,
  Alert,
} from "@mui/material";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { useTranslation } from "@/hooks/useTranslation";
import { config } from "@/config";

export default function SignIn() {
  const [domain, setDomain] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [showCredentials, setShowCredentials] = useState(false);
  const router = useRouter();
  const { t } = useTranslation();

  const handleDomainVerification = async () => {
    try {
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/get-onboard-status`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ client_id: domain }),
        }
      );

      const data = await response.json();

      if (data.status === "DOMAIN_NOT_VERIFIED") {
        setError(t("auth.domainNotVerified"));
      } else if (data.status === "DOMAIN_VERIFIED_ONBOARD_PENDING") {
        setError(t("auth.domainVerifyPending"));
      } else if (data.status === "DOMAIN_ONBOARDED") {
        setShowCredentials(true);
        setError("");
      }
    } catch (err) {
      setError(t("auth.serverError"));
    }
  };

  const handleSignIn = async () => {
    try {
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/signin`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            client_id: domain,
            email,
            password,
          }),
        }
      );

      if (response.status === 200) {
        const data = await response.json();
        localStorage.setItem("tfaToken", data.token);
        router.push("/tfa");
      } else if (response.status === 422) {
        setError(t("auth.accountDisabled"));
      } else if (response.status >= 500) {
        setError(t("auth.serverError"));
      } else {
        setError(t("auth.invalidCredentials"));
      }
    } catch (err) {
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
            autoFocus
            value={domain}
            onChange={(e) => setDomain(e.target.value)}
          />
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
                type="password"
                id="password"
                autoComplete="current-password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
              <Button
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
                onClick={handleSignIn}
              >
                {t("auth.submit")}
              </Button>
            </>
          )}
        </Box>
      </Box>
    </Container>
  );
}
