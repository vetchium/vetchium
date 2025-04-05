import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/achievements.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

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

model ListAchievementsRequest {
    type: AchievementType;

    @doc("If passed, the achievements of that user will alone be returned. If not passed, the logged in hubuser's achievements will be returned.")
    handle?: Handle;
    // TODO: Should we paginate this API ?
}

model DeleteAchievementRequest {
    @doc("The unique identifier of the achievement")
    id: string;
}

@route("/hub/add-achievement")
interface HubAchievements {
    @post
    @useAuth(HubAuth)
    @tag("Achievements")
    addAchievement(@body request: AddAchievementRequest): {
        @statusCode statusCode: 200;
        @body AddAchievementResponse: AddAchievementResponse;
    };
}

@route("/hub/list-achievements")
interface ListAchievements {
    @post
    @useAuth(HubAuth)
    @tag("Achievements")
    listAchievements(@body request: ListAchievementsRequest): {
        @statusCode statusCode: 200;
        @body Achievements: Achievement[];
    };
}

@route("/hub/delete-achievement")
interface DeleteAchievement {
    @post
    @useAuth(HubAuth)
    @tag("Achievements")
    deleteAchievement(@body request: DeleteAchievementRequest): {
        @statusCode statusCode: 200;
    };
}
