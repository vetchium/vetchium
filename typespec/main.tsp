import "@typespec/http";
import "@typespec/openapi3";
import "@typespec/rest";

import "./employer/applications.tsp";
import "./employer/auth.tsp";
import "./employer/costcenters.tsp";
import "./employer/locations.tsp";
import "./employer/orgusers.tsp";
import "./employer/openings.tsp";
import "./employer/candidacy.tsp";
import "./employer/interviews.tsp";

import "./hub/applications.tsp";
import "./hub/hubusers.tsp";
import "./hub/openings.tsp";
import "./hub/candidacy.tsp";
import "./hub/interviews.tsp";
import "./hub/employers.tsp";
import "./hub/workhistory.tsp";
import "./hub/profilepage.tsp";
import "./hub/colleagues.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

@service({
    title: "Vetchi API",
})
@server("http://example.com", "blah")
namespace Vetchi;
