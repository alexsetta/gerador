package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	linhas := 100
	if len(os.Args) > 1 {
		n, err := strconv.Atoi(os.Args[1])
		if err == nil && n > 0 {
			linhas = n
		}
	}

	if err := geraArquivoTeste(linhas); err != nil {
		log.Fatalf("erro: %v", err)
	}
}

const dirGeracao = "."

func geraArquivoTeste(linhas int) error {
	log.Printf("Gerando arquivo de teste com %d linhas", linhas)
	if err := os.MkdirAll(dirGeracao, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório: %w", err)
	}

	nome := fmt.Sprintf("%s/teste_%d.txt", dirGeracao, linhas)
	f, err := os.Create(nome)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	base := time.Now()

	labels := make([]string, linhas)
	for i := range labels {
		labels[i] = fmt.Sprintf("AB%09dBR", rand.Intn(999_999_999)+1)
	}

	if linhas > 100 {
		numDup := linhas / 1000
		if numDup < 1 {
			numDup = 1
		}
		if numDup > linhas-1 {
			numDup = linhas - 1
		}
		for _, t := range rand.Perm(linhas - 1)[:numDup] {
			targetIdx := t + 1
			labels[targetIdx] = labels[rand.Intn(targetIdx)]
		}
	}

	for i := 0; i < linhas; i++ {
		timestamp := base.Add(time.Duration(i+1) * time.Millisecond).Format("2006-01-02 15:04:05.0000")
		etiqueta := labels[i]

		status := 0
		if rand.Intn(100) >= 95 {
			status = rand.Intn(9) + 1
		}

		h := sha256.Sum256([]byte(timestamp + etiqueta + fmt.Sprintf("%d", status)))
		hash := fmt.Sprintf("%x", h)

		fmt.Fprintf(w, "%s;%s;%d;%s\n", timestamp, etiqueta, status, hash)
	}

	err = w.Flush()
	if err != nil {
		return fmt.Errorf("erro ao flushar o arquivo: %w", err)
	}

	log.Printf("Arquivo %s gerado com sucesso", nome)
	return nil
}
