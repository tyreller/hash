package diccionario

import (
	"fmt"
	"hash/crc32"
	"tdas/lista"
	TDALista "tdas/lista"
)

// Tamaño Inicial del diccionario, numero primo
const initialSize int = 19

type parClaveValor[K comparable, V any] struct {
	clav K
	dat  V
}

type hashAbierto[K comparable, V any] struct {
	tabla    []TDALista.Lista[parClaveValor[K, V]]
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
	bytes := convertirABytes(clave)
	hash := crc32.ChecksumIEEE(bytes)
	index := hash % uint32(h.tam)
	return int(index)

}

func CrearHash[K comparable, V any]() Diccionario[K, V] {
	tabla := make([]TDALista.Lista[parClaveValor[K, V]], initialSize)

	for i := 0; i < initialSize; i++ {
		tabla[i] = lista.CrearListaEnlazada[parClaveValor[K, V]]()
	}
	return &hashAbierto[K, V]{
		tabla:    tabla,
		tam:      initialSize,
		cantidad: 0,
	}
}

func redimensionar[K comparable, V any](h *hashAbierto[K, V], newSize int) {
	if newSize < initialSize {
		return
	}
	newTabla := make([]TDALista.Lista[parClaveValor[K, V]], newSize)
	for i := 0; i < newSize; i++ {
		newTabla[i] = lista.CrearListaEnlazada[parClaveValor[K, V]]()
	}

	newHash := hashAbierto[K, V]{
		tabla:    newTabla,
		tam:      newSize,
		cantidad: 0,
	}

	for i := 0; i < h.tam; i++ {
		iterLista := h.tabla[i].Iterador()
		for iterLista.HaySiguiente() {
			iterClave := iterLista.VerActual().clav
			iterDato := iterLista.VerActual().dat
			newHash.Guardar(iterClave, iterDato)
			iterLista.Siguiente()
		}
	}
	*h = newHash
}

func (h *hashAbierto[K, V]) Guardar(clave K, dato V) {
	if h.cantidad > 2*h.tam {
		redimensionar(h, 4*h.tam)
	}

	indice := h.hashFuncIndice(clave)
	lista := h.tabla[indice]
	listaIter := lista.Iterador()
	par := parClaveValor[K, V]{clav: clave, dat: dato}

	for listaIter.HaySiguiente() {
		if listaIter.VerActual().clav == clave {
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
	if h.cantidad < h.tam/4 {
		redimensionar(h, h.tam/2)
	}

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

func (h *hashAbierto[K, V]) Iterar(auxFunction func(clave K, dato V) bool) {
	for i := 0; i < h.tam; i++ {
		listaIter := h.tabla[i].Iterador()
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
	primerIndice := 0

	for primerIndice < h.tam && h.tabla[primerIndice].EstaVacia() {
		primerIndice++
	}
	if primerIndice == h.tam {
		//Si se cumple, significa que esta toda la lista vacia
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

	for i := iterHash.indice + 1; i < iterHash.dict.tam; i++ {
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
	return par.clav, par.dat
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

	for i := iterHash.indice + 1; i < iterHash.dict.tam; i++ {
		if !(tablaHash[i].EstaVacia()) {
			iterHash.indice = i
			iterHash.posIndice = 0
			return
		}
	}
	iterHash.posIndice++
}
