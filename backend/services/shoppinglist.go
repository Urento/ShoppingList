package services

type Shoppinglist struct {
	//Id           string
	Title        string
	Items        []string
	Owner        string
	Participants []string
}

func (s *Shoppinglist) Create() error {

	return nil
}
