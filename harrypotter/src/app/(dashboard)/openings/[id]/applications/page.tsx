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
import OpenInNewIcon from "@mui/icons-material/OpenInNew";

interface Application {
  id: string;
  cover_letter?: string;
  created_at: string;
  filename: string;
  hub_user_handle: string;
  hub_user_name: string;
  hub_user_short_bio: string;
  hub_user_last_employer_domains?: string[];
  resume: string;
  state: ApplicationState;
  color_tag?: ApplicationColorTag;
  resumeUrl?: string;
  endorsers: {
    full_name: string;
    short_bio: string;
    handle: string;
    current_company_domains?: string[];
  }[];
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
  const [loadingResumes, setLoadingResumes] = useState<{
    [key: string]: boolean;
  }>({});

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
        setError(t("common.serverError"));
      }
    } catch (err) {
      setError(t("common.serverError"));
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
        setError(t("common.serverError"));
      }
    } catch (err) {
      setError(t("common.serverError"));
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
        setError(t("common.serverError"));
      }
    } catch (err) {
      setError(t("common.serverError"));
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

  const fetchResume = async (application: Application) => {
    if (application.resumeUrl) return; // Already fetched

    setLoadingResumes((prev) => ({ ...prev, [application.id]: true }));
    setError("");

    try {
      const sessionToken = Cookies.get("session_token");
      if (!sessionToken) {
        setError(t("auth.unauthorized"));
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/get-resume`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${sessionToken}`,
          },
          body: JSON.stringify({
            application_id: application.id,
          }),
        }
      );

      if (response.status === 200) {
        const blob = await response.blob();
        const url = URL.createObjectURL(blob);
        setApplications((prevApps) =>
          prevApps.map((app) =>
            app.id === application.id ? { ...app, resumeUrl: url } : app
          )
        );
      } else if (response.status === 401) {
        setError(t("auth.unauthorized"));
      } else {
        setError(t("common.serverError"));
      }
    } catch (err) {
      setError(t("common.serverError"));
    } finally {
      setLoadingResumes((prev) => ({ ...prev, [application.id]: false }));
    }
  };

  // Fetch resumes for visible applications
  useEffect(() => {
    applications.forEach((application) => {
      if (!application.resumeUrl && !loadingResumes[application.id]) {
        fetchResume(application);
      }
    });
  }, [applications]);

  // Cleanup object URLs when component unmounts
  useEffect(() => {
    return () => {
      applications.forEach((app) => {
        if (app.resumeUrl) {
          URL.revokeObjectURL(app.resumeUrl);
        }
      });
    };
  }, []);

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
                <Box>
                  <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                    <Typography variant="h6">
                      {application.hub_user_name} (@
                      {application.hub_user_handle})
                    </Typography>
                    <IconButton
                      size="small"
                      href={`/u/${application.hub_user_handle}`}
                      target="_blank"
                      component="a"
                      sx={{ ml: 1 }}
                    >
                      <OpenInNewIcon fontSize="small" />
                    </IconButton>
                  </Box>
                  {application.hub_user_last_employer_domains && (
                    <Typography variant="body2" color="text.secondary">
                      {t("applications.lastEmployer")}:{" "}
                      {application.hub_user_last_employer_domains.join(", ")}
                    </Typography>
                  )}
                </Box>
              </Box>

              <Box
                sx={{
                  display: "flex",
                  flexDirection: "column",
                  alignItems: "center",
                  gap: 3,
                  mb: 3,
                }}
              >
                <Box
                  sx={{
                    position: "relative",
                    width: 100,
                    height: 140,
                    border: "1px solid",
                    borderColor: "divider",
                    borderRadius: 1,
                    overflow: "hidden",
                    cursor: "pointer",
                    "&:hover": {
                      boxShadow: 1,
                    },
                  }}
                  onClick={() =>
                    setSelectedResume(application.resumeUrl || null)
                  }
                >
                  {loadingResumes[application.id] ? (
                    <Box
                      sx={{
                        display: "flex",
                        alignItems: "center",
                        justifyContent: "center",
                        height: "100%",
                      }}
                    >
                      <CircularProgress size={24} />
                    </Box>
                  ) : application.resumeUrl ? (
                    <object
                      data={application.resumeUrl}
                      type="application/pdf"
                      width="100%"
                      height="100%"
                      style={{ pointerEvents: "none" }}
                    >
                      <Box
                        sx={{
                          display: "flex",
                          alignItems: "center",
                          justifyContent: "center",
                          height: "100%",
                          bgcolor: "action.hover",
                          color: "text.secondary",
                          p: 1,
                          textAlign: "center",
                        }}
                      >
                        <Typography variant="caption">
                          {t("applications.pdfPreviewNotAvailable")}
                        </Typography>
                      </Box>
                    </object>
                  ) : (
                    <Box
                      sx={{
                        display: "flex",
                        flexDirection: "column",
                        alignItems: "center",
                        justifyContent: "center",
                        height: "100%",
                        bgcolor: "action.hover",
                        color: "text.secondary",
                        p: 1,
                      }}
                    >
                      <CircularProgress size={24} />
                    </Box>
                  )}
                </Box>

                <Box
                  sx={{
                    display: "flex",
                    width: "100%",
                    justifyContent: "space-between",
                    alignItems: "center",
                    px: 4, // Add horizontal padding for better spacing from edges
                  }}
                >
                  <LoadingButton
                    variant="contained"
                    color="primary"
                    onClick={() => handleAction(application.id, "shortlist")}
                    loading={isActionLoading}
                    sx={{ minWidth: 120 }} // Ensure consistent button width
                  >
                    {t("applications.shortlist")}
                  </LoadingButton>

                  <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                    {application.color_tag ? (
                      <Box
                        sx={{
                          display: "flex",
                          alignItems: "center",
                          gap: 0.5,
                          border: "1px solid",
                          borderColor: "divider",
                          borderRadius: 1,
                          padding: "2px 8px", // Increased horizontal padding
                        }}
                      >
                        <FiberManualRecordIcon
                          sx={{
                            color:
                              application.color_tag === "GREEN"
                                ? "success.main"
                                : application.color_tag === "YELLOW"
                                ? "warning.main"
                                : "error.main",
                          }}
                        />
                        <IconButton
                          size="small"
                          onClick={() => handleColorTag(application.id, null)}
                        >
                          <CloseIcon fontSize="small" />
                        </IconButton>
                      </Box>
                    ) : (
                      <Tooltip title={t("applications.setColor")}>
                        <IconButton
                          size="medium" // Increased button size
                          onClick={(event) => {
                            setAnchorEl(event.currentTarget);
                            setSelectedApplicationId(application.id);
                          }}
                          sx={{
                            border: "1px solid",
                            borderColor: "divider",
                            p: 1,
                          }}
                        >
                          <ColorLensIcon />
                        </IconButton>
                      </Tooltip>
                    )}
                  </Box>

                  <LoadingButton
                    variant="contained"
                    color="error"
                    onClick={() => handleAction(application.id, "reject")}
                    loading={isActionLoading}
                    sx={{ minWidth: 120 }} // Ensure consistent button width
                  >
                    {t("applications.reject")}
                  </LoadingButton>
                </Box>
              </Box>

              {application.cover_letter && (
                <Typography sx={{ mt: 2 }}>
                  {application.cover_letter}
                </Typography>
              )}

              {/* Add Endorsers Section */}
              {application.endorsers && application.endorsers.length > 0 && (
                <Box sx={{ mt: 3 }}>
                  <Typography
                    variant="subtitle1"
                    sx={{ mb: 2, fontWeight: 500 }}
                  >
                    {t("applications.endorsers")}
                  </Typography>
                  <Box
                    sx={{ display: "flex", flexDirection: "column", gap: 2 }}
                  >
                    {application.endorsers.map((endorser, index) => (
                      <Paper
                        key={index}
                        sx={{ p: 2, bgcolor: "background.default" }}
                      >
                        <Box
                          sx={{
                            display: "flex",
                            justifyContent: "space-between",
                            alignItems: "flex-start",
                          }}
                        >
                          <Box>
                            <Typography variant="subtitle2">
                              {endorser.full_name} (@{endorser.handle})
                            </Typography>
                            <Typography variant="body2" color="text.secondary">
                              {endorser.short_bio}
                            </Typography>
                          </Box>
                          {endorser.current_company_domains &&
                            endorser.current_company_domains.length > 0 && (
                              <Box
                                sx={{
                                  display: "flex",
                                  gap: 1,
                                  flexWrap: "wrap",
                                }}
                              >
                                {endorser.current_company_domains.map(
                                  (domain, idx) => (
                                    <Chip
                                      key={idx}
                                      label={domain}
                                      size="small"
                                      variant="outlined"
                                      sx={{ bgcolor: "background.paper" }}
                                    />
                                  )
                                )}
                              </Box>
                            )}
                        </Box>
                      </Paper>
                    ))}
                  </Box>
                </Box>
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
        onClose={() => {
          setSelectedResume(null);
        }}
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
            <IconButton
              onClick={() => {
                setSelectedResume(null);
              }}
            >
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
                title={t("applications.resumePreview")}
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
