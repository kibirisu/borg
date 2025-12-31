package ap

import "github.com/kibirisu/borg/internal/domain"

type Collectioner interface {
	Objecter[Collection]
}

type CollectionPager[T any] interface {
	Objecter[CollectionPage[T]]
}

type NoteCollectionPager interface {
	CollectionPager[Noter]
}

type ActorCollectionPager interface {
	CollectionPager[Actorer]
}

type Collection struct {
	ID    string
	Type  string
	First CollectionPager[any]
}

type CollectionPage[T any] struct {
	ID     string
	Type   string
	Next   CollectionPager[T]
	PartOf Collectioner
	Items  []Objecter[T]
}

type collection struct {
	object
}

type collectionPage struct {
	object
}

type actorCollectionPage struct {
	collectionPage
}

type noteCollectionPage struct {
	collectionPage
}

var (
	_ Collectioner         = (*collection)(nil)
	_ CollectionPager[any] = (*collectionPage)(nil)
	_ ActorCollectionPager = (*actorCollectionPage)(nil)
	_ NoteCollectionPager  = (*noteCollectionPage)(nil)
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
	c.raw = &domain.ObjectOrLink{
		Object: &domain.Object{
			ID:   collection.ID,
			Type: collection.Type,
			Collection: &domain.Collection{
				First: *collection.First.GetRaw(),
			},
		},
	}
}

// GetObject implements CollectionPager.
// Subtle: this method shadows the method (object).GetObject of collectionPage.object.
func (c *collectionPage) GetObject() CollectionPage[any] {
	obj := c.raw.Object
	items := []Objecter[any]{}
	for _, item := range c.raw.Object.CollectionPage.Items {
		items = append(items, &object{&item})
	}
	return CollectionPage[any]{
		ID:     obj.ID,
		Type:   obj.Type,
		Next:   &collectionPage{object{&obj.CollectionPage.Next}},
		PartOf: &collection{object{&obj.CollectionPage.PartOf}},
		Items:  items,
	}
}

// SetObject implements CollectionPager.
// Subtle: this method shadows the method (object).SetObject of collectionPage.object.
func (c *collectionPage) SetObject(page CollectionPage[any]) {
	items := []domain.ObjectOrLink{}
	for _, item := range page.Items {
		items = append(items, *item.GetRaw())
	}
	c.raw = &domain.ObjectOrLink{
		Object: &domain.Object{
			ID:   page.ID,
			Type: page.Type,
			CollectionPage: &domain.CollectionPage{
				Next:   *page.Next.GetRaw(),
				PartOf: *page.PartOf.GetRaw(),
				Items:  items,
			},
		},
	}
}

// GetObject implements ActorCollectionPager.
// Subtle: this method shadows the method (collectionPage).GetObject of actorCollectionPage.collectionPage.
func (a *actorCollectionPage) GetObject() CollectionPage[Actorer] {
	panic("unimplemented")
}

// SetObject implements ActorCollectionPager.
// Subtle: this method shadows the method (collectionPage).SetObject of actorCollectionPage.collectionPage.
func (a *actorCollectionPage) SetObject(CollectionPage[Actorer]) {
	panic("unimplemented")
}

// GetObject implements NoteCollectionPager.
// Subtle: this method shadows the method (object).GetObject of noteCollectionPage.object.
func (n *noteCollectionPage) GetObject() CollectionPage[Noter] {
	panic("unimplemented")
}

// SetObject implements NoteCollectionPager.
// Subtle: this method shadows the method (object).SetObject of noteCollectionPage.object.
func (n *noteCollectionPage) SetObject(CollectionPage[Noter]) {
	panic("unimplemented")
}
