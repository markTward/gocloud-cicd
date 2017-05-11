package config

type Helm struct {
	Name      string
	Version   string
	Release   string
	Namespace string
	ChartPath string
}
