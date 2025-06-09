"use client";

import { useTranslation } from "@/hooks/useTranslation";
import { Alert, Box } from "@mui/material";
import { VTagID } from "@vetchium/typespec";
import { useRouter } from "next/navigation";
import TagSelector from "./TagSelector";

interface BrowseTabProps {
  onError: (error: string) => void;
}

export default function BrowseTab({ onError }: BrowseTabProps) {
  const { t } = useTranslation();
  const router = useRouter();

  const handleTagSelect = (tagId: VTagID) => {
    if (tagId) {
      // Navigate to the tag-specific page
      router.push(`/incognito-posts/tag/${tagId}`);
    }
  };

  return (
    <Box>
      {/* Tag Selector */}
      <Box sx={{ mb: 3, display: "flex", gap: 2, flexWrap: "wrap" }}>
        <Box sx={{ minWidth: 200 }}>
          <TagSelector
            selectedTag=""
            onTagSelect={handleTagSelect}
            onError={onError}
          />
        </Box>
      </Box>

      {/* Content */}
      <Alert severity="info">{t("incognitoPosts.feed.selectTagFirst")}</Alert>
    </Box>
  );
}
