"use client";

import { useState } from "react";
import Box from "@mui/material/Box";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import Typography from "@mui/material/Typography";
import Paper from "@mui/material/Paper";
import SearchIcon from "@mui/icons-material/Search";
import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import Select from "@mui/material/Select";
import MenuItem from "@mui/material/MenuItem";
import FormControl from "@mui/material/FormControl";
import InputLabel from "@mui/material/InputLabel";
import { FindHubOpeningsRequest, HubOpening } from "@psankar/vetchi-typespec";
import countries from "@psankar/vetchi-typespec/common/countries.json";
import Cookies from "js-cookie";

interface Country {
  country_code: string;
  en: string;
}

export default function FindOpeningsPage() {
  const [searchQuery, setSearchQuery] = useState("");
  const [countryCode, setCountryCode] = useState("");
  const [searchResults, setSearchResults] = useState<HubOpening[]>([]);
  const [error, setError] = useState<string | null>(null);

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSearchResults([]); // Clear results before new search

    const token = Cookies.get("session_token");
    if (!token) {
      setError("Not authenticated. Please log in again.");
      return;
    }

    const request: FindHubOpeningsRequest = {
      country_code: countryCode || undefined,
      limit: 40,
    };

    try {
      const response = await fetch("/api/hub/find-openings", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ request }),
      });

      if (!response.ok) {
        if (response.status === 401) {
          setError("Session expired. Please log in again.");
          // Clear the invalid session token
          Cookies.remove("session_token", { path: "/" });
          return;
        }
        throw new Error(`Failed to fetch openings: ${response.statusText}`);
      }

      const data = await response.json();

      // Safely check if data exists and has openings property
      if (data && typeof data === "object") {
        const openings = data.openings || [];
        setSearchResults(Array.isArray(openings) ? openings : []);
      } else {
        console.error("Invalid response format:", data);
        setError("Received invalid data format from server");
      }
    } catch (error) {
      console.error("Error searching openings:", error);
      setError("Failed to fetch openings. Please try again.");
    }
  };

  return (
    <AuthenticatedLayout>
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
            <TextField
              fullWidth
              variant="outlined"
              placeholder="Search for job titles, skills, or keywords"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              InputProps={{
                endAdornment: (
                  <Button
                    variant="contained"
                    sx={{ ml: 1 }}
                    type="submit"
                    startIcon={<SearchIcon />}
                  >
                    Search
                  </Button>
                ),
              }}
            />
          </Box>
        </Paper>
        {/* Search results will be displayed here */}
        <Box sx={{ mt: 4 }}>
          {searchResults.length > 0 ? (
            searchResults.map((opening) => (
              <Paper
                key={opening.opening_id_within_company}
                sx={{ p: 2, mb: 2 }}
              >
                <Typography variant="h6">{opening.job_title}</Typography>
                <Typography variant="subtitle1">
                  {opening.company_name}
                </Typography>
                <Typography variant="body1">{opening.jd}</Typography>
              </Paper>
            ))
          ) : (
            <Typography variant="body1" color="text.secondary" align="center">
              No openings found
            </Typography>
          )}
        </Box>
      </Box>
    </AuthenticatedLayout>
  );
}
