package httpserver_test

import (
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/ozaitsev92/tododdd/pkg/httpserver"
)

func freeAddr(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	return l.Addr().String()
}

func waitListening(t *testing.T, addr string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if c, err := net.Dial("tcp", addr); err == nil {
			c.Close()
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("server at %s never started", addr)
}

func TestServer_NewDoesNotListen(t *testing.T) {
	t.Parallel()
	addr := freeAddr(t)
	_ = httpserver.New(http.NotFoundHandler(), httpserver.Addr(addr))

	if c, err := net.DialTimeout("tcp", addr, 100*time.Millisecond); err == nil {
		c.Close()
		t.Fatal("server listens before Start")
	}
}

func TestServer_Shutdown_DrainsInFlightRequest(t *testing.T) {
	t.Parallel()
	addr := freeAddr(t)
	started := make(chan struct{})
	release := make(chan struct{})

	s := httpserver.New(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		close(started)
		<-release
		_, _ = io.WriteString(w, "done")
	}), httpserver.Addr(addr), httpserver.ShutdownTimeout(5*time.Second))
	s.Start()
	waitListening(t, addr)

	type result struct {
		body string
		err  error
	}
	resCh := make(chan result, 1)
	go func() {
		resp, err := http.Get("http://" + addr)
		if err != nil {
			resCh <- result{err: err}
			return
		}
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		resCh <- result{string(b), err}
	}()

	<-started
	shutdownDone := make(chan error, 1)
	go func() { shutdownDone <- s.Shutdown() }()

	select {
	case err := <-shutdownDone:
		t.Fatalf("Shutdown returned (%v) before the in-flight request finished", err)
	case <-time.After(100 * time.Millisecond):
	}

	close(release)
	if res := <-resCh; res.err != nil || res.body != "done" {
		t.Fatalf("in-flight request: %q, %v", res.body, res.err)
	}
	if err := <-shutdownDone; err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
	if _, open := <-s.Notify(); open {
		t.Fatal("Notify delivered an error after graceful shutdown")
	}
}

func TestServer_ListenError_IsNotified(t *testing.T) {
	t.Parallel()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	s := httpserver.New(http.NotFoundHandler(), httpserver.Addr(l.Addr().String()))
	s.Start()

	select {
	case err := <-s.Notify():
		if err == nil {
			t.Fatal("want listen error, got nil")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("no error for an address already in use")
	}
}
