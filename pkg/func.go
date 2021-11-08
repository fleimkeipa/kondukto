package pkg

import (
	"os"

	"github.com/fleimkeipa/kondukto/utils"
	"github.com/go-git/go-git/v5"
	"github.com/google/uuid"
)

func ScanFunc(url string) (string, error) {
	utils.Info("git clone " + url)

	_, err := git.PlainClone("../../tmp/src", false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		return "", err
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	return uuid.String(), nil
}
