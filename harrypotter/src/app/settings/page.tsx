"use client";

import { useTranslation } from "@/hooks/useTranslation";
import {
  Alert,
  Box,
  Button,
  Card,
  CardContent,
  CircularProgress,
  Container,
  Snackbar,
  TextField,
  Typography,
} from "@mui/material";
import { ChangeCoolOffPeriodRequest } from "@psankar/vetchi-typespec";
import { useEffect, useState } from "react";

export default function SettingsPage() {
  const { t } = useTranslation();
  const [coolOffPeriod, setCoolOffPeriod] = useState<number | "">("");
  const [currentCoolOffPeriod, setCurrentCoolOffPeriod] = useState<
    number | null
  >(null);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [isUpdating, setIsUpdating] = useState(false);

  useEffect(() => {
    fetchCoolOffPeriod();
  }, []);

  const fetchCoolOffPeriod = async () => {
    setIsLoading(true);
    try {
      const response = await fetch("/api/employer/get-cool-off-period");
      if (!response.ok) {
        throw new Error(t("settings.coolOffPeriod.fetchError"));
      }
      const data = await response.json();
      setCurrentCoolOffPeriod(data.coolOffPeriod);
      setCoolOffPeriod(data.coolOffPeriod);
    } catch (err) {
      setError(t("settings.coolOffPeriod.fetchError"));
    } finally {
      setIsLoading(false);
    }
  };

  const handleUpdateCoolOffPeriod = async () => {
    if (coolOffPeriod === "") return;

    setIsUpdating(true);
    try {
      const request: ChangeCoolOffPeriodRequest = {
        coolOffPeriod: Number(coolOffPeriod),
      };

      const response = await fetch("/api/employer/change-cool-off-period", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(request),
      });

      if (!response.ok) {
        throw new Error(t("settings.coolOffPeriod.error"));
      }

      setSuccess(t("settings.coolOffPeriod.success"));
      setCurrentCoolOffPeriod(Number(coolOffPeriod));
    } catch (err) {
      setError(t("settings.coolOffPeriod.error"));
    } finally {
      setIsUpdating(false);
    }
  };

  if (isLoading) {
    return (
      <Container maxWidth="lg">
        <Box
          sx={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            height: "50vh",
          }}
        >
          <CircularProgress />
        </Box>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg">
      <Box sx={{ p: 3 }}>
        <Typography variant="h4" gutterBottom>
          {t("settings.title")}
        </Typography>

        <Card sx={{ mt: 3 }}>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              {t("settings.coolOffPeriod.title")}
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
              {t("settings.coolOffPeriod.description")}
            </Typography>

            {currentCoolOffPeriod !== null && (
              <Typography variant="body1" sx={{ mb: 2 }}>
                {t("settings.coolOffPeriod.current", {
                  days: currentCoolOffPeriod.toString(),
                })}
              </Typography>
            )}

            <Box sx={{ display: "flex", gap: 2, alignItems: "center" }}>
              <TextField
                type="number"
                value={coolOffPeriod}
                onChange={(e) =>
                  setCoolOffPeriod(
                    e.target.value === "" ? "" : Number(e.target.value)
                  )
                }
                inputProps={{ min: 0, max: 365 }}
                sx={{ width: 200 }}
                disabled={isUpdating}
              />
              <Button
                variant="contained"
                onClick={handleUpdateCoolOffPeriod}
                disabled={
                  coolOffPeriod === "" ||
                  coolOffPeriod === currentCoolOffPeriod ||
                  isUpdating
                }
                startIcon={
                  isUpdating ? (
                    <CircularProgress size={20} color="inherit" />
                  ) : null
                }
              >
                {t("settings.coolOffPeriod.update")}
              </Button>
            </Box>
          </CardContent>
        </Card>

        <Snackbar
          open={error !== null}
          autoHideDuration={6000}
          onClose={() => setError(null)}
        >
          <Alert severity="error" onClose={() => setError(null)}>
            {error}
          </Alert>
        </Snackbar>

        <Snackbar
          open={success !== null}
          autoHideDuration={6000}
          onClose={() => setSuccess(null)}
        >
          <Alert severity="success" onClose={() => setSuccess(null)}>
            {success}
          </Alert>
        </Snackbar>
      </Box>
    </Container>
  );
}
