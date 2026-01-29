package gerador_codigo

import "strconv"

func (gerador *Gerador) GerarAlocacao(quantidade int) {
	gerador.emissor.Emitir("ALME " + strconv.Itoa(quantidade))
}

func (gerador *Gerador) GerarDesalocacao(quantidade int) {
	gerador.emissor.Emitir("DESM " + strconv.Itoa(quantidade))
}
