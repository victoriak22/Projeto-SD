const zeromq = require('zeromq');
const readline = require('readline');

// Configuração dos sockets ZeroMQ
const reqSocket = new zeromq.Request();
const subSocket = new zeromq.Subscriber();

// Interface para input do usuário
const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout
});

let currentUser = null;
let connected = false;
let subscribedChannels = new Set();

// Função para enviar requisição e receber resposta
async function sendRequest(request) {
  try {
    await reqSocket.send(JSON.stringify(request));
    const [response] = await reqSocket.receive();
    return JSON.parse(response.toString());
  } catch (error) {
    console.error('Erro na comunicação:', error.message);
    return null;
  }
}

// Função para receber mensagens do broker (em background)
async function receiveMessages() {
  for await (const [topic, msg] of subSocket) {
    try {
      const topicStr = topic.toString();
      const data = JSON.parse(msg.toString());
      const timestamp = new Date(data.timestamp * 1000).toLocaleString();
      
      // Mensagem de canal
      if (subscribedChannels.has(topicStr)) {
        console.log(`\n📺 [#${topicStr}] ${data.user}: ${data.message}`);
        console.log(`   ⏰ ${timestamp}`);
      } 
      // Mensagem direta
      else if (topicStr === currentUser) {
        console.log(`\n💬 [DM de ${data.from}]: ${data.message}`);
        console.log(`   ⏰ ${timestamp}`);
      }
      
      // Reexibir prompt
      process.stdout.write('\n➡️  Escolha uma opção: ');
    } catch (error) {
      console.error('Erro ao processar mensagem:', error.message);
    }
  }
}

// Função de login
async function login(username) {
  const request = {
    service: 'login',
    data: {
      user: username,
      timestamp: Math.floor(Date.now() / 1000)
    }
  };

  console.log('\n📤 Enviando requisição de login...');
  const response = await sendRequest(request);

  if (response) {
    if (response.data.status === 'sucesso') {
      currentUser = username;
      
      // Inscrever-se para receber mensagens diretas
      subSocket.subscribe(username);
      console.log(`✅ Login realizado com sucesso! Bem-vindo, ${username}!`);
      console.log(`📬 Inscrito para receber mensagens diretas`);
      console.log(`⏰ Timestamp: ${new Date(response.data.timestamp * 1000).toLocaleString()}`);
      return true;
    } else {
      console.log(`❌ Erro no login: ${response.data.description}`);
      return false;
    }
  }
  return false;
}

// Função para listar usuários
async function listUsers() {
  const request = {
    service: 'users',
    data: {
      timestamp: Math.floor(Date.now() / 1000)
    }
  };

  console.log('\n📤 Buscando lista de usuários...');
  const response = await sendRequest(request);

  if (response && response.data.users) {
    console.log('\n👥 Usuários cadastrados:');
    if (response.data.users.length === 0) {
      console.log('   (Nenhum usuário cadastrado ainda)');
    } else {
      response.data.users.forEach((user, index) => {
        const marker = user === currentUser ? '(você)' : '';
        console.log(`   ${index + 1}. ${user} ${marker}`);
      });
    }
    console.log(`⏰ Timestamp: ${new Date(response.data.timestamp * 1000).toLocaleString()}`);
  }
}

// Função para criar canal
async function createChannel(channelName) {
  const request = {
    service: 'channel',
    data: {
      channel: channelName,
      timestamp: Math.floor(Date.now() / 1000)
    }
  };

  console.log('\n📤 Criando canal...');
  const response = await sendRequest(request);

  if (response) {
    if (response.data.status === 'sucesso') {
      console.log(`✅ Canal "${channelName}" criado com sucesso!`);
      console.log(`⏰ Timestamp: ${new Date(response.data.timestamp * 1000).toLocaleString()}`);
      return true;
    } else {
      console.log(`❌ Erro ao criar canal: ${response.data.description}`);
      return false;
    }
  }
  return false;
}

// Função para listar canais
async function listChannels() {
  const request = {
    service: 'channels',
    data: {
      timestamp: Math.floor(Date.now() / 1000)
    }
  };

  console.log('\n📤 Buscando lista de canais...');
  const response = await sendRequest(request);

  if (response && response.data.channels) {
    console.log('\n📺 Canais disponíveis:');
    if (response.data.channels.length === 0) {
      console.log('   (Nenhum canal criado ainda)');
    } else {
      response.data.channels.forEach((channel, index) => {
        const subscribed = subscribedChannels.has(channel) ? '✓ inscrito' : '';
        console.log(`   ${index + 1}. #${channel} ${subscribed}`);
      });
    }
    console.log(`⏰ Timestamp: ${new Date(response.data.timestamp * 1000).toLocaleString()}`);
  }
}

// Função para inscrever em canal
async function subscribeChannel(channelName) {
  // Verificar se canal existe
  const channelsReq = {
    service: 'channels',
    data: { timestamp: Math.floor(Date.now() / 1000) }
  };
  
  const response = await sendRequest(channelsReq);
  
  if (response && response.data.channels.includes(channelName)) {
    subSocket.subscribe(channelName);
    subscribedChannels.add(channelName);
    console.log(`✅ Inscrito no canal #${channelName}`);
    return true;
  } else {
    console.log(`❌ Canal #${channelName} não existe`);
    return false;
  }
}

// Função para publicar em canal
async function publishMessage(channelName, message) {
  const request = {
    service: 'publish',
    data: {
      user: currentUser,
      channel: channelName,
      message: message,
      timestamp: Math.floor(Date.now() / 1000)
    }
  };

  console.log('\n📤 Publicando mensagem...');
  const response = await sendRequest(request);

  if (response) {
    if (response.data.status === 'OK') {
      console.log(`✅ Mensagem publicada no canal #${channelName}`);
      console.log(`⏰ Timestamp: ${new Date(response.data.timestamp * 1000).toLocaleString()}`);
      return true;
    } else {
      console.log(`❌ Erro: ${response.data.message}`);
      return false;
    }
  }
  return false;
}

// Função para enviar mensagem direta
async function sendDirectMessage(dstUser, message) {
  const request = {
    service: 'message',
    data: {
      src: currentUser,
      dst: dstUser,
      message: message,
      timestamp: Math.floor(Date.now() / 1000)
    }
  };

  console.log('\n📤 Enviando mensagem direta...');
  const response = await sendRequest(request);

  if (response) {
    if (response.data.status === 'OK') {
      console.log(`✅ Mensagem enviada para ${dstUser}`);
      console.log(`⏰ Timestamp: ${new Date(response.data.timestamp * 1000).toLocaleString()}`);
      return true;
    } else {
      console.log(`❌ Erro: ${response.data.message}`);
      return false;
    }
  }
  return false;
}

// Menu principal
function showMenu() {
  console.log('\n' + '='.repeat(60));
  console.log('📱 SISTEMA DE MENSAGENS - MENU PRINCIPAL');
  console.log('='.repeat(60));
  if (currentUser) {
    console.log(`👤 Usuário: ${currentUser}`);
    console.log(`📬 Canais inscritos: ${Array.from(subscribedChannels).join(', ') || 'nenhum'}`);
  }
  console.log('\nOpções:');
  if (!currentUser) {
    console.log('  1. Fazer login');
  } else {
    console.log('  2. Listar usuários cadastrados');
    console.log('  3. Criar novo canal');
    console.log('  4. Listar canais disponíveis');
    console.log('  5. Inscrever em canal');
    console.log('  6. Publicar mensagem em canal');
    console.log('  7. Enviar mensagem direta');
  }
  console.log('  0. Sair');
  console.log('='.repeat(60));
}

// Função para processar a escolha do usuário
function processChoice(choice) {
  switch (choice) {
    case '1':
      if (!currentUser) {
        rl.question('\n📝 Digite seu nome de usuário: ', async (username) => {
          if (username.trim()) {
            await login(username.trim());
          } else {
            console.log('❌ Nome de usuário não pode ser vazio!');
          }
          showMenuAndPrompt();
        });
        return;
      }
      break;
    
    case '2':
      if (currentUser) {
        listUsers().then(() => showMenuAndPrompt());
        return;
      }
      break;
    
    case '3':
      if (currentUser) {
        rl.question('\n📝 Digite o nome do canal a criar: ', async (channelName) => {
          if (channelName.trim()) {
            await createChannel(channelName.trim());
          } else {
            console.log('❌ Nome do canal não pode ser vazio!');
          }
          showMenuAndPrompt();
        });
        return;
      }
      break;
    
    case '4':
      if (currentUser) {
        listChannels().then(() => showMenuAndPrompt());
        return;
      }
      break;
    
    case '5':
      if (currentUser) {
        rl.question('\n📝 Digite o nome do canal para se inscrever: ', async (channelName) => {
          if (channelName.trim()) {
            await subscribeChannel(channelName.trim());
          } else {
            console.log('❌ Nome do canal não pode ser vazio!');
          }
          showMenuAndPrompt();
        });
        return;
      }
      break;
    
    case '6':
      if (currentUser) {
        rl.question('\n📝 Canal: ', (channel) => {
          if (!channel.trim()) {
            console.log('❌ Nome do canal não pode ser vazio!');
            showMenuAndPrompt();
            return;
          }
          rl.question('📝 Mensagem: ', async (message) => {
            if (message.trim()) {
              await publishMessage(channel.trim(), message.trim());
            } else {
              console.log('❌ Mensagem não pode ser vazia!');
            }
            showMenuAndPrompt();
          });
        });
        return;
      }
      break;
    
    case '7':
      if (currentUser) {
        rl.question('\n📝 Destinatário: ', (dst) => {
          if (!dst.trim()) {
            console.log('❌ Nome do usuário não pode ser vazio!');
            showMenuAndPrompt();
            return;
          }
          rl.question('📝 Mensagem: ', async (message) => {
            if (message.trim()) {
              await sendDirectMessage(dst.trim(), message.trim());
            } else {
              console.log('❌ Mensagem não pode ser vazia!');
            }
            showMenuAndPrompt();
          });
        });
        return;
      }
      break;
    
    case '0':
      console.log('\n👋 Encerrando cliente... Até logo!');
      reqSocket.close();
      subSocket.close();
      rl.close();
      process.exit(0);
      return;
    
    default:
      console.log('\n❌ Opção inválida!');
      break;
  }
  
  showMenuAndPrompt();
}

// Função para mostrar menu e esperar input
function showMenuAndPrompt() {
  showMenu();
  rl.question('\n➡️  Escolha uma opção: ', processChoice);
}

// Inicialização
async function init() {
  console.log('\n🚀 Iniciando cliente de mensagens...');
  
  const serverUrl = process.env.SERVER_URL || 'tcp://server:5555';
  const brokerUrl = process.env.BROKER_URL || 'tcp://broker:5558';
  
  console.log(`🔌 Conectando ao servidor: ${serverUrl}`);
  console.log(`🔌 Conectando ao broker: ${brokerUrl}`);
  
  try {
    await reqSocket.connect(serverUrl);
    await subSocket.connect(brokerUrl);
    connected = true;
    console.log('✅ Conectado ao servidor e broker com sucesso!');
    
    // Iniciar recebimento de mensagens em background
    receiveMessages().catch(err => {
      console.error('Erro no recebimento de mensagens:', err.message);
    });
    
    // Aguardar um pouco para garantir conexão
    setTimeout(() => {
      showMenuAndPrompt();
    }, 500);
    
  } catch (error) {
    console.error('❌ Erro ao conectar:', error.message);
    console.log('💡 Verifique se o servidor e broker estão rodando.');
    process.exit(1);
  }
}

// Tratamento de sinais de término
process.on('SIGINT', () => {
  console.log('\n\n👋 Recebido sinal de término. Encerrando...');
  reqSocket.close();
  subSocket.close();
  rl.close();
  process.exit(0);
});

// Iniciar aplicação
init();