import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Cookies from "js-cookie";
import { config } from "@/config";

interface UseMyHandleResult {
  myHandle: string | null;
  isLoading: boolean;
  error: Error | null;
}

export function useMyHandle(): UseMyHandleResult {
  const router = useRouter();
  const [myHandle, setMyHandle] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    async function fetchMyHandle() {
      try {
        const token = Cookies.get("session_token");
        if (!token) {
          router.push("/login");
          return;
        }

        const response = await fetch(
          `${config.API_SERVER_PREFIX}/hub/get-my-handle`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${token}`,
            },
          }
        );

        if (response.status === 401) {
          Cookies.remove("session_token");
          router.push("/login");
          return;
        }

        if (!response.ok) {
          throw new Error("Failed to fetch user handle");
        }

        const data = await response.json();
        setMyHandle(data.handle);
      } catch (err) {
        setError(err instanceof Error ? err : new Error("Unknown error"));
      } finally {
        setIsLoading(false);
      }
    }

    fetchMyHandle();
  }, [router]);

  return { myHandle, isLoading, error };
}
