package pkg

import (
	"os"

	"github.com/fleimkeipa/kondukto/utils"
	"github.com/go-git/go-git/v5"
)

func ScanFunc(url string) (string, error) {
	utils.Info("git clone " + url)

	_, err := git.PlainClone("./", false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})

	if err != nil {
		return "", err
	}

	return "", nil
}
