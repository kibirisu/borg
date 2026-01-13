package ap

import "github.com/kibirisu/borg/internal/domain"

type Collectioner[T any] interface {
	Objecter[Collection[T]]
}

type ActorCollectioner interface {
	Collectioner[Actor]
}

type NoteCollectioner interface {
	Collectioner[Note]
}

type CollectionPager[T any] interface {
	Objecter[CollectionPage[T]]
}

type ActorCollectionPager interface {
	CollectionPager[Actor]
}

type NoteCollectionPager interface {
	CollectionPager[Note]
}

type Collection[T any] struct {
	ID    string
	Type  string
	First CollectionPager[T]
}

type CollectionPage[T any] struct {
	ID     string
	Type   string
	Next   CollectionPager[T]
	PartOf Collectioner[T]
	Items  []Objecter[T]
}

type collection struct {
	object
}

type actorCollection struct {
	collection
}

type noteCollection struct {
	collection
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
	_ Collectioner[any]    = (*collection)(nil)
	_ ActorCollectioner    = (*actorCollection)(nil)
	_ NoteCollectioner     = (*noteCollection)(nil)
	_ CollectionPager[any] = (*collectionPage)(nil)
	_ ActorCollectionPager = (*actorCollectionPage)(nil)
	_ NoteCollectionPager  = (*noteCollectionPage)(nil)
)

func NewActorCollection(from *domain.ObjectOrLink) ActorCollectioner {
	return &actorCollection{collection{object{from}}}
}

func NewEmptyActorCollection() ActorCollectioner {
	return &actorCollection{collection{object{}}}
}

func NewNoteCollection(from *domain.ObjectOrLink) NoteCollectioner {
	return &noteCollection{collection{object{from}}}
}

func NewEmptyNoteCollection() NoteCollectioner {
	return &noteCollection{collection{object{}}}
}

func NewActorCollectionPage(from *domain.ObjectOrLink) ActorCollectionPager {
	return &actorCollectionPage{collectionPage{object{from}}}
}

func NewEmptyActorCollectionPage() ActorCollectionPager {
	return &actorCollectionPage{collectionPage{object{}}}
}

func NewNoteCollectionPage(from *domain.ObjectOrLink) NoteCollectionPager {
	return &noteCollectionPage{collectionPage{object{from}}}
}

func NewEmptyNoteCollectionPage() NoteCollectionPager {
	return &noteCollectionPage{collectionPage{object{}}}
}

// GetObject implements Collectioner.
// Subtle: this method shadows the method (object).GetObject of collection.object.
func (c *collection) GetObject() Collection[any] {
	obj := c.raw.Object
	return Collection[any]{
		ID:    obj.ID,
		Type:  obj.Type,
		First: &collectionPage{object{&obj.Collection.First}},
	}
}

// SetObject implements Collectioner.
// Subtle: this method shadows the method (object).SetObject of collection.object.
func (c *collection) SetObject(collection Collection[any]) {
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

// WithLink implements Collectioner.
// Subtle: this method shadows the method (object).WithLink of collection.object.
func (c *collection) WithLink(link string) Objecter[Collection[any]] {
	c.SetLink(link)
	return c
}

// WithObject implements Collectioner.
// Subtle: this method shadows the method (object).WithObject of collection.object.
func (c *collection) WithObject(collection Collection[any]) Objecter[Collection[any]] {
	c.SetObject(collection)
	return c
}

// GetObject implements ActorCollectioner.
// Subtle: this method shadows the method (collection).GetObject of actorCollection.collection.
func (a *actorCollection) GetObject() Collection[Actor] {
	return Collection[Actor]{
		ID:    a.raw.Object.ID,
		Type:  a.raw.Object.Type,
		First: &actorCollectionPage{collectionPage{object{&a.raw.Object.Collection.First}}},
	}
}

// SetObject implements ActorCollectioner.
// Subtle: this method shadows the method (collection).SetObject of actorCollection.collection.
func (a *actorCollection) SetObject(collection Collection[Actor]) {
	a.raw = &domain.ObjectOrLink{
		Object: &domain.Object{
			ID:   collection.ID,
			Type: collection.Type,
			Collection: &domain.Collection{
				First: *collection.First.GetRaw(),
			},
		},
	}
}

// WithLink implements ActorCollectioner.
// Subtle: this method shadows the method (collection).WithLink of actorCollection.collection.
func (a *actorCollection) WithLink(link string) Objecter[Collection[Actor]] {
	a.SetLink(link)
	return a
}

// WithObject implements ActorCollectioner.
// Subtle: this method shadows the method (collection).WithObject of actorCollection.collection.
func (a *actorCollection) WithObject(collection Collection[Actor]) Objecter[Collection[Actor]] {
	a.SetObject(collection)
	return a
}

// GetObject implements NoteCollectioner.
// Subtle: this method shadows the method (collection).GetObject of noteCollection.collection.
func (n *noteCollection) GetObject() Collection[Note] {
	return Collection[Note]{
		ID:    n.raw.Object.ID,
		Type:  n.raw.Object.Type,
		First: &noteCollectionPage{collectionPage{object{&n.raw.Object.Collection.First}}},
	}
}

// SetObject implements NoteCollectioner.
// Subtle: this method shadows the method (collection).SetObject of noteCollection.collection.
func (n *noteCollection) SetObject(collection Collection[Note]) {
	n.raw = &domain.ObjectOrLink{
		Object: &domain.Object{
			ID:   collection.ID,
			Type: collection.Type,
			Collection: &domain.Collection{
				First: *collection.First.GetRaw(),
			},
		},
	}
}

// WithLink implements NoteCollectioner.
// Subtle: this method shadows the method (collection).WithLink of noteCollection.collection.
func (n *noteCollection) WithLink(link string) Objecter[Collection[Note]] {
	n.SetLink(link)
	return n
}

// WithObject implements NoteCollectioner.
// Subtle: this method shadows the method (collection).WithObject of noteCollection.collection.
func (n *noteCollection) WithObject(collection Collection[Note]) Objecter[Collection[Note]] {
	n.SetObject(collection)
	return n
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
		Next:   &collectionPage{object{obj.CollectionPage.Next}},
		PartOf: &collection{object{&obj.CollectionPage.PartOf}},
		Items:  items,
	}
}

// SetObject implements CollectionPager.
// Subtle: this method shadows the method (object).SetObject of collectionPage.object.
func (c *collectionPage) SetObject(page CollectionPage[any]) {
	items := mapToRaw(page.Items)
	c.raw = &domain.ObjectOrLink{
		Object: &domain.Object{
			ID:   page.ID,
			Type: page.Type,
			CollectionPage: &domain.CollectionPage{
				Next:   page.Next.GetRaw(),
				PartOf: *page.PartOf.GetRaw(),
				Items:  items,
			},
		},
	}
}

// WithLink implements CollectionPager.
// Subtle: this method shadows the method (object).WithLink of collectionPage.object.
func (c *collectionPage) WithLink(link string) Objecter[CollectionPage[any]] {
	c.SetLink(link)
	return c
}

// WithObject implements CollectionPager.
// Subtle: this method shadows the method (object).WithObject of collectionPage.object.
func (c *collectionPage) WithObject(
	page CollectionPage[any],
) Objecter[CollectionPage[any]] {
	c.SetObject(page)
	return c
}

// GetObject implements ActorCollectionPager.
// Subtle: this method shadows the method (collectionPage).GetObject of actorCollectionPage.collectionPage.
func (a *actorCollectionPage) GetObject() CollectionPage[Actor] {
	obj := a.raw.Object
	items := []Objecter[Actor]{}
	for _, item := range a.raw.Object.CollectionPage.Items {
		items = append(items, &actor{object: object{&item}})
	}
	return CollectionPage[Actor]{
		ID:     obj.ID,
		Type:   obj.Type,
		Next:   &actorCollectionPage{collectionPage{object{obj.CollectionPage.Next}}},
		PartOf: &actorCollection{collection{object{&obj.CollectionPage.PartOf}}},
		Items:  items,
	}
}

// SetObject implements ActorCollectionPager.
// Subtle: this method shadows the method (collectionPage).SetObject of actorCollectionPage.collectionPage.
func (a *actorCollectionPage) SetObject(page CollectionPage[Actor]) {
	items := mapToRaw(page.Items)
	a.raw = &domain.ObjectOrLink{
		Object: &domain.Object{
			ID:   page.ID,
			Type: page.Type,
			CollectionPage: &domain.CollectionPage{
				Next:   page.Next.GetRaw(),
				PartOf: *page.PartOf.GetRaw(),
				Items:  items,
			},
		},
	}
}

// WithLink implements ActorCollectionPager.
// Subtle: this method shadows the method (collectionPage).WithLink of actorCollectionPage.collectionPage.
func (a *actorCollectionPage) WithLink(link string) Objecter[CollectionPage[Actor]] {
	a.SetLink(link)
	return a
}

// WithObject implements ActorCollectionPager.
// Subtle: this method shadows the method (collectionPage).WithObject of actorCollectionPage.collectionPage.
func (a *actorCollectionPage) WithObject(
	page CollectionPage[Actor],
) Objecter[CollectionPage[Actor]] {
	a.SetObject(page)
	return a
}

// GetObject implements NoteCollectionPager.
// Subtle: this method shadows the method (object).GetObject of noteCollectionPage.object.
func (n *noteCollectionPage) GetObject() CollectionPage[Note] {
	obj := n.raw.Object
	items := []Objecter[Note]{}
	for _, item := range n.raw.Object.CollectionPage.Items {
		items = append(items, &note{object: object{&item}})
	}
	return CollectionPage[Note]{
		ID:     obj.ID,
		Type:   obj.Type,
		Next:   &noteCollectionPage{collectionPage{object{obj.CollectionPage.Next}}},
		PartOf: &noteCollection{collection{object{&obj.CollectionPage.PartOf}}},
		Items:  items,
	}
}

// SetObject implements NoteCollectionPager.
// Subtle: this method shadows the method (object).SetObject of noteCollectionPage.object.
func (n *noteCollectionPage) SetObject(page CollectionPage[Note]) {
	items := mapToRaw(page.Items)
	n.raw = &domain.ObjectOrLink{
		Object: &domain.Object{
			ID:   page.ID,
			Type: page.Type,
			CollectionPage: &domain.CollectionPage{
				Next:   page.Next.GetRaw(),
				PartOf: *page.PartOf.GetRaw(),
				Items:  items,
			},
		},
	}
}

// WithLink implements NoteCollectionPager.
// Subtle: this method shadows the method (collectionPage).WithLink of noteCollectionPage.collectionPage.
func (n *noteCollectionPage) WithLink(link string) Objecter[CollectionPage[Note]] {
	n.SetLink(link)
	return n
}

// WithObject implements NoteCollectionPager.
// Subtle: this method shadows the method (collectionPage).WithObject of noteCollectionPage.collectionPage.
func (n *noteCollectionPage) WithObject(page CollectionPage[Note]) Objecter[CollectionPage[Note]] {
	n.SetObject(page)
	return n
}

func mapToRaw[T any](objects []Objecter[T]) []domain.ObjectOrLink {
	items := make([]domain.ObjectOrLink, len(objects))
	for idx, item := range objects {
		items[idx] = *item.GetRaw()
	}
	return items
}
