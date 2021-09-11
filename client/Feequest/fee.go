package feequest

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/fenggshu/transport/msg"
	"google.golang.org/grpc"
)

const mtu = 512 * 1024

func GetLocalIdentity() string {
	url := "sqlserver://ynstation:stn@1522@localhost?database=ynstation&connection+timeout=30"
	db, err := sql.Open("mssql", url)

	if err != nil {
		println("Open Error:", err)
	}
	defer db.Close()
	var a string
	rows, err := db.Query("select paravalue from ts_sysparadic where para=1;")

	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		rows.Scan(&a)
	}
	defer rows.Close()

	return a
}

func GetFeeFile(conn *grpc.ClientConn, sid string) []*msg.FeeInfo {

	mc := msg.NewFeerequestClient(conn)
	res, _ := mc.ReqFeeInfo(context.Background(), &msg.Station{
		Id: sid,
	})
	var fis []*msg.FeeInfo

	for {
		feeinfo, err := res.Recv()

		//fmt.Println(feeinfo.GetFileName() + "\t" + strconv.FormatInt(feeinfo.GetSize(), 10) + "\t" + feeinfo.GetMd5())
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
		fis = append(fis, feeinfo)
	}
	return fis

}

func CheckFile(m *msg.FeeInfo) bool {
	_, filename := filepath.Split(m.GetFileName())
	basedir, _ := os.Getwd()
	_, err := os.Open(filepath.Join(basedir, filename))
	if err != nil {
		return false
	}
	fstat, _ := os.Stat(filename)
	if fstat.Size() != m.GetSize() {
		fmt.Println("size not ok")
		return false
	}

	fh, _ := os.Open(filename)

	hash := md5.New()
	io.Copy(hash, fh)
	if hex.EncodeToString(hash.Sum(nil)) == m.GetMd5() {
		fmt.Println("Files: " + filename + " found and no need to download")
	}

	return hex.EncodeToString(hash.Sum(nil)) == m.GetMd5()

}

func GetFeeData(conn *grpc.ClientConn, m *msg.FeeInfo) {
	totalpart := m.GetSize() / mtu
	if m.GetSize()%mtu != 0 {
		totalpart = totalpart + 1
	}

	fc := msg.NewFeerequestClient(conn)
loop:
	var file []*msg.PartData
	for i := 0; i < int(totalpart); i++ {
		var msize int64 = mtu
		if i == int(totalpart)-1 && m.GetSize()%mtu != 0 {
			msize = m.GetSize() % mtu
		}

		data, err := fc.ReqFilePart(context.Background(), &msg.PartInfo{
			Filename: m.GetFileName(),
			Partsize: msize,
			Partid:   int64(i),
		})

		if err != nil {
			fmt.Println(err)
			goto loop
		}
		file = append(file, data)
	}
	_, fn := filepath.Split(m.GetFileName())
	fh, _ := os.Create(fn)
	for _, v := range file {
		fh.Write(v.GetData())
	}

}
