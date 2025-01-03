package vovanDB

func updateExecutor(s UpdateQuery) error {
	err := validateUpdateQuery(s)

	if err != nil {
		return err
	}

	return nil
}
