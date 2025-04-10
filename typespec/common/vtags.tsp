import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

scalar VTagID extends string;

@maxLength(32)
scalar VTagName extends string;

model VTag {
    @doc("Client should not assume anything about the format of the string. Could be UUID or XID or something random.")
    id: VTagID;

    name: VTagName;

    // In future we can add more fields like Icons
}
