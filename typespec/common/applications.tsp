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
