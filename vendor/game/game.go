package game

var (
	//Players its clients
	Players []*Player

	//Render flag if it false, graphics elements(bars,aims,trails,etc...) should not be initialized.
	Render bool
)

func LookupPlayer(name string) (*Player, bool) {
	for _, p := range Players {
		if p.Name == name {
			return p, true
		}
	}

	return nil, false
}
