"use client";

import { AchievementType, Handle } from "@psankar/vetchi-typespec";
import { AchievementSection } from "./Achievement";

interface PublicationsProps {
  userHandle: Handle;
  canEdit: boolean;
}

export function Publications({ userHandle, canEdit }: PublicationsProps) {
  return (
    <AchievementSection
      userHandle={userHandle}
      achievementType={AchievementType.PUBLICATION}
      canEdit={canEdit}
    />
  );
}
