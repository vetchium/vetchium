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
  FormControlLabel,
  Switch,
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
  GlobalCountryCode,
  OpeningTag,
  OpeningTagID,
} from "@psankar/vetchi-typespec";
import { Location, LocationStates } from "@psankar/vetchi-typespec";
import countries from "@psankar/vetchi-typespec/common/countries.json";
import { FeatureFlags } from "@/config/features";

interface Country {
  country_code: string;
  en: string;
}

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
  const [remoteCountries, setRemoteCountries] = useState<string[]>([]);
  const [isGloballyRemote, setIsGloballyRemote] = useState(false);

  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [orgUsers, setOrgUsers] = useState<OrgUserShort[]>([]);
  const [costCenters, setCostCenters] = useState<string[]>([]);
  const [locations, setLocations] = useState<Location[]>([]);
  const [selectedTags, setSelectedTags] = useState<OpeningTag[]>([]);
  const [newTags, setNewTags] = useState<string[]>([]);
  const [availableTags, setAvailableTags] = useState<OpeningTag[]>([]);
  const [tagSearchQuery, setTagSearchQuery] = useState("");

  const router = useRouter();
  const { t } = useTranslation();

  // Initialize from sessionStorage on mount
  useEffect(() => {
    const saved = sessionStorage.getItem("isGloballyRemote");
    if (saved !== null) {
      setIsGloballyRemote(JSON.parse(saved));
    }
  }, []);

  // Persist isGloballyRemote state to sessionStorage
  useEffect(() => {
    sessionStorage.setItem(
      "isGloballyRemote",
      JSON.stringify(isGloballyRemote)
    );
  }, [isGloballyRemote]);

  useEffect(() => {
    fetchOrgUsers();
    fetchCostCenters();
    fetchLocations();
    fetchTags();
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

  const fetchTags = async (searchPrefix?: string) => {
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/filter-opening-tags`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            prefix: searchPrefix,
          }),
        }
      );

      if (!response.ok) {
        throw new Error(t("openings.fetchTagsError"));
      }

      const data = await response.json();
      setAvailableTags(data || []);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("openings.fetchTagsError")
      );
    }
  };

  const handleTagSearch = (query: string) => {
    setTagSearchQuery(query);
    fetchTags(query);
  };

  const handleAddNewTag = (newTag: string) => {
    if (selectedTags.length + newTags.length >= 3) {
      setError(t("openings.maxTagsError"));
      return;
    }
    if (newTag.trim() && !newTags.includes(newTag.trim())) {
      setNewTags([...newTags, newTag.trim()]);
    }
  };

  const handleRemoveNewTag = (tagToRemove: string) => {
    setNewTags(newTags.filter((tag) => tag !== tagToRemove));
  };

  const handleRemoveSelectedTag = (tagToRemove: OpeningTag) => {
    setSelectedTags(selectedTags.filter((tag) => tag.id !== tagToRemove.id));
  };

  const handleTagSelect = (tag: OpeningTag | null) => {
    if (tag && !selectedTags.find((t) => t.id === tag.id)) {
      if (selectedTags.length + newTags.length >= 3) {
        setError(t("openings.maxTagsError"));
        return;
      }
      setSelectedTags([...selectedTags, tag]);
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

      if (
        !isGloballyRemote &&
        selectedLocations.length === 0 &&
        remoteTimezones.length === 0 &&
        remoteCountries.length === 0
      ) {
        setError(t("openings.locationRequiredError"));
        return;
      }

      if (selectedTags.length + newTags.length === 0) {
        setError(t("openings.tagsRequiredError"));
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
        remote_country_codes: isGloballyRemote
          ? [GlobalCountryCode]
          : remoteCountries.length > 0
          ? remoteCountries
          : undefined,
        tags: selectedTags.map((tag) => tag.id),
        new_tags: newTags.length > 0 ? newTags : undefined,
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

          <FormControlLabel
            control={
              <Switch
                checked={isGloballyRemote}
                onChange={(e) => {
                  setIsGloballyRemote(e.target.checked);
                  if (e.target.checked) {
                    setRemoteCountries([]);
                    setRemoteTimezones([]);
                  }
                }}
                color="primary"
              />
            }
            label={t("openings.globallyRemote")}
            sx={{ mt: 2, mb: 1, display: "block" }}
          />

          <Autocomplete
            multiple
            options={countries}
            getOptionLabel={(option) => option.en}
            value={countries.filter((c) =>
              remoteCountries.includes(c.country_code)
            )}
            onChange={(_, newValue) =>
              setRemoteCountries(newValue.map((c) => c.country_code))
            }
            disabled={isGloballyRemote}
            renderInput={(params) => (
              <TextField
                {...params}
                margin="normal"
                label={t("openings.remoteCountries")}
                helperText={t("openings.remoteCountriesHelp")}
              />
            )}
            renderTags={(value, getTagProps) =>
              value.map((option, index) => (
                <Chip
                  label={option.en}
                  {...getTagProps({ index })}
                  key={option.country_code}
                />
              ))
            }
          />
          {FeatureFlags.TimezonesForCreateOpening && (
            <Autocomplete
              multiple
              options={Array.from(validTimezones)}
              value={remoteTimezones}
              onChange={(_, newValue) => setRemoteTimezones(newValue)}
              disabled={isGloballyRemote}
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
                  <Chip
                    label={option}
                    {...getTagProps({ index })}
                    key={option}
                  />
                ))
              }
            />
          )}
          <Autocomplete
            multiple
            options={locations
              .filter((loc) => loc.state === LocationStates.ACTIVE)
              .map((loc) => loc.title)}
            value={selectedLocations}
            onChange={(_, newValue) => setSelectedLocations(newValue)}
            disabled={isGloballyRemote}
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

          <Box sx={{ mt: 2 }}>
            <Typography variant="subtitle1" gutterBottom>
              {t("openings.tags")}
            </Typography>
            <Autocomplete
              options={availableTags}
              getOptionLabel={(option) => option.name}
              value={null}
              onChange={(_, newValue) => handleTagSelect(newValue)}
              onInputChange={(_, value) => handleTagSearch(value)}
              renderInput={(params) => (
                <TextField
                  {...params}
                  label={t("openings.selectTags")}
                  helperText={t("openings.tagsHelp")}
                />
              )}
            />
            <Box sx={{ mt: 1, mb: 2 }}>
              {selectedTags.map((tag) => (
                <Chip
                  key={tag.id}
                  label={tag.name}
                  onDelete={() => handleRemoveSelectedTag(tag)}
                  sx={{ m: 0.5 }}
                />
              ))}
            </Box>

            <Box sx={{ mt: 2 }}>
              <TextField
                fullWidth
                label={t("openings.addNewTag")}
                value={tagSearchQuery}
                onChange={(e) => setTagSearchQuery(e.target.value)}
                InputProps={{
                  endAdornment: (
                    <Button
                      onClick={() => handleAddNewTag(tagSearchQuery)}
                      disabled={
                        !tagSearchQuery.trim() ||
                        selectedTags.length + newTags.length >= 3
                      }
                    >
                      {t("common.add")}
                    </Button>
                  ),
                }}
                helperText={t("openings.newTagHelp")}
              />
            </Box>
            <Box sx={{ mt: 1 }}>
              {newTags.map((tag) => (
                <Chip
                  key={tag}
                  label={tag}
                  onDelete={() => handleRemoveNewTag(tag)}
                  color="primary"
                  variant="outlined"
                  sx={{ m: 0.5 }}
                />
              ))}
            </Box>
          </Box>

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
                yoeMax > 100 ||
                (!isGloballyRemote &&
                  selectedLocations.length === 0 &&
                  remoteTimezones.length === 0 &&
                  remoteCountries.length === 0) ||
                selectedTags.length + newTags.length === 0
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
