package pkg

import (
	"os"

	"github.com/fleimkeipa/kondukto/utils"
	"github.com/go-git/go-git/v5"
)

func ScanFunc(url string) (*git.Repository, error) {
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
