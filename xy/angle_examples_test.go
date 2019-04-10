package xy_test

import (
	"fmt"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
)

func ExampleAngle() {
	p1 := geom.Coord{-4.007890598483777e8, 7.149034067497588e8, -4.122305737303918e7}
	p2 := geom.Coord{6.452880325856061e8, -7.013452035812421e7, 6.060122721006607e8}

	angle := xy.Angle(p1, p2)
	fmt.Println(angle)
	// Output: -0.6437947786359727
}

func ExampleAngleFromOrigin() {
	p1 := geom.Coord{-643891.5406414514, 6.214131154131615e8, -9.241166163738243e7}
	angle := xy.AngleFromOrigin(p1)
	fmt.Println(angle)
	// Output: 1.571832499502282
}

func ExampleIsAcute() {
	p1 := geom.Coord{-2.9746056181996536e8, 1.283116247239797e9, 3.0124856147872955e8}
	p2 := geom.Coord{2.9337112870686615e8, -1.0822405666887188e9, 9.613329966907622e7}
	p3 := geom.Coord{-3.402935182393674e7, -8.477260955562395e8, 2.4474783489619292e7}

	isAcute := xy.IsAcute(p1, p2, p3)
	fmt.Println(isAcute)
	// Output: true
}

func ExampleIsObtuse() {
	p1 := geom.Coord{-6.581881182734076e8, -5.1226495000032324e8, 4.942792920863176e8}
	p2 := geom.Coord{-2.8760338491412956e8, -2.7637897930097174e7, -1.3120283887929991e8}
	p3 := geom.Coord{-7.253118635362322e8, 2.854840728999085e8, -3.3865131338040566e8}

	isObtuse := xy.IsObtuse(p1, p2, p3)
	fmt.Println(isObtuse)
	// Output: false
}

func ExampleAngleBetween() {
	p1 := geom.Coord{-8.6092078831365e7, -1.2832262246888882e8, -5.39892066777803e8}
	p2 := geom.Coord{-4.125610572401442e7, 3.097372706101881e8, 1.5483271373430803e8}
	p3 := geom.Coord{1.641532856745057e8, 3.949735922042323e7, 1.9570089185263705e8}

	angle := xy.AngleBetween(p1, p2, p3)
	fmt.Println(angle)
	// Output: 0.7519299818333081
}

func ExampleAngleBetweenOriented() {
	p1 := geom.Coord{-1.3799002832563987e9, 5.999590771085212e8, -4.693581090182036e8}
	p2 := geom.Coord{6.826007948791102e7, -8.657386626766933e8, -1.493830309099963e9}
	p3 := geom.Coord{-6.183224805123262e8, 2.4666014745222422e8, 7271369.117346094}

	angle := xy.AngleBetweenOriented(p1, p2, p3)
	fmt.Println(angle)
	// Output: -0.22640245255136904
}

func ExampleInteriorAngle() {
	p1 := geom.Coord{9.339625086270301e7, 9.494327011462314e8, -8.832231914445356e8}
	p2 := geom.Coord{-8.685036396637098e7, -9827198.1341636, -5.130707858094123e8}
	p3 := geom.Coord{5.48739535964397e8, 8.532792391532723e8, 2.8251807396930236e8}

	angle := xy.InteriorAngle(p1, p2, p3)
	fmt.Println(angle)
	// Output: 0.44900284899855447
}

func ExampleAngleOrientation() {
	p1 := 1.5973282539123574e8
	p2 := 1.0509666695558771e9

	orient := xy.AngleOrientation(p1, p2)
	fmt.Println(orient)
	// Output: Clockwise
}

func ExampleNormalize() {
	p1 := 7.089301226008829e8

	normalized := xy.Normalize(p1)
	fmt.Println(normalized)
	// Output: 0.7579033437162295
}

func ExampleNormalizePositive() {
	p1 := -2.269415841413788e8

	normalized := xy.NormalizePositive(p1)
	fmt.Println(normalized)
	// Output: 0.4870605702066726
}

func ExampleDiff() {
	p1 := -5.976261773911254e7
	p2 := 1.5847324519716722e8

	diff := xy.Diff(p1, p2)
	fmt.Println(diff)
	// Output: -2.1823585665309447e+08
}
