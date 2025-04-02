"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import { config } from "@/config";
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
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";

export default function SettingsPage() {
  const { t } = useTranslation();
  const router = useRouter();
  const [coolOffPeriod, setCoolOffPeriod] = useState<number | "">("");
  const [currentCoolOffPeriod, setCurrentCoolOffPeriod] = useState<
    number | null
  >(null);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [isUpdating, setIsUpdating] = useState(false);

  const fetchCoolOffPeriod = useCallback(async () => {
    setIsLoading(true);
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/get-cool-off-period`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/signin");
        return;
      }

      if (!response.ok) {
        throw new Error(t("settings.coolOffPeriod.fetchError"));
      }

      const periodText = await response.text();
      console.log("Response text:", periodText);

      const coolOffPeriod = parseInt(periodText, 10);
      if (isNaN(coolOffPeriod)) {
        console.error("Invalid cool off period in response:", periodText);
        throw new Error(t("settings.coolOffPeriod.fetchError"));
      }

      setCurrentCoolOffPeriod(coolOffPeriod);
      setCoolOffPeriod(coolOffPeriod);
    } catch {
      setError(t("settings.coolOffPeriod.fetchError"));
    } finally {
      setIsLoading(false);
    }
  }, [t, router]);

  useEffect(() => {
    fetchCoolOffPeriod();
  }, [fetchCoolOffPeriod]);

  const handleUpdateCoolOffPeriod = async () => {
    if (coolOffPeriod === "") return;

    setIsUpdating(true);
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const request: ChangeCoolOffPeriodRequest = {
        cool_off_period_days: Number(coolOffPeriod),
      };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/change-cool-off-period`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(request),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/signin");
        return;
      }

      if (!response.ok) {
        throw new Error(t("settings.coolOffPeriod.error"));
      }

      setSuccess(t("settings.coolOffPeriod.success"));
      setCurrentCoolOffPeriod(Number(coolOffPeriod));
    } catch {
      setError(t("settings.coolOffPeriod.error"));
    } finally {
      setIsUpdating(false);
    }
  };

  const content = (
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

  if (isLoading) {
    return (
      <AuthenticatedLayout>
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
      </AuthenticatedLayout>
    );
  }

  return <AuthenticatedLayout>{content}</AuthenticatedLayout>;
}
