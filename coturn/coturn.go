package coturn

import (
	"context"
	"time"

	"github.com/tvi/coturn_exporter/coturn/types"

	"github.com/sirupsen/logrus"
)

type SessionStore struct {
	state map[string]internalSession
	f     Fetcher
}

func NewSessionStore(f Fetcher) (*SessionStore, error) {
	return &SessionStore{
		map[string]internalSession{},
		f,
	}, nil
}

type internalSession struct {
	*types.Session
	firstSeen time.Time
}

func (s *SessionStore) ReloadSessions(ctx context.Context) error {
	t, sessions, err := s.f.Fetch(ctx)
	if err != nil {
		return err
	}
	convertActiveSessions(t, sessions)

	now := time.Now()
	for id, session := range sessions {
		lastSessState, ok := s.state[id]
		if ok {
			// if now.After(lastSessState.firstSeen) {}
			// TODO(tvi): Add session reuse delta.
			lastSessState.Session = &session
			s.state[id] = lastSessState
		} else {
			s.state[id] = internalSession{
				Session:   &session,
				firstSeen: now,
			}
		}
	}

	for id, stoppedSession := range s.state {
		if _, ok := sessions[id]; !ok {
			sessionsAge.Observe(float64(stoppedSession.Age))
			sessionsSentBytes.Observe(float64(stoppedSession.Packets.SentBytes))
			sessionsRecvBytes.Observe(float64(stoppedSession.Packets.RecvBytes))
			sessionsTotalBytes.Observe(float64(stoppedSession.Packets.SentBytes + stoppedSession.Packets.RecvBytes))

			logrus.Debugf("Removing stopped session: %+#v\n", stoppedSession.Session)
			delete(s.state, id)
		}
	}
	return nil
}

func convertActiveSessions(total int, sessions map[string]types.Session) {
	activeSessions.Set(float64(total))

	factiveSessionsAge := float64(0)
	factiveSessionsExpiration := float64(0)

	factiveSessionsSentBytes := float64(0)
	factiveSessionsRecvBytes := float64(0)
	factiveSessionsTotalBytes := float64(0)

	for _, session := range sessions {
		factiveSessionsAge += float64(session.Age)
		factiveSessionsExpiration += float64(session.Expires)

		factiveSessionsSentBytes += float64(session.Packets.SentBytes)
		factiveSessionsRecvBytes += float64(session.Packets.RecvBytes)
		factiveSessionsTotalBytes += float64(session.Packets.SentBytes + session.Packets.RecvBytes)
	}

	lSess := float64(len(sessions))
	activeSessionsAge.Set(factiveSessionsAge / lSess)
	activeSessionsExpiration.Set(factiveSessionsExpiration / lSess)

	activeSessionsSentBytes.Set(factiveSessionsSentBytes / lSess)
	activeSessionsRecvBytes.Set(factiveSessionsRecvBytes / lSess)
	activeSessionsTotalBytes.Set(factiveSessionsTotalBytes / lSess)
}
