import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/openings.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

model AddProfessionalEmailRequest {
    email: EmailAddress;
}

model VerifyProfessionalEmailRequest {
    email: EmailAddress;

    @doc("The verification code sent to the email address")
    code: string;
}

model TriggerVerificationRequest {
    email: EmailAddress;
}

model DeleteProfessionalEmailRequest {
    email: EmailAddress;
}

model ProfessionalEmail {
    email: EmailAddress;

    @doc("The date and time when the email address was last verified. If null, the email address was never verified.")
    lastVerifiedAt?: utcDateTime;
}

@route("/hub/add-professional-email")
interface AddProfessionalEmail {
    @tag("HubProfile")
    @post
    @useAuth(HubAuth)
    addProfessionalEmail(@body request: AddProfessionalEmailRequest): {
        @statusCode statusCode: 200;
    };
}

@route("/hub/verify-professional-email")
interface VerifyProfessionalEmail {
    @tag("HubProfile")
    @post
    @useAuth(HubAuth)
    verifyProfessionalEmail(@body request: VerifyProfessionalEmailRequest): {};
}

@route("/hub/trigger-verification")
interface TriggerVerification {
    @tag("HubProfile")
    @post
    @useAuth(HubAuth)
    triggerVerification(@body request: TriggerVerificationRequest): {
        @statusCode statusCode: 200;
    } | {
        @doc("Email address verification not triggered, because its prior verification is still valid")
        @statusCode
        statusCode: 422;
    };
}

@route("/hub/delete-professional-email")
interface DeleteProfessionalEmail {
    @tag("HubProfile")
    @post
    @useAuth(HubAuth)
    deleteProfessionalEmail(@body request: DeleteProfessionalEmailRequest): {};
}

@route("/hub/my-professional-emails")
interface MyProfessionalEmails {
    @tag("HubProfile")
    @get
    @useAuth(HubAuth)
    myProfessionalEmails(): {
        @statusCode statusCode: 200;
        @body emails: ProfessionalEmail[];
    };
}
