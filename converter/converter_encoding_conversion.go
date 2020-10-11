package converter

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/toudi/jpk_vat/common"
)

func (c *Converter) PrepareEncodingConversionTable(conversion_file string) {
	fileBytes, err := ioutil.ReadFile(conversion_file)

	if err != nil {
		logger.Errorf("Nie udało się otworzyć pliku z konwersją znaków")
		return
	}

	encodingConversion = make(map[byte]string)

	for _, line := range strings.Split(string(fileBytes), common.LineBreak) {
		mapping := strings.Split(line, ":")

		if len(mapping) == 2 {
			srcByteHex := strings.Trim(mapping[0], " ")
			dstChar := strings.Trim(mapping[1], " \r")

			srcByte, err := strconv.ParseUint(srcByteHex, 0, 8)
			if err == nil {
				encodingConversion[byte(srcByte)] = dstChar
			}
		}
	}

	logger.Debugf("Odczytano tablicę konwersji znaków: %v\n", encodingConversion)
}

func convertEncoding(input string) string {
	inputBytes := []byte(input)
	outputBytes := []byte{}
	dstByte := []byte{}

	for _, srcByte := range inputBytes {
		dstByte = []byte{srcByte}
		if dstChar, exists := encodingConversion[srcByte]; exists {
			dstByte = bytes.Replace(dstByte, []byte{srcByte}, []byte(dstChar), 1)
		}
		outputBytes = append(outputBytes, dstByte...)
	}

	return string(outputBytes)
}
