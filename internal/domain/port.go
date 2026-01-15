package domain

type Port struct {
	id          string
	name        string
	code        string
	city        string
	country     string
	alias       []string
	regions     []string
	coordinates []float64
	province    string
	timezone    string
	unlocs      []string
}

func (p *Port) ID() string {
	return p.id
}

func (p *Port) Name() string {
	return p.name
}

func (p *Port) Code() string {
	return p.code
}

func (p *Port) City() string {
	return p.city
}

func (p *Port) Country() string {
	return p.country
}

func (p *Port) Alias() []string {
	return p.alias
}

func (p *Port) Regions() []string {
	return p.regions
}

func (p *Port) Coordinates() []float64 {
	return p.coordinates
}

func (p *Port) Province() string {
	return p.province
}

func (p *Port) Timezone() string {
	return p.timezone
}

func (p *Port) Unlocs() []string {
	return p.unlocs
}
