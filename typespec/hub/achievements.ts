import { AchievementType } from "../common/achievements";
import { Handle } from "../common/common";

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

export interface ListAchievementsRequest {
  type: AchievementType;
  handle?: Handle;
  // TODO: Should we paginate this API ?
}

export interface DeleteAchievementRequest {
  id: string;
}
