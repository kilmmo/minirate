package checkrate

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fenggshu/transport/msg"
)

type Fileinfos []msg.FeeInfo

func ProduceFileInfo(basedir string, stationid string) []msg.FeeInfo {
	var fs Fileinfos
	filepath.Walk(basedir, func(path string, info os.FileInfo, err error) error {

		if strings.HasPrefix(info.Name(), stationid) {
			tmp, _ := os.Open(path)
			var fi msg.FeeInfo
			fullpath := path
			fi.FileName = fullpath
			fi.Size = info.Size()
			hash := md5.New()
			io.Copy(hash, tmp)
			fi.Md5 = hex.EncodeToString(hash.Sum(nil))
			fs = append(fs, fi)
		}
		return nil
	})

	return fs

}
