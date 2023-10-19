package diccionario

import "fmt"

func convertirABytes[K comparable](clave K) []byte {
	return []byte(fmt.Sprintf("%v", clave))
}

func CrearHash[K comparable, V any]() Diccionario[K, V] {

}

func Guardar(clave K, dato V) {

}
func Pertenece(clave K) bool {

}
func Obtener(clave K) V {

}
func Borrar(clave K) V {

}
func Cantidad() int {

}
func Iterar(func(clave K, dato V) bool) {

}

func Iterador() IterDiccionario[K, V] {

}

func HaySiguiente() bool {

}
func VerActual() (K, V) {
	if VerActual == false {
		panic("El iterador termino de iterar")
	}

}
func Siguiente() {
	if VerActual == false {
		panic("El iterador termino de iterar")
	}
}
