import React, { useState } from "react";
import { useRouter } from "next/router";
import {
  Container,
  Paper,
  Typography,
  TextField,
  Button,
  FormControl,
  Box,
  Alert,
  CircularProgress,
  RadioGroup,
  FormControlLabel,
  Radio,
} from "@mui/material";
import { config } from "../../config";
import { CountrySelect } from "../../components/CountrySelect";
import { useTranslation } from "@/hooks/useTranslation";

interface OnboardFormData {
  full_name: string;
  resident_country_code: string;
  password: string;
  selected_tier: "FREE_TIER" | "PAID_TIER";
  short_bio?: string;
  long_bio?: string;
}

interface OnboardResponse {
  session_token: string;
  generated_handle: string;
}

export default function SignupHubUser() {
  const router = useRouter();
  const { token } = router.query;
  const { t } = useTranslation();
  const [existingSession, setExistingSession] = useState<boolean>(false);

  React.useEffect(() => {
    const sessionToken = localStorage.getItem("sessionToken");
    if (sessionToken) {
      setExistingSession(true);
    }
  }, []);

  const [formData, setFormData] = useState<OnboardFormData>({
    full_name: "",
    resident_country_code: "",
    password: "",
    selected_tier: "FREE_TIER",
  });

  const [errors, setErrors] = useState<
    Partial<Record<keyof OnboardFormData, string>>
  >({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [apiError, setApiError] = useState<string | null>(null);
  const [success, setSuccess] = useState<OnboardResponse | null>(null);

  const validateForm = (): boolean => {
    const newErrors: Partial<Record<keyof OnboardFormData, string>> = {};

    if (!formData.full_name.trim()) {
      newErrors.full_name = t("hubUserOnboarding.error.requiredField");
    }

    if (!formData.resident_country_code) {
      newErrors.resident_country_code = t(
        "hubUserOnboarding.error.requiredField"
      );
    }

    if (!formData.password || formData.password.length < 8) {
      newErrors.password = t("hubUserOnboarding.error.passwordTooShort");
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm() || !token || typeof token !== "string") {
      return;
    }

    setIsSubmitting(true);
    setApiError(null);

    try {
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/onboard-user`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            ...formData,
            token,
            preferred_language: "en", // Hardcoded value
          }),
        }
      );

      if (!response.ok) {
        throw new Error(t("hubUserOnboarding.error.onboardingFailed"));
      }

      const data: OnboardResponse = await response.json();
      setSuccess(data);

      // Store the session token
      localStorage.setItem("sessionToken", data.session_token);

      // Redirect after 3 seconds
      setTimeout(() => {
        router.push("/dashboard");
      }, 3000);
    } catch (error) {
      setApiError(
        error instanceof Error
          ? error.message
          : t("hubUserOnboarding.error.onboardingFailed")
      );
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleInputChange =
    (field: keyof OnboardFormData) =>
    (
      e: React.ChangeEvent<HTMLInputElement | { name?: string; value: unknown }>
    ) => {
      setFormData((prev) => ({
        ...prev,
        [field]: e.target.value,
      }));
      // Clear error when user starts typing
      if (errors[field]) {
        setErrors((prev) => ({
          ...prev,
          [field]: undefined,
        }));
      }
    };

  if (existingSession) {
    return (
      <Container maxWidth="sm">
        <Paper elevation={3} sx={{ p: 4, mt: 4 }}>
          <Alert severity="warning" sx={{ mb: 2 }}>
            <Typography variant="h6" gutterBottom>
              Another User is Already Signed In
            </Typography>
            <Typography>
              Please sign out of the current account before creating a new one.
            </Typography>
          </Alert>
          <Button
            variant="contained"
            color="primary"
            fullWidth
            onClick={() => {
              localStorage.removeItem("sessionToken");
              window.location.reload();
            }}
          >
            Sign Out
          </Button>
        </Paper>
      </Container>
    );
  }

  if (success) {
    return (
      <Container maxWidth="sm">
        <Paper elevation={3} sx={{ p: 4, mt: 4 }}>
          <Typography variant="h5" gutterBottom>
            {t("hubUserOnboarding.success.title")}
          </Typography>
          <Typography paragraph>
            {t("hubUserOnboarding.success.description")}
          </Typography>
          <Typography paragraph>
            {t("hubUserOnboarding.success.handle", {
              handle: success.generated_handle,
            })}
          </Typography>
          <Typography color="textSecondary">
            {t("hubUserOnboarding.success.redirecting")}
          </Typography>
          <CircularProgress sx={{ mt: 2 }} />
        </Paper>
      </Container>
    );
  }

  return (
    <Container maxWidth="sm">
      <Paper elevation={3} sx={{ p: 4, mt: 4 }}>
        <Typography variant="h4" gutterBottom>
          {t("hubUserOnboarding.title")}
        </Typography>
        <Typography variant="subtitle1" gutterBottom>
          {t("hubUserOnboarding.subtitle")}
        </Typography>

        {apiError && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {apiError}
          </Alert>
        )}

        <form onSubmit={handleSubmit}>
          <TextField
            fullWidth
            label={t("hubUserOnboarding.form.fullName")}
            placeholder={t("hubUserOnboarding.form.fullNamePlaceholder")}
            value={formData.full_name}
            onChange={handleInputChange("full_name")}
            error={!!errors.full_name}
            helperText={errors.full_name}
            margin="normal"
            required
          />

          <TextField
            fullWidth
            type="password"
            label={t("hubUserOnboarding.form.password")}
            placeholder={t("hubUserOnboarding.form.passwordPlaceholder")}
            value={formData.password}
            onChange={handleInputChange("password")}
            error={!!errors.password}
            helperText={errors.password}
            margin="normal"
            required
          />

          <CountrySelect
            value={formData.resident_country_code}
            onChange={(value: string) => {
              setFormData((prev) => ({
                ...prev,
                resident_country_code: value,
              }));
            }}
            error={!!errors.resident_country_code}
            helperText={errors.resident_country_code}
          />

          <FormControl component="fieldset" margin="normal" fullWidth>
            <Typography variant="subtitle2" gutterBottom>
              {t("hubUserOnboarding.form.tier.label")}
            </Typography>
            <RadioGroup
              value={formData.selected_tier}
              onChange={handleInputChange("selected_tier")}
            >
              <FormControlLabel
                value="FREE_TIER"
                control={<Radio />}
                label={
                  <Box>
                    <Typography variant="body1">
                      {t("hubUserOnboarding.form.tier.free")}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      {t("hubUserOnboarding.form.tier.freeDescription")}
                    </Typography>
                  </Box>
                }
              />
              <FormControlLabel
                value="PAID_TIER"
                control={<Radio />}
                label={
                  <Box>
                    <Typography variant="body1">
                      {t("hubUserOnboarding.form.tier.paid")}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      {t("hubUserOnboarding.form.tier.paidDescription")}
                    </Typography>
                  </Box>
                }
              />
            </RadioGroup>
          </FormControl>

          <TextField
            fullWidth
            multiline
            rows={2}
            label={t("hubUserOnboarding.form.shortBio")}
            placeholder={t("hubUserOnboarding.form.shortBioPlaceholder")}
            value={formData.short_bio || ""}
            onChange={handleInputChange("short_bio")}
            margin="normal"
            required
            error={!!errors.short_bio}
            helperText={
              errors.short_bio ||
              `${(formData.short_bio || "").length}/64 characters`
            }
            inputProps={{ maxLength: 64 }}
          />

          <TextField
            fullWidth
            multiline
            rows={4}
            label={t("hubUserOnboarding.form.longBio")}
            placeholder={t("hubUserOnboarding.form.longBioPlaceholder")}
            value={formData.long_bio || ""}
            onChange={handleInputChange("long_bio")}
            margin="normal"
            required
            error={!!errors.long_bio}
            helperText={
              errors.long_bio ||
              `${(formData.long_bio || "").length}/1024 characters`
            }
            inputProps={{ maxLength: 1024 }}
          />

          <Button
            type="submit"
            variant="contained"
            color="primary"
            fullWidth
            size="large"
            disabled={isSubmitting}
            sx={{ mt: 3 }}
          >
            {isSubmitting ? (
              <CircularProgress size={24} />
            ) : (
              t("hubUserOnboarding.form.submit")
            )}
          </Button>
        </form>
      </Paper>
    </Container>
  );
}
