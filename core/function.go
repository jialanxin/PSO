package core

import "math"

func FunctionToCalc(x float64, params []float64) float64 {
	a := params[0]
	b := params[1]
	return a*math.Pow(x, 2) + b
}

func CalcForXPoints(xList []float64, params []float64) []float64 {
	yList := make([]float64, len(xList))
	for ix := range xList {
		yList[ix] = FunctionToCalc(xList[ix], params)
	}
	return yList
}

func MSELoss(xList []float64, yTarget []float64, params []float64) float64 {
	yPredict := CalcForXPoints(xList, params)
	MSE := 0.0
	for ix := range yPredict {
		MSE += math.Pow(yPredict[ix]-yTarget[ix], 2)
	}
	return MSE
}
