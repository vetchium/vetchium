"use client";

import { useState, useEffect } from "react";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import Paper from "@mui/material/Paper";
import IconButton from "@mui/material/IconButton";
import DeleteIcon from "@mui/icons-material/Delete";
import VerifiedIcon from "@mui/icons-material/Verified";
import CircularProgress from "@mui/material/CircularProgress";
import Alert from "@mui/material/Alert";
import Dialog from "@mui/material/Dialog";
import DialogTitle from "@mui/material/DialogTitle";
import DialogContent from "@mui/material/DialogContent";
import DialogActions from "@mui/material/DialogActions";
import { useTranslation } from "@/hooks/useTranslation";
import { config } from "@/config";
import Cookies from "js-cookie";

interface OfficialEmail {
  email: string;
  last_verified_at: string | null;
  verify_in_progress: boolean;
}

export default function OfficialEmails() {
  const { t } = useTranslation();
  const [emails, setEmails] = useState<OfficialEmail[]>([]);
  const [loading, setLoading] = useState(true);
  const [listError, setListError] = useState("");
  const [addError, setAddError] = useState("");
  const [newEmail, setNewEmail] = useState("");
  const [verificationCode, setVerificationCode] = useState("");
  const [verifyingEmail, setVerifyingEmail] = useState<string | null>(null);
  const [verifyError, setVerifyError] = useState("");
  const [addingEmail, setAddingEmail] = useState(false);
  const [deletingEmail, setDeletingEmail] = useState<string | null>(null);
  const [showAddForm, setShowAddForm] = useState(false);
  const [isValidEmail, setIsValidEmail] = useState(false);

  // Fetch official emails
  const fetchEmails = async () => {
    try {
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/my-official-emails`,
        {
          headers: {
            Authorization: `Bearer ${Cookies.get("session_token")}`,
          },
        }
      );

      if (!response.ok) {
        throw new Error(t("officialEmails.errors.loadFailed"));
      }

      const data = await response.json();
      setEmails(data);
      setListError("");
    } catch (err) {
      setListError(t("officialEmails.errors.loadFailed"));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchEmails();
  }, []);

  const validateEmail = (email: string) => {
    const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    return emailRegex.test(email);
  };

  const handleEmailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const email = e.target.value;
    setNewEmail(email);
    setIsValidEmail(validateEmail(email));
  };

  // Add new email
  const handleAddEmail = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newEmail) return;

    setAddingEmail(true);
    setAddError("");

    try {
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/add-official-email`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${Cookies.get("session_token")}`,
          },
          body: JSON.stringify({ email: newEmail }),
        }
      );

      if (!response.ok) {
        switch (response.status) {
          case 409:
            throw new Error(t("officialEmails.errors.emailExists"));
          case 422:
            throw new Error(t("officialEmails.errors.domainNotEmployer"));
          default:
            throw new Error(t("officialEmails.errors.addFailed"));
        }
      }

      setNewEmail("");
      setShowAddForm(false);
      await fetchEmails();
    } catch (err) {
      setAddError(
        err instanceof Error
          ? err.message
          : t("officialEmails.errors.addFailed")
      );
    } finally {
      setAddingEmail(false);
    }
  };

  // Delete email
  const handleDeleteEmail = async (email: string) => {
    setDeletingEmail(email);
    setListError("");

    try {
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/delete-official-email`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${Cookies.get("session_token")}`,
          },
          body: JSON.stringify({ email }),
        }
      );

      if (!response.ok) {
        throw new Error(t("officialEmails.errors.deleteFailed"));
      }

      await fetchEmails();
    } catch (err) {
      setListError(t("officialEmails.errors.deleteFailed"));
    } finally {
      setDeletingEmail(null);
    }
  };

  // Trigger verification
  const handleTriggerVerification = async (email: string) => {
    setListError("");

    try {
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/trigger-verification`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${Cookies.get("session_token")}`,
          },
          body: JSON.stringify({ email }),
        }
      );

      if (!response.ok) {
        throw new Error(t("officialEmails.errors.triggerFailed"));
      }

      setVerifyingEmail(email);
    } catch (err) {
      setListError(t("officialEmails.errors.triggerFailed"));
    }
  };

  // Verify email
  const handleVerifyEmail = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!verifyingEmail || !verificationCode) return;

    setVerifyError("");

    try {
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/verify-official-email`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${Cookies.get("session_token")}`,
          },
          body: JSON.stringify({
            email: verifyingEmail,
            code: verificationCode,
          }),
        }
      );

      if (!response.ok) {
        throw new Error(t("officialEmails.errors.invalidCode"));
      }

      setVerifyingEmail(null);
      setVerificationCode("");
      await fetchEmails();
    } catch (err) {
      setVerifyError(t("officialEmails.errors.invalidCode"));
    }
  };

  const isVerificationExpired = (lastVerifiedAt: string | null) => {
    if (!lastVerifiedAt) return true;
    const verifiedDate = new Date(lastVerifiedAt);
    const ninetyDaysAgo = new Date();
    ninetyDaysAgo.setDate(ninetyDaysAgo.getDate() - 90);
    return verifiedDate < ninetyDaysAgo;
  };

  const formatVerificationDate = (date: string) => {
    return new Date(date).toLocaleDateString(undefined, {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  };

  if (loading) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", mt: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <>
      <Box
        sx={{
          display: "flex",
          justifyContent: "space-between",
          mb: 3,
          alignItems: "center",
        }}
      >
        <Typography variant="h5">{t("officialEmails.title")}</Typography>
        {!showAddForm && !verifyingEmail && (
          <Button
            variant="contained"
            color="primary"
            onClick={() => {
              setAddError("");
              setShowAddForm(true);
            }}
          >
            {t("officialEmails.addEmail")}
          </Button>
        )}
      </Box>

      {showAddForm && (
        <Paper sx={{ p: 3, mb: 4 }}>
          {addError && (
            <Alert severity="error" sx={{ mb: 3 }}>
              {addError}
            </Alert>
          )}
          <Box
            component="form"
            onSubmit={handleAddEmail}
            sx={{
              display: "flex",
              flexDirection: "column",
              gap: 3,
            }}
          >
            <TextField
              fullWidth
              type="email"
              label={t("officialEmails.addEmail")}
              value={newEmail}
              onChange={handleEmailChange}
              disabled={addingEmail}
              size="medium"
              required
              error={newEmail !== "" && !isValidEmail}
              helperText={
                newEmail !== "" && !isValidEmail
                  ? "Please enter a valid email address"
                  : ""
              }
              inputProps={{
                pattern: "[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
              }}
            />
            <Box sx={{ display: "flex", gap: 2 }}>
              <Button
                type="submit"
                variant="contained"
                color="primary"
                disabled={!isValidEmail || addingEmail}
                aria-label={t("officialEmails.addEmailSubmit")}
              >
                {addingEmail ? (
                  <CircularProgress size={24} color="inherit" />
                ) : (
                  t("officialEmails.addEmail")
                )}
              </Button>
              <Button
                onClick={() => {
                  setShowAddForm(false);
                  setNewEmail("");
                  setAddError("");
                }}
                variant="outlined"
                color="inherit"
              >
                {t("common.cancel")}
              </Button>
            </Box>
          </Box>
        </Paper>
      )}

      {verifyingEmail && (
        <Paper sx={{ p: 3, mb: 4 }}>
          {verifyError && (
            <Alert severity="error" sx={{ mb: 3 }}>
              {verifyError}
            </Alert>
          )}
          <Box
            component="form"
            onSubmit={handleVerifyEmail}
            sx={{
              display: "flex",
              flexDirection: "column",
              gap: 3,
            }}
          >
            <Typography variant="body1" sx={{ mb: 1 }}>
              {t("officialEmails.enterVerificationCode", {
                email: verifyingEmail,
              })}
            </Typography>
            <TextField
              fullWidth
              type="text"
              label={t("officialEmails.verificationCode")}
              value={verificationCode}
              onChange={(e) => {
                const value = e.target.value.replace(/[^0-9a-zA-Z]/g, "");
                if (value.length <= 4) {
                  setVerificationCode(value);
                }
              }}
              autoFocus
              required
              error={verificationCode !== "" && verificationCode.length !== 4}
              helperText={
                verificationCode !== "" && verificationCode.length !== 4
                  ? t("officialEmails.errors.invalidCodeLength")
                  : ""
              }
              inputProps={{
                maxLength: 4,
              }}
            />
            <Box sx={{ display: "flex", gap: 2 }}>
              <Button
                type="submit"
                variant="contained"
                color="primary"
                disabled={!verificationCode || verificationCode.length !== 4}
                aria-label={t("officialEmails.verifyEmailSubmit")}
              >
                {t("officialEmails.verifyButton")}
              </Button>
              <Button
                onClick={() => {
                  setVerifyingEmail(null);
                  setVerificationCode("");
                  setVerifyError("");
                }}
                variant="outlined"
                color="inherit"
              >
                {t("common.cancel")}
              </Button>
            </Box>
          </Box>
        </Paper>
      )}

      {!showAddForm && !verifyingEmail && (
        <Paper sx={{ p: 4 }}>
          {listError && (
            <Alert severity="error" sx={{ mb: 3 }}>
              {listError}
            </Alert>
          )}
          {emails.length === 0 ? (
            <Box
              sx={{
                p: 3,
                textAlign: "center",
                bgcolor: "action.hover",
                borderRadius: 1,
              }}
            >
              <Typography color="text.secondary">
                {t("officialEmails.noEmails")}
              </Typography>
            </Box>
          ) : (
            emails.map((email) => (
              <Box
                key={email.email}
                sx={{
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "space-between",
                  p: 2,
                  mb: 1,
                  borderRadius: 1,
                  bgcolor: "background.default",
                  "&:last-child": { mb: 0 },
                }}
              >
                <Box sx={{ display: "flex", flexDirection: "column" }}>
                  <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                    <Typography>{email.email}</Typography>
                    {email.last_verified_at &&
                      !isVerificationExpired(email.last_verified_at) && (
                        <VerifiedIcon color="success" sx={{ fontSize: 20 }} />
                      )}
                  </Box>
                  {email.last_verified_at && (
                    <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                      <Typography variant="caption" color="text.secondary">
                        {t("officialEmails.verifiedOn", {
                          date: formatVerificationDate(email.last_verified_at),
                        })}
                      </Typography>
                      {isVerificationExpired(email.last_verified_at) && (
                        <Typography variant="caption" color="warning.main">
                          {t("officialEmails.verificationExpired")}
                        </Typography>
                      )}
                    </Box>
                  )}
                </Box>
                <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                  {(!email.last_verified_at ||
                    isVerificationExpired(email.last_verified_at)) &&
                    !email.verify_in_progress && (
                      <Button
                        size="small"
                        variant="outlined"
                        onClick={() => handleTriggerVerification(email.email)}
                        aria-label={t("officialEmails.sendCode")}
                      >
                        {t("officialEmails.sendCode")}
                      </Button>
                    )}
                  {email.verify_in_progress && (
                    <Button
                      size="small"
                      variant="outlined"
                      color="info"
                      onClick={() => setVerifyingEmail(email.email)}
                      aria-label={t("officialEmails.enterCode")}
                    >
                      {t("officialEmails.enterCode")}
                    </Button>
                  )}
                  <IconButton
                    onClick={() => handleDeleteEmail(email.email)}
                    disabled={deletingEmail === email.email}
                    color="error"
                    size="small"
                    aria-label={t("officialEmails.deleteEmail")}
                  >
                    {deletingEmail === email.email ? (
                      <CircularProgress size={20} />
                    ) : (
                      <DeleteIcon />
                    )}
                  </IconButton>
                </Box>
              </Box>
            ))
          )}
        </Paper>
      )}
    </>
  );
}
