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
  CostCenter,
  GetCostCentersRequest,
  AddCostCenterRequest,
  UpdateCostCenterRequest,
  DefunctCostCenterRequest,
} from "@psankar/vetchi-typespec";

export default function CostCentersPage() {
  const [costCenters, setCostCenters] = useState<CostCenter[]>([]);
  const [openDialog, setOpenDialog] = useState(false);
  const [editingCostCenter, setEditingCostCenter] = useState<CostCenter | null>(
    null
  );
  const [name, setName] = useState("");
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();
  const { t } = useTranslation();

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

      if (response.status === 401) {
        router.push("/signin");
        return;
      }

      if (!response.ok) {
        throw new Error(t("costCenters.fetchError"));
      }

      const data = await response.json();
      setCostCenters(data || []);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("costCenters.fetchError")
      );
    }
  };

  useEffect(() => {
    fetchCostCenters();
  }, []);

  const handleAddClick = () => {
    setEditingCostCenter(null);
    setName("");
    setOpenDialog(true);
  };

  const handleEditClick = (costCenter: CostCenter) => {
    setEditingCostCenter(costCenter);
    setName(costCenter.name);
    setOpenDialog(true);
  };

  const handleClose = () => {
    setOpenDialog(false);
    setName("");
    setEditingCostCenter(null);
  };

  const handleSave = async () => {
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const url = editingCostCenter
        ? `${config.API_SERVER_PREFIX}/employer/update-cost-center`
        : `${config.API_SERVER_PREFIX}/employer/add-cost-center`;

      const response = await fetch(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          name,
        }),
      });

      if (response.status === 401) {
        router.push("/signin");
        return;
      }

      if (!response.ok) {
        throw new Error(
          editingCostCenter
            ? "Failed to update cost center"
            : "Failed to add cost center"
        );
      }

      handleClose();
      fetchCostCenters();
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    }
  };

  const handleDelete = async (costCenter: CostCenter) => {
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/defunct-cost-center`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            name: costCenter.name,
          }),
        }
      );

      if (response.status === 401) {
        router.push("/signin");
        return;
      }

      if (!response.ok) {
        throw new Error("Failed to delete cost center");
      }

      fetchCostCenters();
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
          {t("costCenters.title")}
        </Typography>
        <Button variant="contained" onClick={handleAddClick}>
          {t("costCenters.add")}
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
              <TableCell>{t("costCenters.name")}</TableCell>
              <TableCell align="right">{t("common.actions")}</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {costCenters.map((costCenter) => (
              <TableRow key={costCenter.name}>
                <TableCell>{costCenter.name}</TableCell>
                <TableCell align="right">
                  <IconButton onClick={() => handleEditClick(costCenter)}>
                    <EditIcon />
                  </IconButton>
                  <IconButton onClick={() => handleDelete(costCenter)}>
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
          {editingCostCenter
            ? t("costCenters.editTitle")
            : t("costCenters.addTitle")}
        </DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label={t("costCenters.name")}
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
