import "@typespec/http";
import "@typespec/openapi3";
import "@typespec/rest";

import "./common/achievements.tsp";
import "./common/applications.tsp";
import "./common/candidacies.tsp";
import "./common/common.tsp";
import "./common/education.tsp";
import "./common/interviews.tsp";
import "./common/openings.tsp";
import "./common/vtags.tsp";

import "./employer/achievements.tsp";
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
import "./employer/settings.tsp";

import "./hub/achievements.tsp";
import "./hub/applications.tsp";
import "./hub/candidacy.tsp";
import "./hub/colleagues.tsp";
import "./hub/education.tsp";
import "./hub/employers.tsp";
import "./hub/hubusers.tsp";
import "./hub/interviews.tsp";
import "./hub/openings.tsp";
import "./hub/posts.tsp";
import "./hub/profilepage.tsp";
import "./hub/workhistory.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

@service({
    title: "Vetchium API",
})
@server("http://example.com", "blah")
namespace Vetchium;
