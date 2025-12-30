package ap

type Collectioner interface {
	Objecter[Collection]
}

type CollectionPager interface {
	Collectioner
}

type Collection struct {
	// TODO
}

type collection struct {
	object
}

// GetObject implements Collectioner.
// Subtle: this method shadows the method (object).GetObject of collection.object.
func (c *collection) GetObject() Collection {
	panic("unimplemented")
}

// SetObject implements Collectioner.
// Subtle: this method shadows the method (object).SetObject of collection.object.
func (c *collection) SetObject(Collection) {
	panic("unimplemented")
}

var _ Collectioner = (*collection)(nil)
