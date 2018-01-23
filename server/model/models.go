package model

type User struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatarURL"`
}

type Experiment struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Assignment struct {
	ID       int        `json:"id"`
	PairID   int        `json:"pairId"`
	Answer   NullString `json:"answer"`
	Duration NullInt64  `json:"duration,omitempty"`
}

type FilePair struct {
	ID   int    `json:"id"`
	Diff string `json:"diff"`
}
