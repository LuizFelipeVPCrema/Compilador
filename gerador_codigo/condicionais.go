package gerador_codigo

import (
	"strconv"
)

// GerarIf: condição, DSVF (senão), bloco verdadeiro, DSVI (fim), [bloco falso].
func (gerador *Gerador) GerarIf(
	gerarCondicao func() error,
	blocoVerdadeiro func() error,
	blocoFalso func() error,
) error {
	if err := gerarCondicao(); err != nil {
		return err
	}
	indiceDSVF := gerador.EmitirComMarcador("DSVF")

	if err := blocoVerdadeiro(); err != nil {
		return err
	}
	indiceDSVI := gerador.EmitirComMarcador("DSVI")
	linhaSenao := gerador.ProximaLinha()

	if blocoFalso != nil {
		if err := blocoFalso(); err != nil {
			return err
		}
	}
	linhaFim := gerador.ProximaLinha()

	gerador.Preencher(indiceDSVI, linhaFim, "DSVI")
	gerador.Preencher(indiceDSVF, linhaSenao, "DSVF")
	return nil
}

// GerarWhile gera código para while: [início], condição, DSVF (fim), bloco, DSVI (início).
func (gerador *Gerador) GerarWhile(
	gerarCondicao func() error,
	bloco func() error,
) error {
	linhaInicio := gerador.ProximaLinha()

	if err := gerarCondicao(); err != nil {
		return err
	}
	indiceDSVF := gerador.EmitirComMarcador("DSVF")

	if err := bloco(); err != nil {
		return err
	}
	gerador.emissor.Emitir("DSVI " + strconv.Itoa(linhaInicio))

	linhaFim := gerador.ProximaLinha()
	gerador.Preencher(indiceDSVF, linhaFim, "DSVF")
	return nil
}
