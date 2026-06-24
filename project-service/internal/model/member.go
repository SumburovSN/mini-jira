package model

type Role string

const (
	RoleOwner  Role = "owner"
	RoleEditor Role = "editor"
	RoleViewer Role = "viewer"
)

type Member struct {
	ID        int
	ProjectID int
	UserID    int
	Role      Role
}

func (r Role) CanEdit() bool {
	return r == RoleOwner || r == RoleEditor
}

func (r Role) CanView() bool {
	return r == RoleOwner ||
		r == RoleEditor ||
		r == RoleViewer
}

func (r Role) IsOwner() bool {
	return r == RoleOwner
}

func (r Role) IsValid() bool {
	switch r {
	case RoleOwner, RoleEditor, RoleViewer:
		return true
	default:
		return false
	}
}
