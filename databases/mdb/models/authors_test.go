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

func testAuthors(t *testing.T) {
	t.Parallel()

	query := Authors()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testAuthorsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Author{}
	if err = randomize.Struct(seed, o, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
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

	count, err := Authors().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testAuthorsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Author{}
	if err = randomize.Struct(seed, o, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := Authors().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Authors().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testAuthorsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Author{}
	if err = randomize.Struct(seed, o, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := AuthorSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Authors().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testAuthorsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Author{}
	if err = randomize.Struct(seed, o, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := AuthorExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if Author exists: %s", err)
	}
	if !e {
		t.Errorf("Expected AuthorExists to return true, but got false.")
	}
}

func testAuthorsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Author{}
	if err = randomize.Struct(seed, o, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	authorFound, err := FindAuthor(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if authorFound == nil {
		t.Error("want a record, got nil")
	}
}

func testAuthorsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Author{}
	if err = randomize.Struct(seed, o, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = Authors().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testAuthorsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Author{}
	if err = randomize.Struct(seed, o, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := Authors().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testAuthorsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	authorOne := &Author{}
	authorTwo := &Author{}
	if err = randomize.Struct(seed, authorOne, authorDBTypes, false, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}
	if err = randomize.Struct(seed, authorTwo, authorDBTypes, false, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = authorOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = authorTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Authors().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testAuthorsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	authorOne := &Author{}
	authorTwo := &Author{}
	if err = randomize.Struct(seed, authorOne, authorDBTypes, false, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}
	if err = randomize.Struct(seed, authorTwo, authorDBTypes, false, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = authorOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = authorTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Authors().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testAuthorsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Author{}
	if err = randomize.Struct(seed, o, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Authors().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testAuthorsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Author{}
	if err = randomize.Struct(seed, o, authorDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(authorColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := Authors().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testAuthorToManyAuthorI18ns(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Author
	var b, c AuthorI18n

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, authorI18nDBTypes, false, authorI18nColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, authorI18nDBTypes, false, authorI18nColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	b.AuthorID = a.ID
	c.AuthorID = a.ID

	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.AuthorI18ns().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.AuthorID == b.AuthorID {
			bFound = true
		}
		if v.AuthorID == c.AuthorID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := AuthorSlice{&a}
	if err = a.L.LoadAuthorI18ns(ctx, tx, false, (*[]*Author)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.AuthorI18ns); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.AuthorI18ns = nil
	if err = a.L.LoadAuthorI18ns(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.AuthorI18ns); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testAuthorToManySources(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Author
	var b, c Source

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, sourceDBTypes, false, sourceColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, sourceDBTypes, false, sourceColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	_, err = tx.Exec("insert into \"authors_sources\" (\"author_id\", \"source_id\") values ($1, $2)", a.ID, b.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec("insert into \"authors_sources\" (\"author_id\", \"source_id\") values ($1, $2)", a.ID, c.ID)
	if err != nil {
		t.Fatal(err)
	}

	check, err := a.Sources().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.ID == b.ID {
			bFound = true
		}
		if v.ID == c.ID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := AuthorSlice{&a}
	if err = a.L.LoadSources(ctx, tx, false, (*[]*Author)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Sources); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Sources = nil
	if err = a.L.LoadSources(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Sources); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testAuthorToManyAddOpAuthorI18ns(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Author
	var b, c, d, e AuthorI18n

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, authorDBTypes, false, strmangle.SetComplement(authorPrimaryKeyColumns, authorColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*AuthorI18n{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, authorI18nDBTypes, false, strmangle.SetComplement(authorI18nPrimaryKeyColumns, authorI18nColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*AuthorI18n{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddAuthorI18ns(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.AuthorID {
			t.Error("foreign key was wrong value", a.ID, first.AuthorID)
		}
		if a.ID != second.AuthorID {
			t.Error("foreign key was wrong value", a.ID, second.AuthorID)
		}

		if first.R.Author != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Author != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.AuthorI18ns[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.AuthorI18ns[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.AuthorI18ns().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testAuthorToManyAddOpSources(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Author
	var b, c, d, e Source

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, authorDBTypes, false, strmangle.SetComplement(authorPrimaryKeyColumns, authorColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Source{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*Source{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddSources(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if first.R.Authors[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}
		if second.R.Authors[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}

		if a.R.Sources[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Sources[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Sources().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testAuthorToManySetOpSources(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Author
	var b, c, d, e Source

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, authorDBTypes, false, strmangle.SetComplement(authorPrimaryKeyColumns, authorColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Source{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err = a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	err = a.SetSources(ctx, tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Sources().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetSources(ctx, tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Sources().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	// The following checks cannot be implemented since we have no handle
	// to these when we call Set(). Leaving them here as wishful thinking
	// and to let people know there's dragons.
	//
	// if len(b.R.Authors) != 0 {
	// 	t.Error("relationship was not removed properly from the slice")
	// }
	// if len(c.R.Authors) != 0 {
	// 	t.Error("relationship was not removed properly from the slice")
	// }
	if d.R.Authors[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}
	if e.R.Authors[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}

	if a.R.Sources[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.Sources[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testAuthorToManyRemoveOpSources(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Author
	var b, c, d, e Source

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, authorDBTypes, false, strmangle.SetComplement(authorPrimaryKeyColumns, authorColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Source{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	err = a.AddSources(ctx, tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Sources().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveSources(ctx, tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Sources().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if len(b.R.Authors) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.Authors) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.Authors[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Authors[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if len(a.R.Sources) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.Sources[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.Sources[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testAuthorsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Author{}
	if err = randomize.Struct(seed, o, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
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

func testAuthorsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Author{}
	if err = randomize.Struct(seed, o, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := AuthorSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testAuthorsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Author{}
	if err = randomize.Struct(seed, o, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Authors().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	authorDBTypes = map[string]string{`ID`: `bigint`, `Code`: `character`, `Name`: `character varying`, `FullName`: `character varying`, `CreatedAt`: `timestamp with time zone`}
	_             = bytes.MinRead
)

func testAuthorsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(authorPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(authorAllColumns) == len(authorPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Author{}
	if err = randomize.Struct(seed, o, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Authors().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, authorDBTypes, true, authorPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testAuthorsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(authorAllColumns) == len(authorPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Author{}
	if err = randomize.Struct(seed, o, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Authors().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, authorDBTypes, true, authorPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(authorAllColumns, authorPrimaryKeyColumns) {
		fields = authorAllColumns
	} else {
		fields = strmangle.SetComplement(
			authorAllColumns,
			authorPrimaryKeyColumns,
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

	slice := AuthorSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testAuthorsUpsert(t *testing.T) {
	t.Parallel()

	if len(authorAllColumns) == len(authorPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := Author{}
	if err = randomize.Struct(seed, &o, authorDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Author: %s", err)
	}

	count, err := Authors().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, authorDBTypes, false, authorPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Author: %s", err)
	}

	count, err = Authors().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}