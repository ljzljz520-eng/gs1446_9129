package ledger

type Category struct {
	Name        string
	TaxCode     string
	Deductible  bool
	DefaultRate int64
}

var Categories = []Category{
	{Name: "consulting", TaxCode: "C01", Deductible: true, DefaultRate: 100},
	{Name: "design", TaxCode: "D01", Deductible: true, DefaultRate: 100},
	{Name: "development", TaxCode: "S01", Deductible: true, DefaultRate: 100},
	{Name: "writing", TaxCode: "W01", Deductible: true, DefaultRate: 100},
	{Name: "photography", TaxCode: "M01", Deductible: true, DefaultRate: 80},
	{Name: "travel", TaxCode: "T01", Deductible: true, DefaultRate: 50},
	{Name: "equipment", TaxCode: "E01", Deductible: true, DefaultRate: 100},
	{Name: "software", TaxCode: "E02", Deductible: true, DefaultRate: 100},
	{Name: "education", TaxCode: "E03", Deductible: true, DefaultRate: 100},
	{Name: "insurance", TaxCode: "I01", Deductible: true, DefaultRate: 100},
	{Name: "office", TaxCode: "O01", Deductible: true, DefaultRate: 100},
	{Name: "phone", TaxCode: "O02", Deductible: true, DefaultRate: 70},
	{Name: "internet", TaxCode: "O03", Deductible: true, DefaultRate: 70},
	{Name: "advertising", TaxCode: "A01", Deductible: true, DefaultRate: 100},
	{Name: "banking", TaxCode: "B01", Deductible: true, DefaultRate: 100},
	{Name: "legal", TaxCode: "L01", Deductible: true, DefaultRate: 100},
	{Name: "accounting", TaxCode: "L02", Deductible: true, DefaultRate: 100},
	{Name: "membership", TaxCode: "M02", Deductible: true, DefaultRate: 50},
	{Name: "meals", TaxCode: "T02", Deductible: true, DefaultRate: 50},
	{Name: "lodging", TaxCode: "T03", Deductible: true, DefaultRate: 50},
	{Name: "mileage", TaxCode: "T04", Deductible: true, DefaultRate: 67},
	{Name: "shipping", TaxCode: "S02", Deductible: true, DefaultRate: 100},
	{Name: "printing", TaxCode: "P01", Deductible: true, DefaultRate: 100},
	{Name: "subcontractor", TaxCode: "S03", Deductible: true, DefaultRate: 100},
	{Name: "commission", TaxCode: "S04", Deductible: true, DefaultRate: 100},
	{Name: "rent", TaxCode: "R01", Deductible: true, DefaultRate: 100},
	{Name: "utilities", TaxCode: "R02", Deductible: true, DefaultRate: 100},
	{Name: "charity", TaxCode: "G01", Deductible: false, DefaultRate: 0},
	{Name: "personal", TaxCode: "X01", Deductible: false, DefaultRate: 0},
	{Name: "refund", TaxCode: "X02", Deductible: false, DefaultRate: 0},
}

func FindCategory(name string) (Category, bool) {
	for _, c := range Categories {
		if c.Name == name {
			return c, true
		}
	}
	return Category{}, false
}
func Deduction(name string, amount int64) int64 {
	c, ok := FindCategory(name)
	if !ok || !c.Deductible {
		return 0
	}
	return amount * c.DefaultRate / 100
}
func TaxCode(name string) string {
	if c, ok := FindCategory(name); ok {
		return c.TaxCode
	}
	return "UNK"
}
