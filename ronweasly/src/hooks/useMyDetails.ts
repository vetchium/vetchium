import { config } from "@/config";
import { Handle, HubUserTier } from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useEffect, useState } from "react";

// TODO: Replace with imported type when available in @vetchium/typespec
interface MyDetails {
  handle: Handle;
  full_name: string;
  tier: HubUserTier;
}

interface UseMyDetailsResult {
  details: MyDetails | null;
  isLoading: boolean;
  error: Error | null;
}

export const useMyDetails = (): UseMyDetailsResult => {
  const [details, setDetails] = useState<MyDetails | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    async function fetchDetails() {
      setIsLoading(true);
      setError(null);
      try {
        const token = Cookies.get("session_token");
        if (!token) {
          throw new Error("User not authenticated");
        }

        const response = await fetch(
          `${config.API_SERVER_PREFIX}/hub/get-my-details`,
          {
            method: "GET",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${token}`,
            },
          }
        );

        if (!response.ok) {
          throw new Error(
            `Failed to fetch user details: ${response.statusText}`
          );
        }

        const data: MyDetails = await response.json();
        setDetails(data);
      } catch (err) {
        setError(
          err instanceof Error
            ? err
            : new Error("Unknown error fetching user details")
        );
      } finally {
        setIsLoading(false);
      }
    }

    fetchDetails();
  }, []);

  return {
    details,
    isLoading,
    error,
  };
};
