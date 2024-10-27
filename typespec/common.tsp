@error
model ValidationErrors {
    @doc("List of fields that failed validation")
    errors: string[];
}

enum OrgUserRole {
    Admin: "ADMIN",
    CostCentersCRUD: "COST_CENTERS_CRUD",
    CostCentersViewer: "COST_CENTERS_VIEWER",
}
