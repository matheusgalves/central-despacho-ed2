// Bot da etapa 1: despacha ambulâncias na grade aberta.
package main

import (
	"bufio"
	"fmt"
	"os"
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

var (
	leitor   *bufio.Scanner
	escritor *bufio.Writer

	hospitais      []Hospital
	chamadosComuns KDTree
	chamadosGraves KDTree

	chamadosPorID hashChamados
	expiraEm      [][]int
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

		chamadosRecolhidos := inteiro(lerCampos()[1]) // RECOLHIDOS r
		for i := 0; i < chamadosRecolhidos; i++ {
			idChamado := inteiro(lerCampos()[0])

			removerPorID(idChamado)
		}

		removerExpirados(turno)

		// é mais simples atualizar as ambulancias 1x por turno
		ambulanciasServidor := inteiro(lerCampos()[1]) // AMBULANCIAS k
		// reaproveita o mesmo slice a cada turno em vez de alocar um novo
		if cap(bufAmbulancias) < ambulanciasServidor {
			bufAmbulancias = make([]Ambulancia, ambulanciasServidor)
		}
		ambulancias := bufAmbulancias[:ambulanciasServidor]

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

var bufAmbulancias []Ambulancia

func jogarTurno(turno int, ambulancias []Ambulancia) {
	for _, ambulancia := range ambulancias {
		mov := "P"
		if ambulancia.Livre {
			melhorChamadoComum := chamadosComuns.buscaMelhor(ambulancia.ponto(), turno)
			melhorChamadoGrave := chamadosGraves.buscaMelhor(ambulancia.ponto(), turno)

			escolhido := melhorDosMelhores(melhorChamadoComum, melhorChamadoGrave, ambulancia)
			if escolhido != nil {
				mov = direcao(ambulancia.X, ambulancia.Y, escolhido.X, escolhido.Y)
			}
		} else {
			hospital := hospitalMaisProximo(ambulancia.ponto())
			mov = direcao(ambulancia.X, ambulancia.Y, hospital.X, hospital.Y)
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
		Livre: string(campos[3]) == "LIVRE", // a conversão só para comparar não aloca
	}
}

// camposBuf é reaproveitado por lerCampos a cada linha.
var camposBuf [][]byte

// lerCampos lê a próxima linha da entrada e a divide em campos, sem alocar:
// cada campo é uma fatia dos bytes do próprio Scanner (leitor.Bytes() não
// copia a linha, ao contrário de leitor.Text()). Por isso o resultado só vale
// até a próxima chamada; converta os campos antes de ler outra linha.
// Devolve nil quando a entrada termina.
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

// inteiro converte os bytes de um número decimal (com sinal opcional).
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
	for i := 0; i < k; i++ {
		lerCampos() // id x y
	}
	b := inteiro(lerCampos()[1])
	for i := 0; i < b; i++ {
		lerCampos() // x1 y1 x2 y2
	}
	lerCampos() // FIM
}

// removerPorID usa as coordenadas da hash para localizar o chamado na árvore
func removerPorID(id int) {
	ponto, encontrado := chamadosPorID.retira(id)
	if !encontrado {
		return
	}
	chamadosGraves.remove(Chamado{ID: id, X: ponto.X, Y: ponto.Y})
	chamadosComuns.remove(Chamado{ID: id, X: ponto.X, Y: ponto.Y})
}

// removerExpirados processa o lote que venceu ao fim do turno anterior.
// O ID sempre sai da hash. Da árvore, só os graves são removidos de fato:
// a árvore de graves é pequena (remover é barato) e, cheia de mortos, faria
// a busca sem candidato percorrer milhares de nós. Os comuns vencidos ficam
// na árvore (remoção preguiçosa): a busca já os descarta pelo prazo, e
// remover da árvore grande de comuns custava caro.
func removerExpirados(turno int) {
	prazoVencido := turno - 1
	for _, id := range expiraEm[prazoVencido] {
		if ponto, encontrado := chamadosPorID.retira(id); encontrado {
			// se o chamado era comum, remove não acha o ID e não faz nada
			chamadosGraves.remove(Chamado{ID: id, X: ponto.X, Y: ponto.Y})
		}
	}
	expiraEm[prazoVencido] = nil
}

// melhorDosMelhores escolhe entre o melhor chamado comum e o melhor grave pela
// nota. Qualquer um dos dois pode ser nil; devolve nil se ambos forem.
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
	limiar := 40 // chamados graves: 60 a 100 pontos

	if chamado.Valor >= limiar {
		return true
	}

	return false
}
