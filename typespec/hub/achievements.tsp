import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/achievements.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

model AddAchievementRequest {
    type: AchievementType;

    @minLength(3)
    @maxLength(128)
    title: string;

    @minLength(3)
    @maxLength(1024)
    description?: string;

    @minLength(3)
    @maxLength(1024)
    url?: string;

    at?: utcDateTime;
}

model AddAchievementResponse {
    id: string;
}

@route("/hub/add-achievement")
interface HubAchievements {
    @post
    @useAuth(HubAuth)
    addAchievement(@body request: AddAchievementRequest): {
        @statusCode statusCode: 200;
        @body AddAchievementResponse: AddAchievementResponse;
    };
}
