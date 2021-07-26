package goshell

import (
	"bytes"
	"testing"
)

func TestNewSecureShell(t *testing.T) {
	var buf bytes.Buffer
	tl := []struct {
		host     string
		port     int
		user     string
		password string

		want string
	}{
		{"127.0.0.1", 22, "root", "********", "retry times was exceeded 10"},
	}
	for _, v := range tl {
		_, err := NewSecureShell(&buf, v.host, v.user, v.password, v.port)
		if err != nil {
			if v.want == err.Error() {
				continue
			}
			t.Fatalf("can not connect to host %s:%d : %v", v.host, v.port, err)
		}
	}
}
