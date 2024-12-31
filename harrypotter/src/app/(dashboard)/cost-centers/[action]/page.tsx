"use client";

import {
  Box,
  Button,
  Container,
  TextField,
  Typography,
  Alert,
  Paper,
} from "@mui/material";
import { useState, useEffect } from "react";
import { useRouter, useParams } from "next/navigation";
import { useTranslation } from "@/hooks/useTranslation";
import { config } from "@/config";
import Cookies from "js-cookie";
import {
  CostCenter,
  GetCostCenterRequest,
  AddCostCenterRequest,
  UpdateCostCenterRequest,
} from "@psankar/vetchi-typespec/employer/costcenters";

export default function CostCenterActionPage() {
  const [name, setName] = useState("");
  const [notes, setNotes] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [initialCostCenter, setInitialCostCenter] = useState<CostCenter | null>(
    null
  );

  const router = useRouter();
  const params = useParams();
  const { t } = useTranslation();
  const isEdit = params.action === "edit";

  useEffect(() => {
    if (isEdit) {
      const costCenterName = window.location.search.split("name=")[1];
      if (costCenterName) {
        fetchCostCenter(decodeURIComponent(costCenterName));
      } else {
        router.push("/cost-centers");
      }
    }
  }, [isEdit]);

  const fetchCostCenter = async (costCenterName: string) => {
    try {
      setIsLoading(true);
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const request: GetCostCenterRequest = {
        name: costCenterName,
      };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/get-cost-center`,
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
        throw new Error(t("costCenters.fetchError"));
      }

      const costCenter = await response.json();
      setInitialCostCenter(costCenter);
      setName(costCenter.name);
      setNotes(costCenter.notes || "");
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("costCenters.fetchError")
      );
    } finally {
      setIsLoading(false);
    }
  };

  const handleSave = async () => {
    try {
      setError(null);
      setIsLoading(true);
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const url = isEdit
        ? `${config.API_SERVER_PREFIX}/employer/update-cost-center`
        : `${config.API_SERVER_PREFIX}/employer/add-cost-center`;

      const requestBody: AddCostCenterRequest | UpdateCostCenterRequest = {
        name,
        notes: notes || undefined,
      };

      const response = await fetch(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(requestBody),
      });

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/signin");
        return;
      }

      if (!response.ok) {
        const data = await response.json();
        throw new Error(
          data.message ||
            (isEdit ? t("costCenters.updateError") : t("costCenters.addError"))
        );
      }

      router.push("/cost-centers");
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    } finally {
      setIsLoading(false);
    }
  };

  if (isLoading && isEdit) {
    return (
      <Container maxWidth="md">
        <Typography>{t("common.loading")}</Typography>
      </Container>
    );
  }

  return (
    <Container maxWidth="md">
      <Paper sx={{ p: 4 }}>
        <Box sx={{ mb: 4 }}>
          <Typography variant="h4" component="h1" gutterBottom>
            {isEdit ? t("costCenters.editTitle") : t("costCenters.addTitle")}
          </Typography>
          {error && (
            <Alert severity="error" sx={{ mt: 2 }}>
              {error}
            </Alert>
          )}
        </Box>

        <Box component="form" noValidate sx={{ mt: 1 }}>
          <TextField
            margin="normal"
            required
            fullWidth
            label={t("costCenters.name")}
            value={name}
            onChange={(e) => setName(e.target.value)}
            disabled={isEdit}
          />

          <TextField
            margin="normal"
            fullWidth
            multiline
            rows={4}
            label={t("costCenters.notes")}
            value={notes}
            onChange={(e) => setNotes(e.target.value)}
          />

          <Box sx={{ mt: 4, display: "flex", gap: 2 }}>
            <Button
              variant="outlined"
              onClick={() => router.push("/cost-centers")}
            >
              {t("common.cancel")}
            </Button>
            <Button
              variant="contained"
              onClick={handleSave}
              disabled={isLoading || !name}
            >
              {isLoading ? t("common.loading") : t("common.save")}
            </Button>
          </Box>
        </Box>
      </Paper>
    </Container>
  );
}
