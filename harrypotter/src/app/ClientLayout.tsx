"use client";

import { AppRouterCacheProvider } from "@mui/material-nextjs/v14-appRouter";
import ThemeRegistry from "@/components/ThemeRegistry";
import { AuthProvider } from "@/contexts/AuthContext";

export default function ClientLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <AppRouterCacheProvider>
      <ThemeRegistry>
        <AuthProvider>{children}</AuthProvider>
      </ThemeRegistry>
    </AppRouterCacheProvider>
  );
}
