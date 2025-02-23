import { useState } from "react";
import { config } from "@/config";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useTranslation } from "@/hooks/useTranslation";
import type {
  HubUserShort,
  MyColleagueApprovals,
} from "@psankar/vetchi-typespec";

interface UseColleagueApprovalsResult {
  approvals: MyColleagueApprovals | null;
  isLoading: boolean;
  error: Error | null;
  fetchApprovals: (paginationKey?: string, limit?: number) => Promise<void>;
}

export function useColleagueApprovals(): UseColleagueApprovalsResult {
  const router = useRouter();
  const { t } = useTranslation();
  const [approvals, setApprovals] = useState<MyColleagueApprovals | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const fetchApprovals = async (paginationKey?: string, limit?: number) => {
    try {
      setIsLoading(true);
      setError(null);

      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/my-colleague-approvals`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            pagination_key: paginationKey,
            limit,
          }),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/login");
        return;
      }

      if (!response.ok) {
        throw new Error(t("approvals.error.fetchFailed"));
      }

      const data = await response.json();
      setApprovals(data);
    } catch (err) {
      setError(
        err instanceof Error ? err : new Error(t("common.error.serverError"))
      );
      throw err;
    } finally {
      setIsLoading(false);
    }
  };

  return {
    approvals,
    isLoading,
    error,
    fetchApprovals,
  };
}
