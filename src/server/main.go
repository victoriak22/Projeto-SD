package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
	"strings"

	zmq "github.com/pebbe/zmq4"
	"github.com/vmihailenco/msgpack/v5"
)

// Relógio lógico global
var (
	logicalClock int64
	clockMutex   sync.Mutex
)

// Funções do relógio lógico
func incrementClock() int64 {
	clockMutex.Lock()
	defer clockMutex.Unlock()
	logicalClock++
	return logicalClock
}

func updateClock(receivedClock int64) int64 {
	clockMutex.Lock()
	defer clockMutex.Unlock()
	if receivedClock > logicalClock {
		logicalClock = receivedClock
	}
	logicalClock++
	return logicalClock
}

func getClock() int64 {
	clockMutex.Lock()
	defer clockMutex.Unlock()
	return logicalClock
}

// Funções para relógio físico ajustado
func getAdjustedTime() int64 {
	return time.Now().Unix() + timeOffset
}

func adjustTime(adjustment int64) {
	timeOffset += adjustment
	log.Printf("⏰ Relógio ajustado em %ds (offset total: %ds)", adjustment, timeOffset)
}

// Estruturas de dados
type LoginRequest struct {
	Service string `json:"service"`
	Data    struct {
		User      string `json:"user"`
		Timestamp int64  `json:"timestamp"`
		Clock     int64  `json:"clock"`
	} `json:"data"`
}

type LoginResponse struct {
	Service string `json:"service"`
	Data    struct {
		Status      string `json:"status"`
		Timestamp   int64  `json:"timestamp"`
		Clock       int64  `json:"clock"`
		Description string `json:"description,omitempty"`
	} `json:"data"`
}

type UsersRequest struct {
	Service string `json:"service"`
	Data    struct {
		Timestamp int64 `json:"timestamp"`
		Clock     int64 `json:"clock"`
	} `json:"data"`
}

type UsersResponse struct {
	Service string `json:"service"`
	Data    struct {
		Timestamp int64    `json:"timestamp"`
		Clock     int64    `json:"clock"`
		Users     []string `json:"users"`
	} `json:"data"`
}

type ChannelRequest struct {
	Service string `json:"service"`
	Data    struct {
		Channel   string `json:"channel"`
		Timestamp int64  `json:"timestamp"`
		Clock     int64  `json:"clock"`
	} `json:"data"`
}

type ChannelResponse struct {
	Service string `json:"service"`
	Data    struct {
		Status      string `json:"status"`
		Timestamp   int64  `json:"timestamp"`
		Clock       int64  `json:"clock"`
		Description string `json:"description,omitempty"`
	} `json:"data"`
}

type ChannelsRequest struct {
	Service string `json:"service"`
	Data    struct {
		Timestamp int64 `json:"timestamp"`
		Clock     int64 `json:"clock"`
	} `json:"data"`
}

type ChannelsResponse struct {
	Service string `json:"service"`
	Data    struct {
		Timestamp int64    `json:"timestamp"`
		Clock     int64    `json:"clock"`
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
		Clock     int64  `json:"clock"`
	} `json:"data"`
}

type PublishResponse struct {
	Service string `json:"service"`
	Data    struct {
		Status    string `json:"status"`
		Message   string `json:"message,omitempty"`
		Timestamp int64  `json:"timestamp"`
		Clock     int64  `json:"clock"`
	} `json:"data"`
}

type MessageRequest struct {
	Service string `json:"service"`
	Data    struct {
		Src       string `json:"src"`
		Dst       string `json:"dst"`
		Message   string `json:"message"`
		Timestamp int64  `json:"timestamp"`
		Clock     int64  `json:"clock"`
	} `json:"data"`
}

type MessageResponse struct {
	Service string `json:"service"`
	Data    struct {
		Status    string `json:"status"`
		Message   string `json:"message,omitempty"`
		Timestamp int64  `json:"timestamp"`
		Clock     int64  `json:"clock"`
	} `json:"data"`
}

// Estrutura para publicação no broker
type Publication struct {
	User      string `json:"user"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
	Clock     int64  `json:"clock"`
}

type DirectMessage struct {
	From      string `json:"from"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
	Clock     int64  `json:"clock"`
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

// Estruturas para comunicação com o servidor de referência
type RankRequest struct {
	Service string `json:"service"`
	Data    struct {
		User      string `json:"user"`
		Timestamp int64  `json:"timestamp"`
		Clock     int64  `json:"clock"`
	} `json:"data"`
}

type RankResponse struct {
	Service string `json:"service"`
	Data    struct {
		Rank      int   `json:"rank"`
		Timestamp int64 `json:"timestamp"`
		Clock     int64 `json:"clock"`
	} `json:"data"`
}

type HeartbeatRequest struct {
	Service string `json:"service"`
	Data    struct {
		User      string `json:"user"`
		Timestamp int64  `json:"timestamp"`
		Clock     int64  `json:"clock"`
	} `json:"data"`
}

// Estruturas para sincronização Berkeley
type ClockRequest struct {
	Service string `json:"service"`
	Data    struct {
		Timestamp int64 `json:"timestamp"`
		Clock     int64 `json:"clock"`
	} `json:"data"`
}

type ClockResponse struct {
	Service string `json:"service"`
	Data    struct {
		Time      int64 `json:"time"`
		Timestamp int64 `json:"timestamp"`
		Clock     int64 `json:"clock"`
	} `json:"data"`
}

type ClockAdjustment struct {
	Service string `json:"service"`
	Data    struct {
		Adjustment int64 `json:"adjustment"`
		Timestamp  int64 `json:"timestamp"`
		Clock      int64 `json:"clock"`
	} `json:"data"`
}

// Estrutura para listagem de servidores
type ListRequest struct {
	Service string `json:"service"`
	Data    struct {
		Timestamp int64 `json:"timestamp"`
		Clock     int64 `json:"clock"`
	} `json:"data"`
}

type ListResponse struct {
	Service string `json:"service"`
	Data    struct {
		List      []ServerInfo `json:"list"`
		Timestamp int64        `json:"timestamp"`
		Clock     int64        `json:"clock"`
	} `json:"data"`
}

type ServerInfo struct {
	Name string `json:"name"`
	Rank int    `json:"rank"`
}

// Estruturas para eleição Bully
type ElectionRequest struct {
	Service string `json:"service"`
	Data    struct {
		Timestamp int64 `json:"timestamp"`
		Clock     int64 `json:"clock"`
	} `json:"data"`
}

type ElectionResponse struct {
	Service string `json:"service"`
	Data    struct {
		Election  string `json:"election"`
		Timestamp int64  `json:"timestamp"`
		Clock     int64  `json:"clock"`
	} `json:"data"`
}

type CoordinatorAnnouncement struct {
	Service string `json:"service"`
	Data    struct {
		Coordinator string `json:"coordinator"`
		Timestamp   int64  `json:"timestamp"`
		Clock       int64  `json:"clock"`
	} `json:"data"`
}

const dataFile = "/data/server_data.json"

var data PersistentData
var pubSocket *zmq.Socket
var serverName string
var serverRank int
var coordinatorName string
var messageCounter int
var timeOffset int64 // Ajuste do relógio físico (Berkeley)

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

// Funções para comunicação com o servidor de referência
func registerWithReference(refSocket *zmq.Socket) error {
	req := RankRequest{Service: "rank"}
	req.Data.User = serverName
	req.Data.Timestamp = time.Now().Unix()
	req.Data.Clock = incrementClock()

	reqData, err := msgpack.Marshal(req)
	if err != nil {
		return fmt.Errorf("erro ao serializar requisição: %v", err)
	}

	if err := refSocket.SendBytes(reqData, 0); err != nil {
		return fmt.Errorf("erro ao enviar requisição: %v", err)
	}

	respData, err := refSocket.RecvBytes(0)
	if err != nil {
		return fmt.Errorf("erro ao receber resposta: %v", err)
	}

	var resp RankResponse
	if err := msgpack.Unmarshal(respData, &resp); err != nil {
		return fmt.Errorf("erro ao deserializar resposta: %v", err)
	}

	updateClock(resp.Data.Clock)
	serverRank = resp.Data.Rank
	log.Printf("✅ Servidor registrado com rank: %d", serverRank)

	return nil
}

func sendHeartbeat(refSocket *zmq.Socket) error {
	req := HeartbeatRequest{Service: "heartbeat"}
	req.Data.User = serverName
	req.Data.Timestamp = time.Now().Unix()
	req.Data.Clock = incrementClock()

	reqData, err := msgpack.Marshal(req)
	if err != nil {
		return err
	}

	if err := refSocket.SendBytes(reqData, 0); err != nil {
		return err
	}

	respData, err := refSocket.RecvBytes(0)
	if err != nil {
		return err
	}

	var resp struct {
		Service string `json:"service"`
		Data    struct {
			Status    string `json:"status"`
			Timestamp int64  `json:"timestamp"`
			Clock     int64  `json:"clock"`
		} `json:"data"`
	}

	if err := msgpack.Unmarshal(respData, &resp); err != nil {
		return err
	}

	updateClock(resp.Data.Clock)
	log.Printf("💓 Heartbeat enviado (rank: %d, clock: %d)", serverRank, resp.Data.Clock)

	return nil
}

// Funções para obter lista de servidores do reference
func getServerList(refSocket *zmq.Socket) ([]ServerInfo, error) {
	req := ListRequest{Service: "list"}
	req.Data.Timestamp = getAdjustedTime()
	req.Data.Clock = incrementClock()

	reqData, err := msgpack.Marshal(req)
	if err != nil {
		return nil, err
	}

	if err := refSocket.SendBytes(reqData, 0); err != nil {
		return nil, err
	}

	respData, err := refSocket.RecvBytes(0)
	if err != nil {
		return nil, err
	}

	var resp ListResponse
	if err := msgpack.Unmarshal(respData, &resp); err != nil {
		return nil, err
	}

	updateClock(resp.Data.Clock)
	return resp.Data.List, nil
}

// Função para criar socket REQ temporário para comunicação entre servidores
func createServerSocket(serverURL string) (*zmq.Socket, error) {
	socket, err := zmq.NewSocket(zmq.REQ)
	if err != nil {
		return nil, err
	}

	if err := socket.Connect(serverURL); err != nil {
		socket.Close()
		return nil, err
	}

	return socket, nil
}

// Sincronização Berkeley - Coordenador coleta timestamps
func berkeleyCoordinator(refSocket *zmq.Socket) error {
	log.Printf("🎯 Iniciando sincronização Berkeley como COORDENADOR")

	// Obter lista de servidores
	servers, err := getServerList(refSocket)
	if err != nil {
		return fmt.Errorf("erro ao obter lista de servidores: %v", err)
	}

	if len(servers) <= 1 {
		log.Printf("⚠️  Apenas 1 servidor ativo, sincronização não necessária")
		return nil
	}

	// Coletar timestamps de todos os servidores
	timestamps := make(map[string]int64)
	timestamps[serverName] = getAdjustedTime()

	log.Printf("📊 Coletando timestamps de %d servidores...", len(servers))

	for _, server := range servers {
		if server.Name == serverName {
			continue // Skip self
		}

		serverURL := fmt.Sprintf("tcp://%s:5555", server.Name)
		socket, err := createServerSocket(serverURL)
		if err != nil {
			log.Printf("⚠️  Erro ao conectar a %s: %v", server.Name, err)
			continue
		}

		// Enviar requisição de clock
		req := ClockRequest{Service: "clock"}
		req.Data.Timestamp = getAdjustedTime()
		req.Data.Clock = incrementClock()

		reqData, _ := msgpack.Marshal(req)
		socket.SendBytes(reqData, 0)

		// Receber resposta
		respData, err := socket.RecvBytes(0)
		socket.Close()

		if err != nil {
			log.Printf("⚠️  Erro ao receber de %s: %v", server.Name, err)
			continue
		}

		var resp ClockResponse
		if err := msgpack.Unmarshal(respData, &resp); err != nil {
			log.Printf("⚠️  Erro ao parsear resposta de %s: %v", server.Name, err)
			continue
		}

		updateClock(resp.Data.Clock)
		timestamps[server.Name] = resp.Data.Time
		log.Printf("   📥 %s: %d", server.Name, resp.Data.Time)
	}

	// Calcular tempo médio
	var sum int64
	for _, t := range timestamps {
		sum += t
	}
	avgTime := sum / int64(len(timestamps))
	log.Printf("📊 Tempo médio calculado: %d", avgTime)

	// Distribuir ajustes
	for _, server := range servers {
		if server.Name == serverName {
			// Ajustar próprio relógio
			adjustment := avgTime - timestamps[serverName]
			if adjustment != 0 {
				adjustTime(adjustment)
			}
			continue
		}

		serverURL := fmt.Sprintf("tcp://%s:5555", server.Name)
		socket, err := createServerSocket(serverURL)
		if err != nil {
			log.Printf("⚠️  Erro ao conectar a %s para ajuste: %v", server.Name, err)
			continue
		}

		// Calcular e enviar ajuste
		adjustment := avgTime - timestamps[server.Name]
		
		adj := ClockAdjustment{Service: "adjust"}
		adj.Data.Adjustment = adjustment
		adj.Data.Timestamp = getAdjustedTime()
		adj.Data.Clock = incrementClock()

		adjData, _ := msgpack.Marshal(adj)
		socket.SendBytes(adjData, 0)
		
		// Aguardar confirmação
		socket.RecvBytes(0)
		socket.Close()

		log.Printf("   📤 Enviado ajuste de %ds para %s", adjustment, server.Name)
	}

	log.Printf("✅ Sincronização Berkeley concluída")
	return nil
}

// Handler para requisição de clock (coordenador pedindo meu tempo)
func handleClockRequest(msg []byte) ([]byte, error) {
	var req ClockRequest
	if err := msgpack.Unmarshal(msg, &req); err != nil {
		return nil, err
	}

	updateClock(req.Data.Clock)

	resp := ClockResponse{Service: "clock"}
	resp.Data.Time = getAdjustedTime()
	resp.Data.Timestamp = getAdjustedTime()
	resp.Data.Clock = incrementClock()

	return msgpack.Marshal(resp)
}

// Handler para ajuste de relógio (coordenador mandando ajuste)
func handleClockAdjustment(msg []byte) ([]byte, error) {
	var req ClockAdjustment
	if err := msgpack.Unmarshal(msg, &req); err != nil {
		return nil, err
	}

	updateClock(req.Data.Clock)

	// Aplicar ajuste
	adjustTime(req.Data.Adjustment)

	// Responder OK
	resp := struct {
		Service string `json:"service"`
		Data    struct {
			Status    string `json:"status"`
			Timestamp int64  `json:"timestamp"`
			Clock     int64  `json:"clock"`
		} `json:"data"`
	}{Service: "adjust"}
	
	resp.Data.Status = "OK"
	resp.Data.Timestamp = getAdjustedTime()
	resp.Data.Clock = incrementClock()

	return msgpack.Marshal(resp)
}

// Determinar coordenador (servidor com maior rank)
func determineCoordinator(refSocket *zmq.Socket) (string, error) {
	servers, err := getServerList(refSocket)
	if err != nil {
		return "", err
	}

	if len(servers) == 0 {
		return serverName, nil
	}

	// Encontrar servidor com maior rank
	maxRank := 0
	coordinator := serverName

	for _, server := range servers {
		if server.Rank > maxRank {
			maxRank = server.Rank
			coordinator = server.Name
		}
	}

	return coordinator, nil
}

// Verificar se deve sincronizar (a cada 10 mensagens)
func checkAndSyncIfNeeded(refSocket *zmq.Socket) {
	messageCounter++
	
	if messageCounter >= 10 {
		messageCounter = 0
		
		// Determinar coordenador
		coordinator, err := determineCoordinator(refSocket)
		if err != nil {
			log.Printf("⚠️  Erro ao determinar coordenador: %v", err)
			return
		}

		coordinatorName = coordinator

		// Se sou coordenador, sincronizar
		if coordinator == serverName {
			go func() {
				if err := berkeleyCoordinator(refSocket); err != nil {
					log.Printf("⚠️  Erro na sincronização: %v", err)
				}
			}()
		}
	}

func startHeartbeatRoutine(refSocket *zmq.Socket) {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		heartbeatCount := 0

		for range ticker.C {
			if err := sendHeartbeat(refSocket); err != nil {
				log.Printf("⚠️  Erro ao enviar heartbeat: %v", err)
			}
			
			heartbeatCount++
			
			// A cada 3 heartbeats (30s), verificar coordenador
			if heartbeatCount >= 3 {
				heartbeatCount = 0
				checkCoordinatorHealth(refSocket)
			}
		}
	}()
}

// Handlers da Parte 1
func handleLogin(msg []byte) ([]byte, error) {
	var req LoginRequest
	if err := msgpack.Unmarshal(msg, &req); err != nil {
		return nil, err
	}

	// Atualizar relógio lógico ao receber mensagem
	updateClock(req.Data.Clock)

	resp := LoginResponse{Service: "login"}
	resp.Data.Timestamp = time.Now().Unix()
	resp.Data.Clock = incrementClock() // Incrementar antes de enviar

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
			log.Printf("✅ Novo usuário cadastrado: %s (clock: %d)", req.Data.User, resp.Data.Clock)
		}
	}

	return msgpack.Marshal(resp)
}

func handleUsers(msg []byte) ([]byte, error) {
	var req UsersRequest
	if err := msgpack.Unmarshal(msg, &req); err != nil {
		return nil, err
	}

	updateClock(req.Data.Clock)

	resp := UsersResponse{Service: "users"}
	resp.Data.Timestamp = time.Now().Unix()
	resp.Data.Clock = incrementClock()
	resp.Data.Users = getUniqueUsers()

	return msgpack.Marshal(resp)
}

func handleChannel(msg []byte) ([]byte, error) {
	var req ChannelRequest
	if err := msgpack.Unmarshal(msg, &req); err != nil {
		return nil, err
	}

	updateClock(req.Data.Clock)

	resp := ChannelResponse{Service: "channel"}
	resp.Data.Timestamp = time.Now().Unix()
	resp.Data.Clock = incrementClock()

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
			log.Printf("✅ Novo canal criado: %s (clock: %d)", req.Data.Channel, resp.Data.Clock)
		}
	}

	return msgpack.Marshal(resp)
}

func handleChannels(msg []byte) ([]byte, error) {
	var req ChannelsRequest
	if err := msgpack.Unmarshal(msg, &req); err != nil {
		return nil, err
	}

	updateClock(req.Data.Clock)

	resp := ChannelsResponse{Service: "channels"}
	resp.Data.Timestamp = time.Now().Unix()
	resp.Data.Clock = incrementClock()
	resp.Data.Channels = data.Channels

	return msgpack.Marshal(resp)
}

// Novos handlers da Parte 2
func handlePublish(msg []byte) ([]byte, error) {
	var req PublishRequest
	if err := msgpack.Unmarshal(msg, &req); err != nil {
		return nil, err
	}

	updateClock(req.Data.Clock)

	resp := PublishResponse{Service: "publish"}
	resp.Data.Timestamp = time.Now().Unix()
	resp.Data.Clock = incrementClock()

	// Validações
	if !channelExists(req.Data.Channel) {
		resp.Data.Status = "erro"
		resp.Data.Message = "Canal não existe"
		return msgpack.Marshal(resp)
	}

	if req.Data.Message == "" {
		resp.Data.Status = "erro"
		resp.Data.Message = "Mensagem não pode ser vazia"
		return msgpack.Marshal(resp)
	}

	// Criar publicação com relógio lógico
	pub := Publication{
		User:      req.Data.User,
		Message:   req.Data.Message,
		Timestamp: req.Data.Timestamp,
		Clock:     incrementClock(),
	}

	pubData, err := msgpack.Marshal(pub)
	if err != nil {
		resp.Data.Status = "erro"
		resp.Data.Message = "Erro ao serializar mensagem"
		return msgpack.Marshal(resp)
	}

	// Publicar no broker (tópico = nome do canal)
	topic := req.Data.Channel
	if err := pubSocket.SendMessage(topic, pubData); err != nil {
		resp.Data.Status = "erro"
		resp.Data.Message = "Erro ao publicar mensagem: " + err.Error()
		log.Printf("❌ Erro ao publicar no canal %s: %v", topic, err)
		return msgpack.Marshal(resp)
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
	log.Printf("📤 Publicação no canal #%s por %s (clock: %d)", req.Data.Channel, req.Data.User, pub.Clock)

	return msgpack.Marshal(resp)
}

func handleMessage(msg []byte) ([]byte, error) {
	var req MessageRequest
	if err := msgpack.Unmarshal(msg, &req); err != nil {
		return nil, err
	}

	updateClock(req.Data.Clock)

	resp := MessageResponse{Service: "message"}
	resp.Data.Timestamp = time.Now().Unix()
	resp.Data.Clock = incrementClock()

	// Validações
	if !userExists(req.Data.Dst) {
		resp.Data.Status = "erro"
		resp.Data.Message = "Usuário de destino não existe"
		return msgpack.Marshal(resp)
	}

	if req.Data.Message == "" {
		resp.Data.Status = "erro"
		resp.Data.Message = "Mensagem não pode ser vazia"
		return msgpack.Marshal(resp)
	}

	// Criar mensagem direta com relógio lógico
	dm := DirectMessage{
		From:      req.Data.Src,
		Message:   req.Data.Message,
		Timestamp: req.Data.Timestamp,
		Clock:     incrementClock(),
	}

	dmData, err := msgpack.Marshal(dm)
	if err != nil {
		resp.Data.Status = "erro"
		resp.Data.Message = "Erro ao serializar mensagem"
		return msgpack.Marshal(resp)
	}

	// Publicar no broker (tópico = nome do usuário de destino)
	topic := req.Data.Dst
	if err := pubSocket.SendMessage(topic, dmData); err != nil {
		resp.Data.Status = "erro"
		resp.Data.Message = "Erro ao enviar mensagem: " + err.Error()
		log.Printf("❌ Erro ao enviar mensagem para %s: %v", topic, err)
		return msgpack.Marshal(resp)
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
	log.Printf("💬 Mensagem de %s para %s (clock: %d)", req.Data.Src, req.Data.Dst, dm.Clock)

	return msgpack.Marshal(resp)
}

func main() {
	log.Println("🚀 Iniciando servidor...")

	// Obter nome do servidor da variável de ambiente
	serverName = os.Getenv("SERVER_NAME")
	if serverName == "" {
		serverName = "server-default"
	}
	log.Printf("📛 Nome do servidor: %s", serverName)

	// Carregar dados persistentes
	if err := loadData(); err != nil {
		log.Fatalf("❌ Erro ao carregar dados: %v", err)
	}
	log.Printf("📊 Dados carregados: %d logins, %d canais, %d msgs canal, %d msgs usuário", 
		len(data.Logins), len(data.Channels), len(data.ChannelMessages), len(data.UserMessages))

	// Conectar ao servidor de referência
	refURL := os.Getenv("REFERENCE_URL")
	if refURL == "" {
		refURL = "tcp://reference:5559"
	}
	
	refSocket, err := zmq.NewSocket(zmq.REQ)
	if err != nil {
		log.Fatalf("❌ Erro ao criar socket de referência: %v", err)
	}
	defer refSocket.Close()

	err = refSocket.Connect(refURL)
	if err != nil {
		log.Fatalf("❌ Erro ao conectar ao servidor de referência: %v", err)
	}
	log.Printf("🔌 Conectado ao servidor de referência: %s", refURL)

	// Aguardar um pouco para garantir conexão
	time.Sleep(2 * time.Second)

	// Registrar no servidor de referência e obter rank
	if err := registerWithReference(refSocket); err != nil {
		log.Fatalf("❌ Erro ao registrar no servidor de referência: %v", err)
	}

	// Iniciar rotina de heartbeat
	startHeartbeatRoutine(refSocket)

	// Iniciar goroutine para receber anúncios de coordenador
	go subscribeToCoordinatorAnnouncements()

	// Aguardar um pouco para subscrição
	time.Sleep(2 * time.Second)

	// Determinar coordenador inicial
	initialCoordinator, err := determineCoordinator(refSocket)
	if err != nil {
		log.Printf("⚠️  Erro ao determinar coordenador inicial: %v", err)
	} else {
		coordinatorName = initialCoordinator
		log.Printf("👑 Coordenador inicial: %s", coordinatorName)
		
		// Se sou o coordenador, anunciar
		if coordinatorName == serverName {
			becomeCoordinator()
		}
	}

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

	log.Printf("✅ Servidor '%s' (rank %d) pronto para receber requisições!", serverName, serverRank)
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
		if err := msgpack.Unmarshal(msg, &baseReq); err != nil {
			log.Printf("❌ Erro ao parsear mensagem: %v", err)
			errorResp, _ := msgpack.Marshal(map[string]string{"error": "Formato de mensagem inválido"})
			repSocket.SendBytes(errorResp, 0)
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
		case "clock":
			response, err = handleClockRequest(msg)
		case "adjust":
			response, err = handleClockAdjustment(msg)
		case "election":
			response, err = handleElectionRequest(msg)
		default:
			response, _ = msgpack.Marshal(map[string]string{"error": fmt.Sprintf("Serviço desconhecido: %s", baseReq.Service)})
		}

		if err != nil {
			log.Printf("❌ Erro ao processar requisição: %v", err)
			response, _ = msgpack.Marshal(map[string]string{"error": err.Error()})
		}

		repSocket.SendBytes(response, 0)

		// Verificar e sincronizar se necessário (a cada 10 mensagens)
		checkAndSyncIfNeeded(refSocket)
	}
}