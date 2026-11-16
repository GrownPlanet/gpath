package geo

import "math"

type Point struct{ Lat, Lon float64 }

func HaversineDistance(a, b Point) float64 {
	const R = 6371000.0
	rad := math.Pi / 180
	p1, p2 := a.Lat*rad, b.Lat*rad
	dp, dl := (b.Lat-a.Lat)*rad, (b.Lon-a.Lon)*rad
	h := math.Sin(dp/2)*math.Sin(dp/2) +
		math.Cos(p1)*math.Cos(p2)*math.Sin(dl/2)*math.Sin(dl/2)
	return 2 * R * math.Asin(math.Sqrt(h))
}
