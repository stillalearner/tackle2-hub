package application2

import (
	"testing"

	"github.com/konveyor/tackle2-hub/api"
	. "github.com/smartystreets/goconvey/convey"
)

func TestApplicationCreateGetDelete(t *testing.T) {
	Convey("Create, Get, and Delete an application", t, func() {
		for _, r := range Samples {
			Convey(r.Name, func() {
				So(Application.Create(&r), ShouldBeTrue)

				// Try get.
				got, err := Application.Get(r.ID)
				So(err, ShouldBeNil)

				// Assert the get response.
				So(got, ShouldResemble, r)

				// Try list.
				gotList, err := Application.List()
				So(err, ShouldBeNil)

				// Assert the list response.
				foundR := api.Application{}
				for _, listR := range gotList {
					if listR.Name == r.Name && listR.ID == r.ID {
						foundR = listR
						break
					}
				}
				So(foundR, ShouldResemble, r)

				// Try delete.
				So(Application.Delete(got.ID), ShouldBeTrue)

				// Check the created application was deleted.
				_, err = Application.Get(r.ID)
				So(err, ShouldNotBeNil)
			})
		}
	})
}

func TestApplicationNotCreateDuplicates(t *testing.T) {
	Convey("Not create duplicate applications", t, func() {
		r := Minimal

		// Create sample.
		So(Application.Create(&r), ShouldBeTrue)

		// Prepare Application with duplicate Name.
		dup := &api.Application{
			Name: r.Name,
		}

		// Try create the duplicate.
		err := Application.Create(dup)
		So(err, ShouldNotBeNil)

		// Clean the duplicate.
		So(Application.Delete(dup.ID), ShouldBeTrue)

		// Clean.
		So(Application.Delete(r.ID), ShouldBeTrue)
	})
}

func TestApplicationNotCreateWithoutName(t *testing.T) {
	Convey("Not create application without name", t, func() {
		// Prepare Application without Name.
		r := &api.Application{
			Name: "",
		}

		// Try create the duplicate Application.
		err := Application.Create(r)
		So(err, ShouldNotBeNil)

		// Clean.
		So(Application.Delete(r.ID), ShouldBeTrue)
	})
}
