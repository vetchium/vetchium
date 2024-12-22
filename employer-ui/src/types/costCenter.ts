export interface CostCenter {
  name: string;
  notes?: string;
  state: CostCenterState;
}

export enum CostCenterState {
  ACTIVE_CC = 'ACTIVE_CC',
  DEFUNCT_CC = 'DEFUNCT_CC',
}

export interface AddCostCenterRequest {
  name: string;
  notes?: string;
}

export interface UpdateCostCenterRequest {
  name: string;
  notes: string;
}

export interface GetCostCentersRequest {
  pagination_key?: string;
  limit?: number;
}

export interface GetCostCenterRequest {
  name: string;
}

export interface DefunctCostCenterRequest {
  name: string;
}

export interface RenameCostCenterRequest {
  old_name: string;
  new_name: string;
} 