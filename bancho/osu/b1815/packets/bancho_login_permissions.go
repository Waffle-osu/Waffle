package packets

const (
	UserPermissionsRegular   = 1
	UserPermissionsBAT       = 2
	UserPermissionsSupporter = 4
	UserPermissionsFriend    = 8
)

func BanchoSendLoginPermissions(packetQueue chan BanchoPacket, permissions int32) {
	BanchoSendIntPacket(packetQueue, BanchoLoginPermissions, permissions)
}
