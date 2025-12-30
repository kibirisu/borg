package ap

import "github.com/kibirisu/borg/internal/domain"

type Collectioner interface {
	Objecter[Collection]
}

type CollectionPager interface {
	Objecter[CollectionPage]
}

type Collection struct {
	ID    string
	Type  string
	First CollectionPager
}

type CollectionPage struct {
	ID     string
	Type   string
	Next   CollectionPager
	PartOf Collectioner
	Items  []Noter
}

type collection struct {
	object
}

type collectionPage struct {
	object
}

var (
	_ Collectioner    = (*collection)(nil)
	_ CollectionPager = (*collectionPage)(nil)
)

// GetObject implements Collectioner.
// Subtle: this method shadows the method (object).GetObject of collection.object.
func (c *collection) GetObject() Collection {
	obj := c.raw.Object
	return Collection{
		ID:    obj.ID,
		Type:  obj.Type,
		First: &collectionPage{object{&obj.Collection.First}},
	}
}

// SetObject implements Collectioner.
// Subtle: this method shadows the method (object).SetObject of collection.object.
func (c *collection) SetObject(collection Collection) {
	c.raw.Object.ID = collection.ID
	c.raw.Object.Type = collection.Type
	c.raw.Object.Collection.First = *collection.First.GetRaw()
}

// GetObject implements CollectionPager.
// Subtle: this method shadows the method (object).GetObject of collectionPage.object.
func (c *collectionPage) GetObject() CollectionPage {
	obj := c.raw.Object
	items := []Noter{}
	for _, item := range c.raw.Object.CollectionPage.Items {
		items = append(items, &note{object{&item}})
	}
	return CollectionPage{
		ID:     obj.ID,
		Type:   obj.Type,
		Next:   &collectionPage{object{&obj.CollectionPage.Next}},
		PartOf: &collection{object{&obj.CollectionPage.PartOf}},
		Items:  items,
	}
}

// SetObject implements CollectionPager.
// Subtle: this method shadows the method (object).SetObject of collectionPage.object.
func (c *collectionPage) SetObject(page CollectionPage) {
	items := []domain.ObjectOrLink{}
	for _, item := range page.Items {
		items = append(items, *item.GetRaw())
	}
	c.raw.Object.ID = page.ID
	c.raw.Object.Type = page.Type
	c.raw.Object.CollectionPage.Next = *page.Next.GetRaw()
	c.raw.Object.CollectionPage.PartOf = *page.PartOf.GetRaw()
	c.raw.Object.CollectionPage.Items = items
}
