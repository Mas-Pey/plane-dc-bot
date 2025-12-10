package playground

func dummy_main() {
	l := Lingkaran{Radius: 7.0}
	s := Segitiga{Alas: 2.0, Tinggi: 3.0}

	PrintLuas(PersegiPanjang{Panjang: 2.0, Lebar: 3.0})
	PrintLuas(l)
	PrintLuas(s)

	persegi := HitungPersegiBiasa(3)
	println("Persegi: ", persegi)
}
