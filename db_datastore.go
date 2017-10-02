// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package lorefarm

import (
	"fmt"

	"cloud.google.com/go/datastore"

	"golang.org/x/net/context"
)

// datastoreDB persists books to Cloud Datastore.
// https://cloud.google.com/datastore/docs/concepts/overview
type datastoreDB struct {
	client *datastore.Client
}

// Ensure datastoreDB conforms to the BookDatabase interface.
var _ PageDatabase = &datastoreDB{}

// newDatastoreDB creates a new BookDatabase backed by Cloud Datastore.
// See the datastore and google packages for details on creating a suitable Client:
// https://godoc.org/cloud.google.com/go/datastore
func newDatastoreDB(client *datastore.Client) (PageDatabase, error) {
	ctx := context.Background()
	// Verify that we can communicate and authenticate with the datastore service.
	t, err := client.NewTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("datastoredb: could not connect: %v", err)
	}
	if err := t.Rollback(); err != nil {
		return nil, fmt.Errorf("datastoredb: could not connect: %v", err)
	}
	return &datastoreDB{
		client: client,
	}, nil
}

// Close closes the database.
func (db *datastoreDB) Close() {
	// No op.
}

func (db *datastoreDB) PageId(id int64) *datastore.Key {
	return datastore.IDKey("Page", id, nil)
}

// GetBook retrieves a book by its ID.
func (db *datastoreDB) GetPage(id int64) (*Page, error) {
	ctx := context.Background()
	k := db.PageId(id)
	page := &Page{}
	if err := db.client.Get(ctx, k, page); err != nil {
		return nil, fmt.Errorf("datastoredb: could not get Page: %v", err)
	}
	page.Id = id
	return page, nil
}

func (db *datastoreDB) GetPageTemplateData(id int64) (*PageTemplateData, error) {
	var page, err = db.GetPage(id)
	if(err != nil) {
		return nil, err
	}
	var result = PageTemplateData {
		Page: page,
		ChildPages: []*Page{},
	}

	result.ChildPages, err = db.ListChildren(id)
	if(err != nil) {
		return nil, err
	}

	return &result, nil;
}

func (db *datastoreDB) AddPage(b *Page) (id int64, err error) {
	ctx := context.Background()
	k := datastore.IncompleteKey("Page", nil)
	k, err = db.client.Put(ctx, k, b)
	if err != nil {
		return 0, fmt.Errorf("datastoredb: could not put Page: %v", err)
	}
	return k.ID, nil
}

// ListBooks returns a list of books, ordered by title.
func (db *datastoreDB) ListChildren(id int64) ([]*Page, error) {
	ctx := context.Background()
	children := make([]*Page, 0)
	q := datastore.NewQuery("Page").Filter("ParentId =", db.PageId(id))

	keys, err := db.client.GetAll(ctx, q, &children)

	if err != nil {
		return nil, fmt.Errorf("datastoredb: could not list children: %v", err)
	}

	for i, k := range keys {
		children[i].Id = k.ID
	}

	return children, nil
}
