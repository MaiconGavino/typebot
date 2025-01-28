package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func handlerDeletQuest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete{
		index := r.URL.Path[len("/quest/"):]
		var idx int
		_, err := fmt.Sscanf(index, "%d", &idx)
		if err != nil || idx < 0 || idx >= len(quests) {
			http.Error(w, "Index invalido ", http.StatusBadRequest)
			return
		}
		quests  = append(quests[:idx],quests[idx+1:]...)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Questão excluida com sucesso!")
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

func enableCORS(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Acess-Control-Allow-Origin", "*")
		w.Header().Set("Acess-Control-Allow-Methods", "POST, GET, DELETE, OPTIONS")
		w.Header().Set("Acess-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions{
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w,r)
	})
}

// Função principal para inicializar o servidor
func main() {
	fileserve := http.FileServer(http.Dir("./template"))
	http.Handle("/", http.StripPrefix("/", fileserve))

	http.Handle("/convert", enableCORS(http.HandlerFunc(handlerQuest)))

	//http.HandleFunc("/quest", handlerQuest)  // Rota para receber questões
	http.HandleFunc("/quests", handlerListQuest) // Rota para listar questões
	http.HandleFunc("/quest/", handlerDeletQuest)
	fmt.Println("Servidor iniciado na porta 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Erro ao iniciar o servidor: %v\n", err)
	}
}
