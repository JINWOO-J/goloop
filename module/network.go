package module

import "fmt"

type NetworkManager interface {
	Start() error
	Term()

	GetPeers() []PeerID

	RegisterReactor(name string, pi ProtocolInfo, reactor Reactor, piList []ProtocolInfo, priority uint8) (ProtocolHandler, error)
	RegisterReactorForStreams(name string, pi ProtocolInfo, reactor Reactor, piList []ProtocolInfo, priority uint8) (ProtocolHandler, error)
	UnregisterReactor(reactor Reactor) error

	SetRole(version int64, role Role, peers ...PeerID)
	GetPeersByRole(role Role) []PeerID
	AddRole(role Role, peers ...PeerID)
	RemoveRole(role Role, peers ...PeerID)
	HasRole(role Role, id PeerID) bool
	Roles(id PeerID) []Role

	SetTrustSeeds(seeds string)
	SetInitialRoles(roles ...Role)
}

type Reactor interface {
	//case broadcast and multicast, if return (true,nil) then rebroadcast
	OnReceive(pi ProtocolInfo, b []byte, id PeerID) (bool, error)
	OnFailure(err error, pi ProtocolInfo, b []byte)
	OnJoin(id PeerID)
	OnLeave(id PeerID)
}

type ProtocolHandler interface {
	Broadcast(pi ProtocolInfo, b []byte, bt BroadcastType) error
	Multicast(pi ProtocolInfo, b []byte, role Role) error
	Unicast(pi ProtocolInfo, b []byte, id PeerID) error
}

type BroadcastType byte
type Role string

const (
	BROADCAST_ALL      BroadcastType = 0
	BROADCAST_NEIGHBOR BroadcastType = 1
	ROLE_VALIDATOR     Role          = "validator"
	ROLE_SEED          Role          = "seed"
	ROLE_NORMAL        Role          = "normal"
)

type PeerID interface {
	Bytes() []byte
	Equal(PeerID) bool
	String() string
}

const (
	ProtoUnknown ProtocolInfo = iota << 8
	ProtoStateSync
	ProtoTransaction
	ProtoConsensus
	ProtoFastSync
	ProtoConsensusSync
)

type ProtocolInfo uint16

func NewProtocolInfo(id byte, version byte) ProtocolInfo {
	return ProtocolInfo(int(id)<<8 | int(version))
}
func (pi ProtocolInfo) ID() byte {
	return byte(pi >> 8)
}
func (pi ProtocolInfo) Version() byte {
	return byte(pi)
}
func (pi ProtocolInfo) String() string {
	return fmt.Sprintf("{%#04x}", pi.Uint16())
}
func (pi ProtocolInfo) Uint16() uint16 {
	return uint16(pi)
}

type NetworkTransport interface {
	Listen() error
	Close() error
	Dial(address string, channel string) error
	PeerID() PeerID
	Address() string
	SetListenAddress(address string) error
	GetListenAddress() string
	SetSecureSuites(channel string, secureSuites string) error
	GetSecureSuites(channel string) string
	SetSecureAeads(channel string, secureAeads string) error
	GetSecureAeads(channel string) string
}

//TODO remove interface and implement network.IsTemporaryError(error) bool
type NetworkError interface {
	error
	Temporary() bool // Is the error temporary?
}
