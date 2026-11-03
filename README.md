# central-despacho-ed2

# Decisões que o seu projeto precisa tomar

obs: tabela hash com id do chamado para as coordenadas (x, y) na memória, para apagar rapidamente

## Achar os chamados abertos próximos de cada ambulância

Estrutura: grade de blocos. O mapa é dividido em blocos de B x B células.
Um slice de blocos, indexado por (y/B)*blocosPorLinha + x/B, guarda em
cada posição um slice com os chamados daquele bloco.

Custo por turno:
- Inserir: calcular o bloco da célula e fazer append no slice dele. O(1) amortizado.
- Remover: o bloco vem da posição do chamado. Achar o chamado
  no bloco custa O(chamados por bloco); remover com swap-remove é O(1).
- Achar o chamado mais próximo, em anéis de blocos (anel 0 = bloco da
  ambulância, anel 1 = os 8 blocos vizinhos, e assim por diante):
  1) calcular o bloco da ambulância;
  2) abrir o anel k, examinar os chamados de cada bloco, descartar os que
     não dão tempo (distância > prazo - t) e guardar o de menor distância;
  3) parar se melhorDistancia <= k*B + 1 (todo chamado ainda não visto
     está a pelo menos k*B + 1 de distância);
  4) parar também se o anel passar das bordas do mapa.
  Custo: O(blocos visitados x chamados por bloco).

Escolha de B: mirar em dezenas de chamados por bloco (no estresse,
B = 16 dá ~32). Chamados concentrados deixam alguns blocos mais cheios
que a média, o que aumenta o custo de examiná-los.

Alternativa descartada: AVL ordenada por x. A consulta por faixa visita
colunas inteiras do mapa, e o custo por ambulância cresce com o número de
chamados na coluna, o que deve estourar o p99 no estresse.


## Descartar os chamados que expiraram sem varrer todos eles
Estrutura: lista por prazo (slice indexado pelo turno do prazo), cada
posição com os chamados que expiram naquele turno.

Custo por turno: inserir um chamado é O(1) (append na lista do seu
prazo). No turno t, processar a lista do prazo t-1: para cada chamado,
achar o bloco pela posição e remover com swap-remove (O(chamados por
bloco) para localizar, O(1) para remover). Turnos sem expiração custam
O(1). Chamados já recolhidos são ignorados.

Alternativa descartada: remover só quando a ambulância encontra o
chamado vencido durante a busca. É barato em média, mas blocos fora do
alcance das ambulâncias acumulam mortos (memória) e a limpeza vira um
pico de tempo no turno em que o bloco é visitado, o que prejudica o p99.



## Remover um chamado que outra equipe recolheu, sabendo só o id
Estrutura: tabela hash própria (encadeamento), de id para (x, y) do
chamado. Sem usar map.

Custo por turno: para cada id de RECOLHIDOS, consultar a tabela
(O(1) esperado), calcular o bloco pelo (x, y), achar o id dentro do
bloco e remover com swap-remove (O(chamados por bloco) para achar, O(1)
para remover) e apagar a entrada da tabela (O(1) esperado). Se o id não
estiver na tabela (já removido), ignora. Ao final, ambulâncias que o alvo
foi recolhido ficam sem alvo e buscam outro

Alternativa descartada: guardar (bloco, índice) na tabela. A remoção
seria O(1), mas o swap-remove mudaria o índice do último chamado do bloco
e exigiria atualizá-lo na tabela a cada remoção.

## Evitar que duas ambulâncias da equipe sigam para o mesmo chamado

Estrutura: campo "reservado" (id da ambulância dona) em cada chamado,
dentro do bloco, mais o alvo atual de cada ambulância (id, x, y).

Custo por turno: reservar e liberar são O(1), pois o chamado já foi
localizado na busca. A busca ignora chamados reservados por outra
ambulância. As ambulâncias são processadas uma por vez, então a segunda
já enxerga a reserva da primeira. A reserva cai quando o chamado é
removido do bloco ou quando a ambulância troca de alvo.

Alternativa descartada: conjunto separado de ids reservados (tabela
hash), que exigiria uma consulta e uma estrutura a mais para manter
consistente com os blocos, sem ganho de custo.

## Achar o hospital mais próximo de uma ambulância ocupada

Estrutura: slice (lista) com as coordenadas dos hospitais, lido no INICIO.

Custo por turno: quando uma ambulância passa a ocupada, varrer os H
hospitais (H <= 20) e guardar o de menor distância como
alvo dela. O(H) por ambulância que acabou de recolher, e nada nos demais
turnos, pois na grade aberta o hospital escolhido continua sendo o mais
próximo durante o trajeto.

Alternativa descartada: grade (como fizemos com os chamamados) para os hospitais. Com H <= 20
a varredura é mais simples e não mais lenta; uma estrutura espacial só
compensaria com muitos hospitais.