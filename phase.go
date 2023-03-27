package bpm

import (
	"bpm/log"
	"context"
	"fmt"
	dockerCtx "github.com/docker/distribution/context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"os"
)

type Context struct {
	Phase       Phase
	Cli         *client.Client
	ProcCtx     context.Context
	ContainerId string
}
type Runner interface {
	Run(ctx Context) error
}

func (phase Phase) Run(ctx Context) error {
	// set up the agent, then run all steps in the agent
	if phase.Name != "" {
		log.Infof("Running phase '%s'", phase.Name)
		procCtx := dockerCtx.Background()
		ctx.ProcCtx = procCtx
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		ctx.Cli = cli
		if err != nil {
			return fmt.Errorf("error creating Docker Cli client: %s", err)
		}
		id, err := startContainer(ctx)
		if err != nil {
			return err
		}
		ctx.ContainerId = id
		for _, step := range phase.Steps {
			err := step.Run(ctx)
			if err != nil {
				return err
			}
		}
		if err := stopContainer(ctx); err != nil {
			fmt.Printf("error while stopping container '%s': %v\n", ctx.ContainerId, err)
		}
	}
	return nil
}

func stopContainer(ctx Context) error {
	log.Infof("stopping container '%s'", ctx.ContainerId)
	return ctx.Cli.ContainerStop(ctx.ProcCtx, ctx.ContainerId, container.StopOptions{})
}

func startContainer(ctx Context) (string, error) {
	resp, err := ctx.Cli.ContainerCreate(ctx.ProcCtx, &container.Config{
		Image: string(ctx.Phase.Agent),
		Tty:   false,
	}, nil, nil, nil, ctx.Phase.Id.String())
	if err != nil {
		return "", err
	}
	if err := ctx.Cli.ContainerStart(ctx.ProcCtx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}
	log.Infof("Started container '%s'", resp.ID)
	return resp.ID, nil
}

func (s Step) Run(ctx Context) error {
	log.Infof("Running command Step '%s' using container '%s'", s.Name, ctx.ContainerId)
	for _, c := range s.Commands {
		resp, err := ctx.Cli.ContainerExecCreate(ctx.ProcCtx, ctx.ContainerId, types.ExecConfig{
			Env:          formatEnv(s.Env),
			Cmd:          c,
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Detach:       false,
			Tty:          true,
		})
		if err != nil {
			return fmt.Errorf("could no exec create: %s", err)
		}
		attach, err := ctx.Cli.ContainerExecAttach(ctx.ProcCtx, resp.ID, types.ExecStartCheck{
			Detach: false,
			Tty:    true,
		})
		if err != nil {
			return fmt.Errorf("could no exec start: %s", err)
		}
		_, err = io.Copy(os.Stdout, attach.Reader)
		if err != nil {
			return err
		}
		attach.Close()
	}
	return nil
}

func formatEnv(env map[string]string) []string {
	formattedEnv := make([]string, len(env))
	for k, v := range env {
		if env[k] == "" {
			formattedEnv = append(formattedEnv, v)
		} else {
			formattedEnv = append(formattedEnv, fmt.Sprintf("%s=%s", k, v))
		}
	}
	return formattedEnv
}

func (p Parallel) Run(ctx Context) error {
	log.Infof("Running Parallel Step '%s'", p.Name)
	return nil
}

func (s Script) Run(ctx Context) error {
	log.Infof("Running Script Step '%s'", s.Name)
	return nil
}
