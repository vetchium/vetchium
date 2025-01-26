"use client";

import { useParams, useRouter } from "next/navigation";
import { useState, useEffect } from "react";
import {
  Box,
  Typography,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  TextField,
  Alert,
  CircularProgress,
} from "@mui/material";
import { useTranslation } from "@/hooks/useTranslation";
import { config } from "@/config";
import Cookies from "js-cookie";
import {
  Assessment,
  GetAssessmentRequest,
  InterviewersDecision,
} from "@psankar/vetchi-typespec";

export default function InterviewDetailPage() {
  const params = useParams();
  const interviewId = params.id as string;
  const { t, tObject } = useTranslation();
  const router = useRouter();

  const [assessment, setAssessment] = useState<Assessment | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  useEffect(() => {
    fetchAssessment();
  }, [interviewId]);

  const fetchAssessment = async () => {
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/get-assessment`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            interview_id: interviewId,
          } satisfies GetAssessmentRequest),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/signin");
        return;
      }

      if (!response.ok) {
        throw new Error(t("interviews.assessment.fetchError"));
      }

      const data = await response.json();
      setAssessment(data || { interview_id: interviewId });
    } catch (err) {
      setError(err instanceof Error ? err.message : t("common.error"));
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    if (!assessment) return;

    try {
      setSaving(true);
      setError(null);
      setSuccessMessage(null);

      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/put-assessment`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(assessment),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/signin");
        return;
      }

      if (!response.ok) {
        throw new Error(t("interviews.assessment.saveError"));
      }

      setSuccessMessage(t("interviews.assessment.saveSuccess"));
    } catch (err) {
      setError(err instanceof Error ? err.message : t("common.error"));
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", p: 3 }}>
        <CircularProgress />
      </Box>
    );
  }

  const ratings = tObject("interviews.assessment.ratings");

  return (
    <Box sx={{ p: 3 }}>
      <Box sx={{ display: "flex", justifyContent: "space-between", mb: 3 }}>
        <Typography variant="h4">{t("interviews.manageInterview")}</Typography>
        <Button variant="outlined" onClick={() => router.back()}>
          {t("common.back")}
        </Button>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      {successMessage && (
        <Alert
          severity="success"
          sx={{ mb: 2 }}
          onClose={() => setSuccessMessage(null)}
        >
          {successMessage}
        </Alert>
      )}

      <Box sx={{ maxWidth: 600 }}>
        <Typography variant="h5" sx={{ mb: 3 }}>
          {t("interviews.assessment.title")}
        </Typography>

        <FormControl fullWidth sx={{ mb: 3 }}>
          <InputLabel id="rating-label">
            {t("interviews.assessment.rating")}
          </InputLabel>
          <Select
            labelId="rating-label"
            value={assessment?.decision || "NEUTRAL"}
            label={t("interviews.assessment.rating")}
            onChange={(e) =>
              setAssessment((prev) => ({
                ...prev!,
                decision: e.target.value as InterviewersDecision,
              }))
            }
          >
            {Object.entries(ratings).map(([key, label]) => (
              <MenuItem key={key} value={key}>
                {label}
              </MenuItem>
            ))}
          </Select>
        </FormControl>

        <TextField
          fullWidth
          multiline
          rows={4}
          label={t("interviews.assessment.feedback")}
          placeholder={t("interviews.assessment.feedbackPlaceholder")}
          value={assessment?.feedback_to_candidate || ""}
          onChange={(e) =>
            setAssessment((prev) => ({
              ...prev!,
              feedback_to_candidate: e.target.value,
            }))
          }
          sx={{ mb: 3 }}
        />

        <Button
          variant="contained"
          onClick={handleSave}
          disabled={saving}
          sx={{ minWidth: 120 }}
        >
          {saving ? (
            <CircularProgress size={24} color="inherit" />
          ) : (
            t("interviews.assessment.save")
          )}
        </Button>
      </Box>
    </Box>
  );
}
