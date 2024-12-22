export interface Location {
  title: string;
  country_code: string;
  postal_address: string;
  postal_code: string;
  openstreetmap_url?: string;
  city_aka?: string[];
  state: LocationState;
}

export enum LocationState {
  ACTIVE_LOCATION = 'ACTIVE_LOCATION',
  DEFUNCT_LOCATION = 'DEFUNCT_LOCATION',
}

export interface AddLocationRequest {
  title: string;
  country_code: string;
  postal_address: string;
  postal_code: string;
  openstreetmap_url?: string;
  city_aka?: string[];
}

export interface UpdateLocationRequest extends AddLocationRequest {
  title: string;
}

export interface GetLocationsRequest {
  pagination_key?: string;
  limit?: number;
}

export interface GetLocationRequest {
  title: string;
}

export interface DefunctLocationRequest {
  title: string;
} 