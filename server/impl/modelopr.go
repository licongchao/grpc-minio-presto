package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	modelpb "da/modelopr/pb/inventory"

	"github.com/pkg/errors"
)

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

// rpc saveModel(stream UploadModelObjReq) returns (UploadModelObjResp) {}
func (s *ServerGRPC) SaveModel(stream modelpb.ModelOprService_SaveModelClient) (err error) {
	firstChunk := true
	var fp *os.File

	var fileData *modelpb.UploadRequestType

	var filename string
	for {

		fileData, err = stream.Recv() //ignoring the data  TO-Do save files received

		if err != nil {
			if err == io.EOF {
				break
			}

			err = errors.Wrapf(err,
				"failed unexpectadely while reading chunks from stream")
			return
		}

		if firstChunk { //first chunk contains file name

			if fileData.Filename != "" { //create file

				fp, err = os.Create(path.Join(s.destDir, filepath.Base(fileData.Filename)))

				if err != nil {
					s.logger.Error().Msg("Unable to create file  :" + fileData.Filename)
					stream.SendAndClose(&modelpb.UploadResponseType{
						Message: "Unable to create file :" + fileData.Filename,
						Code:    modelpb.UploadStatusCode_Failed,
					})
					return
				}
				defer fp.Close()
			} else {
				s.logger.Error().Msg("FileName not provided in first chunk  :" + fileData.Filename)
				stream.SendAndClose(&modelpb.UploadResponseType{
					Message: "FileName not provided in first chunk:" + fileData.Filename,
					Code:    modelpb.UploadStatusCode_Failed,
				})
				return

			}
			filename = fileData.Filename
			firstChunk = false
		}

		err = writeToFp(fp, fileData.Content)
		if err != nil {
			s.logger.Error().Msg("Unable to write chunk of filename :" + fileData.Filename + " " + err.Error())
			stream.SendAndClose(&modelpb.UploadResponseType{
				Message: "Unable to write chunk of filename :" + fileData.Filename,
				Code:    modelpb.UploadStatusCode_Failed,
			})
			return
		}
	}

	//s.logger.Info().Msg("upload received")
	err = stream.SendAndClose(&modelpb.UploadResponseType{
		Message: "Upload received with success",
		Code:    modelpb.UploadStatusCode_Ok,
	})
	if err != nil {
		err = errors.Wrapf(err,
			"failed to send status code")
		return
	}
	fmt.Println("Successfully received and stored the file :" + filename + " in " + s.destDir)
	return
}
