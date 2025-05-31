"use client";

import { useAuth } from "@/hooks/useAuth";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function HomePage() {
  const router = useRouter();
  useAuth(); // Check authentication and redirect if not authenticated

  useEffect(() => {
    // Redirect to posts page since it's now the default
    router.replace("/posts");
  }, [router]);

  // Return null since we're redirecting
  return null;
}
