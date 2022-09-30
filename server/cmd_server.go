package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/urfave/cli"

	"da/httpserver"
	modelpb "da/pb/inventory"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"google.golang.org/grpc"
)

// type modelserver struct {
// 	modelpb.UnimplementedModelOprServiceServer
// }

type Server interface {
	Listen() (err error)
	Close()
}
type ServerGRPC struct {
	logger  zerolog.Logger
	server  *grpc.Server
	Address string

	certificate string
	key         string
	destDir     string

	modelpb.UnimplementedModelOprServiceServer
}

type ServerGRPCConfig struct {
	Certificate string
	Key         string
	Address     string
	DestDir     string
}

//writeToFp takes in a file pointer and byte array and writes the byte array into the file
//returns error if pointer is nil or error in writing to file

func (s *ServerGRPC) Close() {
	if s.server != nil {
		s.server.Stop()
	}
}

// string Filename = 1;
// string Latestver = 2;
// string Stagingver = 3;
func (s *ServerGRPC) GetFilesVer(context.Context, *empty.Empty) (*modelpb.FileInfoResponse, error) {
	fileInfo1 := &modelpb.FileInfo{
		Filename:  "1.text",
		Latestver: "123",
	}
	return &modelpb.FileInfoResponse{FileInfo: []*modelpb.FileInfo{fileInfo1}}, nil
}

func (s *ServerGRPC) UploadStandardVer(stream modelpb.ModelOprService_UploadStandardVerServer) (err error) {
	firstChunk := true

	var fp *os.File
	var fileData *modelpb.FileUploadRequest
	var filename string

	for {
		fileData, err = stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			err = errors.Wrapf(err, "failed unexpectedly while reading chunks from stream")
			return
		}

		if firstChunk {
			if fileData.Filename != "" { //create file
				fp, err = os.Create(path.Join(s.destDir, filepath.Base(fileData.Filename)))

				if err != nil {
					s.logger.Error().Msg("Unable to create file  :" + fileData.Filename)
					stream.SendAndClose(&modelpb.FileUploadResponse{
						Message: "Unable to create file :" + fileData.Filename,
						Status:  modelpb.Status_FAILED,
					})
					return
				}
				defer fp.Close()
			} else {
				s.logger.Error().Msg("FileName not provided in first chunk  :" + fileData.Filename)
				stream.SendAndClose(&modelpb.FileUploadResponse{
					Message: "FileName not provided in first chunk:" + fileData.Filename,
					Status:  modelpb.Status_FAILED,
				})
				return
			}
			filename = fileData.Filename
			firstChunk = false
		}

		err = writeToFp(fp, fileData.Content)
		if err != nil {
			s.logger.Error().Msg("Unable to write chunk of filename :" + fileData.Filename + " " + err.Error())
			stream.SendAndClose(&modelpb.FileUploadResponse{
				Message: "Unable to write chunk of filename :" + fileData.Filename,
				Status:  modelpb.Status_FAILED,
			})
			return
		}
	}
	err = stream.SendAndClose(&modelpb.FileUploadResponse{
		Message: "Upload received with success",
		Status:  modelpb.Status_SUCCESS,
	})
	if err != nil {
		err = errors.Wrapf(err,
			"failed to send status code")
		return
	}
	fmt.Println("Successfully received and stored the file :" + filename + " in " + s.destDir)
	return
}

func writeToFp(fp *os.File, data []byte) error {
	w := 0
	n := len(data)
	for {
		nw, err := fp.Write(data[w:])
		if err != nil {
			return err
		}
		w += nw
		if nw >= n {
			return nil
		}
	}

}

func NewServerGRPC(cfg ServerGRPCConfig) (s ServerGRPC, err error) {
	s.logger = zerolog.New(os.Stdout).
		With().
		Str("from", "server").
		Logger()

	if cfg.Address == "" {
		err = errors.Errorf("Address must be specified")
		s.logger.Error().Msg("Address must be specified")
		return
	}

	s.Address = cfg.Address
	s.certificate = cfg.Certificate
	s.key = cfg.Key

	if _, err = os.Stat(cfg.DestDir); err != nil {
		s.logger.Error().Msg("Directory doesnt exist")
		return
	}

	s.destDir = cfg.DestDir
	return
}

func (s *ServerGRPC) Listen() (err error) {
	// var (
	// 	listener  net.Listener
	// 	grpcOpts  = []grpc.ServerOption{}
	// 	grpcCreds credentials.TransportCredentials
	// )

	// listener, err = net.Listen("tcp", s.Address)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to listen on  %d",
			s.Address)
		return
	}

	// if s.certificate != "" && s.key != "" {
	// 	grpcCreds, err = credentials.NewServerTLSFromFile(
	// 		s.certificate, s.key)
	// 	if err != nil {
	// 		err = errors.Wrapf(err,
	// 			"failed to create tls grpc server using cert %s and key %s",
	// 			s.certificate, s.key)
	// 		return
	// 	}

	// 	grpcOpts = append(grpcOpts, grpc.Creds(grpcCreds))
	// }

	mux := httpserver.GetHTTPServeMux()

	// s.server = grpc.NewServer(grpcOpts...)
	s.server = grpc.NewServer()
	modelpb.RegisterModelOprServiceServer(s.server, s)

	// err = s.server.Serve(listener)
	// if err != nil {
	// 	err = errors.Wrapf(err, "errored listening for grpc connections")
	// 	return
	// }

	http.ListenAndServe(":10000",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
				s.server.ServeHTTP(w, r)
			} else {
				mux.ServeHTTP(w, r)
			}
		}),
	)
	return nil
}

func StartServerCommand() cli.Command {

	return cli.Command{
		Name:  "serve",
		Usage: "initiates a gRPC server",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "a",
				Usage: "Address to listen",
				Value: "localhost:5000",
			},

			&cli.StringFlag{
				Name:  "key",
				Usage: "path to TLS certificate",
			},
			&cli.StringFlag{
				Name:  "certificate",
				Usage: "path to TLS certificate",
			},
			&cli.StringFlag{
				Name:  "d",
				Usage: "Destrination directory Default is /tmp",
				Value: "/tmp",
			},
		},
		Action: func(c *cli.Context) error {
			grpcServer, err := NewServerGRPC(ServerGRPCConfig{
				Address:     c.String("a"),
				Certificate: c.String("certificate"),
				Key:         c.String("key"),
				DestDir:     c.String("d"),
			})
			if err != nil {
				fmt.Println("error is creating server")

				return err
			}
			server := &grpcServer
			err = server.Listen()
			if err != nil {
				return err
			}

			defer server.Close()
			return nil
		},
	}

}
