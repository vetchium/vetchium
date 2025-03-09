"use client";

import {
  Box,
  Button,
  Container,
  TextField,
  Typography,
  Alert,
  Paper,
  IconButton,
} from "@mui/material";
import { Delete as DeleteIcon } from "@mui/icons-material";
import { useState, useEffect, useCallback } from "react";
import { useRouter, useParams } from "next/navigation";
import { useTranslation } from "@/hooks/useTranslation";
import { config } from "@/config";
import Cookies from "js-cookie";
import {
  GetLocationRequest,
  AddLocationRequest,
  UpdateLocationRequest,
} from "@psankar/vetchi-typespec";

export default function LocationActionPage() {
  const [title, setTitle] = useState("");
  const [countryCode, setCountryCode] = useState("");
  const [postalAddress, setPostalAddress] = useState("");
  const [postalCode, setPostalCode] = useState("");
  const [openstreetmapUrl, setOpenstreetmapUrl] = useState("");
  const [cityAka, setCityAka] = useState<string[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [cityAkaInput, setCityAkaInput] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const router = useRouter();
  const params = useParams();
  const { t } = useTranslation();
  const isEdit = params.action === "edit";

  const fetchLocation = useCallback(async (locationTitle: string) => {
    try {
      setIsLoading(true);
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const request: GetLocationRequest = {
        title: locationTitle,
      };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/get-location`,
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

      const location = await response.json();
      setTitle(location.title);
      setCountryCode(location.country_code);
      setPostalAddress(location.postal_address);
      setPostalCode(location.postal_code);
      setOpenstreetmapUrl(location.openstreetmap_url || "");
      setCityAka(location.city_aka || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : t("locations.fetchError"));
    } finally {
      setIsLoading(false);
    }
  }, [router, t]);

  useEffect(() => {
    if (isEdit) {
      const locationId = window.location.search.split("id=")[1];
      if (locationId) {
        fetchLocation(locationId);
      } else {
        router.push("/locations");
      }
    }
  }, [isEdit, fetchLocation, router]);

  const handleAddCityAka = () => {
    if (cityAkaInput.trim() && cityAka.length < 3) {
      setCityAka([...cityAka, cityAkaInput.trim()]);
      setCityAkaInput("");
    }
  };

  const handleRemoveCityAka = (index: number) => {
    setCityAka(cityAka.filter((_, i) => i !== index));
  };

  const handleSave = async () => {
    try {
      setError(null);
      setIsLoading(true);
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const url = isEdit
        ? `${config.API_SERVER_PREFIX}/employer/update-location`
        : `${config.API_SERVER_PREFIX}/employer/add-location`;

      const requestBody: AddLocationRequest | UpdateLocationRequest = {
        title,
        country_code: countryCode,
        postal_address: postalAddress,
        postal_code: postalCode,
        openstreetmap_url: openstreetmapUrl || undefined,
        city_aka: cityAka.length > 0 ? cityAka : undefined,
      };

      const response = await fetch(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(requestBody),
      });

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/signin");
        return;
      }

      if (!response.ok) {
        const data = await response.json();
        throw new Error(
          data.message ||
            (isEdit ? t("locations.updateError") : t("locations.addError"))
        );
      }

      router.push("/locations");
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    } finally {
      setIsLoading(false);
    }
  };

  if (isLoading && isEdit) {
    return (
      <Container maxWidth="md">
        <Typography>{t("common.loading")}</Typography>
      </Container>
    );
  }

  return (
    <Container maxWidth="md">
      <Paper sx={{ p: 4 }}>
        <Box sx={{ mb: 4 }}>
          <Typography variant="h4" component="h1" gutterBottom>
            {isEdit ? t("locations.editTitle") : t("locations.addTitle")}
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
            label={t("locations.locationTitle")}
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            disabled={isEdit}
          />

          <TextField
            margin="normal"
            required
            fullWidth
            label={t("locations.countryCode")}
            value={countryCode}
            onChange={(e) => setCountryCode(e.target.value.toUpperCase())}
            inputProps={{ maxLength: 3 }}
            helperText={t("locations.countryCodeHelp")}
          />

          <TextField
            margin="normal"
            required
            fullWidth
            multiline
            rows={4}
            label={t("locations.postalAddress")}
            value={postalAddress}
            onChange={(e) => setPostalAddress(e.target.value)}
          />

          <TextField
            margin="normal"
            required
            fullWidth
            label={t("locations.postalCode")}
            value={postalCode}
            onChange={(e) => setPostalCode(e.target.value)}
          />

          <TextField
            margin="normal"
            fullWidth
            label={t("locations.mapUrl")}
            type="url"
            value={openstreetmapUrl}
            onChange={(e) => setOpenstreetmapUrl(e.target.value)}
          />

          <Box sx={{ mt: 3 }}>
            <Typography variant="subtitle1" gutterBottom>
              {t("locations.cityAka")} ({cityAka.length}/3)
            </Typography>
            <Box sx={{ display: "flex", gap: 1, mb: 2 }}>
              <TextField
                size="small"
                fullWidth
                value={cityAkaInput}
                onChange={(e) => setCityAkaInput(e.target.value)}
                disabled={cityAka.length >= 3}
                placeholder={t("locations.cityAkaPlaceholder")}
              />
              <Button
                variant="outlined"
                onClick={handleAddCityAka}
                disabled={!cityAkaInput.trim() || cityAka.length >= 3}
              >
                {t("common.add")}
              </Button>
            </Box>
            {cityAka.map((city, index) => (
              <Box
                key={index}
                sx={{
                  display: "flex",
                  alignItems: "center",
                  gap: 1,
                  mb: 1,
                  p: 1,
                  bgcolor: "grey.50",
                  borderRadius: 1,
                }}
              >
                <Typography>{city}</Typography>
                <IconButton
                  size="small"
                  onClick={() => handleRemoveCityAka(index)}
                >
                  <DeleteIcon fontSize="small" />
                </IconButton>
              </Box>
            ))}
          </Box>

          <Box sx={{ mt: 4, display: "flex", gap: 2 }}>
            <Button
              variant="outlined"
              onClick={() => router.push("/locations")}
            >
              {t("common.cancel")}
            </Button>
            <Button
              variant="contained"
              onClick={handleSave}
              disabled={
                isLoading ||
                !title ||
                !countryCode ||
                !postalAddress ||
                !postalCode
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
