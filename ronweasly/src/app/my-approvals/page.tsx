"use client";

import { useEffect, useRef, useState } from "react";
import Box from "@mui/material/Box";
import Paper from "@mui/material/Paper";
import Typography from "@mui/material/Typography";
import CircularProgress from "@mui/material/CircularProgress";
import Alert from "@mui/material/Alert";
import Stack from "@mui/material/Stack";
import Button from "@mui/material/Button";
import Link from "next/link";
import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import { useTranslation } from "@/hooks/useTranslation";
import { useColleagueApprovals } from "@/hooks/useColleagueApprovals";
import { useColleagues } from "@/hooks/useColleagues";
import type { HubUserShort } from "@psankar/vetchi-typespec";
import { config } from "@/config";
import ProfilePicture from "@/components/ProfilePicture";

const PAGE_SIZE = 20;

export default function MyApprovalsPage() {
  const { t } = useTranslation();
  const [paginationKey, setPaginationKey] = useState<string | undefined>();
  const [hasMore, setHasMore] = useState(true);
  const [approvalsList, setApprovalsList] = useState<HubUserShort[]>([]);
  const [isInitialLoad, setIsInitialLoad] = useState(true);
  const loadingRef = useRef<HTMLDivElement>(null);
  const observerRef = useRef<IntersectionObserver | null>(null);

  const { isLoading, error, fetchApprovals } = useColleagueApprovals();
  const { approveColleague, rejectColleague, isApproving, isRejecting } =
    useColleagues();

  const loadMore = async () => {
    if (!hasMore || isLoading) return;

    try {
      const result = await fetchApprovals(paginationKey, PAGE_SIZE);
      if (result?.approvals) {
        setApprovalsList((prev) =>
          isInitialLoad ? result.approvals : [...prev, ...result.approvals]
        );
        setIsInitialLoad(false);

        // Get the last item's handle as the next pagination key
        const lastItem = result.approvals[result.approvals.length - 1];
        setPaginationKey(lastItem?.handle);
        setHasMore(result.approvals.length === PAGE_SIZE);
      } else {
        setHasMore(false);
      }
    } catch (err) {
      // Error handling is done in the hook
      setHasMore(false);
    }
  };

  // Setup intersection observer
  useEffect(() => {
    observerRef.current = new IntersectionObserver(
      (entries) => {
        if (
          entries[0].isIntersecting &&
          hasMore &&
          !isLoading &&
          !isInitialLoad
        ) {
          loadMore();
        }
      },
      { threshold: 1.0 }
    );

    return () => {
      if (observerRef.current) {
        observerRef.current.disconnect();
      }
    };
  }, [hasMore, isLoading, isInitialLoad]);

  // Handle observer connection/disconnection
  useEffect(() => {
    if (loadingRef.current && observerRef.current) {
      observerRef.current.observe(loadingRef.current);
    }

    return () => {
      if (observerRef.current) {
        observerRef.current.disconnect();
      }
    };
  }, [loadingRef.current]);

  // Initial load
  useEffect(() => {
    setIsInitialLoad(true);
    setPaginationKey(undefined);
    setApprovalsList([]);
    setHasMore(true);
    loadMore();

    // Cleanup function
    return () => {
      setApprovalsList([]);
      setPaginationKey(undefined);
      setHasMore(true);
      setIsInitialLoad(true);
    };
  }, []);

  const handleApprove = async (handle: string) => {
    try {
      await approveColleague(handle);
      // Remove the approved user from the list
      setApprovalsList((prev) => prev.filter((user) => user.handle !== handle));

      // If the list is now empty or getting too short, fetch more
      if (approvalsList.length <= 5) {
        const result = await fetchApprovals(paginationKey, PAGE_SIZE);
        if (result?.approvals) {
          setApprovalsList((prev) => [...prev, ...result.approvals]);
          const lastItem = result.approvals[result.approvals.length - 1];
          setPaginationKey(lastItem?.handle);
          setHasMore(result.approvals.length === PAGE_SIZE);
        }
      }
    } catch (err) {
      // Error handling is done in the hook
    }
  };

  const handleReject = async (handle: string) => {
    try {
      await rejectColleague(handle);
      // Remove the rejected user from the list
      setApprovalsList((prev) => prev.filter((user) => user.handle !== handle));

      // If the list is now empty or getting too short, fetch more
      if (approvalsList.length <= 5) {
        const result = await fetchApprovals(paginationKey, PAGE_SIZE);
        if (result?.approvals) {
          setApprovalsList((prev) => [...prev, ...result.approvals]);
          const lastItem = result.approvals[result.approvals.length - 1];
          setPaginationKey(lastItem?.handle);
          setHasMore(result.approvals.length === PAGE_SIZE);
        }
      }
    } catch (err) {
      // Error handling is done in the hook
    }
  };

  return (
    <AuthenticatedLayout>
      <Box sx={{ maxWidth: 800, mx: "auto", mt: 4, px: 2 }}>
        <Typography variant="h4" gutterBottom>
          {t("approvals.title")}
        </Typography>

        <Paper sx={{ p: 3, mb: 4 }}>
          <Typography variant="h6" gutterBottom>
            {t("approvals.colleagueApprovals")}
          </Typography>

          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error.message}
            </Alert>
          )}

          {approvalsList.length === 0 && !isLoading ? (
            <Typography
              color="text.secondary"
              sx={{ textAlign: "center", p: 3 }}
            >
              {t("approvals.noApprovals")}
            </Typography>
          ) : (
            <Stack spacing={2}>
              {approvalsList.map((user: HubUserShort) => (
                <Paper key={user.handle} variant="outlined" sx={{ p: 2 }}>
                  <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
                    <ProfilePicture
                      imageUrl={`${config.API_SERVER_PREFIX}/hub/profile-picture/${user.handle}`}
                      size={40}
                    />
                    <Box sx={{ flex: 1 }}>
                      <Link
                        href={`/u/${user.handle}`}
                        style={{ textDecoration: "none", color: "inherit" }}
                      >
                        <Typography variant="subtitle1" component="div">
                          {user.name}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          @{user.handle}
                        </Typography>
                      </Link>
                      <Typography variant="body2" sx={{ mt: 1 }}>
                        {user.short_bio}
                      </Typography>
                    </Box>
                    <Stack direction="row" spacing={1}>
                      <Button
                        variant="contained"
                        color="primary"
                        onClick={() => handleApprove(user.handle)}
                        disabled={isApproving || isRejecting}
                      >
                        {t("common.approve")}
                      </Button>
                      <Button
                        variant="outlined"
                        color="error"
                        onClick={() => handleReject(user.handle)}
                        disabled={isApproving || isRejecting}
                      >
                        {t("common.reject")}
                      </Button>
                    </Stack>
                  </Box>
                </Paper>
              ))}

              {/* Loading indicator */}
              <Box
                ref={loadingRef}
                sx={{ height: 20, display: "flex", justifyContent: "center" }}
              >
                {isLoading && <CircularProgress size={20} />}
              </Box>
            </Stack>
          )}
        </Paper>
      </Box>
    </AuthenticatedLayout>
  );
}
