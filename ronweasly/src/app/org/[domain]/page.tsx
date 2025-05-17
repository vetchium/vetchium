"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import Alert from "@mui/material/Alert";
import AlertTitle from "@mui/material/AlertTitle";
import Box from "@mui/material/Box";
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
        <Box sx={{ maxWidth: 800, mx: "auto", mt: 4, p: 2 }}>
          <Alert severity="error">
            <AlertTitle>{t("common.error.invalidParams")}</AlertTitle>
            {t("common.error.invalidParamsDetail")}
          </Alert>
        </Box>
      </AuthenticatedLayout>
    );
  }

  const companyDomain = params.domain as string;

  useEffect(() => {
    const fetchEmployerDetails = async () => {
      setLoading(true);
      setError(null);
      setEmployerDetails(null);
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
            router.push("/login");
            return;
          } else if (response.status === 404) {
            // For 404, we don't set an error message here.
            // employerDetails will remain null (set at the start of fetchEmployerDetails),
            // and the specific "if (!employerDetails)" block will handle the UI.
          } else {
            // For other HTTP errors (500, etc.)
            setError(t("employerDetails.error.loadFailed"));
          }
          // setEmployerDetails(null) is already done at the start of the function or if error occurs
        } else {
          const data: HubEmployerDetails = await response.json();
          setEmployerDetails(data);
          // setError(null) is already done at the start of the function
        }
      } catch (err) {
        console.error("Error fetching employer details:", err);
        setError(t("employerDetails.error.loadFailed"));
        // setEmployerDetails(null) is already done at the start
      } finally {
        setLoading(false);
      }
    };

    if (companyDomain) {
      fetchEmployerDetails();
    }
  }, [companyDomain, t, router]);

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
        <Box sx={{ maxWidth: 800, mx: "auto", mt: 4, p: 2 }}>
          <Alert severity="error">
            <AlertTitle>
              {error === t("common.error.notAuthenticated") ||
              error === t("common.error.sessionExpired")
                ? t("common.error.authenticationNeededTitle") // Assuming a new title for auth issues
                : t("common.error.serverError")}
            </AlertTitle>
            {error}
          </Alert>
        </Box>
      </AuthenticatedLayout>
    );
  }

  if (!employerDetails) {
    return (
      <AuthenticatedLayout>
        <Box sx={{ maxWidth: 800, mx: "auto", mt: 4, p: 2 }}>
          <Alert severity="warning">
            <AlertTitle>{t("employerDetails.notFoundTitle")}</AlertTitle>
            {t("employerDetails.notFound")}
          </Alert>
        </Box>
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
