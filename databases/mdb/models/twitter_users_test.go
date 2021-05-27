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

func testTwitterUsers(t *testing.T) {
	t.Parallel()

	query := TwitterUsers()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testTwitterUsersDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TwitterUser{}
	if err = randomize.Struct(seed, o, twitterUserDBTypes, true, twitterUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterUser struct: %s", err)
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

	count, err := TwitterUsers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testTwitterUsersQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TwitterUser{}
	if err = randomize.Struct(seed, o, twitterUserDBTypes, true, twitterUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := TwitterUsers().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := TwitterUsers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testTwitterUsersSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TwitterUser{}
	if err = randomize.Struct(seed, o, twitterUserDBTypes, true, twitterUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := TwitterUserSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := TwitterUsers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testTwitterUsersExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TwitterUser{}
	if err = randomize.Struct(seed, o, twitterUserDBTypes, true, twitterUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := TwitterUserExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if TwitterUser exists: %s", err)
	}
	if !e {
		t.Errorf("Expected TwitterUserExists to return true, but got false.")
	}
}

func testTwitterUsersFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TwitterUser{}
	if err = randomize.Struct(seed, o, twitterUserDBTypes, true, twitterUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	twitterUserFound, err := FindTwitterUser(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if twitterUserFound == nil {
		t.Error("want a record, got nil")
	}
}

func testTwitterUsersBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TwitterUser{}
	if err = randomize.Struct(seed, o, twitterUserDBTypes, true, twitterUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = TwitterUsers().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testTwitterUsersOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TwitterUser{}
	if err = randomize.Struct(seed, o, twitterUserDBTypes, true, twitterUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := TwitterUsers().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testTwitterUsersAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	twitterUserOne := &TwitterUser{}
	twitterUserTwo := &TwitterUser{}
	if err = randomize.Struct(seed, twitterUserOne, twitterUserDBTypes, false, twitterUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterUser struct: %s", err)
	}
	if err = randomize.Struct(seed, twitterUserTwo, twitterUserDBTypes, false, twitterUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = twitterUserOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = twitterUserTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := TwitterUsers().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testTwitterUsersCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	twitterUserOne := &TwitterUser{}
	twitterUserTwo := &TwitterUser{}
	if err = randomize.Struct(seed, twitterUserOne, twitterUserDBTypes, false, twitterUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterUser struct: %s", err)
	}
	if err = randomize.Struct(seed, twitterUserTwo, twitterUserDBTypes, false, twitterUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = twitterUserOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = twitterUserTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := TwitterUsers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testTwitterUsersInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TwitterUser{}
	if err = randomize.Struct(seed, o, twitterUserDBTypes, true, twitterUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := TwitterUsers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testTwitterUsersInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TwitterUser{}
	if err = randomize.Struct(seed, o, twitterUserDBTypes, true); err != nil {
		t.Errorf("Unable to randomize TwitterUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(twitterUserColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := TwitterUsers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testTwitterUserToManyUserTwitterTweets(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a TwitterUser
	var b, c TwitterTweet

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, twitterUserDBTypes, true, twitterUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterUser struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, twitterTweetDBTypes, false, twitterTweetColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, twitterTweetDBTypes, false, twitterTweetColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	b.UserID = a.ID
	c.UserID = a.ID

	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.UserTwitterTweets().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.UserID == b.UserID {
			bFound = true
		}
		if v.UserID == c.UserID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := TwitterUserSlice{&a}
	if err = a.L.LoadUserTwitterTweets(ctx, tx, false, (*[]*TwitterUser)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.UserTwitterTweets); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.UserTwitterTweets = nil
	if err = a.L.LoadUserTwitterTweets(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.UserTwitterTweets); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testTwitterUserToManyAddOpUserTwitterTweets(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a TwitterUser
	var b, c, d, e TwitterTweet

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, twitterUserDBTypes, false, strmangle.SetComplement(twitterUserPrimaryKeyColumns, twitterUserColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*TwitterTweet{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, twitterTweetDBTypes, false, strmangle.SetComplement(twitterTweetPrimaryKeyColumns, twitterTweetColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*TwitterTweet{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddUserTwitterTweets(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.UserID {
			t.Error("foreign key was wrong value", a.ID, first.UserID)
		}
		if a.ID != second.UserID {
			t.Error("foreign key was wrong value", a.ID, second.UserID)
		}

		if first.R.User != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.User != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.UserTwitterTweets[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.UserTwitterTweets[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.UserTwitterTweets().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testTwitterUsersReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TwitterUser{}
	if err = randomize.Struct(seed, o, twitterUserDBTypes, true, twitterUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterUser struct: %s", err)
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

func testTwitterUsersReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TwitterUser{}
	if err = randomize.Struct(seed, o, twitterUserDBTypes, true, twitterUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := TwitterUserSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testTwitterUsersSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TwitterUser{}
	if err = randomize.Struct(seed, o, twitterUserDBTypes, true, twitterUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := TwitterUsers().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	twitterUserDBTypes = map[string]string{`ID`: `bigint`, `Username`: `character varying`, `AccountID`: `character varying`, `DisplayName`: `character varying`}
	_                  = bytes.MinRead
)

func testTwitterUsersUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(twitterUserPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(twitterUserAllColumns) == len(twitterUserPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &TwitterUser{}
	if err = randomize.Struct(seed, o, twitterUserDBTypes, true, twitterUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := TwitterUsers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, twitterUserDBTypes, true, twitterUserPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize TwitterUser struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testTwitterUsersSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(twitterUserAllColumns) == len(twitterUserPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &TwitterUser{}
	if err = randomize.Struct(seed, o, twitterUserDBTypes, true, twitterUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := TwitterUsers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, twitterUserDBTypes, true, twitterUserPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize TwitterUser struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(twitterUserAllColumns, twitterUserPrimaryKeyColumns) {
		fields = twitterUserAllColumns
	} else {
		fields = strmangle.SetComplement(
			twitterUserAllColumns,
			twitterUserPrimaryKeyColumns,
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

	slice := TwitterUserSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testTwitterUsersUpsert(t *testing.T) {
	t.Parallel()

	if len(twitterUserAllColumns) == len(twitterUserPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := TwitterUser{}
	if err = randomize.Struct(seed, &o, twitterUserDBTypes, true); err != nil {
		t.Errorf("Unable to randomize TwitterUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert TwitterUser: %s", err)
	}

	count, err := TwitterUsers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, twitterUserDBTypes, false, twitterUserPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize TwitterUser struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert TwitterUser: %s", err)
	}

	count, err = TwitterUsers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
