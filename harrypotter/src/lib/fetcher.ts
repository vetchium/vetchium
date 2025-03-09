interface FetcherOptions extends RequestInit {
  body?: string;
}

interface APIError extends Error {
  status?: number;
  data?: unknown;
}

export async function fetcher(url: string, options?: FetcherOptions) {
  const response = await fetch(url, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options?.headers,
    },
  });

  if (!response.ok) {
    const error = new Error("API request failed") as APIError;
    const data = await response.json().catch(() => ({}));
    error.status = response.status;
    error.data = data;
    throw error;
  }

  return response.json();
}
