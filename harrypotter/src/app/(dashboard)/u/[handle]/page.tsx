"use client";

import { useParams, useRouter } from "next/navigation";
import { useEffect, useState, useCallback } from "react";
import type { EmployerViewBio } from "@psankar/vetchi-typespec";
import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import {
  Box,
  Typography,
  Alert,
  CircularProgress,
  Chip,
  Button,
  Paper,
  Divider,
} from "@mui/material";
import Cookies from "js-cookie";

export default function UserProfilePage() {
  const params = useParams();
  const router = useRouter();
  const handle = params.handle as string;
  const [bio, setBio] = useState<EmployerViewBio | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { t } = useTranslation();

  const fetchUserBio = useCallback(async () => {
    try {
      const sessionToken = Cookies.get("session_token");
      if (!sessionToken) {
        Cookies.remove("session_token");
        router.push("/signin");
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/get-hub-user-bio`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${sessionToken}`,
          },
          body: JSON.stringify({ handle }),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/signin");
        return;
      }

      if (!response.ok) {
        throw new Error(t("hubUsers.fetchError"));
      }

      const data = await response.json();
      setBio(data);
      setLoading(false);
    } catch (error) {
      console.error("Error fetching user bio:", error);
      setError(
        error instanceof Error ? error.message : t("errors.serverError")
      );
      setLoading(false);
    }
  }, [handle, router]);

  useEffect(() => {
    fetchUserBio();
  }, [fetchUserBio]);

  if (loading) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", my: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Box sx={{ p: 3 }}>
        <Box sx={{ display: "flex", justifyContent: "space-between", mb: 3 }}>
          <Typography variant="h4">{t("hubUsers.profile")}</Typography>
          <Button variant="outlined" onClick={() => router.back()}>
            {t("common.back")}
          </Button>
        </Box>
        <Alert severity="error">{error}</Alert>
      </Box>
    );
  }

  if (!bio) {
    return (
      <Box sx={{ p: 3 }}>
        <Box sx={{ display: "flex", justifyContent: "space-between", mb: 3 }}>
          <Typography variant="h4">{t("hubUsers.profile")}</Typography>
          <Button variant="outlined" onClick={() => router.back()}>
            {t("common.back")}
          </Button>
        </Box>
        <Alert severity="info">{t("hubUsers.notFound")}</Alert>
      </Box>
    );
  }

  return (
    <Box sx={{ p: 3 }}>
      <Box sx={{ display: "flex", justifyContent: "space-between", mb: 3 }}>
        <Box>
          <Typography variant="h4" gutterBottom>
            {bio.full_name}
          </Typography>
          <Typography variant="subtitle1" color="text.secondary">
            @{bio.handle}
          </Typography>
        </Box>
        <Button variant="outlined" onClick={() => router.back()}>
          {t("common.back")}
        </Button>
      </Box>

      <Paper sx={{ p: 3, mb: 3 }}>
        {bio.short_bio && (
          <>
            <Typography variant="h6" gutterBottom color="primary">
              {t("hubUsers.shortBio")}
            </Typography>
            <Typography paragraph>{bio.short_bio}</Typography>
          </>
        )}

        {bio.long_bio && (
          <>
            <Divider sx={{ my: 3 }} />
            <Typography variant="h6" gutterBottom color="primary">
              {t("hubUsers.longBio")}
            </Typography>
            <Typography>{bio.long_bio}</Typography>
          </>
        )}
      </Paper>

      {bio.verified_mail_domains && bio.verified_mail_domains.length > 0 && (
        <Paper sx={{ p: 3 }}>
          <Typography variant="h6" gutterBottom color="primary">
            {t("hubUsers.verifiedDomains")}
          </Typography>
          <Box sx={{ display: "flex", gap: 1, flexWrap: "wrap" }}>
            {bio.verified_mail_domains.map((domain: string, index: number) => (
              <Chip
                key={index}
                label={domain}
                color="primary"
                variant="outlined"
                size="small"
              />
            ))}
          </Box>
        </Paper>
      )}
    </Box>
  );
}
