"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import Autocomplete from "@mui/material/Autocomplete";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Chip from "@mui/material/Chip";
import CircularProgress from "@mui/material/CircularProgress";
import Paper from "@mui/material/Paper";
import Stack from "@mui/material/Stack";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import {
  ApplyForOpeningRequest,
  EducationLevel,
  EducationLevels,
  GetHubOpeningDetailsRequest,
  HubOpeningDetails,
  HubUserShort,
  OpeningState,
  OpeningStates,
  OpeningType,
  OpeningTypes,
} from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";

const formatEducationLevel = (
  level: EducationLevel,
  t: (key: string) => string
) => {
  switch (level) {
    case EducationLevels.BACHELOR:
      return t("openingDetails.educationLevel.bachelor");
    case EducationLevels.MASTER:
      return t("openingDetails.educationLevel.master");
    case EducationLevels.DOCTORATE:
      return t("openingDetails.educationLevel.doctorate");
    case EducationLevels.NOT_MATTERS:
      return t("openingDetails.educationLevel.notMatters");
    case EducationLevels.UNSPECIFIED:
      return t("openingDetails.educationLevel.unspecified");
    default:
      return level;
  }
};

const formatOpeningType = (type: OpeningType, t: (key: string) => string) => {
  switch (type) {
    case OpeningTypes.FULL_TIME:
      return t("openingDetails.openingType.fullTime");
    case OpeningTypes.PART_TIME:
      return t("openingDetails.openingType.partTime");
    case OpeningTypes.CONTRACT:
      return t("openingDetails.openingType.contract");
    case OpeningTypes.INTERNSHIP:
      return t("openingDetails.openingType.internship");
    case OpeningTypes.UNSPECIFIED:
      return t("openingDetails.openingType.unspecified");
    default:
      return type;
  }
};

export default function OpeningDetailsPage() {
  const { t } = useTranslation();
  const params = useParams();
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [opening, setOpening] = useState<HubOpeningDetails | null>(null);
  const [resumeFile, setResumeFile] = useState<File | null>(null);
  const [uploading, setUploading] = useState(false);
  const [selectedEndorsers, setSelectedEndorsers] = useState<HubUserShort[]>(
    []
  );
  const [colleagueSearchInput, setColleagueSearchInput] = useState("");
  const [colleagueOptions, setColleagueOptions] = useState<HubUserShort[]>([]);
  const [loadingColleagues, setLoadingColleagues] = useState(false);
  const [isPopupOpen, setIsPopupOpen] = useState(false);

  if (!params?.domain || !params?.openingId) {
    return (
      <Box sx={{ p: 3 }}>
        <Typography color="error">{t("common.error.invalidParams")}</Typography>
        <Button
          variant="contained"
          onClick={() => router.back()}
          sx={{ mt: 2 }}
        >
          {t("common.back")}
        </Button>
      </Box>
    );
  }

  const companyDomain = params.domain as string;
  const openingId = params.openingId as string;

  useEffect(() => {
    const fetchOpeningDetails = async () => {
      const token = Cookies.get("session_token");
      if (!token) {
        setError(t("common.error.notAuthenticated"));
        return;
      }

      try {
        const request: GetHubOpeningDetailsRequest = {
          company_domain: companyDomain,
          opening_id_within_company: openingId,
        };

        const response = await fetch(
          `${config.API_SERVER_PREFIX}/hub/get-opening-details`,
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
          if (response.status === 401) {
            setError(t("common.error.sessionExpired"));
            Cookies.remove("session_token", { path: "/" });
            return;
          }
          throw new Error(
            `Failed to fetch opening details: ${response.statusText}`
          );
        }

        const data = await response.json();
        setOpening(data);
      } catch (error) {
        console.error("Error fetching opening details:", error);
        setError(t("openingDetails.error.loadFailed"));
      } finally {
        setLoading(false);
      }
    };

    fetchOpeningDetails();
  }, [companyDomain, openingId]);

  // New function to search for colleagues
  const searchColleagues = async (prefix: string) => {
    if (!prefix || prefix.length < 1) {
      setColleagueOptions([]);
      return;
    }

    setLoadingColleagues(true);
    const token = Cookies.get("session_token");
    if (!token) {
      setError(t("common.error.notAuthenticated"));
      return;
    }

    try {
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/filter-colleagues`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            prefix: prefix,
            limit: 6,
          }),
        }
      );

      if (!response.ok) {
        if (response.status === 401) {
          setError(t("common.error.sessionExpired"));
          Cookies.remove("session_token", { path: "/" });
          return;
        }
        throw new Error(`Failed to fetch colleagues: ${response.statusText}`);
      }

      const data = await response.json();
      setColleagueOptions(data);
    } catch (error) {
      console.error("Error fetching colleagues:", error);
    } finally {
      setLoadingColleagues(false);
    }
  };

  // Debounce search input
  useEffect(() => {
    const timer = setTimeout(() => {
      if (colleagueSearchInput && colleagueSearchInput.length >= 1) {
        searchColleagues(colleagueSearchInput);
      }
    }, 300);

    return () => clearTimeout(timer);
  }, [colleagueSearchInput]);

  const handleEndorserAdd = (newEndorser: HubUserShort | null) => {
    if (newEndorser && selectedEndorsers.length < 5) {
      setSelectedEndorsers([...selectedEndorsers, newEndorser]);
      setColleagueSearchInput("");
      setColleagueOptions([]);
    }
  };

  const handleEndorserRemove = (handle: string) => {
    setSelectedEndorsers(
      selectedEndorsers.filter((endorser) => endorser.handle !== handle)
    );
  };

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files && event.target.files[0]) {
      const file = event.target.files[0];
      if (file.type !== "application/pdf") {
        setError(t("openingDetails.error.pdfOnly"));
        return;
      }
      if (file.size > 5 * 1024 * 1024) {
        // 5MB limit
        setError(t("openingDetails.error.fileTooLarge"));
        return;
      }
      setResumeFile(file);
      setError(null);
    }
  };

  const handleApply = async () => {
    if (!resumeFile) {
      setError(t("openingDetails.error.noResume"));
      return;
    }

    const token = Cookies.get("session_token");
    if (!token) {
      setError(t("common.error.notAuthenticated"));
      return;
    }

    setUploading(true);
    try {
      // Convert file to base64
      const reader = new FileReader();
      const base64Promise = new Promise<string>((resolve, reject) => {
        reader.onload = () => {
          const base64String = (reader.result as string).split(",")[1];
          resolve(base64String);
        };
        reader.onerror = reject;
      });
      reader.readAsDataURL(resumeFile);

      const base64Resume = await base64Promise;

      const request: ApplyForOpeningRequest = {
        opening_id_within_company: openingId,
        company_domain: companyDomain,
        resume: base64Resume,
        filename: resumeFile.name,
        endorser_handles: selectedEndorsers.map((endorser) => endorser.handle),
      };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/apply-for-opening`,
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
        if (response.status === 401) {
          setError(t("common.error.sessionExpired"));
          Cookies.remove("session_token", { path: "/" });
          return;
        }

        if (response.status === 422) {
          setError(t("openingDetails.cannotApply"));
          return;
        }

        const errorData = await response.json();
        if (response.status === 400 && errorData.errors) {
          setError(errorData.errors.join(", "));
        } else {
          setError(t("openingDetails.error.applyFailed"));
        }
        return;
      }

      // Success - redirect to applications page
      router.push("/my-applications");
    } catch (error) {
      console.error("Error applying for opening:", error);
      setError(t("openingDetails.error.applyFailed"));
    } finally {
      setUploading(false);
    }
  };

  const canApply = opening?.state === OpeningStates.ACTIVE;

  const getOpeningStateMessage = (state: OpeningState) => {
    switch (state) {
      case OpeningStates.DRAFT:
        return t("openingDetails.state.draft");
      case OpeningStates.SUSPENDED:
        return t("openingDetails.state.suspended");
      case OpeningStates.CLOSED:
        return t("openingDetails.state.closed");
      default:
        return null;
    }
  };

  if (loading) {
    return (
      <AuthenticatedLayout>
        <Box sx={{ display: "flex", justifyContent: "center", mt: 4 }}>
          <CircularProgress />
        </Box>
      </AuthenticatedLayout>
    );
  }

  if (error) {
    return (
      <AuthenticatedLayout>
        <Paper sx={{ p: 2, mb: 2, bgcolor: "error.light" }}>
          <Typography color="error" align="center">
            {error}
          </Typography>
        </Paper>
      </AuthenticatedLayout>
    );
  }

  if (!opening) {
    return (
      <AuthenticatedLayout>
        <Paper sx={{ p: 2, mb: 2 }}>
          <Typography align="center">{t("openingDetails.notFound")}</Typography>
        </Paper>
      </AuthenticatedLayout>
    );
  }

  return (
    <AuthenticatedLayout>
      <Box sx={{ maxWidth: 800, mx: "auto", mt: 4 }}>
        <Paper sx={{ p: 4 }}>
          <Typography variant="h4" gutterBottom>
            {opening.job_title}
          </Typography>
          <Typography variant="h6" color="text.secondary" gutterBottom>
            {opening.company_name}
          </Typography>

          <Stack direction="row" spacing={1} sx={{ mb: 3 }}>
            {opening.opening_type && (
              <Chip label={formatOpeningType(opening.opening_type, t)} />
            )}
            {opening.education_level && (
              <Chip label={formatEducationLevel(opening.education_level, t)} />
            )}
            {opening.yoe_min !== undefined && opening.yoe_max !== undefined && (
              <Chip
                label={t("openingDetails.yearsExperience", {
                  min: opening.yoe_min,
                  max: opening.yoe_max,
                })}
              />
            )}
          </Stack>

          <Typography variant="body1" paragraph>
            {opening.jd}
          </Typography>

          {opening.hiring_manager_name && (
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
              {t("openingDetails.hiringManager")}: {opening.hiring_manager_name}
            </Typography>
          )}

          <Box sx={{ mt: 4 }}>
            {canApply ? (
              opening.is_appliable ? (
                <>
                  <input
                    accept="application/pdf"
                    style={{ display: "none" }}
                    id="resume-file"
                    type="file"
                    onChange={handleFileChange}
                  />
                  <label htmlFor="resume-file">
                    <Button
                      component="span"
                      variant="outlined"
                      color="primary"
                      fullWidth
                      sx={{ mb: 2 }}
                      disabled={uploading}
                    >
                      {resumeFile
                        ? t("openingDetails.resumeSelected", {
                            name: resumeFile.name,
                          })
                        : t("openingDetails.selectResume")}
                    </Button>
                  </label>

                  {/* Endorsers Section */}
                  <Box sx={{ mb: 3 }}>
                    <Typography variant="h6" gutterBottom>
                      {t("openingDetails.endorsers.title")}
                    </Typography>
                    <Typography
                      variant="body2"
                      color="text.secondary"
                      sx={{ mb: 2 }}
                    >
                      {t("openingDetails.endorsers.description")}
                    </Typography>

                    <Autocomplete
                      multiple
                      id="endorsers-autocomplete"
                      options={colleagueOptions}
                      value={selectedEndorsers}
                      loading={loadingColleagues}
                      disabled={uploading}
                      inputValue={colleagueSearchInput}
                      onInputChange={(event, newValue, reason) => {
                        if (selectedEndorsers.length < 5) {
                          setColleagueSearchInput(newValue);
                        }
                      }}
                      getOptionLabel={(option) =>
                        `${option.name} (@${option.handle})`
                      }
                      isOptionEqualToValue={(option, value) =>
                        option.handle === value.handle
                      }
                      open={
                        isPopupOpen &&
                        selectedEndorsers.length < 5 &&
                        colleagueOptions.length > 0
                      }
                      onOpen={() => {
                        if (selectedEndorsers.length < 5) {
                          setIsPopupOpen(true);
                        }
                      }}
                      onClose={() => setIsPopupOpen(false)}
                      renderInput={(params) => (
                        <TextField
                          {...params}
                          label={t("openingDetails.endorsers.search")}
                          variant="outlined"
                          InputProps={{
                            ...params.InputProps,
                            endAdornment: (
                              <>
                                {loadingColleagues ? (
                                  <CircularProgress color="inherit" size={20} />
                                ) : null}
                                {params.InputProps.endAdornment}
                              </>
                            ),
                          }}
                          helperText={
                            selectedEndorsers.length >= 5
                              ? t("openingDetails.endorsers.maxReached")
                              : t("openingDetails.endorsers.remaining", {
                                  count: 5 - selectedEndorsers.length,
                                })
                          }
                        />
                      )}
                      onChange={(_, value) => {
                        if (value.length <= 5) {
                          setSelectedEndorsers(value);
                          setColleagueSearchInput("");
                          setIsPopupOpen(false);
                        }
                      }}
                      sx={{ mb: 2 }}
                    />
                  </Box>

                  <Button
                    variant="contained"
                    color="primary"
                    size="large"
                    onClick={handleApply}
                    fullWidth
                    disabled={!resumeFile || uploading}
                  >
                    {uploading ? (
                      <CircularProgress size={24} color="inherit" />
                    ) : (
                      t("openingDetails.apply")
                    )}
                  </Button>
                </>
              ) : (
                <Typography
                  variant="body1"
                  color="error"
                  align="center"
                  sx={{ mb: 2 }}
                >
                  {t("openingDetails.cannotApply")}
                </Typography>
              )
            ) : (
              <Typography
                variant="body1"
                color="error"
                align="center"
                sx={{ mb: 2 }}
              >
                {getOpeningStateMessage(opening.state)}
              </Typography>
            )}
          </Box>
        </Paper>
      </Box>
    </AuthenticatedLayout>
  );
}
