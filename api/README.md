# Meeting Stone API（Go）

本目录是 `meeting-stone-web` 的后端服务骨架（对应 `docs/02_ARCHITECTURE.md` 的 `api/`）。

## 本地启动

在 `meeting-stone-web/api` 目录执行：

```bash
go run ./cmd/server
```

默认监听：`0.0.0.0:8080`

可选环境变量：

- `PORT`: 端口（默认 8080）
- `HOST`: 监听地址（默认 `0.0.0.0`）
- `APP_ENV`: `dev|prod`（默认 `dev`）
- `WEB_ORIGIN`: 允许的前端 Origin（用于 dev 跨域直连，例如 `http://localhost:5173`）

## 已提供的最小接口（用于前端联调）

- `GET /healthz`
- `GET /api/v1/health`
- `GET /api/v1/games`
- `GET /api/v1/games/{gameKey}/manifest`
- `GET /api/v1/games/{gameKey}/activities`
- `GET /api/v1/games/{gameKey}/lobby/activities`

## 下一步（按 docs 落地）

- 接入 DB（activities / signups / lineups / templates / logs / outbox）
- 鉴权：JWT + Refresh Cookie（与 `docs/02_ARCHITECTURE.md` 对齐）
- 访问控制：访问码 hash + 邀请链接 token（`docs/04_DATA_MODEL.md`）

