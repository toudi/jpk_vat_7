package uploader

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type requestedFileHeader struct {
	Key   string
	Value string
}
type requestedFile struct {
	URL        string `json:"Url"`
	HeaderList []requestedFileHeader
	BlobName   string
}

type initUploadSignedResponseType struct {
	ReferenceNumber         string
	RequestToUploadFileList []requestedFile
}

var initUploadSignedResponse initUploadSignedResponseType

type statusType struct {
	Code        int
	Description string
	Details     string
}

var documentStatus statusType

func (u *Uploader) initUploadSigned() (*initUploadSignedResponseType, error) {
	var err error
	var reader *bytes.Reader

	saftMetadataContent, err := ioutil.ReadFile(u.sourceFile)
	reader = bytes.NewReader(saftMetadataContent)

	err = xml.Unmarshal(saftMetadataContent, saftMetadata)
	if err != nil {
		return nil, fmt.Errorf("Nie udało się odczytać metadanych pliku JPK: %v", err)
	}

	request, err := http.NewRequest("POST", u.gatewayURL()+"/api/Storage/InitUploadSigned", reader)

	if err != nil {
		return nil, fmt.Errorf("Nie udało się stworzyć żądania HTTP: %v", err)
	}

	request.Header.Set("Content-Type", "application/xml")
	response, err := httpClient.Do(request)
	if response.Body != nil {
		defer response.Body.Close()
	}
	content, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode != 200 || err != nil {
		return nil, fmt.Errorf("Nie udało się wysłać żądania HTTP: %v; kod odpowiedzi: %v; błąd: %v", string(content), response.StatusCode, err)
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&initUploadSignedResponse)

	if err != nil {
		return nil, fmt.Errorf("Nie udało się sparsować odpowiedzi: %v", err)
	}

	return &initUploadSignedResponse, nil
}
