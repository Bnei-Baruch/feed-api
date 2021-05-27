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

func testOperations(t *testing.T) {
	t.Parallel()

	query := Operations()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testOperationsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Operation{}
	if err = randomize.Struct(seed, o, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
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

	count, err := Operations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testOperationsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Operation{}
	if err = randomize.Struct(seed, o, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := Operations().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Operations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testOperationsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Operation{}
	if err = randomize.Struct(seed, o, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := OperationSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Operations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testOperationsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Operation{}
	if err = randomize.Struct(seed, o, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := OperationExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if Operation exists: %s", err)
	}
	if !e {
		t.Errorf("Expected OperationExists to return true, but got false.")
	}
}

func testOperationsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Operation{}
	if err = randomize.Struct(seed, o, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	operationFound, err := FindOperation(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if operationFound == nil {
		t.Error("want a record, got nil")
	}
}

func testOperationsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Operation{}
	if err = randomize.Struct(seed, o, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = Operations().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testOperationsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Operation{}
	if err = randomize.Struct(seed, o, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := Operations().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testOperationsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operationOne := &Operation{}
	operationTwo := &Operation{}
	if err = randomize.Struct(seed, operationOne, operationDBTypes, false, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}
	if err = randomize.Struct(seed, operationTwo, operationDBTypes, false, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = operationOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = operationTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Operations().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testOperationsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	operationOne := &Operation{}
	operationTwo := &Operation{}
	if err = randomize.Struct(seed, operationOne, operationDBTypes, false, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}
	if err = randomize.Struct(seed, operationTwo, operationDBTypes, false, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = operationOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = operationTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Operations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testOperationsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Operation{}
	if err = randomize.Struct(seed, o, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Operations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testOperationsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Operation{}
	if err = randomize.Struct(seed, o, operationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(operationColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := Operations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testOperationToManyFiles(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Operation
	var b, c File

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, fileDBTypes, false, fileColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, fileDBTypes, false, fileColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	_, err = tx.Exec("insert into \"files_operations\" (\"operation_id\", \"file_id\") values ($1, $2)", a.ID, b.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec("insert into \"files_operations\" (\"operation_id\", \"file_id\") values ($1, $2)", a.ID, c.ID)
	if err != nil {
		t.Fatal(err)
	}

	check, err := a.Files().All(ctx, tx)
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

	slice := OperationSlice{&a}
	if err = a.L.LoadFiles(ctx, tx, false, (*[]*Operation)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Files); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Files = nil
	if err = a.L.LoadFiles(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Files); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testOperationToManyAddOpFiles(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Operation
	var b, c, d, e File

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, operationDBTypes, false, strmangle.SetComplement(operationPrimaryKeyColumns, operationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*File{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*File{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddFiles(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if first.R.Operations[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}
		if second.R.Operations[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}

		if a.R.Files[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Files[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Files().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testOperationToManySetOpFiles(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Operation
	var b, c, d, e File

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, operationDBTypes, false, strmangle.SetComplement(operationPrimaryKeyColumns, operationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*File{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
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

	err = a.SetFiles(ctx, tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Files().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetFiles(ctx, tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Files().Count(ctx, tx)
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
	// if len(b.R.Operations) != 0 {
	// 	t.Error("relationship was not removed properly from the slice")
	// }
	// if len(c.R.Operations) != 0 {
	// 	t.Error("relationship was not removed properly from the slice")
	// }
	if d.R.Operations[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}
	if e.R.Operations[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}

	if a.R.Files[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.Files[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testOperationToManyRemoveOpFiles(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Operation
	var b, c, d, e File

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, operationDBTypes, false, strmangle.SetComplement(operationPrimaryKeyColumns, operationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*File{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	err = a.AddFiles(ctx, tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Files().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveFiles(ctx, tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Files().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if len(b.R.Operations) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.Operations) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.Operations[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Operations[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if len(a.R.Files) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.Files[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.Files[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testOperationToOneOperationTypeUsingType(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local Operation
	var foreign OperationType

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, operationDBTypes, false, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, operationTypeDBTypes, false, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.TypeID = foreign.ID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Type().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := OperationSlice{&local}
	if err = local.L.LoadType(ctx, tx, false, (*[]*Operation)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Type == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Type = nil
	if err = local.L.LoadType(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Type == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testOperationToOneUserUsingUser(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local Operation
	var foreign User

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, userDBTypes, false, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	queries.Assign(&local.UserID, foreign.ID)
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.User().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if !queries.Equal(check.ID, foreign.ID) {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := OperationSlice{&local}
	if err = local.L.LoadUser(ctx, tx, false, (*[]*Operation)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.User == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.User = nil
	if err = local.L.LoadUser(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.User == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testOperationToOneSetOpOperationTypeUsingType(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Operation
	var b, c OperationType

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, operationDBTypes, false, strmangle.SetComplement(operationPrimaryKeyColumns, operationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, operationTypeDBTypes, false, strmangle.SetComplement(operationTypePrimaryKeyColumns, operationTypeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, operationTypeDBTypes, false, strmangle.SetComplement(operationTypePrimaryKeyColumns, operationTypeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*OperationType{&b, &c} {
		err = a.SetType(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Type != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.TypeOperations[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.TypeID != x.ID {
			t.Error("foreign key was wrong value", a.TypeID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.TypeID))
		reflect.Indirect(reflect.ValueOf(&a.TypeID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.TypeID != x.ID {
			t.Error("foreign key was wrong value", a.TypeID, x.ID)
		}
	}
}
func testOperationToOneSetOpUserUsingUser(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Operation
	var b, c User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, operationDBTypes, false, strmangle.SetComplement(operationPrimaryKeyColumns, operationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*User{&b, &c} {
		err = a.SetUser(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.User != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.Operations[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if !queries.Equal(a.UserID, x.ID) {
			t.Error("foreign key was wrong value", a.UserID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.UserID))
		reflect.Indirect(reflect.ValueOf(&a.UserID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if !queries.Equal(a.UserID, x.ID) {
			t.Error("foreign key was wrong value", a.UserID, x.ID)
		}
	}
}

func testOperationToOneRemoveOpUserUsingUser(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Operation
	var b User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, operationDBTypes, false, strmangle.SetComplement(operationPrimaryKeyColumns, operationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = a.SetUser(ctx, tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveUser(ctx, tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.User().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.User != nil {
		t.Error("R struct entry should be nil")
	}

	if !queries.IsValuerNil(a.UserID) {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.Operations) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testOperationsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Operation{}
	if err = randomize.Struct(seed, o, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
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

func testOperationsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Operation{}
	if err = randomize.Struct(seed, o, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := OperationSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testOperationsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Operation{}
	if err = randomize.Struct(seed, o, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Operations().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	operationDBTypes = map[string]string{`ID`: `bigint`, `UID`: `character`, `TypeID`: `bigint`, `CreatedAt`: `timestamp with time zone`, `Station`: `character varying`, `UserID`: `bigint`, `Details`: `character varying`, `Properties`: `jsonb`}
	_                = bytes.MinRead
)

func testOperationsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(operationPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(operationAllColumns) == len(operationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Operation{}
	if err = randomize.Struct(seed, o, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Operations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, operationDBTypes, true, operationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testOperationsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(operationAllColumns) == len(operationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Operation{}
	if err = randomize.Struct(seed, o, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Operations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, operationDBTypes, true, operationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(operationAllColumns, operationPrimaryKeyColumns) {
		fields = operationAllColumns
	} else {
		fields = strmangle.SetComplement(
			operationAllColumns,
			operationPrimaryKeyColumns,
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

	slice := OperationSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testOperationsUpsert(t *testing.T) {
	t.Parallel()

	if len(operationAllColumns) == len(operationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := Operation{}
	if err = randomize.Struct(seed, &o, operationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Operation: %s", err)
	}

	count, err := Operations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, operationDBTypes, false, operationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Operation: %s", err)
	}

	count, err = Operations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
