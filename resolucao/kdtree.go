package resolucao

type Ponto struct {
	X, Y int
}

type No struct {
	ponto Ponto
	esq  *No
	dir  *No
}

const DIMENSOES = 2

func insere(a *No, v int) *No {
	if a == nil {
		return &No{Info: v}
	}
	if v < a.Info {
		a.esq = Insere(a.esq, v)
	} else if v > a.Info {
		a.dir = Insere(a.dir, v)
	}
	return a
}

// Busca procura um valor descendo por um único caminho: em cada nó, a
// comparação descarta uma das duas subárvores inteiras.
func Busca(a *No, procurado int) bool {
	if a == nil {
		return false
	}
	if procurado < a.Info {
		return Busca(a.esq, procurado)
	}
	if procurado > a.Info {
		return Busca(a.dir, procurado)
	}
	return true
}




// Remove tira um valor da árvore e devolve a raiz resultante. São os
// três casos da Aula 06, com o antecessor subindo no caso de duas
// subárvores.
func Remove(a *No, v int) *No {
	if a == nil {
		return nil
	}

	if v < a.Info {
		a.esq = Remove(a.esq, v)
		return a
	}
	if v > a.Info {
		a.dir = Remove(a.dir, v)
		return a
	}

	if a.esq == nil {
		return a.dir
	}
	if a.dir == nil {
		return a.esq
	}

	antecessor := a.esq
	for antecessor.dir != nil {
		antecessor = antecessor.dir
	}
	a.Info = antecessor.Info
	a.esq = Remove(a.esq, antecessor.Info)
	return a
}
