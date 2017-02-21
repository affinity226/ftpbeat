package beater

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/affinity226/ftpbeat/config"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	// ftp module
	"github.com/jlaffaye/ftp"
)

// Ftpbeat is a struct to hold the beat config & info
type Ftpbeat struct {
	beatConfig      *config.Config
	done            chan struct{}
	period          time.Duration
	connnectType    string
	hostname        string
	port            string
	username        string
	password        string
	remoteDirectory string
	executeType     string
	files           []string
}

const (
	// default values
	defaultPeriod      = "10s"
	defaultHostname    = "127.0.0.1"
	defaultPort        = "21"
	defaultConnectType = "ftp"
	defaultUsername    = "ftpbeat_user"
	defaultPassword    = "ftpbeat_pass"
	defaultDirectory   = "~/"
	defaultExecuteType = "get"

	// supported Connect types
	ctFTP  = "ftp"
	ctSFTP = "sftp"

	etRead = "read"
	etGet  = "get"
)

// New Creates beater
func New() *Ftpbeat {
	return &Ftpbeat{
		done: make(chan struct{}),
	}
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
	fmt.Println("===========================================================")
	fmt.Printf("Period          : %v\n", bt.beatConfig.Ftpbeat.Period)
	fmt.Printf("ConnectType     : %v\n", bt.beatConfig.Ftpbeat.ConnectType)
	fmt.Printf("Hostname        : %v\n", bt.beatConfig.Ftpbeat.Hostname)
	fmt.Printf("Port            : %v\n", bt.beatConfig.Ftpbeat.Port)
	fmt.Printf("Username        : %v\n", bt.beatConfig.Ftpbeat.Username)
	fmt.Printf("RemoteDirectory : %v\n", bt.beatConfig.Ftpbeat.RemoteDirectory)
	fmt.Printf("Files           : %v\n", bt.beatConfig.Ftpbeat.Files)
	fmt.Printf("ExecuteType     : %v\n", bt.beatConfig.Ftpbeat.ExecuteType)
	fmt.Println("===========================================================")
}

// Setup is a function to setup all beat config & info into the beat struct
func (bt *Ftpbeat) Setup(b *beat.Beat) error {

	// Config errors handling
	switch bt.beatConfig.Ftpbeat.ConnectType {
	case ctFTP:
		break
	default:
		err := fmt.Errorf("Unknown [%s] Connection type, supported types: `ftp`, ``", bt.beatConfig.Ftpbeat.ConnectType)
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
		logp.Info("Remote Directory not selected, proceeding with '%v' as default", defaultDirectory)
		bt.beatConfig.Ftpbeat.RemoteDirectory = defaultDirectory
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
	bt.connnectType = bt.beatConfig.Ftpbeat.ConnectType
	bt.hostname = bt.beatConfig.Ftpbeat.Hostname
	bt.port = bt.beatConfig.Ftpbeat.Port
	bt.username = bt.beatConfig.Ftpbeat.Username
	bt.password = bt.beatConfig.Ftpbeat.Password

	bt.files = bt.beatConfig.Ftpbeat.Files
	bt.remoteDirectory = bt.beatConfig.Ftpbeat.RemoteDirectory
	bt.executeType = bt.beatConfig.Ftpbeat.ExecuteType

	logp.Info("Total # of files to get : %d", len(bt.files))
	for index, file := range bt.files {
		logp.Info("Read #%d : %s", index+1, file)
	}

	return nil
}

// Run is a functions that runs the beat
func (bt *Ftpbeat) Run(b *beat.Beat) error {
	logp.Info("ftpbeat is running! Hit CTRL-C to stop it.")

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
	close(bt.done)
}

///*** sqlbeat methods ***///

// beat is a function that iterate over the query array, generate and publish events
func (bt *Ftpbeat) beat(b *beat.Beat) error {
	con, err := ftp.DialTimeout(fmt.Sprintf("%s:%s", bt.hostname, bt.port), 5*time.Second)
	if err != nil {
		fmt.Println("What")
		return err
	}
	defer con.Quit()
	err = con.Login(bt.username, bt.password)
	if err != nil {
		return err
	}
	err = con.ChangeDir(bt.remoteDirectory)
	if err != nil {
		return err
	}

	if bt.executeType == "read" {
	LoopReadFiles:
		for _, file := range bt.files {
			var event common.MapStr
			event = common.MapStr{
				"@timestamp": common.Time(time.Now()),
				"type":       bt.connnectType,
			}
			r, err := con.Retr(file)
			if err != nil {
				continue LoopReadFiles
			} else {
				scan := bufio.NewScanner(r)

				if err := scan.Err(); err != nil {
					return err
					continue LoopReadFiles
				}
				for scan.Scan() {
					event["message"] = scan.Text()
				}
				r.Close()

				logp.Info("event sent")

				b.Events.PublishEvent(event)
				event = nil
			}
		}
	} else { //"get"
	LoopGetFiles:
		for _, file := range bt.files {
			r, err := con.Retr(file)
			if err != nil {
				continue LoopGetFiles
			} else {
				//buf, err := ioutil.ReadAll(r)

				outf, err := os.Create(file)
				if err != nil {
					r.Close()
					continue LoopGetFiles
				}
				io.Copy(outf, r)
			}
		}
	}
	// Great success!
	return nil
}
