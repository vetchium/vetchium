import { useState } from "react";
import { config } from "@/config";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useTranslation } from "@/hooks/useTranslation";
import type {
  MyEndorseApproval,
  MyEndorseApprovalsResponse,
} from "@psankar/vetchi-typespec";
import { EndorsementState } from "@psankar/vetchi-typespec";

interface UseEndorseApprovalsResult {
  endorsements: MyEndorseApprovalsResponse | null;
  isLoading: boolean;
  error: Error | null;
  fetchEndorsements: (
    pagination_key?: string,
    limit?: number,
    state?: EndorsementState[]
  ) => Promise<MyEndorseApprovalsResponse>;
  approveEndorsement: (applicationId: string) => Promise<void>;
  rejectEndorsement: (applicationId: string) => Promise<void>;
}

export function useEndorseApprovals(): UseEndorseApprovalsResult {
  const router = useRouter();
  const { t } = useTranslation();
  const [endorsements, setEndorsements] =
    useState<MyEndorseApprovalsResponse | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [isProcessing, setIsProcessing] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const fetchEndorsements = async (
    pagination_key?: string,
    limit: number = 20,
    state: EndorsementState[] = [EndorsementState.SoughtEndorsement]
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
        `${config.API_SERVER_PREFIX}/hub/my-endorse-approvals`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            pagination_key,
            limit,
            state,
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
      setEndorsements(data);
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

  const approveEndorsement = async (applicationId: string) => {
    try {
      setIsProcessing(true);
      setError(null);

      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        throw new Error("No session token");
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/approve-endorsement`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            application_id: applicationId,
          }),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/login");
        throw new Error("Unauthorized");
      }

      if (!response.ok) {
        throw new Error(t("approvals.error.endorsementActionFailed"));
      }
    } catch (err) {
      const error =
        err instanceof Error ? err : new Error(t("common.error.serverError"));
      setError(error);
      throw error;
    } finally {
      setIsProcessing(false);
    }
  };

  const rejectEndorsement = async (applicationId: string) => {
    try {
      setIsProcessing(true);
      setError(null);

      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        throw new Error("No session token");
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/reject-endorsement`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            application_id: applicationId,
          }),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/login");
        throw new Error("Unauthorized");
      }

      if (!response.ok) {
        throw new Error(t("approvals.error.endorsementActionFailed"));
      }
    } catch (err) {
      const error =
        err instanceof Error ? err : new Error(t("common.error.serverError"));
      setError(error);
      throw error;
    } finally {
      setIsProcessing(false);
    }
  };

  return {
    endorsements,
    isLoading,
    error,
    fetchEndorsements,
    approveEndorsement,
    rejectEndorsement,
  };
}
