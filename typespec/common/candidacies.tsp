import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "./common.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

model GetCandidacyInfoRequest {
    candidacy_id: string;
}

model GetCandidacyCommentsRequest {
    candidacyID: string;
    // TODO: Should we paginate ?
}

union CommenterType {
    ORG_USER: "ORG_USER",
    HUB_USER: "HUB_USER",
}

model CandidacyComment {
    commentID: string;
    commenterName: string;
    commenterType: CommenterType;
    content: string;
    createdAt: utcDateTime;
}
