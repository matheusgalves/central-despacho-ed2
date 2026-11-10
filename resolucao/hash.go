package main

type hashChamados struct {
	chaves  []int
	valores []Ponto
	ocupada []bool
	n       int
}

func novaHashChamados(m int) hashChamados {
	return hashChamados{
		chaves:  make([]int, m),
		valores: make([]Ponto, m),
		ocupada: make([]bool, m),
	}
}

func (t *hashChamados) m() int {
	return len(t.chaves)
}

func (t *hashChamados) hash(k int) int {
	h := k % t.m()
	if h < 0 {
		h += t.m()
	}
	return h
}

func (t *hashChamados) distancia(de, ate int) int {
	return (ate - de + t.m()) % t.m()
}

func (t *hashChamados) posicao(k int) int {
	p := t.hash(k)
	for t.ocupada[p] && t.chaves[p] != k {
		p = (p + 1) % t.m()
	}
	return p
}

func (t *hashChamados) insere(k int, v Ponto) {
	if t.m() == 0 {
		*t = novaHashChamados(1031)
	}

	p := t.posicao(k)
	if t.ocupada[p] {
		t.valores[p] = v
		return
	}

	t.chaves[p] = k
	t.valores[p] = v
	t.ocupada[p] = true
	t.n++

	if 10*t.n > 7*t.m() {
		t.redimensiona(proximoPrimo(2 * t.m()))
	}
}

func (t *hashChamados) redimensiona(m int) {
	nova := novaHashChamados(m)
	for i, k := range t.chaves {
		if t.ocupada[i] {
			nova.insere(k, t.valores[i])
		}
	}
	*t = nova
}

func (t *hashChamados) retira(k int) (Ponto, bool) {
	if t.m() == 0 {
		return Ponto{}, false
	}
	p := t.posicao(k)
	if !t.ocupada[p] {
		return Ponto{}, false
	}

	valor := t.valores[p]
	t.ocupada[p] = false
	t.n--

	buraco := p
	for p = (p + 1) % t.m(); t.ocupada[p]; p = (p + 1) % t.m() {
		if t.distancia(t.hash(t.chaves[p]), p) >= t.distancia(buraco, p) {
			t.chaves[buraco] = t.chaves[p]
			t.valores[buraco] = t.valores[p]
			t.ocupada[buraco] = true
			t.ocupada[p] = false
			buraco = p
		}
	}
	return valor, true
}

func proximoPrimo(k int) int {
	for {
		primo := k >= 2
		for d := 2; d*d <= k; d++ {
			if k%d == 0 {
				primo = false
				break
			}
		}
		if primo {
			return k
		}
		k++
	}
}
