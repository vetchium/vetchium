import "@typespec/http";
import "@typespec/openapi3";
import "@typespec/rest";

import "./common.tsp";
import "./costcenters.tsp";
import "./locations.tsp";
import "./orgusers.tsp";
import "./openings.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

@service({
    title: "Vetchi API",
})
@server("http://example.com", "blah")
namespace Vetchi;
