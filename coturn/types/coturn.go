package types

type Packets struct {
	RecvPackets int64
	RecvBytes   int64
	SentPackets int64
	SentBytes   int64
}
type Rate struct {
	RecvBytesPerSec  int64
	SentBytesPerSec  int64
	TotalBytesPerSec int64
}

type Session struct {
	Idx     int
	ID      string
	User    string
	Realm   string
	Origin  string
	Age     int
	Expires int

	ClientProto string
	RelayProto  string
	ClientAddr  string
	ServerAddr  string
	RelayAddrV4 string
	RelayAddrV6 string

	Fingerprints bool
	Mobile       bool

	TLSVers   string
	TLSCipher string

	BPS int

	Packets Packets
	Rate    Rate

	Peers string
}
