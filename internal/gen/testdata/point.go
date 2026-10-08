package models

type Point struct {
	X    int      `moo:"default=0"`
	Y    int      `moo:"default=0;required"`
	Name string   `moo:"default=\"origin\""`
	ID   int64    `moo:"readonly"`
	Tags []string `moo:""`
	skip int
	Ignore int `moo:"-"`
}
