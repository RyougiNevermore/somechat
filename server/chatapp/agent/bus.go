package agent

import "sync"

type AgentLocalBus struct {
	Agents map[string]*Agent
	lock *sync.Mutex
}

func NewAgentLocalBus() *AgentLocalBus {
	bus := AgentLocalBus{}
	bus.Agents = make(map[string]*Agent)
	bus.lock = new(sync.Mutex)
	return &bus
}

func (bus *AgentLocalBus) Register(agent *Agent) {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	bus.Agents[agent.Id] = agent
}

func (bus *AgentLocalBus) Find(agentId string) (*Agent, bool) {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	agentInstance, has := bus.Agents[agentId]
	return agentInstance, has
}

func (bus *AgentLocalBus) UnRegister(agentId string)  {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	delete(bus.Agents, agentId)
}

