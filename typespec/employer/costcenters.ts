export type CostCenterName = string;

export type CostCenterState = 'ACTIVE_CC' | 'DEFUNCT_CC';

export const CostCenterStates = {
    ACTIVE: 'ACTIVE_CC' as CostCenterState,
    DEFUNCT: 'DEFUNCT_CC' as CostCenterState,
} as const;

export interface CostCenter {
    name: CostCenterName;
    notes?: string;
    state: CostCenterState;
}

export interface AddCostCenterRequest {
    name: CostCenterName;
    notes?: string;
}

export interface DefunctCostCenterRequest {
    name: CostCenterName;
}

export interface GetCostCenterRequest {
    name: CostCenterName;
}

export interface GetCostCentersRequest {
    limit?: number;
    pagination_key?: CostCenterName;
    states?: CostCenterState[];
}

export interface RenameCostCenterRequest {
    old_name: CostCenterName;
    new_name: CostCenterName;
}

export interface UpdateCostCenterRequest {
    name: CostCenterName;
    notes: string;
}

export function isValidCostCenterState(state: string): state is CostCenterState {
    return Object.values(CostCenterStates).includes(state as CostCenterState);
} 