import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import { Box, Button, Chip, Paper, TextField, Typography } from "@mui/material";
import Cookies from "js-cookie";
import { KeyboardEvent, useState } from "react";

interface ComposeSectionProps {
  onPostCreated: () => void;
  onError: (error: string) => void;
  onSuccess: (message: string) => void;
}

export default function ComposeSection({
  onPostCreated,
  onError,
  onSuccess,
}: ComposeSectionProps) {
  const { t } = useTranslation();
  const [content, setContent] = useState("");
  const [tagIds, setTagIds] = useState<string[]>([]);
  const [newTags, setNewTags] = useState<string[]>([]);
  const [loading, setLoading] = useState(false);
  const [tagInput, setTagInput] = useState("");

  const handleContentChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setContent(event.target.value);
  };

  const handleSubmit = async () => {
    if (!content.trim()) {
      onError(t("posts.contentRequired"));
      return;
    }

    if (content.length > 4096) {
      onError(t("posts.contentTooLong"));
      return;
    }

    setLoading(true);

    try {
      const token = Cookies.get("session_token");
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/add-post`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            content: content.trim(),
            tag_ids: tagIds,
            new_tags: newTags,
          }),
        }
      );

      if (!response.ok) {
        throw new Error(t("posts.createError"));
      }

      setContent("");
      setTagIds([]);
      setNewTags([]);
      onPostCreated();
      onSuccess(t("posts.createSuccess"));
    } catch (err) {
      console.debug("create post error", err);
      onError(t("posts.createError"));
    } finally {
      setLoading(false);
    }
  };

  const handleTagDelete = (tagToDelete: string) => {
    setNewTags((prev) => prev.filter((tag) => tag !== tagToDelete));
  };

  const handleTagInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setTagInput(event.target.value);
  };

  const handleTagKeyDown = (event: KeyboardEvent<HTMLDivElement>) => {
    if (event.key === "Enter" && tagInput.trim()) {
      if (tagIds.length + newTags.length >= 3) {
        onError(t("posts.maxTagsReached"));
        return;
      }

      if (!newTags.includes(tagInput.trim())) {
        setNewTags([...newTags, tagInput.trim()]);
        setTagInput("");
        event.preventDefault();
      }
    }
  };

  return (
    <Paper sx={{ p: 3, mb: 3 }}>
      <Typography variant="h6" gutterBottom>
        {t("posts.compose")}
      </Typography>
      <TextField
        fullWidth
        multiline
        rows={4}
        value={content}
        onChange={handleContentChange}
        placeholder={t("posts.content")}
        variant="outlined"
        sx={{ mb: 2 }}
      />

      <Box sx={{ mb: 2 }}>
        <Typography variant="subtitle2" gutterBottom>
          {t("posts.tags")}:
        </Typography>
        <Box sx={{ display: "flex", flexWrap: "wrap", gap: 1, mb: 1 }}>
          {newTags.map((tag) => (
            <Chip
              key={tag}
              label={tag}
              onDelete={() => handleTagDelete(tag)}
              color="primary"
              variant="outlined"
              size="small"
            />
          ))}
        </Box>
        <TextField
          fullWidth
          placeholder={t("posts.addTags")}
          value={tagInput}
          onChange={handleTagInputChange}
          onKeyDown={handleTagKeyDown}
          helperText={`${t("posts.maxTagsReached")} (${
            tagIds.length + newTags.length
          }/3)`}
          disabled={tagIds.length + newTags.length >= 3}
          size="small"
        />
      </Box>

      <Box sx={{ display: "flex", justifyContent: "flex-end" }}>
        <Button
          variant="contained"
          color="primary"
          onClick={handleSubmit}
          disabled={loading || !content.trim()}
        >
          {t("posts.post")}
        </Button>
      </Box>
    </Paper>
  );
}
