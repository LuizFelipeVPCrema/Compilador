# Compilador (LALG *pseudoPython* → código objeto)

Este repositório contém um compilador didático para uma linguagem simples (LALG), dividido em etapas clássicas:

## Estrutura do projeto

- **lexer/** – analisador léxico (tokens)
- **sintatico/** – analisador sintático (parser) + chamadas ao gerador
- **semantico/** – tabela de símbolos (variáveis e funções)
- **gerador_codigo/** – emissão das instruções (INPP, ALME, CRCT, CRVL, ARMZ, SOMA, SUBT, MULT, DIVI, LEIT, IMPR, DSVF, DSVI, PUSHER, PARAM, CHPR, RTPR, DESM, PARA)

Arquivos de referência em **descricao/**:

- `lalg-python.txt` – gramática
- `correto.python.txt` – programa fonte de exemplo
- `codigo.objeto.txt` – exemplo de código objeto gerado

## Como rodar

```bash
go run . <arquivo_fonte>
```

Exemplo:

```bash
go run . descricao/correto.python.txt
```

Gera o arquivo **codigo.objeto.txt** na raiz do projeto.

## O que já foi feito

1. **Integração main → lexer, tabela, gerador, parser**  
   O `main.go` monta o pipeline e o parser chama o gerador durante a análise.

2. **Declarações e comandos**  
   Programa, Corpo, DC (declarações) e Comandos seguem a gramática. Variáveis globais: `DeclararVariavel` + `GerarAlocacao(1)`.

3. **Expressões**  
   Fator emite CRCT (número) ou CRVL (variável). Termo emite MULT/DIVI. Expressão emite SOMA/SUBT.

4. **Atribuição**  
   Após avaliar a expressão, chama `GerarAtribuicao(endereco)` (ARMZ).

5. **print(expressao)**  
   Avalia a expressão e emite IMPR.

6. **if e while**  
   Usam `GerarIf` e `GerarWhile`. Condição gera comparação (CMAI, CPMI, etc.) + DSVF/DSVI. Emissor tem buffer, contador de linha e backpatching (EmitirComMarcador, Preencher) para preencher os alvos dos saltos depois.

7. **Funções**  
   - Declaração: registra função na tabela, emite DSVI para pular o corpo, grava endereço de entrada com `AtualizarEnderecoFuncao`; parâmetros com `DeclararParametro` (endereços base, base+1, …); ao sair do corpo, emite DESM n e RTPR.  
   - Chamada: emite PUSHER (linha de retorno), PARAM para cada argumento (endereço), CHPR (endereço da função).  
   - Tabela: um único mapa; ao sair da função, remove só os nomes da função atual (`adicionadosNaFuncao`) para não misturar main com corpo de função.

8. **Indentação de blocos**  
   Dentro de função/if/while, o bloco só continua se a próxima linha tiver indentação maior ou igual à do bloco. Se tiver menos (ou nenhuma), o bloco termina. Assim o corpo principal do programa precisa ter **menos indentação** que o corpo das funções (ex.: main sem espaço, funções com 2 espaços).


## O que falta fazer

1. **Executador (máquina virtual)**  
   Nenhum módulo ainda lê o `codigo.objeto.txt` e executa as instruções. Falta implementar: ler o arquivo, interpretar cada linha (INPP, ALME, CRCT, CRVL, ARMZ, SOMA, SUBT, MULT, DIVI, LEIT, IMPR, DSVF, DSVI, PUSHER, PARAM, CHPR, RTPR, DESM, PARA), manter memória e pilha, executar até PARA.

2. **Bater com o gabarito**  
   O `descricao/codigo.objeto.txt` não emite código para inicializações com 0.0 (só ALME). O compilador atual emite CRCT 0.0 (e ARMZ onde aplicável). Para ficar igual ao gabarito seria preciso uma regra extra: quando a expressão for só 0 ou 0.0, emitir só ALME e não emitir CRCT/ARMZ. E o programa fonte do main deve ter menos indentação que o corpo das funções (ver `CORRECOES_CODIGO_OBJETO.md` se existir).
