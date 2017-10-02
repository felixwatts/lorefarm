// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package lorefarm

import "cloud.google.com/go/datastore"

type Page struct {
	Id            int64
	ParentId *datastore.Key
	Content       string
}

type PageTemplateData struct {
	Page *Page
	ChildPages []*Page;
}

// BookDatabase provides thread-safe access to a database of books.
type PageDatabase interface {

	PageId(id int64) *datastore.Key
	
	ListChildren(pageId int64) ([]*Page, error)

	GetPage(id int64) (*Page, error)

	GetPageTemplateData(id int64) (*PageTemplateData, error)

	AddPage(b *Page) (id int64, err error)

	// Close closes the database, freeing up any available resources.
	// TODO(cbro): Close() should return an error.
	Close()
}
