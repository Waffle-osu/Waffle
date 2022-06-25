package requests

type Request interface {
	IrcReplyCode() int
	b1816PacketId() int

	IrcData() []byte
	b1816Data() []byte
}
