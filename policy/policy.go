package policy

type Policy struct {
	ID          string   `yaml:"id"`
	Title       string   `yaml:"title"`
	Category    string   `yaml:"category"`
	Subcategory string   `yaml:"subcategory"`
	Severity    string   `yaml:"severity"`
	Levels      []string `yaml:"levels"`
	Tags        []string `yaml:"tags"`

	CheckCmd     string `yaml:"check"`
	RemediateCmd string `yaml:"remediate"`
}

type PolicySet struct {
	Policies []Policy `yaml:"policies"`
}
