package uploader

type saftMetadataStruct struct {
	FileName string `xml:"DocumentList>Document>FileSignatureList>FileSignature>FileName"`
}

var saftMetadata = &saftMetadataStruct{}
