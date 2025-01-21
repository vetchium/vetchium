import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "./common.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

model GetCandidacyInfoRequest {
    candidacy_id: string;
}

model GetCandidacyCommentsRequest {
    candidacyID: string;
    // TODO: Should we paginate ?
}

model CandidacyComment {
    commentID: string;
    commenterName: string;
    commenterType: string;
    content: string;
    createdAt: utcDateTime;
}
