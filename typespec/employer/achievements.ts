import { AchievementType } from "../common/achievements";
import { Handle } from "../common/common";

export interface ListHubUserAchievementsRequest {
  handle: Handle;
  type: AchievementType;
}
