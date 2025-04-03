import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import DeleteIcon from "@mui/icons-material/Delete";
import EditIcon from "@mui/icons-material/Edit";
import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import CircularProgress from "@mui/material/CircularProgress";
import IconButton from "@mui/material/IconButton";
import Paper from "@mui/material/Paper";
import Stack from "@mui/material/Stack";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import {
  AddWorkHistoryRequest,
  DeleteWorkHistoryRequest,
  ListWorkHistoryRequest,
  UpdateWorkHistoryRequest,
  WorkHistory as WorkHistoryType,
} from "@psankar/vetchi-typespec";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

interface WorkHistoryProps {
  userHandle: string;
  canEdit: boolean;
}

type WorkHistoryFormData = Omit<AddWorkHistoryRequest, "id">;

function formatDate(dateString: string): string {
  return new Intl.DateTimeFormat(undefined, {
    year: "numeric",
    month: "long",
    day: "numeric",
  }).format(new Date(dateString));
}

function formatDateForInput(dateString: string): string {
  const date = new Date(dateString);
  // HTML date inputs require YYYY-MM-DD format as per spec
  return date.toISOString().split("T")[0];
}

export function WorkHistory({ userHandle, canEdit }: WorkHistoryProps) {
  const router = useRouter();
  const { t } = useTranslation();
  const [workHistory, setWorkHistory] = useState<WorkHistoryType[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isEditing, setIsEditing] = useState<string | null>(null);
  const [isAddingNew, setIsAddingNew] = useState(false);
  const [formData, setFormData] = useState<WorkHistoryFormData>({
    employer_domain: "",
    title: "",
    start_date: "",
  });

  useEffect(() => {
    fetchWorkHistory();
  }, [userHandle]);

  async function fetchWorkHistory() {
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        return;
      }

      const request: ListWorkHistoryRequest = { user_handle: userHandle };
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/list-work-history`,
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
        router.push("/login");
        return;
      }

      if (!response.ok) {
        throw new Error(t("workHistory.error.fetchFailed"));
      }

      const data = await response.json();
      setWorkHistory(data || []);
      setError(null);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("common.error.serverError")
      );
      setWorkHistory([]);
    } finally {
      setIsLoading(false);
    }
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();

    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        return;
      }

      const endpoint = isEditing
        ? `${config.API_SERVER_PREFIX}/hub/update-work-history`
        : `${config.API_SERVER_PREFIX}/hub/add-work-history`;

      const method = isEditing ? "PUT" : "POST";
      const body = isEditing
        ? ({ ...formData, id: isEditing } as UpdateWorkHistoryRequest)
        : (formData as AddWorkHistoryRequest);

      const response = await fetch(endpoint, {
        method,
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(body),
      });

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/login");
        return;
      }

      if (!response.ok) {
        throw new Error(t("workHistory.error.saveFailed"));
      }

      await fetchWorkHistory();
      setFormData({
        employer_domain: "",
        title: "",
        start_date: "",
      });
      setIsEditing(null);
      setIsAddingNew(false);
      setError(null);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("common.error.serverError")
      );
    }
  }

  async function handleDelete(id: string) {
    if (!confirm(t("workHistory.deleteConfirm"))) {
      return;
    }

    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        return;
      }

      const request: DeleteWorkHistoryRequest = { id };
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/delete-work-history`,
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
        router.push("/login");
        return;
      }

      if (!response.ok) {
        throw new Error(t("workHistory.error.deleteFailed"));
      }

      await fetchWorkHistory();
      setError(null);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("common.error.serverError")
      );
    }
  }

  if (isLoading) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", p: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box>
      {error && (
        <Alert severity="error" sx={{ mb: 4 }}>
          {error}
        </Alert>
      )}

      <Box
        sx={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          mb: 3,
        }}
      >
        <Typography variant="h5">{t("workHistory.title")}</Typography>
        {canEdit && !isAddingNew && !isEditing && (
          <Button
            variant="contained"
            color="primary"
            onClick={() => setIsAddingNew(true)}
          >
            {t("workHistory.addExperience")}
          </Button>
        )}
      </Box>

      {(isAddingNew || isEditing) && canEdit && (
        <Paper sx={{ p: 3, mb: 4 }}>
          <form onSubmit={handleSubmit}>
            <Stack spacing={3}>
              <TextField
                label={t("workHistory.companyDomain")}
                value={formData.employer_domain}
                onChange={(e) =>
                  setFormData({ ...formData, employer_domain: e.target.value })
                }
                required
                fullWidth
              />
              <TextField
                label={t("workHistory.jobTitle")}
                value={formData.title}
                onChange={(e) =>
                  setFormData({ ...formData, title: e.target.value })
                }
                required
                fullWidth
              />
              <TextField
                label={t("workHistory.startDate")}
                type="date"
                value={
                  formData.start_date
                    ? formatDateForInput(formData.start_date)
                    : ""
                }
                onChange={(e) =>
                  setFormData({ ...formData, start_date: e.target.value })
                }
                required
                fullWidth
                InputLabelProps={{ shrink: true }}
              />
              <TextField
                label={t("workHistory.endDate")}
                type="date"
                value={
                  formData.end_date ? formatDateForInput(formData.end_date) : ""
                }
                onChange={(e) =>
                  setFormData({ ...formData, end_date: e.target.value })
                }
                fullWidth
                InputLabelProps={{ shrink: true }}
              />
              <TextField
                label={t("workHistory.description")}
                value={formData.description || ""}
                onChange={(e) =>
                  setFormData({ ...formData, description: e.target.value })
                }
                multiline
                rows={4}
                fullWidth
              />
              <Box sx={{ display: "flex", gap: 2 }}>
                <Button type="submit" variant="contained" color="primary">
                  {isEditing
                    ? t("workHistory.actions.save")
                    : t("workHistory.addExperience")}
                </Button>
                <Button
                  variant="outlined"
                  color="inherit"
                  onClick={() => {
                    setIsEditing(null);
                    setIsAddingNew(false);
                    setFormData({
                      employer_domain: "",
                      title: "",
                      start_date: "",
                    });
                  }}
                >
                  {t("workHistory.actions.cancel")}
                </Button>
              </Box>
            </Stack>
          </form>
        </Paper>
      )}

      {!workHistory || workHistory.length === 0 ? (
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
            {t("workHistory.noEntries")}
          </Typography>
        </Paper>
      ) : (
        <Stack spacing={2}>
          {workHistory.map((entry) => (
            <Paper
              key={entry.id}
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
                    {entry.title}
                  </Typography>
                  <Typography
                    variant="subtitle1"
                    sx={{
                      color: "text.primary",
                      mb: 1,
                    }}
                  >
                    {entry.employer_name || entry.employer_domain}
                  </Typography>
                  <Typography
                    variant="body2"
                    sx={{
                      color: "text.secondary",
                      mb: entry.description ? 2 : 0,
                    }}
                  >
                    {formatDate(entry.start_date)} -{" "}
                    {entry.end_date
                      ? formatDate(entry.end_date)
                      : t("workHistory.present")}
                  </Typography>
                  {entry.description && (
                    <Typography
                      variant="body2"
                      sx={{
                        color: "text.primary",
                        whiteSpace: "pre-wrap",
                        lineHeight: 1.6,
                      }}
                    >
                      {entry.description}
                    </Typography>
                  )}
                </Box>
                {canEdit && !isEditing && !isAddingNew && (
                  <Box sx={{ display: "flex", gap: 1 }}>
                    <IconButton
                      onClick={() => {
                        setIsEditing(entry.id);
                        setFormData({
                          employer_domain: entry.employer_domain,
                          title: entry.title,
                          start_date: entry.start_date,
                          end_date: entry.end_date,
                          description: entry.description,
                        });
                      }}
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
                      onClick={() => handleDelete(entry.id)}
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
    </Box>
  );
}
