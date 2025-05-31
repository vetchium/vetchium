"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import { Visibility, VisibilityOff } from "@mui/icons-material";
import {
  Alert,
  Box,
  Button,
  Card,
  CardContent,
  CircularProgress,
  Container,
  IconButton,
  InputAdornment,
  Snackbar,
  TextField,
  Typography,
} from "@mui/material";
import {
  ChangeCoolOffPeriodRequest,
  EmployerChangePasswordRequest,
} from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";

export default function SettingsPage() {
  const { t } = useTranslation();
  const router = useRouter();
  const [coolOffPeriod, setCoolOffPeriod] = useState<number | "">("");
  const [currentCoolOffPeriod, setCurrentCoolOffPeriod] = useState<
    number | null
  >(null);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [isUpdating, setIsUpdating] = useState(false);

  // Password change state
  const [passwordData, setPasswordData] = useState({
    oldPassword: "",
    newPassword: "",
    confirmPassword: "",
  });
  const [passwordErrors, setPasswordErrors] = useState({
    oldPassword: "",
    newPassword: "",
    confirmPassword: "",
  });
  const [showPasswords, setShowPasswords] = useState({
    oldPassword: false,
    newPassword: false,
    confirmPassword: false,
  });
  const [isChangingPassword, setIsChangingPassword] = useState(false);

  const fetchCoolOffPeriod = useCallback(async () => {
    setIsLoading(true);
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/get-cool-off-period`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/signin");
        return;
      }

      if (!response.ok) {
        throw new Error(t("settings.coolOffPeriod.fetchError"));
      }

      const periodText = await response.text();
      console.log("Response text:", periodText);

      const coolOffPeriod = parseInt(periodText, 10);
      if (isNaN(coolOffPeriod)) {
        console.error("Invalid cool off period in response:", periodText);
        throw new Error(t("settings.coolOffPeriod.fetchError"));
      }

      setCurrentCoolOffPeriod(coolOffPeriod);
      setCoolOffPeriod(coolOffPeriod);
    } catch {
      setError(t("settings.coolOffPeriod.fetchError"));
    } finally {
      setIsLoading(false);
    }
  }, [t, router]);

  useEffect(() => {
    fetchCoolOffPeriod();
  }, [fetchCoolOffPeriod]);

  const handleUpdateCoolOffPeriod = async () => {
    if (coolOffPeriod === "") return;

    setIsUpdating(true);
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const request: ChangeCoolOffPeriodRequest = {
        cool_off_period_days: Number(coolOffPeriod),
      };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/change-cool-off-period`,
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
        throw new Error(t("settings.coolOffPeriod.error"));
      }

      setSuccess(t("settings.coolOffPeriod.success"));
      setCurrentCoolOffPeriod(Number(coolOffPeriod));
    } catch {
      setError(t("settings.coolOffPeriod.error"));
    } finally {
      setIsUpdating(false);
    }
  };

  // Password validation function
  const validatePassword = (password: string): string => {
    if (password.length < 8) {
      return t("settings.password.errors.tooShort");
    }
    if (!/(?=.*[a-z])/.test(password)) {
      return t("settings.password.errors.missingLowercase");
    }
    if (!/(?=.*[A-Z])/.test(password)) {
      return t("settings.password.errors.missingUppercase");
    }
    if (!/(?=.*\d)/.test(password)) {
      return t("settings.password.errors.missingNumber");
    }
    if (!/(?=.*[!@#$%^&*])/.test(password)) {
      return t("settings.password.errors.missingSpecial");
    }
    return "";
  };

  const handlePasswordChange = async () => {
    // Reset errors
    setPasswordErrors({
      oldPassword: "",
      newPassword: "",
      confirmPassword: "",
    });

    // Validate inputs
    let hasErrors = false;
    const newErrors = { ...passwordErrors };

    if (!passwordData.oldPassword) {
      newErrors.oldPassword = t("settings.password.errors.required");
      hasErrors = true;
    }

    if (!passwordData.newPassword) {
      newErrors.newPassword = t("settings.password.errors.required");
      hasErrors = true;
    } else {
      const passwordValidationError = validatePassword(
        passwordData.newPassword
      );
      if (passwordValidationError) {
        newErrors.newPassword = passwordValidationError;
        hasErrors = true;
      }
    }

    if (!passwordData.confirmPassword) {
      newErrors.confirmPassword = t("settings.password.errors.required");
      hasErrors = true;
    } else if (passwordData.newPassword !== passwordData.confirmPassword) {
      newErrors.confirmPassword = t("settings.password.errors.mismatch");
      hasErrors = true;
    }

    if (hasErrors) {
      setPasswordErrors(newErrors);
      return;
    }

    setIsChangingPassword(true);
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const request: EmployerChangePasswordRequest = {
        old_password: passwordData.oldPassword,
        new_password: passwordData.newPassword,
      };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/change-password`,
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
        if (
          response.headers.get("content-type")?.includes("application/json")
        ) {
          // This might be a validation error or incorrect old password
          try {
            const errorData = await response.json();
            if (errorData.errors && errorData.errors.includes("old_password")) {
              setPasswordErrors({
                ...passwordErrors,
                oldPassword: t("settings.password.errors.incorrectOldPassword"),
              });
              return;
            }
          } catch {
            // Fall through to token expiry handling
          }
        }
        // Token expired
        Cookies.remove("session_token");
        router.push("/signin");
        return;
      }

      if (response.status === 400) {
        // Validation errors
        try {
          const errorData = await response.json();
          const newErrors = { ...passwordErrors };
          if (errorData.errors) {
            if (errorData.errors.includes("old_password")) {
              newErrors.oldPassword = t("settings.password.errors.invalid");
            }
            if (errorData.errors.includes("new_password")) {
              newErrors.newPassword = t("settings.password.errors.invalid");
            }
          }
          setPasswordErrors(newErrors);
          return;
        } catch {
          throw new Error(t("settings.password.errors.changeError"));
        }
      }

      if (!response.ok) {
        throw new Error(t("settings.password.errors.changeError"));
      }

      // Success
      setSuccess(t("settings.password.success"));
      setPasswordData({
        oldPassword: "",
        newPassword: "",
        confirmPassword: "",
      });
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : t("settings.password.errors.changeError")
      );
    } finally {
      setIsChangingPassword(false);
    }
  };

  const togglePasswordVisibility = (field: keyof typeof showPasswords) => {
    setShowPasswords((prev) => ({
      ...prev,
      [field]: !prev[field],
    }));
  };

  const content = (
    <Container maxWidth="lg">
      <Box sx={{ p: 3 }}>
        <Typography variant="h4" gutterBottom>
          {t("settings.title")}
        </Typography>

        <Card sx={{ mt: 3 }}>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              {t("settings.coolOffPeriod.title")}
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
              {t("settings.coolOffPeriod.description")}
            </Typography>

            <Box sx={{ display: "flex", gap: 2, alignItems: "center" }}>
              <TextField
                type="number"
                value={coolOffPeriod}
                onChange={(e) =>
                  setCoolOffPeriod(
                    e.target.value === "" ? "" : Number(e.target.value)
                  )
                }
                inputProps={{ min: 0, max: 365 }}
                sx={{ width: 200 }}
                disabled={isUpdating}
              />
              <Button
                variant="contained"
                onClick={handleUpdateCoolOffPeriod}
                disabled={
                  coolOffPeriod === "" ||
                  coolOffPeriod === currentCoolOffPeriod ||
                  isUpdating
                }
                startIcon={
                  isUpdating ? (
                    <CircularProgress size={20} color="inherit" />
                  ) : null
                }
              >
                {t("settings.coolOffPeriod.update")}
              </Button>
            </Box>
          </CardContent>
        </Card>

        {/* Password Change Section */}
        <Card sx={{ mt: 3 }}>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              {t("settings.password.title")}
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
              {t("settings.password.description")}
            </Typography>

            <Box sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
              <TextField
                type={showPasswords.oldPassword ? "text" : "password"}
                label={t("settings.password.oldPassword")}
                value={passwordData.oldPassword}
                onChange={(e) =>
                  setPasswordData((prev) => ({
                    ...prev,
                    oldPassword: e.target.value,
                  }))
                }
                error={!!passwordErrors.oldPassword}
                helperText={passwordErrors.oldPassword}
                disabled={isChangingPassword}
                InputProps={{
                  endAdornment: (
                    <InputAdornment position="end">
                      <IconButton
                        onClick={() => togglePasswordVisibility("oldPassword")}
                        edge="end"
                      >
                        {showPasswords.oldPassword ? (
                          <VisibilityOff />
                        ) : (
                          <Visibility />
                        )}
                      </IconButton>
                    </InputAdornment>
                  ),
                }}
              />
              <TextField
                type={showPasswords.newPassword ? "text" : "password"}
                label={t("settings.password.newPassword")}
                value={passwordData.newPassword}
                onChange={(e) =>
                  setPasswordData((prev) => ({
                    ...prev,
                    newPassword: e.target.value,
                  }))
                }
                error={!!passwordErrors.newPassword}
                helperText={
                  passwordErrors.newPassword ||
                  t("settings.password.requirements")
                }
                disabled={isChangingPassword}
                InputProps={{
                  endAdornment: (
                    <InputAdornment position="end">
                      <IconButton
                        onClick={() => togglePasswordVisibility("newPassword")}
                        edge="end"
                      >
                        {showPasswords.newPassword ? (
                          <VisibilityOff />
                        ) : (
                          <Visibility />
                        )}
                      </IconButton>
                    </InputAdornment>
                  ),
                }}
              />
              <TextField
                type={showPasswords.confirmPassword ? "text" : "password"}
                label={t("settings.password.confirmPassword")}
                value={passwordData.confirmPassword}
                onChange={(e) =>
                  setPasswordData((prev) => ({
                    ...prev,
                    confirmPassword: e.target.value,
                  }))
                }
                error={!!passwordErrors.confirmPassword}
                helperText={passwordErrors.confirmPassword}
                disabled={isChangingPassword}
                InputProps={{
                  endAdornment: (
                    <InputAdornment position="end">
                      <IconButton
                        onClick={() =>
                          togglePasswordVisibility("confirmPassword")
                        }
                        edge="end"
                      >
                        {showPasswords.confirmPassword ? (
                          <VisibilityOff />
                        ) : (
                          <Visibility />
                        )}
                      </IconButton>
                    </InputAdornment>
                  ),
                }}
              />
              <Box sx={{ mt: 2 }}>
                <Button
                  variant="contained"
                  onClick={handlePasswordChange}
                  disabled={
                    !passwordData.oldPassword ||
                    !passwordData.newPassword ||
                    !passwordData.confirmPassword ||
                    isChangingPassword
                  }
                  startIcon={
                    isChangingPassword ? (
                      <CircularProgress size={20} color="inherit" />
                    ) : null
                  }
                >
                  {t("settings.password.change")}
                </Button>
              </Box>
            </Box>
          </CardContent>
        </Card>

        <Snackbar
          open={error !== null}
          autoHideDuration={6000}
          onClose={() => setError(null)}
        >
          <Alert severity="error" onClose={() => setError(null)}>
            {error}
          </Alert>
        </Snackbar>

        <Snackbar
          open={success !== null}
          autoHideDuration={6000}
          onClose={() => setSuccess(null)}
        >
          <Alert severity="success" onClose={() => setSuccess(null)}>
            {success}
          </Alert>
        </Snackbar>
      </Box>
    </Container>
  );

  if (isLoading) {
    return (
      <AuthenticatedLayout>
        <Container maxWidth="lg">
          <Box
            sx={{
              display: "flex",
              justifyContent: "center",
              alignItems: "center",
              height: "50vh",
            }}
          >
            <CircularProgress />
          </Box>
        </Container>
      </AuthenticatedLayout>
    );
  }

  return <AuthenticatedLayout>{content}</AuthenticatedLayout>;
}
