import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

union OpeningType {
    FullTimeOpening: "FULL_TIME_OPENING",
    PartTimeOpening: "PART_TIME_OPENING",
    ContractOpening: "CONTRACT_OPENING",
    InternshipOpening: "INTERNSHIP_OPENING",
    UnspecifiedOpening: "UNSPECIFIED_OPENING",
}

union EducationLevel {
    BachelorEducation: "BACHELOR_EDUCATION",
    MasterEducation: "MASTER_EDUCATION",
    DoctorateEducation: "DOCTORATE_EDUCATION",
    NotMattersEducation: "NOT_MATTERS_EDUCATION",
    UnspecifiedEducation: "UNSPECIFIED_EDUCATION",
}
