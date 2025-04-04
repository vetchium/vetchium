"use client";

import { useOrgUsers } from "@/hooks/useOrgUsers";
import { useTranslation } from "@/hooks/useTranslation";
import {
  Block as BlockIcon,
  CheckCircle as CheckCircleIcon,
  Warning as WarningIcon,
} from "@mui/icons-material";
import {
  Alert,
  Button,
  Chip,
  CircularProgress,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  FormControlLabel,
  Paper,
  Snackbar,
  Switch,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField,
} from "@mui/material";
import { OrgUser } from "@vetchium/typespec";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

export default function OrgUsersPage() {
  const { t } = useTranslation();
  const [searchQuery, setSearchQuery] = useState("");
  const [disableUserDialogOpen, setDisableUserDialogOpen] = useState(false);
  const [selectedUserEmail, setSelectedUserEmail] = useState<string | null>(
    null
  );
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [includeDisabled, setIncludeDisabled] = useState(false);
  const router = useRouter();

  const { users, isLoading, error, disableUser, enableUser, fetchUsers } =
    useOrgUsers();

  // Load saved preference on mount
  useEffect(() => {
    const savedValue = localStorage.getItem("includeDisabledUsers");
    if (savedValue) {
      setIncludeDisabled(savedValue === "true");
    }
  }, []);

  // Fetch users when includeDisabled changes
  useEffect(() => {
    fetchUsers(includeDisabled);
  }, [includeDisabled, fetchUsers]); // Added fetchUsers to dependencies

  const filteredUsers = users?.filter(
    (user: OrgUser) =>
      user.email.toLowerCase().includes(searchQuery.toLowerCase()) ||
      user.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const handleDisableUser = async () => {
    if (!selectedUserEmail) return;
    try {
      await disableUser(selectedUserEmail);
      setDisableUserDialogOpen(false);
      setSelectedUserEmail(null);
      // Refetch users with current filter state
      fetchUsers(includeDisabled);
    } catch (error) {
      console.error("Failed to disable user:", error);
    }
  };

  const handleEnableUser = async (email: string) => {
    try {
      await enableUser(email);
      // Refetch users with current filter state
      fetchUsers(includeDisabled);
    } catch (error) {
      console.error("Failed to enable user:", error);
    }
  };

  if (isLoading) {
    return (
      <div
        style={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          height: "100%",
        }}
      >
        <CircularProgress />
      </div>
    );
  }

  if (error) {
    return (
      <div style={{ padding: "1rem" }}>
        <Alert severity="error">{t("orgUsers.fetchError")}</Alert>
      </div>
    );
  }

  return (
    <div style={{ padding: "1rem" }}>
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          marginBottom: "1rem",
        }}
      >
        <h1 style={{ fontSize: "1.5rem", fontWeight: "bold" }}>
          {t("orgUsers.title")}
        </h1>
        <Button
          variant="contained"
          onClick={() => router.push("/org-users/create")}
        >
          {t("orgUsers.add")}
        </Button>
      </div>

      <Paper style={{ padding: "1rem" }}>
        <div
          style={{
            display: "flex",
            gap: "1rem",
            alignItems: "center",
            marginBottom: "1rem",
          }}
        >
          <TextField
            placeholder={t("orgUsers.searchPlaceholder")}
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            style={{ maxWidth: "20rem" }}
            fullWidth
          />
          <FormControlLabel
            control={
              <Switch
                checked={includeDisabled}
                onChange={(e) => {
                  setIncludeDisabled(e.target.checked);
                  localStorage.setItem(
                    "includeDisabledUsers",
                    e.target.checked.toString()
                  );
                }}
              />
            }
            label={t("orgUsers.includeDisabled")}
          />
        </div>

        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>{t("orgUsers.email")}</TableCell>
                <TableCell>{t("orgUsers.name")}</TableCell>
                <TableCell>{t("orgUsers.rolesList")}</TableCell>
                <TableCell>{t("orgUsers.state")}</TableCell>
                <TableCell>{t("common.actions")}</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {filteredUsers?.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={5} align="center">
                    {t("orgUsers.noUsers")}
                  </TableCell>
                </TableRow>
              ) : (
                filteredUsers?.map((user: OrgUser) => (
                  <TableRow key={user.email}>
                    <TableCell>{user.email}</TableCell>
                    <TableCell>{user.name}</TableCell>
                    <TableCell>
                      <div
                        style={{
                          display: "flex",
                          flexWrap: "wrap",
                          gap: "0.25rem",
                        }}
                      >
                        {user.roles.map((role: string) => (
                          <Chip
                            key={role}
                            label={t(`orgUsers.roles.${role}`)}
                            size="small"
                          />
                        ))}
                      </div>
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={t(`orgUsers.states.${user.state}`)}
                        color={
                          user.state === "DISABLED_ORG_USER"
                            ? "error"
                            : "default"
                        }
                        size="small"
                      />
                    </TableCell>
                    <TableCell>
                      {user.state === "DISABLED_ORG_USER" ? (
                        <Button
                          variant="outlined"
                          size="small"
                          onClick={() => handleEnableUser(user.email)}
                          startIcon={<CheckCircleIcon />}
                          color="success"
                        >
                          {t("orgUsers.enable")}
                        </Button>
                      ) : (
                        <Button
                          variant="outlined"
                          size="small"
                          onClick={() => {
                            setSelectedUserEmail(user.email);
                            setDisableUserDialogOpen(true);
                          }}
                          startIcon={<BlockIcon />}
                          color="error"
                        >
                          {t("orgUsers.disable")}
                        </Button>
                      )}
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>

      <Dialog
        open={disableUserDialogOpen}
        onClose={() => setDisableUserDialogOpen(false)}
      >
        <DialogTitle sx={{ display: "flex", alignItems: "center", gap: 1 }}>
          <WarningIcon color="error" />
          {t("orgUsers.confirmDisable.modalTitle")}
        </DialogTitle>
        <DialogContent>
          <DialogContentText>
            {t("orgUsers.confirmDisable.message")}
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDisableUserDialogOpen(false)}>
            {t("orgUsers.confirmDisable.cancelButton")}
          </Button>
          <Button
            onClick={handleDisableUser}
            color="error"
            variant="contained"
            startIcon={<BlockIcon />}
          >
            {t("orgUsers.confirmDisable.confirmButton")}
          </Button>
        </DialogActions>
      </Dialog>

      <Snackbar
        open={!!successMessage}
        autoHideDuration={6000}
        onClose={() => setSuccessMessage(null)}
        anchorOrigin={{ vertical: "bottom", horizontal: "center" }}
      >
        <Alert onClose={() => setSuccessMessage(null)} severity="success">
          {successMessage}
        </Alert>
      </Snackbar>
    </div>
  );
}
