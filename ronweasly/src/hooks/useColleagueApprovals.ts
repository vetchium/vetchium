import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import type { MyColleagueApprovals } from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useState } from "react";

interface UseColleagueApprovalsResult {
  approvals: MyColleagueApprovals | null;
  isLoading: boolean;
  error: Error | null;
  fetchApprovals: (
    pagination_key?: string,
    limit?: number
  ) => Promise<MyColleagueApprovals>;
}

export function useColleagueApprovals(): UseColleagueApprovalsResult {
  const router = useRouter();
  const { t } = useTranslation();
  const [approvals, setApprovals] = useState<MyColleagueApprovals | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const fetchApprovals = async (
    pagination_key?: string,
    limit: number = 20
  ) => {
    try {
      setIsLoading(true);
      setError(null);

      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        throw new Error("No session token");
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
            pagination_key,
            limit,
          }),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/login");
        throw new Error("Unauthorized");
      }

      if (!response.ok) {
        throw new Error(t("approvals.error.fetchFailed"));
      }

      const data = await response.json();
      setApprovals(data);
      return data;
    } catch (err) {
      const error =
        err instanceof Error ? err : new Error(t("common.error.serverError"));
      setError(error);
      throw error;
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
