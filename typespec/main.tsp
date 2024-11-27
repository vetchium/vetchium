import "@typespec/http";
import "@typespec/openapi3";
import "@typespec/rest";

import "./employer/auth.tsp";
import "./employer/costcenters.tsp";
import "./employer/locations.tsp";
import "./employer/orgusers.tsp";
import "./employer/openings.tsp";

import "./hub/hubusers.tsp";
import "./hub/openings.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

@service({
    title: "Vetchi API",
})
@server("http://example.com", "blah")
namespace Vetchi;
