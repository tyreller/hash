package diccionario

import (
	"fmt"
	"hash/crc32"
	"tdas/lista"
	TDALista "tdas/lista"
)

// Tamaño Inicial del diccionario, numero primo
const INITIAL_SIZE int = 19

const BIGGER_SIZE_THRESHOLD int = 8
const BIGGER_HASH_FACTOR int = 2

const SMALLER_SIZE_THRESHOLD int = 8
const SMALLER_HASH_FACTOR int = 2

type parClaveValor[K comparable, V any] struct {
	clave K
	dato  V
}

type hashAbierto[K comparable, V any] struct {
	tabla    []TDALista.Lista[parClaveValor[K, V]]
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
	bytes := convertirABytes(clave)
	hash := crc32.ChecksumIEEE(bytes)
	index := hash % uint32(len(h.tabla))
	return int(index)

}

func crearTabla[K comparable, V any](size int) []TDALista.Lista[parClaveValor[K, V]] {
	tabla := make([]TDALista.Lista[parClaveValor[K, V]], size)
	for i := 0; i < size; i++ {
		tabla[i] = lista.CrearListaEnlazada[parClaveValor[K, V]]()
	}
	return tabla
}

func CrearHash[K comparable, V any]() Diccionario[K, V] {
	tabla := crearTabla[K, V](INITIAL_SIZE)

	return &hashAbierto[K, V]{
		tabla:    tabla,
		cantidad: 0,
	}
}

func redimensionar[K comparable, V any](h *hashAbierto[K, V], newSize int) {
	if newSize < INITIAL_SIZE {
		return
	}

	//Save old values
	oldTabla := h.tabla
	oldSize := len(h.tabla)

	//Make new table and reset elements counter
	h.tabla = crearTabla[K, V](newSize)
	h.cantidad = 0

	for i := 0; i < oldSize; i++ {
		iterLista := oldTabla[i].Iterador()
		for iterLista.HaySiguiente() {
			h.Guardar(iterLista.VerActual().clave, iterLista.VerActual().dato)
			iterLista.Siguiente()
		}
	}
}

func (h *hashAbierto[K, V]) Guardar(clave K, dato V) {
	if h.cantidad > len(h.tabla)*BIGGER_SIZE_THRESHOLD {
		redimensionar(h, len(h.tabla)*BIGGER_HASH_FACTOR)
	}

	indice := h.hashFuncIndice(clave)
	lista := h.tabla[indice]
	listaIter := lista.Iterador()
	par := parClaveValor[K, V]{clave: clave, dato: dato}

	for listaIter.HaySiguiente() {
		if listaIter.VerActual().clave == clave {
			listaIter.Borrar()
			listaIter.Insertar(par)
			return
		}
		listaIter.Siguiente()
	}

	lista.InsertarUltimo(par)
	h.cantidad++
}

func (h *hashAbierto[K, V]) Pertenece(clave K) bool {
	indice := h.hashFuncIndice(clave)
	lista := h.tabla[indice]
	listaIter := lista.Iterador()
	for listaIter.HaySiguiente() {
		par := listaIter.VerActual()
		if par.clave == clave {
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
		if par.clave == clave {
			return par.dato
		}
		listaIter.Siguiente()
	}
	panic("La clave no pertenece al diccionario")
}

func (h *hashAbierto[K, V]) Borrar(clave K) V {
	if h.cantidad < len(h.tabla)/SMALLER_SIZE_THRESHOLD {
		redimensionar(h, len(h.tabla)/SMALLER_HASH_FACTOR)
	}

	indice := h.hashFuncIndice(clave)
	listaIter := h.tabla[indice].Iterador()

	for listaIter.HaySiguiente() {
		par := listaIter.VerActual()
		if par.clave == clave {
			listaIter.Borrar()
			h.cantidad--
			return par.dato
		}
		listaIter.Siguiente()

	}
	panic("La clave no pertenece al diccionario")
}

func (h *hashAbierto[K, V]) Cantidad() int {
	return h.cantidad
}

func (h *hashAbierto[K, V]) Iterar(auxFunction func(clave K, dato V) bool) {
	for i := 0; i < len(h.tabla); i++ {
		listaIter := h.tabla[i].Iterador()
		for listaIter.HaySiguiente() {
			par := listaIter.VerActual()
			//Continua hasta que auxFunction devuelva True
			if !auxFunction(par.clave, par.dato) {
				return
			}
			listaIter.Siguiente()
		}
	}
}

func (h *hashAbierto[K, V]) Iterador() IterDiccionario[K, V] {
	primerIndice := 0

	for primerIndice < len(h.tabla) && h.tabla[primerIndice].EstaVacia() {
		primerIndice++
	}
	if primerIndice == len(h.tabla) {
		//Si se cumple, significa que esta toda el hash esta vacia
		primerIndice = 0
	}
	//Busca donde esta la primera celda no-vacia de la tabla
	return &iterHashAbierto[K, V]{h, primerIndice, 0}
}

func (iterHash *iterHashAbierto[K, V]) HaySiguiente() bool {
	tablaHash := iterHash.dict.tabla
	lista := tablaHash[iterHash.indice]

	if lista.Largo() > iterHash.posIndice {
		return true
	}

	for i := iterHash.indice + 1; i < len(iterHash.dict.tabla); i++ {
		if !(tablaHash[i].EstaVacia()) {
			return true
		}
	}
	return false
}

func (iterHash *iterHashAbierto[K, V]) VerActual() (K, V) {
	if !iterHash.HaySiguiente() {
		panic("El iterador termino de iterar")
	}

	listaIter := iterHash.dict.tabla[iterHash.indice].Iterador()

	for i := 0; i < iterHash.posIndice; i++ {
		listaIter.Siguiente()
	}
	par := listaIter.VerActual()
	return par.clave, par.dato
}

func (iterHash *iterHashAbierto[K, V]) Siguiente() {
	if !iterHash.HaySiguiente() {
		panic("El iterador termino de iterar")
	}

	tablaHash := iterHash.dict.tabla
	lista := tablaHash[iterHash.indice]

	if lista.Largo() > iterHash.posIndice+1 {
		iterHash.posIndice++
		return
	}

	for i := iterHash.indice + 1; i < len(iterHash.dict.tabla); i++ {
		if !(tablaHash[i].EstaVacia()) {
			iterHash.indice = i
			iterHash.posIndice = 0
			return
		}
	}
	iterHash.posIndice++
}
