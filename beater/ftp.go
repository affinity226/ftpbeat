package beater

import (
	"bufio"
	"fmt"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/jlaffaye/ftp"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type stFTP struct {
	con *ftp.ServerConn
}

func (f *stFTP) Init(bt *Ftpbeat) error {
	var err error
	fmt.Println("FTP init")
	f.con, err = ftp.DialTimeout(fmt.Sprintf("%s:%s", bt.hostname, bt.port), 5*time.Second)
	if err != nil {
		logp.Err(fmt.Sprintf("%v", err))
		return err
	}
	return nil
}

func (f *stFTP) Login(bt *Ftpbeat) error {
	fmt.Println("FTP Login")
	err := f.con.Login(bt.username, bt.password)
	if err != nil {
		logp.Err(fmt.Sprintf("%v", err))
	}
	return err

}

func (f *stFTP) CheckFiles(bt *Ftpbeat) error {
	var err error
	fmt.Println("FTP CheckFiles")
	err = f.con.ChangeDir(bt.remoteDirectory)
	if err != nil {
		logp.Err(fmt.Sprintf("%v", err))
		return err
	}

	var temp []string
	for _, fn := range bt.files {
		if strings.ContainsAny(fn, "* | ?") {
			list, err := f.con.NameList(fn)
			if err == nil {
				temp = append(temp, list...)
			} else {
				logp.Err(fmt.Sprintf("%v", err))
			}
		} else {
			temp = append(temp, fn)
		}
	}
	bt.files = temp
	logp.Info("Files : ", bt.files)
	return nil

}

func (f *stFTP) GenEvent(file string, bt *Ftpbeat, b *beat.Beat) error {
	fmt.Println("FTP GenEvent")
	var event common.MapStr
	r, err := f.con.Retr(file)
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
			b.Events.PublishEvent(event)
			event = nil
		}
		r.Close()
	}
	return nil

}

func (f *stFTP) CopyFiles(file string, bt *Ftpbeat) error {
	fmt.Println("FTP CopyFiles")
	r, err := f.con.Retr(file)
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

func (f *stFTP) Quit() {
	fmt.Println("FTP Quit")
	if f.con != nil {
		f.con.Quit()
	}
}
