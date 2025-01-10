"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import Box from "@mui/material/Box";
import Paper from "@mui/material/Paper";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import Chip from "@mui/material/Chip";
import Stack from "@mui/material/Stack";
import CircularProgress from "@mui/material/CircularProgress";
import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import {
  HubOpeningDetails,
  OpeningType,
  EducationLevel,
  OpeningTypes,
  EducationLevels,
  GetHubOpeningDetailsRequest,
  OpeningState,
  OpeningStates,
} from "@psankar/vetchi-typespec";
import { config } from "@/config";
import Cookies from "js-cookie";
import { useTranslation } from "@/hooks/useTranslation";

const formatEducationLevel = (
  level: EducationLevel,
  t: (key: string) => string
) => {
  switch (level) {
    case EducationLevels.BACHELOR:
      return t("openingDetails.educationLevel.bachelor");
    case EducationLevels.MASTER:
      return t("openingDetails.educationLevel.master");
    case EducationLevels.DOCTORATE:
      return t("openingDetails.educationLevel.doctorate");
    case EducationLevels.NOT_MATTERS:
      return t("openingDetails.educationLevel.notMatters");
    case EducationLevels.UNSPECIFIED:
      return t("openingDetails.educationLevel.unspecified");
    default:
      return level;
  }
};

const formatOpeningType = (type: OpeningType, t: (key: string) => string) => {
  switch (type) {
    case OpeningTypes.FULL_TIME:
      return t("openingDetails.openingType.fullTime");
    case OpeningTypes.PART_TIME:
      return t("openingDetails.openingType.partTime");
    case OpeningTypes.CONTRACT:
      return t("openingDetails.openingType.contract");
    case OpeningTypes.INTERNSHIP:
      return t("openingDetails.openingType.internship");
    case OpeningTypes.UNSPECIFIED:
      return t("openingDetails.openingType.unspecified");
    default:
      return type;
  }
};

export default function OpeningDetailsPage() {
  const { t } = useTranslation();
  const params = useParams();
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [opening, setOpening] = useState<HubOpeningDetails | null>(null);

  useEffect(() => {
    const fetchOpeningDetails = async () => {
      const token = Cookies.get("session_token");
      if (!token) {
        setError(t("common.error.notAuthenticated"));
        return;
      }

      try {
        const request: GetHubOpeningDetailsRequest = {
          company_domain: params.domain as string,
          opening_id_within_company: params.openingId as string,
        };

        const response = await fetch(
          `${config.API_SERVER_PREFIX}/hub/get-opening-details`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify(request),
          }
        );

        if (!response.ok) {
          if (response.status === 401) {
            setError(t("common.error.sessionExpired"));
            Cookies.remove("session_token", { path: "/" });
            return;
          }
          throw new Error(
            `Failed to fetch opening details: ${response.statusText}`
          );
        }

        const data = await response.json();
        setOpening(data);
      } catch (error) {
        console.error("Error fetching opening details:", error);
        setError(t("openingDetails.error.loadFailed"));
      } finally {
        setLoading(false);
      }
    };

    fetchOpeningDetails();
  }, [params.domain, params.openingId, t]);

  const handleApply = async () => {
    // TODO: Implement apply functionality
    alert("Apply functionality will be implemented soon!");
  };

  const canApply = opening?.state === OpeningStates.ACTIVE;

  const getOpeningStateMessage = (state: OpeningState) => {
    switch (state) {
      case OpeningStates.DRAFT:
        return t("openingDetails.state.draft");
      case OpeningStates.SUSPENDED:
        return t("openingDetails.state.suspended");
      case OpeningStates.CLOSED:
        return t("openingDetails.state.closed");
      default:
        return null;
    }
  };

  if (loading) {
    return (
      <AuthenticatedLayout>
        <Box sx={{ display: "flex", justifyContent: "center", mt: 4 }}>
          <CircularProgress />
        </Box>
      </AuthenticatedLayout>
    );
  }

  if (error) {
    return (
      <AuthenticatedLayout>
        <Paper sx={{ p: 2, mb: 2, bgcolor: "error.light" }}>
          <Typography color="error" align="center">
            {error}
          </Typography>
        </Paper>
      </AuthenticatedLayout>
    );
  }

  if (!opening) {
    return (
      <AuthenticatedLayout>
        <Paper sx={{ p: 2, mb: 2 }}>
          <Typography align="center">{t("openingDetails.notFound")}</Typography>
        </Paper>
      </AuthenticatedLayout>
    );
  }

  return (
    <AuthenticatedLayout>
      <Box sx={{ maxWidth: 800, mx: "auto", mt: 4 }}>
        <Paper sx={{ p: 4 }}>
          <Typography variant="h4" gutterBottom>
            {opening.job_title}
          </Typography>
          <Typography variant="h6" color="text.secondary" gutterBottom>
            {opening.company_name}
          </Typography>

          <Stack direction="row" spacing={1} sx={{ mb: 3 }}>
            {opening.opening_type && (
              <Chip label={formatOpeningType(opening.opening_type, t)} />
            )}
            {opening.education_level && (
              <Chip label={formatEducationLevel(opening.education_level, t)} />
            )}
            {opening.yoe_min !== undefined && opening.yoe_max !== undefined && (
              <Chip
                label={t("openingDetails.yearsExperience", {
                  min: opening.yoe_min,
                  max: opening.yoe_max,
                })}
              />
            )}
          </Stack>

          <Typography variant="body1" paragraph>
            {opening.jd}
          </Typography>

          {opening.hiring_manager_name && (
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
              {t("openingDetails.hiringManager")}: {opening.hiring_manager_name}
            </Typography>
          )}

          <Box sx={{ mt: 4 }}>
            {canApply ? (
              <Button
                variant="contained"
                color="primary"
                size="large"
                onClick={handleApply}
                fullWidth
              >
                {t("openingDetails.apply")}
              </Button>
            ) : (
              <Typography
                variant="body1"
                color="error"
                align="center"
                sx={{ mb: 2 }}
              >
                {getOpeningStateMessage(opening.state)}
              </Typography>
            )}
          </Box>
        </Paper>
      </Box>
    </AuthenticatedLayout>
  );
}
