package auth

import (
	"context"
	"errors"

	"github.com/smallnest/rpcx/protocol"
)

func Auth(ctx context.Context, req *protocol.Message, token string) error {

	if token == "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IjIxMjMyZjI5N2E1N2E1YTc0Mzg5NGEwZTRhODAxZmMzIiwicGFzc3dvcmQiOiIyMTIzMmYyOTdhNTdhNWE3NDM4OTRhMGU0YTgwMWZjMyIsImV4cCI6MTYyMTQ2MTQ5MSwiaXNzIjoiZ2luLWJsb2cifQ.2AzzE_5rEZGgIC75yI0BdoKt7H_hJ8uQwuZaNCzRhy4" {
		return nil
	}
	return errors.New("invalid token")
}
