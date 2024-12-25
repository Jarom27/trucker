package commands

type MessageBuilder interface {
	BuildChecksum(Checksum) []byte
	BuildFlags(flag []byte)
	BuildBody(body []byte)
	BuildHeaders(headers []byte)
	GetResult() []byte
}
