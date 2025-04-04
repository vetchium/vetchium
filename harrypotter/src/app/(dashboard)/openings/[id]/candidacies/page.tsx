"use client";

import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import {
  FilterList as FilterListIcon,
  OpenInNew as OpenInNewIcon,
  Visibility as VisibilityIcon,
} from "@mui/icons-material";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import {
  Alert,
  Autocomplete,
  Box,
  Button,
  Card,
  CardContent,
  Chip,
  CircularProgress,
  IconButton,
  Paper,
  Popover,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField,
  Tooltip,
  Typography,
} from "@mui/material";
import {
  Candidacy,
  CandidacyState,
  FilterCandidacyInfosRequest,
} from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";

interface Endorser {
  full_name: string;
  short_bio: string;
  handle: string;
  current_company_domains?: string[];
}

interface CandidacyWithEndorsers extends Candidacy {
  endorsers?: Endorser[];
}

function CandidacyStateLabel({
  state,
  t,
}: {
  state: CandidacyState;
  t: (key: string) => string;
}) {
  let color:
    | "primary"
    | "secondary"
    | "error"
    | "info"
    | "success"
    | "warning" = "info";
  switch (state) {
    case "INTERVIEWING":
      color = "info";
      break;
    case "OFFERED":
      color = "warning";
      break;
    case "OFFER_ACCEPTED":
      color = "success";
      break;
    case "OFFER_DECLINED":
    case "CANDIDATE_UNSUITABLE":
    case "CANDIDATE_NOT_RESPONDING":
    case "CANDIDATE_WITHDREW":
    case "EMPLOYER_DEFUNCT":
      color = "error";
      break;
  }
  return (
    <Chip label={t(`candidacies.states.${state}`)} color={color} size="small" />
  );
}

function OpeningDetailsCard({
  opening,
  t,
}: {
  opening: Candidacy;
  t: (key: string) => string;
}) {
  return (
    <Card sx={{ mb: 3 }}>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          {t("openings.details")}
        </Typography>
        <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
          <Box>
            <Typography variant="subtitle2" component="span">
              Opening ID:{" "}
            </Typography>
            <Typography component="span">{opening.opening_id}</Typography>
          </Box>
          <Box>
            <Typography variant="subtitle2" component="span">
              Title:{" "}
            </Typography>
            <Typography component="span">{opening.opening_title}</Typography>
          </Box>
          <Box>
            <Typography variant="subtitle2" gutterBottom>
              Description:
            </Typography>
            <Typography>{opening.opening_description}</Typography>
          </Box>
        </Box>
      </CardContent>
    </Card>
  );
}

function ColumnFilter({
  anchorEl,
  onClose,
  value,
  onChange,
  placeholder,
  options,
}: {
  anchorEl: HTMLElement | null;
  onClose: () => void;
  value: string;
  onChange: (value: string) => void;
  placeholder: string;
  options?: { value: string; label: string }[];
}) {
  return (
    <Popover
      open={Boolean(anchorEl)}
      anchorEl={anchorEl}
      onClose={onClose}
      anchorOrigin={{
        vertical: "bottom",
        horizontal: "left",
      }}
      transformOrigin={{
        vertical: "top",
        horizontal: "left",
      }}
    >
      <Box sx={{ p: 2, minWidth: 220 }}>
        {options ? (
          <Autocomplete
            size="small"
            options={options}
            getOptionLabel={(option) => option.label}
            renderInput={(params) => (
              <TextField {...params} placeholder={placeholder} />
            )}
            value={options.find((opt) => opt.value === value) || null}
            onChange={(_, newValue) => onChange(newValue?.value || "")}
            autoComplete
            autoHighlight
            autoFocus
          />
        ) : (
          <TextField
            size="small"
            placeholder={placeholder}
            value={value}
            onChange={(e) => onChange(e.target.value)}
            autoFocus
          />
        )}
      </Box>
    </Popover>
  );
}

function CandidaciesTable({
  candidacies,
  t,
}: {
  candidacies: CandidacyWithEndorsers[];
  t: (key: string) => string;
}) {
  const router = useRouter();
  const [filters, setFilters] = useState({
    applicantName: "",
    handle: "",
    state: "",
  });
  const [filterAnchors, setFilterAnchors] = useState<{
    [key: string]: HTMLElement | null;
  }>({
    applicantName: null,
    handle: null,
    state: null,
  });

  // Get unique states from candidacies and create options
  const stateOptions = Array.from(
    new Set(candidacies.map((c) => c.candidacy_state))
  ).map((state) => ({
    value: state,
    label: t(`candidacies.states.${state}`),
  }));

  const handleFilterClick = (
    event: React.MouseEvent<HTMLElement>,
    field: string
  ) => {
    setFilterAnchors((prev) => ({
      ...prev,
      [field]: event.currentTarget,
    }));
  };

  const handleFilterClose = (field: string) => {
    setFilterAnchors((prev) => ({
      ...prev,
      [field]: null,
    }));
  };

  const filteredCandidacies = candidacies.filter((candidacy) => {
    const nameMatch = candidacy.applicant_name
      .toLowerCase()
      .includes(filters.applicantName.toLowerCase());
    const handleMatch = candidacy.applicant_handle
      .toLowerCase()
      .includes(filters.handle.toLowerCase());
    const stateMatch = filters.state
      ? candidacy.candidacy_state === filters.state
      : true;

    return nameMatch && handleMatch && stateMatch;
  });

  return (
    <TableContainer component={Paper}>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>
              <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                <Typography>{t("candidacies.applicantName")}</Typography>
                <IconButton
                  size="small"
                  onClick={(e) => handleFilterClick(e, "applicantName")}
                  color={filters.applicantName ? "primary" : "default"}
                >
                  <FilterListIcon />
                </IconButton>
              </Box>
              <ColumnFilter
                anchorEl={filterAnchors.applicantName}
                onClose={() => handleFilterClose("applicantName")}
                value={filters.applicantName}
                onChange={(value) =>
                  setFilters((prev) => ({ ...prev, applicantName: value }))
                }
                placeholder={`${t("candidacies.filterPlaceholder")} ${t(
                  "candidacies.applicantName"
                ).toLowerCase()}`}
              />
            </TableCell>
            <TableCell>
              <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                <Typography>{t("candidacies.handle")}</Typography>
                <IconButton
                  size="small"
                  onClick={(e) => handleFilterClick(e, "handle")}
                  color={filters.handle ? "primary" : "default"}
                >
                  <FilterListIcon />
                </IconButton>
              </Box>
              <ColumnFilter
                anchorEl={filterAnchors.handle}
                onClose={() => handleFilterClose("handle")}
                value={filters.handle}
                onChange={(value) =>
                  setFilters((prev) => ({ ...prev, handle: value }))
                }
                placeholder={`${t("candidacies.filterPlaceholder")} ${t(
                  "candidacies.handle"
                ).toLowerCase()}`}
              />
            </TableCell>
            <TableCell>
              <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                <Typography>{t("candidacies.state")}</Typography>
                <IconButton
                  size="small"
                  onClick={(e) => handleFilterClick(e, "state")}
                  color={filters.state ? "primary" : "default"}
                >
                  <FilterListIcon />
                </IconButton>
              </Box>
              <ColumnFilter
                anchorEl={filterAnchors.state}
                onClose={() => handleFilterClose("state")}
                value={filters.state}
                onChange={(value) =>
                  setFilters((prev) => ({ ...prev, state: value }))
                }
                placeholder={`${t("candidacies.filterPlaceholder")} ${t(
                  "candidacies.state"
                ).toLowerCase()}`}
                options={stateOptions}
              />
            </TableCell>
            <TableCell>
              <Typography>{t("candidacies.endorsers")}</Typography>
            </TableCell>
            <TableCell>
              <Typography>{t("common.actions")}</Typography>
            </TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {filteredCandidacies.map((candidacy) => (
            <TableRow key={candidacy.candidacy_id}>
              <TableCell>{candidacy.applicant_name}</TableCell>
              <TableCell>
                {candidacy.applicant_handle}
                <IconButton
                  size="small"
                  href={`/u/${candidacy.applicant_handle}`}
                  target="_blank"
                  component="a"
                  sx={{ ml: 1 }}
                >
                  <OpenInNewIcon fontSize="small" />
                </IconButton>
              </TableCell>
              <TableCell>
                <CandidacyStateLabel state={candidacy.candidacy_state} t={t} />
              </TableCell>
              <TableCell>
                {candidacy.endorsers && candidacy.endorsers.length > 0 ? (
                  <Tooltip
                    title={
                      <Box>
                        {candidacy.endorsers.map((endorser, idx) => (
                          <Typography key={idx} variant="body2">
                            {endorser.full_name} (@{endorser.handle})
                            {endorser.current_company_domains &&
                              endorser.current_company_domains.length > 0 && (
                                <Typography
                                  variant="caption"
                                  component="div"
                                  color="text.secondary"
                                >
                                  {endorser.current_company_domains.join(", ")}
                                </Typography>
                              )}
                          </Typography>
                        ))}
                      </Box>
                    }
                  >
                    <Chip
                      label={`${candidacy.endorsers.length} ${
                        candidacy.endorsers.length === 1
                          ? t("candidacies.endorser")
                          : t("candidacies.endorsers")
                      }`}
                      size="small"
                      color="primary"
                    />
                  </Tooltip>
                ) : (
                  <Typography variant="body2" color="text.secondary">
                    {t("candidacies.noEndorsers")}
                  </Typography>
                )}
              </TableCell>
              <TableCell>
                <IconButton
                  color="primary"
                  onClick={() =>
                    router.push(`/candidacy/${candidacy.candidacy_id}`)
                  }
                  title={t("candidacies.viewCandidacy")}
                >
                  <VisibilityIcon />
                </IconButton>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
}

function LoadingSkeleton() {
  return (
    <Box sx={{ display: "flex", justifyContent: "center", my: 4 }}>
      <CircularProgress />
    </Box>
  );
}

export default function CandidaciesPage() {
  const params = useParams();
  const openingId = params.id as string;
  const [candidacies, setCandidacies] = useState<CandidacyWithEndorsers[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { t } = useTranslation();
  const router = useRouter();

  // Scroll to top when component mounts
  useEffect(() => {
    window.scrollTo(0, 0);
  }, []);

  useEffect(() => {
    let isMounted = true;

    async function fetchCandidacies() {
      try {
        const sessionToken = Cookies.get("session_token");
        if (!sessionToken) {
          if (isMounted) {
            setError(t("auth.unauthorized"));
          }
          return;
        }

        const request: FilterCandidacyInfosRequest = {
          opening_id: openingId,
          limit: 40,
        };

        const response = await fetch(
          `${config.API_SERVER_PREFIX}/employer/filter-candidacy-infos`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${sessionToken}`,
            },
            body: JSON.stringify(request),
          }
        );

        if (!isMounted) return;

        if (response.status === 401) {
          Cookies.remove("session_token");
          router.push("/signin");
          return;
        }

        if (!response.ok) {
          throw new Error(t("candidacies.fetchError"));
        }

        const data = await response.json();
        if (isMounted) {
          setCandidacies(data);
        }
      } catch (error) {
        if (isMounted) {
          console.error("Error fetching candidacies:", error);
          setError(
            error instanceof Error ? error.message : t("common.serverError")
          );
        }
      } finally {
        if (isMounted) {
          setLoading(false);
        }
      }
    }

    fetchCandidacies();

    return () => {
      isMounted = false;
    };
    // We only want to refetch when the openingId changes
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [openingId]);

  if (loading) {
    return <LoadingSkeleton />;
  }

  if (error) {
    return <Alert severity="error">{error}</Alert>;
  }

  if (candidacies.length === 0) {
    return (
      <Box sx={{ p: 3 }}>
        <Button
          variant="text"
          startIcon={<ArrowBackIcon />}
          onClick={() => router.push(`/openings/${openingId}`)}
          sx={{ mb: 3 }}
          size="small"
        >
          {t("openings.backToOpening")}
        </Button>
        <Typography color="text.secondary">
          {t("candidacies.noCandidacies")}
        </Typography>
      </Box>
    );
  }

  return (
    <Box sx={{ p: 3 }}>
      <Button
        variant="text"
        startIcon={<ArrowBackIcon />}
        onClick={() => router.push(`/openings/${openingId}`)}
        sx={{ mb: 3 }}
        size="small"
      >
        {t("openings.backToOpening")}
      </Button>

      {/* Opening Details */}
      <OpeningDetailsCard opening={candidacies[0]} t={t} />

      {/* Candidacies Table */}
      <CandidaciesTable candidacies={candidacies} t={t} />
    </Box>
  );
}
