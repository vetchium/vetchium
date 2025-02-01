"use client";

import {
  Box,
  Button,
  Container,
  TextField,
  Typography,
  Alert,
  Paper,
  Chip,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  OutlinedInput,
} from "@mui/material";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { useTranslation } from "@/hooks/useTranslation";
import { config } from "@/config";
import Cookies from "js-cookie";

const AVAILABLE_ROLES = [
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
];

export default function CreateOrgUserPage() {
  const [email, setEmail] = useState("");
  const [name, setName] = useState("");
  const [roles, setRoles] = useState<string[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const router = useRouter();
  const { t } = useTranslation();

  const handleSave = async () => {
    try {
      setIsLoading(true);
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/add-org-user`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            email,
            name,
            roles,
          }),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/signin");
        return;
      }

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || t("orgUsers.addError"));
      }

      router.push("/org-users");
    } catch (err) {
      setError(err instanceof Error ? err.message : t("orgUsers.addError"));
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Container maxWidth="md">
      <Paper sx={{ p: 4 }}>
        <Box sx={{ mb: 4 }}>
          <Typography variant="h4" component="h1" gutterBottom>
            {t("orgUsers.addTitle")}
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
            label={t("orgUsers.email")}
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            error={email.length > 0 && !email.includes("@")}
            helperText={
              email.length > 0 && !email.includes("@")
                ? t("validation.email.invalid")
                : ""
            }
          />

          <TextField
            margin="normal"
            required
            fullWidth
            label={t("orgUsers.name")}
            value={name}
            onChange={(e) => setName(e.target.value)}
            inputProps={{ minLength: 2, maxLength: 64 }}
            error={name.length > 0 && (name.length < 2 || name.length > 64)}
            helperText={
              name.length > 0 && (name.length < 2 || name.length > 64)
                ? "Name must be between 2 and 64 characters"
                : ""
            }
          />

          <FormControl fullWidth margin="normal" required>
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
              input={<OutlinedInput label={t("orgUsers.rolesList")} />}
              renderValue={(selected) => (
                <Box sx={{ display: "flex", flexWrap: "wrap", gap: 0.5 }}>
                  {selected.map((value) => (
                    <Chip key={value} label={t(`orgUsers.roles.${value}`)} />
                  ))}
                </Box>
              )}
            >
              {AVAILABLE_ROLES.map((role) => (
                <MenuItem key={role} value={role}>
                  {t(`orgUsers.roles.${role}`)}
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          <Box sx={{ mt: 4, display: "flex", gap: 2 }}>
            <Button
              variant="outlined"
              onClick={() => router.push("/org-users")}
            >
              {t("common.cancel")}
            </Button>
            <Button
              variant="contained"
              onClick={handleSave}
              disabled={
                isLoading ||
                !email ||
                !email.includes("@") ||
                !name ||
                name.length < 2 ||
                name.length > 64 ||
                roles.length === 0
              }
            >
              {isLoading ? t("common.loading") : t("common.save")}
            </Button>
          </Box>
        </Box>
      </Paper>
    </Container>
  );
}
