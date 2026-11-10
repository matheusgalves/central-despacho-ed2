package main

import (
	"bufio"
	"fmt"
	"os"
)

const limiarGrave = 40

type Chamado struct {
	ID    int
	X, Y  int
	Valor int
	Prazo int
}

type Ambulancia struct {
	ID    int
	X, Y  int
	Livre bool
}

type Hospital struct {
	X, Y int
}

type destinoHospital struct {
	hospital Hospital
	definido bool
}

var (
	leitor   *bufio.Scanner
	escritor *bufio.Writer

	hospitais      []Hospital
	chamadosComuns KDTree
	chamadosGraves KDTree

	chamadosPorID hashChamados
	expiraEm      [][]int
	destinos      []destinoHospital

	bufAmbulancias []Ambulancia
	camposBuf      [][]byte
)

func main() {
	leitor = bufio.NewScanner(os.Stdin)
	leitor.Buffer(make([]byte, 1<<16), 1<<20)
	escritor = bufio.NewWriter(os.Stdout)

	lerInicio()
	fmt.Fprintln(escritor, "PRONTO")
	escritor.Flush()

	for {
		c := lerCampos()
		if c == nil {
			return
		}
		turno := inteiro(c[1])

		novosChamados := inteiro(lerCampos()[1])
		for i := 0; i < novosChamados; i++ {
			chamado := montaChamado(lerCampos())

			if chamadoEhGrave(chamado) {
				chamadosGraves.insere(chamado)
			} else {
				chamadosComuns.insere(chamado)
			}

			chamadosPorID.insere(chamado.ID, chamado.ponto())

			if chamado.Prazo < len(expiraEm) {
				expiraEm[chamado.Prazo] = append(expiraEm[chamado.Prazo], chamado.ID)
			}
		}

		chamadosRecolhidos := inteiro(lerCampos()[1])
		for i := 0; i < chamadosRecolhidos; i++ {
			removerPorID(inteiro(lerCampos()[0]))
		}

		removerExpirados(turno)

		ambulanciasServidor := inteiro(lerCampos()[1])
		if cap(bufAmbulancias) < ambulanciasServidor {
			bufAmbulancias = make([]Ambulancia, ambulanciasServidor)
		}
		ambulancias := bufAmbulancias[:ambulanciasServidor]
		for i := 0; i < ambulanciasServidor; i++ {
			ambulancias[i] = montaAmbulancia(lerCampos())
		}

		placar := inteiro(lerCampos()[1])
		for i := 0; i < placar; i++ {
			lerCampos()
		}

		lerCampos()

		jogarTurno(turno, ambulancias)
	}
}

func jogarTurno(turno int, ambulancias []Ambulancia) {
	for _, ambulancia := range ambulancias {
		mov := "P"
		if ambulancia.Livre {
			destinos[ambulancia.ID] = destinoHospital{}
			melhorChamadoComum := chamadosComuns.buscaMelhor(ambulancia.ponto(), turno)
			melhorChamadoGrave := chamadosGraves.buscaMelhor(ambulancia.ponto(), turno)

			escolhido := melhorDosMelhores(melhorChamadoComum, melhorChamadoGrave, ambulancia)
			if escolhido != nil {
				mov = direcao(ambulancia.X, ambulancia.Y, escolhido.X, escolhido.Y)
			}
		} else {
			destino := &destinos[ambulancia.ID]
			if !destino.definido {
				destino.hospital = hospitalMaisProximo(ambulancia.ponto())
				destino.definido = true
			}
			mov = direcao(ambulancia.X, ambulancia.Y, destino.hospital.X, destino.hospital.Y)
		}
		fmt.Fprintln(escritor, ambulancia.ID, mov)
	}
	fmt.Fprintln(escritor, "FIM")
	escritor.Flush()
}

func montaChamado(campos [][]byte) Chamado {
	return Chamado{
		ID:    inteiro(campos[0]),
		X:     inteiro(campos[1]),
		Y:     inteiro(campos[2]),
		Valor: inteiro(campos[3]),
		Prazo: inteiro(campos[4]),
	}
}

func montaAmbulancia(campos [][]byte) Ambulancia {
	return Ambulancia{
		ID:    inteiro(campos[0]),
		X:     inteiro(campos[1]),
		Y:     inteiro(campos[2]),
		Livre: string(campos[3]) == "LIVRE",
	}
}

func lerCampos() [][]byte {
	if !leitor.Scan() {
		return nil
	}
	linha := leitor.Bytes()
	camposBuf = camposBuf[:0]
	i := 0
	for i < len(linha) {
		for i < len(linha) && ehEspaco(linha[i]) {
			i++
		}
		inicio := i
		for i < len(linha) && !ehEspaco(linha[i]) {
			i++
		}
		if i > inicio {
			camposBuf = append(camposBuf, linha[inicio:i])
		}
	}
	return camposBuf
}

func ehEspaco(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r'
}

func inteiro(campo []byte) int {
	n, negativo := 0, false
	for i, c := range campo {
		if i == 0 && c == '-' {
			negativo = true
			continue
		}
		if c < '0' || c > '9' {
			panic("número inválido na entrada: " + string(campo))
		}
		n = n*10 + int(c-'0')
	}
	if negativo {
		return -n
	}
	return n
}

func direcao(x, y, alvoX, alvoY int) string {
	switch {
	case x < alvoX:
		return "L"
	case x > alvoX:
		return "O"
	case y < alvoY:
		return "S"
	case y > alvoY:
		return "N"
	}
	return "P"
}

func lerInicio() {
	inicio := lerCampos()
	expiraEm = make([][]int, inteiro(inicio[3])+1)

	h := inteiro(lerCampos()[1])
	for i := 0; i < h; i++ {
		c := lerCampos()
		hospitais = append(hospitais, Hospital{inteiro(c[0]), inteiro(c[1])})
	}

	k := inteiro(lerCampos()[1])
	destinos = make([]destinoHospital, k)
	for i := 0; i < k; i++ {
		lerCampos()
	}

	b := inteiro(lerCampos()[1])
	for i := 0; i < b; i++ {
		lerCampos()
	}

	lerCampos()
}

func removerPorID(id int) {
	ponto, encontrado := chamadosPorID.retira(id)
	if !encontrado {
		return
	}
	chamadosGraves.remove(Chamado{ID: id, X: ponto.X, Y: ponto.Y})
	chamadosComuns.remove(Chamado{ID: id, X: ponto.X, Y: ponto.Y})
}

func removerExpirados(turno int) {
	prazoVencido := turno - 1
	for _, id := range expiraEm[prazoVencido] {
		if ponto, encontrado := chamadosPorID.retira(id); encontrado {
			chamadosGraves.remove(Chamado{ID: id, X: ponto.X, Y: ponto.Y})
		}
	}
	expiraEm[prazoVencido] = nil
}

func melhorDosMelhores(chamadoComum, chamadoGrave *Chamado, ambulancia Ambulancia) *Chamado {
	if chamadoComum == nil {
		return chamadoGrave
	}
	if chamadoGrave == nil {
		return chamadoComum
	}

	if calculaNota(*chamadoGrave, ambulancia) >= calculaNota(*chamadoComum, ambulancia) {
		return chamadoGrave
	}

	return chamadoComum
}

func calculaNota(chamado Chamado, ambulancia Ambulancia) float64 {
	distanciaAteChamado := ambulancia.ponto().manhattan(chamado.ponto())
	distanciaChamadoAteHospital := hospitalMaisProximo(chamado.ponto()).ponto().manhattan(chamado.ponto())

	return float64(chamado.Valor) / float64(distanciaAteChamado+distanciaChamadoAteHospital+1)
}

func hospitalMaisProximo(ponto Ponto) Hospital {
	melhor := hospitais[0]
	melhorDistancia := ponto.manhattan(melhor.ponto())

	for _, h := range hospitais[1:] {
		if d := ponto.manhattan(h.ponto()); d < melhorDistancia {
			melhor = h
			melhorDistancia = d
		}
	}

	return melhor
}

func chamadoEhGrave(chamado Chamado) bool {
	return chamado.Valor >= limiarGrave
}
