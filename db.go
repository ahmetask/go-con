package gocon

type Db interface {
	Save(config Config) error
	Read(appName string) (*Config, error)
	ReadSpecs() ([]Spec, error)
}
