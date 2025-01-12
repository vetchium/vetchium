"use client";

import { useEffect, useState, use } from "react";
import {
  Box,
  Paper,
  Typography,
  Alert,
  CircularProgress,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  IconButton,
  Card,
  CardContent,
  Chip,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Pagination,
  Tooltip,
  Menu,
  Stack,
} from "@mui/material";
import { useTranslation } from "@/hooks/useTranslation";
import { config } from "@/config";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import {
  ApplicationState,
  ApplicationColorTag,
} from "@psankar/vetchi-typespec";
import CloseIcon from "@mui/icons-material/Close";
import VisibilityIcon from "@mui/icons-material/Visibility";
import { LoadingButton } from "@mui/lab";
import FiberManualRecordIcon from "@mui/icons-material/FiberManualRecord";
import ColorLensIcon from "@mui/icons-material/ColorLens";

interface Application {
  id: string;
  cover_letter?: string;
  created_at: string;
  filename: string;
  hub_user_handle: string;
  hub_user_last_employer_domain?: string;
  resume: string;
  state: ApplicationState;
  color_tag?: ApplicationColorTag;
}

interface PageProps {
  params: Promise<{
    id: string;
  }>;
}

const ITEMS_PER_PAGE = 10;

export default function ApplicationsPage({ params }: PageProps) {
  const { id } = use(params);
  const [applications, setApplications] = useState<Application[]>([]);
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const [selectedResume, setSelectedResume] = useState<string | null>(null);
  const [isActionLoading, setIsActionLoading] = useState(false);
  const [colorTagFilter, setColorTagFilter] = useState<
    ApplicationColorTag | ""
  >("");
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [paginationKey, setPaginationKey] = useState<string | null>(null);
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [selectedApplicationId, setSelectedApplicationId] = useState<
    string | null
  >(null);

  const { t } = useTranslation();
  const router = useRouter();

  const fetchApplications = async (newPage?: number) => {
    setIsLoading(true);
    setError("");

    try {
      const sessionToken = Cookies.get("session_token");
      if (!sessionToken) {
        setError(t("auth.unauthorized"));
        setIsLoading(false);
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/get-applications`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${sessionToken}`,
          },
          body: JSON.stringify({
            opening_id: id,
            state: "APPLIED",
            limit: ITEMS_PER_PAGE,
            pagination_key: newPage === 1 ? null : paginationKey,
            color_tag_filter: colorTagFilter || undefined,
          }),
        }
      );

      if (response.status === 200) {
        const data = await response.json();
        setApplications(data || []);
        if (data?.length === ITEMS_PER_PAGE) {
          setPaginationKey(data[data.length - 1].id);
        }
      } else if (response.status === 401) {
        setError(t("auth.unauthorized"));
      } else {
        setError(t("common.error"));
      }
    } catch (err) {
      setError(t("common.error"));
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchApplications(1);
  }, [id, colorTagFilter]);

  const handleAction = async (
    applicationId: string,
    action: "shortlist" | "reject"
  ) => {
    setIsActionLoading(true);
    setError("");

    try {
      const sessionToken = Cookies.get("session_token");
      if (!sessionToken) {
        setError(t("auth.unauthorized"));
        return;
      }

      const endpoint =
        action === "shortlist" ? "shortlist-application" : "reject-application";
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/${endpoint}`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${sessionToken}`,
          },
          body: JSON.stringify({
            application_id: applicationId,
          }),
        }
      );

      if (response.status === 200) {
        fetchApplications(page);
      } else if (response.status === 401) {
        setError(t("auth.unauthorized"));
      } else {
        setError(t("common.error"));
      }
    } catch (err) {
      setError(t("common.error"));
    } finally {
      setIsActionLoading(false);
    }
  };

  const handleColorTag = async (
    applicationId: string,
    colorTag: ApplicationColorTag | null
  ) => {
    setIsActionLoading(true);
    setError("");

    try {
      const sessionToken = Cookies.get("session_token");
      if (!sessionToken) {
        setError(t("auth.unauthorized"));
        return;
      }

      const endpoint = colorTag
        ? "set-application-color-tag"
        : "remove-application-color-tag";
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/${endpoint}`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${sessionToken}`,
          },
          body: JSON.stringify(
            colorTag
              ? { application_id: applicationId, color_tag: colorTag }
              : { application_id: applicationId }
          ),
        }
      );

      if (response.status === 200) {
        fetchApplications(page);
      } else if (response.status === 401) {
        setError(t("auth.unauthorized"));
      } else {
        setError(t("common.error"));
      }
    } catch (err) {
      setError(t("common.error"));
    } finally {
      setIsActionLoading(false);
    }
  };

  const handlePageChange = (
    event: React.ChangeEvent<unknown>,
    value: number
  ) => {
    setPage(value);
    fetchApplications(value);
  };

  if (isLoading) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", my: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ width: "100%", p: 3 }}>
      <Box sx={{ display: "flex", justifyContent: "space-between", mb: 3 }}>
        <Typography variant="h4">{t("applications.title")}</Typography>
        <Button variant="outlined" onClick={() => router.back()}>
          {t("common.back")}
        </Button>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError("")}>
          {error}
        </Alert>
      )}

      <Box sx={{ mb: 3 }}>
        <FormControl sx={{ minWidth: 200 }}>
          <InputLabel>{t("applications.filterByColor")}</InputLabel>
          <Select
            value={colorTagFilter}
            onChange={(e) =>
              setColorTagFilter(e.target.value as ApplicationColorTag | "")
            }
            label={t("applications.filterByColor")}
          >
            <MenuItem value="">
              <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                {t("applications.allColors")}
              </Box>
            </MenuItem>
            <MenuItem value="GREEN">
              <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                <FiberManualRecordIcon sx={{ color: "success.main" }} />
                {t("applications.colorGreen")}
              </Box>
            </MenuItem>
            <MenuItem value="YELLOW">
              <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                <FiberManualRecordIcon sx={{ color: "warning.main" }} />
                {t("applications.colorYellow")}
              </Box>
            </MenuItem>
            <MenuItem value="RED">
              <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                <FiberManualRecordIcon sx={{ color: "error.main" }} />
                {t("applications.colorRed")}
              </Box>
            </MenuItem>
          </Select>
        </FormControl>
      </Box>

      <Stack spacing={3}>
        {applications?.map((application) => (
          <Card key={application.id}>
            <CardContent>
              <Box
                sx={{
                  display: "flex",
                  justifyContent: "space-between",
                  mb: 2,
                }}
              >
                <Typography variant="h6">
                  {application.hub_user_handle}
                </Typography>
                <Box sx={{ display: "flex", gap: 1, alignItems: "center" }}>
                  {application.color_tag && (
                    <Box
                      sx={{ display: "flex", alignItems: "center", gap: 0.5 }}
                    >
                      <FiberManualRecordIcon
                        sx={{
                          color:
                            application.color_tag === "GREEN"
                              ? "success.main"
                              : application.color_tag === "YELLOW"
                              ? "warning.main"
                              : "error.main",
                          fontSize: "small",
                        }}
                      />
                      <IconButton
                        size="small"
                        onClick={() => handleColorTag(application.id, null)}
                        sx={{
                          padding: 0.5,
                          color: "text.secondary",
                          "&:hover": { bgcolor: "action.hover" },
                        }}
                      >
                        <CloseIcon sx={{ fontSize: "small" }} />
                      </IconButton>
                    </Box>
                  )}
                  {!application.color_tag && (
                    <Tooltip title={t("applications.setColor")}>
                      <IconButton
                        size="small"
                        onClick={(event) => {
                          setAnchorEl(event.currentTarget);
                          setSelectedApplicationId(application.id);
                        }}
                        sx={{
                          color: "text.secondary",
                          "&:hover": { bgcolor: "action.hover" },
                        }}
                      >
                        <ColorLensIcon fontSize="small" />
                      </IconButton>
                    </Tooltip>
                  )}
                  <IconButton
                    onClick={() => setSelectedResume(application.resume)}
                    color="primary"
                  >
                    <VisibilityIcon />
                  </IconButton>
                </Box>
              </Box>

              <Box sx={{ display: "flex", gap: 2 }}>
                <LoadingButton
                  variant="contained"
                  color="primary"
                  onClick={() => handleAction(application.id, "shortlist")}
                  loading={isActionLoading}
                >
                  {t("applications.shortlist")}
                </LoadingButton>
                <LoadingButton
                  variant="contained"
                  color="error"
                  onClick={() => handleAction(application.id, "reject")}
                  loading={isActionLoading}
                >
                  {t("applications.reject")}
                </LoadingButton>
              </Box>

              {application.cover_letter && (
                <Typography sx={{ mt: 2 }}>
                  {application.cover_letter}
                </Typography>
              )}
            </CardContent>
          </Card>
        ))}
      </Stack>

      {applications.length > 0 && (
        <Box sx={{ display: "flex", justifyContent: "center", mt: 3 }}>
          <Pagination
            count={totalPages}
            page={page}
            onChange={handlePageChange}
            color="primary"
          />
        </Box>
      )}

      <Dialog
        open={!!selectedResume}
        onClose={() => setSelectedResume(null)}
        maxWidth="lg"
        fullWidth
      >
        <DialogTitle>
          <Box
            sx={{
              display: "flex",
              justifyContent: "space-between",
              alignItems: "center",
            }}
          >
            {t("applications.resumePreview")}
            <IconButton onClick={() => setSelectedResume(null)}>
              <CloseIcon />
            </IconButton>
          </Box>
        </DialogTitle>
        <DialogContent>
          {selectedResume && (
            <Box sx={{ height: "80vh" }}>
              <iframe
                src={selectedResume}
                style={{ width: "100%", height: "100%", border: "none" }}
                title="Resume Preview"
              />
            </Box>
          )}
        </DialogContent>
      </Dialog>

      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={() => {
          setAnchorEl(null);
          setSelectedApplicationId(null);
        }}
      >
        <MenuItem
          onClick={() => {
            handleColorTag(selectedApplicationId!, "GREEN");
            setAnchorEl(null);
            setSelectedApplicationId(null);
          }}
        >
          <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
            <FiberManualRecordIcon sx={{ color: "success.main" }} />
            {t("applications.colorGreen")}
          </Box>
        </MenuItem>
        <MenuItem
          onClick={() => {
            handleColorTag(selectedApplicationId!, "YELLOW");
            setAnchorEl(null);
            setSelectedApplicationId(null);
          }}
        >
          <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
            <FiberManualRecordIcon sx={{ color: "warning.main" }} />
            {t("applications.colorYellow")}
          </Box>
        </MenuItem>
        <MenuItem
          onClick={() => {
            handleColorTag(selectedApplicationId!, "RED");
            setAnchorEl(null);
            setSelectedApplicationId(null);
          }}
        >
          <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
            <FiberManualRecordIcon sx={{ color: "error.main" }} />
            {t("applications.colorRed")}
          </Box>
        </MenuItem>
        {selectedApplicationId &&
          applications.find((a) => a.id === selectedApplicationId)
            ?.color_tag && (
            <MenuItem
              onClick={() => {
                handleColorTag(selectedApplicationId!, null);
                setAnchorEl(null);
                setSelectedApplicationId(null);
              }}
            >
              <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                <CloseIcon fontSize="small" />
                {t("applications.removeColor")}
              </Box>
            </MenuItem>
          )}
      </Menu>
    </Box>
  );
}
