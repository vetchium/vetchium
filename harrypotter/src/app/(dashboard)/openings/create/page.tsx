"use client";

import {
  Box,
  Button,
  Container,
  TextField,
  Typography,
  Alert,
  Paper,
  Autocomplete,
  MenuItem,
  Select,
  FormControl,
  InputLabel,
  Chip,
} from "@mui/material";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useTranslation } from "@/hooks/useTranslation";
import { config } from "@/config";
import Cookies from "js-cookie";
import {
  CreateOpeningRequest,
  OrgUserShort,
  OpeningType,
  EducationLevel,
  validTimezones,
  TimeZone,
  OpeningTypes,
  EducationLevels,
} from "@psankar/vetchi-typespec";
import { Location, LocationStates } from "@psankar/vetchi-typespec";

export default function CreateOpeningPage() {
  const [title, setTitle] = useState("");
  const [positions, setPositions] = useState(1);
  const [jd, setJd] = useState("");
  const [recruiter, setRecruiter] = useState<OrgUserShort | null>(null);
  const [hiringManager, setHiringManager] = useState<OrgUserShort | null>(null);
  const [costCenterName, setCostCenterName] = useState("");
  const [openingType, setOpeningType] = useState<OpeningType>(
    OpeningTypes.FULL_TIME
  );
  const [yoeMin, setYoeMin] = useState(1);
  const [yoeMax, setYoeMax] = useState(80);
  const [minEducationLevel, setMinEducationLevel] = useState<EducationLevel>(
    EducationLevels.UNSPECIFIED
  );
  const [employerNotes, setEmployerNotes] = useState("");
  const [remoteTimezones, setRemoteTimezones] = useState<string[]>([]);
  const [selectedLocations, setSelectedLocations] = useState<string[]>([]);

  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [orgUsers, setOrgUsers] = useState<OrgUserShort[]>([]);
  const [costCenters, setCostCenters] = useState<string[]>([]);
  const [locations, setLocations] = useState<Location[]>([]);

  const router = useRouter();
  const { t } = useTranslation();

  useEffect(() => {
    fetchOrgUsers();
    fetchCostCenters();
    fetchLocations();
  }, []);

  const fetchOrgUsers = async (searchPrefix?: string) => {
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/filter-org-users`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            prefix: searchPrefix,
            limit: 40,
          }),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/signin");
        return;
      }

      if (!response.ok) {
        throw new Error(t("openings.fetchError"));
      }

      const data = await response.json();
      setOrgUsers(data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : t("openings.fetchError"));
    }
  };

  const fetchCostCenters = async () => {
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/get-cost-centers`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({}),
        }
      );

      if (!response.ok) {
        throw new Error(t("openings.fetchCostCentersError"));
      }

      const data = await response.json();
      setCostCenters(data.map((cc: { name: string }) => cc.name) || []);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("openings.fetchCostCentersError")
      );
    }
  };

  const fetchLocations = async () => {
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/get-locations`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({}),
        }
      );

      if (!response.ok) {
        throw new Error(t("openings.fetchLocationsError"));
      }

      const data = await response.json();
      setLocations(data || []);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("openings.fetchLocationsError")
      );
    }
  };

  const handleSave = async () => {
    try {
      setIsLoading(true);
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      if (!recruiter || !hiringManager) {
        setError(t("openings.missingUserError"));
        return;
      }

      const request: CreateOpeningRequest = {
        title,
        positions,
        jd,
        recruiter: recruiter.email,
        hiring_manager: hiringManager.email,
        cost_center_name: costCenterName,
        opening_type: openingType,
        yoe_min: yoeMin,
        yoe_max: yoeMax,
        min_education_level: minEducationLevel,
        employer_notes: employerNotes || undefined,
        remote_timezones:
          remoteTimezones.length > 0 ? remoteTimezones : undefined,
        location_titles:
          selectedLocations.length > 0 ? selectedLocations : undefined,
      };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/create-opening`,
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
        const errorData = await response.json();
        throw new Error(errorData.message || t("openings.createError"));
      }

      router.push("/openings");
    } catch (err) {
      setError(err instanceof Error ? err.message : t("openings.createError"));
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Container maxWidth="md">
      <Paper sx={{ p: 4 }}>
        <Box sx={{ mb: 4 }}>
          <Typography variant="h4" component="h1" gutterBottom>
            {t("openings.createTitle")}
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
            label={t("openings.openingTitle")}
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            inputProps={{ minLength: 3, maxLength: 32 }}
            error={title.length > 0 && (title.length < 3 || title.length > 32)}
            helperText={
              title.length > 0 && (title.length < 3 || title.length > 32)
                ? t("validation.title.length.3.32")
                : ""
            }
          />

          <TextField
            margin="normal"
            required
            fullWidth
            type="number"
            label={t("openings.positions")}
            value={positions}
            onChange={(e) => setPositions(parseInt(e.target.value, 10))}
            inputProps={{ min: 1, max: 20 }}
            error={positions < 1 || positions > 20}
            helperText={
              positions < 1 || positions > 20
                ? t("validation.positions.range.1.20")
                : ""
            }
          />

          <TextField
            margin="normal"
            required
            fullWidth
            multiline
            rows={4}
            label={t("openings.jobDescription")}
            value={jd}
            onChange={(e) => setJd(e.target.value)}
            inputProps={{ minLength: 10, maxLength: 1024 }}
            error={jd.length > 0 && (jd.length < 10 || jd.length > 1024)}
            helperText={
              jd.length > 0 && (jd.length < 10 || jd.length > 1024)
                ? t("validation.jobDescription.length.10.1024")
                : ""
            }
          />

          <Autocomplete
            options={orgUsers}
            getOptionLabel={(option) => `${option.name} (${option.email})`}
            value={recruiter}
            onChange={(_, newValue) => setRecruiter(newValue)}
            onInputChange={(_, value) => fetchOrgUsers(value)}
            renderInput={(params) => (
              <TextField
                {...params}
                required
                margin="normal"
                label={t("openings.recruiter")}
              />
            )}
          />

          <Autocomplete
            options={orgUsers}
            getOptionLabel={(option) => `${option.name} (${option.email})`}
            value={hiringManager}
            onChange={(_, newValue) => setHiringManager(newValue)}
            onInputChange={(_, value) => fetchOrgUsers(value)}
            renderInput={(params) => (
              <TextField
                {...params}
                required
                margin="normal"
                label={t("openings.hiringManager")}
              />
            )}
          />

          <Autocomplete
            options={costCenters}
            value={costCenterName}
            onChange={(_, newValue) => setCostCenterName(newValue || "")}
            renderInput={(params) => (
              <TextField
                {...params}
                required
                margin="normal"
                label={t("openings.costCenter")}
              />
            )}
          />

          <FormControl fullWidth margin="normal" required>
            <InputLabel>{t("openings.type")}</InputLabel>
            <Select
              value={openingType}
              onChange={(e) => setOpeningType(e.target.value as OpeningType)}
              label={t("openings.type")}
            >
              {Object.values(OpeningTypes).map((type) => (
                <MenuItem key={type} value={type}>
                  {t(`openings.types.${type}`)}
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          <Box sx={{ display: "flex", gap: 2, mt: 2 }}>
            <TextField
              margin="normal"
              required
              fullWidth
              type="number"
              label={t("openings.minYoe")}
              value={yoeMin}
              onChange={(e) => setYoeMin(parseInt(e.target.value, 10))}
              inputProps={{ min: 0, max: 100 }}
            />

            <TextField
              margin="normal"
              required
              fullWidth
              type="number"
              label={t("openings.maxYoe")}
              value={yoeMax}
              onChange={(e) => setYoeMax(parseInt(e.target.value, 10))}
              inputProps={{ min: 1, max: 100 }}
            />
          </Box>

          <FormControl fullWidth margin="normal">
            <InputLabel>{t("openings.minEducation")}</InputLabel>
            <Select
              value={minEducationLevel}
              onChange={(e) =>
                setMinEducationLevel(e.target.value as EducationLevel)
              }
              label={t("openings.minEducation")}
            >
              {Object.values(EducationLevels).map((level) => (
                <MenuItem key={level} value={level}>
                  {t(`openings.education.${level}`)}
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          <Autocomplete
            multiple
            options={Array.from(validTimezones)}
            value={remoteTimezones}
            onChange={(_, newValue) => setRemoteTimezones(newValue)}
            renderInput={(params) => (
              <TextField
                {...params}
                margin="normal"
                label={t("openings.remoteTimezones")}
                helperText={t("openings.remoteTimezonesHelp")}
              />
            )}
            renderTags={(value, getTagProps) =>
              value.map((option, index) => (
                <Chip label={option} {...getTagProps({ index })} key={option} />
              ))
            }
          />

          <Autocomplete
            multiple
            options={locations
              .filter((loc) => loc.state === LocationStates.ACTIVE)
              .map((loc) => loc.title)}
            value={selectedLocations}
            onChange={(_, newValue) => setSelectedLocations(newValue)}
            renderInput={(params) => (
              <TextField
                {...params}
                margin="normal"
                label={t("openings.officeLocations")}
                helperText={t("openings.officeLocationsHelp")}
              />
            )}
            renderTags={(value, getTagProps) =>
              value.map((option, index) => (
                <Chip label={option} {...getTagProps({ index })} key={option} />
              ))
            }
          />

          <TextField
            margin="normal"
            fullWidth
            multiline
            rows={4}
            label={t("openings.employerNotes")}
            value={employerNotes}
            onChange={(e) => setEmployerNotes(e.target.value)}
            inputProps={{ maxLength: 1024 }}
            error={employerNotes.length > 1024}
            helperText={
              employerNotes.length > 1024
                ? t("validation.employerNotes.maxLength.1024")
                : ""
            }
          />

          <Box sx={{ mt: 4, display: "flex", gap: 2 }}>
            <Button variant="outlined" onClick={() => router.push("/openings")}>
              {t("common.cancel")}
            </Button>
            <Button
              variant="contained"
              onClick={handleSave}
              disabled={
                isLoading ||
                !title ||
                title.length < 3 ||
                title.length > 32 ||
                !positions ||
                positions < 1 ||
                positions > 20 ||
                !jd ||
                jd.length < 10 ||
                jd.length > 1024 ||
                !recruiter ||
                !hiringManager ||
                !costCenterName ||
                !openingType ||
                yoeMin < 0 ||
                yoeMax <= yoeMin ||
                yoeMin > 100 ||
                yoeMax > 100
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
