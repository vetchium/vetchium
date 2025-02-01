interface FetcherOptions extends RequestInit {
  body?: string;
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
    const error = new Error("API request failed");
    const data = await response.json().catch(() => ({}));
    (error as any).status = response.status;
    (error as any).data = data;
    throw error;
  }

  return response.json();
}
