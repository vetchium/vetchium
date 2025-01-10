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
} from "@psankar/vetchi-typespec";
import { config } from "@/config";
import Cookies from "js-cookie";

const formatEducationLevel = (level: EducationLevel) => {
  switch (level) {
    case EducationLevels.BACHELOR:
      return "Bachelor's Degree";
    case EducationLevels.MASTER:
      return "Master's Degree";
    case EducationLevels.DOCTORATE:
      return "Doctorate";
    case EducationLevels.NOT_MATTERS:
      return "Any Education Level";
    case EducationLevels.UNSPECIFIED:
      return "Not Specified";
    default:
      return level;
  }
};

const formatOpeningType = (type: OpeningType) => {
  switch (type) {
    case OpeningTypes.FULL_TIME:
      return "Full Time";
    case OpeningTypes.PART_TIME:
      return "Part Time";
    case OpeningTypes.CONTRACT:
      return "Contract";
    case OpeningTypes.INTERNSHIP:
      return "Internship";
    case OpeningTypes.UNSPECIFIED:
      return "Not Specified";
    default:
      return type;
  }
};

export default function OpeningDetailsPage() {
  const params = useParams();
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [opening, setOpening] = useState<HubOpeningDetails | null>(null);

  useEffect(() => {
    const fetchOpeningDetails = async () => {
      const token = Cookies.get("session_token");
      if (!token) {
        setError("Not authenticated. Please log in again.");
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
            setError("Session expired. Please log in again.");
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
        setError("Failed to load opening details. Please try again later.");
      } finally {
        setLoading(false);
      }
    };

    fetchOpeningDetails();
  }, [params.domain, params.openingId]);

  const handleApply = async () => {
    // TODO: Implement apply functionality
    alert("Apply functionality will be implemented soon!");
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
          <Typography align="center">Opening not found</Typography>
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
              <Chip label={formatOpeningType(opening.opening_type)} />
            )}
            {opening.education_level && (
              <Chip label={formatEducationLevel(opening.education_level)} />
            )}
            {opening.yoe_min !== undefined && opening.yoe_max !== undefined && (
              <Chip
                label={`${opening.yoe_min}-${opening.yoe_max} years experience`}
              />
            )}
          </Stack>

          <Typography variant="body1" paragraph>
            {opening.jd}
          </Typography>

          {opening.hiring_manager_name && (
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
              Hiring Manager: {opening.hiring_manager_name}
            </Typography>
          )}

          <Box sx={{ mt: 4 }}>
            <Button
              variant="contained"
              color="primary"
              size="large"
              onClick={handleApply}
              fullWidth
            >
              Apply for this Opening
            </Button>
          </Box>
        </Paper>
      </Box>
    </AuthenticatedLayout>
  );
}
