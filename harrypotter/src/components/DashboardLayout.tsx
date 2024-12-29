"use client";

import AuthenticatedLayout from "./AuthenticatedLayout";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <AuthenticatedLayout>{children}</AuthenticatedLayout>;
}
