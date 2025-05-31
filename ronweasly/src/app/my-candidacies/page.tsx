"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import { config } from "@/config";
import { useAuth } from "@/hooks/useAuth";
import { useTranslation } from "@/hooks/useTranslation";
import ChevronRightIcon from "@mui/icons-material/ChevronRight";
import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import Chip from "@mui/material/Chip";
import CircularProgress from "@mui/material/CircularProgress";
import IconButton from "@mui/material/IconButton";
import Paper from "@mui/material/Paper";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import {
  CandidacyState,
  MyCandidaciesRequest,
  MyCandidacy,
} from "@vetchium/typespec";
import Cookies from "js-cookie";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

const getCandidacyStateColor = (
  state: CandidacyState
): "primary" | "success" | "error" | "warning" => {
  switch (state) {
    case "INTERVIEWING":
      return "primary";
    case "OFFERED":
      return "warning";
    case "OFFER_ACCEPTED":
      return "success";
    case "OFFER_DECLINED":
    case "CANDIDATE_UNSUITABLE":
    case "CANDIDATE_NOT_RESPONDING":
    case "CANDIDATE_WITHDREW":
    case "EMPLOYER_DEFUNCT":
      return "error";
    default:
      return "primary";
  }
};

export default function MyCandidaciesPage() {
  const router = useRouter();
  const { t } = useTranslation();
  useAuth(); // Check authentication and redirect if not authenticated
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [candidacies, setCandidacies] = useState<MyCandidacy[]>([]);
  const [paginationKey, setPaginationKey] = useState<string | null>(null);

  const fetchCandidacies = async () => {
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/get-my-candidacies`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            pagination_key: paginationKey,
            limit: 40,
          } as MyCandidaciesRequest),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/login");
        return;
      }

      if (!response.ok) {
        throw new Error(t("common.error.serverError"));
      }

      const data = await response.json();
      setCandidacies((prevCandidacies) =>
        paginationKey ? [...prevCandidacies, ...data] : data
      );
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("common.error.serverError")
      );
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCandidacies();
  }, [paginationKey]);

  if (loading && !paginationKey) {
    return (
      <AuthenticatedLayout>
        <Box sx={{ display: "flex", justifyContent: "center", mt: 4 }}>
          <CircularProgress />
        </Box>
      </AuthenticatedLayout>
    );
  }

  return (
    <AuthenticatedLayout>
      <Box sx={{ maxWidth: 800, mx: "auto", mt: 4, px: 2 }}>
        <Typography variant="h4" gutterBottom>
          {t("navigation.myCandidacies")}
        </Typography>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        {candidacies.length === 0 ? (
          <Paper sx={{ p: 4, textAlign: "center" }}>
            <Typography color="text.secondary">
              {t("candidacies.noCandidacies")}
            </Typography>
          </Paper>
        ) : (
          <Stack spacing={2}>
            {candidacies.map((candidacy) => (
              <Paper key={candidacy.candidacy_id} sx={{ p: 3 }}>
                <Box
                  sx={{
                    display: "flex",
                    justifyContent: "space-between",
                    alignItems: "flex-start",
                  }}
                >
                  <Box sx={{ flex: 1, mr: 2 }}>
                    <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                      <Typography variant="h6" gutterBottom>
                        {candidacy.opening_title}
                      </Typography>
                      <IconButton
                        component={Link}
                        href={`/candidacy/${candidacy.candidacy_id}`}
                        size="small"
                        color="primary"
                        aria-label={t("candidacies.viewDetails")}
                        sx={{
                          opacity: 0.7,
                          "&:hover": {
                            opacity: 1,
                            transform: "translateX(2px)",
                          },
                          transition: "transform 0.2s",
                        }}
                      >
                        <ChevronRightIcon />
                      </IconButton>
                    </Box>
                    <Typography
                      variant="subtitle1"
                      color="text.secondary"
                      gutterBottom
                    >
                      {candidacy.company_name}
                    </Typography>
                    <Typography variant="body2" sx={{ mt: 1 }}>
                      {candidacy.opening_description}
                    </Typography>
                  </Box>
                  <Chip
                    label={t(`candidacies.states.${candidacy.candidacy_state}`)}
                    color={getCandidacyStateColor(candidacy.candidacy_state)}
                    sx={{ ml: 2, flexShrink: 0 }}
                  />
                </Box>
              </Paper>
            ))}
          </Stack>
        )}

        {loading && paginationKey && (
          <Box sx={{ display: "flex", justifyContent: "center", mt: 2 }}>
            <CircularProgress />
          </Box>
        )}
      </Box>
    </AuthenticatedLayout>
  );
}
