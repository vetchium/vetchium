"use client";

import { useEffect, useState } from "react";
import {
  Box,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
  Alert,
  CircularProgress,
  Button,
} from "@mui/material";
import { useTranslation } from "@/hooks/useTranslation";
import { config } from "@/config";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";

interface Opening {
  id: string;
  title: string;
  positions: number;
  filled_positions: number;
  recruiter: {
    name: string;
    email: string;
  };
  hiring_manager: {
    name: string;
    email: string;
  };
  cost_center_name: string;
  opening_type: string;
  state: string;
  created_at: string;
  last_updated_at: string;
}

export default function Openings() {
  const [openings, setOpenings] = useState<Opening[]>([]);
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const { t } = useTranslation();
  const router = useRouter();

  useEffect(() => {
    let isMounted = true;

    const fetchOpenings = async () => {
      try {
        const sessionToken = Cookies.get("session_token");
        if (!sessionToken) {
          if (isMounted) {
            setError(t("auth.unauthorized"));
            setIsLoading(false);
          }
          return;
        }

        const response = await fetch(
          `${config.API_SERVER_PREFIX}/employer/filter-openings`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${sessionToken}`,
            },
            body: JSON.stringify({}),
          }
        );

        if (!isMounted) return;

        if (response.status === 200) {
          const data = await response.json();
          setOpenings(data);
        } else if (response.status === 401) {
          setError(t("auth.unauthorized"));
        } else {
          setError(t("common.error"));
        }
      } catch (err) {
        if (isMounted) {
          setError(t("common.error"));
        }
      } finally {
        if (isMounted) {
          setIsLoading(false);
        }
      }
    };

    fetchOpenings();

    return () => {
      isMounted = false;
    };
    // We intentionally omit 't' from dependencies as it would cause unnecessary refetches
    // The translations are stable enough that this is safe
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <Box sx={{ width: "100%" }}>
      <Box
        sx={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          mb: 3,
        }}
      >
        <Typography variant="h4" gutterBottom>
          {t("openings.title")}
        </Typography>
        <Button
          variant="contained"
          color="primary"
          onClick={() => router.push("/openings/create")}
        >
          {t("openings.create")}
        </Button>
      </Box>
      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}
      {isLoading ? (
        <Box sx={{ display: "flex", justifyContent: "center", my: 4 }}>
          <CircularProgress />
        </Box>
      ) : (
        <TableContainer component={Paper}>
          <Table sx={{ minWidth: 650 }} aria-label="openings table">
            <TableHead>
              <TableRow>
                <TableCell>Title</TableCell>
                <TableCell align="right">Positions</TableCell>
                <TableCell align="right">Filled</TableCell>
                <TableCell>Recruiter</TableCell>
                <TableCell>Hiring Manager</TableCell>
                <TableCell>Cost Center</TableCell>
                <TableCell>Type</TableCell>
                <TableCell>State</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {openings.length > 0 ? (
                openings.map((opening) => (
                  <TableRow
                    key={opening.id}
                    sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
                  >
                    <TableCell component="th" scope="row">
                      {opening.title}
                    </TableCell>
                    <TableCell align="right">{opening.positions}</TableCell>
                    <TableCell align="right">
                      {opening.filled_positions}
                    </TableCell>
                    <TableCell>{opening.recruiter.name}</TableCell>
                    <TableCell>{opening.hiring_manager.name}</TableCell>
                    <TableCell>{opening.cost_center_name}</TableCell>
                    <TableCell>{opening.opening_type}</TableCell>
                    <TableCell>{opening.state}</TableCell>
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell colSpan={8} align="center">
                    {t("openings.noOpenings")}
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </TableContainer>
      )}
    </Box>
  );
}
