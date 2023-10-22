package diccionario

import (
	"fmt"
	TDALista "tdas/lista"
)

type parClaveValor[K comparable, V any] struct {
	clav K
	dat  V
}

type hashAbierto[K comparable, V any] struct {
	tabla    []TDALista.Lista[parClaveValor[K, V]]
	tam      int
	cantidad int
}

func convertirABytes[K comparable](clave K) []byte {
	return []byte(fmt.Sprintf("%v", clave))
}

func (h *hashAbierto[K, V]) hashFuncIndice(clave K) int {
	indiceBytes := convertirABytes(clave)

	hashValue := 0
	for i := 0; i < len(indiceBytes); i++ {
		b := indiceBytes[i]
		hashValue += int(b)
	}

	indice := hashValue % h.tam
	return indice
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

// Tamaño fijo en 151, numero primo
func CrearHash[K comparable, V any]() Diccionario[K, V] {
	tabla := make([]TDALista.Lista[parClaveValor[K, V]], 151)
	return &hashAbierto[K, V]{
		tabla:    tabla,
		tam:      151,
		cantidad: 0,
	}
}

func (h *hashAbierto[K, V]) Guardar(clave K, dato V) {
	indice := h.hashFuncIndice(clave)
	lista := h.tabla[indice]
	par := parClaveValor[K, V]{clav: clave, dat: dato}
	lista.InsertarUltimo(par)
	h.cantidad++
}

func (h *hashAbierto[K, V]) Pertenece(clave K) bool {
	indice := h.hashFuncIndice(clave)
	listaIter := h.tabla[indice].Iterador()

	for listaIter.HaySiguiente() {
		par := listaIter.VerActual()
		if par.clav == clave {
			return true
		}
		listaIter.Siguiente()
	}
	return false
}

func (h *hashAbierto[K, V]) Obtener(clave K) V {
	indice := h.hashFuncIndice(clave)
	listaIter := h.tabla[indice].Iterador()

	for listaIter.HaySiguiente() {
		par := listaIter.VerActual()
		if par.clav == clave {
			return par.dat
		}
		listaIter.Siguiente()
	}
	panic("La clave no pertenece al diccionario")
}

func (h *hashAbierto[K, V]) Borrar(clave K) V {
	indice := h.hashFuncIndice(clave)
	listaIter := h.tabla[indice].Iterador()

	for listaIter.HaySiguiente() {
		par := listaIter.VerActual()
		if par.clav == clave {
			listaIter.Borrar()
			h.cantidad--
			return par.dat
		}
		listaIter.Siguiente()
	}
	panic("La clave no pertenece al diccionario")
}

func (h *hashAbierto[K, V]) Cantidad() int {
	return h.cantidad
}

func (h *hashAbierto[K, V]) Iterar(func(clave K, dato V) bool) {

}

func (h *hashAbierto[K, V]) Iterador() IterDiccionario[K, V] {

}

func (h *hashAbierto[K, V]) HaySiguiente() bool {

}
func (h *hashAbierto[K, V]) VerActual() (K, V) {
	if VerActual == false {
		panic("El iterador termino de iterar")
	}

}
func (h *hashAbierto[K, V]) Siguiente() {
	if VerActual == false {
		panic("El iterador termino de iterar")
	}
}
