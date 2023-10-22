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
	return &iterHashAbierto[K, V]{h, 0, 0}
}

func (iterHash *iterHashAbierto[K, V]) HaySiguiente() bool {
	for i := iterHash.indice; i < iterHash.dict.tam; i++ {
		lista := iterHash.dict.tabla[i]
		if i == iterHash.indice {
			lista = lista.SubListaDesde(iterHash.posIndice)
		}
		listaIter := lista.Iterador()
		if listaIter.HaySiguiente() {
			return true
		}
		iterHash.indice = i + 1
		iterHash.posIndice = 0
	}
	return false
}

func (iterHash *iterHashAbierto[K, V]) VerActual() (K, V) {
	if !iterHash.HaySiguiente() {
		panic("El iterador terminó de iterar")
	}
	lista := iterHash.dict.tabla[iterHash.indice]
	if iterHash.indice == iterHash.dict.tam-1 {
		return lista.Ultimo().(parClaveValor[K, V]).clav, lista.Ultimo().(parClaveValor[K, V]).dat
	}
	if iterHash.indice == iterHash.dict.tam {
		return lista.Ultimo().(parClaveValor[K, V]).clav, lista.Ultimo().(parClaveValor[K, V]).dat
	}
	return lista.Elemento(iterHash.posIndice).(parClaveValor[K, V]).clav, lista.Elemento(iterHash.posIndice).(parClaveValor[K, V]).dat
}

func (iterHash *iterHashAbierto[K, V]) Siguiente() {
	if !iterHash.HaySiguiente() {
		panic("El iterador terminó de iterar")
	}
	iterHash.posIndice++
}
