import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

scalar OrgUserState extends string;

// The user has been added to the organization but not yet signed up by clicking on the link in the invitation email
alias AddedOrgUser = "ADDED_ORG_USER";
// The user is active in the organization
alias ActiveOrgUser = "ACTIVE_ORG_USER";
// The user is no longer active in the organization
alias DisabledOrgUser = "DISABLED_ORG_USER";
// The user is replicated from a different directory service (e.g. LDAP, Google, Microsoft Active Directory, etc.)
alias ReplicatedOrgUser = "REPLICATED_ORG_USER";

model OrgUser {
    email: EmailAddress;
    name: string;
    state: OrgUserState;
    roles: OrgUserRole[];
}

model AddOrgUserRequest {
    email: EmailAddress;

    @minLength(3)
    @maxLength(256)
    name: string;

    @minItems(1)
    roles: OrgUserRole[];
}

model DisableOrgUserRequest {
    email: EmailAddress;
}

model EnableOrgUserRequest {
    email: EmailAddress;
}

model FilterOrgUsersRequest {
    @doc("If specified, only org users with an email address or name that starts with this prefix will be returned, ordered by the email addresses")
    prefix?: string;

    @minLength(3)
    @maxLength(256)
    pagination_key?: EmailAddress;

    @maxValue(40)
    @doc("The maximum number of org users to return; 40 is the default if not specified")
    limit?: integer;

    @doc("If specified, only org users with this state will be returned. If not specified, all users in either of ${ActiveOrgUser} or ${AddedOrgUser} states will be returned")
    state?: OrgUserState[];
}

model OrgUserShort {
    name: string;
    email: EmailAddress;

    @doc("In case the same email address of the org_user is added by a hub profile, its Vetchi handle is populated here")
    vetchi_handle?: string;
}

model UpdateOrgUserRequest {
    @doc("The email address of the org user to update")
    email: EmailAddress;

    @doc("The new name of the org user")
    name: string;

    @doc("The new roles of the org user")
    roles: OrgUserRole[];
}

@route("/employer/add-org-user")
interface AddOrgUser {
    @tag("OrgUsers")
    @doc("Requires any of ${Admin}, ${OrgUsersCRUD} roles")
    @post
    addOrgUser(@body addOrgUserRequest: AddOrgUserRequest): {
        @statusCode statusCode: 200;
    } | {
        @doc("The org user already exists")
        @statusCode
        statusCode: 409;
    } | {
        @statusCode statusCode: 400;
        @body error: ValidationErrors;
    };
}

@route("/employer/disable-org-user")
interface DisableOrgUser {
    @tag("OrgUsers")
    @doc("Requires anyof ${Admin}, ${OrgUsersCRUD} roles. The last OrgUser with an ${Admin} role cannot be updated")
    @post
    disableOrgUser(@body disableOrgUserRequest: DisableOrgUserRequest): {
        @statusCode statusCode: 200;
    } | {
        @statusCode statusCode: 400;
        @body error: ValidationErrors;
    } | {
        @doc("Last active admin cannot be disabled")
        @statusCode
        statusCode: 403;
    };
}

@route("/employer/enable-org-user")
interface EnableOrgUser {
    @tag("OrgUsers")
    @doc("Requires any of ${Admin}, ${OrgUsersCRUD} roles. The account must be in ${DisabledOrgUser} state")
    @post
    enableOrgUser(@body enableOrgUserRequest: EnableOrgUserRequest): {
        @statusCode statusCode: 200;
    };
}

@route("/employer/filter-org-users")
interface FilterOrgUsers {
    @tag("OrgUsers")
    @doc("Requires any of ${Admin}, ${OrgUsersCRUD}, ${OrgUsersViewer} roles")
    @post
    filterOrgUsers(@body filterOrgUsersRequest: FilterOrgUsersRequest): {
        @statusCode statusCode: 200;
    };
}

model SignupOrgUserRequest {
    email: EmailAddress;
    name: string;
    password: Password;
}

@route("/employer/signup-orguser")
interface SignupOrgUser {
    @tag("OrgUsers")
    @post
    signupOrgUser(@body signupOrgUserRequest: SignupOrgUserRequest): {
        @statusCode statusCode: 200;
    } | {
        @doc("Either Invite Token is not found or expired")
        @statusCode
        statusCode: 403;
    } | {
        @statusCode statusCode: 400;
        @body error: ValidationErrors;
    };
}

@route("/employer/update-org-user")
interface UpdateOrgUser {
    @tag("OrgUsers")
    @doc("Requires any of ${Admin}, ${OrgUsersCRUD} roles")
    @post
    updateOrgUser(@body updateOrgUserRequest: UpdateOrgUserRequest):
        | {
              @statusCode statusCode: 200;
          }
        | {
              @statusCode statusCode: 400;
              @body error: ValidationErrors;
          }
        | {
              @doc("User with same email already exist")
              @statusCode
              statusCode: 409;
          }
        | {
              @doc("Last active admin cannot be updated")
              @statusCode
              statusCode: 403;
          };
}
