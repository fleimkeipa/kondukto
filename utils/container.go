package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/go-git/go-git/v5"
)

func RunContainer(cli *client.Client, repo *git.Repository) error {
	start := time.Now()
	ctx := context.Background()

	reader, err := cli.ImagePull(ctx, "opensorcery/bandit", types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, reader)

	fmt.Println("imagePull", time.Since(start))
	workingDir := "/code"

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// docker run --rm -v ${PWD}:/code opensorcery/bandit -r /code
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:      "opensorcery/bandit",
		Cmd:        []string{"-r", workingDir, "-f", "json", "-o", "result.json"},
		WorkingDir: workingDir,
	}, &container.HostConfig{
		AutoRemove: true,
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Target: workingDir,
				Source: cwd + "/tmp/src",
			},
		},
	}, nil, nil, "")
	if err != nil {
		return err
	}
	fmt.Println("containerCreate", time.Since(start))

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	fmt.Println("containerStart", time.Since(start))

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}
	fmt.Println("containerWait", time.Since(start))

	return nil
}
