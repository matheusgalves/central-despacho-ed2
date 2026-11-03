# central-despacho-ed2

## Decisões que o seu projeto precisa tomar
- como achar, a cada turno, os chamados abertos próximos de cada ambulância;
R: Implementar regiões, por primeiro uma grade fixa de 50x50. Parar a busca quando a menor distancia possível de uma região no entorno for maior que o melhor achado
- como descartar os chamados que expiraram sem varrer todos eles;
- como remover um chamado que outra equipe recolheu, sabendo só o identificador;
- como evitar que duas ambulâncias da sua equipe sigam para o mesmo chamado;
- como achar o hospital mais próximo de uma ambulância ocupada.

OBS: O seu bot tem de saber sozinho, pelo prazo, quais chamados ainda valem.

TABELA HASH
é útil para pesquisar e armazenar elementos e tem processamento constante;
perde força se precisar de algo ordenado.

dado -> função hash -> número -> guarda no array[numero] (se colidir, lista ligada no elemento)

ÁRVORES



OBS
chamdos são árvore X dividir mapa em grades
- parar quando a menor distancia possivel das regiões que sobraram for maior que o melhor achado