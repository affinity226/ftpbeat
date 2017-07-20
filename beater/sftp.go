package beater

import (
	"bufio"
	"fmt"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type stSFTP struct {
	con    *ssh.Client
	client *sftp.Client
}

func (f *stSFTP) Init(bt *Ftpbeat) error {
	var err error
	var auths []ssh.AuthMethod
	auths = append(auths, ssh.Password(bt.password))
	config := ssh.ClientConfig{
		User: bt.username,
		Auth: auths,
	}
	f.con, err = ssh.Dial("tcp", fmt.Sprintf("%s:%s", bt.hostname, bt.port), &config)
	if err != nil {
		logp.Err(fmt.Sprintf("%v", err))
		return err
	}

	return nil
}

func (f *stSFTP) Login(bt *Ftpbeat) error {
	var err error
	f.client, err = sftp.NewClient(f.con, sftp.MaxPacket(1<<15))
	if err != nil {
		logp.Err(fmt.Sprintf("%v", err))
	}
	return err

}

func (f *stSFTP) CheckFiles(bt *Ftpbeat) error {
	var temp []string
	for _, fn := range bt.files {
		if strings.ContainsAny(fn, "* | ?") {
			list, err := f.client.Glob(filepath.Join(bt.remoteDirectory, fn))
			if err == nil {
				for _, fPath := range list {
					_, fName := filepath.Split(fPath)
					temp = append(temp, fName)
				}
			} else {
				logp.Err(fmt.Sprintf("%v", err))
				return err
			}
		} else {
			temp = append(temp, fn)
		}
	}
	bt.files = temp
	logp.Info("Files : %v", bt.files)
	return nil

}

func (f *stSFTP) GenEventForLocalFile(file string, bt *Ftpbeat, b *beat.Beat) error {
	var event common.MapStr
	//r, err := f.client.Open(filepath.Join(bt.currentDirectory, file))
	r, err := os.Open(filepath.Join(bt.currentDirectory, file))
	if err != nil {
		logp.Err(fmt.Sprintf("%v", err))
		return err
	} else {
		scan := bufio.NewScanner(r)

		if err := scan.Err(); err != nil {
			logp.Err(fmt.Sprintf("%v", err))
			r.Close()
			return err
		}
		for scan.Scan() {
			event = common.MapStr{
				"@timestamp": common.Time(time.Now()),
				"type":       bt.connectType,
			}
			event["message"] = scan.Text()
			//b.Events.PublishEvent(event)
			bt.client.PublishEvent(event)
			event = nil
		}
		r.Close()
	}
	return nil

}
func (f *stSFTP) GenEvent(file string, bt *Ftpbeat, b *beat.Beat) error {
	var event common.MapStr
	r, err := f.client.Open(filepath.Join(bt.remoteDirectory, file))
	if err != nil {
		logp.Err(fmt.Sprintf("%v", err))
		return err
	} else {
		scan := bufio.NewScanner(r)

		if err := scan.Err(); err != nil {
			logp.Err(fmt.Sprintf("%v", err))
			r.Close()
			return err
		}
		for scan.Scan() {
			event = common.MapStr{
				"@timestamp": common.Time(time.Now()),
				"type":       bt.connectType,
			}
			event["message"] = scan.Text()
			//b.Events.PublishEvent(event)
			bt.client.PublishEvent(event)
			event = nil
		}
		r.Close()
	}
	return nil

}

func (f *stSFTP) CopyFiles(file string, bt *Ftpbeat) error {
	r, err := f.client.Open(filepath.Join(bt.remoteDirectory, file))
	if err != nil {
		logp.Err(fmt.Sprintf("%v : %s", err, file))
		return err
	} else {
		outf, err := os.Create(filepath.Join(bt.currentDirectory, file))
		if err != nil {
			r.Close()
			logp.Err(fmt.Sprintf("%v : %s", err, file))
			return err
		}
		io.Copy(outf, r)
		outf.Close()
		r.Close()
	}
	return nil

}

func (f *stSFTP) Quit() {
	if f.con != nil {
		f.con.Close()
	}
	if f.client != nil {
		f.client.Close()
	}
}
