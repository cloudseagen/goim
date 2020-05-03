package comet

import (
	"log"
	"os"

	"github.com/cloudseagen/goim/internal/comet/conf"
)

var whitelist *Whitelist

// Whitelist .
type Whitelist struct {
	log  *log.Logger
	list map[string]struct{} // whitelist for debug
}

// InitWhitelist a whitelist struct.
func InitWhitelist(c *conf.Whitelist) (err error) {
	var (
		mid string
		f   *os.File
	)
	if f, err = os.OpenFile(c.WhiteLog, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644); err == nil {
		whitelist = new(Whitelist)
		whitelist.log = log.New(f, "", log.LstdFlags)
		whitelist.list = make(map[string]struct{})
		for _, mid = range c.Whitelist {
			whitelist.list[mid] = struct{}{}
		}
	}
	return
}

// Contains whitelist contains a mid or not.
func (w *Whitelist) Contains(mid string) (ok bool) {
	if len(mid) > 0 {
		_, ok = w.list[mid]
	}
	return
}

// Printf calls l.Output to print to the logger.
func (w *Whitelist) Printf(format string, v ...interface{}) {
	w.log.Printf(format, v...)
}
