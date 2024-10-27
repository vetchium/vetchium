import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

@route("/employer/add-location")
interface AddLocation {
    @tag("Locations")
    @post
    addLocation(@body location: Location): {
        @statusCode statusCode: 200;
    };
}

@route("/employer/defunct-location")
interface DefunctLocation {
    @tag("Locations")
    @post
    defunctLocation(@body defuntLocationRequest: DefunctLocationRequest): {
        @statusCode statusCode: 200;
    };
}

@route("/employer/get-location")
interface GetLocation {
    @tag("Locations")
    @post
    getLocation(@body getLocation: GetLocationRequest): {
        @statusCode statusCode: 200;
        @body location: Location;
    };
}

@route("/employer/get-locations")
interface GetLocations {
    @tag("Locations")
    @post
    getLocations(@body getLocationsRequest: GetLocationsRequest): {
        @statusCode statusCode: 200;
        @body locations: Location[];
    };
}

@route("/employer/rename-location")
interface RenameLocation {
    @tag("Locations")
    @post
    renameLocation(@body renameLocationRequest: RenameLocationRequest): {
        @statusCode statusCode: 200;
    };
}

@route("/employer/update-location")
interface UpdateLocation {
    @tag("Locations")
    @post
    updateLocation(@body updateLocation: Location): {
        @statusCode statusCode: 200;
    };
}

@minLength(3)
@maxLength(32)
scalar City extends string;

@minLength(3)
@maxLength(3)
@doc("ISO 3166-1 alpha-3 code for the Country")
scalar CountryCode extends string;

@doc("Office Location of an Employer")
model Location {
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
