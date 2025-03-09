"use client";

import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import BusinessIcon from "@mui/icons-material/Business";
import CloseIcon from "@mui/icons-material/Close";
import ColorLensIcon from "@mui/icons-material/ColorLens";
import EmailIcon from "@mui/icons-material/Email";
import FiberManualRecordIcon from "@mui/icons-material/FiberManualRecord";
import OpenInNewIcon from "@mui/icons-material/OpenInNew";
import PeopleIcon from "@mui/icons-material/People";
import ZoomInIcon from "@mui/icons-material/ZoomIn";
import {
  Alert,
  Box,
  Button,
  Card,
  CardContent,
  Chip,
  CircularProgress,
  Dialog,
  DialogContent,
  DialogTitle,
  Divider,
  FormControl,
  Grid,
  IconButton,
  InputLabel,
  Menu,
  MenuItem,
  Pagination,
  Paper,
  Select,
  Stack,
  Typography,
} from "@mui/material";
import { Application, ApplicationColorTag } from "@psankar/vetchi-typespec";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { use, useEffect, useState, useCallback } from "react";

interface PageProps {
  params: Promise<{
    id: string;
  }>;
}

interface ExtendedApplication extends Application {
  resumeUrl?: string;
}

const ITEMS_PER_PAGE = 10;

const getColorTagStyles = (colorTag: ApplicationColorTag) => {
  switch (colorTag) {
    case "GREEN":
      return {
        bgcolor: "success.main",
        color: "white",
        "& .MuiChip-deleteIcon": {
          color: "white",
          "&:hover": {
            color: "error.light",
          },
        },
      };
    case "YELLOW":
      return {
        bgcolor: "#FFD700", // Pure yellow color
        color: "black", // Black text for better contrast on yellow
        "& .MuiChip-deleteIcon": {
          color: "black",
          "&:hover": {
            color: "error.dark",
          },
        },
      };
    case "RED":
      return {
        bgcolor: "error.main",
        color: "white",
        "& .MuiChip-deleteIcon": {
          color: "white",
          "&:hover": {
            color: "error.light",
          },
        },
      };
  }
};

export default function ApplicationsPage({ params }: PageProps) {
  const { id } = use(params);
  const [applications, setApplications] = useState<ExtendedApplication[]>([]);
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const [zoomedResume, setZoomedResume] = useState<{
    url: string | null;
    applicationId: string | null;
  }>({
    url: null,
    applicationId: null,
  });
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

  const fetchApplications = useCallback(
    async (newPage?: number) => {
      console.log("Fetching applications, page:", newPage);
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
          console.log("Applications fetched:", data);
          console.log("Response data:", data);
          console.log("Setting applications state with:", data);
          setApplications(data ?? []);
          setTotalPages(Math.ceil((data.total_count || 0) / ITEMS_PER_PAGE));
          if (data.length === ITEMS_PER_PAGE) {
            setPaginationKey(data[data.length - 1].id);
          }
        } else if (response.status === 401) {
          setError(t("auth.unauthorized"));
        } else {
          setError(t("common.serverError"));
        }
      } catch {
        setError(t("common.serverError"));
      } finally {
        setIsLoading(false);
      }
    },
    [id, colorTagFilter, paginationKey, t]
  );

  useEffect(() => {
    console.log("Fetching applications for opening ID:", id);
    fetchApplications(1);
  }, [id, colorTagFilter, fetchApplications]);

  const handleViewResume = useCallback(
    async (application: ExtendedApplication) => {
      if (application.resumeUrl) {
        // Don't auto-open the resume in full screen
        return;
      }

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
      } catch {
        setError(t("common.serverError"));
      } finally {
        setLoadingResumes((prev) => ({ ...prev, [application.id]: false }));
      }
    },
    [t]
  );

  useEffect(() => {
    // Remove auto-opening of resume
    applications.forEach((application) => {
      if (!application.resumeUrl && !loadingResumes[application.id]) {
        handleViewResume(application);
      }
    });
  }, [applications]);

  useEffect(() => {
    console.log("Fetching applications, page:", page);
    fetchApplications(page);
  }, [fetchApplications, page]);

  useEffect(() => {
    if (loadingResumes) {
      applications.forEach((application) => {
        if (!application.resumeUrl && !loadingResumes[application.id]) {
          handleViewResume(application);
        }
      });
    }
  }, [applications, loadingResumes, handleViewResume]);

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
    } catch {
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
    } catch {
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

  if (isLoading) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", my: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ width: "100%", p: 3 }}>
      <Button
        variant="text"
        startIcon={<ArrowBackIcon />}
        onClick={() => router.push(`/openings/${id}`)}
        sx={{ mb: 3 }}
        size="small"
      >
        {t("openings.backToOpening")}
      </Button>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      <Box sx={{ mb: 3 }}>
        <FormControl sx={{ minWidth: 200 }}>
          <InputLabel>Filter by Color Tag</InputLabel>
          <Select
            value={colorTagFilter}
            label="Filter by Color Tag"
            onChange={(e) =>
              setColorTagFilter(e.target.value as ApplicationColorTag | "")
            }
          >
            <MenuItem value="">All</MenuItem>
            <MenuItem value="GREEN">
              <FiberManualRecordIcon sx={{ color: "success.main" }} /> Green
            </MenuItem>
            <MenuItem value="YELLOW">
              <FiberManualRecordIcon sx={{ color: "#FFD700" }} /> Yellow
            </MenuItem>
            <MenuItem value="RED">
              <FiberManualRecordIcon sx={{ color: "error.main" }} /> Red
            </MenuItem>
          </Select>
        </FormControl>
      </Box>

      <Grid container spacing={3}>
        {applications.map((application) => (
          <Grid item xs={12} key={application.id}>
            <Card elevation={2}>
              <CardContent>
                <Grid container spacing={2}>
                  <Grid item xs={12}>
                    <Box
                      sx={{
                        display: "flex",
                        justifyContent: "space-between",
                        mb: 2,
                      }}
                    >
                      <Typography variant="h6">
                        {application.hub_user_name}
                      </Typography>
                      <Button
                        variant="outlined"
                        color="error"
                        onClick={() => handleAction(application.id, "reject")}
                        disabled={isActionLoading}
                      >
                        Reject
                      </Button>
                    </Box>
                  </Grid>

                  <Grid item xs={12} md={8}>
                    <Typography
                      variant="body2"
                      color="textSecondary"
                      gutterBottom
                      sx={{ display: "flex", alignItems: "center", gap: 0.5 }}
                    >
                      @{application.hub_user_handle}
                      <IconButton
                        size="small"
                        href={`/u/${application.hub_user_handle}`}
                        target="_blank"
                        sx={{ padding: "2px" }}
                      >
                        <OpenInNewIcon sx={{ fontSize: 16 }} />
                      </IconButton>
                    </Typography>
                    <Typography variant="body1" paragraph>
                      {application.hub_user_short_bio}
                    </Typography>

                    {application.hub_user_last_employer_domains && (
                      <Box
                        sx={{ display: "flex", alignItems: "center", mb: 1 }}
                      >
                        <BusinessIcon sx={{ mr: 1, color: "text.secondary" }} />
                        <Typography variant="body2">
                          Last worked at:{" "}
                          {application.hub_user_last_employer_domains.join(
                            ", "
                          )}
                        </Typography>
                      </Box>
                    )}

                    {application.cover_letter && (
                      <Typography
                        variant="body2"
                        color="text.secondary"
                        sx={{ mt: 2 }}
                      >
                        <strong>Cover Letter:</strong>{" "}
                        {application.cover_letter}
                      </Typography>
                    )}
                  </Grid>

                  <Grid item xs={12} md={4}>
                    {application.resumeUrl ? (
                      <Box
                        sx={{
                          display: "flex",
                          flexDirection: "column",
                          gap: 2,
                        }}
                      >
                        <Box
                          sx={{
                            border: "1px solid #e0e0e0",
                            borderRadius: 1,
                            overflow: "hidden",
                            cursor: "pointer",
                            height: "300px",
                            position: "relative",
                          }}
                          onClick={() =>
                            setZoomedResume({
                              url: application.resumeUrl!,
                              applicationId: application.id,
                            })
                          }
                        >
                          <iframe
                            src={application.resumeUrl}
                            style={{
                              width: "100%",
                              height: "100%",
                              border: "none",
                            }}
                            title="Resume Preview"
                          />
                          <Box
                            sx={{
                              position: "absolute",
                              top: 0,
                              left: 0,
                              right: 0,
                              bottom: 0,
                              display: "flex",
                              alignItems: "center",
                              justifyContent: "center",
                              bgcolor: "rgba(0, 0, 0, 0.1)",
                              opacity: 0,
                              transition: "opacity 0.2s",
                              "&:hover": {
                                opacity: 1,
                              },
                            }}
                          >
                            <ZoomInIcon
                              sx={{ fontSize: 40, color: "common.white" }}
                            />
                          </Box>
                        </Box>
                        <Box
                          sx={{ display: "flex", gap: 1, alignItems: "center" }}
                        >
                          <Button
                            variant="outlined"
                            startIcon={<ColorLensIcon />}
                            onClick={(e) => {
                              setSelectedApplicationId(application.id);
                              setAnchorEl(e.currentTarget);
                            }}
                            fullWidth
                          >
                            {application.color_tag
                              ? "Change Color"
                              : "Color Tag"}
                          </Button>
                          {application.color_tag && (
                            <Chip
                              icon={
                                <FiberManualRecordIcon
                                  sx={{ color: "inherit" }}
                                />
                              }
                              label={application.color_tag}
                              onDelete={() =>
                                handleColorTag(application.id, null)
                              }
                              sx={{
                                borderRadius: 1,
                                ...getColorTagStyles(application.color_tag),
                              }}
                            />
                          )}
                        </Box>
                      </Box>
                    ) : (
                      <Box
                        sx={{
                          height: "300px",
                          display: "flex",
                          alignItems: "center",
                          justifyContent: "center",
                          border: "1px solid #e0e0e0",
                          borderRadius: 1,
                        }}
                      >
                        <CircularProgress />
                      </Box>
                    )}
                  </Grid>

                  {application.endorsers.length > 0 && (
                    <Grid item xs={12}>
                      <Divider sx={{ my: 2 }} />
                      <Paper
                        elevation={0}
                        sx={{
                          p: 2,
                          bgcolor: "grey.50",
                          border: "1px solid",
                          borderColor: "divider",
                          borderRadius: 2,
                        }}
                      >
                        <Typography
                          variant="subtitle2"
                          gutterBottom
                          sx={{
                            display: "flex",
                            alignItems: "center",
                            color: "text.primary",
                            fontWeight: 600,
                          }}
                        >
                          <PeopleIcon sx={{ mr: 1, color: "primary.main" }} />
                          Endorsements ({application.endorsers.length})
                        </Typography>
                        <Stack spacing={2} sx={{ mt: 2 }}>
                          {application.endorsers.map((endorser, index) => (
                            <Paper
                              key={index}
                              elevation={0}
                              sx={{
                                p: 2,
                                bgcolor: "background.paper",
                                borderRadius: 1,
                                border: "1px solid",
                                borderColor: "divider",
                                "&:hover": {
                                  bgcolor: "background.paper",
                                  boxShadow: 1,
                                },
                              }}
                            >
                              <Typography
                                variant="subtitle2"
                                sx={{
                                  display: "flex",
                                  alignItems: "center",
                                  gap: 0.5,
                                }}
                              >
                                {endorser.full_name}
                                <Box
                                  component="span"
                                  sx={{
                                    display: "flex",
                                    alignItems: "center",
                                    gap: 0.5,
                                    color: "text.secondary",
                                  }}
                                >
                                  (@{endorser.handle}
                                  <IconButton
                                    size="small"
                                    href={`/u/${endorser.handle}`}
                                    target="_blank"
                                    sx={{ padding: "2px" }}
                                  >
                                    <OpenInNewIcon sx={{ fontSize: 16 }} />
                                  </IconButton>
                                  )
                                </Box>
                              </Typography>
                              <Typography
                                variant="body2"
                                color="text.secondary"
                                sx={{ mt: 1 }}
                              >
                                {endorser.short_bio}
                              </Typography>
                              {endorser.current_company_domains && (
                                <Box
                                  sx={{
                                    display: "flex",
                                    flexDirection: "column",
                                    gap: 1,
                                    mt: 1,
                                  }}
                                >
                                  <Typography
                                    variant="caption"
                                    color="text.secondary"
                                    sx={{
                                      display: "flex",
                                      alignItems: "center",
                                      gap: 0.5,
                                    }}
                                  >
                                    <EmailIcon sx={{ fontSize: "small" }} />
                                    Confirmed Email Domains
                                  </Typography>
                                  <Box
                                    sx={{
                                      display: "flex",
                                      flexWrap: "wrap",
                                      gap: 1,
                                    }}
                                  >
                                    {endorser.current_company_domains.map(
                                      (domain, idx) => (
                                        <Chip
                                          key={idx}
                                          label={domain}
                                          size="small"
                                          variant="outlined"
                                        />
                                      )
                                    )}
                                  </Box>
                                </Box>
                              )}
                            </Paper>
                          ))}
                        </Stack>
                      </Paper>
                    </Grid>
                  )}

                  <Grid item xs={12}>
                    <Box
                      sx={{
                        display: "flex",
                        justifyContent: "flex-start",
                        mt: 2,
                      }}
                    >
                      <Button
                        variant="contained"
                        color="success"
                        onClick={() =>
                          handleAction(application.id, "shortlist")
                        }
                        disabled={isActionLoading}
                      >
                        Shortlist
                      </Button>
                    </Box>
                  </Grid>
                </Grid>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      {applications.length > 0 && (
        <Box sx={{ display: "flex", justifyContent: "center", mt: 4 }}>
          <Pagination
            count={totalPages}
            page={page}
            onChange={handlePageChange}
            color="primary"
          />
        </Box>
      )}

      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={() => setAnchorEl(null)}
      >
        {(!selectedApplicationId ||
          !applications.find((a) => a.id === selectedApplicationId)
            ?.color_tag ||
          applications.find((a) => a.id === selectedApplicationId)
            ?.color_tag !== "GREEN") && (
          <MenuItem
            onClick={() => {
              handleColorTag(selectedApplicationId!, "GREEN");
              setAnchorEl(null);
            }}
          >
            <FiberManualRecordIcon sx={{ color: "success.main", mr: 1 }} />{" "}
            Green
          </MenuItem>
        )}
        {(!selectedApplicationId ||
          !applications.find((a) => a.id === selectedApplicationId)
            ?.color_tag ||
          applications.find((a) => a.id === selectedApplicationId)
            ?.color_tag !== "YELLOW") && (
          <MenuItem
            onClick={() => {
              handleColorTag(selectedApplicationId!, "YELLOW");
              setAnchorEl(null);
            }}
          >
            <FiberManualRecordIcon sx={{ color: "#FFD700", mr: 1 }} /> Yellow
          </MenuItem>
        )}
        {(!selectedApplicationId ||
          !applications.find((a) => a.id === selectedApplicationId)
            ?.color_tag ||
          applications.find((a) => a.id === selectedApplicationId)
            ?.color_tag !== "RED") && (
          <MenuItem
            onClick={() => {
              handleColorTag(selectedApplicationId!, "RED");
              setAnchorEl(null);
            }}
          >
            <FiberManualRecordIcon sx={{ color: "error.main", mr: 1 }} /> Red
          </MenuItem>
        )}
        {selectedApplicationId &&
          applications.find((a) => a.id === selectedApplicationId)
            ?.color_tag && (
            <MenuItem
              onClick={() => {
                handleColorTag(selectedApplicationId!, null);
                setAnchorEl(null);
              }}
            >
              <CloseIcon sx={{ mr: 1 }} /> Remove Tag
            </MenuItem>
          )}
      </Menu>

      <Dialog
        open={!!zoomedResume.url}
        onClose={() => setZoomedResume({ url: null, applicationId: null })}
        maxWidth={false}
        fullScreen
      >
        <DialogTitle
          sx={{
            m: 0,
            p: 2,
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
          }}
        >
          Resume
          <IconButton
            aria-label="close"
            onClick={() => setZoomedResume({ url: null, applicationId: null })}
            sx={{
              color: (theme) => theme.palette.grey[500],
            }}
          >
            <CloseIcon />
          </IconButton>
        </DialogTitle>
        <DialogContent dividers sx={{ p: 0 }}>
          {zoomedResume.url && (
            <iframe
              src={zoomedResume.url}
              style={{ width: "100%", height: "100%", border: "none" }}
              title="Resume Full Preview"
            />
          )}
        </DialogContent>
      </Dialog>
    </Box>
  );
}
