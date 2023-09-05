package middleware

import (
	"bufio"
	"bytes"
	"context"
	"github.com/google/uuid"
	"github.com/jollyboss123/scs/v2"
	"net"
	"net/http"
	"time"
)

type key string

const (
	KeyID      key = "id"
	KeySession key = "session"
)

func Authenticate(m *scs.SessionManager) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token := m.Token(ctx)
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			_, found, err := m.CtxStore.FindCtx(ctx, token)
			if err != nil {
				return
			}

			if !found {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func LoadAndSave(s *scs.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var token string
			cookie, err := r.Cookie(s.Cookie.Name)
			if err == nil {
				token = cookie.Value
			}

			ctx, err := s.Load(r.Context(), token)
			if err != nil {
				s.ErrorFunc(w, r, err)
				return
			}

			sr := r.WithContext(ctx)
			bw := &bufferedResponseWriter{ResponseWriter: w}
			next.ServeHTTP(bw, sr)

			if sr.MultipartForm != nil {
				_ = sr.MultipartForm.RemoveAll()
			}

			var userID any
			userID, ok := s.Get(ctx, string(KeyID)).(uuid.UUID)
			if !ok {
				userID = nil
			}
			ctx = context.WithValue(ctx, KeyID, userID)

			switch s.Status(ctx) {
			case scs.Modified:
				token, expiry, err := s.Commit(ctx)
				if err != nil {
					s.ErrorFunc(w, r, err)
					return
				}

				s.WriteSessionCookie(ctx, w, token, expiry)
			case scs.Destroyed:
				s.WriteSessionCookie(ctx, w, "", time.Time{})
			}

			w.Header().Add("Vary", "Cookie")

			if bw.code != 0 {
				w.WriteHeader(bw.code)
			}
			_, _ = w.Write(bw.buf.Bytes())
		})
	}
}

type bufferedResponseWriter struct {
	http.ResponseWriter
	buf         bytes.Buffer
	code        int
	wroteHeader bool
}

func (bw *bufferedResponseWriter) Write(b []byte) (int, error) {
	return bw.buf.Write(b)
}

func (bw *bufferedResponseWriter) WriteHeader(code int) {
	if !bw.wroteHeader {
		bw.code = code
		bw.wroteHeader = true
	}
}

func (bw *bufferedResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj := bw.ResponseWriter.(http.Hijacker)
	return hj.Hijack()
}

func (bw *bufferedResponseWriter) Push(target string, opts *http.PushOptions) error {
	if pusher, ok := bw.ResponseWriter.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}
	return http.ErrNotSupported
}
