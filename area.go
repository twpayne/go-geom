package geom

func area1(flatCoords []float64, offset, end, stride int) float64 {
	var doubleArea float64
	for i := offset + stride; i < end; i += stride {
		doubleArea += (flatCoords[i+1] - flatCoords[i+1-stride]) * (flatCoords[i] + flatCoords[i-stride])
	}
	return doubleArea / 2
}

func area2(flatCoords []float64, offset int, ends []int, stride int) float64 {
	var area float64
	for i, end := range ends {
		a := area1(flatCoords, offset, end, stride)
		if i == 0 {
			area = a
		} else {
			area -= a
		}
		offset = end
	}
	return area
}

func area3(flatCoords []float64, offset int, endss [][]int, stride int) float64 {
	var area float64
	for _, ends := range endss {
		area += area2(flatCoords, offset, ends, stride)
		offset = ends[len(ends)-1]
	}
	return area
}
