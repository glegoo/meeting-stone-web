package httpserver

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func handleGetGames(w http.ResponseWriter, _ *http.Request) {
	Ok(w, []GameInfo{
		{GameKey: "wow-classic", Name: "World of Warcraft (Classic)", GameFamily: "wow"},
		{GameKey: "wow-retail", Name: "World of Warcraft (Retail)", GameFamily: "wow"},
		{GameKey: "ff14", Name: "Final Fantasy XIV"},
	})
}

func handleGetManifest(w http.ResponseWriter, r *http.Request) {
	gameKey := chi.URLParam(r, "gameKey")
	manifest, ok := getManifestByGameKey(gameKey)
	if !ok {
		Fail(w, http.StatusNotFound, 404, "not_found", "unknown gameKey")
		return
	}
	Ok(w, manifest)
}

func handleGetActivities(w http.ResponseWriter, r *http.Request) {
	// MVP：先给前端联调用的 mock；后续接 DB 时替换为 repo 查询。
	gameKey := chi.URLParam(r, "gameKey")
	now := time.Now().UTC()

	items := []ActivitySummary{
		{
			ID:         1,
			GameKey:    gameKey,
			Title:      "示例活动：今晚 20:00 集合",
			Deadline:   now.Add(6 * time.Hour).Format(time.RFC3339),
			Size:       25,
			IsPublic:   true,
			AccessMode: "public",
			Owner:      &ActivityOwner{ID: "demo-owner", Name: "Alice"},
		},
	}

	// 保持与前端兼容：返回 {data, hasMore}
	hasMore := false
	Ok(w, ListResponse[ActivitySummary]{Data: items, HasMore: &hasMore})
}

func handleGetLobbyActivities(w http.ResponseWriter, r *http.Request) {
	// MVP：与 /activities 同结构，先返回空列表，避免 UI 报错。
	hasMore := false
	Ok(w, ListResponse[ActivitySummary]{Data: []ActivitySummary{}, HasMore: &hasMore})
}

func getManifestByGameKey(gameKey string) (GameManifest, bool) {
	switch gameKey {
	case "wow-classic":
		return GameManifest{
			GameKey:         "wow-classic",
			GameFamily:      "wow",
			ManifestVersion: 1,
			Roles: []Role{
				{Key: "tank", Label: "坦克"},
				{Key: "healer", Label: "治疗"},
				{Key: "dps", Label: "输出"},
			},
			Features: &Features{
				ServerRegion:    true,
				Faction:         true,
				BuffComposition: false,
				ScoreProvider:   "wcl",
			},
			PartySizes:  []int{10, 25, 40},
			RosterShape: &RosterShape{Type: "grid", Cols: 5},
		}, true
	case "wow-retail":
		return GameManifest{
			GameKey:         "wow-retail",
			GameFamily:      "wow",
			ManifestVersion: 1,
			Roles: []Role{
				{Key: "tank", Label: "Tank"},
				{Key: "healer", Label: "Healer"},
				{Key: "dps", Label: "DPS"},
			},
			Features: &Features{
				ServerRegion:    true,
				Faction:         true,
				BuffComposition: true,
				ScoreProvider:   "wcl",
			},
			PartySizes:  []int{10, 20, 25, 30},
			RosterShape: &RosterShape{Type: "grid", Cols: 5},
		}, true
	case "ff14":
		return GameManifest{
			GameKey:         "ff14",
			ManifestVersion: 1,
			Roles: []Role{
				{Key: "tank", Label: "Tank"},
				{Key: "healer", Label: "Healer"},
				{Key: "dps", Label: "DPS"},
			},
			Features: &Features{
				ServerRegion:    false,
				Faction:         false,
				BuffComposition: false,
				ScoreProvider:   "fflogs",
			},
			PartySizes:  []int{4, 8},
			RosterShape: &RosterShape{Type: "grid", Cols: 4},
		}, true
	default:
		return GameManifest{}, false
	}
}
