"use client";

import { useParams } from "next/navigation";
import { useEffect, useState } from "react";
import {
  Candidacy,
  CandidacyState,
  FilterCandidacyInfosRequest,
} from "@psankar/vetchi-typespec";
import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import {
  Box,
  Card,
  CardContent,
  Typography,
  Alert,
  CircularProgress,
  Chip,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
} from "@mui/material";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";

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

function CandidaciesTable({
  candidacies,
  t,
}: {
  candidacies: Candidacy[];
  t: (key: string) => string;
}) {
  return (
    <TableContainer component={Paper}>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>{t("candidacies.applicantName")}</TableCell>
            <TableCell>{t("candidacies.handle")}</TableCell>
            <TableCell>{t("candidacies.state")}</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {candidacies.map((candidacy) => (
            <TableRow key={candidacy.candidacy_id}>
              <TableCell>{candidacy.applicant_name}</TableCell>
              <TableCell>{candidacy.applicant_handle}</TableCell>
              <TableCell>
                <CandidacyStateLabel state={candidacy.candidacy_state} t={t} />
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
  const [candidacies, setCandidacies] = useState<Candidacy[]>([]);
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
            error instanceof Error ? error.message : t("candidacies.fetchError")
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
        <Box sx={{ display: "flex", justifyContent: "space-between", mb: 3 }}>
          <Typography variant="h4">{t("candidacies.title")}</Typography>
          <Button variant="outlined" onClick={() => router.back()}>
            {t("common.back")}
          </Button>
        </Box>
        <Typography color="text.secondary">
          {t("candidacies.noCandidacies")}
        </Typography>
      </Box>
    );
  }

  return (
    <Box sx={{ p: 3 }}>
      <Box sx={{ display: "flex", justifyContent: "space-between", mb: 3 }}>
        <Typography variant="h4">{t("candidacies.title")}</Typography>
        <Button variant="outlined" onClick={() => router.back()}>
          {t("common.back")}
        </Button>
      </Box>

      {/* Opening Details */}
      <OpeningDetailsCard opening={candidacies[0]} t={t} />

      {/* Candidacies Table */}
      <CandidaciesTable candidacies={candidacies} t={t} />
    </Box>
  );
}
