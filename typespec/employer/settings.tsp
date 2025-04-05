import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/achievements.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

model ChangeCoolOffPeriodRequest {
    @doc("The new cool off period in days. Should be done by users with the role of an admin.")
    @minValue(0)
    @maxValue(365)
    coolOffPeriod: int32;
}

@route("/employer/change-cool-off-period")
interface ChangeCoolOffPeriod {
    @post
    @useAuth(EmployerAuth)
    @tag("Employer Settings")
    changeCoolOffPeriod(@body request: ChangeCoolOffPeriodRequest): void;
}

@route("/employer/get-cool-off-period")
interface GetCoolOffPeriod {
    @get
    @useAuth(EmployerAuth)
    @tag("Employer Settings")
    getCoolOffPeriod(): {
        @doc("The cool off period in days")
        coolOffPeriodDays: int32;
    };
}
