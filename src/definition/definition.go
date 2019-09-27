package definition

type E_Target int

const (
	_ E_Target = iota
	ET_csv
	ET_ini
)

type E_Change int

const (
	_ E_Change = iota
	EC_add_csv
	EC_add_row
	EC_add_column
	EC_change_field
)

type stOneRule struct {
	target    E_Target
	changeway E_Change
	item      string
	content   string
	rule      string
}
