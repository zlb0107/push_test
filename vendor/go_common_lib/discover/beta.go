package discover

import (
	"math"
	"math/rand"
)

// Beta分布描述的是单变量分布，Dirichlet分布描述的是多变量分布，因此，Beta分布可作为二项分布的先验概率，Dirichlet分布可作为多项分布的先验概率。这两个分布都用到了Gamma函数。

// NextBeta 从贝塔分布中取一个值
func NextBeta(α float64, β float64) float64 {
	dα := []float64{α, β}
	return NextDirichlet(dα)[0]
}

// NextDirichlet 从狄利克雷分布取值
func NextDirichlet(α []float64) []float64 {
	x := make([]float64, len(α))
	sum := float64(0.0)
	for i := 0; i < len(α); i++ {
		x[i] = NextGamma(α[i], 1.0)
		sum += x[i]
	}
	for i := 0; i < len(α); i++ {
		x[i] /= sum
	}
	return x
}

// NextGamma 从Gamma分布中取一个数
func NextGamma(α, λ float64) float64 {
	//if α is a small integer, this way is faster on my laptop
	if α == float64(int64(α)) && α <= 15 {
		x := NextExp(λ)
		for i := 1; i < int(α); i++ {
			x += NextExp(λ)
		}
		return x
	}

	if α < 0.75 {
		return RejectionSample(Gamma_PDF(α, λ), Exp_PDF(λ), Exp(λ), 1)
	}

	//Tadikamalla ACM '73
	a := α - 1
	b := 0.5 + 0.5*sqrt(4*α-3)
	c := a * (1 + b) / b
	d := (b - 1) / (a * b)
	s := a / b
	p := 1.0 / (2 - exp(-s))
	var x, y float64
	for i := 1; ; i++ {
		u := NextUniform()
		if u > p {
			var e float64
			for e = -log((1 - u) / (1 - p)); e > s; e = e - a/b {
			}
			x = a - b*e
			y = a - x
		} else {
			x = a - b*log(u/p)
			y = x - a
		}
		u2 := NextUniform()
		if log(u2) <= a*log(d*x)-x+y/b+c {
			break
		}
	}
	return x / λ
}

// 从lambda=λ的指数分布中获取一个float64
func NextExp(λ float64) float64 { return rand.ExpFloat64() / λ }

func RejectionSample(targetDensity func(float64) float64, sourceDensity func(float64) float64, source func() float64, K float64) float64 {
	x := source()
	for ; NextUniform() >= targetDensity(x)/(K*sourceDensity(x)); x = source() {

	}
	return x
}

// 概率密度函数
func Gamma_PDF(k float64, θ float64) func(x float64) float64 {
	return func(x float64) float64 {
		if x < 0 {
			return 0
		}
		return pow(x, k-1) * exp(-x/θ) / (Γ(k) * pow(θ, k))
	}
}

func Exp_PDF(λ float64) func(x float64) float64 {
	return func(x float64) float64 {
		if x < 0 {
			return 0
		}
		return λ * NextExp(-1*λ*x)
	}
}

func Exp(λ float64) func() float64 { return func() float64 { return NextExp(λ) } }

var NextUniform func() float64 = rand.Float64

var log func(float64) float64 = math.Log
var exp func(float64) float64 = math.Exp
var sqrt func(float64) float64 = math.Sqrt
var pow func(float64, float64) float64 = math.Pow
var Γ = math.Gamma
