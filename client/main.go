package main

import (
	"context"
	modelpb "daclient/pb/inventory"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const chunkSize = 64 * 1024

type uploader struct {
	dir         string
	client      modelpb.ModelOprServiceClient
	ctx         context.Context
	wg          sync.WaitGroup
	requests    chan string // each request is a filepath on client accessible to client
	DoneRequest chan string
	FailRequest chan string
}

//NewUploader creates a object of type uploader and creates fixed worker goroutines/threads
func NewUploader(ctx context.Context, client modelpb.ModelOprServiceClient, dir string) *uploader {
	d := &uploader{
		ctx:         ctx,
		client:      client,
		dir:         dir,
		requests:    make(chan string),
		DoneRequest: make(chan string),
		FailRequest: make(chan string),
	}
	for i := 0; i < 1; i++ {
		d.wg.Add(1)
		go d.worker(i + 1)
	}
	return d
}

func UploadFiles(ctx context.Context, client modelpb.ModelOprServiceClient, filepathlist []string, dir string) error {

	d := NewUploader(ctx, client, dir)
	var errorUploadbulk error

	if dir != "" {

		files, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			for _, file := range files {

				if !file.IsDir() {

					d.Do(dir + "/" + file.Name())

				}
			}
		}()

		for _, file := range files {
			if !file.IsDir() {
				select {

				case <-d.DoneRequest:

					//fmt.Println("sucessfully sent :" + req)

				case req := <-d.FailRequest:

					fmt.Println("failed to  send " + req)
					errorUploadbulk = errors.Wrapf(errorUploadbulk, " Failed to send %s", req)

				}
			}
		}
		fmt.Println("All done ")
	} else {

		go func() {
			for _, file := range filepathlist {
				d.Do(file)
			}
		}()

		defer d.Stop()

		for i := 0; i < len(filepathlist); i++ {
			select {

			case <-d.DoneRequest:
			//	fmt.Println("sucessfully sent " + req)
			case req := <-d.FailRequest:
				fmt.Println("failed to  send " + req)
				errorUploadbulk = errors.Wrapf(errorUploadbulk, " Failed to send %s", req)
			}
		}

	}

	return errorUploadbulk
}
func (d *uploader) Do(filepath string) {
	d.requests <- filepath
}
func (d *uploader) Stop() {
	close(d.requests)
	d.wg.Wait()
}
func (d *uploader) worker(workerID int) {
	defer d.wg.Done()
	var (
		buf        []byte
		firstChunk bool
	)
	for request := range d.requests {

		//open
		//.Println("Processsing " + request)
		file, errOpen := os.Open(request)
		if errOpen != nil {
			errOpen = errors.Wrapf(errOpen,
				"failed to open file %s",
				request)
			return
		}

		defer file.Close()

		//start uploader
		streamUploader, err := d.client.UploadStandardVer(d.ctx)
		if err != nil {
			err = errors.Wrapf(err,
				"failed to create upload stream for file %s",
				request)
			return
		}
		defer streamUploader.CloseSend()
		_, errstat := file.Stat()
		if errstat != nil {
			err = errors.Wrapf(err,
				"Unable to get file size  %s",
				request)
			return
		}

		//create a buffer of chunkSize to be streamed
		buf = make([]byte, chunkSize)
		firstChunk = true
		for {
			n, errRead := file.Read(buf)
			if errRead != nil {
				if errRead == io.EOF {
					errRead = nil
					break
				}

				errRead = errors.Wrapf(errRead,
					"errored while copying from file to buf")
				return
			}
			if firstChunk {
				err = streamUploader.Send(&modelpb.FileUploadRequest{
					Filename: request,
					Content:  buf[:n],
				})
				firstChunk = false
			} else {
				err = streamUploader.Send(&modelpb.FileUploadRequest{
					Content: buf[:n],
				})
			}
			if err != nil {
				break
				//bar.Reset(0)
				//return
			}
		}
		status, err := streamUploader.CloseAndRecv()

		if err != nil { //retry needed
			fmt.Println("failed to receive upstream status response")

			d.FailRequest <- request
			return
		}

		if status.Status != modelpb.Status_SUCCESS { //retry needed
			d.FailRequest <- request
			return
		}
		//fmt.Println("writing done for : " + request + " by " + strconv.Itoa(workerID))
		d.DoneRequest <- request
	}

}

func main() {
	options := []grpc.DialOption{}

	options = append(options, grpc.WithInsecure())

	conn, err := grpc.Dial("localhost:5000", options...)
	if err != nil {
		log.Fatalf("cannot connect: %v", err)
	}
	defer conn.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go UploadFiles(context.Background(), modelpb.NewModelOprServiceClient(conn), []string{}, "/home/licongchao/TTemp/2021-09-23")
	wg.Wait()

}
