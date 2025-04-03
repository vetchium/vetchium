"use client";

import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import AddIcon from "@mui/icons-material/Add";
import DeleteIcon from "@mui/icons-material/Delete";
import EditIcon from "@mui/icons-material/Edit";
import LaunchIcon from "@mui/icons-material/Launch";
import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import CircularProgress from "@mui/material/CircularProgress";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogTitle from "@mui/material/DialogTitle";
import IconButton from "@mui/material/IconButton";
import Paper from "@mui/material/Paper";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import {
  Achievement,
  AchievementType,
  AddAchievementRequest,
  AddAchievementResponse,
  DeleteAchievementRequest,
  Handle,
  ListAchievementsRequest,
} from "@psankar/vetchi-typespec";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";

interface AchievementProps {
  userHandle: Handle;
  achievementType: AchievementType;
  canEdit: boolean;
}

export function AchievementSection({
  userHandle,
  achievementType,
  canEdit,
}: AchievementProps) {
  const router = useRouter();
  const { t } = useTranslation();
  const [achievements, setAchievements] = useState<Achievement[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const [openDialog, setOpenDialog] = useState(false);
  const [openDeleteDialog, setOpenDeleteDialog] = useState(false);
  const [openUrlWarningDialog, setOpenUrlWarningDialog] = useState(false);
  const [selectedUrl, setSelectedUrl] = useState("");
  const [selectedAchievement, setSelectedAchievement] =
    useState<Achievement | null>(null);

  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [url, setUrl] = useState("");
  const [date, setDate] = useState("");

  const [titleError, setTitleError] = useState("");
  const [descriptionError, setDescriptionError] = useState("");
  const [urlError, setUrlError] = useState("");
  const [dateError, setDateError] = useState("");

  // Translation keys based on achievement type
  const transSection =
    achievementType === AchievementType.PATENT
      ? "achievements.patents"
      : achievementType === AchievementType.PUBLICATION
      ? "achievements.publications"
      : "achievements.certifications";

  const fetchAchievements = useCallback(async () => {
    setIsLoading(true);
    setError(null);

    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        return;
      }

      const request: ListAchievementsRequest = {
        type: achievementType,
        handle: userHandle,
      };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/list-achievements`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(request),
          credentials: "include",
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/login");
        return;
      }

      if (!response.ok) {
        throw new Error(t(`${transSection}.error.fetchFailed`));
      }

      const data = await response.json();
      setAchievements(data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setIsLoading(false);
    }
  }, [achievementType, userHandle, t, transSection, router]);

  useEffect(() => {
    fetchAchievements();
  }, [fetchAchievements]);

  const resetForm = () => {
    setTitle("");
    setDescription("");
    setUrl("");
    setDate("");
    setTitleError("");
    setDescriptionError("");
    setUrlError("");
    setDateError("");
  };

  const handleOpenDialog = () => {
    resetForm();
    setOpenDialog(true);
  };

  const handleCloseDialog = () => {
    setOpenDialog(false);
    resetForm();
  };

  const handleUrlChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    let inputUrl = event.target.value;

    // Remove any protocol prefix the user might type or paste
    inputUrl = inputUrl.replace(/^(https?:\/\/|ftp:\/\/|\/\/)/i, "");

    setUrl(inputUrl);
    setUrlError("");
  };

  const validateForm = () => {
    let isValid = true;

    // Title validation
    if (!title.trim()) {
      setTitleError(t("common.error.requiredField"));
      isValid = false;
    } else if (title.length < 3 || title.length > 128) {
      setTitleError(t(`${transSection}.error.titleLength`));
      isValid = false;
    } else {
      setTitleError("");
    }

    // Description validation (optional)
    if (description.trim() && description.length > 1024) {
      setDescriptionError(t(`${transSection}.error.descriptionTooLong`));
      isValid = false;
    } else {
      setDescriptionError("");
    }

    // URL validation (optional)
    if (url.trim()) {
      if (url.length > 1024) {
        setUrlError(t(`${transSection}.error.urlTooLong`));
        isValid = false;
      } else {
        try {
          // Validate URL with https protocol
          const urlWithProtocol = `https://${url.trim()}`;
          new URL(urlWithProtocol);
          setUrlError("");
        } catch (e) {
          setUrlError(t(`${transSection}.error.invalidUrl`));
          isValid = false;
        }
      }
    }

    // Date validation (optional)
    if (date.trim()) {
      const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
      if (!dateRegex.test(date)) {
        setDateError(t(`${transSection}.error.invalidDate`));
        isValid = false;
      } else {
        const selectedDate = new Date(date);
        const today = new Date();
        today.setHours(0, 0, 0, 0);

        if (selectedDate > today) {
          setDateError(t(`${transSection}.error.futureDate`));
          isValid = false;
        } else {
          setDateError("");
        }
      }
    }

    return isValid;
  };

  const handleSubmit = async () => {
    if (!validateForm()) {
      return;
    }

    setIsSubmitting(true);

    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        return;
      }

      const request: AddAchievementRequest = {
        type: achievementType,
        title: title.trim(),
      };

      if (description.trim()) {
        request.description = description.trim();
      }

      if (url.trim()) {
        request.url = `https://${url.trim()}`;
      }

      if (date.trim()) {
        request.at = new Date(date);
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/add-achievement`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(request),
          credentials: "include",
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/login");
        return;
      }

      if (!response.ok) {
        throw new Error(t(`${transSection}.error.saveFailed`));
      }

      const data: AddAchievementResponse = await response.json();

      await fetchAchievements();
      handleCloseDialog();
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleOpenDeleteDialog = (achievement: Achievement) => {
    setSelectedAchievement(achievement);
    setOpenDeleteDialog(true);
  };

  const handleCloseDeleteDialog = () => {
    setOpenDeleteDialog(false);
    setSelectedAchievement(null);
  };

  const deleteAchievement = async (id: string, shouldRefresh = true) => {
    setIsSubmitting(true);

    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        return;
      }

      const request: DeleteAchievementRequest = { id };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/delete-achievement`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(request),
          credentials: "include",
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/login");
        return;
      }

      if (!response.ok) {
        throw new Error(t(`${transSection}.error.deleteFailed`));
      }

      if (shouldRefresh) {
        await fetchAchievements();
        handleCloseDeleteDialog();
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDelete = async () => {
    if (selectedAchievement) {
      await deleteAchievement(selectedAchievement.id);
    }
  };

  const formatDate = (dateString?: Date) => {
    if (!dateString) return "";

    const date = new Date(dateString);
    return date.toLocaleDateString();
  };

  const handleUrlClick = (url: string, event: React.MouseEvent) => {
    event.preventDefault();
    setSelectedUrl(url);
    setOpenUrlWarningDialog(true);
  };

  const handleExternalNavigation = () => {
    window.open(selectedUrl, "_blank", "noopener,noreferrer");
    setOpenUrlWarningDialog(false);
  };

  const handleCloseUrlWarningDialog = () => {
    setOpenUrlWarningDialog(false);
    setSelectedUrl("");
  };

  const handleEdit = (achievement: Achievement) => {
    setTitle(achievement.title);
    setDescription(achievement.description || "");
    setUrl(achievement.url ? achievement.url.replace(/^https?:\/\//, "") : "");
    setDate(
      achievement.at ? new Date(achievement.at).toISOString().split("T")[0] : ""
    );
    setSelectedAchievement(achievement);
    setOpenDialog(true);
  };

  return (
    <Box>
      <Box
        sx={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          mb: 2,
        }}
      >
        <Typography variant="h5" component="h2">
          {t(`${transSection}.title`)}
        </Typography>
        {canEdit && (
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={() => handleOpenDialog()}
          >
            {t(
              `${transSection}.add${
                achievementType.charAt(0) +
                achievementType.slice(1).toLowerCase()
              }`
            )}
          </Button>
        )}
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      {isLoading ? (
        <Box sx={{ display: "flex", justifyContent: "center" }}>
          <CircularProgress />
        </Box>
      ) : !achievements || achievements.length === 0 ? (
        <Paper
          elevation={1}
          sx={{
            p: { xs: 3, sm: 4 },
            textAlign: "center",
            bgcolor: "background.paper",
            borderRadius: 2,
          }}
        >
          <Typography color="text.secondary">
            {t(`${transSection}.noEntries`)}
          </Typography>
        </Paper>
      ) : (
        <Stack spacing={2}>
          {achievements.map((achievement) => (
            <Paper
              key={achievement.id}
              elevation={1}
              sx={{
                p: { xs: 3, sm: 4 },
                bgcolor: (theme) =>
                  theme.palette.mode === "light" ? "grey.50" : "grey.900",
                borderRadius: 2,
                transition: "all 0.2s ease-in-out",
                border: "1px solid",
                borderColor: "divider",
                "&:hover": {
                  boxShadow: (theme) => theme.shadows[2],
                  transform: canEdit ? "translateY(-2px)" : "none",
                  bgcolor: (theme) =>
                    theme.palette.mode === "light" ? "#ffffff" : "grey.800",
                },
              }}
            >
              <Box
                sx={{
                  display: "flex",
                  justifyContent: "space-between",
                  gap: 2,
                }}
              >
                <Box sx={{ flex: 1, minWidth: 0 }}>
                  <Typography
                    variant="h6"
                    gutterBottom
                    sx={{
                      color: "primary.main",
                      fontWeight: 600,
                    }}
                  >
                    {achievement.title}
                  </Typography>
                  {achievement.description && (
                    <Typography
                      variant="body2"
                      sx={{
                        color: "text.primary",
                        whiteSpace: "pre-wrap",
                        lineHeight: 1.6,
                        mb: achievement.url ? 2 : 0,
                      }}
                    >
                      {achievement.description}
                    </Typography>
                  )}
                  {achievement.url && (
                    <Button
                      variant="text"
                      color="primary"
                      size="small"
                      onClick={() => {
                        if (achievement.url) {
                          setSelectedUrl(achievement.url);
                          setOpenUrlWarningDialog(true);
                        }
                      }}
                      startIcon={<LaunchIcon />}
                      sx={{
                        mt: 1,
                        textTransform: "none",
                        "&:hover": {
                          bgcolor: "primary.lighter",
                        },
                      }}
                    >
                      {t(`${transSection}.url`)}
                    </Button>
                  )}
                </Box>
                {canEdit && (
                  <Box sx={{ display: "flex", gap: 1 }}>
                    <IconButton
                      onClick={() => handleEdit(achievement)}
                      color="primary"
                      size="small"
                      sx={{
                        "&:hover": {
                          bgcolor: "primary.lighter",
                        },
                      }}
                    >
                      <EditIcon />
                    </IconButton>
                    <IconButton
                      onClick={() => handleOpenDeleteDialog(achievement)}
                      color="error"
                      size="small"
                      sx={{
                        "&:hover": {
                          bgcolor: "error.lighter",
                        },
                      }}
                    >
                      <DeleteIcon />
                    </IconButton>
                  </Box>
                )}
              </Box>
            </Paper>
          ))}
        </Stack>
      )}

      <Dialog
        open={openUrlWarningDialog}
        onClose={() => setOpenUrlWarningDialog(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle sx={{ color: "warning.main" }}>
          {t("common.externalLink.warning")}
        </DialogTitle>
        <DialogContent>
          <Typography>{t("common.externalLink.message")}</Typography>
          <Typography
            variant="body2"
            sx={{
              mt: 2,
              color: "text.secondary",
              wordBreak: "break-all",
            }}
          >
            {selectedUrl}
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenUrlWarningDialog(false)}>
            {t("common.cancel")}
          </Button>
          <Button
            onClick={() => {
              window.open(selectedUrl, "_blank", "noopener,noreferrer");
              setOpenUrlWarningDialog(false);
            }}
            variant="contained"
            color="primary"
            startIcon={<LaunchIcon />}
          >
            {t("common.proceed")}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
