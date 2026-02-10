package main

import (
	"os"

	"github.com/ricky97gr/lpcenter/server/database"
	"github.com/ricky97gr/lpcenter/server/router"
	"github.com/ricky97gr/lpcenter/server/utils"
)

func main() {
	opt := utils.Options{
		FileName:   "",
		Level:      "info",
		ModuleName: "lpcenter",
		W:          os.Stdout,
	}
	utils.Logger = utils.New(opt)

	utils.Logger.Infof("lpcenter Started\n")

	if err := utils.InitRSA(); err != nil {
		utils.Logger.Errorf("Failed to initialize RSA keys: %v", err)
		os.Exit(1)
	}

	if err := database.Start(); err != nil {
		utils.Logger.Errorf("Failed to initialize database: %v", err)
		os.Exit(1)
	}

	go router.Start()
	go router.StartDownloadServer()

	select {}
}
