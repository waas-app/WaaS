package auth

import (
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/volatiletech/authboss/v3"
	"github.com/waas-app/WaaS/config"
)

type SessionState struct {
	session *sessions.Session
}

// Get a key from the session
func (s SessionState) Get(key string) (string, bool) {
	str, ok := s.session.Values[key]
	if !ok {
		return "", false
	}
	value := str.(string)

	return value, ok
}

type SessionStorer struct {
	Name string

	sessions.Store
}

// NewSessionStorer creates a new gorilla sessions.CookieStore
// and uses that for underlying storage.
//
// The sessionName is the name that will be used for the session cookie
// on the client (one session has many values).
//
// Also it takes pairs of keys (hmac auth, encryption), and if doing
// key rotation multiple of these pairs. The second key of the pair
// should be set to nil if encryption isn't desired.
//
// Authentication keys should be 32 or 64 bytes.
// Encryption keys should be 16, 24, or 32 bytes for AES-128, AES-192, and AES-256
// respectively.
//
// These docs are prone to doc-rot since they're copied from the gorilla
// session store documentation.
func NewSessionStorer(sessionName string, keypairs ...[]byte) SessionStorer {

	cs := sessions.NewCookieStore(keypairs...)
	cs.Options.Domain = config.Spec.CookieDomain
	return SessionStorer{
		Name:  sessionName,
		Store: cs,
	}
}

// NewSessionStorerFromExisting takes a store object that's already
// configured and uses it directly. This can be anything that satisfies
// the interface.
//
// sessionName is the name of the cookie/file/whatever on the client or on
// the filesystem etc.
func NewSessionStorerFromExisting(sessionName string, store sessions.Store) SessionStorer {
	return SessionStorer{
		Name:  sessionName,
		Store: store,
	}
}

// ReadState loads the session from the request context
func (s SessionStorer) ReadState(r *http.Request) (authboss.ClientState, error) {
	// Note that implementers of Get in gorilla all return a new session
	session, err := s.Store.Get(r, s.Name)
	if err != nil {
		e, ok := err.(securecookie.Error)
		if ok && !e.IsDecode() {
			// We ignore decoding errors, but nothing else
			return nil, err
		}

		// Get returning a new session even when there's an error is a bit
		// more up in the air, so we force the new session here if we've
		// previously encountered an error.
		session, err = s.Store.New(r, s.Name)
		if err != nil {
			return nil, err
		}
	}

	cs := &SessionState{
		session: session,
	}

	return cs, nil
}

// WriteState to the responsewriter
func (s SessionStorer) WriteState(w http.ResponseWriter, state authboss.ClientState, ev []authboss.ClientStateEvent) error {
	// This should never be nil (despite what authboss.ClientStateReadWriter
	// interface says) because all Get methods return a new session in gorilla.
	// In cases where Get returns an error, we ensure we create a new session
	ses := state.(*SessionState)

	for _, ev := range ev {
		switch ev.Kind {
		case authboss.ClientStateEventPut:
			ses.session.Values[ev.Key] = ev.Value
		case authboss.ClientStateEventDel:
			delete(ses.session.Values, ev.Key)
		}
	}

	return s.Store.Save(nil, w, ses.session)
}
