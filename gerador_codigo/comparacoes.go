package gerador_codigo

func (gerador *Gerador) GerarMenor() {
	gerador.emissor.Emitir("CPME")
}

func (gerador *Gerador) GerarMaior() {
	gerador.emissor.Emitir("CPMA")
}

func (gerador *Gerador) GerarIgual() {
	gerador.emissor.Emitir("CPIG")
}

func (gerador *Gerador) GerarDiferente() {
	gerador.emissor.Emitir("CDES")
}

func (gerador *Gerador) GerarMenorIgual() {
	gerador.emissor.Emitir("CPMI")
}

func (gerador *Gerador) GerarMaiorIgual() {
	gerador.emissor.Emitir("CMAI")
}
