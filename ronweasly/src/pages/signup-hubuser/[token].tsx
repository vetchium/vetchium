import { useTranslation } from "@/hooks/useTranslation";
import { Visibility, VisibilityOff } from "@mui/icons-material";
import {
  Alert,
  Box,
  Button,
  Card,
  CircularProgress,
  Container,
  IconButton,
  InputAdornment,
  Paper,
  TextField,
  Typography,
} from "@mui/material";
import type {
  OnboardHubUserRequest,
  OnboardHubUserResponse,
} from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useRouter } from "next/router";
import React, { useState } from "react";
import { CountrySelect } from "../../components/CountrySelect";
import { config } from "../../config";

interface FormData extends Omit<OnboardHubUserRequest, "token"> {
  confirm_password: string;
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

  const [formData, setFormData] = useState<FormData>({
    full_name: "",
    resident_country_code: "",
    password: "",
    confirm_password: "",
    selected_tier: "PAID_HUB_USER",
  });

  const [errors, setErrors] = useState<Partial<Record<keyof FormData, string>>>(
    {}
  );
  const [touched, setTouched] = useState<
    Partial<Record<keyof FormData, boolean>>
  >({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [apiError, setApiError] = useState<string | null>(null);
  const [success, setSuccess] = useState<OnboardHubUserResponse | null>(null);
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  const validateForm = (fieldsToValidate?: Array<keyof FormData>): boolean => {
    const newErrors: Partial<Record<keyof FormData, string>> = {};

    // If specific fields are provided, only validate those
    // Otherwise, validate only touched fields
    const fields =
      fieldsToValidate ||
      (Object.keys(touched).filter(
        (key) => touched[key as keyof FormData]
      ) as Array<keyof FormData>);

    // If no fields to validate, return true
    if (fields.length === 0) {
      return true;
    }

    // Only validate fields that should be validated
    if (fields.includes("full_name") && !formData.full_name.trim()) {
      newErrors.full_name = t("hubUserOnboarding.error.requiredField");
    }

    if (
      fields.includes("resident_country_code") &&
      !formData.resident_country_code
    ) {
      newErrors.resident_country_code = t(
        "hubUserOnboarding.error.requiredField"
      );
    }

    if (
      fields.includes("password") &&
      (!formData.password ||
        formData.password.length < 12 ||
        formData.password.length > 64)
    ) {
      newErrors.password = t("hubUserOnboarding.error.passwordLength");
    }

    if (
      fields.includes("confirm_password") &&
      formData.password !== formData.confirm_password
    ) {
      newErrors.confirm_password = t(
        "hubUserOnboarding.error.passwordMismatch"
      );
    }

    if (fields.includes("short_bio") && !formData.short_bio?.trim()) {
      newErrors.short_bio = t("hubUserOnboarding.error.requiredField");
    }

    if (fields.includes("long_bio") && !formData.long_bio?.trim()) {
      newErrors.long_bio = t("hubUserOnboarding.error.requiredField");
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const isFormValid = (): boolean => {
    // Check for required fields
    if (
      !formData.full_name.trim() ||
      !formData.resident_country_code ||
      !formData.password ||
      !formData.confirm_password ||
      !formData.short_bio?.trim() ||
      !formData.long_bio?.trim()
    ) {
      return false;
    }

    // Check password length
    if (formData.password.length < 12 || formData.password.length > 64) {
      return false;
    }

    // Check if passwords match
    if (formData.password !== formData.confirm_password) {
      return false;
    }

    return true;
  };

  // Add effect to validate form whenever form data changes
  // But only validate fields that have been touched
  React.useEffect(() => {
    validateForm();
  }, [formData]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // Mark all fields as touched for final validation
    const allFields = Object.keys(formData) as Array<keyof FormData>;
    setTouched(
      allFields.reduce((acc, field) => ({ ...acc, [field]: true }), {})
    );

    // Validate all fields regardless of touch state
    if (!validateForm(allFields) || !token || typeof token !== "string") {
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
        if (response.status === 404) {
          throw new Error(t("hubUserOnboarding.error.invalidToken"));
        }

        if (response.status === 400) {
          const errorData = await response.json();
          throw new Error(
            t("hubUserOnboarding.error.validationError", {
              details: errorData.message,
            })
          );
        }

        throw new Error(t("hubUserOnboarding.error.onboardingFailed"));
      }

      const data: OnboardHubUserResponse = await response.json();
      setSuccess(data);

      // Store the session token in a cookie
      Cookies.set("session_token", data.session_token, { path: "/" });

      // Redirect after 3 seconds
      setTimeout(() => {
        // TODO: If the account type is PAID_HUB_USER, redirect to the payments page
        router.push("/");
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
    (field: keyof FormData) =>
    (
      e: React.ChangeEvent<HTMLInputElement | { name?: string; value: unknown }>
    ) => {
      // Mark field as touched
      if (!touched[field]) {
        setTouched((prev) => ({
          ...prev,
          [field]: true,
        }));
      }

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

  const handleConfirmPasswordBlur = () => {
    // Mark confirm_password as touched
    setTouched((prev) => ({
      ...prev,
      confirm_password: true,
    }));

    if (
      formData.confirm_password &&
      formData.password !== formData.confirm_password
    ) {
      setErrors((prev) => ({
        ...prev,
        confirm_password: t("hubUserOnboarding.error.passwordMismatch"),
      }));
    }
  };

  if (existingSession) {
    return (
      <Box sx={{ bgcolor: "#f5f8fa", minHeight: "100vh", py: 4 }}>
        <Container maxWidth="md">
          <Paper elevation={3} sx={{ p: 4, bgcolor: "white" }}>
            <Alert severity="warning" sx={{ mb: 2 }}>
              <Typography variant="h6" gutterBottom>
                Another User is Already Signed In
              </Typography>
              <Typography>
                Please sign out of the current account before creating a new
                one.
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
      </Box>
    );
  }

  if (success) {
    return (
      <Box sx={{ bgcolor: "#f5f8fa", minHeight: "100vh", py: 4 }}>
        <Container maxWidth="md">
          <Paper elevation={3} sx={{ p: 4, bgcolor: "white" }}>
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
      </Box>
    );
  }

  return (
    <Box sx={{ bgcolor: "#f5f8fa", minHeight: "100vh", py: 4 }}>
      <Container maxWidth="md">
        <Paper elevation={3} sx={{ p: 4, bgcolor: "white" }}>
          <Typography variant="h4" gutterBottom>
            {t("hubUserOnboarding.title")}
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
              type={showPassword ? "text" : "password"}
              label={t("hubUserOnboarding.form.password")}
              placeholder={t("hubUserOnboarding.form.passwordPlaceholder")}
              value={formData.password}
              onChange={handleInputChange("password")}
              error={!!errors.password}
              helperText={errors.password}
              margin="normal"
              required
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton
                      aria-label="toggle password visibility"
                      onClick={() => setShowPassword(!showPassword)}
                      edge="end"
                    >
                      {showPassword ? <VisibilityOff /> : <Visibility />}
                    </IconButton>
                  </InputAdornment>
                ),
              }}
            />

            <TextField
              fullWidth
              type={showConfirmPassword ? "text" : "password"}
              label={t("hubUserOnboarding.form.confirmPassword")}
              placeholder={t(
                "hubUserOnboarding.form.confirmPasswordPlaceholder"
              )}
              value={formData.confirm_password}
              onChange={handleInputChange("confirm_password")}
              onBlur={handleConfirmPasswordBlur}
              error={!!errors.confirm_password}
              helperText={errors.confirm_password}
              margin="normal"
              required
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton
                      aria-label="toggle password visibility"
                      onClick={() =>
                        setShowConfirmPassword(!showConfirmPassword)
                      }
                      edge="end"
                    >
                      {showConfirmPassword ? <VisibilityOff /> : <Visibility />}
                    </IconButton>
                  </InputAdornment>
                ),
              }}
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

            <Box sx={{ my: 4 }}>
              <Box
                sx={{
                  display: "grid",
                  gridTemplateColumns: { xs: "1fr", sm: "1fr 1fr" },
                  gap: 2,
                }}
              >
                <Card
                  sx={{
                    p: 3,
                    cursor: "pointer",
                    border: (theme) =>
                      formData.selected_tier === "FREE_HUB_USER"
                        ? `2px solid ${theme.palette.primary.main}`
                        : "1px solid #e0e0e0",
                    bgcolor: (theme) =>
                      formData.selected_tier === "FREE_HUB_USER"
                        ? theme.palette.primary.light + "20"
                        : "#ffffff",
                    opacity:
                      formData.selected_tier === "FREE_HUB_USER" ? 1 : 0.6,
                    transform:
                      formData.selected_tier === "FREE_HUB_USER"
                        ? "scale(1.02)"
                        : "none",
                    boxShadow:
                      formData.selected_tier === "FREE_HUB_USER" ? 4 : 1,
                    transition: "all 0.2s",
                    "&:hover": {
                      borderColor: "primary.main",
                      transform: "scale(1.02)",
                      boxShadow: 4,
                      opacity: 1,
                    },
                  }}
                  onClick={() =>
                    setFormData((prev) => ({
                      ...prev,
                      selected_tier: "FREE_HUB_USER",
                    }))
                  }
                >
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="h5" gutterBottom>
                      {t("hubUserOnboarding.form.tier.free")}
                    </Typography>
                    <Typography
                      variant="h4"
                      color="primary"
                      sx={{ mb: 2, fontWeight: "bold" }}
                    >
                      Free
                    </Typography>
                  </Box>
                  <Box sx={{ mt: 2 }}>
                    <Typography component="ul" sx={{ pl: 2 }}>
                      <Typography component="li">Profile</Typography>
                      <Typography component="li">Job search</Typography>
                      <Typography component="li">Ads</Typography>
                    </Typography>
                  </Box>
                </Card>

                <Card
                  sx={{
                    p: 3,
                    cursor: "pointer",
                    border: (theme) =>
                      formData.selected_tier === "PAID_HUB_USER"
                        ? `2px solid ${theme.palette.primary.main}`
                        : "1px solid #e0e0e0",
                    bgcolor: (theme) =>
                      formData.selected_tier === "PAID_HUB_USER"
                        ? theme.palette.primary.light + "20"
                        : "#ffffff",
                    opacity:
                      formData.selected_tier === "PAID_HUB_USER" ? 1 : 0.6,
                    transform:
                      formData.selected_tier === "PAID_HUB_USER"
                        ? "scale(1.02)"
                        : "none",
                    boxShadow:
                      formData.selected_tier === "PAID_HUB_USER" ? 4 : 1,
                    transition: "all 0.2s",
                    "&:hover": {
                      borderColor: "primary.main",
                      transform: "scale(1.02)",
                      boxShadow: 4,
                      opacity: 1,
                    },
                  }}
                  onClick={() =>
                    setFormData((prev) => ({
                      ...prev,
                      selected_tier: "PAID_HUB_USER",
                    }))
                  }
                >
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="h5" gutterBottom>
                      {t("hubUserOnboarding.form.tier.paid")}
                    </Typography>
                    <Typography
                      variant="h4"
                      color="primary"
                      sx={{ mb: 2, fontWeight: "bold" }}
                    >
                      99$ per year
                    </Typography>
                  </Box>
                  <Box sx={{ mt: 2 }}>
                    <Typography component="ul" sx={{ pl: 2 }}>
                      <Typography component="li">All in Free Tier</Typography>
                      <Typography component="li">No Ads</Typography>
                      <Typography component="li">
                        Support Open Source Software development
                      </Typography>
                    </Typography>
                  </Box>
                </Card>
              </Box>
            </Box>

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
              disabled={isSubmitting || !isFormValid()}
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
    </Box>
  );
}
