package server

import (
	"context"
	"time"

	"github.com/MariusVanDerWijden/ShareMyRPC/raiden"
)

var pollPause = 100 * time.Millisecond

type Server struct {
	node    *raiden.Raiden
	history *raiden.PaymentHistory
	token   string
	peer    string
}

func NewServer(url, token, peer string) *Server {
	return &Server{
		node: raiden.NewRaiden(url),
	}
}

// PaymentReceived returns true if a payment was received
// within the given timeout.
func (s *Server) PaymentReceived(ctx context.Context, maxTimeout time.Duration) bool {
	start := time.Now()
	for {
		if s.pollHistory() {
			return true
		}
		// If we reached the timeout -> cancel
		if time.Since(start) > maxTimeout {
			break
		}
		// If the context is canceled -> cancel
		select {
		case <-ctx.Done():
			break
		default:
			time.Sleep(pollPause)
		}
	}
	return false
}

func (s *Server) pollHistory() bool {
	h, err := s.node.PaymentHistory(s.token, s.peer)
	if err != nil {
		return false
	}
	// No new entries in history
	if len(*h) == len(*s.history) {
		return false
	}
	// new entries in history, update history
	s.history = h
	return true
}
