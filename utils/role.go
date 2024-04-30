package utils

const (
	DepositorRole = "depositor"
	BankerRole    = "banker"
)

func HasPermission(userRole string, accessibleRoles []string) bool {
	for _, role := range accessibleRoles {
		if userRole == role {
			return true
		}
	}
	return false
}
