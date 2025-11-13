package helpers

import "os"

func GetConfigDir() string {
	userConfigDir, _ := os.UserConfigDir();
  cfgDir := userConfigDir + "/TipAggregator";
	os.MkdirAll(cfgDir, 0755)
	return cfgDir
}
