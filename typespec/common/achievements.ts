export enum AchievementType {
  PATENT = "PATENT",
  PUBLICATION = "PUBLICATION",
  CERTIFICATION = "CERTIFICATION",
}

export interface Achievement {
  id: string;
  type: AchievementType;
  title: string;
  description: string;
  url: string;
}
