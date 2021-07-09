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

func testTwitterTweets(t *testing.T) {
	t.Parallel()

	query := TwitterTweets()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testTwitterTweetsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TwitterTweet{}
	if err = randomize.Struct(seed, o, twitterTweetDBTypes, true, twitterTweetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterTweet struct: %s", err)
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

	count, err := TwitterTweets().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testTwitterTweetsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TwitterTweet{}
	if err = randomize.Struct(seed, o, twitterTweetDBTypes, true, twitterTweetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterTweet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := TwitterTweets().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := TwitterTweets().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testTwitterTweetsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TwitterTweet{}
	if err = randomize.Struct(seed, o, twitterTweetDBTypes, true, twitterTweetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterTweet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := TwitterTweetSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := TwitterTweets().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testTwitterTweetsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TwitterTweet{}
	if err = randomize.Struct(seed, o, twitterTweetDBTypes, true, twitterTweetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterTweet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := TwitterTweetExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if TwitterTweet exists: %s", err)
	}
	if !e {
		t.Errorf("Expected TwitterTweetExists to return true, but got false.")
	}
}

func testTwitterTweetsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TwitterTweet{}
	if err = randomize.Struct(seed, o, twitterTweetDBTypes, true, twitterTweetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterTweet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	twitterTweetFound, err := FindTwitterTweet(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if twitterTweetFound == nil {
		t.Error("want a record, got nil")
	}
}

func testTwitterTweetsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TwitterTweet{}
	if err = randomize.Struct(seed, o, twitterTweetDBTypes, true, twitterTweetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterTweet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = TwitterTweets().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testTwitterTweetsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TwitterTweet{}
	if err = randomize.Struct(seed, o, twitterTweetDBTypes, true, twitterTweetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterTweet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := TwitterTweets().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testTwitterTweetsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	twitterTweetOne := &TwitterTweet{}
	twitterTweetTwo := &TwitterTweet{}
	if err = randomize.Struct(seed, twitterTweetOne, twitterTweetDBTypes, false, twitterTweetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterTweet struct: %s", err)
	}
	if err = randomize.Struct(seed, twitterTweetTwo, twitterTweetDBTypes, false, twitterTweetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterTweet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = twitterTweetOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = twitterTweetTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := TwitterTweets().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testTwitterTweetsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	twitterTweetOne := &TwitterTweet{}
	twitterTweetTwo := &TwitterTweet{}
	if err = randomize.Struct(seed, twitterTweetOne, twitterTweetDBTypes, false, twitterTweetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterTweet struct: %s", err)
	}
	if err = randomize.Struct(seed, twitterTweetTwo, twitterTweetDBTypes, false, twitterTweetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterTweet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = twitterTweetOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = twitterTweetTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := TwitterTweets().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testTwitterTweetsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TwitterTweet{}
	if err = randomize.Struct(seed, o, twitterTweetDBTypes, true, twitterTweetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterTweet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := TwitterTweets().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testTwitterTweetsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TwitterTweet{}
	if err = randomize.Struct(seed, o, twitterTweetDBTypes, true); err != nil {
		t.Errorf("Unable to randomize TwitterTweet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(twitterTweetColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := TwitterTweets().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testTwitterTweetToOneTwitterUserUsingUser(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local TwitterTweet
	var foreign TwitterUser

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, twitterTweetDBTypes, false, twitterTweetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterTweet struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, twitterUserDBTypes, false, twitterUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterUser struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.UserID = foreign.ID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.User().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := TwitterTweetSlice{&local}
	if err = local.L.LoadUser(ctx, tx, false, (*[]*TwitterTweet)(&slice), nil); err != nil {
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

func testTwitterTweetToOneSetOpTwitterUserUsingUser(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a TwitterTweet
	var b, c TwitterUser

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, twitterTweetDBTypes, false, strmangle.SetComplement(twitterTweetPrimaryKeyColumns, twitterTweetColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, twitterUserDBTypes, false, strmangle.SetComplement(twitterUserPrimaryKeyColumns, twitterUserColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, twitterUserDBTypes, false, strmangle.SetComplement(twitterUserPrimaryKeyColumns, twitterUserColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*TwitterUser{&b, &c} {
		err = a.SetUser(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.User != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.UserTwitterTweets[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.UserID != x.ID {
			t.Error("foreign key was wrong value", a.UserID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.UserID))
		reflect.Indirect(reflect.ValueOf(&a.UserID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.UserID != x.ID {
			t.Error("foreign key was wrong value", a.UserID, x.ID)
		}
	}
}

func testTwitterTweetsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TwitterTweet{}
	if err = randomize.Struct(seed, o, twitterTweetDBTypes, true, twitterTweetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterTweet struct: %s", err)
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

func testTwitterTweetsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TwitterTweet{}
	if err = randomize.Struct(seed, o, twitterTweetDBTypes, true, twitterTweetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterTweet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := TwitterTweetSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testTwitterTweetsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TwitterTweet{}
	if err = randomize.Struct(seed, o, twitterTweetDBTypes, true, twitterTweetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterTweet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := TwitterTweets().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	twitterTweetDBTypes = map[string]string{`ID`: `bigint`, `UserID`: `bigint`, `TwitterID`: `character varying`, `FullText`: `text`, `TweetAt`: `timestamp with time zone`, `Raw`: `jsonb`, `CreatedAt`: `timestamp with time zone`}
	_                   = bytes.MinRead
)

func testTwitterTweetsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(twitterTweetPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(twitterTweetAllColumns) == len(twitterTweetPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &TwitterTweet{}
	if err = randomize.Struct(seed, o, twitterTweetDBTypes, true, twitterTweetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterTweet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := TwitterTweets().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, twitterTweetDBTypes, true, twitterTweetPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize TwitterTweet struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testTwitterTweetsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(twitterTweetAllColumns) == len(twitterTweetPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &TwitterTweet{}
	if err = randomize.Struct(seed, o, twitterTweetDBTypes, true, twitterTweetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TwitterTweet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := TwitterTweets().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, twitterTweetDBTypes, true, twitterTweetPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize TwitterTweet struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(twitterTweetAllColumns, twitterTweetPrimaryKeyColumns) {
		fields = twitterTweetAllColumns
	} else {
		fields = strmangle.SetComplement(
			twitterTweetAllColumns,
			twitterTweetPrimaryKeyColumns,
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

	slice := TwitterTweetSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testTwitterTweetsUpsert(t *testing.T) {
	t.Parallel()

	if len(twitterTweetAllColumns) == len(twitterTweetPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := TwitterTweet{}
	if err = randomize.Struct(seed, &o, twitterTweetDBTypes, true); err != nil {
		t.Errorf("Unable to randomize TwitterTweet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert TwitterTweet: %s", err)
	}

	count, err := TwitterTweets().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, twitterTweetDBTypes, false, twitterTweetPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize TwitterTweet struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert TwitterTweet: %s", err)
	}

	count, err = TwitterTweets().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}