# Compilador (LALG *pseudoPython* → código objeto)

Este repositório contém um compilador didático para uma linguagem simples (LALG), dividido em etapas clássicas:

## Estrutura do projeto

- **lexer/** – analisador léxico (tokens)
- **sintatico/** – analisador sintático (parser) + chamadas ao gerador
- **semantico/** – tabela de símbolos (variáveis e funções)
- **gerador_codigo/** – emissão das instruções
- **executor/** – máquina virtual (executa o código objeto)

Arquivos de referência em **descricao/**:

- `lalg-python.txt` – gramática
- `correto.python.txt` – programa fonte de exemplo
- `codigo.objeto.txt` – exemplo de código objeto gerado

## Como rodar

Precisa passar pelo menos o arquivo. Os modos são:

1. **Só compilar** – gera o `codigo.objeto.txt` na pasta do projeto:
   ```
   go run . descricao/correto.python.txt
   ```

2. **Compilar e já executar** – compila e roda o programa (entrada/saída no terminal):
   ```
   go run . descricao/correto.python.txt run
   ```

3. **Só executar** – quando você já tem o `codigo.objeto.txt` e quer só rodar:
   ```
   go run . descricao/codigo.objeto.txt --exec
   ```

4. **Gerar lexemas** – além de compilar, gera o `lexemas.txt`:
   ```
   go run . descricao/correto.python.txt lexemas
   ```

Resumo: primeiro argumento é sempre o arquivo (fonte .python.txt ou objeto). O segundo pode ser `run`, `lexemas` ou `--exec` (só com arquivo objeto).
