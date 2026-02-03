package gerador_codigo

import (
	"fmt"
	"os"
	"strconv"
)

type Emissor struct {
	instrucoes []string
	arquivo    *os.File
	caminho    string
}

func NovoEmissor(caminho string) *Emissor {
	return &Emissor{
		instrucoes: nil,
		arquivo:    nil,
		caminho:    caminho,
	}
}

func (e *Emissor) Emitir(instrucao string) {
	e.instrucoes = append(e.instrucoes, instrucao)
}

func (e *Emissor) LinhaAtual() int {
	return len(e.instrucoes)
}

func (e *Emissor) ProximaLinha() int {
	return len(e.instrucoes) + 1
}

func (e *Emissor) EmitirComMarcador(prefixo string) int {
	e.instrucoes = append(e.instrucoes, prefixo+" 0")
	return len(e.instrucoes) - 1
}

func (e *Emissor) Preencher(indice int, linhaAlvo int, prefixo string) {
	if indice >= 0 && indice < len(e.instrucoes) {
		e.instrucoes[indice] = prefixo + " " + strconv.Itoa(linhaAlvo)
	}
}

func (e *Emissor) Fechar() {
	e.arquivo, _ = os.Create(e.caminho)
	for _, inst := range e.instrucoes {
		fmt.Fprintln(e.arquivo, inst)
	}
	e.arquivo.Close()
}
