import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import type { Bio, GetBioRequest, UpdateBioRequest } from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

export function useProfile(handle: string) {
  const router = useRouter();
  const { t } = useTranslation();
  const [bio, setBio] = useState<Bio | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    if (!handle) {
      setIsLoading(false);
      return;
    }
    fetchBio();
  }, [handle]);

  const fetchBio = async () => {
    if (!handle) {
      return;
    }

    try {
      setIsLoading(true);
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        return;
      }

      const response = await fetch(`${config.API_SERVER_PREFIX}/hub/get-bio`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ handle } satisfies GetBioRequest),
      });

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/login");
        return;
      }

      if (!response.ok) {
        throw new Error(t("profile.bio.error.fetchFailed"));
      }

      const data = await response.json();
      setBio(data);
      setError(null);
    } catch (err) {
      setError(
        err instanceof Error ? err : new Error(t("common.error.serverError"))
      );
    } finally {
      setIsLoading(false);
    }
  };

  const updateBio = async (updatedBio: Bio) => {
    if (!handle) {
      return;
    }

    try {
      setIsSaving(true);
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/update-bio`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            full_name: updatedBio.full_name,
            short_bio: updatedBio.short_bio,
            long_bio: updatedBio.long_bio,
          } satisfies UpdateBioRequest),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/login");
        return;
      }

      if (!response.ok) {
        throw new Error(t("profile.bio.error.updateFailed"));
      }

      setBio(updatedBio);
      setError(null);
    } catch (err) {
      setError(
        err instanceof Error ? err : new Error(t("common.error.serverError"))
      );
      throw err;
    } finally {
      setIsSaving(false);
    }
  };

  const uploadProfilePicture = async (file: File) => {
    if (!handle) {
      return;
    }

    try {
      setIsSaving(true);
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        return;
      }

      const formData = new FormData();
      formData.append("image", file);

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/upload-profile-picture`,
        {
          method: "POST",
          headers: {
            Authorization: `Bearer ${token}`,
          },
          body: formData,
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/login");
        return;
      }

      if (!response.ok) {
        throw new Error(t("profile.bio.error.uploadFailed"));
      }

      setError(null);
    } catch (err) {
      setError(
        err instanceof Error ? err : new Error(t("common.error.serverError"))
      );
      throw err;
    } finally {
      setIsSaving(false);
    }
  };

  return {
    bio,
    isLoading,
    error,
    isSaving,
    updateBio,
    uploadProfilePicture,
    refetch: fetchBio,
  };
}
