package main

import (
	"fmt"
	"time"

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
	for {
	loop:
		a := feequest.GetFeeFile(conn, feequest.GetLocalIdentity())

		for _, v := range a {
			if feequest.CheckFile(v) {
				continue
			} else {
				feequest.GetFeeData(conn, v)
			}

		}
		fmt.Println("rate download completed,program will check after 30 mins")
		timer1 := time.NewTimer(30 * time.Minute)
		<-timer1.C
		fmt.Println("30mins passed! now program will recheck and download files")
		goto loop
	}

}
