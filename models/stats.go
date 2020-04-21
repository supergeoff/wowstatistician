package models

type Stats struct {
	SyncDate      string          `json:"syncdate"`
	Source        string          `json:"source"`
	Overall       int             `json:"overall"`
	Distributions []*Distribution `json:"distributions"`
}

type Distribution struct {
	Class string  `json:"class"`
	Total int     `json:"total"`
	Specs []*Spec `json:"specs"`
}

type Spec struct {
	Spec  string `json:"spec"`
	Count int    `json:"count"`
}

func (s *Stats) FindDistribution(class string) *Distribution {
	for _, v := range s.Distributions {
		if v.Class == class {
			return v
		}
	}
	return nil
}

func (d *Distribution) FindSpec(spec string) *Spec {
	for _, v := range d.Specs {
		if v.Spec == spec {
			return v
		}
	}
	return nil
}
