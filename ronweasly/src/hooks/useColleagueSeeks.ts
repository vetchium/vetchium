import { useState } from "react";
import { useTranslation } from "@/hooks/useTranslation";
import { useRouter } from "next/navigation";
import type { HubUserShort, MyColleagueSeeks } from "@psankar/vetchi-typespec";
import { config } from "@/config";
import Cookies from "js-cookie";

interface UseColleagueSeeksResult {
  seeks: MyColleagueSeeks | null;
  isLoading: boolean;
  error: Error | null;
  fetchSeeks: () => Promise<void>;
}

export function useColleagueSeeks(): UseColleagueSeeksResult {
  const [seeks, setSeeks] = useState<MyColleagueSeeks | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const { t } = useTranslation();
  const router = useRouter();

  const fetchSeeks = async () => {
    setIsLoading(true);
    setError(null);

    const sessionToken = Cookies.get("session_token");
    if (!sessionToken) {
      router.push("/login");
      setError(new Error(t("common.error.notAuthenticated")));
      setIsLoading(false);
      return;
    }

    try {
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/my-colleague-seeks`,
        {
          method: "POST",
          credentials: "include",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${sessionToken}`,
          },
          body: JSON.stringify({
            limit: 40,
          }),
        }
      );

      if (!response.ok) {
        if (response.status === 401) {
          router.push("/login");
          throw new Error(t("common.error.notAuthenticated"));
        }
        throw new Error(t("requisitions.error.fetchFailed"));
      }

      const data = await response.json();
      setSeeks(data);
    } catch (err) {
      setError(err instanceof Error ? err : new Error(String(err)));
    } finally {
      setIsLoading(false);
    }
  };

  return {
    seeks,
    isLoading,
    error,
    fetchSeeks,
  };
}
