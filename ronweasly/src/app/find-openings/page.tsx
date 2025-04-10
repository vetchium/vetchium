"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import { config } from "@/config";
import AddIcon from "@mui/icons-material/Add";
import LocalOfferIcon from "@mui/icons-material/LocalOffer";
import OpenInNewIcon from "@mui/icons-material/OpenInNew";
import SearchIcon from "@mui/icons-material/Search";
import Autocomplete, { createFilterOptions } from "@mui/material/Autocomplete";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import CircularProgress from "@mui/material/CircularProgress";
import FormControl from "@mui/material/FormControl";
import IconButton from "@mui/material/IconButton";
import InputLabel from "@mui/material/InputLabel";
import MenuItem from "@mui/material/MenuItem";
import Paper from "@mui/material/Paper";
import Select from "@mui/material/Select";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import { FindHubOpeningsRequest, HubOpening, VTag } from "@vetchium/typespec";
import countries from "@vetchium/typespec/common/countries.json";
import Cookies from "js-cookie";
import { useRouter, useSearchParams } from "next/navigation";
import { Suspense, useEffect, useState } from "react";

interface Country {
  country_code: string;
  en: string;
}

// Interface for tags including free input
interface TagOption extends VTag {
  inputValue?: string;
}

// Filter configuration for Autocomplete
const filter = createFilterOptions<TagOption>();

// Cache for search results
let searchCache: {
  results: HubOpening[];
  countryCode: string;
  searchQuery: string;
} | null = null;

function FindOpeningsContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [searchQuery, setSearchQuery] = useState("");
  const [countryCode, setCountryCode] = useState("");
  const [searchResults, setSearchResults] = useState<HubOpening[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [tagSuggestions, setTagSuggestions] = useState<VTag[]>([]);
  const [selectedTags, setSelectedTags] = useState<VTag[]>([]);

  // Fetch tag suggestions when user types
  useEffect(() => {
    const fetchTags = async () => {
      if (searchQuery.length >= 2) {
        const token = Cookies.get("session_token");
        if (!token) return;

        try {
          const response = await fetch(
            `${config.API_SERVER_PREFIX}/hub/filter-vtags`,
            {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${token}`,
              },
              body: JSON.stringify({
                prefix: searchQuery,
              }),
            }
          );

          if (!response.ok) {
            if (response.status === 401) {
              setError("Session expired. Please log in again.");
              Cookies.remove("session_token", { path: "/" });
              return;
            }
            throw new Error(`Failed to fetch tags: ${response.statusText}`);
          }

          const data = await response.json();
          setTagSuggestions(Array.isArray(data) ? data : []);
        } catch (error) {
          console.error("Error fetching tags:", error);
          setTagSuggestions([]);
        }
      } else {
        setTagSuggestions([]);
      }
    };

    const debounceTimer = setTimeout(fetchTags, 300);
    return () => clearTimeout(debounceTimer);
  }, [searchQuery]);

  // Restore state from cache on mount
  useEffect(() => {
    if (searchCache) {
      setSearchResults(searchCache.results);
      setCountryCode(searchCache.countryCode);
      setSearchQuery(searchCache.searchQuery);
    }
  }, []);

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    const token = Cookies.get("session_token");
    if (!token) {
      setError("Session expired. Please log in again.");
      return;
    }

    const request: FindHubOpeningsRequest = {
      country_code: countryCode || "",
      terms: searchQuery ? [searchQuery] : undefined,
      tags: selectedTags.map((tag) => tag.id),
    };

    try {
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/find-openings`,
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
          setError("Session expired. Please log in again.");
          Cookies.remove("session_token", { path: "/" });
          return;
        }
        throw new Error(`Failed to fetch openings: ${response.statusText}`);
      }

      const data = await response.json();

      if (Array.isArray(data)) {
        setSearchResults(data);
        // Update cache
        searchCache = {
          results: data,
          countryCode,
          searchQuery,
        };
      } else {
        console.error("Invalid response format:", data);
        setError("Received invalid data format from server");
      }
    } catch (error) {
      console.error("Error searching openings:", error);
      setError("Failed to fetch openings. Please try again.");
    }
  };

  const handleOpeningClick = (opening: HubOpening, newTab?: boolean) => {
    const url = `/org/${opening.company_domain}/opening/${opening.opening_id_within_company}`;
    if (newTab) {
      window.open(url, "_blank");
    } else {
      router.push(url);
    }
  };

  const handleOpeningMouseDown = (e: React.MouseEvent, opening: HubOpening) => {
    // Middle click
    if (e.button === 1) {
      e.preventDefault();
      handleOpeningClick(opening, true);
    }
  };

  return (
    <Box sx={{ maxWidth: 800, mx: "auto", mt: 4 }}>
      <Typography variant="h4" gutterBottom align="center">
        Find Openings
      </Typography>
      {error && (
        <Paper sx={{ p: 2, mb: 2, bgcolor: "error.light" }}>
          <Typography color="error" align="center">
            {error}
          </Typography>
        </Paper>
      )}
      <Paper
        component="form"
        onSubmit={handleSearch}
        sx={{
          p: 4,
          mt: 4,
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
        }}
      >
        <Typography
          variant="body1"
          color="text.secondary"
          gutterBottom
          align="center"
        >
          Search for job openings across all locations
        </Typography>
        <Box
          sx={{
            width: "100%",
            mt: 2,
            display: "flex",
            flexDirection: "column",
            gap: 2,
          }}
        >
          <FormControl fullWidth>
            <InputLabel id="country-select-label">Country</InputLabel>
            <Select
              labelId="country-select-label"
              id="country-select"
              value={countryCode}
              label="Country"
              onChange={(e) => setCountryCode(e.target.value)}
            >
              <MenuItem value="">
                <em>All Countries</em>
              </MenuItem>
              {countries.map((country: Country) => (
                <MenuItem
                  key={country.country_code}
                  value={country.country_code}
                >
                  {country.en}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
          <Autocomplete
            multiple
            freeSolo
            id="tags-search"
            options={tagSuggestions}
            value={selectedTags}
            onChange={(_, newValue) => {
              const processedValues = newValue.map((option) => {
                if (typeof option === "string") {
                  // For free solo values, create a new tag option
                  return {
                    name: option,
                    id: "", // Empty ID indicates it's a new tag
                  };
                }
                return option;
              });
              setSelectedTags(processedValues);
            }}
            getOptionLabel={(option) => {
              if (typeof option === "string") {
                return option;
              }
              return option.name;
            }}
            filterOptions={(options, params) => {
              const filtered = filter(options, params);
              const { inputValue } = params;

              // Only suggest creating a new tag if it's not already in suggestions
              // and not already selected
              const isExisting = options.some(
                (option) => option.name === inputValue
              );
              const isSelected = selectedTags.some(
                (tag) => tag.name === inputValue
              );

              if (inputValue !== "" && !isExisting && !isSelected) {
                filtered.push({
                  name: inputValue,
                  id: "", // Empty ID indicates it's a new tag
                });
              }

              return filtered;
            }}
            renderOption={(props, option) => {
              const { key, ...otherProps } = props;
              return (
                <li key={key} {...otherProps}>
                  {!option.id ? (
                    <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                      <AddIcon fontSize="small" />
                      <span>{option.name}</span>
                    </Box>
                  ) : (
                    <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                      <LocalOfferIcon fontSize="small" />
                      <span>{option.name}</span>
                    </Box>
                  )}
                </li>
              );
            }}
            renderInput={(params) => (
              <TextField
                {...params}
                variant="outlined"
                placeholder="Search for job titles, skills, or keywords"
                onChange={(e) => setSearchQuery(e.target.value)}
                InputProps={{
                  ...params.InputProps,
                  endAdornment: (
                    <>
                      {params.InputProps.endAdornment}
                      <Button
                        variant="contained"
                        sx={{ ml: 1 }}
                        type="submit"
                        startIcon={<SearchIcon />}
                      >
                        Search
                      </Button>
                    </>
                  ),
                }}
              />
            )}
          />
        </Box>
      </Paper>
      <Box sx={{ mt: 4 }}>
        {searchResults.length > 0 ? (
          searchResults.map((opening) => (
            <Paper
              key={`${opening.company_domain}-${opening.opening_id_within_company}`}
              sx={{
                p: 2,
                mb: 2,
                cursor: "pointer",
                "&:hover": { bgcolor: "action.hover" },
              }}
              onClick={() => handleOpeningClick(opening)}
              onMouseDown={(e) => handleOpeningMouseDown(e, opening)}
            >
              <Box
                sx={{
                  display: "flex",
                  justifyContent: "space-between",
                  alignItems: "flex-start",
                }}
              >
                <Box>
                  <Typography variant="h6">{opening.job_title}</Typography>
                  <Typography variant="subtitle1">
                    {opening.company_name}
                  </Typography>
                  <Typography variant="body1">{opening.jd}</Typography>
                </Box>
                <IconButton
                  onClick={(e) => {
                    e.stopPropagation();
                    handleOpeningClick(opening, true);
                  }}
                  size="small"
                  sx={{ ml: 1 }}
                  title="Open in new tab"
                >
                  <OpenInNewIcon />
                </IconButton>
              </Box>
            </Paper>
          ))
        ) : (
          <Typography variant="body1" color="text.secondary" align="center">
            No openings found
          </Typography>
        )}
      </Box>
    </Box>
  );
}

export default function FindOpeningsPage() {
  return (
    <AuthenticatedLayout>
      <Suspense
        fallback={
          <Box
            sx={{
              display: "flex",
              justifyContent: "center",
              alignItems: "center",
              minHeight: "50vh",
            }}
          >
            <CircularProgress />
          </Box>
        }
      >
        <FindOpeningsContent />
      </Suspense>
    </AuthenticatedLayout>
  );
}
