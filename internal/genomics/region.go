package genomics

type Region struct {
	Name, Start, End string
}

func (r *Region) String() string {
	if r.Name == "" || r.Name == "*" || r.Start == "-1" {
		return r.Name
	}
	if r.End == "-1" {
		return r.Name + ":" + r.Start
	} else {
		return r.Name + ":" + r.Start + "-" + r.End
	}
}
