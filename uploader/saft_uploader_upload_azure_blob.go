package uploader

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
)

func (u *Uploader) uploadFileToAzureBlob(saftSessionParams *initUploadSignedResponseType) error {
	var err error

	uploadedFile := path.Join(u.workdir, saftMetadata.FileName+".aes")

	fileBytes, err := ioutil.ReadFile(uploadedFile)
	if err != nil {
		return fmt.Errorf("Nie udało się odczytac pliku do przesłania: %v", err)
	}
	fileBody := bytes.NewReader(fileBytes)

	fileUploadRequest, err := http.NewRequest("PUT", saftSessionParams.RequestToUploadFileList[0].URL, fileBody)

	for _, header := range saftSessionParams.RequestToUploadFileList[0].HeaderList {
		fileUploadRequest.Header.Set(header.Key, header.Value)
	}

	response, err := httpClient.Do(fileUploadRequest)

	if response.StatusCode != 201 {
		return fmt.Errorf("Błąd wysyłania pliku do Azure: %v", response.Status)
	}
	finishUploadPayload.AzureBlobNameList[0] = saftSessionParams.RequestToUploadFileList[0].BlobName
	response.Body.Close()
	return nil
}
