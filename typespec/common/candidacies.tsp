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

model Candidacy {
    candidacy_id: string;
    opening_id: string;
    opening_title: string;
    opening_description: string;
    candidacy_state: CandidacyState;
    applicant_name: string;
    applicant_handle: string;
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
