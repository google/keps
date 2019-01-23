package sections

var sectionOrdering = map[string]int{
	// any user defined sections fit here
	Summary:            1,
	Motivation:         2,
	DeveloperGuide:     3,
	OperatorGuide:      4,
	TeacherGuide:       5,
	GraduationCriteria: 6,
	Readme:             99,
}

type ByOrder []string

func (entries ByOrder) Len() int      { return len(entries) }
func (entries ByOrder) Swap(i, j int) { entries[i], entries[j] = entries[j], entries[i] }
func (entries ByOrder) Less(i, j int) bool {
	return sectionOrdering[entries[i]] < sectionOrdering[entries[j]]
}
