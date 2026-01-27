package v7m_3

import "io"

func (g *v7m_3) Save(writer io.Writer) error {
	if err := g.root.ApplyOrdering(JPK_V7M_3ChildrenOrder); err != nil {
		return err
	}

	return g.root.DumpToWriter(writer, 0)
}
