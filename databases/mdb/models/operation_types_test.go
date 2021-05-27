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

func testOperationTypes(t *testing.T) {
	t.Parallel()

	query := OperationTypes()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testOperationTypesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &OperationType{}
	if err = randomize.Struct(seed, o, operationTypeDBTypes, true, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
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

	count, err := OperationTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testOperationTypesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &OperationType{}
	if err = randomize.Struct(seed, o, operationTypeDBTypes, true, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := OperationTypes().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := OperationTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testOperationTypesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &OperationType{}
	if err = randomize.Struct(seed, o, operationTypeDBTypes, true, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := OperationTypeSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := OperationTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testOperationTypesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &OperationType{}
	if err = randomize.Struct(seed, o, operationTypeDBTypes, true, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := OperationTypeExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if OperationType exists: %s", err)
	}
	if !e {
		t.Errorf("Expected OperationTypeExists to return true, but got false.")
	}
}

func testOperationTypesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &OperationType{}
	if err = randomize.Struct(seed, o, operationTypeDBTypes, true, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	operationTypeFound, err := FindOperationType(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if operationTypeFound == nil {
		t.Error("want a record, got nil")
	}
}

func testOperationTypesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &OperationType{}
	if err = randomize.Struct(seed, o, operationTypeDBTypes, true, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = OperationTypes().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testOperationTypesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &OperationType{}
	if err = randomize.Struct(seed, o, operationTypeDBTypes, true, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := OperationTypes().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testOperationTypesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operationTypeOne := &OperationType{}
	operationTypeTwo := &OperationType{}
	if err = randomize.Struct(seed, operationTypeOne, operationTypeDBTypes, false, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}
	if err = randomize.Struct(seed, operationTypeTwo, operationTypeDBTypes, false, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = operationTypeOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = operationTypeTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := OperationTypes().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testOperationTypesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	operationTypeOne := &OperationType{}
	operationTypeTwo := &OperationType{}
	if err = randomize.Struct(seed, operationTypeOne, operationTypeDBTypes, false, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}
	if err = randomize.Struct(seed, operationTypeTwo, operationTypeDBTypes, false, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = operationTypeOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = operationTypeTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := OperationTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testOperationTypesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &OperationType{}
	if err = randomize.Struct(seed, o, operationTypeDBTypes, true, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := OperationTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testOperationTypesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &OperationType{}
	if err = randomize.Struct(seed, o, operationTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(operationTypeColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := OperationTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testOperationTypeToManyTypeOperations(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a OperationType
	var b, c Operation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, operationTypeDBTypes, true, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, operationDBTypes, false, operationColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, operationDBTypes, false, operationColumnsWithDefault...); err != nil {
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

	check, err := a.TypeOperations().All(ctx, tx)
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

	slice := OperationTypeSlice{&a}
	if err = a.L.LoadTypeOperations(ctx, tx, false, (*[]*OperationType)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.TypeOperations); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.TypeOperations = nil
	if err = a.L.LoadTypeOperations(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.TypeOperations); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testOperationTypeToManyAddOpTypeOperations(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a OperationType
	var b, c, d, e Operation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, operationTypeDBTypes, false, strmangle.SetComplement(operationTypePrimaryKeyColumns, operationTypeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Operation{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, operationDBTypes, false, strmangle.SetComplement(operationPrimaryKeyColumns, operationColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*Operation{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddTypeOperations(ctx, tx, i != 0, x...)
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

		if a.R.TypeOperations[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.TypeOperations[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.TypeOperations().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testOperationTypesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &OperationType{}
	if err = randomize.Struct(seed, o, operationTypeDBTypes, true, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
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

func testOperationTypesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &OperationType{}
	if err = randomize.Struct(seed, o, operationTypeDBTypes, true, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := OperationTypeSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testOperationTypesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &OperationType{}
	if err = randomize.Struct(seed, o, operationTypeDBTypes, true, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := OperationTypes().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	operationTypeDBTypes = map[string]string{`ID`: `bigint`, `Name`: `character varying`, `Description`: `character varying`}
	_                    = bytes.MinRead
)

func testOperationTypesUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(operationTypePrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(operationTypeAllColumns) == len(operationTypePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &OperationType{}
	if err = randomize.Struct(seed, o, operationTypeDBTypes, true, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := OperationTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, operationTypeDBTypes, true, operationTypePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testOperationTypesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(operationTypeAllColumns) == len(operationTypePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &OperationType{}
	if err = randomize.Struct(seed, o, operationTypeDBTypes, true, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := OperationTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, operationTypeDBTypes, true, operationTypePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(operationTypeAllColumns, operationTypePrimaryKeyColumns) {
		fields = operationTypeAllColumns
	} else {
		fields = strmangle.SetComplement(
			operationTypeAllColumns,
			operationTypePrimaryKeyColumns,
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

	slice := OperationTypeSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testOperationTypesUpsert(t *testing.T) {
	t.Parallel()

	if len(operationTypeAllColumns) == len(operationTypePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := OperationType{}
	if err = randomize.Struct(seed, &o, operationTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert OperationType: %s", err)
	}

	count, err := OperationTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, operationTypeDBTypes, false, operationTypePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert OperationType: %s", err)
	}

	count, err = OperationTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
