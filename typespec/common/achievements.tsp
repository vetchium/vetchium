import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

namespace Vetchi;

union AchievementType {
    PATENT: "PATENT",
    PUBLICATION: "PUBLICATION",
    CERTIFICATION: "CERTIFICATION",
}

model Achievement {
    @doc("The unique identifier of the achievement. This will be non-empty only for the hubuser who owns the achievement. For everyone else, it will be empty.")
    id: string;

    type: AchievementType;
    title: string;
    description: string;
    url: string;
    at: utcDateTime;
}
