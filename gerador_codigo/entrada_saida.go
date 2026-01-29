package gerador_codigo

func (gerador *Gerador) GerarLeitura() {
	gerador.emissor.Emitir("LEIT")
}
