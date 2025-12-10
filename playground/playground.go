package playground

import (
	"fmt"
	"math"
)

/*
STRUCT
*/
type PersegiPanjang struct {
	Panjang float64
	Lebar   float64
}

type Lingkaran struct {
	Radius float64
}

type Segitiga struct {
	Alas   float64
	Tinggi float64
}

/*
INTERFACE
*/
type Geometri interface {
	Luas() float64
}

/*
METHOD DECLARATION
*/
func (pp PersegiPanjang) Luas() float64 {
	return pp.Panjang * pp.Lebar
}

func (l Lingkaran) Luas() float64 {
	return math.Pi * l.Radius * l.Radius
}

func (s Segitiga) Luas() float64 {
	return (1.0 / 2.0) * s.Alas * s.Tinggi
}

/*
Fungsi tanpa method (STAND ALONE)
*/

func HitungPersegiBiasa(sisi int) int {
	return sisi * sisi
}

/*
GENERIC FUNCTION
*/
func PrintLuas(g Geometri) {
	fmt.Println("Luas Objek :", g.Luas())
}
