import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

@doc("Represents a post made by an employer")
model EmployerPost {
    @doc("The unique identifier for the employer post")
    id: string;

    @doc("The content of the employer post")
    @minLength(1)
    @maxLength(4096) // Assuming same limits as user posts for now
    content: string;

    @doc("Tags associated with the employer post")
    @maxItems(3) // Assuming same limits as user posts for now
    tags: string[];

    @doc("The name of the employer who made the post")
    employer_name: string;

    @doc("The primary domain name of the employer")
    employer_domain_name: string;

    @doc("The timestamp when the employer post was created")
    created_at: string; // Consider using `utcDateTime` if more specific type is available

    @doc("The timestamp when the employer post was last updated")
    updated_at: string; // Consider using `utcDateTime`
}
