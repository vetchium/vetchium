@error
model ValidationErrors {
    @doc("List of fields that failed validation")
    errors: string[];
}

alias Admin = "ADMIN";
alias CostCentersCRUD = "COST_CENTERS_CRUD";
alias CostCentersViewer = "COST_CENTERS_VIEWER";
alias LocationsCRUD = "LOCATIONS_CRUD";
alias LocationsViewer = "LOCATIONS_VIEWER";
