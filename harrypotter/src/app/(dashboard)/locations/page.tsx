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
  const [name, setName] = useState("");
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();
  const { t } = useTranslation();

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

      if (response.status === 401) {
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
    setName("");
    setOpenDialog(true);
  };

  const handleEditClick = (location: Location) => {
    setEditingLocation(location);
    setName(location.title);
    setOpenDialog(true);
  };

  const handleClose = () => {
    setOpenDialog(false);
    setName("");
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

      const response = await fetch(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          title: name,
        }),
      });

      if (response.status === 401) {
        router.push("/signin");
        return;
      }

      if (!response.ok) {
        throw new Error(
          editingLocation
            ? "Failed to update location"
            : "Failed to add location"
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

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/defunct-location`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            title: location.title,
          }),
        }
      );

      if (response.status === 401) {
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
              <TableCell align="right">{t("common.actions")}</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {locations.map((location) => (
              <TableRow key={location.title}>
                <TableCell>{location.title}</TableCell>
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
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose}>{t("common.cancel")}</Button>
          <Button onClick={handleSave}>{t("common.save")}</Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
}
