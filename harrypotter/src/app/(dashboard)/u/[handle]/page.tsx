"use client";

import { AchievementSection } from "@/components/Achievement";
import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import {
  Alert,
  Avatar,
  Box,
  Button,
  Chip,
  CircularProgress,
  Divider,
  Paper,
  Stack,
  Typography,
} from "@mui/material";
import type { Education, EmployerViewBio } from "@vetchium/typespec";
import { AchievementType } from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useParams, useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";

export default function UserProfilePage() {
  const params = useParams();
  const router = useRouter();
  const handle = params.handle as string;
  const [bio, setBio] = useState<EmployerViewBio | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [profilePicture, setProfilePicture] = useState<string | null>(null);
  const [profilePictureLoading, setProfilePictureLoading] = useState(false);
  const [education, setEducation] = useState<Education[]>([]);
  const [educationLoading, setEducationLoading] = useState(false);
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
  }, [handle, router, t]);

  const fetchEducation = useCallback(async () => {
    if (!sessionToken || !handle) return;

    try {
      setEducationLoading(true);
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/list-hub-user-education`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${sessionToken}`,
          },
          body: JSON.stringify({ handle }),
        }
      );

      if (response.ok) {
        const data = await response.json();
        setEducation(data);
      } else {
        console.debug("Failed to fetch education:", response.status);
        setEducation([]);
      }
    } catch (error) {
      console.error("Error fetching education:", error);
      setEducation([]);
    } finally {
      setEducationLoading(false);
    }
  }, [handle, sessionToken]);

  useEffect(() => {
    fetchUserBio();
  }, [fetchUserBio]);

  useEffect(() => {
    if (bio) {
      fetchProfilePicture();
      fetchEducation();
    }
  }, [bio, fetchProfilePicture, fetchEducation]);

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

      <Stack spacing={3}>
        <Paper sx={{ p: { xs: 2, sm: 3 } }}>
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

        <Paper sx={{ p: { xs: 2, sm: 3 } }}>
          <Typography variant="h6" gutterBottom color="primary">
            {t("hubUsers.workHistory")}
          </Typography>

          {bio.work_history && bio.work_history.length > 0 ? (
            <Box>
              {bio.work_history.map((work, index) => (
                <Box
                  key={work.id}
                  sx={{
                    mb: 3,
                    "&:last-child": { mb: 0 },
                  }}
                >
                  <Typography variant="h6" component="div">
                    {work.title}
                  </Typography>
                  <Typography color="text.secondary">
                    {work.employer_name
                      ? `${work.employer_name} (${work.employer_domain})`
                      : work.employer_domain}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    {new Date(work.start_date).toLocaleDateString(undefined, {
                      month: "long",
                      year: "numeric",
                    })}{" "}
                    -{" "}
                    {work.end_date
                      ? new Date(work.end_date).toLocaleDateString(undefined, {
                          month: "long",
                          year: "numeric",
                        })
                      : t("hubUsers.currentlyWorking")}
                  </Typography>
                  {work.description && (
                    <Typography
                      variant="body2"
                      sx={{
                        mt: 1,
                        whiteSpace: "pre-wrap",
                      }}
                    >
                      {work.description}
                    </Typography>
                  )}
                  {index < bio.work_history.length - 1 && (
                    <Divider sx={{ mt: 3 }} />
                  )}
                </Box>
              ))}
            </Box>
          ) : (
            <Typography color="text.secondary">
              {t("hubUsers.noWorkHistory")}
            </Typography>
          )}
        </Paper>

        <Paper sx={{ p: { xs: 2, sm: 3 } }}>
          <Typography variant="h6" gutterBottom color="primary">
            {t("hubUsers.education")}
          </Typography>

          {educationLoading ? (
            <Box sx={{ display: "flex", justifyContent: "center", my: 2 }}>
              <CircularProgress size={24} />
            </Box>
          ) : education && education.length > 0 ? (
            <Box>
              {education.map((edu, index) => (
                <Box
                  key={edu.id || index}
                  sx={{
                    mb: 3,
                    "&:last-child": { mb: 0 },
                  }}
                >
                  <Typography variant="h6" component="div">
                    {edu.degree || t("hubUsers.educationDegree")}
                  </Typography>
                  <Typography color="text.secondary">
                    {edu.institute_domain}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    {edu.start_date
                      ? new Date(edu.start_date).toLocaleDateString(undefined, {
                          month: "long",
                          year: "numeric",
                        })
                      : ""}{" "}
                    {edu.start_date && "-"}{" "}
                    {edu.end_date
                      ? new Date(edu.end_date).toLocaleDateString(undefined, {
                          month: "long",
                          year: "numeric",
                        })
                      : t("hubUsers.currentlyStudying")}
                  </Typography>
                  {edu.description && (
                    <Typography
                      variant="body2"
                      sx={{
                        mt: 1,
                        whiteSpace: "pre-wrap",
                      }}
                    >
                      {edu.description}
                    </Typography>
                  )}
                  {index < education.length - 1 && <Divider sx={{ mt: 3 }} />}
                </Box>
              ))}
            </Box>
          ) : (
            <Typography color="text.secondary">
              {t("hubUsers.noEducation")}
            </Typography>
          )}
        </Paper>

        <Paper sx={{ p: { xs: 2, sm: 3 } }}>
          <Typography variant="h6" gutterBottom color="primary">
            {t("achievements.patents.title")}
          </Typography>
          <AchievementSection
            userHandle={handle}
            achievementType={AchievementType.PATENT}
          />
        </Paper>

        <Paper sx={{ p: { xs: 2, sm: 3 } }}>
          <Typography variant="h6" gutterBottom color="primary">
            {t("achievements.publications.title")}
          </Typography>
          <AchievementSection
            userHandle={handle}
            achievementType={AchievementType.PUBLICATION}
          />
        </Paper>

        <Paper sx={{ p: { xs: 2, sm: 3 } }}>
          <Typography variant="h6" gutterBottom color="primary">
            {t("achievements.certifications.title")}
          </Typography>
          <AchievementSection
            userHandle={handle}
            achievementType={AchievementType.CERTIFICATION}
          />
        </Paper>
      </Stack>
    </Box>
  );
}
