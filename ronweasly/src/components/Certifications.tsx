"use client";

import { AchievementType, Handle } from "@psankar/vetchi-typespec";
import { AchievementSection } from "./Achievement";

interface CertificationsProps {
  userHandle: Handle;
  canEdit: boolean;
}

export function Certifications({ userHandle, canEdit }: CertificationsProps) {
  return (
    <AchievementSection
      userHandle={userHandle}
      achievementType={AchievementType.CERTIFICATION}
      canEdit={canEdit}
    />
  );
}
