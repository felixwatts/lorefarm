// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package lorefarm

type Page struct {
	Id            int64
	ParentId int64
	ChildIds []int64
	Content       string
}

// BookDatabase provides thread-safe access to a database of books.
type PageDatabase interface {
	
	ListChildren(pageId int64) ([]*Page, error)

	GetPage(id int64) (*Page, error)

	AddPage(b *Page) (id int64, err error)

	// Close closes the database, freeing up any available resources.
	// TODO(cbro): Close() should return an error.
	Close()
}
