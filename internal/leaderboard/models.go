package leaderboard

type SubmitScoreRequest struct {
	PlayerID      string `json:"player_id"`
	PlayerName    string `json:"player_name"`
	Score         int    `json:"score"`
	Mode          string `json:"mode"`
	Skin          string `json:"skin"`
	ClientVersion string `json:"client_version"`
	DurationMS    int64  `json:"duration_ms"`
}

type Entry struct {
	Rank       int    `json:"rank"`
	PlayerName string `json:"player_name"`
	Score      int    `json:"score"`
	Mode       string `json:"mode"`
	Skin       string `json:"skin"`
	UpdatedAt  string `json:"updated_at"`
}

type SubmitScoreResponse struct {
	Accepted bool  `json:"accepted"`
	Improved bool  `json:"improved"`
	Rank     int   `json:"rank"`
	Entry    Entry `json:"entry"`
}

type LeaderboardResponse struct {
	Mode    string  `json:"mode"`
	Entries []Entry `json:"entries"`
}

type VersionInfoResponse struct {
	MinClientVersion    string `json:"min_client_version"`
	LatestClientVersion string `json:"latest_client_version"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
