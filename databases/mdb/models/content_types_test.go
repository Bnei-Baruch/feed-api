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

func testContentTypes(t *testing.T) {
	t.Parallel()

	query := ContentTypes()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testContentTypesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ContentType{}
	if err = randomize.Struct(seed, o, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
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

	count, err := ContentTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testContentTypesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ContentType{}
	if err = randomize.Struct(seed, o, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := ContentTypes().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := ContentTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testContentTypesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ContentType{}
	if err = randomize.Struct(seed, o, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ContentTypeSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := ContentTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testContentTypesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ContentType{}
	if err = randomize.Struct(seed, o, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := ContentTypeExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if ContentType exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ContentTypeExists to return true, but got false.")
	}
}

func testContentTypesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ContentType{}
	if err = randomize.Struct(seed, o, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	contentTypeFound, err := FindContentType(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if contentTypeFound == nil {
		t.Error("want a record, got nil")
	}
}

func testContentTypesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ContentType{}
	if err = randomize.Struct(seed, o, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = ContentTypes().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testContentTypesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ContentType{}
	if err = randomize.Struct(seed, o, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := ContentTypes().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testContentTypesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentTypeOne := &ContentType{}
	contentTypeTwo := &ContentType{}
	if err = randomize.Struct(seed, contentTypeOne, contentTypeDBTypes, false, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}
	if err = randomize.Struct(seed, contentTypeTwo, contentTypeDBTypes, false, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = contentTypeOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = contentTypeTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := ContentTypes().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testContentTypesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	contentTypeOne := &ContentType{}
	contentTypeTwo := &ContentType{}
	if err = randomize.Struct(seed, contentTypeOne, contentTypeDBTypes, false, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}
	if err = randomize.Struct(seed, contentTypeTwo, contentTypeDBTypes, false, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = contentTypeOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = contentTypeTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ContentTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testContentTypesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ContentType{}
	if err = randomize.Struct(seed, o, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ContentTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testContentTypesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ContentType{}
	if err = randomize.Struct(seed, o, contentTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(contentTypeColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := ContentTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testContentTypeToManyTypeCollections(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a ContentType
	var b, c Collection

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, collectionDBTypes, false, collectionColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, collectionDBTypes, false, collectionColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	b.TypeID = a.ID
	c.TypeID = a.ID

	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.TypeCollections().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.TypeID == b.TypeID {
			bFound = true
		}
		if v.TypeID == c.TypeID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := ContentTypeSlice{&a}
	if err = a.L.LoadTypeCollections(ctx, tx, false, (*[]*ContentType)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.TypeCollections); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.TypeCollections = nil
	if err = a.L.LoadTypeCollections(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.TypeCollections); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testContentTypeToManyTypeContentUnits(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a ContentType
	var b, c ContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, contentUnitDBTypes, false, contentUnitColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, contentUnitDBTypes, false, contentUnitColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	b.TypeID = a.ID
	c.TypeID = a.ID

	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.TypeContentUnits().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.TypeID == b.TypeID {
			bFound = true
		}
		if v.TypeID == c.TypeID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := ContentTypeSlice{&a}
	if err = a.L.LoadTypeContentUnits(ctx, tx, false, (*[]*ContentType)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.TypeContentUnits); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.TypeContentUnits = nil
	if err = a.L.LoadTypeContentUnits(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.TypeContentUnits); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testContentTypeToManyAddOpTypeCollections(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a ContentType
	var b, c, d, e Collection

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentTypeDBTypes, false, strmangle.SetComplement(contentTypePrimaryKeyColumns, contentTypeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Collection{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, collectionDBTypes, false, strmangle.SetComplement(collectionPrimaryKeyColumns, collectionColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*Collection{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddTypeCollections(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.TypeID {
			t.Error("foreign key was wrong value", a.ID, first.TypeID)
		}
		if a.ID != second.TypeID {
			t.Error("foreign key was wrong value", a.ID, second.TypeID)
		}

		if first.R.Type != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Type != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.TypeCollections[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.TypeCollections[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.TypeCollections().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testContentTypeToManyAddOpTypeContentUnits(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a ContentType
	var b, c, d, e ContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentTypeDBTypes, false, strmangle.SetComplement(contentTypePrimaryKeyColumns, contentTypeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ContentUnit{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, contentUnitDBTypes, false, strmangle.SetComplement(contentUnitPrimaryKeyColumns, contentUnitColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*ContentUnit{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddTypeContentUnits(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.TypeID {
			t.Error("foreign key was wrong value", a.ID, first.TypeID)
		}
		if a.ID != second.TypeID {
			t.Error("foreign key was wrong value", a.ID, second.TypeID)
		}

		if first.R.Type != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Type != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.TypeContentUnits[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.TypeContentUnits[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.TypeContentUnits().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testContentTypesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ContentType{}
	if err = randomize.Struct(seed, o, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
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

func testContentTypesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ContentType{}
	if err = randomize.Struct(seed, o, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ContentTypeSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testContentTypesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ContentType{}
	if err = randomize.Struct(seed, o, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := ContentTypes().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	contentTypeDBTypes = map[string]string{`ID`: `bigint`, `Name`: `character varying`, `Description`: `character varying`}
	_                  = bytes.MinRead
)

func testContentTypesUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(contentTypePrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(contentTypeAllColumns) == len(contentTypePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &ContentType{}
	if err = randomize.Struct(seed, o, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ContentTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, contentTypeDBTypes, true, contentTypePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testContentTypesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(contentTypeAllColumns) == len(contentTypePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &ContentType{}
	if err = randomize.Struct(seed, o, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ContentTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, contentTypeDBTypes, true, contentTypePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(contentTypeAllColumns, contentTypePrimaryKeyColumns) {
		fields = contentTypeAllColumns
	} else {
		fields = strmangle.SetComplement(
			contentTypeAllColumns,
			contentTypePrimaryKeyColumns,
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

	slice := ContentTypeSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testContentTypesUpsert(t *testing.T) {
	t.Parallel()

	if len(contentTypeAllColumns) == len(contentTypePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := ContentType{}
	if err = randomize.Struct(seed, &o, contentTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert ContentType: %s", err)
	}

	count, err := ContentTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, contentTypeDBTypes, false, contentTypePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert ContentType: %s", err)
	}

	count, err = ContentTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
