package application3

import (
	"testing"

	"github.com/konveyor/tackle2-hub/api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = BeforeSuite(func() {
	// Set up any necessary fixtures or configurations before running the entire test suite
	gotList, err := Application.List()
	Expect(err).To(BeNil())

	for _, listR := range gotList {
		err:= Application.Delete(listR.ID);		
		Expect(err).To(BeNil());
	}

})

var _ = AfterSuite(func() {
	// Set up any necessary fixtures or configurations before running the entire test suite
	gotList, err := Application.List()
	Expect(err).To(BeNil())

	for _, listR := range gotList {
		err:= Application.Delete(listR.ID);		
		Expect(err).To(BeNil());
	}

})

var _ = Describe("Application", func() {
	Describe("Create, Get, and Delete", func() {

		for _, r := range Samples {
			Context("Given an application with name "+r.Name, func() {
				It("Should create the application successfully", func() {
					err := Application.Create(&r)
					Expect(err).To(BeNil())

					got, err := Application.Get(r.ID)
					Expect(err).To(BeNil())
					Expect(*got).To(Equal(r))
				})

				It("Should list the created application successfully", func() {
					gotList, err := Application.List()
					Expect(err).To(BeNil())

					foundR := api.Application{}
					for _, listR := range gotList {
						if listR.Name == r.Name && listR.ID == r.ID {
							foundR = listR
							break
						}
					}
					Expect(foundR).To(Equal(r))
				})

				It("Should delete the application successfully", func() {
					err := Application.Delete(r.ID)
					Expect(err).To(BeNil())

					_, err = Application.Get(r.ID)
					Expect(err).ToNot(BeNil())
				})
			})
		}
	})

	Describe("Duplicate Prevention", func() {
		BeforeEach(func() {
			// Set up any necessary test fixtures or configurations
		})

		AfterEach(func() {
			// Clean up any test fixtures or configurations
		})

		It("Should prevent creating duplicate applications", func() {
			r := Minimal
			Expect(Application.Create(&r)).To(BeNil())
			defer func() {
				Expect(Application.Delete(r.ID)).To(BeNil())
			}()

			dup := &api.Application{
				Name: r.Name,
			}
			err := Application.Create(dup)
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Creation Without Name", func() {
		BeforeEach(func() {
			// Set up any necessary test fixtures or configurations
		})

		AfterEach(func() {
			// Clean up any test fixtures or configurations
		})

		It("Should fail to create an application without a name", func() {
			r := &api.Application{
				Name: "",
			}

			err := Application.Create(r)
			Expect(err).ToNot(BeNil())
		})
	})
})

func TestApplication(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Application Suite")
}
