package resolucao

type Ponto struct {
	X, Y int
}

func (ponto Ponto) valorDoEixo(eixo int) int {
	if eixo == 0 {
		return ponto.X
	}
	return ponto.Y
}

type No struct {
	ponto Ponto
	esq   *No
	dir   *No
}

type KDTree struct {
	raiz *No
}

const DIMENSOES = 2

func (arvore *KDTree) insere(ponto Ponto) {
	arvore.raiz = insereRecursivo(arvore.raiz, ponto, 0)
}

func (arvore *KDTree) remove(ponto Ponto) {
	arvore.raiz = removeRecursivo(arvore.raiz, ponto, 0)
}


func insereRecursivo(no *No, ponto Ponto, profundidade int) *No {
	if no == nil {
		return &No{ponto: ponto}
	}

	eixoAtual := profundidade % DIMENSOES

	if ponto.valorDoEixo(eixoAtual) < no.ponto.valorDoEixo(eixoAtual) {
		no.esq = insereRecursivo(no.esq, ponto, profundidade+1)
	} else {
		no.dir = insereRecursivo(no.dir, ponto, profundidade+1)
	}

	return no
}

func removeRecursivo(no *No, ponto Ponto, profundidade int) *No {
	if no == nil {
		return nil
	}

	eixoAtual := profundidade % DIMENSOES

	if no.ponto == ponto {
		if no.dir != nil {
			minimo := acharMinimo(no.dir, eixoAtual, profundidade + 1)
			no.ponto = minimo.ponto
			no.dir = removeRecursivo(no.dir, minimo.ponto, profundidade + 1)
		} else if no.esq != nil {
			minimo := acharMinimo(no.esq, eixoAtual, profundidade + 1)
			no.ponto = minimo.ponto
			no.esq = removeRecursivo(no.esq, minimo.ponto, profundidade + 1)
		} else {
			return nil // folha
		}
	} else if ponto.valorDoEixo(eixoAtual) < no.ponto.valorDoEixo(eixoAtual) {
		no.esq = removeRecursivo(no.esq, ponto, profundidade + 1)
	} else {
		no.dir = removeRecursivo(no.dir, ponto, profundidade + 1)
	}

	return no
}

func acharMinimo(no *No, eixo int, profundidade int) *No {
	if no == nil {
		return nil
	}

	eixoAtual := profundidade % DIMENSOES

	if eixoAtual == eixo {
		if no.dir == nil {
			return nil
		}
		return acharMinimo(no.esq, eixo, profundidade+1)
	}

	minimo := menorNo(no, acharMinimo(no.esq, eixo, profundidade+1), eixo)
	return menorNo(minimo, acharMinimo(no.dir, eixo, profundidade+1), eixo)
}

func menorNo(a, b *No, eixo int) *No {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if b.ponto.valorDoEixo(eixo) < a.ponto.valorDoEixo(eixo) {
		return b
	}
	return a
}


// Busca procura um valor descendo por um único caminho: em cada nó, a
// comparação descarta uma das duas subárvores inteiras.
// func Busca(a *No, procurado int) bool {
// 	if a == nil {
// 		return false
// 	}
// 	if procurado < a.Info {
// 		return Busca(a.esq, procurado)
// 	}
// 	if procurado > a.Info {
// 		return Busca(a.dir, procurado)
// 	}
// 	return true
// }