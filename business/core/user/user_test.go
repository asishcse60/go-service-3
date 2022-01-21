package user_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/asishcse60/service/business/core/user"
	"github.com/asishcse60/service/business/data/dbtest"
	"github.com/asishcse60/service/business/sys/auth"
	"github.com/asishcse60/service/foundation/docker"
)

var c *docker.Container

func TestMain(m *testing.M) {
	var err error
	c, err = dbtest.StartDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbtest.StopDB(c)

	m.Run()
}

func TestUser(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, c, "testuser")
	t.Cleanup(teardown)

	core := user.NewCore(log, db)
	t.Log("Given the need to work with User records.")
	{
		testID := 0
		ctx := context.Background()
		now := time.Date(2018, time.October, 1, 0, 0, 0, 0, time.UTC)

		nu := user.NewUser{
			Name:            "Bill Kennedy",
			Email:           "bill@ardanlabs.com",
			Roles:           []string{auth.RoleAdmin},
			Password:        "gophers",
			PasswordConfirm: "gophers",
		}

		usr, err := core.Create(ctx, nu, now)
		if err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to create user : %s.", dbtest.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to create user.", dbtest.Success, testID)

		saved, err := core.QueryByID(ctx, usr.ID)
		if err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user by ID: %s.", dbtest.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to retrieve user by ID.", dbtest.Success, testID)

		if diff := cmp.Diff(usr, saved); diff != "" {
			t.Fatalf("\t%s\tTest %d:\tShould get back the same user. Diff:\n%s", dbtest.Failed, testID, diff)
		}
		t.Logf("\t%s\tTest %d:\tShould get back the same user.", dbtest.Success, testID)

		upd := user.UpdateUser{
			Name:  dbtest.StringPointer("Jacob Walker"),
			Email: dbtest.StringPointer("jacob@ardanlabs.com"),
		}

		if err := core.Update(ctx, usr.ID, upd, now); err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to update user : %s.", dbtest.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to update user.", dbtest.Success, testID)

		saved, err = core.QueryByEmail(ctx, *upd.Email)
		if err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user by Email : %s.", dbtest.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to retrieve user by Email.", dbtest.Success, testID)

		if saved.Name != *upd.Name {
			t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Name.", dbtest.Failed, testID)
			t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Name)
			t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Name)
		} else {
			t.Logf("\t%s\tTest %d:\tShould be able to see updates to Name.", dbtest.Success, testID)
		}

		if saved.Email != *upd.Email {
			t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Email.", dbtest.Failed, testID)
			t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Email)
			t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Email)
		} else {
			t.Logf("\t%s\tTest %d:\tShould be able to see updates to Email.", dbtest.Success, testID)
		}

		if err := core.Delete(ctx, usr.ID); err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to delete user : %s.", dbtest.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to delete user.", dbtest.Success, testID)

		_, err = core.QueryByID(ctx, usr.ID)
		if !errors.Is(err, user.ErrNotFound) {
			t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve user : %s.", dbtest.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve user.", dbtest.Success, testID)
	}

}