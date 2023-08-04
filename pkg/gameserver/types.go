package gameserver

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type PlayerState struct {
	Health   int      `json:"health"`
	Position Position `json:"position"`
}

type WSMessage struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}
