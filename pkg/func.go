package pkg

import (
	"os"

	"github.com/fleimkeipa/kondukto/utils"
	"github.com/go-git/go-git/v5"
)

func ScanFunc(url string) (*git.Repository, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	//delete repository and results.json
	err = os.RemoveAll(cwd + "/tmp/src")
	if err != nil {
		return nil, err
	}

	utils.Info("git clone " + url)

	repo, err := git.PlainClone("./tmp/src", false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		return nil, err
	}

	return repo, nil
}
