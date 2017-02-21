// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

type Config struct {
	Ftpbeat FtpbeatConfig
}

type FtpbeatConfig struct {
	Period           string   `yaml:"period"`
	ConnectType      string   `yaml:"connecttype"`
	Hostname         string   `yaml:"hostname"`
	Port             string   `yaml:"port"`
	Username         string   `yaml:"username"`
	Password         string   `yaml:"password"`
	RemoteDirectory  string   `yaml:"remotedirectory"`
	CurrentDirectory string   `yaml:"currentdirectory"`
	Files            []string `yaml:"files"`
	ExecuteType      string   `yaml:"executetype"`
}
