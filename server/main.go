package main

import (
	"flag"
	"os"

	"github.com/urfave/cli"
)

var (
	port = flag.Int("port", 50000, "Model Service port")
)

// type server struct {
// 	modelpb.UnimplementedModelOprServiceServer
// }

// // SayHello implements helloworld.GreeterServer
// func (s *server) saveModel(ctx context.Context, in *modelpb.ModelObjReq) (*modelpb.ModelObjResp, error) {
// 	log.Printf("Received: %v", in.GetName())
// 	return &modelpb.ModelObjResp{}, nil
// }

func main() {
	app := cli.NewApp()
	app.Name = "Model Operation"
	app.Usage = "Model Version Storage"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		StartServerCommand(),
	}
	app.Run(os.Args)

	// flag.Parse()
	// lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }
	// defer lis.Close()

	// //
	// bbto := grocksdb.NewDefaultBlockBasedTableOptions()
	// bbto.SetBlockCache(grocksdb.NewLRUCache(3 << 30))

	// opts := grocksdb.NewDefaultOptions()
	// opts.SetBlockBasedTableFactory(bbto)
	// opts.SetCreateIfMissing(true)

	// db, err := grocksdb.OpenDb(opts, "/home/licongchao/workspace/DataPreparation/apps/golang/data-preparation/server/rocksdb")
	// ro := grocksdb.NewDefaultReadOptions()
	// wo := grocksdb.NewDefaultWriteOptions()

	// // if ro and wo are not used again, be sure to Close them.
	// err = db.Put(wo, []byte("foo"), []byte("bar"))
	// value, err := db.Get(ro, []byte("foo"))
	// fmt.Print(value)
	// defer value.Free()
	// err = db.Delete(wo, []byte("foo"))
	// //
	// rpcServer := grpc.NewServer()
	// modelpb.RegisterModelOprServiceServer(rpcServer, &server{})

	// log.Printf("Data Preparation server start to serve...")
	// if err := rpcServer.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v", err)
	// }
}
