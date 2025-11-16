package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

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
	Service string `msgpack:"service"`
	Data    struct {
		User      string `msgpack:"user"`
		Timestamp int64  `msgpack:"timestamp"`
		Clock     int64  `msgpack:"clock"`
	} `msgpack:"data"`
}

type LoginResponse struct {
	Service string `msgpack:"service"`
	Data    struct {
		Status      string `msgpack:"status"`
		Timestamp   int64  `msgpack:"timestamp"`
		Clock       int64  `msgpack:"clock"`
		Description string `msgpack:"description,omitempty"`
	} `msgpack:"data"`
}

type UsersRequest struct {
	Service string `msgpack:"service"`
	Data    struct {
		Timestamp int64 `msgpack:"timestamp"`
		Clock     int64 `msgpack:"clock"`
	} `msgpack:"data"`
}

type UsersResponse struct {
	Service string `msgpack:"service"`
	Data    struct {
		Timestamp int64    `msgpack:"timestamp"`
		Clock     int64    `msgpack:"clock"`
		Users     []string `msgpack:"users"`
	} `msgpack:"data"`
}

type ChannelRequest struct {
	Service string `msgpack:"service"`
	Data    struct {
		Channel   string `msgpack:"channel"`
		Timestamp int64  `msgpack:"timestamp"`
		Clock     int64  `msgpack:"clock"`
	} `msgpack:"data"`
}

type ChannelResponse struct {
	Service string `msgpack:"service"`
	Data    struct {
		Status      string `msgpack:"status"`
		Timestamp   int64  `msgpack:"timestamp"`
		Clock       int64  `msgpack:"clock"`
		Description string `msgpack:"description,omitempty"`
	} `msgpack:"data"`
}

type ChannelsRequest struct {
	Service string `msgpack:"service"`
	Data    struct {
		Timestamp int64 `msgpack:"timestamp"`
		Clock     int64 `msgpack:"clock"`
	} `msgpack:"data"`
}

type ChannelsResponse struct {
	Service string `msgpack:"service"`
	Data    struct {
		Timestamp int64    `msgpack:"timestamp"`
		Clock     int64    `msgpack:"clock"`
		Channels  []string `msgpack:"channels"`
	} `msgpack:"data"`
}

// Novas estruturas para Parte 2
type PublishRequest struct {
	Service string `msgpack:"service"`
	Data    struct {
		User      string `msgpack:"user"`
		Channel   string `msgpack:"channel"`
		Message   string `msgpack:"message"`
		Timestamp int64  `msgpack:"timestamp"`
		Clock     int64  `msgpack:"clock"`
	} `msgpack:"data"`
}

type PublishResponse struct {
	Service string `msgpack:"service"`
	Data    struct {
		Status    string `msgpack:"status"`
		Message   string `msgpack:"message,omitempty"`
		Timestamp int64  `msgpack:"timestamp"`
		Clock     int64  `msgpack:"clock"`
	} `msgpack:"data"`
}

type MessageRequest struct {
	Service string `msgpack:"service"`
	Data    struct {
		Src       string `msgpack:"src"`
		Dst       string `msgpack:"dst"`
		Message   string `msgpack:"message"`
		Timestamp int64  `msgpack:"timestamp"`
		Clock     int64  `msgpack:"clock"`
	} `msgpack:"data"`
}

type MessageResponse struct {
	Service string `msgpack:"service"`
	Data    struct {
		Status    string `msgpack:"status"`
		Message   string `msgpack:"message,omitempty"`
		Timestamp int64  `msgpack:"timestamp"`
		Clock     int64  `msgpack:"clock"`
	} `msgpack:"data"`
}

// Estrutura para publicação no broker
type Publication struct {
	User      string `msgpack:"user"`
	Message   string `msgpack:"message"`
	Timestamp int64  `msgpack:"timestamp"`
	Clock     int64  `msgpack:"clock"`
}

type DirectMessage struct {
	From      string `msgpack:"from"`
	Message   string `msgpack:"message"`
	Timestamp int64  `msgpack:"timestamp"`
	Clock     int64  `msgpack:"clock"`
}

// Estruturas de persistência
type UserLogin struct {
	Username  string `msgpack:"username"`
	Timestamp int64  `msgpack:"timestamp"`
}

type ChannelMessage struct {
	User      string `msgpack:"user"`
	Channel   string `msgpack:"channel"`
	Message   string `msgpack:"message"`
	Timestamp int64  `msgpack:"timestamp"`
}

type UserMessage struct {
	Src       string `msgpack:"src"`
	Dst       string `msgpack:"dst"`
	Message   string `msgpack:"message"`
	Timestamp int64  `msgpack:"timestamp"`
}

type PersistentData struct {
	Logins          []UserLogin      `msgpack:"logins"`
	Channels        []string         `msgpack:"channels"`
	ChannelMessages []ChannelMessage `msgpack:"channel_messages"`
	UserMessages    []UserMessage    `msgpack:"user_messages"`
}

// Estruturas para comunicação com o servidor de referência
type RankRequest struct {
	Service string `msgpack:"service"`
	Data    struct {
		User      string `msgpack:"user"`
		Timestamp int64  `msgpack:"timestamp"`
		Clock     int64  `msgpack:"clock"`
	} `msgpack:"data"`
}

type RankResponse struct {
	Service string `msgpack:"service"`
	Data    struct {
		Rank      int   `msgpack:"rank"`
		Timestamp int64 `msgpack:"timestamp"`
		Clock     int64 `msgpack:"clock"`
	} `msgpack:"data"`
}

type HeartbeatRequest struct {
	Service string `msgpack:"service"`
	Data    struct {
		User      string `msgpack:"user"`
		Timestamp int64  `msgpack:"timestamp"`
		Clock     int64  `msgpack:"clock"`
	} `msgpack:"data"`
}

// Estruturas para sincronização Berkeley
type ClockRequest struct {
	Service string `msgpack:"service"`
	Data    struct {
		Timestamp int64 `msgpack:"timestamp"`
		Clock     int64 `msgpack:"clock"`
	} `msgpack:"data"`
}

type ClockResponse struct {
	Service string `msgpack:"service"`
	Data    struct {
		Time      int64 `msgpack:"time"`
		Timestamp int64 `msgpack:"timestamp"`
		Clock     int64 `msgpack:"clock"`
	} `msgpack:"data"`
}

type ClockAdjustment struct {
	Service string `msgpack:"service"`
	Data    struct {
		Adjustment int64 `msgpack:"adjustment"`
		Timestamp  int64 `msgpack:"timestamp"`
		Clock      int64 `msgpack:"clock"`
	} `msgpack:"data"`
}

// Estrutura para listagem de servidores
type ListRequest struct {
	Service string `msgpack:"service"`
	Data    struct {
		Timestamp int64 `msgpack:"timestamp"`
		Clock     int64 `msgpack:"clock"`
	} `msgpack:"data"`
}

type ListResponse struct {
	Service string `msgpack:"service"`
	Data    struct {
		List      []ServerInfo `msgpack:"list"`
		Timestamp int64        `msgpack:"timestamp"`
		Clock     int64        `msgpack:"clock"`
	} `msgpack:"data"`
}

type ServerInfo struct {
	Name string `msgpack:"name"`
	Rank int    `msgpack:"rank"`
}

// Estruturas para eleição Bully
type ElectionRequest struct {
	Service string `msgpack:"service"`
	Data    struct {
		Timestamp int64 `msgpack:"timestamp"`
		Clock     int64 `msgpack:"clock"`
	} `msgpack:"data"`
}

type ElectionResponse struct {
	Service string `msgpack:"service"`
	Data    struct {
		Election  string `msgpack:"election"`
		Timestamp int64  `msgpack:"timestamp"`
		Clock     int64  `msgpack:"clock"`
	} `msgpack:"data"`
}

type CoordinatorAnnouncement struct {
	Service string `msgpack:"service"`
	Data    struct {
		Coordinator string `msgpack:"coordinator"`
		Timestamp   int64  `msgpack:"timestamp"`
		Clock       int64  `msgpack:"clock"`
	} `msgpack:"data"`
}

// Estruturas para replicação de dados (Parte 5)
type ReplicationRequest struct {
	Service string `msgpack:"service"`
	Data    struct {
		Type      string      `msgpack:"type"` // "login", "channel", "channel_message", "user_message"
		Content   interface{} `msgpack:"content"`
		Timestamp int64       `msgpack:"timestamp"`
		Clock     int64       `msgpack:"clock"`
	} `msgpack:"data"`
}

type ReplicationResponse struct {
	Service string `msgpack:"service"`
	Data    struct {
		Status    string `msgpack:"status"`
		Timestamp int64  `msgpack:"timestamp"`
		Clock     int64  `msgpack:"clock"`
	} `msgpack:"data"`
}

type SyncRequest struct {
	Service string `msgpack:"service"`
	Data    struct {
		LastSync  int64 `msgpack:"last_sync"`
		Timestamp int64 `msgpack:"timestamp"`
		Clock     int64 `msgpack:"clock"`
	} `msgpack:"data"`
}

type SyncResponse struct {
	Service string `msgpack:"service"`
	Data    struct {
		Logins          []UserLogin      `msgpack:"logins"`
		Channels        []string         `msgpack:"channels"`
		ChannelMessages []ChannelMessage `msgpack:"channel_messages"`
		UserMessages    []UserMessage    `msgpack:"user_messages"`
		Timestamp       int64            `msgpack:"timestamp"`
		Clock           int64            `msgpack:"clock"`
	} `msgpack:"data"`
}

const dataFile = "/data/server_data.json"

var data PersistentData
var pubSocket *zmq.Socket
var serverName string
var serverRank int
var coordinatorName string
var messageCounter int
var timeOffset int64     // Ajuste do relógio físico (Berkeley)
var lastSyncTime int64   // Última sincronização (Parte 5)
var dataMutex sync.Mutex // Proteger acesso aos dados

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

	if _, err := refSocket.SendBytes(reqData, 0); err != nil {
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

	if _, err := refSocket.SendBytes(reqData, 0); err != nil {
		return err
	}

	respData, err := refSocket.RecvBytes(0)
	if err != nil {
		return err
	}

	var resp struct {
		Service string `msgpack:"service"`
		Data    struct {
			Status    string `msgpack:"status"`
			Timestamp int64  `msgpack:"timestamp"`
			Clock     int64  `msgpack:"clock"`
		} `msgpack:"data"`
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

	if _, err := refSocket.SendBytes(reqData, 0); err != nil {
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
		Service string `msgpack:"service"`
		Data    struct {
			Status    string `msgpack:"status"`
			Timestamp int64  `msgpack:"timestamp"`
			Clock     int64  `msgpack:"clock"`
		} `msgpack:"data"`
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

	maxRank := -1
	var coordinator string

	for _, server := range servers {
		if server.Rank > maxRank ||
			(server.Rank == maxRank && server.Name > coordinator) {
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

// Substituir a função handleLogin no main.go

func handleLogin(msg []byte) ([]byte, error) {
	// Log do payload bruto para debug
	log.Printf("🔍 DEBUG: Recebido payload de login (tamanho: %d bytes)", len(msg))
	
	var req LoginRequest
	if err := msgpack.Unmarshal(msg, &req); err != nil {
		log.Printf("❌ Erro ao deserializar LoginRequest: %v", err)
		return nil, err
	}

	// Log detalhado dos campos recebidos
	log.Printf("🔍 DEBUG: Service=%s, User='%s', Timestamp=%d, Clock=%d", 
		req.Service, req.Data.User, req.Data.Timestamp, req.Data.Clock)

	// Atualizar relógio lógico ao receber mensagem
	updateClock(req.Data.Clock)

	resp := LoginResponse{Service: "login"}
	resp.Data.Timestamp = time.Now().Unix()
	resp.Data.Clock = incrementClock()

	// Validação melhorada
	trimmedUser := strings.TrimSpace(req.Data.User)
	
	if trimmedUser == "" {
		log.Printf("⚠️  Login rejeitado: usuário vazio (original: '%s', trimmed: '%s')", 
			req.Data.User, trimmedUser)
		resp.Data.Status = "erro"
		resp.Data.Description = "Nome de usuário não pode ser vazio"
	} else if userExists(trimmedUser) {
		log.Printf("⚠️  Login rejeitado: usuário '%s' já existe", trimmedUser)
		resp.Data.Status = "erro"
		resp.Data.Description = "Usuário já existe"
	} else {
		login := UserLogin{
			Username:  trimmedUser,
			Timestamp: req.Data.Timestamp,
		}

		dataMutex.Lock()
		data.Logins = append(data.Logins, login)
		dataMutex.Unlock()

		if err := saveData(); err != nil {
			log.Printf("❌ Erro ao salvar dados para usuário '%s': %v", trimmedUser, err)
			resp.Data.Status = "erro"
			resp.Data.Description = "Erro ao salvar dados: " + err.Error()
		} else {
			resp.Data.Status = "sucesso"
			log.Printf("✅ Novo usuário cadastrado: '%s' (clock: %d)", trimmedUser, resp.Data.Clock)

			// Replicar para outros servidores (assíncrono)
			go func() {
				refURL := os.Getenv("REFERENCE_URL")
				if refURL == "" {
					refURL = "tcp://reference:5559"
				}
				refSock, err := zmq.NewSocket(zmq.REQ)
				if err == nil {
					defer refSock.Close()
					refSock.Connect(refURL)
					time.Sleep(500 * time.Millisecond)
					replicateData(refSock, "login", login)
				}
			}()
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
		channel := req.Data.Channel

		dataMutex.Lock()
		data.Channels = append(data.Channels, channel)
		dataMutex.Unlock()

		if err := saveData(); err != nil {
			resp.Data.Status = "erro"
			resp.Data.Description = "Erro ao salvar dados: " + err.Error()
		} else {
			resp.Data.Status = "sucesso"
			log.Printf("✅ Novo canal criado: %s (clock: %d)", req.Data.Channel, resp.Data.Clock)

			// Replicar para outros servidores
			go func() {
				refURL := os.Getenv("REFERENCE_URL")
				if refURL == "" {
					refURL = "tcp://reference:5559"
				}
				refSock, err := zmq.NewSocket(zmq.REQ)
				if err == nil {
					defer refSock.Close()
					refSock.Connect(refURL)
					time.Sleep(500 * time.Millisecond)
					replicateData(refSock, "channel", channel)
				}
			}()
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
	if _, err := pubSocket.SendMessage(topic, pubData); err != nil {
		resp.Data.Status = "erro"
		resp.Data.Message = "Erro ao publicar mensagem: " + err.Error()
		log.Printf("❌ Erro ao publicar no canal %s: %v", topic, err)
		return msgpack.Marshal(resp)
	}

	// Salvar na persistência
	channelMsg := ChannelMessage{
		User:      req.Data.User,
		Channel:   req.Data.Channel,
		Message:   req.Data.Message,
		Timestamp: req.Data.Timestamp,
	}

	dataMutex.Lock()
	data.ChannelMessages = append(data.ChannelMessages, channelMsg)
	dataMutex.Unlock()

	if err := saveData(); err != nil {
		log.Printf("⚠️  Aviso: erro ao salvar mensagem: %v", err)
	}

	resp.Data.Status = "OK"
	log.Printf("📤 Publicação no canal #%s por %s (clock: %d)", req.Data.Channel, req.Data.User, pub.Clock)

	// Replicar mensagem para outros servidores
	go func() {
		refURL := os.Getenv("REFERENCE_URL")
		if refURL == "" {
			refURL = "tcp://reference:5559"
		}
		refSock, err := zmq.NewSocket(zmq.REQ)
		if err == nil {
			defer refSock.Close()
			refSock.Connect(refURL)
			time.Sleep(500 * time.Millisecond)
			replicateData(refSock, "channel_message", channelMsg)
		}
	}()

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
	if _, err := pubSocket.SendMessage(topic, dmData); err != nil {
		resp.Data.Status = "erro"
		resp.Data.Message = "Erro ao enviar mensagem: " + err.Error()
		log.Printf("❌ Erro ao enviar mensagem para %s: %v", topic, err)
		return msgpack.Marshal(resp)
	}

	// Salvar na persistência
	userMsg := UserMessage{
		Src:       req.Data.Src,
		Dst:       req.Data.Dst,
		Message:   req.Data.Message,
		Timestamp: req.Data.Timestamp,
	}

	dataMutex.Lock()
	data.UserMessages = append(data.UserMessages, userMsg)
	dataMutex.Unlock()

	if err := saveData(); err != nil {
		log.Printf("⚠️  Aviso: erro ao salvar mensagem: %v", err)
	}

	resp.Data.Status = "OK"
	log.Printf("💬 Mensagem de %s para %s (clock: %d)", req.Data.Src, req.Data.Dst, dm.Clock)

	// Replicar mensagem para outros servidores
	go func() {
		refURL := os.Getenv("REFERENCE_URL")
		if refURL == "" {
			refURL = "tcp://reference:5559"
		}
		refSock, err := zmq.NewSocket(zmq.REQ)
		if err == nil {
			defer refSock.Close()
			refSock.Connect(refURL)
			time.Sleep(500 * time.Millisecond)
			replicateData(refSock, "user_message", userMsg)
		}
	}()

	return msgpack.Marshal(resp)
}

// ----------------------------
// Replicação de dados
// ----------------------------

func replicateData(refSocket *zmq.Socket, dataType string, content interface{}) error {
	// Monta requisição de replicação
	req := ReplicationRequest{Service: "replicate"}
	req.Data.Type = dataType
	req.Data.Content = content
	req.Data.Timestamp = getAdjustedTime()
	req.Data.Clock = incrementClock()

	reqData, err := msgpack.Marshal(req)
	if err != nil {
		return fmt.Errorf("replicateData: erro ao serializar: %v", err)
	}

	// Envia e espera resposta do reference
	if _, err := refSocket.SendBytes(reqData, 0); err != nil {
		return fmt.Errorf("replicateData: erro ao enviar: %v", err)
	}

	respData, err := refSocket.RecvBytes(0)
	if err != nil {
		return fmt.Errorf("replicateData: erro ao receber: %v", err)
	}

	var resp ReplicationResponse
	if err := msgpack.Unmarshal(respData, &resp); err != nil {
		return fmt.Errorf("replicateData: erro ao desserializar resposta: %v", err)
	}

	// Atualiza relógio lógico com o clock retornado
	updateClock(resp.Data.Clock)
	if resp.Data.Status != "OK" {
		return fmt.Errorf("replicateData: reference retornou status != OK")
	}
	return nil
}

// Handler para requisições "replicate" recebidas por este servidor.
// Essa função aplica a réplica localmente (append nos slices) para manter persistência.
func handleReplication(msg []byte) ([]byte, error) {
	var req ReplicationRequest
	if err := msgpack.Unmarshal(msg, &req); err != nil {
		return nil, err
	}

	// Atualiza relógio lógico
	updateClock(req.Data.Clock)

	// Aplicar réplica conforme tipo
	dataMutex.Lock()
	switch req.Data.Type {
	case "login":
		// content -> UserLogin
		raw, _ := msgpack.Marshal(req.Data.Content)
		var ul UserLogin
		if err := msgpack.Unmarshal(raw, &ul); err == nil {
			// evitar duplicatas
			if !userExists(ul.Username) {
				data.Logins = append(data.Logins, ul)
			}
		}
	case "channel":
		raw, _ := msgpack.Marshal(req.Data.Content)
		var ch string
		if err := msgpack.Unmarshal(raw, &ch); err == nil {
			if !channelExists(ch) {
				data.Channels = append(data.Channels, ch)
			}
		}
	case "channel_message":
		raw, _ := msgpack.Marshal(req.Data.Content)
		var cm ChannelMessage
		if err := msgpack.Unmarshal(raw, &cm); err == nil {
			data.ChannelMessages = append(data.ChannelMessages, cm)
		}
	case "user_message":
		raw, _ := msgpack.Marshal(req.Data.Content)
		var um UserMessage
		if err := msgpack.Unmarshal(raw, &um); err == nil {
			data.UserMessages = append(data.UserMessages, um)
		}
	default:
		// tipo desconhecido: apenas log
		log.Printf("handleReplication: tipo desconhecido: %s", req.Data.Type)
	}
	_ = saveData() // tenta persistir (ignorar erro aqui)
	dataMutex.Unlock()

	// Responder OK
	resp := ReplicationResponse{Service: "replicate"}
	resp.Data.Status = "OK"
	resp.Data.Timestamp = getAdjustedTime()
	resp.Data.Clock = incrementClock()
	return msgpack.Marshal(resp)
}

// ----------------------------
// Eleição (Bully simples)
// ----------------------------

// handleElectionRequest — responde a pedidos de eleição.
// Se receber uma eleição, responde com "OK" e inicia sua própria eleição se tiver rank maior.
func handleElectionRequest(msg []byte) ([]byte, error) {
	var req ElectionRequest
	if err := msgpack.Unmarshal(msg, &req); err != nil {
		return nil, err
	}

	updateClock(req.Data.Clock)

	// Responder que recebeu a eleição
	resp := ElectionResponse{Service: "election"}
	resp.Data.Election = "OK"
	resp.Data.Timestamp = getAdjustedTime()
	resp.Data.Clock = incrementClock()

	// Se este servidor tem rank maior que o remetente, inicia sua própria eleição
	// (o payload do request não carrega o nome/ rank do remetente no modelo atual,
	//  então assumimos que o reference será responsável por encaminhar; para simplicidade,
	//  apenas retornamos OK aqui; a lógica de iniciar eleição localmente é feita pela função iniciadora)
	return msgpack.Marshal(resp)
}

// Inicia uma eleição Bully simples usando a lista de servidores do reference.
func initiateElection(refSocket *zmq.Socket) error {
	log.Printf("🏳️ Iniciando eleição Bully...")

	// Obter lista de servidores
	servers, err := getServerList(refSocket)
	if err != nil {
		return fmt.Errorf("initiateElection: erro ao obter lista: %v", err)
	}

	// Se ninguém tem rank maior, me torno coordenador
	higherExists := false
	for _, s := range servers {
		if s.Name == serverName {
			continue
		}
		if s.Rank > serverRank {
			higherExists = true
			// Tentar contatar servidor com rank maior
			serverURL := fmt.Sprintf("tcp://%s:5555", s.Name)
			sock, err := createServerSocket(serverURL)
			if err != nil {
				log.Printf("initiateElection: não consegui conectar %s: %v", s.Name, err)
				continue
			}

			// Envia pedido de eleição
			req := ElectionRequest{Service: "election"}
			req.Data.Timestamp = getAdjustedTime()
			req.Data.Clock = incrementClock()

			reqData, _ := msgpack.Marshal(req)
			if _, err := sock.SendBytes(reqData, 0); err != nil {
				log.Printf("initiateElection: erro ao enviar para %s: %v", s.Name, err)
				sock.Close()
				continue
			}

			// Aguarda resposta breve (non-blocking simple)
			sock.SetRcvtimeo(500 * time.Millisecond)
			respData, err := sock.RecvBytes(0)
			sock.Close()
			if err == nil && len(respData) > 0 {
				// recebeu resposta; então existe servidor superior respondendo
				log.Printf("initiateElection: %s respondeu, candidato superior presente", s.Name)
				// não me torno coordenador; aguardar anúncio
				return nil
			}
		}
	}

	if !higherExists {
		// Não há nenhum com rank maior: torno-me coordenador
		log.Printf("🏆 Nenhum servidor com rank maior encontrado — tornando-me coordenador")
		coordinatorName = serverName
		if err := becomeCoordinator(); err != nil {
			return fmt.Errorf("initiateElection: erro ao anunciar coordenadoria: %v", err)
		}
	}

	return nil
}

// becomeCoordinator — anuncia ao reference (ou à rede) que este servidor é o novo coordenador.
func becomeCoordinator() error {
	// Anuncia para o servidor de reference (REQ -> REFERENCE_URL)
	refURL := os.Getenv("REFERENCE_URL")
	if refURL == "" {
		refURL = "tcp://reference:5559"
	}

	refSock, err := zmq.NewSocket(zmq.REQ)
	if err != nil {
		return fmt.Errorf("becomeCoordinator: erro ao criar socket: %v", err)
	}
	defer refSock.Close()

	if err := refSock.Connect(refURL); err != nil {
		return fmt.Errorf("becomeCoordinator: erro ao conectar reference: %v", err)
	}

	ann := CoordinatorAnnouncement{Service: "coordinator"}
	ann.Data.Coordinator = serverName
	ann.Data.Timestamp = getAdjustedTime()
	ann.Data.Clock = incrementClock()

	data, err := msgpack.Marshal(ann)
	if err != nil {
		return fmt.Errorf("becomeCoordinator: erro ao serializar announcement: %v", err)
	}

	if _, err := refSock.SendBytes(data, 0); err != nil {
		return fmt.Errorf("becomeCoordinator: erro ao enviar announcement: %v", err)
	}

	// opcional: aguarda ACK
	_, _ = refSock.RecvBytes(0)
	log.Printf("👑 Anúncio de coordenador enviado: %s", serverName)
	return nil
}

// ----------------------------
// Checagem de health do coordenador
// ----------------------------

func checkCoordinatorHealth(refSocket *zmq.Socket) {
	// Se não houver coordenador conhecido, tenta determinar
	if coordinatorName == "" || coordinatorName == serverName {
		return
	}

	coordinatorURL := fmt.Sprintf("tcp://%s:5555", coordinatorName)
	sock, err := createServerSocket(coordinatorURL)
	if err != nil {
		log.Printf("checkCoordinatorHealth: não conseguiu conectar ao coordenador %s: %v", coordinatorName, err)
		// iniciar eleição
		go initiateElection(refSocket)
		return
	}
	defer sock.Close()

	// Envia heartbeat (usando o mesmo formato)
	req := HeartbeatRequest{Service: "heartbeat"}
	req.Data.User = serverName
	req.Data.Timestamp = getAdjustedTime()
	req.Data.Clock = incrementClock()

	reqData, _ := msgpack.Marshal(req)
	if _, err := sock.SendBytes(reqData, 0); err != nil {
		log.Printf("checkCoordinatorHealth: erro ao enviar heartbeat para %s: %v", coordinatorName, err)
		go initiateElection(refSocket)
		return
	}

	// Aguarda resposta com timeout curto
	sock.SetRcvtimeo(500 * time.Millisecond)
	respData, err := sock.RecvBytes(0)
	if err != nil || len(respData) == 0 {
		log.Printf("checkCoordinatorHealth: coordenador %s não respondeu, iniciando eleição", coordinatorName)
		go initiateElection(refSocket)
		return
	}

	// Se respondeu, apenas atualize relógio
	var resp struct {
		Service string `msgpack:"service"`
		Data    struct {
			Status    string `msgpack:"status"`
			Timestamp int64  `msgpack:"timestamp"`
			Clock     int64  `msgpack:"clock"`
		} `msgpack:"data"`
	}
	if err := msgpack.Unmarshal(respData, &resp); err == nil {
		updateClock(resp.Data.Clock)
		log.Printf("checkCoordinatorHealth: coordenador %s está ativo (clock: %d)", coordinatorName, resp.Data.Clock)
	}
}

// ----------------------------
// Sincronização periódica (Parte 5)
// ----------------------------

func startSyncRoutine(refSocket *zmq.Socket) {
	go func() {
		ticker := time.NewTicker(60 * time.Second) // intervalo de sync configurável
		defer ticker.Stop()
		for range ticker.C {
			// Se sou coordenador, faço sincronização Berkeley (coletar timestamps)
			if coordinatorName == serverName {
				if err := berkeleyCoordinator(refSocket); err != nil {
					log.Printf("startSyncRoutine: erro na sincronização Berkeley: %v", err)
				}
			}
		}
	}()
}

// ----------------------------
// Subscrição a anúncios de coordenador (PUB/SUB)
// ----------------------------

func subscribeToCoordinatorAnnouncements() {
	announceURL := os.Getenv("COORD_ANNOUNCE_URL")
	if announceURL == "" {
		// padrão: referência publica anúncios em 5560 (ajuste conforme sua infra)
		announceURL = "tcp://reference:5560"
	}

	subSock, err := zmq.NewSocket(zmq.SUB)
	if err != nil {
		log.Printf("subscribeToCoordinatorAnnouncements: erro ao criar SUB socket: %v", err)
		return
	}
	defer subSock.Close()

	if err := subSock.Connect(announceURL); err != nil {
		log.Printf("subscribeToCoordinatorAnnouncements: erro ao conectar %s: %v", announceURL, err)
		return
	}

	// receber tudo
	subSock.SetSubscribe("") // subscribe to all topics

	for {
		msg, err := subSock.RecvBytes(0)
		if err != nil {
			log.Printf("subscribeToCoordinatorAnnouncements: erro ao receber: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		var ann CoordinatorAnnouncement
		if err := msgpack.Unmarshal(msg, &ann); err != nil {
			log.Printf("subscribeToCoordinatorAnnouncements: não foi possível parsear anúncio: %v", err)
			continue
		}

		// Atualiza estado local
		coordinatorName = ann.Data.Coordinator
		updateClock(ann.Data.Clock)
		log.Printf("📣 Recebido anúncio de coordenador: %s (clock: %d)", coordinatorName, ann.Data.Clock)
	}
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

	// Iniciar rotina de sincronização periódica (Parte 5)
	startSyncRoutine(refSocket)

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
	log.Println("=" + strings.Repeat("=", 70))

	// Loop principal
	for {
		msg, err := repSocket.RecvBytes(0)
		if err != nil {
			log.Printf("❌ Erro ao receber mensagem: %v", err)
			continue
		}

		// Identificar o tipo de serviço
		var baseReq struct {
			Service string `msgpack:"service"`
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
		case "replicate":
			response, err = handleReplication(msg)
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
