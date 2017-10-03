// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package lorefarm

const ROOT_PAGE_ID int64 = 5722646637445120

type PageData struct {
	Id            int64
	ParentId int64
	Content       string
}

type Page struct {
	Current *PageData
	Next []*PageData
}

type PageDatabase interface {
	GetPage(id int64) (*Page, error)
	AddPage(b *PageData) (id int64, err error)
}
