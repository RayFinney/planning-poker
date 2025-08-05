package planning

type Planning struct {
	Id            string         `json:"id"`
	LastConnected string         `json:"lastConnected"`
	CreatedAt     string         `json:"created_at"`
	Owner         Player         `json:"owner"`
	Players       []Player       `json:"players"`
	Revealed      bool           `json:"revealed"`
	MyVote        int            `json:"myVote"` // My vote is the vote of the player who is currently connected
	Votes         map[string]int `json:"votes"`  // Vote key is player ID
	HiddenVotes   map[string]int `json:"-"`      // Vote key is player ID
}

type Player struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	IsOwner bool   `json:"-"` // IsOwner indicates if the player is the owner of the planning
}
