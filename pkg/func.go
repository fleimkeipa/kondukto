package pkg

import (
	"fmt"
	"os"

	"github.com/fleimkeipa/kondukto/utils"
	"github.com/go-git/go-git/v5"
	"github.com/google/uuid"
)

func ScanFunc(url string) (string, error) {
	utils.Info("git clone " + url)

	_, err := git.PlainClone("/deneme/src", false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	fmt.Println("cp0")
	if err != nil {
		return "", err
	}
	fmt.Println("cp1")
	uuid, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	fmt.Println("cp2")

	return uuid.String(), nil
}
