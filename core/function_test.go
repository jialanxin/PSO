package core

import (
	"testing"
)

func TestFunctionToCalc(t *testing.T) {
	y := FunctionToCalc(0.0, []float64{1.0, 2.0})
	if y != 2.0 {
		t.Error("Fail of Calc")
	} else {
		t.Log(y)
	}
}
func TestCalcForXPoints(t *testing.T) {
	xList := []float64{-1.0, 0.0, 1.0}
	yList := CalcForXPoints(xList, []float64{1.0, 2.0})
	t.Log(yList)
}
func TestMSELoss(t *testing.T) {
	xList := []float64{-1.0, 0.0, 1.0}
	yTarget := []float64{3.0, 2.0, 3.0}
	params := []float64{1.0, 2.0}
	MSE := MSELoss(xList, yTarget, params)
	t.Log(MSE)
}
