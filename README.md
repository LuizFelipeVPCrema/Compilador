# Compilador (LALG *pseudoPython* → código objeto)

Este repositório contém um compilador didático para uma linguagem simples (LALG), dividido em etapas clássicas:

## Estrutura do projeto

- **lexer/** – analisador léxico (tokens)
- **sintatico/** – analisador sintático (parser) + chamadas ao gerador
- **semantico/** – tabela de símbolos (variáveis e funções)
- **gerador_codigo/** – emissão das instruções

Arquivos de referência em **descricao/**:

- `lalg-python.txt` – gramática
- `correto.python.txt` – programa fonte de exemplo
- `codigo.objeto.txt` – exemplo de código objeto gerado

## Como rodar

```bash
go run . <arquivo_fonte>
```

Exemplo: `go run . descricao/correto.python.txt`  
Gera **codigo.objeto.txt** na raiz.

## O que falta fazer

1. **Executador (máquina virtual)**  
   Nenhum módulo lê o codigo.objeto.txt e executa as instruções. Falta: ler o arquivo, interpretar cada instrução, manter memória e pilha, executar até PARA.

2. **Conferência com o gabarito**  
   O codigo.objeto.txt gerado pode ainda divergir em detalhes do `descricao/codigo.objeto.txt`; vale comparar e ajustar se necessário (por exemplo formato de números ou ordem de instruções em casos de borda).
