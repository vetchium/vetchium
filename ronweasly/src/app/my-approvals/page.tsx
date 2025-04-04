"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import ProfilePicture from "@/components/ProfilePicture";
import { config } from "@/config";
import { useColleagueApprovals } from "@/hooks/useColleagueApprovals";
import { useColleagues } from "@/hooks/useColleagues";
import { useEndorseApprovals } from "@/hooks/useEndorseApprovals";
import { useTranslation } from "@/hooks/useTranslation";
import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Chip from "@mui/material/Chip";
import CircularProgress from "@mui/material/CircularProgress";
import Divider from "@mui/material/Divider";
import Grid from "@mui/material/Grid";
import Paper from "@mui/material/Paper";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import type { HubUserShort, MyEndorseApproval } from "@vetchium/typespec";
import { format } from "date-fns";
import Link from "next/link";
import { useEffect, useRef, useState } from "react";

const PAGE_SIZE = 20;

export default function MyApprovalsPage() {
  const { t } = useTranslation();

  // Colleague approvals state
  const [paginationKey, setPaginationKey] = useState<string | undefined>();
  const [hasMore, setHasMore] = useState(true);
  const [approvalsList, setApprovalsList] = useState<HubUserShort[]>([]);
  const [isInitialLoad, setIsInitialLoad] = useState(true);
  const loadingRef = useRef<HTMLDivElement>(null);
  const observerRef = useRef<IntersectionObserver | null>(null);

  // Endorsement approvals state
  const [endorsePaginationKey, setEndorsePaginationKey] = useState<
    string | undefined
  >();
  const [hasMoreEndorsements, setHasMoreEndorsements] = useState(true);
  const [endorsementsList, setEndorsementsList] = useState<MyEndorseApproval[]>(
    []
  );
  const [isInitialEndorsementLoad, setIsInitialEndorsementLoad] =
    useState(true);
  const endorsementLoadingRef = useRef<HTMLDivElement>(null);
  const endorsementObserverRef = useRef<IntersectionObserver | null>(null);

  const { isLoading, error, fetchApprovals } = useColleagueApprovals();
  const {
    isLoading: isEndorsementsLoading,
    error: endorsementsError,
    fetchEndorsements,
    approveEndorsement,
    rejectEndorsement,
  } = useEndorseApprovals();
  const { approveColleague, rejectColleague, isApproving, isRejecting } =
    useColleagues();

  // Load more colleague approvals
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

  // Load more endorsement approvals
  const loadMoreEndorsements = async () => {
    if (!hasMoreEndorsements || isEndorsementsLoading) return;

    try {
      const result = await fetchEndorsements(endorsePaginationKey, PAGE_SIZE);
      if (result?.endorsements) {
        setEndorsementsList((prev) =>
          isInitialEndorsementLoad
            ? result.endorsements
            : [...prev, ...result.endorsements]
        );
        setIsInitialEndorsementLoad(false);

        setEndorsePaginationKey(result.pagination_key);
        setHasMoreEndorsements(
          result.endorsements.length === PAGE_SIZE && !!result.pagination_key
        );
      } else {
        setHasMoreEndorsements(false);
      }
    } catch (err) {
      // Error handling is done in the hook
      setHasMoreEndorsements(false);
    }
  };

  // Setup intersection observer for colleague approvals
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

  // Setup intersection observer for endorsement approvals
  useEffect(() => {
    endorsementObserverRef.current = new IntersectionObserver(
      (entries) => {
        if (
          entries[0].isIntersecting &&
          hasMoreEndorsements &&
          !isEndorsementsLoading &&
          !isInitialEndorsementLoad
        ) {
          loadMoreEndorsements();
        }
      },
      { threshold: 1.0 }
    );

    return () => {
      if (endorsementObserverRef.current) {
        endorsementObserverRef.current.disconnect();
      }
    };
  }, [hasMoreEndorsements, isEndorsementsLoading, isInitialEndorsementLoad]);

  // Handle observer connection/disconnection for colleague approvals
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

  // Handle observer connection/disconnection for endorsement approvals
  useEffect(() => {
    if (endorsementLoadingRef.current && endorsementObserverRef.current) {
      endorsementObserverRef.current.observe(endorsementLoadingRef.current);
    }

    return () => {
      if (endorsementObserverRef.current) {
        endorsementObserverRef.current.disconnect();
      }
    };
  }, [endorsementLoadingRef.current]);

  // Initial load for colleague approvals
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

  // Initial load for endorsement approvals
  useEffect(() => {
    setIsInitialEndorsementLoad(true);
    setEndorsePaginationKey(undefined);
    setEndorsementsList([]);
    setHasMoreEndorsements(true);
    loadMoreEndorsements();

    // Cleanup function
    return () => {
      setEndorsementsList([]);
      setEndorsePaginationKey(undefined);
      setHasMoreEndorsements(true);
      setIsInitialEndorsementLoad(true);
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

  const handleApproveEndorsement = async (applicationId: string) => {
    try {
      await approveEndorsement(applicationId);
      // Remove the approved endorsement from the list
      setEndorsementsList((prev) =>
        prev.filter(
          (endorsement) => endorsement.application_id !== applicationId
        )
      );

      // If the list is now empty or getting too short, fetch more
      if (endorsementsList.length <= 5) {
        const result = await fetchEndorsements(endorsePaginationKey, PAGE_SIZE);
        if (result?.endorsements) {
          setEndorsementsList((prev) => [...prev, ...result.endorsements]);
          setEndorsePaginationKey(result.pagination_key);
          setHasMoreEndorsements(
            result.endorsements.length === PAGE_SIZE && !!result.pagination_key
          );
        }
      }
    } catch (err) {
      // Error handling is done in the hook
    }
  };

  const handleRejectEndorsement = async (applicationId: string) => {
    try {
      await rejectEndorsement(applicationId);
      // Remove the rejected endorsement from the list
      setEndorsementsList((prev) =>
        prev.filter(
          (endorsement) => endorsement.application_id !== applicationId
        )
      );

      // If the list is now empty or getting too short, fetch more
      if (endorsementsList.length <= 5) {
        const result = await fetchEndorsements(endorsePaginationKey, PAGE_SIZE);
        if (result?.endorsements) {
          setEndorsementsList((prev) => [...prev, ...result.endorsements]);
          setEndorsePaginationKey(result.pagination_key);
          setHasMoreEndorsements(
            result.endorsements.length === PAGE_SIZE && !!result.pagination_key
          );
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

        {/* Colleague Approvals Section */}
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

        {/* Endorsement Approvals Section */}
        <Paper sx={{ p: 3, mb: 4 }}>
          <Typography variant="h6" gutterBottom>
            {t("approvals.endorsementApprovals")}
          </Typography>

          {endorsementsError && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {endorsementsError.message}
            </Alert>
          )}

          {endorsementsList.length === 0 && !isEndorsementsLoading ? (
            <Typography
              color="text.secondary"
              sx={{ textAlign: "center", p: 3 }}
            >
              {t("approvals.noEndorsements")}
            </Typography>
          ) : (
            <Stack spacing={2}>
              {endorsementsList.map((endorsement: MyEndorseApproval) => (
                <Paper
                  key={endorsement.application_id}
                  variant="outlined"
                  sx={{ p: 2 }}
                >
                  <Grid container spacing={2}>
                    <Grid item xs={12}>
                      <Box
                        sx={{
                          display: "flex",
                          alignItems: "center",
                          gap: 2,
                          mb: 1,
                        }}
                      >
                        <ProfilePicture
                          imageUrl={`${config.API_SERVER_PREFIX}/hub/profile-picture/${endorsement.applicant_handle}`}
                          size={40}
                        />
                        <Box sx={{ flex: 1 }}>
                          <Link
                            href={`/u/${endorsement.applicant_handle}`}
                            style={{ textDecoration: "none", color: "inherit" }}
                          >
                            <Typography variant="subtitle1" component="div">
                              {endorsement.applicant_name}
                            </Typography>
                            <Typography variant="body2" color="text.secondary">
                              @{endorsement.applicant_handle}
                            </Typography>
                          </Link>
                        </Box>
                        <Chip
                          label={endorsement.application_status}
                          color="primary"
                          size="small"
                          variant="outlined"
                        />
                      </Box>
                    </Grid>

                    <Grid item xs={12}>
                      <Divider />
                    </Grid>

                    <Grid item xs={12}>
                      <Typography variant="body2" sx={{ mb: 1 }}>
                        {t("approvals.endorsement.from")}{" "}
                        <strong>{endorsement.applicant_name}</strong>{" "}
                        {t("approvals.endorsement.for")}{" "}
                        <strong>{endorsement.opening_title}</strong>{" "}
                        {t("approvals.endorsement.at")}{" "}
                        <strong>{endorsement.employer_name}</strong>
                      </Typography>
                      <Typography
                        variant="body2"
                        color="text.secondary"
                        sx={{ mb: 1 }}
                      >
                        {t("approvals.endorsement.appliedOn", {
                          date: format(
                            new Date(endorsement.application_created_at),
                            "MMM d, yyyy"
                          ),
                        })}
                      </Typography>
                      <Typography variant="body2" sx={{ mb: 2 }}>
                        {endorsement.applicant_short_bio}
                      </Typography>

                      <Box
                        sx={{
                          display: "flex",
                          justifyContent: "space-between",
                          alignItems: "center",
                        }}
                      >
                        <Link
                          href={endorsement.opening_url}
                          target="_blank"
                          rel="noopener noreferrer"
                          style={{ textDecoration: "none" }}
                        >
                          <Button size="small" variant="text">
                            {t("approvals.endorsement.viewOpening")}
                          </Button>
                        </Link>

                        <Stack direction="row" spacing={1}>
                          <Button
                            variant="contained"
                            color="primary"
                            onClick={() =>
                              handleApproveEndorsement(
                                endorsement.application_id
                              )
                            }
                            size="small"
                          >
                            {t("common.approve")}
                          </Button>
                          <Button
                            variant="outlined"
                            color="error"
                            onClick={() =>
                              handleRejectEndorsement(
                                endorsement.application_id
                              )
                            }
                            size="small"
                          >
                            {t("common.reject")}
                          </Button>
                        </Stack>
                      </Box>
                    </Grid>
                  </Grid>
                </Paper>
              ))}

              {/* Loading indicator */}
              <Box
                ref={endorsementLoadingRef}
                sx={{ height: 20, display: "flex", justifyContent: "center" }}
              >
                {isEndorsementsLoading && <CircularProgress size={20} />}
              </Box>
            </Stack>
          )}
        </Paper>
      </Box>
    </AuthenticatedLayout>
  );
}
