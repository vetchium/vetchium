"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import ProfilePicture from "@/components/ProfilePicture";
import { config } from "@/config";
import { useAuth } from "@/hooks/useAuth";
import { useColleagueSeeks } from "@/hooks/useColleagueSeeks";
import { useTranslation } from "@/hooks/useTranslation";
import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import CircularProgress from "@mui/material/CircularProgress";
import Paper from "@mui/material/Paper";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import type { HubUserShort } from "@vetchium/typespec";
import Link from "next/link";
import { useEffect, useRef, useState } from "react";

const PAGE_SIZE = 20;

export default function MyRequisitionsPage() {
  const { t } = useTranslation();
  useAuth(); // Check authentication and redirect if not authenticated
  const [paginationKey, setPaginationKey] = useState<string | undefined>();
  const [hasMore, setHasMore] = useState(true);
  const [seeksList, setSeeksList] = useState<HubUserShort[]>([]);
  const [isInitialLoad, setIsInitialLoad] = useState(true);
  const loadingRef = useRef<HTMLDivElement>(null);
  const observerRef = useRef<IntersectionObserver | null>(null);

  const { isLoading, error, fetchSeeks } = useColleagueSeeks();

  const loadMore = async () => {
    if (!hasMore || isLoading) return;

    try {
      const result = await fetchSeeks(paginationKey, PAGE_SIZE);
      if (result?.seeks) {
        setSeeksList((prev) =>
          isInitialLoad ? result.seeks : [...prev, ...result.seeks]
        );
        setIsInitialLoad(false);

        // Get the last item's handle as the next pagination key
        const lastItem = result.seeks[result.seeks.length - 1];
        setPaginationKey(lastItem?.handle);
        setHasMore(result.seeks.length === PAGE_SIZE);
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
    setSeeksList([]);
    setHasMore(true);
    loadMore();

    // Cleanup function
    return () => {
      setSeeksList([]);
      setPaginationKey(undefined);
      setHasMore(true);
      setIsInitialLoad(true);
    };
  }, []);

  return (
    <AuthenticatedLayout>
      <Box sx={{ maxWidth: 800, mx: "auto", mt: 4, px: 2 }}>
        <Typography variant="h4" gutterBottom>
          {t("requisitions.title")}
        </Typography>

        <Paper sx={{ p: 3, mb: 4 }}>
          <Typography variant="h6" gutterBottom>
            {t("requisitions.colleagueSeeks")}
          </Typography>

          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error.message}
            </Alert>
          )}

          {seeksList.length === 0 && !isLoading ? (
            <Typography
              color="text.secondary"
              sx={{ textAlign: "center", p: 3 }}
            >
              {t("requisitions.noSeeks")}
            </Typography>
          ) : (
            <Stack spacing={2}>
              {seeksList.map((user: HubUserShort) => (
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
