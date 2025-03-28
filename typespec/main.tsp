import "@typespec/http";
import "@typespec/openapi3";
import "@typespec/rest";

import "./employer/applications.tsp";
import "./employer/auth.tsp";
import "./employer/candidacy.tsp";
import "./employer/costcenters.tsp";
import "./employer/education.tsp";
import "./employer/interviews.tsp";
import "./employer/locations.tsp";
import "./employer/openings.tsp";
import "./employer/orgusers.tsp";
import "./employer/profilepage.tsp";

import "./hub/applications.tsp";
import "./hub/candidacy.tsp";
import "./hub/colleagues.tsp";
import "./hub/education.tsp";
import "./hub/employers.tsp";
import "./hub/hubusers.tsp";
import "./hub/interviews.tsp";
import "./hub/openings.tsp";
import "./hub/profilepage.tsp";
import "./hub/workhistory.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

@service({
    title: "Vetchi API",
})
@server("http://example.com", "blah")
namespace Vetchi;
