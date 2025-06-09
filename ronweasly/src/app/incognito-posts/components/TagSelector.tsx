"use client";

import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import {
  CircularProgress,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
} from "@mui/material";
import { VTag, VTagID } from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useEffect, useState } from "react";

interface TagSelectorProps {
  selectedTag: VTagID;
  onTagSelect: (tagId: VTagID) => void;
  onError: (error: string) => void;
}

export default function TagSelector({
  selectedTag,
  onTagSelect,
  onError,
}: TagSelectorProps) {
  const { t } = useTranslation();
  const [tags, setTags] = useState<VTag[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const loadTags = async () => {
      try {
        const token = Cookies.get("session_token");
        if (!token) {
          throw new Error("User not authenticated");
        }

        // Using the existing filter-vtags endpoint
        const response = await fetch(
          `${config.API_SERVER_PREFIX}/hub/filter-vtags`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify({ prefix: "" }), // Empty prefix to get all tags
          }
        );

        if (!response.ok) {
          throw new Error(`Failed to fetch tags: ${response.statusText}`);
        }

        const data: VTag[] = await response.json();
        setTags(data);
      } catch (error) {
        onError(
          error instanceof Error
            ? error.message
            : t("incognitoPosts.errors.tagLoadFailed")
        );
      } finally {
        setIsLoading(false);
      }
    };

    loadTags();
  }, [onError, t]);

  const handleTagChange = (tagId: VTagID) => {
    onTagSelect(tagId);
  };

  return (
    <FormControl fullWidth>
      <InputLabel>{t("incognitoPosts.filters.tagFilter")}</InputLabel>
      <Select
        value={selectedTag}
        label={t("incognitoPosts.filters.tagFilter")}
        onChange={(e) => handleTagChange(e.target.value as VTagID)}
        disabled={isLoading}
        endAdornment={
          isLoading ? <CircularProgress size={20} sx={{ mr: 2 }} /> : null
        }
      >
        <MenuItem value="">
          <em>{t("incognitoPosts.filters.selectTag")}</em>
        </MenuItem>
        {tags.map((tag) => (
          <MenuItem key={tag.id} value={tag.id}>
            {tag.name}
          </MenuItem>
        ))}
      </Select>
    </FormControl>
  );
}
