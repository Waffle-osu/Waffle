package chat

type AdminPrivilegable interface {
	IsOfAdminPrivileges() bool
}
