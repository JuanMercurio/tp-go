package ports

type Patch struct {
	Op         string `json:"op"`
	Path       string `json:"path"`
	NuevoValor any    `json:"value"`
}
