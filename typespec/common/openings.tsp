import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

union OpeningState {
    // Initial state of an Opening. Applicants are not allowed. Occasionally can be reverted to this state by the employer, to prevent new Applicants.
    DraftOpening: "DRAFT_OPENING",

    // Applicants are allowed to apply.
    ActiveOpening: "ACTIVE_OPENING",

    // Vetchi moves an Opening to this state when there are many pending unprocessed Applications for this Opening. New Applicants are not allowed.
    SuspendedOpening: "SUSPENDED_OPENING",

    // Applicants are no longer allowed to apply. All the positions needed were either filled or can no longer be filled.
    ClosedOpening: "CLOSED_OPENING",
}

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

model Salary {
    // salary values are in decimal and not integer because of crypto currencies

    @minValue(0)
    min_amount: decimal;

    @minValue(1)
    max_amount: decimal;

    currency: Currency;
}

model FilterOpeningTagsRequest {
    @doc("The prefix that should be used to filter the tags by name")
    prefix?: string;
}

scalar OpeningTagID extends string;

model OpeningTag {
    @doc("Client should not assume anything about the format of the string. Could be UUID or XID or something random.")
    id: OpeningTagID;

    name: string;

    // In future we can add more fields like Icons
}
