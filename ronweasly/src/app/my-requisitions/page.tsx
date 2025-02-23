"use client";

import { useEffect } from "react";
import Box from "@mui/material/Box";
import Paper from "@mui/material/Paper";
import Typography from "@mui/material/Typography";
import CircularProgress from "@mui/material/CircularProgress";
import Alert from "@mui/material/Alert";
import Stack from "@mui/material/Stack";
import Link from "next/link";
import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import { useTranslation } from "@/hooks/useTranslation";
import { useColleagueSeeks } from "@/hooks/useColleagueSeeks";
import type { HubUserShort } from "@psankar/vetchi-typespec";
import { config } from "@/config";
import ProfilePicture from "@/components/ProfilePicture";

export default function MyRequisitionsPage() {
  const { t } = useTranslation();
  const { seeks, isLoading, error, fetchSeeks } = useColleagueSeeks();

  useEffect(() => {
    fetchSeeks();
  }, []);

  return (
    <AuthenticatedLayout>
      <Box sx={{ maxWidth: 800, mx: "auto", mt: 4, px: 2 }}>
        <Typography variant="h4" gutterBottom>
          {t("requisitions.title")}
        </Typography>

        <Paper sx={{ p: 3, mb: 4 }}>
          <Typography variant="h6" gutterBottom>
            {t("requisitions.colleagueSeeks")}
          </Typography>

          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error.message}
            </Alert>
          )}

          {isLoading ? (
            <Box sx={{ display: "flex", justifyContent: "center", p: 3 }}>
              <CircularProgress />
            </Box>
          ) : seeks?.seeks.length === 0 ? (
            <Typography
              color="text.secondary"
              sx={{ textAlign: "center", p: 3 }}
            >
              {t("requisitions.noSeeks")}
            </Typography>
          ) : (
            <Stack spacing={2}>
              {seeks?.seeks.map((user: HubUserShort) => (
                <Paper key={user.handle} variant="outlined" sx={{ p: 2 }}>
                  <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
                    <ProfilePicture
                      imageUrl={`${config.API_SERVER_PREFIX}/hub/profile-picture/${user.handle}`}
                      size={40}
                    />
                    <Box sx={{ flex: 1 }}>
                      <Link
                        href={`/u/${user.handle}`}
                        style={{ textDecoration: "none", color: "inherit" }}
                      >
                        <Typography variant="subtitle1" component="div">
                          {user.name}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          @{user.handle}
                        </Typography>
                      </Link>
                      <Typography variant="body2" sx={{ mt: 1 }}>
                        {user.short_bio}
                      </Typography>
                    </Box>
                  </Box>
                </Paper>
              ))}
            </Stack>
          )}
        </Paper>
      </Box>
    </AuthenticatedLayout>
  );
}
