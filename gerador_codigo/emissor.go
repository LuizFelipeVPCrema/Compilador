package gerador_codigo

import (
	"fmt"
	"os"
)

type Emissor struct {
	arquivo *os.File
}

func NovoEmissor(caminho string) *Emissor {
	arquivo, _ := os.Create(caminho)
	return &Emissor{arquivo: arquivo}
}

func (e *Emissor) Emitir(instrucao string) {
	fmt.Fprintln(e.arquivo, instrucao)
}

func (e *Emissor) Fechar() {
	e.arquivo.Close()
}
