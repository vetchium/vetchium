"use client";

import { useMyHandle } from "@/hooks/useMyHandle";
import { useParams } from "next/navigation";
import { WorkHistory } from "./WorkHistory";
import { useTranslation } from "@/hooks/useTranslation";
import CircularProgress from "@mui/material/CircularProgress";
import Box from "@mui/material/Box";

export default function ProfilePage() {
  const params = useParams();
  const userHandle = params.handle as string;
  const { myHandle, isLoading: isLoadingHandle } = useMyHandle();
  const { t } = useTranslation();
  const isOwnProfile = myHandle === userHandle;

  if (isLoadingHandle) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", p: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-2xl font-bold mb-8">
        {isOwnProfile ? "My Profile" : `${userHandle}'s Profile`}
      </h1>

      <div className="bg-white rounded-lg shadow-md p-6">
        <WorkHistory userHandle={userHandle} canEdit={isOwnProfile} />
      </div>
    </div>
  );
}
