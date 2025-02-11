import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import {
  WorkHistory as WorkHistoryType,
  AddWorkHistoryRequest,
  UpdateWorkHistoryRequest,
  DeleteWorkHistoryRequest,
  ListWorkHistoryRequest,
} from "@psankar/vetchi-typespec";
import { config } from "@/config";
import Cookies from "js-cookie";
import { useTranslation } from "@/hooks/useTranslation";
import Box from "@mui/material/Box";
import Paper from "@mui/material/Paper";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import Stack from "@mui/material/Stack";
import CircularProgress from "@mui/material/CircularProgress";
import Alert from "@mui/material/Alert";
import IconButton from "@mui/material/IconButton";
import EditIcon from "@mui/icons-material/Edit";
import DeleteIcon from "@mui/icons-material/Delete";

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
                value={formData.start_date}
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
                value={formData.end_date || ""}
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
        <Paper sx={{ p: 4, textAlign: "center" }}>
          <Typography color="text.secondary">
            {t("workHistory.noEntries")}
          </Typography>
        </Paper>
      ) : (
        <Stack spacing={2}>
          {workHistory.map((entry) => (
            <Paper key={entry.id} sx={{ p: 3 }}>
              <Box sx={{ display: "flex", justifyContent: "space-between" }}>
                <Box>
                  <Typography variant="h6" gutterBottom>
                    {entry.title}
                  </Typography>
                  <Typography
                    variant="subtitle1"
                    color="text.secondary"
                    gutterBottom
                  >
                    {entry.employer_name || entry.employer_domain}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    {formatDate(entry.start_date)} -{" "}
                    {entry.end_date
                      ? formatDate(entry.end_date)
                      : t("workHistory.present")}
                  </Typography>
                  {entry.description && (
                    <Typography variant="body2" sx={{ mt: 2 }}>
                      {entry.description}
                    </Typography>
                  )}
                </Box>
                {canEdit && !isEditing && !isAddingNew && (
                  <Box>
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
                    >
                      <EditIcon />
                    </IconButton>
                    <IconButton
                      onClick={() => handleDelete(entry.id)}
                      color="error"
                      size="small"
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
