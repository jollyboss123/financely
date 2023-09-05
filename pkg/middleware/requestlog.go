package middleware

import (
	"github.com/jollyboss123/finance-tracker/pkg/logger"
	"io"
	"net"
	"net/http"
	"time"
)

func RequestLog(l *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			le := &logEntry{
				UserID:            getUserID(r),
				ReceivedTime:      start,
				RequestMethod:     r.Method,
				RequestURL:        r.URL.String(),
				RequestHeaderSize: headerSize(r.Header),
				UserAgent:         r.UserAgent(),
				Referer:           r.Referer(),
				Proto:             r.Proto,
				RemoteIP:          ipFromHostPort(r.RemoteAddr),
			}

			if addr, ok := r.Context().Value(http.LocalAddrContextKey).(net.Addr); ok {
				le.ServerIP = ipFromHostPort(addr.String())
			}
			r2 := new(http.Request)
			*r2 = *r
			rcc := &readCounterCloser{r: r.Body}
			r2.Body = rcc
			w2 := &responseStats{w: w}

			next.ServeHTTP(w2, r2)

			le.Latency = time.Since(start)
			if rcc.err == nil && rcc.r != nil {
				// If the handler hasn't encountered an error in the Body (like EOF),
				// then consume the rest of the Body to provide an accurate rcc.n.
				io.Copy(io.Discard, rcc)
			}
			le.RequestBodySize = rcc.n
			le.Status = w2.code
			if le.Status == 0 {
				le.Status = http.StatusOK
			}
			le.ResponseHeaderSize, le.ResponseBodySize = w2.size()

			l.Info().
				Str("userID", le.UserID.String()).
				Time("received_time", le.ReceivedTime).
				Str("method", le.RequestMethod).
				Str("url", le.RequestURL).
				Int64("header_size", le.RequestHeaderSize).
				Int64("body_size", le.RequestBodySize).
				Str("agent", le.UserAgent).
				Str("referer", le.Referer).
				Str("proto", le.Proto).
				Str("remote_ip", le.RemoteIP).
				Str("server_ip", le.ServerIP).
				Int("status", le.Status).
				Int64("resp_header_size", le.ResponseHeaderSize).
				Int64("resp_body_size", le.ResponseBodySize).
				Dur("latency", le.Latency).
				Msg("")
		})
	}
}
