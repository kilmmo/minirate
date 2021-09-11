package feeresp

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/fenggshu/transport/msg"
	checkrate "github.com/fenggshu/transport/server/Checkrate"
	"google.golang.org/grpc"
)

const mtu = 512 * 1024

type server struct {
	msg.UnimplementedFeerequestServer
}

func (s *server) ReqFeeInfo(station *msg.Station, ms msg.Feerequest_ReqFeeInfoServer) error {
	fmt.Println(station.GetId())
	var infos []msg.FeeInfo
	infos = checkrate.ProduceFileInfo("/home/Toll2021", station.GetId())
	for _, v := range infos {
		fmt.Println(v.GetFileName() + "\t" + strconv.FormatInt(v.GetSize(), 10) + "\t" + v.GetMd5())
		ms.Send(&v)
	}
	return nil
}

func (s *server) ReqFilePart(ctx context.Context, m *msg.PartInfo) (*msg.PartData, error) {
	fmt.Println(m.GetFilename() + "\t part data:" + strconv.FormatInt(m.GetPartid(), 10))
	fh, _ := os.Open(m.GetFilename())
	if m.GetPartsize() == mtu {
		var data = make([]byte, mtu)
		fh.ReadAt(data, mtu*m.GetPartid())
		return &msg.PartData{
			Data: data,
		}, nil
	} else {
		var data = make([]byte, m.GetPartsize())
		fh.ReadAt(data, mtu*(m.GetPartid()))
		return &msg.PartData{
			Data: data,
		}, nil
	}
}

func InitServer() {
	conn, err := net.Listen("tcp", "10.53.8.119:21")
	if err != nil {
		fmt.Print(err)
	}

	gs := grpc.NewServer()
	msg.RegisterFeerequestServer(gs, &server{})
	gs.Serve(conn)
}
