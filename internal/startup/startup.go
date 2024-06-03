package startup

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/Kingfish219/PlaNet/internal/interfaces"
	"github.com/Kingfish219/PlaNet/internal/presets"
	"github.com/Kingfish219/PlaNet/internal/publisher"
	"github.com/Kingfish219/PlaNet/internal/repository"
	"github.com/Kingfish219/PlaNet/internal/ui/menu/systray"
)

type Startup struct {
	userInterfaces []interfaces.UserInterface
}

func New() Startup {
	return Startup{
		userInterfaces: []interfaces.UserInterface{},
	}
}

func (startup *Startup) Initialize() error {
	repoFilePath, err := startup.createRepoFilePath()
	if err != nil {
		return err
	}

	dnsRepository := repository.NewDnsRepository(repoFilePath)
	startup.migrateDb(dnsRepository)

	publisher := publisher.Publisher{}

	// console := console.New(dnsRepository)
	// startup.userInterfaces = append(startup.userInterfaces, console)
	// publisher.UISubscribers = append(publisher.UISubscribers, console)

	systray := systray.New(dnsRepository)
	startup.userInterfaces = append(startup.userInterfaces, systray)
	publisher.UISubscribers = append(publisher.UISubscribers, systray)

	startup.userInterfaces = append(startup.userInterfaces, systray)

	return nil
}

func (startup *Startup) Start() error {
	var err error

	for _, userInterface := range startup.userInterfaces {
		err = userInterface.Initialize()
	}

	return err
}

func (startup *Startup) createRepoFilePath() (string, error) {
	tempDirPath, err := os.UserCacheDir()
	if err != nil {
		tempDirPath = os.TempDir()
	}

	planetTempDirPath := filepath.Join(tempDirPath, "PlaNet")
	err = os.MkdirAll(planetTempDirPath, 0644)
	if err != nil {
		return "", err
	}

	repoFilePath := filepath.Join(planetTempDirPath, "config.json")
	_, err = os.Stat(repoFilePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", nil
		}

		_, err = os.Create(repoFilePath)
		if err != nil {
			return "", err
		}
	}

	return repoFilePath, nil
}

func (startup *Startup) migrateDb(repository interfaces.DnsRepository) {
	presetDnsList := presets.GetDnsPresets()
	for _, pre := range presetDnsList {
		repository.ModifyDnsConfigurations(pre)
	}
}
