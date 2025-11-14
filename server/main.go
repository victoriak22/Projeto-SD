package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	zmq "github.com/pebbe/zmq4"
)

// Estruturas de dados
type LoginRequest struct {
	Service string `json:"service"`
	Data    struct {
		User      string `json:"user"`
		Timestamp int64  `json:"timestamp"`
	} `json:"data"`
}

type LoginResponse struct {
	Service string `json:"service"`
	Data    struct {
		Status      string `json:"status"`
		Timestamp   int64  `json:"timestamp"`
		Description string `json:"description,omitempty"`
	} `json:"data"`
}

type UsersRequest struct {
	Service string `json:"service"`
	Data    struct {
		Timestamp int64 `json:"timestamp"`
	} `json:"data"`
}

type UsersResponse struct {
	Service string `json:"service"`
	Data    struct {
		Timestamp int64    `json:"timestamp"`
		Users     []string `json:"users"`
	} `json:"data"`
}

type ChannelRequest struct {
	Service string `json:"service"`
	Data    struct {
		Channel   string `json:"channel"`
		Timestamp int64  `json:"timestamp"`
	} `json:"data"`
}

type ChannelResponse struct {
	Service string `json:"service"`
	Data    struct {
		Status      string `json:"status"`
		Timestamp   int64  `json:"timestamp"`
		Description string `json:"description,omitempty"`
	} `json:"data"`
}

type ChannelsRequest struct {
	Service string `json:"service"`
	Data    struct {
		Timestamp int64 `json:"timestamp"`
	} `json:"data"`
}

type ChannelsResponse struct {
	Service string `json:"service"`
	Data    struct {
		Timestamp int64    `json:"timestamp"`
		Channels  []string `json:"channels"`
	} `json:"data"`
}

// Novas estruturas para Parte 2
type PublishRequest struct {
	Service string `json:"service"`
	Data    struct {
		User      string `json:"user"`
		Channel   string `json:"channel"`
		Message   string `json:"message"`
		Timestamp int64  `json:"timestamp"`
	} `json:"data"`
}

type PublishResponse struct {
	Service string `json:"service"`
	Data    struct {
		Status    string `json:"status"`
		Message   string `json:"message,omitempty"`
		Timestamp int64  `json:"timestamp"`
	} `json:"data"`
}

type MessageRequest struct {
	Service string `json:"service"`
	Data    struct {
		Src       string `json:"src"`
		Dst       string `json:"dst"`
		Message   string `json:"message"`
		Timestamp int64  `json:"timestamp"`
	} `json:"data"`
}

type MessageResponse struct {
	Service string `json:"service"`
	Data    struct {
		Status    string `json:"status"`
		Message   string `json:"message,omitempty"`
		Timestamp int64  `json:"timestamp"`
	} `json:"data"`
}

// Estrutura para publicação no broker
type Publication struct {
	User      string `json:"user"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type DirectMessage struct {
	From      string `json:"from"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

// Estruturas de persistência
type UserLogin struct {
	Username  string `json:"username"`
	Timestamp int64  `json:"timestamp"`
}

type ChannelMessage struct {
	User      string `json:"user"`
	Channel   string `json:"channel"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type UserMessage struct {
	Src       string `json:"src"`
	Dst       string `json:"dst"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type PersistentData struct {
	Logins          []UserLogin      `json:"logins"`
	Channels        []string         `json:"channels"`
	ChannelMessages []ChannelMessage `json:"channel_messages"`
	UserMessages    []UserMessage    `json:"user_messages"`
}

const dataFile = "/data/server_data.json"

var data PersistentData
var pubSocket *zmq.Socket

// Funções de persistência
func loadData() error {
	file, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			data = PersistentData{
				Logins:          []UserLogin{},
				Channels:        []string{},
				ChannelMessages: []ChannelMessage{},
				UserMessages:    []UserMessage{},
			}
			return saveData()
		}
		return err
	}
	return json.Unmarshal(file, &data)
}

func saveData() error {
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	
	os.MkdirAll("/data", 0755)
	
	return os.WriteFile(dataFile, file, 0644)
}

func userExists(username string) bool {
	for _, login := range data.Logins {
		if login.Username == username {
			return true
		}
	}
	return false
}

func channelExists(channel string) bool {
	for _, ch := range data.Channels {
		if ch == channel {
			return true
		}
	}
	return false
}

func getUniqueUsers() []string {
	userMap := make(map[string]bool)
	for _, login := range data.Logins {
		userMap[login.Username] = true
	}
	
	users := []string{}
	for user := range userMap {
		users = append(users, user)
	}
	return users
}

// Handlers da Parte 1
func handleLogin(msg []byte) ([]byte, error) {
	var req LoginRequest
	if err := json.Unmarshal(msg, &req); err != nil {
		return nil, err
	}

	resp := LoginResponse{Service: "login"}
	resp.Data.Timestamp = time.Now().Unix()

	if req.Data.User == "" {
		resp.Data.Status = "erro"
		resp.Data.Description = "Nome de usuário não pode ser vazio"
	} else if userExists(req.Data.User) {
		resp.Data.Status = "erro"
		resp.Data.Description = "Usuário já existe"
	} else {
		data.Logins = append(data.Logins, UserLogin{
			Username:  req.Data.User,
			Timestamp: req.Data.Timestamp,
		})
		
		if err := saveData(); err != nil {
			resp.Data.Status = "erro"
			resp.Data.Description = "Erro ao salvar dados: " + err.Error()
		} else {
			resp.Data.Status = "sucesso"
			log.Printf("✅ Novo usuário cadastrado: %s", req.Data.User)
		}
	}

	return json.Marshal(resp)
}

func handleUsers(msg []byte) ([]byte, error) {
	var req UsersRequest
	if err := json.Unmarshal(msg, &req); err != nil {
		return nil, err
	}

	resp := UsersResponse{Service: "users"}
	resp.Data.Timestamp = time.Now().Unix()
	resp.Data.Users = getUniqueUsers()

	return json.Marshal(resp)
}

func handleChannel(msg []byte) ([]byte, error) {
	var req ChannelRequest
	if err := json.Unmarshal(msg, &req); err != nil {
		return nil, err
	}

	resp := ChannelResponse{Service: "channel"}
	resp.Data.Timestamp = time.Now().Unix()

	if req.Data.Channel == "" {
		resp.Data.Status = "erro"
		resp.Data.Description = "Nome do canal não pode ser vazio"
	} else if channelExists(req.Data.Channel) {
		resp.Data.Status = "erro"
		resp.Data.Description = "Canal já existe"
	} else {
		data.Channels = append(data.Channels, req.Data.Channel)
		
		if err := saveData(); err != nil {
			resp.Data.Status = "erro"
			resp.Data.Description = "Erro ao salvar dados: " + err.Error()
		} else {
			resp.Data.Status = "sucesso"
			log.Printf("✅ Novo canal criado: %s", req.Data.Channel)
		}
	}

	return json.Marshal(resp)
}

func handleChannels(msg []byte) ([]byte, error) {
	var req ChannelsRequest
	if err := json.Unmarshal(msg, &req); err != nil {
		return nil, err
	}

	resp := ChannelsResponse{Service: "channels"}
	resp.Data.Timestamp = time.Now().Unix()
	resp.Data.Channels = data.Channels

	return json.Marshal(resp)
}

// Novos handlers da Parte 2
func handlePublish(msg []byte) ([]byte, error) {
	var req PublishRequest
	if err := json.Unmarshal(msg, &req); err != nil {
		return nil, err
	}

	resp := PublishResponse{Service: "publish"}
	resp.Data.Timestamp = time.Now().Unix()

	// Validações
	if !channelExists(req.Data.Channel) {
		resp.Data.Status = "erro"
		resp.Data.Message = "Canal não existe"
		return json.Marshal(resp)
	}

	if req.Data.Message == "" {
		resp.Data.Status = "erro"
		resp.Data.Message = "Mensagem não pode ser vazia"
		return json.Marshal(resp)
	}

	// Criar publicação
	pub := Publication{
		User:      req.Data.User,
		Message:   req.Data.Message,
		Timestamp: req.Data.Timestamp,
	}

	pubData, err := json.Marshal(pub)
	if err != nil {
		resp.Data.Status = "erro"
		resp.Data.Message = "Erro ao serializar mensagem"
		return json.Marshal(resp)
	}

	// Publicar no broker (tópico = nome do canal)
	topic := req.Data.Channel
	if err := pubSocket.SendMessage(topic, pubData); err != nil {
		resp.Data.Status = "erro"
		resp.Data.Message = "Erro ao publicar mensagem: " + err.Error()
		log.Printf("❌ Erro ao publicar no canal %s: %v", topic, err)
		return json.Marshal(resp)
	}

	// Salvar na persistência
	data.ChannelMessages = append(data.ChannelMessages, ChannelMessage{
		User:      req.Data.User,
		Channel:   req.Data.Channel,
		Message:   req.Data.Message,
		Timestamp: req.Data.Timestamp,
	})

	if err := saveData(); err != nil {
		log.Printf("⚠️  Aviso: erro ao salvar mensagem: %v", err)
	}

	resp.Data.Status = "OK"
	log.Printf("📤 Publicação no canal #%s por %s: %s", req.Data.Channel, req.Data.User, req.Data.Message)

	return json.Marshal(resp)
}

func handleMessage(msg []byte) ([]byte, error) {
	var req MessageRequest
	if err := json.Unmarshal(msg, &req); err != nil {
		return nil, err
	}

	resp := MessageResponse{Service: "message"}
	resp.Data.Timestamp = time.Now().Unix()

	// Validações
	if !userExists(req.Data.Dst) {
		resp.Data.Status = "erro"
		resp.Data.Message = "Usuário de destino não existe"
		return json.Marshal(resp)
	}

	if req.Data.Message == "" {
		resp.Data.Status = "erro"
		resp.Data.Message = "Mensagem não pode ser vazia"
		return json.Marshal(resp)
	}

	// Criar mensagem direta
	dm := DirectMessage{
		From:      req.Data.Src,
		Message:   req.Data.Message,
		Timestamp: req.Data.Timestamp,
	}

	dmData, err := json.Marshal(dm)
	if err != nil {
		resp.Data.Status = "erro"
		resp.Data.Message = "Erro ao serializar mensagem"
		return json.Marshal(resp)
	}

	// Publicar no broker (tópico = nome do usuário de destino)
	topic := req.Data.Dst
	if err := pubSocket.SendMessage(topic, dmData); err != nil {
		resp.Data.Status = "erro"
		resp.Data.Message = "Erro ao enviar mensagem: " + err.Error()
		log.Printf("❌ Erro ao enviar mensagem para %s: %v", topic, err)
		return json.Marshal(resp)
	}

	// Salvar na persistência
	data.UserMessages = append(data.UserMessages, UserMessage{
		Src:       req.Data.Src,
		Dst:       req.Data.Dst,
		Message:   req.Data.Message,
		Timestamp: req.Data.Timestamp,
	})

	if err := saveData(); err != nil {
		log.Printf("⚠️  Aviso: erro ao salvar mensagem: %v", err)
	}

	resp.Data.Status = "OK"
	log.Printf("💬 Mensagem de %s para %s: %s", req.Data.Src, req.Data.Dst, req.Data.Message)

	return json.Marshal(resp)
}

func main() {
	log.Println("🚀 Iniciando servidor...")

	// Carregar dados persistentes
	if err := loadData(); err != nil {
		log.Fatalf("❌ Erro ao carregar dados: %v", err)
	}
	log.Printf("📊 Dados carregados: %d logins, %d canais, %d msgs canal, %d msgs usuário", 
		len(data.Logins), len(data.Channels), len(data.ChannelMessages), len(data.UserMessages))

	// Configurar socket REQ-REP
	repSocket, err := zmq.NewSocket(zmq.REP)
	if err != nil {
		log.Fatalf("❌ Erro ao criar socket REP: %v", err)
	}
	defer repSocket.Close()

	err = repSocket.Bind("tcp://*:5555")
	if err != nil {
		log.Fatalf("❌ Erro ao fazer bind REP: %v", err)
	}
	log.Println("📡 Socket REP escutando na porta 5555...")

	// Configurar socket PUB (conecta ao broker XSUB)
	pubSocket, err = zmq.NewSocket(zmq.PUB)
	if err != nil {
		log.Fatalf("❌ Erro ao criar socket PUB: %v", err)
	}
	defer pubSocket.Close()

	brokerURL := "tcp://broker:5557"
	err = pubSocket.Connect(brokerURL)
	if err != nil {
		log.Fatalf("❌ Erro ao conectar ao broker: %v", err)
	}
	log.Printf("🔌 Socket PUB conectado ao broker em %s", brokerURL)

	log.Println("✅ Servidor pronto para receber requisições!")
	log.Println("=" + "=" * 70)

	// Loop principal
	for {
		msg, err := repSocket.RecvBytes(0)
		if err != nil {
			log.Printf("❌ Erro ao receber mensagem: %v", err)
			continue
		}

		// Identificar o tipo de serviço
		var baseReq struct {
			Service string `json:"service"`
		}
		if err := json.Unmarshal(msg, &baseReq); err != nil {
			log.Printf("❌ Erro ao parsear mensagem: %v", err)
			repSocket.SendBytes([]byte(`{"error":"Formato de mensagem inválido"}`), 0)
			continue
		}

		var response []byte
		switch baseReq.Service {
		case "login":
			response, err = handleLogin(msg)
		case "users":
			response, err = handleUsers(msg)
		case "channel":
			response, err = handleChannel(msg)
		case "channels":
			response, err = handleChannels(msg)
		case "publish":
			response, err = handlePublish(msg)
		case "message":
			response, err = handleMessage(msg)
		default:
			response = []byte(fmt.Sprintf(`{"error":"Serviço desconhecido: %s"}`, baseReq.Service))
		}

		if err != nil {
			log.Printf("❌ Erro ao processar requisição: %v", err)
			response = []byte(fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		}

		repSocket.SendBytes(response, 0)
	}
}