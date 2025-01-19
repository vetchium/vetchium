"use client";

import {
  Box,
  Button,
  Container,
  Alert,
  Typography,
  FormControlLabel,
  Switch,
  Grid,
  Paper,
  IconButton,
  Chip,
  Link,
} from "@mui/material";
import {
  Edit as EditIcon,
  Delete as DeleteIcon,
  Map as MapIcon,
} from "@mui/icons-material";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useTranslation } from "@/hooks/useTranslation";
import { config } from "@/config";
import Cookies from "js-cookie";
import {
  Location,
  GetLocationsRequest,
  DefunctLocationRequest,
  LocationStates,
} from "@psankar/vetchi-typespec";

export default function LocationsPage() {
  const [locations, setLocations] = useState<Location[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [includeDefunct, setIncludeDefunct] = useState(false);
  const router = useRouter();
  const { t } = useTranslation();

  const fetchLocations = async () => {
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const request: GetLocationsRequest = {
        states: includeDefunct
          ? [LocationStates.ACTIVE, LocationStates.DEFUNCT]
          : [LocationStates.ACTIVE],
      };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/get-locations`,
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
        throw new Error(t("locations.fetchError"));
      }

      const data = await response.json();
      setLocations(data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : t("locations.fetchError"));
    }
  };

  useEffect(() => {
    const savedValue = localStorage.getItem("includeDefunctLocations");
    if (savedValue) {
      setIncludeDefunct(savedValue === "true");
    }
  }, []);

  useEffect(() => {
    fetchLocations();
  }, [includeDefunct]);

  const handleDelete = async (location: Location) => {
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const request: DefunctLocationRequest = {
        title: location.title,
      };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/defunct-location`,
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
        throw new Error(t("locations.defunctError"));
      }

      fetchLocations();
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("locations.defunctError")
      );
    }
  };

  return (
    <Container maxWidth="lg">
      <Box
        sx={{
          mb: 4,
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
        }}
      >
        <Typography variant="h4" component="h1" sx={{ color: "text.primary" }}>
          {t("locations.title")}
        </Typography>
        <Box sx={{ display: "flex", gap: 2, alignItems: "center" }}>
          <FormControlLabel
            control={
              <Switch
                checked={includeDefunct}
                onChange={(e) => {
                  setIncludeDefunct(e.target.checked);
                  localStorage.setItem(
                    "includeDefunctLocations",
                    e.target.checked.toString()
                  );
                }}
              />
            }
            label={t("locations.includeDefunct")}
            sx={{ color: "text.primary" }}
          />
          <Button
            variant="contained"
            onClick={() => router.push("/locations/add")}
          >
            {t("locations.add")}
          </Button>
        </Box>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      {locations.length > 0 ? (
        <Grid container spacing={2}>
          {locations.map((location) => (
            <Grid item xs={12} md={6} lg={4} key={location.title}>
              <Paper
                sx={{
                  p: 2,
                  height: "100%",
                  opacity: location.state === "DEFUNCT_LOCATION" ? 0.7 : 1,
                  color: "text.primary",
                }}
              >
                <Box
                  sx={{
                    display: "flex",
                    justifyContent: "space-between",
                    alignItems: "flex-start",
                    mb: 2,
                  }}
                >
                  <Box>
                    <Typography
                      variant="h6"
                      gutterBottom
                      sx={{ color: "text.primary" }}
                    >
                      {location.title}
                    </Typography>
                    <Typography
                      variant="caption"
                      sx={{
                        backgroundColor:
                          location.state === "DEFUNCT_LOCATION"
                            ? "error.main"
                            : "success.main",
                        color: "white",
                        px: 1,
                        py: 0.5,
                        borderRadius: 1,
                      }}
                    >
                      {location.state === "DEFUNCT_LOCATION"
                        ? t("locations.defunct")
                        : t("locations.active")}
                    </Typography>
                  </Box>
                  <Box>
                    <IconButton
                      onClick={() =>
                        router.push(
                          `/locations/edit?title=${encodeURIComponent(
                            location.title
                          )}`
                        )
                      }
                      size="small"
                      sx={{ color: "text.primary" }}
                    >
                      <EditIcon />
                    </IconButton>
                    <IconButton
                      onClick={() => handleDelete(location)}
                      size="small"
                      disabled={location.state === "DEFUNCT_LOCATION"}
                      sx={{ color: "text.primary" }}
                    >
                      <DeleteIcon />
                    </IconButton>
                  </Box>
                </Box>

                <Box sx={{ mb: 1 }}>
                  <Typography
                    variant="subtitle2"
                    sx={{ color: "text.secondary", mb: 0.5 }}
                  >
                    {t("locations.countryCode")}
                  </Typography>
                  <Typography sx={{ color: "text.primary" }}>
                    {location.country_code}
                  </Typography>
                </Box>

                <Box sx={{ mb: 1 }}>
                  <Typography
                    variant="subtitle2"
                    sx={{ color: "text.secondary", mb: 0.5 }}
                  >
                    {t("locations.postalAddress")}
                  </Typography>
                  <Typography
                    sx={{ color: "text.primary", whiteSpace: "pre-line" }}
                  >
                    {location.postal_address}
                  </Typography>
                </Box>

                <Box sx={{ mb: 1 }}>
                  <Typography
                    variant="subtitle2"
                    sx={{ color: "text.secondary", mb: 0.5 }}
                  >
                    {t("locations.postalCode")}
                  </Typography>
                  <Typography sx={{ color: "text.primary" }}>
                    {location.postal_code}
                  </Typography>
                </Box>

                {location.city_aka && location.city_aka.length > 0 && (
                  <Box sx={{ mb: 1 }}>
                    <Typography
                      variant="subtitle2"
                      sx={{ color: "text.secondary", mb: 0.5 }}
                    >
                      {t("locations.cityAka")}
                    </Typography>
                    <Box sx={{ display: "flex", gap: 1, flexWrap: "wrap" }}>
                      {location.city_aka.map((city, index) => (
                        <Chip
                          key={index}
                          label={city}
                          size="small"
                          variant="outlined"
                          sx={{ color: "text.primary", borderColor: "divider" }}
                        />
                      ))}
                    </Box>
                  </Box>
                )}

                {location.openstreetmap_url && (
                  <Box sx={{ mt: 2 }}>
                    <Link
                      href={location.openstreetmap_url}
                      target="_blank"
                      rel="noopener noreferrer"
                      style={{ textDecoration: "none" }}
                    >
                      <Button
                        variant="outlined"
                        size="small"
                        startIcon={<MapIcon />}
                        sx={{ color: "text.primary", borderColor: "divider" }}
                      >
                        {t("locations.viewMap")}
                      </Button>
                    </Link>
                  </Box>
                )}
              </Paper>
            </Grid>
          ))}
        </Grid>
      ) : (
        <Paper sx={{ p: 4, textAlign: "center" }}>
          <Typography variant="body1" sx={{ color: "text.secondary" }}>
            {t("locations.noLocations")}
          </Typography>
        </Paper>
      )}
    </Container>
  );
}
