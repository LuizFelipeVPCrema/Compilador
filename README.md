# Compilador (LALG *pseudoPython* → código objeto)

Este repositório contém um compilador didático para uma linguagem simples (LALG), dividido em etapas clássicas:

- **Léxico** (`lexer/`): lê o arquivo fonte e gera uma tabela de tokens em `lexemas.txt`.
- **Sintático** (`sintatico/`): consome tokens e valida a estrutura do programa.
- **Semântico** (`semantico/`): mantém a tabela de símbolos (variáveis/funções) e validações associadas.
- **Gerador de código intermediário/objeto** (`gerador_codigo/`): emite instruções de uma máquina virtual.

Arquivos de referência importantes do projeto:

- `lexemas.txt`: saída do analisador léxico (tokens).
- `descricao/codigo.objeto.txt`: gabarito/exemplo de código objeto esperado.
- `descricao/correto.python.txt`: exemplo do programa fonte usado para gerar `lexemas.txt`.
- `descricao/lalg-python.txt`: especificação (gramática) da linguagem.

## Como rodar (estado atual)

No estado atual, o `main.go` executa **somente o léxico** e gera `lexemas.txt`:

```bash
go run . <arquivo_entrada>
```

Exemplo:

```bash
go run . descricao/correto.python.txt
```

Isso deve gerar/atualizar `lexemas.txt`.

## Oque precisa ser feito (para concluir a geração do `codigo.objeto`)

O projeto ainda precisa “fechar o ciclo” completo: **fonte → tokens → parsing/semântica → emissão de código objeto**.

Checklist objetivo do que falta implementar/ajustar:

1. **Integrar o parser com o gerador de código**
   - Hoje o `sintatico/sintatico.go` consome tokens para validar a gramática, mas não chama o `gerador_codigo` enquanto parseia.
   - Próximo passo: durante `Expressao`, `Termo`, `Fator`, `AtribuicaoOuChamada`, `Imprimir`, `Condicional`, `While`, `Funcao` etc., chamar as rotinas do `Gerador` para emitir o código correspondente.

2. **Evoluir a tabela de símbolos para suportar endereços e escopo**
   - O gabarito (`descricao/codigo.objeto.txt`) usa endereços (`ARMZ 12`, `CRVL 9`, etc.) e também precisa distinguir **variáveis globais, locais e parâmetros**.
   - Próximo passo: armazenar no símbolo pelo menos: `Endereco`, `Escopo` (global/função) e metadados de função (endereço de entrada da função e quantidade de parâmetros).

3. **Criar as condicionais saltos/rótulos do `condicionais.go`**
   - Criar `GerarIf`/`GerarWhile`.


4. **Mapear `read()` e `input()` para `LEIT`**
   - Criar as funções de `read()` e `input()`
   - Próximo passo: no parser, ao reconhecer `read()`/`input()` como expressão, emitir `LEIT` e depois, na atribuição, `ARMZ`.

5. **Funções: entrada, parâmetros, retorno e desalocação**
   - O gabarito usa `PUSHER`, `PARAM`, `CHPR`, `RTPR` e `DESM` no final das funções.
   - Próximo passo: definir claramente o protocolo de chamada e emitir a sequência correta durante:
     - declaração de função (pular corpo no início do programa, depois entrar no corpo quando chamada)
     - chamada de função (empilhar retorno, parâmetros, `CHPR`)
     - final de função (`DESM n` + `RTPR`)


