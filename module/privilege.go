package module

type Cost struct {
	CostRMB float64
	CostUSD float64
	CostAUD float64
}

type Membership struct {
	Cost
	Space       int // megabytes
	NCasePerMon int // case limit per month
}
