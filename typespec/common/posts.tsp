import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

model EmployerPost {
    id: string;
    content: string;
    tags: string[];
    company_domain: string;
    created_at: string;
}
