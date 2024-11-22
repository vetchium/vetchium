import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

enum LocationState {
    ActiveLocation: "ACTIVE_LOCATION",
    DefunctLocation: "DEFUNCT_LOCATION",
}

@doc("Office Location of an Employer")
model Location {
    @minLength(3)
    @maxLength(32)
    title: string;

    country_code: CountryCode;

    @maxLength(1024)
    postal_address: string;

    postal_code: string;

    @maxLength(255)
    openstreetmap_url?: url;

    @maxItems(3)
    city_aka?: City[];

    state: LocationState;
}

model AddLocationRequest {
    @minLength(3)
    @maxLength(32)
    title: string;

    country_code: CountryCode;

    @maxLength(1024)
    postal_address: string;

    postal_code: string;
    openstreetmap_url?: url;

    @maxItems(3)
    city_aka?: City[];
}

model DefunctLocationRequest {
    @minLength(3)
    @maxLength(32)
    title: string;
}

model GetLocationRequest {
    @minLength(3)
    @maxLength(32)
    title: string;
}

model GetLocationsRequest {
    @minLength(3)
    @maxLength(32)
    pagination_key?: string;

    @maxValue(100)
    @doc("Number of Locations to return; 100 is the default if not specified")
    limit?: integer;
}

model RenameLocationRequest {
    @minLength(3)
    @maxLength(32)
    old_title: string;

    @minLength(3)
    @maxLength(32)
    new_title: string;
}

model UpdateLocationRequest {
    @minLength(3)
    @maxLength(32)
    title: string;

    country_code: CountryCode;

    @maxLength(1024)
    postal_address: string;

    postal_code: string;
    openstreetmap_url?: url;

    @maxItems(3)
    city_aka?: City[];
}

// TODO Should we add a max limit to the number of locations that can be added ?
@route("/employer/add-location")
interface AddLocation {
    @tag("Locations")
    @doc("Requires any of ${Admin}, ${LocationsCRUD} roles")
    @post
    @useAuth(EmployerAuth)
    addLocation(@body addLocationRequest: AddLocationRequest): {
        @statusCode statusCode: 200;
    } | {
        @doc("Location already exists")
        @statusCode
        statusCode: 409;
    } | {
        @statusCode statusCode: 400;
        @body error: ValidationErrors;
    };
}

@route("/employer/defunct-location")
interface DefunctLocation {
    @tag("Locations")
    @doc("Requires any of ${Admin}, ${LocationsCRUD} roles")
    @post
    @useAuth(EmployerAuth)
    defunctLocation(@body defuntLocationRequest: DefunctLocationRequest): {
        @statusCode statusCode: 200;
    };
}

@route("/employer/get-location")
interface GetLocation {
    @tag("Locations")
    @doc("Requires any of ${Admin}, ${LocationsCRUD}, ${LocationsViewer} roles")
    @post
    @useAuth(EmployerAuth)
    getLocation(@body getLocation: GetLocationRequest): {
        @statusCode statusCode: 200;
        @body location: Location;
    };
}

@route("/employer/get-locations")
interface GetLocations {
    @tag("Locations")
    @doc("Requires any of ${Admin}, ${LocationsCRUD}, ${LocationsViewer} roles")
    @post
    @useAuth(EmployerAuth)
    getLocations(@body getLocationsRequest: GetLocationsRequest): {
        @statusCode statusCode: 200;
        @body locations: Location[];
    };
}

@route("/employer/rename-location")
interface RenameLocation {
    @tag("Locations")
    @post
    @useAuth(EmployerAuth)
    renameLocation(@body renameLocationRequest: RenameLocationRequest): {
        @statusCode statusCode: 200;
    };
}

@route("/employer/update-location")
interface UpdateLocation {
    @tag("Locations")
    @doc("Requires any of ${Admin}, ${LocationsCRUD} roles")
    @post
    @useAuth(EmployerAuth)
    updateLocation(@body updateLocationRequest: UpdateLocationRequest): {
        @statusCode statusCode: 200;
    };
}
