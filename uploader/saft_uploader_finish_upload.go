package uploader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type finishUploadType struct {
	ReferenceNumber   string
	AzureBlobNameList []string
}

type finishUploadResponseType struct {
	Message string
	Errors  []string
}

var finishUploadPayload finishUploadType
var finishUploadResponse finishUploadResponseType

func init() {
	finishUploadPayload.AzureBlobNameList = make([]string, 1)
}

func (u *Uploader) finishUploadSession() error {
	var err error
	var finishUploadPayloadJSON = new(bytes.Buffer)
	if err = json.NewEncoder(finishUploadPayloadJSON).Encode(finishUploadPayload); err != nil {
		return fmt.Errorf("Nie udało się zakodować parametrów wejściowych do metody FinishUpload: %v", err)
	}

	request, err := http.NewRequest("POST", u.gatewayURL()+"/api/Storage/FinishUpload", finishUploadPayloadJSON)
	request.Header.Set("Content-Type", "application/json")
	response, err := httpClient.Do(request)
	if err != nil || response.StatusCode != 200 {
		return fmt.Errorf("Błąd podczas kończenia sesji upload: %v", err)
	}

	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&finishUploadResponse); err != nil {
		return fmt.Errorf("Nie udało się zdekodować odpowiedzi z FinishUpload: %v", err)
	}

	return nil
}
