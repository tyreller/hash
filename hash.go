package diccionario

import (
	"fmt"
	"tdas/lista"
	TDALista "tdas/lista"
)

// Tamaño del diccionario, numero primo
const fixedSize int = 151

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
	tabla := make([]TDALista.Lista[parClaveValor[K, V]], fixedSize)

	for i := 0; i < fixedSize; i++ {
		tabla[i] = lista.CrearListaEnlazada[parClaveValor[K, V]]()
	}
	return &hashAbierto[K, V]{
		tabla:    tabla,
		tam:      fixedSize,
		cantidad: 0,
	}
}

func (h *hashAbierto[K, V]) Guardar(clave K, dato V) {
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
	}
	h.cantidad++
	lista.InsertarUltimo(par)
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
	//Busca donde esta la primera celda no-vacia de la tabla
	return &iterHashAbierto[K, V]{h, primerIndice, 0}
}

func (iterHash *iterHashAbierto[K, V]) HaySiguiente() bool {
	tablaHash := iterHash.dict.tabla
	iterLista := tablaHash[iterHash.indice].Iterador()

	for i := 0; i < iterHash.posIndice; i++ {
		iterLista.Siguiente()
		//Notar que este siguiente es funcion de Lista, no de hash (por lo tanto, no es definicion circular)
	}

	if iterLista.HaySiguiente() {
		return true
	}

	for i := iterHash.indice + 1; i < iterHash.dict.tam; i++ {
		if !tablaHash[i].EstaVacia() {
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
	iterLista := tablaHash[iterHash.indice].Iterador()

	for i := 0; i < iterHash.posIndice; i++ {
		iterLista.Siguiente()
	}

	if iterLista.HaySiguiente() {
		iterHash.posIndice++
	} else {
		for i := iterHash.indice + 1; i < iterHash.dict.tam; i++ {
			if !tablaHash[i].EstaVacia() {
				iterHash.indice = i
				iterHash.posIndice = 0
				return
			}
		}
	}
}
