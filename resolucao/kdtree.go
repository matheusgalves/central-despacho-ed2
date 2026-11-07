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

func (ponto Ponto) manhattan(pontoB Ponto) int {
	return abs(ponto.X-pontoB.X) + abs(ponto.Y-pontoB.Y)
}

func (chamado Chamado) ponto() Ponto {
	return Ponto{X: chamado.X, Y: chamado.Y}
}

func (ambulancia Ambulancia) ponto() Ponto {
	return Ponto{X: ambulancia.X, Y: ambulancia.Y}
}

type No struct {
	chamado Chamado
	esq     *No
	dir     *No
}

type KDTree struct {
	raiz *No
}

const DIMENSOES = 2

// buscaMelhor devolve o chamado mais próximo (Manhattan) de consulta.
// O segundo retorno é false se a árvore estiver vazia.
func (arvore *KDTree) buscaMelhor(consulta Ponto) (Chamado, bool) {
	// por enquanto esse aqui só considera chamados mais perto
	melhor := buscaMaisProximoRecursivo(arvore.raiz, consulta, 0, nil)
	if melhor == nil {
		return Chamado{}, false
	}
	return melhor.chamado, true
}

func (arvore *KDTree) insere(chamado Chamado) {
	arvore.raiz = insereRecursivo(arvore.raiz, chamado, 0)
}

// remove retira o chamado de mesmo ID (a posição guia a descida).
func (arvore *KDTree) remove(chamado Chamado) {
	arvore.raiz = removeRecursivo(arvore.raiz, chamado, 0)
}

func buscaMaisProximoRecursivo(no *No, consulta Ponto, profundidade int, melhor *No) *No {
	if no == nil {
		return melhor
	}

	// 1. o nó atual é candidato?
	if melhor == nil || no.chamado.ponto().manhattan(consulta) < melhor.chamado.ponto().manhattan(consulta) {
		melhor = no
	}

	// 2. desce primeiro pelo lado onde a consulta está
	eixo := profundidade % DIMENSOES
	diferenca := consulta.valorDoEixo(eixo) - no.chamado.ponto().valorDoEixo(eixo)

	proximo, oposto := no.esq, no.dir
	if diferenca >= 0 {
		proximo, oposto = no.dir, no.esq
	}

	melhor = buscaMaisProximoRecursivo(proximo, consulta, profundidade+1, melhor)

	// 3. só visita o outro lado se o plano de corte estiver mais perto
	// do que o melhor já encontrado
	if abs(diferenca) < melhor.chamado.ponto().manhattan(consulta) {
		melhor = buscaMaisProximoRecursivo(oposto, consulta, profundidade+1, melhor)
	}

	return melhor
}

func insereRecursivo(no *No, chamado Chamado, profundidade int) *No {
	if no == nil {
		return &No{chamado: chamado}
	}

	eixoAtual := profundidade % DIMENSOES

	if chamado.ponto().valorDoEixo(eixoAtual) < no.chamado.ponto().valorDoEixo(eixoAtual) {
		no.esq = insereRecursivo(no.esq, chamado, profundidade+1)
	} else {
		no.dir = insereRecursivo(no.dir, chamado, profundidade+1)
	}

	return no
}

func removeRecursivo(no *No, chamado Chamado, profundidade int) *No {
	if no == nil {
		return nil
	}

	eixoAtual := profundidade % DIMENSOES

	if no.chamado.ID == chamado.ID {
		if no.dir != nil {
			minimo := acharMinimo(no.dir, eixoAtual, profundidade+1)
			no.chamado = minimo.chamado
			no.dir = removeRecursivo(no.dir, minimo.chamado, profundidade+1)
		} else if no.esq != nil {
			minimo := acharMinimo(no.esq, eixoAtual, profundidade+1)
			no.chamado = minimo.chamado
			no.dir = removeRecursivo(no.esq, minimo.chamado, profundidade+1)
			no.esq = nil
		} else {
			return nil // folha
		}
	} else if chamado.ponto().valorDoEixo(eixoAtual) < no.chamado.ponto().valorDoEixo(eixoAtual) {
		no.esq = removeRecursivo(no.esq, chamado, profundidade+1)
	} else {
		no.dir = removeRecursivo(no.dir, chamado, profundidade+1)
	}

	return no
}

func acharMinimo(no *No, eixo int, profundidade int) *No {
	if no == nil {
		return nil
	}

	eixoAtual := profundidade % DIMENSOES

	if eixoAtual == eixo {
		// o mínimo só pode estar à esquerda; sem esquerda, é o próprio nó
		if no.esq == nil {
			return no
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
	if b.chamado.ponto().valorDoEixo(eixo) < a.chamado.ponto().valorDoEixo(eixo) {
		return b
	}
	return a
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
