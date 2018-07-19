package main

type Element interface {
	ChooseDiePool() error
}

func ChooseDice(e Element) error {
	err := e.ChooseDiePool()
	if err != nil {
		return err
	}
	return nil
}
