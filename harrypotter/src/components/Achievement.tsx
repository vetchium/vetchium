import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import LinkIcon from "@mui/icons-material/Link";
import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import CircularProgress from "@mui/material/CircularProgress";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";
import Link from "@mui/material/Link";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import { Achievement, AchievementType, Handle } from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";

interface AchievementSectionProps {
  userHandle: Handle;
  achievementType: AchievementType;
}

export function AchievementSection({
  userHandle,
  achievementType,
}: AchievementSectionProps) {
  const router = useRouter();
  const { t } = useTranslation();
  const [achievements, setAchievements] = useState<Achievement[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [openUrlWarningDialog, setOpenUrlWarningDialog] = useState(false);
  const [selectedUrl, setSelectedUrl] = useState("");

  // Translation keys based on achievement type
  const transSection =
    achievementType === AchievementType.PATENT
      ? "achievements.patents"
      : achievementType === AchievementType.PUBLICATION
      ? "achievements.publications"
      : "achievements.certifications";

  const fetchAchievements = useCallback(async () => {
    setIsLoading(true);
    setError(null);

    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const request = {
        handle: userHandle,
        type: achievementType,
      };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/list-hub-user-achievements`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(request),
          credentials: "include",
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/signin");
        return;
      }

      if (!response.ok) {
        throw new Error(t(`${transSection}.error.fetchFailed`));
      }

      const data = await response.json();
      setAchievements(data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setIsLoading(false);
    }
  }, [achievementType, userHandle, t, transSection, router]);

  useEffect(() => {
    fetchAchievements();
  }, [fetchAchievements]);

  const formatDate = (dateString?: Date) => {
    if (!dateString) return "";
    const date = new Date(dateString);
    return date.toLocaleDateString();
  };

  const handleUrlClick = (url: string, event: React.MouseEvent) => {
    event.preventDefault();
    setSelectedUrl(url);
    setOpenUrlWarningDialog(true);
  };

  const handleExternalNavigation = () => {
    window.open(selectedUrl, "_blank", "noopener,noreferrer");
    setOpenUrlWarningDialog(false);
  };

  const handleCloseUrlWarningDialog = () => {
    setOpenUrlWarningDialog(false);
    setSelectedUrl("");
  };

  if (isLoading) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", my: 2 }}>
        <CircularProgress size={24} />
      </Box>
    );
  }

  if (error) {
    return (
      <Alert severity="error" sx={{ mb: 2 }}>
        {error}
      </Alert>
    );
  }

  if (achievements.length === 0) {
    return (
      <Typography color="text.secondary">
        {t(`${transSection}.noEntries`)}
      </Typography>
    );
  }

  return (
    <Box>
      <Stack spacing={2}>
        {achievements.map((achievement) => (
          <Card key={achievement.id} variant="outlined">
            <CardContent>
              <Box>
                <Typography variant="h6" component="div">
                  {achievement.title}
                </Typography>

                {achievement.description && (
                  <Typography
                    variant="body2"
                    sx={{
                      mt: 1,
                      whiteSpace: "pre-wrap",
                    }}
                  >
                    {achievement.description}
                  </Typography>
                )}

                {achievement.at && (
                  <Typography
                    variant="body2"
                    color="text.secondary"
                    sx={{ mt: 1 }}
                  >
                    {formatDate(achievement.at)}
                  </Typography>
                )}

                {achievement.url && (
                  <Link
                    href={achievement.url}
                    target="_blank"
                    rel="noopener"
                    onClick={(e) => handleUrlClick(achievement.url!, e)}
                    sx={{
                      display: "inline-flex",
                      alignItems: "center",
                      mt: 1,
                    }}
                  >
                    <LinkIcon fontSize="small" sx={{ mr: 0.5 }} />
                    {t(`${transSection}.url`)}
                  </Link>
                )}
              </Box>
            </CardContent>
          </Card>
        ))}
      </Stack>

      {/* URL Warning Dialog */}
      <Dialog open={openUrlWarningDialog} onClose={handleCloseUrlWarningDialog}>
        <DialogTitle>{t("common.warning")}</DialogTitle>
        <DialogContent>
          <DialogContentText>
            {t("common.external_url_warning")}
            <Typography component="div" sx={{ mt: 2, wordBreak: "break-all" }}>
              {selectedUrl}
            </Typography>
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseUrlWarningDialog}>
            {t("common.cancel")}
          </Button>
          <Button
            onClick={handleExternalNavigation}
            variant="contained"
            color="primary"
          >
            {t("common.proceed")}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
