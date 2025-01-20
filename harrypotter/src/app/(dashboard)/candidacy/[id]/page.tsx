"use client";

import { useParams } from "next/navigation";
import { Box, Typography, Button } from "@mui/material";
import { useTranslation } from "@/hooks/useTranslation";
import { useRouter } from "next/navigation";

export default function CandidacyDetailPage() {
  const params = useParams();
  const candidacyId = params.id as string;
  const { t } = useTranslation();
  const router = useRouter();

  return (
    <Box sx={{ p: 3 }}>
      <Box sx={{ display: "flex", justifyContent: "space-between", mb: 3 }}>
        <Typography variant="h4">{t("candidacies.viewCandidacy")}</Typography>
        <Button variant="outlined" onClick={() => router.back()}>
          {t("common.back")}
        </Button>
      </Box>
      <Typography>Candidacy ID: {candidacyId}</Typography>
    </Box>
  );
}
