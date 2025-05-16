"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import CircularProgress from "@mui/material/CircularProgress";
import Paper from "@mui/material/Paper";
import Typography from "@mui/material/Typography";
import {
  GetEmployerDetailsRequest,
  HubEmployerDetails,
} from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";

export default function EmployerDetailsPage() {
  const { t } = useTranslation();
  const params = useParams();
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [employerDetails, setEmployerDetails] =
    useState<HubEmployerDetails | null>(null);

  if (!params?.domain) {
    return (
      <AuthenticatedLayout>
        <Box sx={{ p: 3 }}>
          <Typography color="error">
            {t("common.error.invalidParams")}
          </Typography>
          <Button
            variant="contained"
            onClick={() => router.back()}
            sx={{ mt: 2 }}
          >
            {t("common.back")}
          </Button>
        </Box>
      </AuthenticatedLayout>
    );
  }

  const companyDomain = params.domain as string;

  useEffect(() => {
    const fetchEmployerDetails = async () => {
      setLoading(true);
      setError(null);
      const token = Cookies.get("session_token");
      if (!token) {
        setError(t("common.error.notAuthenticated"));
        setLoading(false);
        return;
      }

      try {
        const request: GetEmployerDetailsRequest = {
          domain: companyDomain,
        };

        const response = await fetch(
          `${config.API_SERVER_PREFIX}/hub/get-employer-details`,
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
          } else if (response.status === 404) {
            setError(t("employerDetails.error.notFound"));
          } else {
            throw new Error(
              `Failed to fetch employer details: ${response.statusText}`
            );
          }
          setEmployerDetails(null);
        } else {
          const data: HubEmployerDetails = await response.json();
          setEmployerDetails(data);
        }
      } catch (err) {
        console.error("Error fetching employer details:", err);
        setError(t("employerDetails.error.loadFailed"));
        setEmployerDetails(null);
      } finally {
        setLoading(false);
      }
    };

    if (companyDomain) {
      fetchEmployerDetails();
    }
  }, [companyDomain, t]);

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
        <Paper
          sx={{
            p: 2,
            mb: 2,
            bgcolor: "error.light",
            mx: "auto",
            mt: 4,
            maxWidth: 800,
          }}
        >
          <Typography color="error.main" align="center">
            {error}
          </Typography>
          <Box sx={{ display: "flex", justifyContent: "center", mt: 2 }}>
            <Button variant="outlined" onClick={() => router.back()}>
              {t("common.back")}
            </Button>
          </Box>
        </Paper>
      </AuthenticatedLayout>
    );
  }

  if (!employerDetails) {
    return (
      <AuthenticatedLayout>
        <Paper sx={{ p: 2, mb: 2, mx: "auto", mt: 4, maxWidth: 800 }}>
          <Typography align="center">
            {t("employerDetails.notFound")}
          </Typography>
          <Box sx={{ display: "flex", justifyContent: "center", mt: 2 }}>
            <Button variant="outlined" onClick={() => router.back()}>
              {t("common.back")}
            </Button>
          </Box>
        </Paper>
      </AuthenticatedLayout>
    );
  }

  return (
    <AuthenticatedLayout>
      <Box sx={{ maxWidth: 800, mx: "auto", mt: 4 }}>
        <Paper sx={{ p: 4 }}>
          <Typography variant="h4" gutterBottom>
            {employerDetails.name}
          </Typography>
          <Typography variant="body1" sx={{ mb: 1 }}>
            {t("employerDetails.verifiedEmployees")}:{" "}
            {employerDetails.verified_employees_count}
          </Typography>
          <Typography variant="body1">
            {t("employerDetails.activeOpenings")}:{" "}
            {employerDetails.active_openings_count}
          </Typography>
        </Paper>
      </Box>
    </AuthenticatedLayout>
  );
}
