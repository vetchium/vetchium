import { config } from "@/config";
import { HubUserTier } from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useEffect, useState } from "react";

interface UseMyTierResult {
  tier: HubUserTier | null;
  isLoading: boolean;
  error: Error | null;
}

export const useMyTier = (): UseMyTierResult => {
  const [tier, setTier] = useState<HubUserTier | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    async function fetchTier() {
      setIsLoading(true);
      setError(null);
      try {
        const token = Cookies.get("session_token");
        if (!token) {
          // Handle not logged in state if necessary, maybe redirect or set error
          // For now, just setting an error
          throw new Error("User not authenticated");
        }

        const response = await fetch(
          `${config.API_SERVER_PREFIX}/hub/my-tier`,
          {
            method: "POST", // Assuming POST based on other examples, adjust if GET
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${token}`,
            },
          }
        );

        if (!response.ok) {
          // Handle specific errors like 401 Unauthorized if needed
          throw new Error(`Failed to fetch tier: ${response.statusText}`);
        }

        const data = await response.json();
        setTier(data as HubUserTier); // Assuming the API returns the tier directly
      } catch (err) {
        setError(
          err instanceof Error ? err : new Error("Unknown error fetching tier")
        );
      } finally {
        setIsLoading(false);
      }
    }

    fetchTier();
  }, []);

  return {
    tier,
    isLoading,
    error,
  };
};
