package geom

import "math"

// WGS84 GCJ02 conversion
// https://github.com/andelf/rust-postgis/blob/master/src/mars.rs
type converter struct {
	casm_rr float64
	casm_t1 float64
	casm_t2 float64
	casm_x1 float64
	casm_y1 float64
	casm_x2 float64
	casm_y2 float64
	casm_f  float64
}

func (c *converter) yj_sin2(x float64) float64 {
	var tt float64
	var ss float64
	var ff float64
	var s2 float64
	var cc int
	ff = 0
	if x < 0 {
		x = -x
		ff = 1
	}

	cc = int(x / 6.28318530717959)

	tt = x - float64(cc)*6.28318530717959
	if tt > 3.1415926535897932 {
		tt = tt - 3.1415926535897932
		if ff == 1 {
			ff = 0
		} else if ff == 0 {
			ff = 1
		}
	}
	x = tt
	ss = x
	s2 = x
	tt = tt * tt
	s2 = s2 * tt
	ss = ss - s2*0.166666666666667
	s2 = s2 * tt
	ss = ss + s2*8.33333333333333e-03
	s2 = s2 * tt
	ss = ss - s2*1.98412698412698e-04
	s2 = s2 * tt
	ss = ss + s2*2.75573192239859e-06
	s2 = s2 * tt
	ss = ss - s2*2.50521083854417e-08
	if ff == 1 {
		ss = -ss
	}
	return ss
}

func (c *converter) transform_yj5(x, y float64) float64 {
	var tt float64
	tt = 300 + 1*x + 2*y + 0.1*x*x + 0.1*x*y + 0.1*math.Sqrt(math.Sqrt(x*x))
	tt = tt + (20*c.yj_sin2(18.849555921538764*x)+20*c.yj_sin2(6.283185307179588*x))*0.6667
	tt = tt + (20*c.yj_sin2(3.141592653589794*x)+40*c.yj_sin2(1.047197551196598*x))*0.6667
	tt = tt + (150*c.yj_sin2(0.2617993877991495*x)+300*c.yj_sin2(0.1047197551196598*x))*0.6667
	return tt
}

func (c *converter) transform_jy5(x, xx float64) float64 {
	var n float64
	var a float64
	var e float64
	a = 6378245
	e = 0.00669342
	n = math.Sqrt(1 - e*c.yj_sin2(x*0.0174532925199433)*c.yj_sin2(x*0.0174532925199433))
	n = (xx * 180) / (a / n * math.Cos(x*0.0174532925199433) * 3.1415926)
	return n
}

func (c *converter) transform_jyj5(x, yy float64) float64 {
	var m float64
	var a float64
	var e float64
	var mm float64
	a = 6378245
	e = 0.00669342
	mm = 1 - e*c.yj_sin2(x*0.0174532925199433)*c.yj_sin2(x*0.0174532925199433)
	m = (a * (1 - e)) / (mm * math.Sqrt(mm))
	return (yy * 180) / (m * 3.1415926)
}

func (c *converter) transform_yjy5(x, y float64) float64 {
	var tt float64
	tt = -100 + 2*x + 3*y + 0.2*y*y + 0.1*x*y + 0.2*math.Sqrt(math.Sqrt(x*x))
	tt = tt + (20*c.yj_sin2(18.849555921538764*x)+20*c.yj_sin2(6.283185307179588*x))*0.6667
	tt = tt + (20*c.yj_sin2(3.141592653589794*y)+40*c.yj_sin2(1.047197551196598*y))*0.6667
	tt = tt + (160*c.yj_sin2(0.2617993877991495*y)+320*c.yj_sin2(0.1047197551196598*y))*0.6667
	return tt
}

func (c *converter) random_yj() float64 {
	var t float64
	var casm_a float64 = 314159269
	var casm_c float64 = 453806245
	c.casm_rr = casm_a*c.casm_rr + casm_c
	t = float64(int(c.casm_rr / 2))
	c.casm_rr = c.casm_rr - t*2
	c.casm_rr = c.casm_rr / 2
	return (c.casm_rr)
}

func (c *converter) iniCasm(w_time, w_lng, w_lat float64) {
	var tt float64
	c.casm_t1 = w_time
	c.casm_t2 = w_time
	tt = float64((int)(w_time / 0.357))
	c.casm_rr = w_time - tt*0.357
	if w_time == 0 {
		c.casm_rr = 0.3
	}

	c.casm_x1 = w_lng
	c.casm_y1 = w_lat
	c.casm_x2 = w_lng
	c.casm_y2 = w_lat
	c.casm_f = 3
}

func (c *converter) wgtochina_lb(wg_flag, wg_lng, wg_lat, wg_heit, wg_week, wg_time int) (x, y float64) {
	var x_add float64
	var y_add float64
	var h_add float64
	var x_l float64
	var y_l float64
	var casm_v float64
	var t1_t2 float64
	var x1_x2 float64
	var y1_y2 float64
	x, y = float64(wg_lng), float64(wg_lat)
	if wg_heit > 5000 {
		return
	}
	x_l = float64(wg_lng)
	x_l = x_l / 3686400.0
	y_l = float64(wg_lat)
	y_l = y_l / 3686400.0
	if x_l < 72.004 {
		return
	}
	if x_l > 137.8347 {
		return
	}
	if y_l < 0.8293 {
		return
	}
	if y_l > 55.8271 {
		return
	}
	if wg_flag == 0 {
		c.iniCasm(float64(wg_time), float64(wg_lng), float64(wg_lat))
		return
	}
	c.casm_t2 = float64(wg_time)
	t1_t2 = (c.casm_t2 - c.casm_t1) / 1000.0
	if t1_t2 <= 0 {
		c.casm_t1 = c.casm_t2
		c.casm_f = c.casm_f + 1
		c.casm_x1 = c.casm_x2
		c.casm_f = c.casm_f + 1
		c.casm_y1 = c.casm_y2
		c.casm_f = c.casm_f + 1
	} else {
		if t1_t2 > 120 {
			if c.casm_f == 3 {
				c.casm_f = 0
				c.casm_x2 = float64(wg_lng)
				c.casm_y2 = float64(wg_lat)
				x1_x2 = c.casm_x2 - c.casm_x1
				y1_y2 = c.casm_y2 - c.casm_y1
				casm_v = math.Sqrt(x1_x2*x1_x2+y1_y2*y1_y2) / t1_t2
				if casm_v > 3185 {
					return
				}
			}
			c.casm_t1 = c.casm_t2
			c.casm_f = c.casm_f + 1
			c.casm_x1 = c.casm_x2
			c.casm_f = c.casm_f + 1
			c.casm_y1 = c.casm_y2
			c.casm_f = c.casm_f + 1
		}
	}
	x_add = c.transform_yj5(x_l-105, y_l-35)
	y_add = c.transform_yjy5(x_l-105, y_l-35)
	h_add = float64(wg_heit)
	x_add = x_add + h_add*0.001 + c.yj_sin2(float64(wg_time)*0.0174532925199433) + c.random_yj()
	y_add = y_add + h_add*0.001 + c.yj_sin2(float64(wg_time)*0.0174532925199433) + c.random_yj()
	x = ((x_l + c.transform_jy5(y_l, x_add)) * 3686400)
	y = ((y_l + c.transform_jyj5(y_l, y_add)) * 3686400)
	return
}

// WGS84 coords to GCJ02
func WGS84ToGCJ02(x, y float64) (float64, float64) {
	c := new(converter)
	x1 := x * 3686400.0
	y1 := y * 3686400.0
	gps_week := 0
	gps_week_time := 0
	gps_height := 0
	x, y = c.wgtochina_lb(
		1,
		int(x1),
		int(y1),
		gps_height,
		gps_week,
		gps_week_time,
	)

	x = x / 3686400.0
	y = y / 3686400.0
	return x, y
}

// GCJ02 coords to WGS84
func GCJ02ToWGS84(x, y float64) (float64, float64) {
	var epsilon float64 = 0.00001
	bisection_find_vals := func(
		x,
		y,
		x0,
		y0,
		x1,
		y1,
		epsilon float64,
	) (float64, float64) {
		var x_, y_ float64
		for {
			x_ = (x0 + x1) / 2.0
			y_ = (y0 + y1) / 2.0
			x_e, y_e := WGS84ToGCJ02(x_, y_)

			if math.Abs(x-x_e) <= epsilon && math.Abs(y-y_e) <= epsilon {
				break
			}

			x_e0, y_e0 := WGS84ToGCJ02(x0, y0)
			x_e1, y_e1 := WGS84ToGCJ02(x1, y1)

			// if over some bound
			adjusted := 1

			if x < x_e0 {
				x0 -= x_e0 - x // instead of 0.5
			} else if x > x_e1 {
				x1 += x - x_e1
			} else {
				adjusted = 0
			}

			if y < y_e0 {
				y0 -= y_e0 - y
			} else if y > y_e1 {
				y1 += y - y_e1
			} else {
				adjusted |= adjusted
			}

			if adjusted == 1 {
				continue
			}

			if x_e0 <= x && x <= x_e {
				x1 = x_
			} else if x_e <= x && x <= x_e1 {
				x0 = x_
			}

			if y_e0 <= y && y <= y_e {
				y1 = y_
			} else if y_e <= y && y <= y_e1 {
				y0 = y_
			}

			if x1-x0 < epsilon*0.1 {
				x0 = x0 - x0*0.01
				x1 = x1 + x1*0.01
			}
			if y1-y0 < epsilon*0.1 {
				y0 = y0 - y0*0.01
				y1 = y1 + y1*0.01
			}
		}
		return x_, y_
	}

	return bisection_find_vals(x, y, x-0.1, y-0.1, x+0.1, y+0.1, epsilon)
}
