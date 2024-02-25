# OSwithGO

# Simulador de Gerenciamento de Processos

## Descrição Breve do Projeto:
O simulador de escalonamento de processos é uma ferramenta que permite aos usuários criar, gerenciar e simular o escalonamento de processos em um sistema operacional. Ele oferece a capacidade de criar processos manualmente ou gerar processos aleatórios, escolher entre diferentes algoritmos de escalonamento, Round Robin, selecionar o tempo de quantum para o algoritmo executar a simulação do escalonamento.

### Modelagem dos Processos como Goroutines:
- Os processos são representados como goroutines em Go.
- Cada goroutine representa um processo que pode ser CPU bound ou I/O bound.
  
### Controle de Acesso à CPU e ao I/O:
##### Starvation:
- Starvation é um problema comum em sistemas de escalonamento de processos.
- Ele ocorre quando um processo fica permanentemente impedido de executar devido à priorização de outros processos.
- Isso pode resultar em processos esperando indefinidamente para serem executados, causando atrasos e impactando o desempenho geral do sistema.

##### Evitando o Starvation
- O controle é realizado via o algoritmo Round Robin customizado.
- O algoritmo Round Robin é adaptado para lidar com processos I/O bound e CPU bound simultaneamente.
- Inicializa os primeiros processos I/O e CPU bound na lista de execução.
- Cada processo recebe um tempo de execução conforme o quantum definido.
- Se um processo não terminar sua execução dentro do quantum, ele é retornado ao final da lista de execução, e em seguida, o próximo processo da lista é chamado para execução.
- Esse mecanismo evita o starvation dos processos, pois proporciona uma parcela igual de tempo de uso da CPU e do I/O para todos os processos, garantindo que todos tenham a oportunidade de serem executados. Isso impede que um processo monopolize os recursos e que outros processos fiquem impedidos de progredir, resultando em uma distribuição mais justa e equilibrada do tempo de CPU e I/O entre os processos.

### Recursos Compartilhados:
- Os recursos compartilhados incluem canais `cpuResource` e `ioResource`, utilizados para controlar o acesso à CPU e aos recursos de I/O, respectivamente.
- Os processos compartilham esses recursos para coordenar suas operações de CPU e I/O.
- Esses recursos estão nos arquivos `cpu.go` e `io.go`.
  - Eles funcionam da seguinte forma: criam um arquivo txt referente à CPU ou I/O.
  - Durante a simulação, escrevem o nome do processo, seu ID e o tempo de uso no arquivo correspondente.
  - Após terminar, liberam o recurso para que outros processos possam utilizá-lo.

### Controle da Concorrência:
- A concorrência é controlada utilizando canais para sincronização entre as goroutines que representam os processos.
- Os canais são utilizados para coordenar o acesso aos recursos compartilhados, garantindo que apenas um processo possa acessar a CPU ou o dispositivo de I/O por vez.

### Implementação Parametrizada:
- O algoritmo de escalonamento e o tempo de quantum são parametrizados pelo menu do usuário.
- Os recursos compartilhados e a sincronização são gerenciados globalmente no código, garantindo consistência e controle adequado da concorrência.

## Descrição Detalhada:

1. **Criação de Processos:**
   - Os usuários podem criar processos manualmente, inserindo informações como nome, prioridade, se é I/O bound ou não, e o tempo total de CPU.
   - Também é possível gerar um número especificado de processos aleatórios, onde as informações são geradas de forma aleatória, como nome, prioridade, se é I/O bound ou não, e o tempo total de CPU.

2. **Escolha de Algoritmo de Escalonamento:**
   - Os usuários podem escolher o algoritmo de escalonamento disponível: Round Robin.
   - As ações de escalonamento são definidas de acordo com as regras do algoritmo selecionado.

3. **Seleção de Tempo de Quantum:**
   - Os usuários podem definir o tempo de quantum, que é o intervalo de tempo máximo que um processo pode ocupar a CPU e I/O antes de ser preempedido.

4. **Execução da Simulação:**
   - A simulação é iniciada quando os usuários optam por executar o escalonamento.
   - Durante a simulação, a fila de processos prontos é exibida, e o algoritmo de escalonamento selecionado é aplicado aos processos na fila.
   - Os processos são executados de acordo com as regras do algoritmo escolhido até que todos os processos tenham sido concluídos.
   - Durante a execução da simulação, os processos I/O bound e CPU bound são tratados simultaneamente. 
     - Se um processo estiver em estado de I/O bound, ele iniciará sua operação de I/O enquanto aguarda a CPU.
     - Enquanto isso, os processos CPU bound utilizarão a CPU.
     - Se um processo I/O bound precisar da CPU, ele aguardará até que a CPU esteja disponível.
     - Se um processo CPU bound precisar de I/O, ele aguardará até que o dispositivo de I/O esteja livre.
   - Essa administração acontece na função `executeProcess`, localizada no arquivo `scheduler.go`. 
     - Essa função utiliza `cpuResource = make(chan bool, 1)` para controlar o recurso da CPU e `ioResource = make(chan bool, 1)` para controlar o recurso de I/O.
     - Um canal em Go é uma estrutura de dados que permite a comunicação e sincronização entre goroutines, que são "threads" leves. 
     - A função `executeProcess` utiliza goroutines para controlar a utilização dos recursos de CPU e I/O de forma sincronizada, garantindo um comportamento correto na simulação.
    - As goroutines são como "threads" em Go, que são gerenciadas pelo próprio Go Runtime, permitindo uma concorrência eficiente.
   - Além disso, a variável `wg` do tipo `sync.WaitGroup` é utilizada para coordenar e esperar o término das goroutines, garantindo que a simulação seja concluída corretamente.
.
1. **Cálculo do Tempo Médio de Espera:**
   - Após a conclusão da simulação, o tempo médio de espera dos processos é calculado e exibido.
   - O tempo médio de espera representa o tempo médio que os processos esperam na fila de prontos antes de serem executados.

2. **Menu de Opções Interativo:**
   - Um menu interativo permite que os usuários naveguem pelas diferentes funcionalidades do simulador, como criar processos, escolher algoritmos, selecionar o tempo de quantum, iniciar a simulação, criar processos aleatórios e sair do programa.

