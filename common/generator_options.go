package common

type GeneratorOptions struct {
	Verbose                bool
	EncodingConversionFile string
	UseCurrentDir          bool
	GenerateMetadata       bool
	CSVDelimiter           string
	CSVSections            bool
	XLSXSpreadsheetName    string
	KSeFRegistry           string
}
