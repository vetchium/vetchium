import { City, CountryCode } from '../common/common';

export type LocationState = 'ACTIVE_LOCATION' | 'DEFUNCT_LOCATION';

export const LocationStates = {
    ACTIVE: 'ACTIVE_LOCATION' as LocationState,
    DEFUNCT: 'DEFUNCT_LOCATION' as LocationState,
} as const;

export interface Location {
    title: string;
    country_code: CountryCode;
    postal_address: string;
    postal_code: string;
    openstreetmap_url?: string;
    city_aka?: City[];
    state: LocationState;
}

export interface AddLocationRequest {
    title: string;
    country_code: CountryCode;
    postal_address: string;
    postal_code: string;
    openstreetmap_url?: string;
    city_aka?: City[];
}

export interface DefunctLocationRequest {
    title: string;
}

export interface GetLocationRequest {
    title: string;
}

export interface GetLocationsRequest {
    states?: LocationState[];
    pagination_key?: string;
    limit?: number;
}

export interface RenameLocationRequest {
    old_title: string;
    new_title: string;
}

export interface UpdateLocationRequest {
    title: string;
    country_code: CountryCode;
    postal_address: string;
    postal_code: string;
    openstreetmap_url?: string;
    city_aka?: City[];
}

export function isValidLocationState(state: string): state is LocationState {
    return Object.values(LocationStates).includes(state as LocationState);
} 