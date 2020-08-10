package uploader

import (
	"fmt"
	"io/ioutil"
)

func (u *Uploader) uploadSAFTFile() error {
	var err error

	initUploadSignedResponse, err := u.initUploadSigned()
	if err != nil {
		return fmt.Errorf("Nie udało się wysłać pliku metadanych: %v", err)
	}

	logger.Debugf("numer referencyjny dokumentu: %s", initUploadSignedResponse.ReferenceNumber)
	finishUploadPayload.ReferenceNumber = initUploadSignedResponse.ReferenceNumber

	statusCheckURL := u.gatewayURL() + "/api/Storage/Status/" + initUploadSignedResponse.ReferenceNumber
	err = ioutil.WriteFile(u.saftRefNoFile(), []byte(statusCheckURL), 0644)
	if err != nil {
		return fmt.Errorf("Nie udało się zapisać numeru referencyjnego: %v", err)
	}

	err = u.uploadFileToAzureBlob(initUploadSignedResponse)

	if err != nil {
		return fmt.Errorf("Nie udało się wysłać pliku do Azure: %v", err)
	}

	if err = u.finishUploadSession(); err != nil {
		return fmt.Errorf("Nie udało się zakończyć sesji: %v", err)
	}

	return nil
}
