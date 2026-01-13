package httpserver

type GameInfo struct {
	GameKey    string `json:"gameKey"`
	Name       string `json:"name"`
	GameFamily string `json:"gameFamily,omitempty"`
}

type GameManifest struct {
	GameKey         string       `json:"gameKey"`
	GameFamily      string       `json:"gameFamily,omitempty"`
	ManifestVersion int          `json:"manifestVersion,omitempty"`
	Roles           []Role       `json:"roles,omitempty"`
	Jobs            []Job        `json:"jobs,omitempty"`
	Specs           []Spec       `json:"specs,omitempty"`
	Features        *Features    `json:"features,omitempty"`
	PartySizes      []int        `json:"partySizes,omitempty"`
	RosterShape     *RosterShape `json:"rosterShape,omitempty"`
	BuffComposition any          `json:"buffComposition,omitempty"`
}

type Role struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

type Job struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Color string `json:"color,omitempty"`
}

type Spec struct {
	JobKey string `json:"jobKey"`
	Key    string `json:"key"`
	Label  string `json:"label"`
}

type Features struct {
	ServerRegion    bool   `json:"serverRegion,omitempty"`
	Faction         bool   `json:"faction,omitempty"`
	BuffComposition bool   `json:"buffComposition,omitempty"`
	ScoreProvider   string `json:"scoreProvider,omitempty"`
}

type RosterShape struct {
	Type string `json:"type"`
	Cols int    `json:"cols"`
}

type ActivityOwner struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ActivitySummary struct {
	ID         int64          `json:"id"`
	GameKey    string         `json:"gameKey"`
	Title      string         `json:"title"`
	Deadline   string         `json:"deadline"`
	Size       int            `json:"size"`
	IsPublic   bool           `json:"isPublic"`
	AccessMode string         `json:"accessMode,omitempty"`
	Owner      *ActivityOwner `json:"owner,omitempty"`
}

type ListResponse[T any] struct {
	Data    []T   `json:"data"`
	HasMore *bool `json:"hasMore,omitempty"`
}
