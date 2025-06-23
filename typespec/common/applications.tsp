import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

union ApplicationState {
    Applied: "APPLIED",
    Rejected: "REJECTED",
    Shortlisted: "SHORTLISTED",
    Withdrawn: "WITHDRAWN",
    Expired: "EXPIRED",
}

@doc("Score from a specific model")
model ModelScore {
    @doc("Name of the model that generated the score")
    model_name: string;

    @doc("Score value from 0 to 100")
    @minValue(0)
    @maxValue(100)
    score: int32;
}
