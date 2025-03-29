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
    id: string;
    type: AchievementType;
    title: string;
    description: string;
    url: string;
    at: utcDateTime;
}
