package beater

import (
	"fmt"
	"time"

	"github.com/affinity226/ftpbeat/config"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"
)

// Ftpbeat is a struct to hold the beat config & info
type Ftpbeat struct {
	beatConfig       *config.Config
	done             chan struct{}
	period           time.Duration
	connectType      string
	hostname         string
	port             string
	username         string
	password         string
	remoteDirectory  string
	currentDirectory string
	executeType      string
	files            []string
	//runner           interface{}
	runner integratedFunc
	client publisher.Client
}

const (
	// default values
	defaultPeriod          = "10s"
	defaultHostname        = "127.0.0.1"
	defaultPort            = "21"
	defaultConnectType     = "ftp"
	defaultUsername        = "ftpbeat_user"
	defaultPassword        = "ftpbeat_pass"
	defaultRemoteDirectory = "~/"
	defaultCurrDirectory   = "./"
	defaultExecuteType     = "get"

	// supported Connect types
	ctFTP  = "ftp"
	ctSFTP = "sftp"

	etRead = "read"
	etGet  = "get"
)

// New Creates beater
/*func New() *Ftpbeat {
	return &Ftpbeat{
		done: make(chan struct{}),
	}
}*/
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	bt := &Ftpbeat{
		done: make(chan struct{}),
	}
	err := cfgfile.Read(&bt.beatConfig, "")
	if err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}
	err = bt.Setup(b)
	if err != nil {
		return nil, fmt.Errorf("Error setting config file: %v", err)
	}
	bt.PrintConfig()
	return bt, err
}

type integratedFunc interface {
	Init(bt *Ftpbeat) error
	Login(bt *Ftpbeat) error
	CheckFiles(bt *Ftpbeat) error
	GenEvent(file string, bt *Ftpbeat, b *beat.Beat) error
	GenEventForLocalFile(file string, bt *Ftpbeat, b *beat.Beat) error
	CopyFiles(file string, bt *Ftpbeat) error
	Quit()
}

///*** Beater interface methods ***///

// Config is a function to read config file
func (bt *Ftpbeat) Config(b *beat.Beat) error {

	// Load beater beatConfig
	err := cfgfile.Read(&bt.beatConfig, "")
	if err != nil {
		return fmt.Errorf("Error reading config file: %v", err)
	}
	bt.PrintConfig()

	return nil
}

func (bt *Ftpbeat) PrintConfig() {
	logp.Info("===========================================================")
	logp.Info("Period           : %v", bt.beatConfig.Ftpbeat.Period)
	logp.Info("ConnectType      : %v", bt.beatConfig.Ftpbeat.ConnectType)
	logp.Info("Hostname         : %v", bt.beatConfig.Ftpbeat.Hostname)
	logp.Info("Port             : %v", bt.beatConfig.Ftpbeat.Port)
	logp.Info("Username         : %v", bt.beatConfig.Ftpbeat.Username)
	logp.Info("RemoteDirectory  : %v", bt.beatConfig.Ftpbeat.RemoteDirectory)
	logp.Info("CurrentDirectory : %v", bt.beatConfig.Ftpbeat.CurrentDirectory)
	logp.Info("Files            : %v", bt.beatConfig.Ftpbeat.Files)
	logp.Info("ExecuteType      : %v", bt.beatConfig.Ftpbeat.ExecuteType)
	logp.Info("===========================================================")
}

// Setup is a function to setup all beat config & info into the beat struct
func (bt *Ftpbeat) Setup(b *beat.Beat) error {

	// Config errors handling
	switch bt.beatConfig.Ftpbeat.ConnectType {
	case ctFTP, ctSFTP:
		break
	default:
		err := fmt.Errorf("Unknown [%s] Connection type, supported types: `ftp`, `sftp`", bt.beatConfig.Ftpbeat.ConnectType)
		return err
	}

	if len(bt.beatConfig.Ftpbeat.Files) < 1 {
		err := fmt.Errorf("There are no files to get")
		return err
	}

	// Setting defaults for missing config
	if bt.beatConfig.Ftpbeat.Period == "" {
		logp.Info("Period not selected, proceeding with '%v' as default", defaultPeriod)
		bt.beatConfig.Ftpbeat.Period = defaultPeriod
	}

	if bt.beatConfig.Ftpbeat.ConnectType == "" {
		logp.Info("Connection Type not selected, proceeding with '%v' as default", defaultConnectType)
		bt.beatConfig.Ftpbeat.ConnectType = defaultConnectType
	}

	if bt.beatConfig.Ftpbeat.Hostname == "" {
		logp.Info("Hostname not selected, proceeding with '%v' as default", defaultHostname)
		bt.beatConfig.Ftpbeat.Hostname = defaultHostname
	}

	if bt.beatConfig.Ftpbeat.Port == "" {
		logp.Info("Port not selected, proceeding with '%v' as default", bt.beatConfig.Ftpbeat.Port)
		bt.beatConfig.Ftpbeat.Port = defaultPort
	}

	if bt.beatConfig.Ftpbeat.Username == "" {
		logp.Info("Username not selected, proceeding with '%v' as default", defaultUsername)
		bt.beatConfig.Ftpbeat.Username = defaultUsername
	}

	if bt.beatConfig.Ftpbeat.Password == "" {
		logp.Info("Password not selected, proceeding with default password")
		bt.beatConfig.Ftpbeat.Password = defaultPassword
	}

	if bt.beatConfig.Ftpbeat.RemoteDirectory == "" {
		logp.Info("Remote Directory not selected, proceeding with '%v' as default", defaultRemoteDirectory)
		bt.beatConfig.Ftpbeat.RemoteDirectory = defaultRemoteDirectory
	}

	if bt.beatConfig.Ftpbeat.CurrentDirectory == "" {
		logp.Info("Current Directory not selected, proceeding with '%v' as default", defaultCurrDirectory)
		bt.beatConfig.Ftpbeat.CurrentDirectory = defaultCurrDirectory
	}

	if bt.beatConfig.Ftpbeat.ExecuteType == "" {
		logp.Info("Execute Type not selected, proceeding with '%v' as default", defaultExecuteType)
		bt.beatConfig.Ftpbeat.ExecuteType = defaultExecuteType
	}

	// Config errors handling
	switch bt.beatConfig.Ftpbeat.ExecuteType {
	case etGet, etRead:
		break
	default:
		err := fmt.Errorf("Unknown [%s] Execute type, supported types: `read`, `get`", bt.beatConfig.Ftpbeat.ConnectType)
		return err
	}

	// Parse the Period string
	var durationParseError error
	bt.period, durationParseError = time.ParseDuration(bt.beatConfig.Ftpbeat.Period)
	if durationParseError != nil {
		return durationParseError
	}

	// Handle password decryption and save in the bt
	if bt.beatConfig.Ftpbeat.Password != "" {
	}

	// Save config values to the bt
	bt.connectType = bt.beatConfig.Ftpbeat.ConnectType
	bt.hostname = bt.beatConfig.Ftpbeat.Hostname
	bt.port = bt.beatConfig.Ftpbeat.Port
	bt.username = bt.beatConfig.Ftpbeat.Username
	bt.password = bt.beatConfig.Ftpbeat.Password

	bt.files = bt.beatConfig.Ftpbeat.Files
	bt.remoteDirectory = bt.beatConfig.Ftpbeat.RemoteDirectory
	bt.currentDirectory = bt.beatConfig.Ftpbeat.CurrentDirectory
	bt.executeType = bt.beatConfig.Ftpbeat.ExecuteType

	logp.Info("Total # of files to get : %d", len(bt.files))
	for index, file := range bt.files {
		logp.Info("Read #%d : %s", index+1, file)
	}

	switch bt.connectType {
	case ctFTP:
		bt.runner = new(stFTP)
		break
	case ctSFTP:
		bt.runner = new(stSFTP)
		break
	}

	return nil
}

// Run is a functions that runs the beat
func (bt *Ftpbeat) Run(b *beat.Beat) error {
	logp.Info("ftpbeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()

	ticker := time.NewTicker(bt.period)
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		err := bt.beat(b)
		if err != nil {
			return err
		}
	}
}

// Cleanup is a function that does nothing on this beat :)
func (bt *Ftpbeat) Cleanup(b *beat.Beat) error {
	return nil
}

// Stop is a function that runs once the beat is stopped
func (bt *Ftpbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

// beat is a function that iterate over the query array, generate and publish events
func (bt *Ftpbeat) beat(b *beat.Beat) error {
	logp.Info("Run Beat Periodically")

	err := bt.runner.Init(bt)
	if err != nil {
		return err
	}
	defer bt.runner.Quit()

	err = bt.runner.Login(bt)
	if err != nil {
		return err
	}

	err = bt.runner.CheckFiles(bt)
	if err != nil {
		return err
	}
	for _, file := range bt.files {
		if bt.executeType == etRead {
			bt.runner.GenEvent(file, bt, b)
		} else {
			bt.runner.CopyFiles(file, bt)
			bt.runner.GenEventForLocalFile(file, bt, b)
		}
	}
	// Great success!
	return nil
}
