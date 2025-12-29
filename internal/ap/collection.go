package ap

type Collectioner interface {
	ObjectOrLink[Collection]
}

type Collection struct {
	ID    string
	First Objecter
}
