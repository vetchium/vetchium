"use client";

import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import AddIcon from "@mui/icons-material/Add";
import DeleteIcon from "@mui/icons-material/Delete";
import LinkIcon from "@mui/icons-material/Link";
import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import CircularProgress from "@mui/material/CircularProgress";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";
import IconButton from "@mui/material/IconButton";
import Link from "@mui/material/Link";
import Stack from "@mui/material/Stack";
import TextField from "@mui/material/TextField";
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
          // Basic URL validation
          new URL(url.startsWith("http") ? url : `https://${url}`);
          setUrlError("");
        } catch (e) {
          setUrlError(t(`${transSection}.error.invalidUrl`));
          isValid = false;
        }
      }
    } else {
      setUrlError("");
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
    } else {
      setDateError("");
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
        request.url = url.trim().startsWith("http")
          ? url.trim()
          : `https://${url.trim()}`;
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
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      {isLoading ? (
        <Box sx={{ display: "flex", justifyContent: "center", my: 4 }}>
          <CircularProgress />
        </Box>
      ) : achievements.length === 0 ? (
        <Typography sx={{ mt: 2, mb: 4, fontStyle: "italic" }}>
          {t(`${transSection}.noEntries`)}
        </Typography>
      ) : (
        <Stack spacing={2} sx={{ mt: 2 }}>
          {achievements.map((achievement) => (
            <Card key={achievement.id} variant="outlined">
              <CardContent>
                <Box
                  sx={{
                    display: "flex",
                    justifyContent: "space-between",
                    alignItems: "flex-start",
                  }}
                >
                  <Box>
                    <Typography variant="h6" component="div">
                      {achievement.title}
                    </Typography>

                    {achievement.description && (
                      <Typography
                        variant="body2"
                        color="text.secondary"
                        sx={{ mt: 1 }}
                      >
                        {achievement.description}
                      </Typography>
                    )}

                    {achievement.at && (
                      <Typography
                        variant="body2"
                        color="text.secondary"
                        sx={{ mt: 1 }}
                      >
                        {t(`${transSection}.date`)}:{" "}
                        {formatDate(achievement.at)}
                      </Typography>
                    )}

                    {achievement.url && (
                      <Link
                        href={achievement.url}
                        target="_blank"
                        rel="noopener"
                        sx={{
                          display: "inline-flex",
                          alignItems: "center",
                          mt: 1,
                        }}
                      >
                        <LinkIcon fontSize="small" sx={{ mr: 0.5 }} />
                        {t(`${transSection}.url`)}
                      </Link>
                    )}
                  </Box>

                  {canEdit && (
                    <Box>
                      <IconButton
                        aria-label={t("achievements.actions.delete")}
                        onClick={() => handleOpenDeleteDialog(achievement)}
                      >
                        <DeleteIcon />
                      </IconButton>
                    </Box>
                  )}
                </Box>
              </CardContent>
            </Card>
          ))}
        </Stack>
      )}

      {/* Add Dialog */}
      <Dialog
        open={openDialog}
        onClose={handleCloseDialog}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>
          {t(
            `${transSection}.add${
              achievementType.charAt(0) + achievementType.slice(1).toLowerCase()
            }`
          )}
        </DialogTitle>
        <DialogContent>
          <TextField
            margin="dense"
            label={t(`${transSection}.title_field`)}
            type="text"
            fullWidth
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            error={!!titleError}
            helperText={titleError}
            required
            sx={{ mb: 2, mt: 1 }}
          />
          <TextField
            margin="dense"
            label={t(`${transSection}.description`)}
            type="text"
            fullWidth
            multiline
            rows={3}
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            error={!!descriptionError}
            helperText={descriptionError}
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            label={t(`${transSection}.url`)}
            type="text"
            fullWidth
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            error={!!urlError}
            helperText={urlError}
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            label={t(`${transSection}.date`)}
            type="date"
            fullWidth
            value={date}
            onChange={(e) => setDate(e.target.value)}
            error={!!dateError}
            helperText={dateError}
            InputLabelProps={{ shrink: true }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseDialog} disabled={isSubmitting}>
            {t("achievements.actions.cancel")}
          </Button>
          <Button
            onClick={handleSubmit}
            variant="contained"
            disabled={isSubmitting}
            startIcon={isSubmitting ? <CircularProgress size={20} /> : null}
          >
            {t("achievements.actions.save")}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <Dialog open={openDeleteDialog} onClose={handleCloseDeleteDialog}>
        <DialogTitle>{t("achievements.actions.delete")}</DialogTitle>
        <DialogContent>
          <DialogContentText>
            {t(`${transSection}.deleteConfirm`)}
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseDeleteDialog} disabled={isSubmitting}>
            {t("achievements.actions.cancel")}
          </Button>
          <Button
            onClick={handleDelete}
            color="error"
            disabled={isSubmitting}
            startIcon={isSubmitting ? <CircularProgress size={20} /> : null}
          >
            {t("achievements.actions.delete")}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
