"use client";

import { useParams, useRouter } from "next/navigation";
import { Box, Typography, Button } from "@mui/material";
import { useTranslation } from "@/hooks/useTranslation";

export default function InterviewDetailPage() {
  const params = useParams();
  const interviewId = params.id as string;
  const { t } = useTranslation();
  const router = useRouter();

  return (
    <Box sx={{ p: 3 }}>
      <Box sx={{ display: "flex", justifyContent: "space-between", mb: 3 }}>
        <Typography variant="h4">{t("interviews.manageInterview")}</Typography>
        <Button variant="outlined" onClick={() => router.back()}>
          {t("common.back")}
        </Button>
      </Box>

      <Typography>
        {t("interviews.placeholder")} (ID: {interviewId})
      </Typography>
    </Box>
  );
}
