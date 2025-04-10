export interface FilterVTagsRequest {
  prefix?: string;
}

export type VTagID = string;

export type VTagName = string;

export interface VTag {
  id: VTagID;
  name: VTagName;
}
