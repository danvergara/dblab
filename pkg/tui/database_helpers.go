package tui

func (t *Tui) showTables() ([]string, error) {
	tables, err := t.c.ShowTables()
	if err != nil {
		return nil, err
	}

	return tables, nil
}
