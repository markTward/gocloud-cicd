package cicd

type GKE struct {
	Name         string
	Project      string
	Cluster      string
	Computezone  string
	Keyfile      string
	Requirements reqs
	Context     string
}
