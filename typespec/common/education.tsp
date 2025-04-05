import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

namespace Vetchium;

model Institute {
    domain: Domain;
    name: string;
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
