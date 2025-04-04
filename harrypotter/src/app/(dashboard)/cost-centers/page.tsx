"use client";

import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import { Delete as DeleteIcon, Edit as EditIcon } from "@mui/icons-material";
import {
  Alert,
  Box,
  Button,
  CircularProgress,
  Container,
  FormControlLabel,
  Grid,
  IconButton,
  Paper,
  Switch,
  Typography,
} from "@mui/material";
import {
  CostCenter,
  DefunctCostCenterRequest,
  GetCostCentersRequest,
} from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useMemo, useState } from "react";

export default function CostCentersPage() {
  const [costCenters, setCostCenters] = useState<CostCenter[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [includeDefunct, setIncludeDefunct] = useState(() => {
    if (typeof window === "undefined") return false;
    return localStorage.getItem("includeDefunctCostCenters") === "true";
  });
  const [isLoading, setIsLoading] = useState(false);

  const router = useRouter();
  const { t } = useTranslation();

  // Memoize error messages
  const errorMessages = useMemo(
    () => ({
      unauthorized: t("auth.unauthorized"),
      commonError: t("common.error"),
    }),
    [t]
  );

  const fetchCostCenters = useCallback(async () => {
    try {
      setIsLoading(true);
      const sessionToken = Cookies.get("session_token");
      if (!sessionToken) {
        setError(errorMessages.unauthorized);
        setIsLoading(false);
        return;
      }

      const request: GetCostCentersRequest = {
        states: includeDefunct ? ["ACTIVE_CC", "DEFUNCT_CC"] : ["ACTIVE_CC"],
      };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/get-cost-centers`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${sessionToken}`,
          },
          body: JSON.stringify(request),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/signin");
        return;
      }

      if (!response.ok) {
        throw new Error(errorMessages.commonError);
      }

      const data = await response.json();
      setCostCenters(data || []);
    } catch {
      setError(errorMessages.commonError);
    } finally {
      setIsLoading(false);
    }
  }, [includeDefunct, router, errorMessages]);

  useEffect(() => {
    fetchCostCenters();
  }, [fetchCostCenters]);

  const handleDelete = async (costCenter: CostCenter) => {
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const request: DefunctCostCenterRequest = {
        name: costCenter.name,
      };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/defunct-cost-center`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(request),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/signin");
        return;
      }

      if (!response.ok) {
        throw new Error(t("costCenters.defunctError"));
      }

      fetchCostCenters();
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("costCenters.defunctError")
      );
    }
  };

  return (
    <Container maxWidth="lg">
      <Box
        sx={{
          mb: 4,
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
        }}
      >
        <Typography variant="h4" component="h1" sx={{ color: "text.primary" }}>
          {t("costCenters.title")}
        </Typography>
        <Box sx={{ display: "flex", gap: 2, alignItems: "center" }}>
          <FormControlLabel
            control={
              <Switch
                checked={includeDefunct}
                onChange={(e) => {
                  setIncludeDefunct(e.target.checked);
                  localStorage.setItem(
                    "includeDefunctCostCenters",
                    e.target.checked.toString()
                  );
                }}
              />
            }
            label={t("costCenters.includeDefunct")}
            sx={{ color: "text.primary" }}
          />
          <Button
            variant="contained"
            onClick={() => router.push("/cost-centers/add")}
          >
            {t("costCenters.add")}
          </Button>
        </Box>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      {isLoading ? (
        <Box sx={{ display: "flex", justifyContent: "center", p: 4 }}>
          <CircularProgress />
        </Box>
      ) : costCenters.length > 0 ? (
        <Grid container spacing={2}>
          {costCenters.map((costCenter) => (
            <Grid item xs={12} md={6} lg={4} key={costCenter.name}>
              <Paper
                sx={{
                  p: 2,
                  height: "100%",
                  opacity: costCenter.state === "DEFUNCT_CC" ? 0.7 : 1,
                  color: "text.primary",
                }}
              >
                <Box
                  sx={{
                    display: "flex",
                    justifyContent: "space-between",
                    alignItems: "flex-start",
                    mb: 2,
                  }}
                >
                  <Box>
                    <Typography
                      variant="h6"
                      gutterBottom
                      sx={{ color: "text.primary" }}
                    >
                      {costCenter.name}
                    </Typography>
                    <Typography
                      variant="caption"
                      sx={{
                        backgroundColor:
                          costCenter.state === "DEFUNCT_CC"
                            ? "error.main"
                            : "success.main",
                        color: "white",
                        px: 1,
                        py: 0.5,
                        borderRadius: 1,
                      }}
                    >
                      {costCenter.state === "DEFUNCT_CC"
                        ? t("costCenters.defunct")
                        : t("costCenters.active")}
                    </Typography>
                  </Box>
                  <Box>
                    <IconButton
                      onClick={() =>
                        router.push(
                          `/cost-centers/edit?name=${encodeURIComponent(
                            costCenter.name
                          )}`
                        )
                      }
                      size="small"
                      sx={{ color: "text.primary" }}
                    >
                      <EditIcon />
                    </IconButton>
                    <IconButton
                      onClick={() => handleDelete(costCenter)}
                      size="small"
                      disabled={costCenter.state === "DEFUNCT_CC"}
                      sx={{ color: "text.primary" }}
                    >
                      <DeleteIcon />
                    </IconButton>
                  </Box>
                </Box>

                {costCenter.notes && (
                  <Box sx={{ mb: 1 }}>
                    <Typography
                      variant="subtitle2"
                      sx={{ color: "text.secondary", mb: 0.5 }}
                    >
                      {t("costCenters.notes")}
                    </Typography>
                    <Typography
                      sx={{ color: "text.primary", whiteSpace: "pre-line" }}
                    >
                      {costCenter.notes}
                    </Typography>
                  </Box>
                )}
              </Paper>
            </Grid>
          ))}
        </Grid>
      ) : (
        <Paper sx={{ p: 4, textAlign: "center" }}>
          <Typography variant="body1" sx={{ color: "text.secondary" }}>
            {t("costCenters.noCostCenters")}
          </Typography>
        </Paper>
      )}
    </Container>
  );
}
