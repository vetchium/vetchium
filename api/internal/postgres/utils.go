package postgres

import (
	"fmt"

	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (p *PG) convertToOrgUserRoles(
	dbRoles []string,
) ([]vetchi.OrgUserRole, error) {
	var roles []vetchi.OrgUserRole
	for _, str := range dbRoles {
		role := vetchi.OrgUserRole(str)
		switch role {
		case vetchi.Admin,
			vetchi.CostCentersCRUD,
			vetchi.CostCentersViewer,
			vetchi.LocationsCRUD,
			vetchi.LocationsViewer,
			vetchi.OpeningsCRUD,
			vetchi.OpeningsViewer,
			vetchi.OrgUsersCRUD,
			vetchi.OrgUsersViewer:
			roles = append(roles, role)
		default:
			p.log.Error("invalid role in the database", "role", str)
			return nil, fmt.Errorf("invalid role: %s", str)
		}
	}
	return roles, nil
}

func convertOrgUserRolesToStringArray(
	roles []vetchi.OrgUserRole,
) []string {
	var strRoles []string
	for _, role := range roles {
		strRoles = append(strRoles, string(role))
	}
	return strRoles
}
