package main

import (
	_ "expvar"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"context"

	"golang.org/x/net/trace"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/micro/examples/booking/api/demo/proto"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
)

// Demo is data type
type Demo struct {
	Client client.Client
}

// Testapi is interface api
func (d *Demo) Testapi(ctx context.Context, req *demo.Request, rsp *demo.Response) error {
	result := new(demo.Result)

	if req.Name == "uwei" {
		result.Msg = "Hi " + req.Name
		result.Code = 1
	} else {
		result.Msg = req.Name
		result.Code = -1
	}

	rsp.Result = result

	return nil

}

func main() {
	result := new(demo.Result)
	result.Msg = "this is from demo api"
	result.Code = -1

	fmt.Print(result)

	db, err := sql.Open("mysql", "uwei:uwei@/uwei")
	if err != nil {
		fmt.Print(err)
	} else {
		defer db.Close()
		rows, _ := db.Query("select * from uwei_table")
		for rows.Next() {
			var name string
			var id int
			if err := rows.Scan(&id, &name); err != nil {
				fmt.Print(err)
			}
			fmt.Printf("%s is %d\n", name, id)
		}
		if err := rows.Err(); err != nil {
			fmt.Print(err)
		}
	}

	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) {
		return true, true
	}
	service := micro.NewService(
		micro.Name("go.micro.api.demo"),
	)
	service.Init()
	demo.RegisterDemoHandler(service.Server(), &Demo{service.Client()})
	service.Run()

}
