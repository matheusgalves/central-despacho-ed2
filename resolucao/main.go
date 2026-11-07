// Bot da etapa 1: despacha ambulâncias na grade aberta.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Chamado struct {
	ID                 int
	X, Y               int
	Valor              int
	Prazo              int
	AmbulanciaACaminho bool
}

type Ambulancia struct {
	ID    int
	X, Y  int
	Livre bool
}

type Hospital struct {
	X, Y int
}

type alvoAmbulancia struct {
	id    int // zero indica que a ambulância não tem alvo
	ponto Ponto
	prazo int
}

var (
	leitor   *bufio.Scanner
	escritor *bufio.Writer

	hospitais     []Hospital
	chamados      KDTree
	chamadosPorID hashChamados
	expiraEm      [][]int
	alvos         []alvoAmbulancia
)

func main() {
	leitor = bufio.NewScanner(os.Stdin)
	leitor.Buffer(make([]byte, 1<<16), 1<<20)
	escritor = bufio.NewWriter(os.Stdout)

	lerInicio()
	fmt.Fprintln(escritor, "PRONTO")
	escritor.Flush()

	for {
		c := lerCampos() // TURNO t
		if c == nil {
			return // fim da partida
		}
		turno := inteiro(c[1])

		novosChamados := inteiro(lerCampos()[1]) // CHAMADOS n

		for i := 0; i < novosChamados; i++ {
			chamado := montaChamado(lerCampos())
			chamados.insere(chamado)
			chamadosPorID.insere(chamado.ID, chamado.ponto())
			if chamado.Prazo < len(expiraEm) {
				expiraEm[chamado.Prazo] = append(expiraEm[chamado.Prazo], chamado.ID)
			}
		}

		chamadosRecolhidos := inteiro(lerCampos()[1]) // RECOLHIDOS r
		for i := 0; i < chamadosRecolhidos; i++ {
			idChamado := inteiro(lerCampos()[0])

			removerPorID(idChamado)
		}

		removerExpirados(turno)

		// é mais simples atualizar as ambulancias 1x por turno
		ambulanciasServidor := inteiro(lerCampos()[1]) // AMBULANCIAS k
		ambulancias := make([]Ambulancia, ambulanciasServidor)

		for i := 0; i < ambulanciasServidor; i++ {
			ambulancias[i] = montaAmbulancia(lerCampos())
		}

		placar := inteiro(lerCampos()[1]) // PLACAR
		for i := 0; i < placar; i++ {
			lerCampos()
		}

		lerCampos() // FIM

		jogarTurno(turno, ambulancias)
	}
}

func montaChamado(campos []string) Chamado {
	return Chamado{
		ID:    inteiro(campos[0]),
		X:     inteiro(campos[1]),
		Y:     inteiro(campos[2]),
		Valor: inteiro(campos[3]),
		Prazo: inteiro(campos[4]),
	}
}

func montaAmbulancia(campos []string) Ambulancia {
	return Ambulancia{
		ID:    inteiro(campos[0]),
		X:     inteiro(campos[1]),
		Y:     inteiro(campos[2]),
		Livre: campos[3] == "LIVRE",
	}
}

// lerCampos lê a próxima linha da entrada e a divide em campos.
// Devolve nil quando a entrada termina.
func lerCampos() []string {
	if !leitor.Scan() {
		return nil
	}
	return strings.Fields(leitor.Text())
}

func inteiro(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		panic("número inválido na entrada: " + s)
	}
	return v
}

// TODO: ETAPA 2 - Ajustar para considerar bloqueios
func distancia(x1, y1, x2, y2 int) int {
	return abs(x1-x2) + abs(y1-y2)
}

// direcao devolve o movimento que aproxima (x, y) de (alvoX, alvoY),
// andando primeiro na horizontal.
// TODO: ETAPA 2 - Pensar em como otimizar para gerar o melhor caminho
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

// lerInicio lê a mensagem de início e prepara o estado da partida.
func lerInicio() {
	inicio := lerCampos() // INICIO largura altura turnos limite_turno limite_preparo equipe num_equipes
	expiraEm = make([][]int, inteiro(inicio[3])+1)
	h := inteiro(lerCampos()[1])
	for i := 0; i < h; i++ {
		c := lerCampos()
		hospitais = append(hospitais, Hospital{inteiro(c[0]), inteiro(c[1])})
	}
	k := inteiro(lerCampos()[1])
	alvos = make([]alvoAmbulancia, k)
	for i := 0; i < k; i++ {
		lerCampos() // id x y
	}
	b := inteiro(lerCampos()[1])
	for i := 0; i < b; i++ {
		lerCampos() // x1 y1 x2 y2
	}
	lerCampos() // FIM
}

// removerPorID usa as coordenadas da hash para localizar o chamado na árvore.
func removerPorID(id int) {
	ponto, encontrado := chamadosPorID.retira(id)
	if !encontrado {
		return
	}
	chamados.remove(Chamado{ID: id, X: ponto.X, Y: ponto.Y})
}

// removerExpirados processa o lote que venceu ao fim do turno anterior.
func removerExpirados(turno int) {
	prazoVencido := turno - 1
	for _, id := range expiraEm[prazoVencido] {
		removerPorID(id)
	}
	expiraEm[prazoVencido] = nil
}

// liberarAlvo encerra a reserva, mesmo se o chamado já tiver saído da árvore.
func liberarAlvo(idAmbulancia int) {
	alvo := alvos[idAmbulancia]
	if chamado := chamados.buscaPorID(alvo.id, alvo.ponto); chamado != nil {
		chamado.AmbulanciaACaminho = false
	}
	alvos[idAmbulancia] = alvoAmbulancia{}
}

func jogarTurno(turno int, ambulancias []Ambulancia) {
	// Libera todos os alvos inválidos antes de distribuir os novos.
	for _, a := range ambulancias {
		alvo := alvos[a.ID]
		if alvo.id == 0 {
			continue
		}
		posicao := Ponto{X: a.X, Y: a.Y}
		alcancavel := chamadoAlcancavel(turno, alvo.prazo, posicao.manhattan(alvo.ponto))
		if !a.Livre || !alcancavel || !chamadosPorID.contem(alvo.id) {
			liberarAlvo(a.ID)
		}
	}

	for _, a := range ambulancias {
		mov := "P"
		if a.Livre {
			alvo := alvos[a.ID]
			if alvo.id == 0 {
				chamado := chamados.buscaMaisProximo(Ponto{X: a.X, Y: a.Y}, turno)
				if chamado != nil {
					chamado.AmbulanciaACaminho = true
					alvo = alvoAmbulancia{id: chamado.ID, ponto: chamado.ponto(), prazo: chamado.Prazo}
					alvos[a.ID] = alvo
				}
			}
			if alvo.id != 0 {
				mov = direcao(a.X, a.Y, alvo.ponto.X, alvo.ponto.Y)
			}
		} else {
			melhor := 0
			//TODO: ETAPA 1 - Armazenar hospital mais próximo
			for i, h := range hospitais {
				if distancia(a.X, a.Y, h.X, h.Y) < distancia(a.X, a.Y, hospitais[melhor].X, hospitais[melhor].Y) {
					melhor = i
				}
			}
			mov = direcao(a.X, a.Y, hospitais[melhor].X, hospitais[melhor].Y)
		}
		fmt.Fprintln(escritor, a.ID, mov)
	}
	fmt.Fprintln(escritor, "FIM")
	escritor.Flush()
}
