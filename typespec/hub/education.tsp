import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/openings.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

model AddEducationRequest {
    @doc("The domain of the institute")
    institute_domain: Domain;

    @minLength(3)
    @maxLength(64)
    @doc("TODO: In future, we could make the institutes to add the unique list of degrees that they offer and allow users to only select from that list. But that can be done only after Vetchi becomes a big platform.")
    degree?: string;

    @doc("If specified, the start date should be earlier than the end date")
    start_date?: plainDate;

    @doc("If specified, the end date should be later than the start date")
    end_date?: plainDate;

    @maxLength(1024)
    description?: string;
}

model AddEducationResponse {
    education_id: string;
}

model Institute {
    domain: Domain;
    name: string;
}

model FilterInstitutesRequest {
    @minLength(3)
    @maxLength(12)
    @doc("The prefix of the institute name or domain")
    prefix?: string;
}

model DeleteEducationRequest {
    education_id: string;
}

model Education {
    @doc("This would be non-empty only for the ListEducationResponse of the logged in user. For other users, this would be empty")
    id?: string;

    institute_domain: Domain;
    degree?: string;
    start_date?: plainDate;
    end_date?: plainDate;
    description?: string;
}

model ListEducationRequest {
    @doc("The user handle of the user to list the education for. If not provided, the education of the current user will be listed.")
    user_handle?: string;
}

@route("/hub/add-education")
interface AddEducation {
    @post
    @useAuth(HubAuth)
    @tag("Education")
    addEducation(@body request: AddEducationRequest): {
        @statusCode statusCode: 200;
        @body AddEducationResponse: AddEducationResponse;
    } | {
        @doc("Max limit of education reached")
        @statusCode
        statusCode: 422;
    } | {
        @statusCode statusCode: 400;
        @body ValidationErrors: ValidationErrors;
    };
}

@route("/hub/filter-institutes")
interface FilterInstitutes {
    @post
    @useAuth(HubAuth)
    @tag("Education")
    filterInstitutes(@body request: FilterInstitutesRequest): {
        @statusCode statusCode: 200;
        @body Institutes: Institute[];
    };
}

@route("/hub/delete-education")
interface DeleteEducation {
    @post
    @useAuth(HubAuth)
    @tag("Education")
    deleteEducation(@body request: DeleteEducationRequest): {
        @statusCode statusCode: 200;
    };
}

@route("/hub/list-education")
interface ListEducation {
    @post
    @useAuth(HubAuth)
    @tag("Education")
    listEducation(@body request: ListEducationRequest): {
        @statusCode statusCode: 200;
        @body Education: Education[];
    };
}
