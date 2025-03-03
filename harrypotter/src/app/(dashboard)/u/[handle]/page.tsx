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
  Avatar,
} from "@mui/material";
import Cookies from "js-cookie";

export default function UserProfilePage() {
  const params = useParams();
  const router = useRouter();
  const handle = params.handle as string;
  const [bio, setBio] = useState<EmployerViewBio | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [profilePicture, setProfilePicture] = useState<string | null>(null);
  const [profilePictureLoading, setProfilePictureLoading] = useState(false);
  const { t } = useTranslation();
  const sessionToken = Cookies.get("session_token");

  const fetchProfilePicture = useCallback(async () => {
    if (!sessionToken || !handle) return;

    try {
      setProfilePictureLoading(true);
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/get-hub-user-profile-picture/${handle}`,
        {
          headers: {
            Authorization: `Bearer ${sessionToken}`,
          },
        }
      );

      if (response.ok) {
        const blob = await response.blob();
        if (blob.size > 0) {
          const dataUrl = URL.createObjectURL(blob);
          setProfilePicture(dataUrl);
        } else {
          setProfilePicture(null);
        }
      } else {
        // The User may not have set a profile picture
        console.debug("Failed to fetch profile picture:", response.status);
        setProfilePicture(null);
      }
    } catch (error) {
      console.error("Error fetching profile picture:", error);
      setProfilePicture(null);
    } finally {
      setProfilePictureLoading(false);
    }
  }, [handle, sessionToken]);

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

  useEffect(() => {
    if (bio) {
      fetchProfilePicture();
    }
  }, [bio, fetchProfilePicture]);

  // Cleanup object URL on unmount
  useEffect(() => {
    return () => {
      if (profilePicture) {
        URL.revokeObjectURL(profilePicture);
      }
    };
  }, [profilePicture]);

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
    <Box
      sx={{
        p: { xs: 2, sm: 3 },
        maxWidth: "1200px",
        mx: "auto",
      }}
    >
      <Box
        sx={{
          display: "flex",
          flexDirection: { xs: "column", md: "row" },
          alignItems: { xs: "center", md: "flex-start" },
          gap: 4,
          mb: 4,
        }}
      >
        <Avatar
          src={profilePicture || ""}
          alt={bio.full_name}
          sx={{
            width: { xs: 200, md: 280 },
            height: { xs: 200, md: 280 },
            bgcolor: profilePictureLoading ? "grey.300" : "primary.main",
          }}
        >
          {profilePictureLoading ? (
            <CircularProgress size={100} />
          ) : (
            bio.full_name.charAt(0).toUpperCase()
          )}
        </Avatar>
        <Box
          sx={{
            textAlign: { xs: "center", md: "left" },
            flex: 1,
          }}
        >
          <Typography variant="h3" gutterBottom>
            {bio.full_name}
          </Typography>
          <Typography variant="h6" color="text.secondary" gutterBottom>
            @{bio.handle}
          </Typography>
          {bio.verified_mail_domains &&
            bio.verified_mail_domains.length > 0 && (
              <>
                <Typography
                  variant="subtitle1"
                  color="primary"
                  gutterBottom
                  sx={{ mt: 2 }}
                >
                  {t("hubUsers.verifiedOfficialEmailDomains")}
                </Typography>
                <Box
                  sx={{
                    display: "flex",
                    gap: 1,
                    flexWrap: "wrap",
                    justifyContent: { xs: "center", md: "flex-start" },
                  }}
                >
                  {bio.verified_mail_domains.map((domain, index) => (
                    <Chip
                      key={index}
                      label={domain}
                      color="primary"
                      variant="outlined"
                      size="small"
                    />
                  ))}
                </Box>
              </>
            )}
        </Box>
      </Box>

      <Paper sx={{ p: { xs: 2, sm: 3 }, mb: 3 }}>
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
    </Box>
  );
}
