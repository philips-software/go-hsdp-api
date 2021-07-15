package cdl

const (
	ROLE_STUDYMANAGER   = "STUDYMANAGER"
	ROLE_MONITOR        = "MONITOR"
	ROLE_UPLOADER       = "UPLOADER"
	ROLE_DATA_SCIENTIST = "DATASCIENTIST"
)

type RoleRequest struct {
	IAMUserUUID string `json:"IAMuserUUID" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Role        string `json:"role" validate:"required"`
	InstituteID string `json:"instituteID"`
}

type RoleAssignment struct {
	IAMUserUUID string `json:"IAMuserUUID"`
	Roles       []struct {
		Role string `json:"role"`
	} `json:"roles"`
}

type RoleAssignmentResult []RoleAssignment

func (r RoleAssignmentResult) Roles(userUUID string) []string {
	for _, role := range r {
		if role.IAMUserUUID == userUUID {
			var roles []string
			for _, foundRole := range role.Roles {
				roles = append(roles, foundRole.Role)
			}
			return roles
		}
	}
	return []string{}
}
