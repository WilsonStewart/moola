package processor

type Account struct {
	Name    string
	Kind    string
	Balance float64
}

var validAccountKinds = []string{
	"cash",
	"envelope",
}
