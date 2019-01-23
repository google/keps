package sections

func Locations(entries []Entry) []string {
	locs := []string{}

	for i := range entries {
		locs = append(locs, entries[i].Filename())
	}

	return locs
}
