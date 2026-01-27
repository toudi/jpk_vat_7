package abstract

import "io"

type Generator interface {
	SetKSeFRegistryFile(registryFile string) error
	SetData(section string, data map[string]string) error
	Save(writer io.Writer) error
}

type GeneratorFactory func() Generator
