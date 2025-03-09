"use client";

import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import {
  Alert,
  Box,
  Button,
  CircularProgress,
  FormControlLabel,
  Paper,
  Switch,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
} from "@mui/material";
import {
  OpeningState,
  OpeningStates,
} from "@psankar/vetchi-typespec/common/openings";
import Cookies from "js-cookie";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useMemo, useState } from "react";

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
  state: OpeningState;
  created_at: string;
  last_updated_at: string;
}

export default function Openings() {
  const [openings, setOpenings] = useState<Opening[]>([]);
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const [showClosed, setShowClosed] = useState(false);

  const { t } = useTranslation();
  const router = useRouter();

  useEffect(() => {
    const saved = localStorage.getItem("showClosedOpenings");
    if (saved) {
      setShowClosed(JSON.parse(saved));
    }
  }, []);

  // Memoize the states array
  const states = useMemo(() => {
    const baseStates = [
      OpeningStates.DRAFT,
      OpeningStates.ACTIVE,
      OpeningStates.SUSPENDED,
    ];
    return showClosed ? [...baseStates, OpeningStates.CLOSED] : baseStates;
  }, [showClosed]);

  // Handle switch change
  const handleShowClosedChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const newValue = e.target.checked;
      setShowClosed(newValue);
      localStorage.setItem("showClosedOpenings", JSON.stringify(newValue));
    },
    [] // Remove unnecessary t dependency
  );

  // Fetch openings function
  const fetchOpenings = useCallback(
    async (isMounted: boolean) => {
      try {
        const sessionToken = Cookies.get("session_token");
        if (!sessionToken) {
          if (isMounted) {
            setError("auth.unauthorized");
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
            body: JSON.stringify({ state: states }),
          }
        );

        if (!isMounted) return;

        if (response.status === 200) {
          const data = await response.json();
          setOpenings(data);
        } else if (response.status === 401) {
          setError("auth.unauthorized");
        } else {
          setError("common.error");
        }
      } catch {
        if (isMounted) {
          setError(t("openings.fetchError"));
        }
      } finally {
        if (isMounted) {
          setIsLoading(false);
        }
      }
    },
    [states, t] // Add missing t dependency
  );

  useEffect(() => {
    let isMounted = true;
    fetchOpenings(isMounted);
    return () => {
      isMounted = false;
    };
  }, [fetchOpenings]);

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
        <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
          <FormControlLabel
            control={
              <Switch checked={showClosed} onChange={handleShowClosedChange} />
            }
            label={t("openings.showClosed")}
          />
          <Button
            variant="contained"
            color="primary"
            onClick={() => router.push("/openings/create")}
          >
            {t("openings.create")}
          </Button>
        </Box>
      </Box>
      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {t(error)}
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
                <TableCell>ID</TableCell>
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
                    <TableCell>
                      <Link
                        href={`/openings/${opening.id}`}
                        style={{ color: "blue", textDecoration: "underline" }}
                      >
                        {opening.id}
                      </Link>
                    </TableCell>
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
                    <TableCell>
                      {t(`openings.types.${opening.opening_type}`)}
                    </TableCell>
                    <TableCell>
                      {t(`openings.state.${opening.state}`)}
                    </TableCell>
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
