import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/achievements.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

model ListHubUserAchievementsRequest {
    handle: Handle;
    type: AchievementType;
}

@route("/employer/list-hub-user-achievements")
interface ListHubUserAchievements {
    @post
    @useAuth(EmployerAuth)
    listHubUserAchievements(@body request: ListHubUserAchievementsRequest): {
        @statusCode statusCode: 200;
        @body Achievements: Achievement[];
    };
}
