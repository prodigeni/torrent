package tracker

import (
	"errors"
	"net"
	"net/url"
)

type AnnounceRequest struct {
	InfoHash   [20]byte
	PeerId     [20]byte
	Downloaded int64
	Left       int64
	Uploaded   int64
	Event      AnnounceEvent
	IPAddress  int32
	Key        int32
	NumWant    int32 // How many peer addresses are desired. -1 for default.
	Port       int16
} // 82 bytes

type AnnounceResponse struct {
	Interval int32 // Minimum seconds the local peer should wait before next announce.
	Leechers int32
	Seeders  int32
	Peers    []Peer
}

type AnnounceEvent int32

type Peer struct {
	IP   net.IP
	Port int
}

const (
	None      AnnounceEvent = iota
	Completed               // The local peer just completed the torrent.
	Started                 // The local peer has just resumed this torrent.
	Stopped                 // The local peer is leaving the swarm.
)

type Client interface {
	// Returns ErrNotConnected if Connect needs to be called.
	Announce(*AnnounceRequest) (AnnounceResponse, error)
	Connect() error
	String() string
	URL() string
}

var (
	ErrNotConnected = errors.New("not connected")
	ErrBadScheme    = errors.New("unknown scheme")

	schemes = make(map[string]func(*url.URL) Client)
)

func RegisterClientScheme(scheme string, newFunc func(*url.URL) Client) {
	schemes[scheme] = newFunc
}

// Returns ErrBadScheme if the tracker scheme isn't recognised.
func New(rawurl string) (cl Client, err error) {
	url_s, err := url.Parse(rawurl)
	if err != nil {
		return
	}
	newFunc, ok := schemes[url_s.Scheme]
	if !ok {
		err = ErrBadScheme
		return
	}
	cl = newFunc(url_s)
	return
}
