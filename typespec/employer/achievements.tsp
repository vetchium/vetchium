import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/achievements.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

model ListAchievementsRequest {
    handle: Handle;
}

@route("/employer/list-hub-user-achievements")
interface EmployerAchievements {
    @post
    @useAuth(EmployerAuth)
    listAchievements(@body request: ListAchievementsRequest): {
        @statusCode statusCode: 200;
        @body Achievement: Achievement[];
    };
}
