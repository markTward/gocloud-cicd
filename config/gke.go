package config

type GKE struct {
	Name        string
	Enabled     bool
	Project     string
	Cluster     string
	Computezone string
	Keyfile     string
}
