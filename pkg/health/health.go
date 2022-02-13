package health

import (
	"context"
	"fmt"
)

type Health string

type HealthArgs struct {
	BasePath string
}

type HealthReply struct {
	Status uint
	Error  string
}

func (h *Health) Check(ctx context.Context, args *HealthArgs, reply *HealthReply) error {
	reply.Status = 1
	return nil
}

func (h *Health) Reload(ctx context.Context, args *HealthArgs, reply *HealthReply) error {
	fmt.Println("test")
	return nil
}
