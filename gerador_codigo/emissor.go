package gerador_codigo

import (
	"fmt"
	"os"
	"strconv"
)

type Emissor struct {
	instrucoes  []string
	arquivo     *os.File
	caminho     string
	bufferAtivo bool
	buffer      []string
}

func NovoEmissor(caminho string) *Emissor {
	return &Emissor{
		instrucoes: nil,
		arquivo:    nil,
		caminho:    caminho,
	}
}

func (e *Emissor) Emitir(instrucao string) {
	if e.bufferAtivo {
		e.buffer = append(e.buffer, instrucao)
		return
	}
	e.instrucoes = append(e.instrucoes, instrucao)
}

// StartBuffer faz as próximas Emitir irem para um buffer (para reordenar expressões).
func (e *Emissor) StartBuffer() {
	e.bufferAtivo = true
	e.buffer = nil
}

// GetBufferAndClear retorna o buffer e desativa o buffering.
func (e *Emissor) GetBufferAndClear() []string {
	e.bufferAtivo = false
	b := e.buffer
	e.buffer = nil
	return b
}

// EmitBuffer emite as instruções do buffer na ordem dada (termo esquerdo após termo direito).
func (e *Emissor) EmitBuffer(instrucoes []string) {
	for _, s := range instrucoes {
		e.instrucoes = append(e.instrucoes, s)
	}
}

// RemoverUltima remove e retorna a última instrução emitida
func (e *Emissor) RemoverUltima() string {
	if len(e.instrucoes) == 0 {
		return ""
	}
	ultima := e.instrucoes[len(e.instrucoes)-1]
	e.instrucoes = e.instrucoes[:len(e.instrucoes)-1]
	return ultima
}

func (e *Emissor) LinhaAtual() int {
	return len(e.instrucoes) - 1
}

func (e *Emissor) ProximaLinha() int {
	return len(e.instrucoes)
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
