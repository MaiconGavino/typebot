package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Estrutura para armazenar uma questão recebida
type QuestResponse struct {
	Name  string `json:"name"`
	Quest string `json:"quest"`
}

// Armazena todas as questões enviadas
var quests []QuestResponse

// Handler para receber uma questão (rota: /quest)
func handlerQuest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var qt QuestResponse
		// Decodifica o corpo da requisição JSON em uma estrutura QuestResponse
		err := json.NewDecoder(r.Body).Decode(&qt)
		if err != nil {
			http.Error(w, "Dados inválidos", http.StatusBadRequest)
			return
		}
		// Adiciona a questão à lista
		quests = append(quests, qt)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Questão recebida com sucesso!")
		return
	}
	http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
}

// Handler para excluir uma questão (rota: /quest/{index})
func handlerDeleteQuest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		// Obtém o índice da URL
		indexStr := strings.TrimPrefix(r.URL.Path, "/quest/")
		idx, err := strconv.Atoi(indexStr)
		if err != nil || idx < 0 || idx >= len(quests) {
			http.Error(w, "Índice inválido", http.StatusBadRequest)
			return
		}
		// Remove a questão pelo índice
		quests = append(quests[:idx], quests[idx+1:]...)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Questão excluída com sucesso!")
		return
	}
	http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
}

// Handler para listar todas as questões (rota: /quests)
func handlerListQuest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Codifica a lista de questões em JSON e envia como resposta
	err := json.NewEncoder(w).Encode(quests)
	if err != nil {
		http.Error(w, "Erro ao processar os dados", http.StatusInternalServerError)
		return
	}
}

// Middleware para configurar CORS
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Função principal para inicializar o servidor
func main() {
	// Configuração para servir arquivos estáticos
	fileserve := http.FileServer(http.Dir("./template"))
	http.Handle("/", http.StripPrefix("/", fileserve))

	// Rotas da API com middleware de CORS
	http.Handle("/quest", enableCORS(http.HandlerFunc(handlerQuest)))         // Adicionar questões
	http.Handle("/quests", enableCORS(http.HandlerFunc(handlerListQuest)))    // Listar questões
	http.Handle("/quest/", enableCORS(http.HandlerFunc(handlerDeleteQuest))) // Excluir questões

	fmt.Println("Servidor iniciado na porta 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Erro ao iniciar o servidor: %v\n", err)
	}
}
