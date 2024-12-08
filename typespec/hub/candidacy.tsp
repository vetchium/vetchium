import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/applications.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

model AddHubCandidacyCommentRequest {
    candidacy_id: string;

    @maxLength(2048)
    comment: string;
}

@route("/hub/add-candidacy-comment")
interface AddHubCandidacyComment {
    @post
    addComment(@body request: AddHubCandidacyCommentRequest): {
        @statusCode statusCode: 200;
    };
}
