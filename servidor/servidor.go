package servidor

import (
	"database-mysql/banco"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type usuario struct {
	ID    uint32 `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	bodyReq, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Falha ao ler o corpo da requisição"))
		return
	}

	var user usuario

	if err = json.Unmarshal(bodyReq, &user); err != nil {
		w.Write([]byte("Falha ao converter o usuário para struct"))
		return
	}

	fmt.Printf("Novo usuário a ser inserido: %+v\n", user)

	db, err := banco.Conectar()
	if err != nil {
		w.Write([]byte("Falha ao conectar ao banco de dados"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("INSERT INTO usuarios (nome, email) VALUES (?,?)")
	if err != nil {
		w.Write([]byte("Falha ao preparar o statement"))
		return
	}
	defer statement.Close()

	insercao, err := statement.Exec(user.Nome, user.Email)
	if err != nil {
		w.Write([]byte("Falha ao executar a inserção"))
		return
	}

	idInserido, err := insercao.LastInsertId()
	if err != nil {
		w.Write([]byte("Falha ao obter o ID do usuário"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Usuário criado com sucesso! ID: %d", idInserido)))
}

func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	db, err := banco.Conectar()
	if err != nil {
		w.Write([]byte("Falha ao conectar ao banco de dados"))
		return
	}
	defer db.Close()

	linhas, err := db.Query("SELECT * FROM usuarios")
	if err != nil {
		w.Write([]byte("Falha ao buscar os usuários"))
		return
	}
	defer linhas.Close()

	var usuarios []usuario

	for linhas.Next() {
		var usuario usuario

		if err := linhas.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); err != nil {
			w.Write([]byte("Falha ao escanear a linha"))
			return
		}

		usuarios = append(usuarios, usuario)
	}

	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(usuarios); err != nil {
		w.Write([]byte("Falha ao converter os usuarios para JSON"))
		return
	}
}

func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ID inválido"))
		return
	}

	db, err := banco.Conectar()
	if err != nil {
		w.Write([]byte("Falha ao conectar ao banco de dados"))
		return
	}
	defer db.Close()

	var usuario usuario
	linha, err := db.Query(fmt.Sprintf("SELECT * FROM usuarios WHERE id = %d", ID))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Falha ao buscar o usuário"))
		return
	}

	if linha.Next() {
		if err := linha.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Falha ao escanear a linha"))
			return
		}
	}
	if usuario.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Usuário não encontrado"))
		return
	}

	if err = json.NewEncoder(w).Encode(usuario); err != nil {
		w.Write([]byte("Falha ao converter o usuário para JSON"))
		return
	}
}

func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Erro ao converter o ID para uint32"))
		return
	}

	bodyReq, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Falha ao ler o corpo da requisição"))
		return
	}

	var usuario usuario
	if err := json.Unmarshal(bodyReq, &usuario); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Erro ao converter o usuario para struct"))
		return
	}

	db, err := banco.Conectar()
	if err != nil {
		w.Write([]byte("Falha ao conectar ao banco de dados"))
		return
	}

	statement, err := db.Prepare("UPDATE usuarios SET nome = ?, email = ? WHERE id = ?")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Falha ao preparar o statement"))
		return
	}

	_, err = statement.Exec(usuario.Nome, usuario.Email, ID)
	if err != nil {
		w.Write([]byte("Falha ao executar a atualização"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Erro ao converter o ID para uint32"))
		return
	}

	db, err := banco.Conectar()
	if err != nil {
		w.Write([]byte("Falha ao conectar ao banco de dados"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("DELETE FROM usuarios WHERE id = ?")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Falha ao preparar o statement"))
		return
	}
	defer statement.Close()

	_, err = statement.Exec(ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Falha ao tentar deletar o usuário"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
