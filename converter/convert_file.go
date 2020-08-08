package converter

import "fmt"

func (c *Converter) convertFile() error {
	var err error

	logger.Debugf("Tryb konwersji pliku")

	if err = jpk.parsujCSV(c.source); err != nil {
		return fmt.Errorf("Błąd parsowania pliku CSV: %v", err)
	}

	return nil
}
