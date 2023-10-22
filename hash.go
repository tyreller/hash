package diccionario

import (
	"fmt"
	"tdas/lista"
	TDALista "tdas/lista"
)

type parClaveValor[K comparable, V any] struct {
	clav K
	dat  V
}

type hashAbierto[K comparable, V any] struct {
	celda    []TDALista.Lista[parClaveValor[K, V]]
	tam      int
	cantidad int
}

type iterHashAbierto[K comparable, V any] struct {
	dict      *hashAbierto[K, V]
	indice    int
	posIndice int
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

// Tamaño fijo en 151, numero primo
func CrearHash[K comparable, V any]() Diccionario[K, V] {
	size := 151
	celda := make([]TDALista.Lista[parClaveValor[K, V]], size)
	for i := 0; i < size; i++ {
		celda[i] = lista.CrearListaEnlazada[parClaveValor[K, V]]()
	}
	return &hashAbierto[K, V]{
		celda:    celda,
		tam:      size,
		cantidad: 0,
	}
}

func (h *hashAbierto[K, V]) Guardar(clave K, dato V) {
	indice := h.hashFuncIndice(clave)
	lista := h.celda[indice]
	par := parClaveValor[K, V]{clav: clave, dat: dato}
	lista.InsertarUltimo(par)
	h.cantidad++
}

func (h *hashAbierto[K, V]) Pertenece(clave K) bool {
	indice := h.hashFuncIndice(clave)
	lista := h.celda[indice]
	listaIter := lista.Iterador()

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
	listaIter := h.celda[indice].Iterador()

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
	listaIter := h.celda[indice].Iterador()

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

func (h *hashAbierto[K, V]) Iterar(auxFunction func(clave K, dato V) bool) {
	for i := 0; i < h.tam; i++ {
		listaIter := h.celda[i].Iterador()
		for listaIter.HaySiguiente() {
			par := listaIter.VerActual()
			//Continua hasta que auxFunction devuelva True
			if !auxFunction(par.clav, par.dat) {
				return
			}
			listaIter.Siguiente()
		}
	}
}

func (h *hashAbierto[K, V]) Iterador() IterDiccionario[K, V] {
	return &iterHashAbierto[K, V]{h, 0, 0}
}

func (iterHash *iterHashAbierto[K, V]) HaySiguiente() bool {
	lista := iterHash.dict.celda[iterHash.indice]
	iterLista := lista.Iterador()

	if iterLista.HaySiguiente() {
		return true
	}

	return false
}

func (iterHash *iterHashAbierto[K, V]) VerActual() (K, V) {
	if !iterHash.HaySiguiente() {
		panic("El iterador terminó de iterar")
	}

	lista := iterHash.dict.celda[iterHash.indice]
	listaIter := lista.Iterador()
	for i := 0; i < iterHash.posIndice; i++ {
		listaIter.Siguiente()
	}
	par := listaIter.VerActual()
	return par.clav, par.dat
}

func (iterHash *iterHashAbierto[K, V]) Siguiente() {
	if !iterHash.HaySiguiente() {
		panic("El iterador terminó de iterar")
	}
	iterHash.posIndice++
}
