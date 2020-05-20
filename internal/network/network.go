package network

import (
	"os"

	"github.com/rmanna/ado-pipeline-creator/internal/fileutils"
	logger "github.com/rmanna/ado-pipeline-creator/internal/logger"
)

// SetPort exported
func SetPort() string {
	config := fileutils.ReadYamlConfig("./internal", "config")
	p := os.Getenv("PORT")
	if p != "" {
		logger.Log.InfoArg("Port :" + p + " set by environment")
		return ":" + p
	}
	logger.Log.InfoArg("Port :" + config.Server.Port + " set by configuration file")
	return ":" + config.Server.Port
}
