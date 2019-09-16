package options

const (
	BetaParamsGlobalKey = "global"
)

type Options struct {
	MergeOp *MergeOption
}

type MergeOption struct {
	Method        string
	BetaParamsMap map[string]BetaParams
}

type BetaParams struct {
	Alpha float64
	Beta  float64
}
