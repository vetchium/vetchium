"use client";

import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export function useAuth() {
  const router = useRouter();

  useEffect(() => {
    const token = Cookies.get("session_token");
    if (!token) {
      router.push("/login");
    }
  }, [router]);

  const isAuthenticated = () => {
    return !!Cookies.get("session_token");
  };

  return { isAuthenticated };
}
