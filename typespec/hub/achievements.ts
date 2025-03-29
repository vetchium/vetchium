import { AchievementType } from "../common/achievements";

export interface AddAchievementRequest {
  type: AchievementType;
  title: string;
  description?: string;
  url?: string;
  at?: Date;
}

export interface AddAchievementResponse {
  id: string;
}
