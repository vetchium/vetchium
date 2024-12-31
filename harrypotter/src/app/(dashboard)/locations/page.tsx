"use client";

import {
  Box,
  Button,
  Container,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  IconButton,
  Alert,
  Typography,
} from "@mui/material";
import { Edit as EditIcon, Delete as DeleteIcon } from "@mui/icons-material";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useTranslation } from "@/hooks/useTranslation";
import { config } from "@/config";
import Cookies from "js-cookie";
import {
  Location,
  GetLocationsRequest,
  AddLocationRequest,
  UpdateLocationRequest,
  DefunctLocationRequest,
} from "@psankar/vetchi-typespec";

export default function LocationsPage() {
  const [locations, setLocations] = useState<Location[]>([]);
  const [openDialog, setOpenDialog] = useState(false);
  const [editingLocation, setEditingLocation] = useState<Location | null>(null);
  const [title, setTitle] = useState("");
  const [countryCode, setCountryCode] = useState("");
  const [postalAddress, setPostalAddress] = useState("");
  const [postalCode, setPostalCode] = useState("");
  const [openstreetmapUrl, setOpenstreetmapUrl] = useState("");
  const [cityAka, setCityAka] = useState<string[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [cityAkaInput, setCityAkaInput] = useState("");
  const router = useRouter();
  const { t } = useTranslation();

  const fetchLocations = async () => {
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const request: GetLocationsRequest = {};

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
    fetchLocations();
  }, []);

  const handleAddClick = () => {
    setEditingLocation(null);
    setTitle("");
    setCountryCode("");
    setPostalAddress("");
    setPostalCode("");
    setOpenstreetmapUrl("");
    setCityAka([]);
    setOpenDialog(true);
  };

  const handleEditClick = (location: Location) => {
    setEditingLocation(location);
    setTitle(location.title);
    setCountryCode(location.country_code);
    setPostalAddress(location.postal_address);
    setPostalCode(location.postal_code);
    setOpenstreetmapUrl(location.openstreetmap_url || "");
    setCityAka(location.city_aka || []);
    setOpenDialog(true);
  };

  const handleClose = () => {
    setOpenDialog(false);
    setTitle("");
    setCountryCode("");
    setPostalAddress("");
    setPostalCode("");
    setOpenstreetmapUrl("");
    setCityAka([]);
    setEditingLocation(null);
  };

  const handleSave = async () => {
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const url = editingLocation
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
        throw new Error(
          editingLocation ? t("locations.updateError") : t("locations.addError")
        );
      }

      handleClose();
      fetchLocations();
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    }
  };

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
        throw new Error("Failed to delete location");
      }

      fetchLocations();
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    }
  };

  const handleAddCityAka = () => {
    if (cityAkaInput.trim() && cityAka.length < 3) {
      setCityAka([...cityAka, cityAkaInput.trim()]);
      setCityAkaInput("");
    }
  };

  const handleRemoveCityAka = (index: number) => {
    setCityAka(cityAka.filter((_, i) => i !== index));
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
        <Typography variant="h4" component="h1">
          {t("locations.title")}
        </Typography>
        <Button variant="contained" onClick={handleAddClick}>
          {t("locations.add")}
        </Button>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>{t("locations.locationTitle")}</TableCell>
              <TableCell>{t("locations.countryCode")}</TableCell>
              <TableCell>{t("locations.postalCode")}</TableCell>
              <TableCell align="right">{t("common.actions")}</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {locations.map((location) => (
              <TableRow key={location.title}>
                <TableCell>{location.title}</TableCell>
                <TableCell>{location.country_code}</TableCell>
                <TableCell>{location.postal_code}</TableCell>
                <TableCell align="right">
                  <IconButton onClick={() => handleEditClick(location)}>
                    <EditIcon />
                  </IconButton>
                  <IconButton onClick={() => handleDelete(location)}>
                    <DeleteIcon />
                  </IconButton>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      <Dialog open={openDialog} onClose={handleClose}>
        <DialogTitle>
          {editingLocation ? t("locations.editTitle") : t("locations.addTitle")}
        </DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label={t("locations.locationTitle")}
            type="text"
            fullWidth
            required
            value={title}
            onChange={(e) => setTitle(e.target.value)}
          />
          <TextField
            margin="dense"
            label={t("locations.countryCode")}
            type="text"
            fullWidth
            required
            value={countryCode}
            onChange={(e) => setCountryCode(e.target.value.toUpperCase())}
            inputProps={{ maxLength: 3 }}
            helperText={t("locations.countryCodeHelp")}
          />
          <TextField
            margin="dense"
            label={t("locations.postalAddress")}
            type="text"
            fullWidth
            required
            multiline
            rows={3}
            value={postalAddress}
            onChange={(e) => setPostalAddress(e.target.value)}
          />
          <TextField
            margin="dense"
            label={t("locations.postalCode")}
            type="text"
            fullWidth
            required
            value={postalCode}
            onChange={(e) => setPostalCode(e.target.value)}
          />
          <TextField
            margin="dense"
            label={t("locations.mapUrl")}
            type="url"
            fullWidth
            value={openstreetmapUrl}
            onChange={(e) => setOpenstreetmapUrl(e.target.value)}
          />
          <Box sx={{ mt: 2 }}>
            <Typography variant="subtitle2" gutterBottom>
              {t("locations.cityAka")} ({cityAka.length}/3)
            </Typography>
            <Box sx={{ display: "flex", gap: 1, mb: 1 }}>
              <TextField
                size="small"
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
                  mb: 0.5,
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
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose}>{t("common.cancel")}</Button>
          <Button
            onClick={handleSave}
            disabled={!title || !countryCode || !postalAddress || !postalCode}
          >
            {t("common.save")}
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
}
