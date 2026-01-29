package gerador_codigo

type Gerador struct {
	emissor *Emissor
	rotulo  int
}

func NovoGerador(saida string) *Gerador {
	return &Gerador{
		emissor: NovoEmissor(saida),
		rotulo:  0,
	}
}

func (gerador *Gerador) NovoRotulo() int {
	gerador.rotulo++
	return gerador.rotulo
}

func (gerador *Gerador) InicioPrograma() {
	gerador.emissor.Emitir("INPP")
}

func (gerador *Gerador) FimPrograma() {
	gerador.emissor.Emitir("PARA")
	gerador.emissor.Fechar()
}
