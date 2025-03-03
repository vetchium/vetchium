import "@typespec/http";

using TypeSpec.Http;

@service({
    title: "Employer Profile APIs",
})
namespace Vetchi.Employer;

model GetHubUserBioRequest {
    @doc("The handle of the hub user whose bio is being requested")
    handle: string;
}

model EmployerWorkHistory {
    id: string;
    employer_domain: string;

    @doc("In case the employer is not yet onboarded, the field will be empty")
    employer_name?: string;

    @doc("The title of the Job")
    title: string;

    start_date: plainDate;
    end_date?: plainDate;

    @maxLength(1024)
    description?: string;
}

@doc("Bio information of a hub user from employer's perspective")
model EmployerViewBio {
    @doc("The handle of the hub user")
    handle: string;

    @doc("The full name of the hub user")
    full_name: string;

    @doc("The short bio of the user. The one that gets displayed below the user's name.")
    short_bio: string;

    @doc("The long bio of the user. Something like 'About Me' of a user")
    long_bio: string;

    @doc("The list of verified mail domains for the user.")
    verified_mail_domains?: string[];

    @doc("The work history of the hubuser")
    work_history: EmployerWorkHistory[];
}

@route("/employer/get-hub-user-bio")
interface GetHubUserBio {
    @tag("HubProfile")
    @post
    @useAuth(EmployerAuth)
    @doc("Get a hub user's bio information. Requires any of ADMIN, APPLICATIONS_CRUD or APPLICATIONS_VIEWER roles")
    getHubUserBio(@body request: GetHubUserBioRequest): {
        @body bio: EmployerViewBio;
        @statusCode statusCode: 200;
    };
}

@route("/employer/get-hub-user-profile-picture/{handle}")
interface GetHubUserProfilePicture {
    @tag("HubProfile")
    @get
    @useAuth(EmployerAuth)
    @doc("Get a hub user's profile picture. Requires any of ADMIN, APPLICATIONS_CRUD or APPLICATIONS_VIEWER roles")
    getHubUserProfilePicture(@path handle: string): {
        @statusCode statusCode: 200;
        @body image: bytes;
    } | {
        @doc("User not found or user has no profile picture")
        @statusCode
        statusCode: 404;
    };
}
