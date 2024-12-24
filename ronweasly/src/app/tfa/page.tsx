"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Container from "@mui/material/Container";
import Box from "@mui/material/Box";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import Typography from "@mui/material/Typography";
import Alert from "@mui/material/Alert";
import Paper from "@mui/material/Paper";
import { config } from "@/config";
import Cookies from "js-cookie";

export default function TFAPage() {
  const router = useRouter();
  const [tfaCode, setTfaCode] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    const tfaToken = Cookies.get("tfa_token");
    if (!tfaToken) {
      router.replace("/login");
      return;
    }

    try {
      const response = await fetch(`${config.API_SERVER_PREFIX}/hub/tfa`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          tfa_token: tfaToken,
          tfa_code: tfaCode,
          remember_me: true,
        }),
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.message || "TFA verification failed");
      }

      const data = await response.json();
      // Store the session token in a cookie and remove TFA token
      Cookies.set("session_token", data.session_token, { path: "/" });
      Cookies.remove("tfa_token", { path: "/" });
      router.push("/");
    } catch (err) {
      setError(err instanceof Error ? err.message : "TFA verification failed");
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
          <Typography component="h1" variant="h5" gutterBottom>
            Two-Factor Authentication
          </Typography>
          <Typography
            variant="body2"
            color="text.secondary"
            align="center"
            sx={{ mb: 3 }}
          >
            Please enter the verification code sent to your email
          </Typography>
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
              id="tfa-code"
              label="Verification Code"
              name="tfa-code"
              autoComplete="one-time-code"
              autoFocus
              value={tfaCode}
              onChange={(e) => setTfaCode(e.target.value)}
            />
            <Button
              type="submit"
              fullWidth
              variant="contained"
              sx={{ mt: 3, mb: 2 }}
            >
              Verify
            </Button>
          </Box>
        </Paper>
      </Box>
    </Container>
  );
}
