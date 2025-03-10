"use client";

import {
  Box,
  Button,
  Container,
  TextField,
  Typography,
  Alert,
  FormControlLabel,
  Switch,
} from "@mui/material";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useTranslation } from "@/hooks/useTranslation";
import { config } from "@/config";
import Cookies from "js-cookie";
import {
  EmployerTFARequest,
  EmployerTFAResponse,
} from "@psankar/vetchi-typespec";

export default function TFA() {
  const router = useRouter();
  const { t } = useTranslation();
  const [tfaCode, setTfaCode] = useState("");
  const [error, setError] = useState("");
  const [rememberMe, setRememberMe] = useState(false);

  // Initialize from localStorage on mount
  useEffect(() => {
    const saved = localStorage.getItem("rememberMe");
    if (saved !== null) {
      setRememberMe(JSON.parse(saved));
    }
  }, []);

  // Persist rememberMe state to localStorage
  useEffect(() => {
    localStorage.setItem("rememberMe", JSON.stringify(rememberMe));
  }, [rememberMe]);

  useEffect(() => {
    const token = Cookies.get("tfa_token");
    if (!token) {
      router.replace("/signin");
    }
  }, [router]);

  const handleVerify = async () => {
    try {
      const token = Cookies.get("tfa_token");
      const requestBody: EmployerTFARequest = {
        tfa_code: tfaCode,
        tfa_token: token!,
        remember_me: rememberMe,
      };

      const response = await fetch(`${config.API_SERVER_PREFIX}/employer/tfa`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(requestBody),
      });

      if (response.status === 200) {
        const data: EmployerTFAResponse = await response.json();
        Cookies.set("session_token", data.session_token, { path: "/" });
        Cookies.remove("tfa_token", { path: "/" });
        router.push("/");
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
          {t("auth.tfa")}
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
            id="tfaCode"
            label={t("auth.tfaCode")}
            name="tfaCode"
            autoFocus
            value={tfaCode}
            onChange={(e) => setTfaCode(e.target.value)}
          />
          <FormControlLabel
            control={
              <Switch
                checked={rememberMe}
                onChange={(e) => setRememberMe(e.target.checked)}
                name="rememberMe"
                color="primary"
              />
            }
            label={t("auth.rememberMe")}
          />
          <Button
            fullWidth
            variant="contained"
            sx={{ mt: 3, mb: 2 }}
            onClick={handleVerify}
          >
            {t("auth.verify")}
          </Button>
        </Box>
      </Box>
    </Container>
  );
}
