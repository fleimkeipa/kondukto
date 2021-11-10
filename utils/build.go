package utils

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/go-git/go-git/v5"
)

func RunContainer(cli *client.Client, repo *git.Repository) error {
	ctx := context.Background()

	reader, err := cli.ImagePull(ctx, "opensorcery/bandit", types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, reader)

	workingDir := "/code"

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	// fmt.Println("cwd", cwd)

	// docker run --rm -v ${PWD}:/code opensorcery/bandit -r /code
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:      "opensorcery/bandit",
		Cmd:        []string{"-r", workingDir, "-f", "json", "-o", "results.json"},
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

	//delete repository
	err = os.RemoveAll(cwd + "/tmp/src")
	if err != nil {
		return err
	}

	return nil
}

// type ErrorLine struct {
// 	Error       string      `json:"error"`
// 	ErrorDetail ErrorDetail `json:"errorDetail"`
// }

// type ErrorDetail struct {
// 	Message string `json:"message"`
// }

// func ImageBuild(cli *client.Client) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
// 	defer cancel()

// 	tar, err := archive.TarWithOptions("../../", &archive.TarOptions{})
// 	if err != nil {
// 		return err
// 	}

// 	opts := types.ImageBuildOptions{
// 		Tags:       []string{dockerRegistryUserID + "bandit"},
// 		NoCache:    true,
// 		Remove:     true,
// 		Dockerfile: "Dockerfile",
// 		Outputs:    []types.ImageBuildOutput{},
// 	}

// 	res, err := cli.ImageBuild(ctx, tar, opts)
// 	if err != nil {
// 		return err
// 	}

// 	defer res.Body.Close()

// 	err = print(res.Body)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func print(rd io.Reader) error {
// 	var lastLine string

// 	scanner := bufio.NewScanner(rd)
// 	for scanner.Scan() {
// 		lastLine = scanner.Text()
// 		fmt.Println(scanner.Text())
// 	}

// 	errLine := &ErrorLine{}
// 	json.Unmarshal([]byte(lastLine), errLine)
// 	if errLine.Error != "" {
// 		return errors.New(errLine.Error)
// 	}

// 	if err := scanner.Err(); err != nil {
// 		return err
// 	}

// 	return nil
// }
