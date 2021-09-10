package main

import (
	"fmt"

	feequest "github.com/fenggshu/transport/client/Feequest"
	"google.golang.org/grpc"
)

func main() {
	sid := feequest.GetLocalIdentity()
	fmt.Println(sid)

	conn, err := grpc.Dial("10.53.8.119:21", grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		fmt.Println(err)
	}
	a := feequest.GetFeeFile(conn, feequest.GetLocalIdentity())
	for _, v := range a {
		feequest.GetFeeData(conn, v)
	}

}
