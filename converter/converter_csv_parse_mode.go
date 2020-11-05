package converter

const (
	// Pojedynczy plik CSV, gdzie sekcje zdefiniowane są w
	// linii nr. 1 a dane do sekcji rozdzielane są nowymi wierszami
	ParserModeSingleFile = iota
	// Pojedynczy plik CSV, gdzie jego struktura jest następująca:
	// SEKCJA;nazwa-sekcji
	// Kolumna;kolumna;kolumna
	// dane;dane;dane;dane;
	// SEKCJA;nazwa-sekcji-1
	// kolumna;kolumna; ...
	ParserModeSingleFileWithSections = iota
)
