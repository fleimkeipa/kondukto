package utils

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/stdcopy"
)

var dockerRegistryUserID = "1149"

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

func ImageBuild(cli *client.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	tar, err := archive.TarWithOptions("../../", &archive.TarOptions{})
	if err != nil {
		return err
	}

	opts := types.ImageBuildOptions{
		Tags:       []string{dockerRegistryUserID + "bandit"},
		NoCache:    true,
		Remove:     true,
		Dockerfile: "Dockerfile",
		Outputs:    []types.ImageBuildOutput{},
	}

	res, err := cli.ImageBuild(ctx, tar, opts)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	err = print(res.Body)
	if err != nil {
		return err
	}

	return nil
}

func print(rd io.Reader) error {
	var lastLine string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lastLine = scanner.Text()
		fmt.Println(scanner.Text())
	}

	errLine := &ErrorLine{}
	json.Unmarshal([]byte(lastLine), errLine)
	if errLine.Error != "" {
		return errors.New(errLine.Error)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func RunContainer(cli *client.Client) error {
	ctx := context.Background()

	// reader, err := cli.ImagePull(ctx, "opensorcery/bandit", types.ImagePullOptions{})
	// if err != nil {
	// 	return err
	// }
	// io.Copy(os.Stdout, reader)

	// docker run --rm -v ${PWD}:/code opensorcery/bandit -r /code
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		User:       dockerRegistryUserID,
		WorkingDir: "/code",
		// Cmd:        []string{"bandit", "-r", "/code"},
		Image:   "1149bandit",
		Volumes: map[string]struct{}{},
	}, nil, nil, nil, "")
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return err
	}
	// Demultiplex stdout
	// from the container logs
	stdoutput := new(bytes.Buffer)
	stdcopy.StdCopy(stdoutput, nil, out)
	if err != nil {
		panic(err)
	}
	stdoutput.ReadFrom(out)
	out1 := stdoutput.String()
	b, err := json.Marshal(out1)
	if err != nil {
		fmt.Println("marshall error", err)
		return err
	}
	ioutil.WriteFile("../result.json", b, 0777)
	fmt.Println(out1)

	// if err := cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{}); err != nil {
	// 	return err
	// }

	// stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	return nil
}
