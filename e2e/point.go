package e2e

//go:generate moo point.go

type Point struct {
	X    int      `moo:"default=0"`
	Y    int      `moo:"default=0;required"`
	Name string   `moo:"default=\"origin\""`
	ID   int64    `moo:"readonly"`
	Tags []string `moo:""`
	secret int
	Ignore int `moo:"-"`
}
