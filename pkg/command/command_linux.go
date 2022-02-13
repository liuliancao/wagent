package command

import (
	"bytes"
	"context"
	"log"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
	"time"
)

type CmdArgs struct {
	RunAs    string
	Cmd      string
	Async    bool
	Timeout  int
	Callback string
}

type CmdReply struct {
	Result string
	Error  string
	Code   int
}

type CmdExecutor string

//SetCmdUser: set cmd run as
func (c *CmdExecutor) SetCmdUser(cmd *exec.Cmd, username string) error {
	sysuser, err := user.Lookup(username)
	if err != nil {
		return err
	}

	uid, err := strconv.Atoi(sysuser.Uid)
	gid, err := strconv.Atoi(sysuser.Gid)

	cmd.SysProcAttr.Credential = &syscall.Credential{
		Uid: uint32(uid),
		Gid: uint32(gid),
	}
	return nil
}

//ExecuteCommand excutes command real
func (c *CmdExecutor) ExecuteCommand(ctx context.Context, cargs CmdArgs, reply *CmdReply) (code int, result, erro string) {
	var cmd *exec.Cmd
	if cargs.Async {
		cmd = exec.Command("/bin/bash", "-c", cargs.Cmd)
	} else {
		timeout := time.Duration(cargs.Timeout)
		ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
		cmd = exec.CommandContext(ctx, "/bin/bash", "-c", cargs.Cmd)
		defer cancel()
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	err := c.SetCmdUser(cmd, cargs.RunAs)
	if err != nil {
		log.Printf("set user command error: %v", err)
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if cargs.Async {
		go func() {
			if err := cmd.Start(); err != nil {
				log.Fatalf("async executes error :%v", err)
			}
			cmd.Wait()
			log.Println("call back", cargs.Callback)
		}()
	} else {
		if err := cmd.Run(); err != nil {
			log.Println(err)
			return 20005, stdout.String(), err.Error() + stderr.String()
		}
		return 200, stdout.String(), stderr.String()
		/*che := make(chan error, 1)
		defer close(che)
		go func() {
			if err := cmd.Run(); err != nil {
				che <- err
				return
			}
			che <- nil
		}()
		timeout := time.Duration(cargs.Timeout)
		select {
		case _ = <-che:
			return stdout.String(), stderr.String()
		case <-time.After(timeout * time.Second):
			fmt.Println("will kill")
			_ = cmd.Process.Kill()
			return "", "sync executes with timeout"
		}*/
	}

	return 200, "async triggered", ""
}

func (c *CmdExecutor) Execute(ctx context.Context, cargs CmdArgs, reply *CmdReply) error {
	// match the cmd and split to parts
	//stdout, err := c.ExecuteCommand(cargs)
	reply.Code, reply.Result, reply.Error = c.ExecuteCommand(ctx, cargs, reply)
	return nil
}
