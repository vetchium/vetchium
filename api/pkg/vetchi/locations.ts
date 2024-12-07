import { CountryCode } from "./common";

export type LocationState = "ACTIVE_LOCATION" | "DEFUNCT_LOCATION";

export interface Location {
  title: string;
  country_code: string;
  postal_address: string;
  postal_code: string;
  openstreetmap_url: string;
  city_aka: string[];
  state: LocationState;
}

export interface AddLocationRequest {
  title: string;
  country_code: CountryCode;
  postal_address: string;
  postal_code: string;
  openstreetmap_url?: string;
  city_aka?: string[];
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
  city_aka?: string[];
}
