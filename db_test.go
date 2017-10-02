// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package lorefarm

import (
	"testing"

	"cloud.google.com/go/datastore"

	//"cloud.google.com/go/datastore"

	//"golang.org/x/net/context"

	//"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func testDB(t *testing.T, db PageDatabase) {
	defer db.Close()

	b := &Page{
		ParentId: datastore.IDKey("Page", 42, nil),
		Content: "miaow",
	}

	id, err := db.AddPage(b)
	if err != nil {
		t.Fatal(err)
	}

	gotPage, err := db.GetPage(id)
	if err != nil {
		t.Error(err)
	}
	if got, want := gotPage.ParentId, b.ParentId; got != want {
		t.Errorf("page.ParentId description: got %q, want %q", got, want)
	}
	if got, want := gotPage.Content, b.Content; got != want {
		t.Errorf("page.Content description: got %q, want %q", got, want)
	}

	if _, err := db.GetPage(1000); err == nil {
		t.Error("want non-nil err")
	}
}

// func TestDatastoreDB(t *testing.T) {
// 	tc := testutil.SystemTest(t)
// 	ctx := context.Background()

// 	client, err := datastore.NewClient(ctx, tc.ProjectID)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer client.Close()

// 	db, err := newDatastoreDB(client)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	testDB(t, db)
// }
