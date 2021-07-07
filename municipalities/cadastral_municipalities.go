package municipalities

// CadastralMunicipality represents a cadastral municipality.
type CadastralMunicipality struct {
	ID   int64        `json:"id"`
	Name string       `json:"name"`
	Type CadastreType `json:"type"`
}

// CadastreType is a numeric code identifying the cadastre type.
type CadastreType uint8

// Known cadastre types.
// Should be updated as more types appear in the wild.
const (
	RealEstateCadastre        CadastreType = 3 // Katastar nepokretnosti
	PartialRealEstateCadastre CadastreType = 5 // Katastar nepokretnosti na delu katastarske opštine
	FormingRealEstateCadastre CadastreType = 9 // Osnivanje katastra nepokretnosti kroz postupak komasacije
)
