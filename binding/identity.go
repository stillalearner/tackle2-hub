package binding

import (
	"github.com/konveyor/tackle2-hub/api"
)

//
// Identity API.
type Identity struct {
	// hub API client.
	client *Client
}

//
// Create a Identity.
func (h *Identity) Create(r *api.Identity) (err error) {
	err = h.client.Post(api.IdentitiesRoot, &r)
	return
}

//
// Get a Identity by ID.
func (h *Identity) Get(id uint) (r *api.Identity, err error) {
	r = &api.Identity{}
	path := Path(api.IdentityRoot).Inject(Params{api.ID: id})
	err = h.client.Get(path, r)
	return
}

//
// List Identities.
func (h *Identity) List() (list []api.Identity, err error) {
	list = []api.Identity{}
	err = h.client.Get(api.IdentitiesRoot, &list)
	return
}

//
// Update a Identity.
func (h *Identity) Update(r *api.Identity) (err error) {
	path := Path(api.IdentityRoot).Inject(Params{api.ID: r.ID})
	err = h.client.Put(path, r)
	return
}

//
// Delete a Identity.
func (h *Identity) Delete(id uint) (err error) {
	err = h.client.Delete(Path(api.IdentityRoot).Inject(Params{api.ID: id}))
	return
}
