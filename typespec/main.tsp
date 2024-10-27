import "@typespec/http";
import "@typespec/openapi3";
import "@typespec/rest";

import "./common.tsp";
import "./costcenters.tsp";
import "./locations.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

@service({
    title: "Vetchi API",
})
@server("http://example.com", "blah")
namespace Vetchi;
