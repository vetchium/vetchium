"use client";

import { useParams, useRouter } from "next/navigation";
import { useState, useEffect } from "react";
import {
  Box,
  Button,
  Container,
  TextField,
  Typography,
  Alert,
  Paper,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Stack,
  FormControlLabel,
  Switch,
} from "@mui/material";
import { useTranslation } from "@/hooks/useTranslation";
import {
  InterviewType,
  InterviewTypes,
  TimeZone,
  validTimezones,
  AddInterviewRequest,
} from "@psankar/vetchi-typespec";
import { config } from "@/config";
import Cookies from "js-cookie";
import { DateTimePicker } from "@mui/x-date-pickers/DateTimePicker";
import { LocalizationProvider } from "@mui/x-date-pickers/LocalizationProvider";
import { AdapterDateFns } from "@mui/x-date-pickers/AdapterDateFns";

export default function AddInterviewPage() {
  const params = useParams();
  const candidacyId = params.id as string;
  const { t } = useTranslation();
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  // Get user's timezone and find closest matching TimeZone enum value
  const userTimezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
  const userOffset = new Date().getTimezoneOffset();
  const offsetHours = Math.floor(Math.abs(userOffset) / 60);
  const offsetMinutes = Math.abs(userOffset) % 60;
  const offsetStr = `${userOffset <= 0 ? "+" : "-"}${offsetHours
    .toString()
    .padStart(2, "0")}${offsetMinutes.toString().padStart(2, "0")}`;

  // Find the closest matching timezone from validTimezones
  const defaultTimezone =
    Array.from(validTimezones).find((tz) => tz.includes(`GMT${offsetStr}`)) ||
    "UTC Coordinated Universal Time GMT+0000";

  const [interview, setInterview] = useState<{
    startTime: string;
    endTime: string;
    type: InterviewType;
    description: string;
    timezone: TimeZone;
  }>({
    startTime: "",
    endTime: "",
    type: InterviewTypes.VIDEO_CALL,
    description: "",
    timezone: defaultTimezone,
  });

  const [allowPastDates, setAllowPastDates] = useState(false);
  const [use24HourFormat, setUse24HourFormat] = useState(true);

  // Handle localStorage in useEffect to avoid SSR issues
  useEffect(() => {
    const saved = localStorage.getItem("create_interview_24hour_format");
    if (saved !== null) {
      setUse24HourFormat(saved === "true");
    }
  }, []);

  // Update localStorage when time format preference changes
  useEffect(() => {
    localStorage.setItem(
      "create_interview_24hour_format",
      use24HourFormat.toString()
    );
  }, [use24HourFormat]);

  const handleAddInterview = async () => {
    try {
      setIsLoading(true);
      setError(null);

      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      // Convert local times to UTC for API
      const startDate = new Date(interview.startTime);
      const endDate = new Date(interview.endTime);

      // Get timezone offset from selected timezone
      const tzMatch = interview.timezone.match(/GMT([+-]\d{4})/);
      const tzOffset = tzMatch ? tzMatch[1] : "+0000";
      const tzHours = parseInt(tzOffset.slice(1, 3));
      const tzMinutes = parseInt(tzOffset.slice(3));
      const offsetMillis =
        (tzHours * 60 + tzMinutes) * 60 * 1000 * (tzOffset[0] === "+" ? -1 : 1);

      // Adjust dates to UTC
      const utcStartDate = new Date(startDate.getTime() + offsetMillis);
      const utcEndDate = new Date(endDate.getTime() + offsetMillis);

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/add-interview`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            candidacy_id: candidacyId,
            start_time: utcStartDate,
            end_time: utcEndDate,
            interview_type: interview.type,
            description: interview.description,
          } satisfies AddInterviewRequest),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/signin");
        return;
      }

      if (!response.ok) throw new Error(t("interviews.addError"));

      router.push(`/candidacy/${candidacyId}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : t("common.serverError"));
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Container maxWidth="md">
      <Paper sx={{ p: 4 }}>
        <Box sx={{ mb: 4 }}>
          <Typography variant="h4" component="h1" gutterBottom>
            {t("interviews.addNew")}
          </Typography>
          {error && (
            <Alert severity="error" sx={{ mt: 2 }}>
              {error}
            </Alert>
          )}
        </Box>

        <Box component="form" noValidate sx={{ mt: 1 }}>
          <Stack spacing={3}>
            <FormControl fullWidth>
              <InputLabel>{t("interviews.type")}</InputLabel>
              <Select
                value={interview.type}
                label={t("interviews.type")}
                onChange={(e) =>
                  setInterview({
                    ...interview,
                    type: e.target.value as InterviewType,
                  })
                }
              >
                <MenuItem value={InterviewTypes.VIDEO_CALL}>
                  {t("interviews.types.VIDEO_CALL")}
                </MenuItem>
                <MenuItem value={InterviewTypes.IN_PERSON}>
                  {t("interviews.types.IN_PERSON")}
                </MenuItem>
                <MenuItem value={InterviewTypes.TAKE_HOME}>
                  {t("interviews.types.TAKE_HOME")}
                </MenuItem>
              </Select>
            </FormControl>

            <FormControl fullWidth>
              <InputLabel>{t("interviews.timezone")}</InputLabel>
              <Select
                value={interview.timezone}
                label={t("interviews.timezone")}
                onChange={(e) =>
                  setInterview({
                    ...interview,
                    timezone: e.target.value as TimeZone,
                  })
                }
              >
                {Array.from(validTimezones).map((tz) => (
                  <MenuItem key={tz} value={tz}>
                    {tz}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>

            <FormControlLabel
              control={
                <Switch
                  checked={allowPastDates}
                  onChange={(e) => setAllowPastDates(e.target.checked)}
                />
              }
              label={t("interviews.allowPastDates")}
            />

            <FormControlLabel
              control={
                <Switch
                  checked={use24HourFormat}
                  onChange={(e) => setUse24HourFormat(e.target.checked)}
                />
              }
              label={t("interviews.use24HourFormat")}
            />

            <LocalizationProvider dateAdapter={AdapterDateFns}>
              <DateTimePicker
                label={t("interviews.startTime")}
                value={
                  interview.startTime ? new Date(interview.startTime) : null
                }
                onChange={(newValue: Date | null) => {
                  if (newValue) {
                    const startDate = newValue;
                    // Set end time to 1 hour after start time
                    const endDate = new Date(startDate);
                    endDate.setHours(startDate.getHours() + 1);

                    setInterview({
                      ...interview,
                      startTime: startDate.toISOString(),
                      endTime: endDate.toISOString(),
                    });
                  }
                }}
                views={["year", "month", "day", "hours", "minutes"]}
                ampm={!use24HourFormat}
                format={
                  use24HourFormat
                    ? "MMMM dd, yyyy HH:mm"
                    : "MMMM dd, yyyy hh:mm a"
                }
                minDateTime={allowPastDates ? undefined : new Date()}
                slotProps={{
                  textField: {
                    fullWidth: true,
                  },
                }}
              />

              <DateTimePicker
                label={t("interviews.endTime")}
                value={interview.endTime ? new Date(interview.endTime) : null}
                onChange={(newValue: Date | null) => {
                  if (newValue) {
                    // Only update if end time is after start time
                    if (
                      interview.startTime &&
                      new Date(newValue) <= new Date(interview.startTime)
                    ) {
                      setError(t("interviews.endTimeBeforeStart"));
                      return;
                    }
                    setInterview({
                      ...interview,
                      endTime: newValue.toISOString(),
                    });
                  }
                }}
                views={["year", "month", "day", "hours", "minutes"]}
                ampm={!use24HourFormat}
                format={
                  use24HourFormat
                    ? "MMMM dd, yyyy HH:mm"
                    : "MMMM dd, yyyy hh:mm a"
                }
                minDateTime={
                  interview.startTime
                    ? new Date(interview.startTime)
                    : undefined
                }
                slotProps={{
                  textField: {
                    fullWidth: true,
                  },
                }}
              />
            </LocalizationProvider>

            <TextField
              label={t("interviews.description")}
              multiline
              rows={4}
              value={interview.description}
              onChange={(e) =>
                setInterview({
                  ...interview,
                  description: e.target.value,
                })
              }
              fullWidth
            />

            <Box sx={{ mt: 4, display: "flex", gap: 2 }}>
              <Button
                variant="outlined"
                onClick={() => router.push(`/candidacy/${candidacyId}`)}
              >
                {t("common.cancel")}
              </Button>
              <Button
                variant="contained"
                onClick={handleAddInterview}
                disabled={
                  isLoading ||
                  !interview.startTime ||
                  !interview.endTime ||
                  !interview.description
                }
              >
                {isLoading ? t("common.loading") : t("common.save")}
              </Button>
            </Box>
          </Stack>
        </Box>
      </Paper>
    </Container>
  );
}
