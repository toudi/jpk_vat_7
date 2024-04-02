package uploader

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func (u *Uploader) uploadFileToAzureBlob(saftSessionParams *initUploadSignedResponseType) error {
	var err error

	uploadedFile := filepath.Join(u.workdir, saftMetadata.FileName)

	fileBytes, err := os.ReadFile(uploadedFile)
	if err != nil {
		return fmt.Errorf("Nie udało się odczytac pliku do przesłania: %v", err)
	}
	fileBody := bytes.NewReader(fileBytes)

	fileUploadRequest, err := http.NewRequest(
		"PUT",
		saftSessionParams.RequestToUploadFileList[0].URL,
		fileBody,
	)

	if err != nil {
		return errors.Join(err, fmt.Errorf("błąd tworzenia requestu do azure"))
	}

	for _, header := range saftSessionParams.RequestToUploadFileList[0].HeaderList {
		fileUploadRequest.Header.Set(header.Key, header.Value)
	}

	response, err := httpClient.Do(fileUploadRequest)
	if err != nil {
		return errors.Join(err, fmt.Errorf("błąd wykonania requestu PUT do azure"))
	}

	if response.StatusCode != 201 {
		return fmt.Errorf("Błąd wysyłania pliku do Azure: %v", response.Status)
	}
	finishUploadPayload.AzureBlobNameList[0] = saftSessionParams.RequestToUploadFileList[0].BlobName
	response.Body.Close()
	return nil
}
