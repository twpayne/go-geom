package xy_test

import (
	"testing"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
	"github.com/twpayne/go-geom/xy/orientation"
)

func TestAngle(t *testing.T) {
	for i, tc := range []struct {
		p1, p2 geom.Coord
		result float64
	}{
		{
			p1:     geom.Coord{-4.007890598483777e8, 7.149034067497588e8, -4.122305737303918e7},
			p2:     geom.Coord{6.452880325856061e8, -7.013452035812421e7, 6.060122721006607e8},
			result: -0.6437947786359727,
		},
		{
			p1:     geom.Coord{4.415559940389009e8, -1.9410956330428556e7, -3.4032011177462184e8},
			p2:     geom.Coord{-2.4046479004409158e8, -1.495553321588844e9, 3.801260331473494e8},
			result: -2.0036085354443243,
		},
		{
			p1:     geom.Coord{-4.661617106595113e8, 2.5040156355098817e8, 3.097086861435584e8},
			p2:     geom.Coord{1.3712169076859632e8, 7.234387287330664e8, 2.4094721366674533e8},
			result: 0.6649730674276739,
		},
		{
			p1:     geom.Coord{-1.7284360981816053e9, -2.361372303896285e8, -4.184641346325376e8},
			p2:     geom.Coord{-2.420562335141231e8, 6.118145442669621e7, -4.134947093381834e8},
			result: 0.19742318999485298,
		},
		{
			p1:     geom.Coord{-1.6550190854088405e8, -1.2218781891720397e9, -5.787102293280704e8},
			p2:     geom.Coord{9.876584504327146e8, 2.725822820782923e8, -7.19405957696415e8},
			result: 0.9135993912770102,
		},
	} {
		calculated := xy.Angle(tc.p1, tc.p2)

		if calculated != tc.result {
			t.Errorf("Test %v failed: expected xy.Angle(%v, %v) expected \n\t%v but got \n\t%v", i+1, tc.p1, tc.p2, tc.result, calculated)
		}

	}
}

func TestAngleFromOrigin(t *testing.T) {
	for i, tc := range []struct {
		p1     geom.Coord
		result float64
	}{
		{
			p1:     geom.Coord{-643891.5406414514, 6.214131154131615e8, -9.241166163738243e7},
			result: 1.571832499502282,
		},
		{
			p1:     geom.Coord{-5.526240186938026e8, -4.1654756589198244e8, -3.904115882281978e8},
			result: -2.495687539636821,
		},
		{
			p1:     geom.Coord{-9775211.937969998, -7.988444321540045e8, -2.9062555922294575e8},
			result: -1.583032406419488,
		},
		{
			p1:     geom.Coord{-7.170140718747358e8, -5.5130056931151845e7, -5672967.701280272},
			result: -3.064855246212557,
		},
		{
			p1:     geom.Coord{-3.4201112699568516e8, -7.256864044916959e8, 1.5383556733916698e9},
			result: -2.0112159720228164,
		},
	} {
		calculated := xy.AngleFromOrigin(tc.p1)

		if calculated != tc.result {
			t.Errorf("Test %v failed: expected xy.AngleFromOrigin(%v) expected \n\t%v but got \n\t%v", i+1, tc.p1, tc.result, calculated)
		}

	}
}

func TestIsAcute(t *testing.T) {
	for i, tc := range []struct {
		p1, p2, p3 geom.Coord
		result     bool
	}{
		{
			p1:     geom.Coord{-2.9746056181996536e8, 1.283116247239797e9, 3.0124856147872955e8},
			p2:     geom.Coord{2.9337112870686615e8, -1.0822405666887188e9, 9.613329966907622e7},
			p3:     geom.Coord{-3.402935182393674e7, -8.477260955562395e8, 2.4474783489619292e7},
			result: true,
		},
		{
			p1:     geom.Coord{1.2441498441622052e9, -1.9039620247337012e9, 1.3258053125928226e8},
			p2:     geom.Coord{-8.34728749413481e8, 3.979772507634378e8, 5.111888830951517e8},
			p3:     geom.Coord{6.087108620010223e8, 1.8734617987205285e8, -1.0570348250682911e8},
			result: true,
		},
		{
			p1:     geom.Coord{-5.0915064274566126e8, -1.4456369240713427e9, 2.1506319910428783e8},
			p2:     geom.Coord{6.405668498559644e8, -3.791562031465599e8, 5.596300821687293e8},
			p3:     geom.Coord{8.241172353750097e8, -3.9414469756236546e7, -2.702165842686878e8},
			result: false,
		},
		{
			p1:     geom.Coord{-1.435496848555126e9, 3.4072911256794184e7, 2.459210259260985e8},
			p2:     geom.Coord{-1.8459206790266247e9, -1.7220237003056505e9, -7.026074366858591e8},
			p3:     geom.Coord{-1.1784100863898702e9, -3.7082065759031725e8, 3577102.337059896},
			result: true,
		},
		{
			p1:     geom.Coord{1.3200492259293087e8, -1.400507993538053e9, 1.0397860978589308e8},
			p2:     geom.Coord{-1.4174043745880973e8, -2.855324865806007e8, 1.154853523604694e9},
			p3:     geom.Coord{-1.6406303327700076e9, 4.8617091926175547e8, -8.222288062865702e8},
			result: false,
		},
	} {
		calculated := xy.IsAcute(tc.p1, tc.p2, tc.p3)

		if calculated != tc.result {
			t.Errorf("Test %v failed: expected xy.IsAcute(%v, %v, %v) expected %v but got %v", i+1, tc.p1, tc.p2, tc.p3, tc.result, calculated)
		}

	}
}

func TestIsObtuse(t *testing.T) {
	for i, tc := range []struct {
		p1, p2, p3 geom.Coord
		result     bool
	}{
		{
			p1:     geom.Coord{-6.581881182734076e8, -5.1226495000032324e8, 4.942792920863176e8},
			p2:     geom.Coord{-2.8760338491412956e8, -2.7637897930097174e7, -1.3120283887929991e8},
			p3:     geom.Coord{-7.253118635362322e8, 2.854840728999085e8, -3.3865131338040566e8},
			result: false,
		},
		{
			p1:     geom.Coord{-6.052601027752758e8, -1.0390522973193089e9, -5.487930680078092e8},
			p2:     geom.Coord{8.843340231350782e8, -8.723399162019621e8, -3321691.634961795},
			p3:     geom.Coord{-7.543599435337427e8, 1.5808204538931034e9, 1.0818796276370132e9},
			result: false,
		},
		{
			p1:     geom.Coord{-6.193065013327997e8, -6.35194942114364e7, 4.3272539543963164e7},
			p2:     geom.Coord{-5.95973359223499e8, 5945981.053576445, 9.226238629036537e8},
			p3:     geom.Coord{3.9272480109009665e8, 3.088998415162513e8, 6.645348620149242e7},
			result: true,
		},
		{
			p1:     geom.Coord{7.610654287692245e8, -6.658609134050195e8, 1.3293491844735564e8},
			p2:     geom.Coord{4.2667262625053006e8, -3.3481316736032414e8, 1.6475762301338202e8},
			p3:     geom.Coord{-4.199827597001981e8, 7.482292086773602e8, 9.971765404694296e8},
			result: true,
		},
		{
			p1:     geom.Coord{-7.643350938452588e8, -2.7699133391444945e8, 2.702299133834568e8},
			p2:     geom.Coord{-1.7382158607827853e7, 5398823.811261921, 1.5609158933203138e7},
			p3:     geom.Coord{-2.4190532792687505e8, -5.756084856732128e8, -1.460466219293458e8},
			result: false,
		},
	} {
		calculated := xy.IsObtuse(tc.p1, tc.p2, tc.p3)

		if calculated != tc.result {
			t.Errorf("Test %v failed: expected xy.IsObtuse(%v, %v, %v) expected %v but got %v", i+1, tc.p1, tc.p2, tc.p3, tc.result, calculated)
		}

	}
}

func TestAngleBetween(t *testing.T) {
	for i, tc := range []struct {
		p1, p2, p3 geom.Coord
		result     float64
	}{
		{
			p1:     geom.Coord{-8.6092078831365e7, -1.2832262246888882e8, -5.39892066777803e8},
			p2:     geom.Coord{-4.125610572401442e7, 3.097372706101881e8, 1.5483271373430803e8},
			p3:     geom.Coord{1.641532856745057e8, 3.949735922042323e7, 1.9570089185263705e8},
			result: 0.7519299818333081,
		},
		{
			p1:     geom.Coord{-2.8151546579932548e7, -3.18057858177073e8, 4.651812237590953e8},
			p2:     geom.Coord{-3.362790282579993e8, 921376.339215076, 3.993733502580851e8},
			p3:     geom.Coord{-4.3757589855782084e7, 2.7736682744679105e8, 7.852890296262044e8},
			result: 1.5598518932475245,
		},
		{
			p1:     geom.Coord{-2.1434525095170313e8, -3.9586869555708617e8, 8.53673374777788e8},
			p2:     geom.Coord{1.6475387708561451e9, 1.3332417513595498e9, 4.7034371208287525e8},
			p3:     geom.Coord{-1.127313286995323e9, -6.606057228728307e8, -2717521.243700768},
			result: 0.1253788551617605,
		},
		{
			p1:     geom.Coord{1.0664500465302559e9, 8.475985637345538e7, -1.621500824133781e9},
			p2:     geom.Coord{4559347.7496108785, 5.161084478242324e7, -1842932.5795175508},
			p3:     geom.Coord{-4.2862346563618964e8, -8.308086105874093e8, -6.966296470909512e8},
			result: 2.058347140220315,
		},
		{
			p1:     geom.Coord{1.3687909198725855e8, -1.1203973392664804e9, -2.45804716717005e8},
			p2:     geom.Coord{-1.0056483813015188e9, 1.6751128488153452e8, -1.8167151284755492e8},
			p3:     geom.Coord{-2.228428636191809e8, -1.1812102896854641e9, 4.388310794147439e8},
			result: 0.19976523282195835,
		},
	} {
		calculated := xy.AngleBetween(tc.p1, tc.p2, tc.p3)

		if calculated != tc.result {
			t.Errorf("Test %v failed: expected xy.AngleBetween(%v, %v, %v) expected \n\t\t%v but got \n\t\t%v", i+1, tc.p1, tc.p2, tc.p3, tc.result, calculated)
		}

	}
}

func TestAngleBetweenOriented(t *testing.T) {
	for i, tc := range []struct {
		p1, p2, p3 geom.Coord
		result     float64
	}{
		{
			p1:     geom.Coord{-1.3799002832563987e9, 5.999590771085212e8, -4.693581090182036e8},
			p2:     geom.Coord{6.826007948791102e7, -8.657386626766933e8, -1.493830309099963e9},
			p3:     geom.Coord{-6.183224805123262e8, 2.4666014745222422e8, 7271369.117346094},
			result: -0.22640245255136904,
		},
		{
			p1:     geom.Coord{6.796487221259736e7, 1.4775165450533025e9, 3258059.847120839},
			p2:     geom.Coord{-6.803421390423136e8, -1.0234495740416303e9, -5.470859926941457e8},
			p3:     geom.Coord{6.443781032426777e8, -1.810385570385187e8, 6.070318143319839e8},
			result: -0.7136563927685474,
		},
		{
			p1:     geom.Coord{5.120536476740612e7, -2.7176934954242444e8, -2.7027023064203584e8},
			p2:     geom.Coord{8.332976211782128e7, -4.67914336571098e8, -1.24317898024329e9},
			p3:     geom.Coord{-1.2179566171482772e8, -1.7824466580072454e8, -4.298802275705581e8},
			result: 0.4538277103800854,
		},
		{
			p1:     geom.Coord{-8.202691782975099e8, 8.782971263839295e8, -9.219191553882729e8},
			p2:     geom.Coord{-1.2725212616954826e8, 2.2006225859706864e8, -1.9247200296977368e8},
			p3:     geom.Coord{4.0049870580738544e8, 4.591976832016299e7, -2.1777764388295308e8},
			result: -2.7006509608972893,
		},
		{
			p1:     geom.Coord{-7.134986212152288e8, -5.527091163926333e8, 1.256171186717098e9},
			p2:     geom.Coord{7.722824262322676e7, 1.2972244051461305e8, 1.2943775785668051e8},
			p3:     geom.Coord{5.426394733747559e8, -2323555.25265493, -6.024980080960876e7},
			result: 2.1531208615200885,
		},
	} {
		calculated := xy.AngleBetweenOriented(tc.p1, tc.p2, tc.p3)

		if calculated != tc.result {
			t.Errorf("Test %v failed: expected xy.AngleBetweenOriented(%v, %v, %v) expected \n\t%v but got \n\t%v", i+1, tc.p1, tc.p2, tc.p3, tc.result, calculated)
		}

	}
}

func TestInteriorAngle(t *testing.T) {
	for i, tc := range []struct {
		p1, p2, p3 geom.Coord
		result     float64
	}{
		{
			p1:     geom.Coord{9.339625086270301e7, 9.494327011462314e8, -8.832231914445356e8},
			p2:     geom.Coord{-8.685036396637098e7, -9827198.1341636, -5.130707858094123e8},
			p3:     geom.Coord{5.48739535964397e8, 8.532792391532723e8, 2.8251807396930236e8},
			result: 0.44900284899855447,
		},
		{
			p1:     geom.Coord{6.523521917718492e8, -1.7481105701895738e8, 1.381806851427019e9},
			p2:     geom.Coord{-8.91688057161475e7, -1.6987404322706103e9, -2.166188234151498e8},
			p3:     geom.Coord{-5.438779575706835e8, 1.7904042826669493e9, -9.194009291344139e7},
			result: 0.5824487823865407,
		},
		{
			p1:     geom.Coord{-2.115808427748782e8, 4.0164370121586424e8, -4.843953798053123e8},
			p2:     geom.Coord{2.2232659336159042e8, -1.4901190499371336e9, 4.8436342680557925e8},
			p3:     geom.Coord{7.506740282650052e8, -4.8757491165846115e8, -2.1487242670012325e7},
			result: 0.7104856314869243,
		},
		{
			p1:     geom.Coord{-1.3806701997111824e8, 4.733218140107204e7, 5.980208692031132e8},
			p2:     geom.Coord{-3.4264253869461334e8, -5.818205740522029e8, -3.6896886549013627e8},
			p3:     geom.Coord{6.63086981247813e8, -1.9734552701705813e9, -5.945945639340445e8},
			result: 2.2014190828743176,
		},
		{
			p1:     geom.Coord{1.9437329983711197e9, 1.6100156972127568e7, 8.719154732991188e8},
			p2:     geom.Coord{7.108626995403899e8, 1.3066388554032483e9, -4.715294366047639e7},
			p3:     geom.Coord{-3.4631579594085485e8, -4.448719942414226e8, 9.847856755232031e8},
			result: 1.3055971267227189,
		},
	} {
		calculated := xy.InteriorAngle(tc.p1, tc.p2, tc.p3)

		if calculated != tc.result {
			t.Errorf("Test %v failed: expected xy.InteriorAngle(%v, %v, %v) expected \n\t%v but got \n\t%v", i+1, tc.p1, tc.p2, tc.p3, tc.result, calculated)
		}

	}
}

func TestGetAngleOrientation(t *testing.T) {
	for i, tc := range []struct {
		p1, p2 float64
		result orientation.Type
	}{
		{
			p1:     1.5973282539123574e8,
			p2:     1.0509666695558771e9,
			result: orientation.Clockwise,
		},
		{
			p1:     -1.9743974140799935e9,
			p2:     1.690220700227534e8,
			result: orientation.CounterClockwise,
		},
		{
			p1:     1.758686954900797e7,
			p2:     2.27491156028423e7,
			result: orientation.Clockwise,
		},
		{
			p1:     1.6512245510554624e8,
			p2:     3.581973387733263e8,
			result: orientation.CounterClockwise,
		},
		{
			p1:     1.1606004655250182e9,
			p2:     3.8888292684591454e8,
			result: orientation.CounterClockwise,
		},
	} {
		calculated := xy.AngleOrientation(tc.p1, tc.p2)

		if calculated != tc.result {
			t.Errorf("Test %v failed: expected xy.GetAngleOrientation(%v, %v) expected %v but got %v", i+1, tc.p1, tc.p2, tc.result, calculated)
		}

	}
}

func TestNormalize(t *testing.T) {
	for i, tc := range []struct {
		p1     float64
		result float64
	}{
		{
			p1:     7.089301226008829e8,
			result: 0.7579033437162295,
		},
		{
			p1:     1.6423604211038163e8,
			result: -0.3600960607195205,
		},
		{
			p1:     9.606844105626652e8,
			result: 0.8766870561033144,
		},
		{
			p1:     5239293.964126772,
			result: -2.9361486826719343,
		},
		{
			p1:     6.136421257534456e8,
			result: -2.8945550816760957,
		},
	} {
		calculated := xy.Normalize(tc.p1)

		if calculated != tc.result {
			t.Errorf("Test %v failed: expected xy.Normalize(%v) expected %v but got %v", i+1, tc.p1, tc.result, calculated)
		}

	}
}

func TestNormalizePositive(t *testing.T) {
	for i, tc := range []struct {
		p1     float64
		result float64
	}{
		{
			p1:     -2.269415841413788e8,
			result: 0.4870605702066726,
		},
		{
			p1:     4.680315524842384e7,
			result: 3.198674730205582,
		},
		{
			p1:     4.5465330578180933e8,
			result: 0.2790471976134583,
		},
		{
			p1:     4.18319606111153e7,
			result: 1.9473086960627342,
		},
		{
			p1:     -1.0427918153375134e8,
			result: 5.003804592005487,
		},
	} {
		calculated := xy.NormalizePositive(tc.p1)

		if calculated != tc.result {
			t.Errorf("Test %v failed: expected xy.NormalizePositive(%v) expected %v but got %v", i+1, tc.p1, tc.result, calculated)
		}

	}
}

func TestDiff(t *testing.T) {
	for i, tc := range []struct {
		p1, p2 float64
		result float64
	}{
		{
			p1:     -5.976261773911254e7,
			p2:     1.5847324519716722e8,
			result: -2.1823585665309447e8,
		},
		{
			p1:     -1.019120645031252e8,
			p2:     -1.5011529441975794e9,
			result: -1.399240873411269e9,
		},
		{
			p1:     -8.346336466770616e8,
			p2:     -8.035798233809209e8,
			result: -3.1053817012955364e7,
		},
		{
			p1:     1.7851664990995303e8,
			p2:     -2.371991702990724e8,
			result: -4.1571581392584014e8,
		},
		{
			p1:     -1.855711053106832e7,
			p2:     4.015083173132894e8,
			result: -4.200654215611724e8,
		},
	} {
		calculated := xy.Diff(tc.p1, tc.p2)

		if calculated != tc.result {
			t.Errorf("Test %v failed: expected xy.Diff(%v, %v) expected %v but got %v", i+1, tc.p1, tc.p2, tc.result, calculated)
		}

	}
}
