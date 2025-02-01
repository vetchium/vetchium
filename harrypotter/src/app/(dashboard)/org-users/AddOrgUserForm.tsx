"use client";

import { useState } from "react";
import { useTranslation } from "@/hooks/useTranslation";
import {
  TextField,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Box,
  FormHelperText,
  Chip,
  CircularProgress,
  Alert,
} from "@mui/material";

interface AddOrgUserFormProps {
  onSubmit: (data: {
    email: string;
    name: string;
    roles: string[];
  }) => Promise<void>;
  onCancel: () => void;
}

const availableRoles = [
  "ADMIN",
  "ORG_USERS_CRUD",
  "ORG_USERS_VIEWER",
  "COST_CENTERS_CRUD",
  "COST_CENTERS_VIEWER",
  "LOCATIONS_CRUD",
  "LOCATIONS_VIEWER",
  "OPENINGS_CRUD",
  "OPENINGS_VIEWER",
  "APPLICATIONS_CRUD",
  "APPLICATIONS_VIEWER",
] as const;

export function AddOrgUserForm({ onSubmit, onCancel }: AddOrgUserFormProps) {
  const { t } = useTranslation();
  const [email, setEmail] = useState("");
  const [name, setName] = useState("");
  const [roles, setRoles] = useState<string[]>([]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [errors, setErrors] = useState<{
    email?: string;
    name?: string;
    roles?: string;
  }>({});

  const validateForm = () => {
    const newErrors: typeof errors = {};

    if (!email) {
      newErrors.email = t("validation.email.required");
    } else if (!/^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i.test(email)) {
      newErrors.email = t("validation.email.invalid");
    }

    if (!name) {
      newErrors.name = t("validation.name.required");
    } else if (name.length < 2 || name.length > 64) {
      newErrors.name = t("validation.name.length.2.64");
    }

    if (roles.length === 0) {
      newErrors.roles = t("validation.roles.required");
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validateForm()) return;

    setIsSubmitting(true);
    setSubmitError(null);

    try {
      await onSubmit({ email, name, roles });
    } catch (error) {
      console.error("Failed to submit form:", error);
      setSubmitError(t("orgUsers.addError"));
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Box component="form" onSubmit={handleSubmit} sx={{ mt: 2 }}>
      {submitError && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {submitError}
        </Alert>
      )}

      <TextField
        label={t("orgUsers.email")}
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        error={!!errors.email}
        helperText={errors.email}
        fullWidth
        margin="normal"
        disabled={isSubmitting}
        required
      />

      <TextField
        label={t("orgUsers.name")}
        value={name}
        onChange={(e) => setName(e.target.value)}
        error={!!errors.name}
        helperText={errors.name}
        fullWidth
        margin="normal"
        disabled={isSubmitting}
        required
        inputProps={{ minLength: 2, maxLength: 64 }}
      />

      <FormControl fullWidth margin="normal" error={!!errors.roles} required>
        <InputLabel>{t("orgUsers.rolesList")}</InputLabel>
        <Select
          multiple
          value={roles}
          onChange={(e) =>
            setRoles(
              typeof e.target.value === "string"
                ? [e.target.value]
                : e.target.value
            )
          }
          renderValue={(selected) => (
            <Box sx={{ display: "flex", flexWrap: "wrap", gap: 0.5 }}>
              {selected.map((role) => (
                <Chip
                  key={role}
                  label={t(`orgUsers.roles.${role}`)}
                  size="small"
                />
              ))}
            </Box>
          )}
          disabled={isSubmitting}
        >
          {availableRoles.map((role) => (
            <MenuItem key={role} value={role}>
              {t(`orgUsers.roles.${role}`)}
            </MenuItem>
          ))}
        </Select>
        {errors.roles && <FormHelperText>{errors.roles}</FormHelperText>}
      </FormControl>

      <Box sx={{ mt: 3, display: "flex", justifyContent: "flex-end", gap: 1 }}>
        <Button onClick={onCancel} disabled={isSubmitting}>
          {t("common.cancel")}
        </Button>
        <Button
          type="submit"
          variant="contained"
          disabled={isSubmitting}
          startIcon={isSubmitting ? <CircularProgress size={20} /> : null}
        >
          {t("orgUsers.add")}
        </Button>
      </Box>
    </Box>
  );
}
