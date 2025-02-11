import { useState, useEffect } from "react";

interface UseMyHandleResult {
  myHandle: string | null;
  isLoading: boolean;
  error: Error | null;
}

export function useMyHandle(): UseMyHandleResult {
  const [myHandle, setMyHandle] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    async function fetchMyHandle() {
      try {
        const response = await fetch("/api/hub/get-my-handle");

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
  }, []);

  return { myHandle, isLoading, error };
}
