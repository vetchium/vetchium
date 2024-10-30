@error
model ValidationErrors {
    @doc("List of fields that failed validation")
    errors: string[];
}

@minLength(3)
@maxLength(256)
@doc("A valid email address")
scalar EmailAddress extends string;

scalar OrgUserRole extends string;

alias Admin = "ADMIN";
alias CostCentersCRUD = "COST_CENTERS_CRUD";
alias CostCentersViewer = "COST_CENTERS_VIEWER";
alias LocationsCRUD = "LOCATIONS_CRUD";
alias LocationsViewer = "LOCATIONS_VIEWER";
alias OrgUsersCRUD = "ORG_USERS_CRUD";
alias OrgUsersViewer = "ORG_USERS_VIEWER";
