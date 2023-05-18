package application2

import (
	"github.com/konveyor/tackle2-hub/api"
)

// Set of valid Application resources for tests and reuse.
// Important: initialize test application from this samples, not use it directly to not affect other tests.
var (
    Minimal = api.Application{
        Name: "Minimal application",
    }
    PathfinderGit = api.Application{
        Name:        "Pathfinder",
        Description: "Tackle Pathfinder application.",
        Repository: &api.Repository{
            Kind:   "git",
            URL:    "https://github.com/konveyor/tackle-pathfinder.git",
            Branch: "1.2.0",
        },
    }
    Samples = []api.Application{Minimal, PathfinderGit}
)

