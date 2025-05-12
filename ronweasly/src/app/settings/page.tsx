"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import UserInvite from "@/components/UserInvite";
import { useTranslation } from "@/hooks/useTranslation";
import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Container from "@mui/material/Container";
import Divider from "@mui/material/Divider";
import FormControl from "@mui/material/FormControl";
import FormHelperText from "@mui/material/FormHelperText";
import Grid from "@mui/material/Grid";
import Paper from "@mui/material/Paper";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import {
  ChangeEmailAddressRequest,
  ChangePasswordRequest,
  validateEmailAddress,
  validatePassword,
} from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

// API functions for changing password and email
const changePassword = async (
  data: ChangePasswordRequest
): Promise<Response> => {
  const token = Cookies.get("session_token");
  return fetch("/api/hub/change-password", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  });
};

const changeEmailAddress = async (
  data: ChangeEmailAddressRequest
): Promise<Response> => {
  const token = Cookies.get("session_token");
  return fetch("/api/hub/change-email-address", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  });
};

export default function Settings() {
  const { t } = useTranslation();
  const router = useRouter();

  // Password change state
  const [oldPassword, setOldPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [passwordError, setPasswordError] = useState<string | null>(null);
  const [passwordSuccess, setPasswordSuccess] = useState(false);
  const [passwordLoading, setPasswordLoading] = useState(false);
  const [isNewPasswordValid, setIsNewPasswordValid] = useState(false);

  // Email change state
  const [newEmail, setNewEmail] = useState("");
  const [confirmEmail, setConfirmEmail] = useState("");
  const [emailError, setEmailError] = useState<string | null>(null);
  const [emailSuccess, setEmailSuccess] = useState(false);
  const [emailLoading, setEmailLoading] = useState(false);
  const [isNewEmailValid, setIsNewEmailValid] = useState(false);

  // Auth check
  useEffect(() => {
    const token = Cookies.get("session_token");
    if (!token) {
      router.push("/login");
    }
  }, [router]);

  // Handle new password change
  const handleNewPasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const password = e.target.value;
    setNewPassword(password);

    // Only reset confirmation if we had a valid password before
    if (isNewPasswordValid) {
      setConfirmPassword("");
    }

    // Clear errors while typing
    if (passwordError) {
      setPasswordError(null);
    }
  };

  // Validate password on blur
  const validateNewPassword = () => {
    if (newPassword) {
      const validation = validatePassword(newPassword);
      setIsNewPasswordValid(validation.isValid);
      if (!validation.isValid) {
        setPasswordError(validation.error || "Invalid password");
      } else {
        setPasswordError(null);
      }
    } else {
      setIsNewPasswordValid(false);
      setPasswordError(null);
    }
  };

  // Handle new email change
  const handleNewEmailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const email = e.target.value;
    setNewEmail(email);

    // Only reset confirmation if we had a valid email before
    if (isNewEmailValid) {
      setConfirmEmail("");
    }

    // Clear errors while typing
    if (emailError) {
      setEmailError(null);
    }
  };

  // Validate email on blur
  const validateNewEmail = () => {
    if (newEmail) {
      const isValid = validateEmailAddress(newEmail);
      setIsNewEmailValid(isValid);
      if (!isValid) {
        setEmailError("Please enter a valid email address");
      } else {
        setEmailError(null);
      }
    } else {
      setIsNewEmailValid(false);
      setEmailError(null);
    }
  };

  // Handle password change
  const handlePasswordChange = async (e: React.FormEvent) => {
    e.preventDefault();
    setPasswordError(null);
    setPasswordSuccess(false);

    // Validate new password
    const passwordValidation = validatePassword(newPassword);
    if (!passwordValidation.isValid) {
      setPasswordError(passwordValidation.error || "Invalid password");
      return;
    }

    // Check if passwords match
    if (newPassword !== confirmPassword) {
      setPasswordError("Passwords do not match");
      return;
    }

    // Prepare request data
    const data: ChangePasswordRequest = {
      old_password: oldPassword,
      new_password: newPassword,
    };

    try {
      setPasswordLoading(true);
      const response = await changePassword(data);

      if (response.ok) {
        setPasswordSuccess(true);
        setOldPassword("");
        setNewPassword("");
        setConfirmPassword("");
      } else {
        if (response.status === 401) {
          setPasswordError(t("settings.changePassword.error.invalidPassword"));
        } else {
          setPasswordError(t("settings.changePassword.error.failed"));
        }
      }
    } catch (error) {
      setPasswordError(t("settings.changePassword.error.failed"));
    } finally {
      setPasswordLoading(false);
    }
  };

  // Handle email change
  const handleEmailChange = async (e: React.FormEvent) => {
    e.preventDefault();
    setEmailError(null);
    setEmailSuccess(false);

    // Check if emails match
    if (newEmail !== confirmEmail) {
      setEmailError("Email addresses do not match");
      return;
    }

    // Validate email
    if (!validateEmailAddress(newEmail)) {
      setEmailError("Please enter a valid email address");
      return;
    }

    // Prepare request data
    const data: ChangeEmailAddressRequest = {
      email: newEmail,
    };

    try {
      setEmailLoading(true);
      const response = await changeEmailAddress(data);

      if (response.ok) {
        setEmailSuccess(true);
        setNewEmail("");
        setConfirmEmail("");
      } else {
        if (response.status === 409) {
          setEmailError(t("settings.changeEmail.error.emailInUse"));
        } else {
          const errorData = await response.json().catch(() => ({}));
          setEmailError(
            errorData.message || t("settings.changeEmail.error.failed")
          );
        }
      }
    } catch (error) {
      setEmailError("An error occurred while changing email address");
      console.error("Email change error:", error);
    } finally {
      setEmailLoading(false);
    }
  };

  return (
    <AuthenticatedLayout>
      <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          {t("settings.title")}
        </Typography>

        <Grid container spacing={3}>
          {/* Change Password Section */}
          <Grid item xs={12} md={6}>
            <Paper sx={{ p: 3, mb: 3 }}>
              <Typography variant="h6" gutterBottom>
                Change Password
              </Typography>
              <Divider sx={{ mb: 2 }} />

              {passwordSuccess && (
                <Alert severity="success" sx={{ mb: 2 }}>
                  Password changed successfully!
                </Alert>
              )}

              {passwordError && (
                <Alert severity="error" sx={{ mb: 2 }}>
                  {passwordError}
                </Alert>
              )}

              <Box
                component="form"
                onSubmit={handlePasswordChange}
                noValidate
                tabIndex={0} // This creates a new tab sequence
              >
                <FormControl fullWidth margin="normal">
                  <TextField
                    required
                    id="old-password"
                    label="Current Password"
                    type="password"
                    value={oldPassword}
                    onChange={(e) => setOldPassword(e.target.value)}
                    disabled={passwordLoading}
                  />
                </FormControl>

                <FormControl fullWidth margin="normal">
                  <TextField
                    required
                    id="new-password"
                    label="New Password"
                    type="password"
                    value={newPassword}
                    onChange={handleNewPasswordChange}
                    onBlur={validateNewPassword}
                    disabled={passwordLoading}
                  />
                  <FormHelperText>
                    Password must be at least 12 characters long and include
                    uppercase, lowercase, numbers, and special characters.
                  </FormHelperText>
                </FormControl>

                <FormControl fullWidth margin="normal">
                  <TextField
                    required
                    id="confirm-password"
                    label="Confirm New Password"
                    type="password"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    disabled={passwordLoading}
                    error={!isNewPasswordValid && confirmPassword !== ""}
                    helperText={
                      !isNewPasswordValid && confirmPassword !== ""
                        ? "Please enter a valid password above first"
                        : ""
                    }
                  />
                </FormControl>

                <Button
                  type="submit"
                  variant="contained"
                  color="primary"
                  disabled={
                    passwordLoading ||
                    !oldPassword ||
                    !newPassword ||
                    !confirmPassword
                  }
                  sx={{ mt: 2 }}
                >
                  {passwordLoading ? "Changing..." : "Change Password"}
                </Button>
              </Box>
            </Paper>
          </Grid>

          {/* Change Email Section */}
          <Grid item xs={12} md={6}>
            <Paper sx={{ p: 3, mb: 3 }}>
              <Typography variant="h6" gutterBottom>
                Change Email Address
              </Typography>
              <Divider sx={{ mb: 2 }} />

              {emailSuccess && (
                <Alert severity="success" sx={{ mb: 2 }}>
                  Email address changed successfully!
                </Alert>
              )}

              {emailError && (
                <Alert severity="error" sx={{ mb: 2 }}>
                  {emailError}
                </Alert>
              )}

              <Box
                component="form"
                onSubmit={handleEmailChange}
                noValidate
                tabIndex={0} // This creates a new tab sequence
              >
                <Alert severity="warning" sx={{ mb: 2 }}>
                  <strong>Warning:</strong> No verification email will be sent.
                  Please double-check your new email address carefully.
                </Alert>

                <FormControl fullWidth margin="normal">
                  <TextField
                    required
                    id="new-email"
                    label="New Email Address"
                    type="email"
                    value={newEmail}
                    onChange={handleNewEmailChange}
                    onBlur={validateNewEmail}
                    disabled={emailLoading}
                  />
                </FormControl>

                <FormControl fullWidth margin="normal">
                  <TextField
                    required
                    id="confirm-email"
                    label="Confirm New Email Address"
                    type="password"
                    autoComplete="new-email"
                    value={confirmEmail}
                    onChange={(e) => setConfirmEmail(e.target.value)}
                    disabled={emailLoading}
                    error={!isNewEmailValid && confirmEmail !== ""}
                    helperText={
                      !isNewEmailValid && confirmEmail !== ""
                        ? "Please enter a valid email address above first"
                        : ""
                    }
                  />
                  <FormHelperText>
                    Please confirm your new email address. Type carefully as the
                    text is masked.
                  </FormHelperText>
                </FormControl>

                <Button
                  type="submit"
                  variant="contained"
                  color="primary"
                  disabled={emailLoading || !newEmail || !confirmEmail}
                  sx={{ mt: 2 }}
                >
                  {emailLoading ? "Changing..." : "Change Email Address"}
                </Button>
              </Box>
            </Paper>
          </Grid>

          {/* Invite User Section */}
          <Grid item xs={12}>
            <Paper sx={{ p: 3 }}>
              <Typography variant="h6" gutterBottom>
                Invite User
              </Typography>
              <Divider sx={{ mb: 2 }} />
              <Box tabIndex={0}>
                {" "}
                {/* This creates a new tab sequence */}
                <UserInvite />
              </Box>
            </Paper>
          </Grid>
        </Grid>
      </Container>
    </AuthenticatedLayout>
  );
}
