"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import UserInvite from "@/components/UserInvite";
import { useTranslation } from "@/hooks/useTranslation";
import Container from "@mui/material/Container";
import Typography from "@mui/material/Typography";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function Settings() {
  const { t } = useTranslation();
  const router = useRouter();

  // Auth check
  useEffect(() => {
    const token = Cookies.get("session_token");
    if (!token) {
      router.push("/login");
    }
  }, [router]);

  return (
    <AuthenticatedLayout>
      <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          {t("settings.title")}
        </Typography>

        {/* Invite User Section */}
        <UserInvite />
      </Container>
    </AuthenticatedLayout>
  );
}
