# central-despacho-ed2

# Decisões que o seu projeto precisa tomar

## Achar os chamados abertos próximos de cada ambulância

Estrutura: duas KD-Trees de 2 dimensões (x, y), uma para os chamados
graves (valor >= 40) e outra para os comuns. A KD-Tree é uma árvore
binária de busca que alterna a coordenada comparada a cada nível: a raiz
separa por x, os filhos por y, os netos por x de novo, e assim por diante.

A separação em duas árvores vem do critério de escolha: a ambulância vai
ao chamado de maior nota, valor / (distância até o chamado + distância do
chamado ao hospital mais próximo + 1). Um grave (60 a 100) vale
mais que vários comuns (1 a 20) mesmo estando mais longe, mas o mais
próximo numa árvore única quase nunca é grave. Com duas árvores,
cada busca devolve direto o comum e o grave mais próximos e desempatamos os dois pela nota.

Custo por turno:
- Inserir: descer da raiz comparando a coordenada do nível até achar um
  filho vazio. O(profundidade), O(log n) esperado. Medido no estresse
  (~1 milhão de nós): profundidade média 27, o esperado para uma árvore
  montada em ordem aleatória.
- Remover (recolhidos e graves vencidos): descer pela posição até achar o
  id, O(log n) esperado; se o nó tem filhos, ele é trocado pelo mínimo da
  coordenada do nível numa das subárvores, o que pode visitar parte dela.
  Os comuns vencidos não são removidos (ver a seção de expiração).
- Achar o chamado mais próximo de uma ambulância:
  1) descer primeiro pelo lado da consulta, guardando o melhor candidato
     (prazo ainda alcançável e não reservado neste turno);
  2) na volta, só visitar o outro lado se a distância até que o corte esperado
     seja menor que a do melhor candidato (poda);
  3) na árvore de graves, a busca começa com um limite de distância: se o
     melhor comum tem nota n, um grave a mais de 100/n células nunca ganha
     dele, então esses ramos nem são visitados.
  Custo: O(log n) típico (medido: ~45 nós visitados por busca). O pior
  caso é O(n), quando nenhum nó passa nos filtros e a poda nunca dispara;
  o limite e a remoção dos graves vencidos vão evitar esse caso na prática.

Alternativa descartada: BST comum (ou AVL) ordenada por uma coordenada,
por exemplo x. A distância de Manhattan depende das duas coordenadas: um
chamado com x próximo pode estar longe em y. A BST só consegue restringir
a busca a uma faixa |x - x_ambulância| <= d e precisa examinar todos os
chamados dessa faixa, que no estresse são colunas inteiras do mapa com
milhares de chamados. O custo por ambulância cresceria com o número de
chamados na faixa e estouraria o p99. A KD-Tree mantém a inserção e a
descida de uma BST, mas, ao alternar x e y, cada nível corta o plano ao
meio numa direção diferente, e a poda funciona nas duas dimensões.


## Descartar os chamados que expiraram sem varrer todos eles

Estrutura: lista por prazo (slice criado no INICIO com turnos + 1
posições, indexado pelo turno do prazo), cada posição com os ids que
vencem naquele turno.

Custo por turno: inserir um chamado é O(1) amortizado (append na lista do
seu prazo). No turno t, processar a lista do prazo t-1: para cada id,
retirá-lo da tabela hash (O(1) esperado; se não estiver lá, já foi
recolhido e é ignorado) e removê-lo da árvore de graves (O(log n)
esperado, e a árvore é pequena; se o chamado era comum, a remoção não
acha o id e custa só a descida). Turnos sem expiração custam O(1).

Os comuns vencidos continuam na árvore, pois a busca
já descarta chamados com prazo vencido, então eles não mudam nenhuma
decisão. Com isso, economizamos remoções na árvore, que custam mais. Os graves precisam ser removidos, pois como têm prazo curto, a árvore de
graves viraria quase só chamados mortos, e uma busca sem candidato válido percorreria dezenas de milhares de nós.

Alternativa descartada: remover os vencidos também da árvore de comuns.
Medido no estresse e com ~1000 vencimentos por turno, a remoção (que troca
o nó pelo mínimo de uma subárvore) era um terço do processador e gerava
turnos de 6 a 11 ms, o que prejudicava o p99.

## Remover um chamado que outra equipe recolheu, sabendo só o id

Estrutura: tabela hash própria, de id para (x, y) do chamado, com
endereçamento aberto e sondagem linear. Sem usar map. A remoção desloca
as chaves seguintes para trás (sem lápides), e a tabela cresce para o
próximo primo acima do dobro quando passa de 70% de ocupação.

Custo por turno: para cada id de RECOLHIDOS, retirá-lo da tabela (O(1)
esperado) e, com o (x, y), remover o chamado da KD-Tree descendo pela
posição até achar o id (O(log n) esperado). A tabela não guarda se o
chamado é grave, então a remoção é tentada nas duas árvores; na errada
custa só a descida. Se o id não estiver na tabela, ignora. Como as
ambulâncias reavaliam o alvo a cada turno, não há mais nada a atualizar.

Alternativa descartada: guardar na tabela um ponteiro para o nó da
árvore. Quando um nó com filhos é removido, o conteúdo do mínimo da
subárvore é copiado para ele, então o ponteiro de outro chamado passaria
a apontar para o lugar errado. A posição (x, y) não muda e sempre guia a
descida até o nó certo.

## Evitar que duas ambulâncias da equipe sigam para o mesmo chamado

Estrutura: campo de reserva em cada chamado (no próprio nó da árvore),
mais uma lista com os chamados reservados no turno.

Custo por turno: as ambulâncias são processadas uma por vez. O chamado
escolhido é marcado (O(1), pois a busca devolve um ponteiro para o nó),
e as buscas das ambulâncias seguintes ignoram os marcados. No fim do
turno, todas as marcas são apagadas (O(ambulâncias)) e, no turno
seguinte, tudo é reavaliado.

Alternativa descartada: guardar o alvo de cada ambulância entre turnos.
Medido no medio, essa versão fez ~30% menos pontos que o bot ingênuo, pois a
ambulância seguia para o alvo antigo e ignorava os chamados novos que
surgiam mais perto no caminho.

## Achar o hospital mais próximo de uma ambulância ocupada

Estrutura: slice (lista) com as coordenadas dos hospitais, lido no INICIO.

Custo por turno: varrer os H hospitais (H <= 20) e escolher o de menor
distância, O(H). Isso é feito para cada ambulância ocupada e, na nota,
para os dois candidatos (comum e grave) de cada ambulância livre. No
estresse são no máximo 100 ambulâncias x 2 x 20 comparações por turno.

Alternativa descartada: calcular o hospital mais próximo de cada chamado
na inserção. Também é O(H), mas para os ~1000 chamados novos de cada turno
do estresse, enquanto a varredura sob demanda só olha os candidatos que
realmente são avaliados.
