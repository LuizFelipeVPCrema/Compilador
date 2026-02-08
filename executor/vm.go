package executor

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Guarda o programa, a "memória" (variáveis), as pilhas e de onde estamos.
type VM struct {
	instrucoes []string  // lista de instruções do código objeto (cada linha = uma instrução)
	mem        []float64 // memória global: índice = endereço da variável, valor = float (tudo vira real)
	pilha      []float64 // pilha de dados: expressões (ex: 2 + 3 empilha 2, 3, depois SOMA deixa 5)
	retorno    []int     // pilha de retorno: quando chama função (CHPR), guarda pra onde voltar (0-based)
	params     []float64 // parâmetros que o caller passou; no CHPR a gente copia isso pra mem[8], [9], ...
	pc         int       // program counter: índice da próxima instrução (0-based, igual o gerador)
	entrada    *bufio.Reader
	saida      io.Writer
}

// NovaVM cria uma VM. Se não passar entrada/saída, usa stdin/stdout (pra poder redirecionar depois).
func NovaVM(entrada io.Reader, saida io.Writer) *VM {
	if entrada == nil {
		entrada = os.Stdin
	}
	if saida == nil {
		saida = os.Stdout
	}
	return &VM{
		mem:     make([]float64, 0, 64), // capacidade inicial pra não ficar realocando toda hora
		pilha:   make([]float64, 0, 32),
		retorno: make([]int, 0, 16),
		params:  make([]float64, 0, 8),
		entrada: bufio.NewReader(entrada),
		saida:   saida,
	}
}

// CarregarPrograma abre o arquivo de código objeto e devolve uma linha por instrução.
func CarregarPrograma(caminho string) ([]string, error) {
	f, err := os.Open(caminho)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var inst []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		linha := strings.TrimSpace(sc.Text())
		if linha == "" {
			continue // para ignorar as linhas vazias
		}
		inst = append(inst, linha)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return inst, nil
}

// Executar carrega o programa do arquivo, zera PC e as pilhas (reuso de slice com [:0]),
// e chama rodar(). Assim dá pra rodar o mesmo programa de novo sem criar VM nova.
func (vm *VM) Executar(caminhoArquivo string) error {
	inst, err := CarregarPrograma(caminhoArquivo)
	if err != nil {
		return fmt.Errorf("carregar programa: %w", err)
	}
	vm.instrucoes = inst
	vm.pc = 0
	vm.pilha = vm.pilha[:0]
	vm.retorno = vm.retorno[:0]
	vm.params = vm.params[:0]
	return vm.rodar()
}

// rodar é o loop principal: pega a instrução no PC, separa instrução e argumento (ex: "ARMZ 5" -> op=ARMZ, arg=5),
// avança o PC (porque já "consumiu" essa instrução) e executa. PARA faz pc = len(instrucoes) e sai do loop.
func (vm *VM) rodar() error {
	for vm.pc >= 0 && vm.pc < len(vm.instrucoes) {
		linha := vm.instrucoes[vm.pc]
		partes := strings.SplitN(linha, " ", 2) // no máximo 2 partes: "OP" e "resto"
		op := partes[0]
		arg := ""
		if len(partes) > 1 {
			arg = strings.TrimSpace(partes[1])
		}

		vm.pc++ // avança antes de executar; assim desvios (DSVI, etc.) sobrescrevem o pc certo

		if err := vm.executarUma(op, arg); err != nil {
			return err
		}
	}
	return nil
}

// garantirMemoria: o compilador pode usar endereços altos (ex: 12, 15). Em vez de alocar um vetor gigante,
// vai crescendo o slice conforme precisa (append de 0 até o addr existir).
func (vm *VM) garantirMemoria(addr int) {
	for len(vm.mem) <= addr {
		vm.mem = append(vm.mem, 0)
	}
}

// executarUma interpreta uma instrução e execta:
func (vm *VM) executarUma(op, arg string) error {
	switch op {
	// --- Início / memória (o gerador emite ALME/DESM mas a gente não precisa fazer nada;
	//     os endereços já são fixos na tabela de símbolos, só garantimos memória no ARMZ/CRVL) ---
	case "INPP":
		return nil
	case "ALME":
		return nil
	case "DESM":
		return nil
	// constante, carregar variável, armazenar
	case "CRCT":
		val, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			return fmt.Errorf("linha %d: CRCT valor inválido %q: %w", vm.pc, arg, err)
		}
		vm.pilha = append(vm.pilha, val)
		return nil
	case "CRVL":
		addr, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("linha %d: CRVL endereço inválido %q: %w", vm.pc, arg, err)
		}
		vm.garantirMemoria(addr)
		vm.pilha = append(vm.pilha, vm.mem[addr])
		return nil
	case "ARMZ":
		addr, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("linha %d: ARMZ endereço inválido %q: %w", vm.pc, arg, err)
		}
		if len(vm.pilha) == 0 {
			return fmt.Errorf("linha %d: ARMZ com pilha vazia", vm.pc)
		}
		val := vm.pilha[len(vm.pilha)-1]
		vm.pilha = vm.pilha[:len(vm.pilha)-1]
		vm.garantirMemoria(addr)
		vm.mem[addr] = val
		return nil
	// desempilha 2, aplica op, empilha resultado
	case "SOMA":
		return vm.binop(func(a, b float64) float64 { return a + b })
	case "SUBT":
		return vm.binop(func(a, b float64) float64 { return a - b })
	case "MULT":
		return vm.binop(func(a, b float64) float64 { return a * b })
	case "DIVI":
		return vm.binop(func(a, b float64) float64 {
			if b == 0 {
				return 0 // evita panic
			}
			return a / b
		})
	case "CPME": // menor
		return vm.cmp(func(a, b float64) bool { return a < b })
	case "CPMA": // maior
		return vm.cmp(func(a, b float64) bool { return a > b })
	case "CPIG": // igual
		return vm.cmp(func(a, b float64) bool { return a == b })
	case "CDES": // diferente
		return vm.cmp(func(a, b float64) bool { return a != b })
	case "CPMI": // menor ou igual
		return vm.cmp(func(a, b float64) bool { return a <= b })
	case "CMAI": // maior ou igual
		return vm.cmp(func(a, b float64) bool { return a >= b })
	//  Desvios
	case "DSVI":
		alvo, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("linha %d: DSVI alvo inválido %q: %w", vm.pc, arg, err)
		}

		vm.pc = alvo
		return nil
	case "DSVF":
		alvo, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("linha %d: DSVF alvo inválido %q: %w", vm.pc, arg, err)
		}
		if len(vm.pilha) == 0 {
			return fmt.Errorf("linha %d: DSVF com pilha vazia", vm.pc)
		}
		topo := vm.pilha[len(vm.pilha)-1]
		vm.pilha = vm.pilha[:len(vm.pilha)-1]
		if topo == 0 {
			vm.pc = alvo
		}
		return nil
	// LEIT lê uma linha e empilha; IMPR desempilha e imprime
	case "LEIT":
		linha, err := vm.entrada.ReadString('\n')
		if err != nil && err != io.EOF {
			return fmt.Errorf("linha %d: LEIT: %w", vm.pc, err)
		}
		linha = strings.TrimSpace(linha)
		if linha == "" {
			if err == io.EOF {
				return fmt.Errorf("linha %d: LEIT: entrada insuficiente (EOF)", vm.pc)
			}
			vm.pilha = append(vm.pilha, 0) // linha vazia = 0 (evita erro)
			return nil
		}
		val, err := strconv.ParseFloat(linha, 64)
		if err != nil {
			return fmt.Errorf("linha %d: LEIT: valor inválido %q: %w", vm.pc, linha, err)
		}
		vm.pilha = append(vm.pilha, val)
		return nil
	case "IMPR":
		if len(vm.pilha) == 0 {
			return fmt.Errorf("linha %d: IMPR com pilha vazia", vm.pc)
		}
		val := vm.pilha[len(vm.pilha)-1]
		vm.pilha = vm.pilha[:len(vm.pilha)-1]
		fmt.Fprintln(vm.saida, formatNum(val))
		return nil
	// PUSHER guarda endereço de retorno;

	case "PUSHER":
		addr, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("linha %d: PUSHER endereço inválido %q: %w", vm.pc, arg, err)
		}
		vm.retorno = append(vm.retorno, addr) // caller já colocou o endereço certo (0-based) no código
		return nil
	// PARAM copia valor da memória pra lista;
	case "PARAM":
		addr, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("linha %d: PARAM endereço inválido %q: %w", vm.pc, arg, err)
		}
		vm.garantirMemoria(addr)
		vm.params = append(vm.params, vm.mem[addr]) // vamos passar esses valores pra função
		return nil
	// CHPR copia params pra mem[8].. e pula pro início da função;
	case "CHPR":
		alvo, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("linha %d: CHPR alvo inválido %q: %w", vm.pc, arg, err)
		}
		// Convenção do compilador: parâmetros da função ficam em mem[8], mem[9], ...
		base := 8
		for i, v := range vm.params {
			vm.garantirMemoria(base + i)
			vm.mem[base+i] = v
		}
		vm.params = vm.params[:0]
		vm.pc = alvo // pula pro início da função; o retorno já foi empilhado pelo PUSHER
		return nil
	// RTPR volta pro endereço guardado
	case "RTPR":
		if len(vm.retorno) == 0 {
			return fmt.Errorf("linha %d: RTPR com pilha de retorno vazia", vm.pc)
		}
		addr := vm.retorno[len(vm.retorno)-1]
		vm.retorno = vm.retorno[:len(vm.retorno)-1]
		vm.pc = addr // volta pra instrução logo depois do CHPR
		return nil
	case "PARA":
		vm.pc = len(vm.instrucoes) // sai do loop do rodar()
		return nil
	default:
		return fmt.Errorf("linha %d: instrução desconhecida %q", vm.pc, op)
	}
}

// binop: desempilha os dois topos (b depois a), aplica a função (soma, subtração, etc.) e empilha o resultado.
func (vm *VM) binop(f func(a, b float64) float64) error {
	if len(vm.pilha) < 2 {
		return fmt.Errorf("linha %d: operação binária com pilha insuficiente", vm.pc)
	}
	b := vm.pilha[len(vm.pilha)-1]
	a := vm.pilha[len(vm.pilha)-2]
	vm.pilha = vm.pilha[:len(vm.pilha)-2]
	vm.pilha = append(vm.pilha, f(a, b))
	return nil
}

// cmp: compara os dois topos; empilha 1 se verdadeiro, 0 se falso (pra DSVF usar depois).
func (vm *VM) cmp(f func(a, b float64) bool) error {
	if len(vm.pilha) < 2 {
		return fmt.Errorf("linha %d: comparação com pilha insuficiente", vm.pc)
	}
	b := vm.pilha[len(vm.pilha)-1]
	a := vm.pilha[len(vm.pilha)-2]
	vm.pilha = vm.pilha[:len(vm.pilha)-2]
	if f(a, b) {
		vm.pilha = append(vm.pilha, 1)
	} else {
		vm.pilha = append(vm.pilha, 0)
	}
	return nil
}

// formatNum: se o valor for "inteiro" (ex: 5.0), imprime sem casa decimal; senão imprime como real.
func formatNum(v float64) string {
	if v == float64(int64(v)) {
		return strconv.FormatInt(int64(v), 10)
	}
	return strconv.FormatFloat(v, 'f', -1, 64)
}
