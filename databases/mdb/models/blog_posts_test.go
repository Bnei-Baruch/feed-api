// Code generated by SQLBoiler 3.7.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/randomize"
	"github.com/volatiletech/sqlboiler/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testBlogPosts(t *testing.T) {
	t.Parallel()

	query := BlogPosts()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testBlogPostsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BlogPost{}
	if err = randomize.Struct(seed, o, blogPostDBTypes, true, blogPostColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BlogPost struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := BlogPosts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testBlogPostsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BlogPost{}
	if err = randomize.Struct(seed, o, blogPostDBTypes, true, blogPostColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BlogPost struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := BlogPosts().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := BlogPosts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testBlogPostsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BlogPost{}
	if err = randomize.Struct(seed, o, blogPostDBTypes, true, blogPostColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BlogPost struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := BlogPostSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := BlogPosts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testBlogPostsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BlogPost{}
	if err = randomize.Struct(seed, o, blogPostDBTypes, true, blogPostColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BlogPost struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := BlogPostExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if BlogPost exists: %s", err)
	}
	if !e {
		t.Errorf("Expected BlogPostExists to return true, but got false.")
	}
}

func testBlogPostsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BlogPost{}
	if err = randomize.Struct(seed, o, blogPostDBTypes, true, blogPostColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BlogPost struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	blogPostFound, err := FindBlogPost(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if blogPostFound == nil {
		t.Error("want a record, got nil")
	}
}

func testBlogPostsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BlogPost{}
	if err = randomize.Struct(seed, o, blogPostDBTypes, true, blogPostColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BlogPost struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = BlogPosts().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testBlogPostsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BlogPost{}
	if err = randomize.Struct(seed, o, blogPostDBTypes, true, blogPostColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BlogPost struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := BlogPosts().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testBlogPostsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	blogPostOne := &BlogPost{}
	blogPostTwo := &BlogPost{}
	if err = randomize.Struct(seed, blogPostOne, blogPostDBTypes, false, blogPostColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BlogPost struct: %s", err)
	}
	if err = randomize.Struct(seed, blogPostTwo, blogPostDBTypes, false, blogPostColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BlogPost struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = blogPostOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = blogPostTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := BlogPosts().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testBlogPostsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	blogPostOne := &BlogPost{}
	blogPostTwo := &BlogPost{}
	if err = randomize.Struct(seed, blogPostOne, blogPostDBTypes, false, blogPostColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BlogPost struct: %s", err)
	}
	if err = randomize.Struct(seed, blogPostTwo, blogPostDBTypes, false, blogPostColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BlogPost struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = blogPostOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = blogPostTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := BlogPosts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testBlogPostsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BlogPost{}
	if err = randomize.Struct(seed, o, blogPostDBTypes, true, blogPostColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BlogPost struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := BlogPosts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testBlogPostsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BlogPost{}
	if err = randomize.Struct(seed, o, blogPostDBTypes, true); err != nil {
		t.Errorf("Unable to randomize BlogPost struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(blogPostColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := BlogPosts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testBlogPostToOneBlogUsingBlog(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local BlogPost
	var foreign Blog

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, blogPostDBTypes, false, blogPostColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BlogPost struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, blogDBTypes, false, blogColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Blog struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.BlogID = foreign.ID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Blog().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := BlogPostSlice{&local}
	if err = local.L.LoadBlog(ctx, tx, false, (*[]*BlogPost)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Blog == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Blog = nil
	if err = local.L.LoadBlog(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Blog == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testBlogPostToOneSetOpBlogUsingBlog(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a BlogPost
	var b, c Blog

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, blogPostDBTypes, false, strmangle.SetComplement(blogPostPrimaryKeyColumns, blogPostColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, blogDBTypes, false, strmangle.SetComplement(blogPrimaryKeyColumns, blogColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, blogDBTypes, false, strmangle.SetComplement(blogPrimaryKeyColumns, blogColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Blog{&b, &c} {
		err = a.SetBlog(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Blog != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.BlogPosts[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.BlogID != x.ID {
			t.Error("foreign key was wrong value", a.BlogID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.BlogID))
		reflect.Indirect(reflect.ValueOf(&a.BlogID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.BlogID != x.ID {
			t.Error("foreign key was wrong value", a.BlogID, x.ID)
		}
	}
}

func testBlogPostsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BlogPost{}
	if err = randomize.Struct(seed, o, blogPostDBTypes, true, blogPostColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BlogPost struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testBlogPostsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BlogPost{}
	if err = randomize.Struct(seed, o, blogPostDBTypes, true, blogPostColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BlogPost struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := BlogPostSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testBlogPostsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BlogPost{}
	if err = randomize.Struct(seed, o, blogPostDBTypes, true, blogPostColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BlogPost struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := BlogPosts().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	blogPostDBTypes = map[string]string{`ID`: `bigint`, `BlogID`: `bigint`, `WPID`: `bigint`, `Title`: `text`, `Content`: `text`, `PostedAt`: `timestamp without time zone`, `CreatedAt`: `timestamp with time zone`, `Link`: `character varying`, `Filtered`: `boolean`}
	_               = bytes.MinRead
)

func testBlogPostsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(blogPostPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(blogPostAllColumns) == len(blogPostPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &BlogPost{}
	if err = randomize.Struct(seed, o, blogPostDBTypes, true, blogPostColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BlogPost struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := BlogPosts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, blogPostDBTypes, true, blogPostPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize BlogPost struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testBlogPostsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(blogPostAllColumns) == len(blogPostPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &BlogPost{}
	if err = randomize.Struct(seed, o, blogPostDBTypes, true, blogPostColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BlogPost struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := BlogPosts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, blogPostDBTypes, true, blogPostPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize BlogPost struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(blogPostAllColumns, blogPostPrimaryKeyColumns) {
		fields = blogPostAllColumns
	} else {
		fields = strmangle.SetComplement(
			blogPostAllColumns,
			blogPostPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := BlogPostSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testBlogPostsUpsert(t *testing.T) {
	t.Parallel()

	if len(blogPostAllColumns) == len(blogPostPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := BlogPost{}
	if err = randomize.Struct(seed, &o, blogPostDBTypes, true); err != nil {
		t.Errorf("Unable to randomize BlogPost struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert BlogPost: %s", err)
	}

	count, err := BlogPosts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, blogPostDBTypes, false, blogPostPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize BlogPost struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert BlogPost: %s", err)
	}

	count, err = BlogPosts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}