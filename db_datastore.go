package lorefarm

import (
	"fmt"
	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
)

type database struct {
	client *datastore.Client
}

// Ensure database conforms to the BookDatabase interface.
var _ PageDatabase = &database{}

func (db *database) AddPage(b *PageData) (id int64, err error) {
	ctx := context.Background()
	k := datastore.IncompleteKey("Page", nil)
	k, err = db.client.Put(ctx, k, b)
	if err != nil {
		return 0, fmt.Errorf("database: could not put Page: %v", err)
	}
	return k.ID, nil
}

func (db *database) GetPage(id int64) (*Page, error) {
	var pageData, err = db.getPageData(id)
	if(err != nil) {
		return nil, err
	}
	var result = Page {
		Current: pageData,
		Next: []*PageData{},
	}

	result.Next, err = db.getChildren(id)
	if(err != nil) {
		return nil, err
	}

	return &result, nil;
}

func newDatabase(client *datastore.Client) (PageDatabase, error) {
	ctx := context.Background()
	// Verify that we can communicate and authenticate with the datastore service.
	t, err := client.NewTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("database: could not connect: %v", err)
	}
	if err := t.Rollback(); err != nil {
		return nil, fmt.Errorf("database: could not connect: %v", err)
	}
	return &database{
		client: client,
	}, nil
}

func (db *database) pageKey(id int64) *datastore.Key {
	return datastore.IDKey("Page", id, nil)
}

func (db *database) getPageData(id int64) (*PageData, error) {
	ctx := context.Background()
	k := db.pageKey(id)
	page := &PageData{}
	if err := db.client.Get(ctx, k, page); err != nil {
		return nil, fmt.Errorf("database: could not get Page: %v", err)
	}
	page.Id = id
	return page, nil
}

func (db *database) getChildren(id int64) ([]*PageData, error) {
	ctx := context.Background()
	children := make([]*PageData, 0)
	q := datastore.NewQuery("Page").Filter("ParentId =", id)

	keys, err := db.client.GetAll(ctx, q, &children)

	if err != nil {
		return nil, fmt.Errorf("database: could not list children: %v", err)
	}

	for i, k := range keys {
		children[i].Id = k.ID
	}

	return children, nil
}
