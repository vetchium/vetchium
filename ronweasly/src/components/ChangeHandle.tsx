import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import CircularProgress from "@mui/material/CircularProgress";
import Link from "@mui/material/Link";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import {
  CheckHandleAvailabilityRequest,
  CheckHandleAvailabilityResponse,
  HubUserTier,
  HubUserTiers,
  SetHandleRequest,
} from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useState } from "react";

interface ChangeHandleProps {
  currentHandle: string;
  userTier: HubUserTier;
  onSuccess?: () => void;
}

export default function ChangeHandle({
  currentHandle,
  userTier,
  onSuccess,
}: ChangeHandleProps) {
  const { t } = useTranslation();
  const [newHandle, setNewHandle] = useState("");
  const [isLoadingCheck, setIsLoadingCheck] = useState(false);
  const [isLoadingSet, setIsLoadingSet] = useState(false);
  const [availabilityResult, setAvailabilityResult] =
    useState<CheckHandleAvailabilityResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [handleValidationError, setHandleValidationError] = useState<
    string | null
  >(null);

  const isPaidUser = userTier === HubUserTiers.PaidHubUserTier;

  const validateHandle = (handle: string): boolean => {
    // Basic validation (adjust regex as needed based on actual rules)
    const handleRegex = /^[a-zA-Z0-9_]{3,32}$/;
    if (!handleRegex.test(handle)) {
      setHandleValidationError(t("profile.changeHandle.error.invalidFormat"));
      return false;
    } else {
      setHandleValidationError(null);
      return true;
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setNewHandle(value);
    setAvailabilityResult(null); // Reset availability when handle changes
    setError(null);
    setSuccess(null);
    validateHandle(value); // Validate on change
  };

  const handleCheckAvailability = async () => {
    if (!validateHandle(newHandle) || !newHandle) {
      return;
    }
    setError(null);
    setSuccess(null);
    setIsLoadingCheck(true);
    setAvailabilityResult(null);
    const token = Cookies.get("session_token");

    if (!token) {
      setError(t("common.error.notAuthenticated")); // Or redirect
      setIsLoadingCheck(false);
      return;
    }

    try {
      const request: CheckHandleAvailabilityRequest = { handle: newHandle };
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/check-handle-availability`,
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
        // Handle specific errors like 401, 403 (if free user tries)
        throw new Error(t("profile.changeHandle.error.checkFailed"));
      }

      const result: CheckHandleAvailabilityResponse = await response.json();
      setAvailabilityResult(result);
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : t("profile.changeHandle.error.checkFailed")
      );
    } finally {
      setIsLoadingCheck(false);
    }
  };

  const handleSetHandle = async () => {
    if (!validateHandle(newHandle) || !availabilityResult?.is_available) {
      setError(t("profile.changeHandle.error.notAvailableOrInvalid"));
      return;
    }
    setError(null);
    setSuccess(null);
    setIsLoadingSet(true);
    const token = Cookies.get("session_token");

    if (!token) {
      setError(t("common.error.notAuthenticated")); // Or redirect
      setIsLoadingSet(false);
      return;
    }

    try {
      const request: SetHandleRequest = { handle: newHandle };
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/set-handle`,
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
        // Handle specific errors: 401, 403 (Unpaid), 409 (Conflict)
        if (response.status === 409) {
          throw new Error(t("settings.changeHandle.error.conflict"));
        }
        throw new Error(t("settings.changeHandle.error.setFailed"));
      }

      setSuccess(t("profile.changeHandle.success"));

      // Clear the input and results
      setNewHandle("");
      setAvailabilityResult(null);

      // Notify parent component of success and refresh page after a short delay
      setTimeout(() => {
        if (onSuccess) {
          onSuccess();
        }
        window.location.reload();
      }, 1500); // Give user time to see success message
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : t("profile.changeHandle.error.setFailed")
      );
    } finally {
      setIsLoadingSet(false);
    }
  };

  return (
    <Box>
      <Typography variant="body1" paragraph>
        {t("profile.changeHandle.currentHandle")}:{" "}
        <strong>{currentHandle}</strong>
      </Typography>

      {!isPaidUser && (
        <Alert severity="info">
          {t("profile.upgradeRequired.message")}{" "}
          <Link href="/upgrade" underline="hover">
            {t("profile.upgradeRequired.upgradeButton")}
          </Link>
        </Alert>
      )}

      {isPaidUser && (
        <Box
          component="form"
          noValidate
          sx={{ mt: 2 }}
          onSubmit={(e) => {
            e.preventDefault();
            // Decide action based on state? Or separate buttons?
            // Assuming separate buttons for check and set
          }}
        >
          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}
          {success && (
            <Alert severity="success" sx={{ mb: 2 }}>
              {success}
            </Alert>
          )}
          <TextField
            fullWidth
            label={t("profile.changeHandle.newHandleLabel")}
            placeholder={t("profile.changeHandle.newHandlePlaceholder")}
            value={newHandle}
            onChange={handleInputChange}
            disabled={isLoadingCheck || isLoadingSet}
            required
            error={!!handleValidationError}
            helperText={
              handleValidationError || t("profile.changeHandle.formatHelp")
            }
            margin="normal"
            InputProps={{
              endAdornment: (
                <Button
                  onClick={handleCheckAvailability}
                  disabled={
                    !newHandle || isLoadingCheck || !!handleValidationError
                  }
                  size="small"
                  sx={{ ml: 1 }}
                >
                  {isLoadingCheck ? (
                    <CircularProgress size={20} />
                  ) : (
                    t("profile.changeHandle.checkAvailabilityButton")
                  )}
                </Button>
              ),
            }}
          />

          {availabilityResult && (
            <Alert
              severity={availabilityResult.is_available ? "success" : "warning"}
              sx={{ mt: 1, mb: 2 }}
            >
              {availabilityResult.is_available
                ? t("profile.changeHandle.available")
                : t("profile.changeHandle.notAvailable")}
              {!availabilityResult.is_available &&
                availabilityResult.suggested_alternatives &&
                availabilityResult.suggested_alternatives.length > 0 && (
                  <>
                    <br />
                    {t("profile.changeHandle.suggestions")}:{" "}
                    {availabilityResult.suggested_alternatives.join(", ")}
                  </>
                )}
            </Alert>
          )}

          <Box sx={{ mt: 2, display: "flex", justifyContent: "flex-end" }}>
            <Button
              variant="contained"
              onClick={handleSetHandle}
              disabled={
                !newHandle ||
                !availabilityResult?.is_available ||
                isLoadingSet ||
                !!handleValidationError
              }
              startIcon={
                isLoadingSet ? <CircularProgress size={20} /> : undefined
              }
            >
              {t("profile.changeHandle.setHandleButton")}
            </Button>
          </Box>
        </Box>
      )}
    </Box>
  );
}
