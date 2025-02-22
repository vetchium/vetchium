import { useState } from "react";
import { config } from "@/config";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useTranslation } from "@/hooks/useTranslation";

export function useColleagues() {
  const router = useRouter();
  const { t } = useTranslation();
  const [isConnecting, setIsConnecting] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const connectColleague = async (handle: string) => {
    try {
      setIsConnecting(true);
      setError(null);

      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/connect-colleague`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({ handle }),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/login");
        return;
      }

      if (!response.ok) {
        if (response.status === 404) {
          throw new Error("profile.error.userNotFound");
        }
        if (response.status === 422) {
          throw new Error("profile.error.cannotConnect");
        }
        throw new Error("profile.error.connectionFailed");
      }
    } catch (err) {
      setError(
        err instanceof Error ? err : new Error(t("common.error.serverError"))
      );
      throw err;
    } finally {
      setIsConnecting(false);
    }
  };

  return {
    connectColleague,
    isConnecting,
    error,
  };
}
