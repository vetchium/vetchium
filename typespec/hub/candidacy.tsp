import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/applications.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

model AddHubCandidacyCommentRequest {
    candidacy_id: string;

    @maxLength(2048)
    comment: string;
}

model MyCandidacy {
    candidacy_id: string;
    company_name: string;
    company_domain: string;
    opening_id: string;
    opening_title: string;
    opening_description: string;
    candidacy_state: CandidacyState;
}

model MyCandidaciesRequest {
    @doc("The candidacy_id of the last candidacy. Candidacies are returned in reverse chronological order and if two candidacies have the same timestamp, they are further ordered by id.")
    pagination_key?: string;

    @doc("The maximum number of candidacies to return. Defaults to 40.")
    @maxValue(40)
    limit?: integer;
}

@route("/hub/add-candidacy-comment")
interface AddHubCandidacyComment {
    @tag("Candidacies")
    @post
    addComment(@body request: AddHubCandidacyCommentRequest): {
        @statusCode statusCode: 200;
    };
}

@route("/hub/get-my-candidacies")
interface GetMyCandidacies {
    @tag("Candidacies")
    @post
    getMyCandidacies(@body request: MyCandidaciesRequest): {
        @statusCode statusCode: 200;
        @body response: MyCandidacy[];
    };
}

@route("/hub/get-candidacy-info")
interface GetHubCandidacyInfo {
    @tag("Candidacies")
    @post
    getHubCandidacyInfo(@body request: GetCandidacyInfoRequest): {
        @statusCode statusCode: 200;
        @body response: Candidacy;
    };
}

@route("/hub/get-candidacy-comments")
interface GetHubCandidacyComments {
    @tag("Candidacies")
    @post
    getCandidacyComments(@body request: GetCandidacyCommentsRequest): {
        @doc("If an invalid candidacy ID is provided, an empty list is returned.")
        @statusCode
        statusCode: 200;

        @body response: CandidacyComment[];
    };
}
