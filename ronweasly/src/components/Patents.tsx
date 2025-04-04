"use client";

import { AchievementType, Handle } from "@vetchium/typespec";
import { AchievementSection } from "./Achievement";

interface PatentsProps {
  userHandle: Handle;
  canEdit: boolean;
}

export function Patents({ userHandle, canEdit }: PatentsProps) {
  return (
    <AchievementSection
      userHandle={userHandle}
      achievementType={AchievementType.PATENT}
      canEdit={canEdit}
    />
  );
}
