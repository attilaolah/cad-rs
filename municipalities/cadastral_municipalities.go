package municipalities

// CadastralMunicipality represents a cadastral municipality.
type CadastralMunicipality struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}
