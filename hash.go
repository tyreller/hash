package diccionario

import (
	"fmt"
	TDALista "tdas/lista"
)

type parClaveValor[K comparable, V any] struct {
	clave K
	dato  V
}

type hashAbierto[K comparable, V any] struct {
	tabla    []TDALista.Lista[parClaveValor[K, V]]
	tam      int
	cantidad int
}

func convertirABytes[K comparable](clave K) []byte {
	return []byte(fmt.Sprintf("%v", clave))
}

//Tamaño variable
// func CrearHash[K comparable, V any](tam int) Diccionario[K, V] {
// 	tabla := make([]TDALista.Lista[parClaveValor[K, V]], tam)
// 	return &hashAbierto[K, V]{
// 		tabla:    tabla,
// 		tam:      tam,
// 		cantidad: 0,
// 	}
// }

// Tamaño fijo en 151
func CrearHash[K comparable, V any]() Diccionario[K, V] {
	tabla := make([]TDALista.Lista[parClaveValor[K, V]], 151)
	return &hashAbierto[K, V]{
		tabla:    tabla,
		tam:      151,
		cantidad: 0,
	}
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
