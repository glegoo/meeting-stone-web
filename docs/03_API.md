# 03 API 设计（按 gameKey 隔离）

> 目标：对前端提供稳定、可扩展（多游戏）的 HTTP API。所有业务接口以 `gameKey` 作为路径隔离维度，鉴权以登录态为准（不信任前端传 `uid`）。

## 0. 统一约定

### 0.1 Base URL

- `/api/v1`

### 0.2 鉴权

- `Authorization: Bearer <access_token>`（建议）
- Refresh Token 通过 HttpOnly Cookie（建议）

### 0.3 响应格式（建议）

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

错误：
```json
{
  "code": 400,
  "message": "validation_error",
  "error": "具体错误信息"
}
```

### 0.4 gameKey

示例：
- `wow-retail`
- `wow-classic`
- `ff14`

说明：
- `gameKey` 用于表达“规则集（ruleset）”。同一 IP（比如 WoW）不同版本/赛季若规则差异显著，应拆成不同 `gameKey`。
- 可通过 `gameFamily`/`parentGameKey` 把多个 `gameKey` 聚合到同一大类（用于 UI/统计，不参与业务隔离）。

## 1. Games 与 Manifest

### 1.1 获取支持的游戏列表

- `GET /api/v1/games`

返回（示例）：
```json
{
  "code": 200,
  "data": [
    { "gameKey": "wow-retail", "name": "World of Warcraft (Retail)", "gameFamily": "wow" },
    { "gameKey": "wow-classic", "name": "World of Warcraft (Classic)", "gameFamily": "wow" },
    { "gameKey": "ff14", "name": "Final Fantasy XIV" }
  ]
}
```

### 1.2 获取游戏 Manifest（前端渲染契约）

- `GET /api/v1/games/{gameKey}/manifest`

Manifest 建议字段（示例，按需裁剪）：
```json
{
  "code": 200,
  "data": {
    "gameKey": "wow-retail",
    "gameFamily": "wow",
    "manifestVersion": 1,
    "roles": [{ "key": "tank", "label": "Tank" }],
    "jobs": [{ "key": "Paladin", "label": "Paladin", "color": "#F58CBA" }],
    "specs": [{ "jobKey": "Paladin", "key": "Protection", "label": "Protection" }],
    "features": {
      "serverRegion": true,
      "faction": true,
      "buffComposition": true,
      "scoreProvider": "wcl"
    },
    "partySizes": [10, 20, 25, 40],
    "rosterShape": { "type": "grid", "cols": 5 }
  }
}
```

## 2. 活动（Activities）

### 2.1 获取活动列表（我创建的/我参与的）

- `GET /api/v1/games/{gameKey}/activities?type=created|joined&page=1&pageSize=10&sort=desc&includeExpire=false`

返回（简化示例）：
```json
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "gameKey": "wow",
      "title": "XX团",
      "deadline": "2026-01-13T12:00:00Z",
      "size": 25,
      "isPublic": true,
      "accessMode": "code",
      "owner": { "id": "uuid", "name": "Alice" }
    }
  ],
  "hasMore": true
}
```

### 2.2 大厅搜索（公开活动）

- `GET /api/v1/games/{gameKey}/lobby/activities?server=...&region=...&faction=1&date=2026-01-13&size=25&keyword=xx&page=1`

说明：
- server/region/faction 是否启用由 manifest 的 `features.serverRegion/faction` 决定。
- date 建议按“用户时区的当天”解释（PRD 需定）。

### 2.3 创建活动

- `POST /api/v1/games/{gameKey}/activities`

Body（示例）：
```json
{
  "title": "2026-01-13 Tue XX",
  "description": "集合时间…",
  "deadline": "2026-01-13T12:00:00Z",
  "size": 25,
  "visibility": { "isPublic": true },
  "access": { "mode": "code", "code": "123456" },
  "serverRegion": { "server": "Firemaw", "region": "EU", "faction": 1 },
  "lineupList": ["第1车", "第2车"],
  "playerLimits": [
    { "type": "total", "limit": 25 },
    { "type": "role", "role": "tank", "limit": 2 }
  ],
  "rules": {}
}
```

返回：
- 创建后的活动摘要（不返回明文 code）

### 2.4 更新活动

- `PATCH /api/v1/games/{gameKey}/activities/{id}`

说明：
- 变更触发通知：activity updated（按频控/去抖）
- 访问码变更会使旧访问码失效

### 2.5 删除活动

- `DELETE /api/v1/games/{gameKey}/activities/{id}`

## 2.5 模板（Templates）

### 2.5.1 获取模板列表

- `GET /api/v1/games/{gameKey}/templates`

返回（示例）：
```json
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "gameKey": "wow-classic",
      "title": "ICC25H周六团",
      "description": "准备合剂、食物...",
      "size": 25,
      "lineupList": [
        { "dayOfWeek": 4, "time": "20:00", "name": "ICC25H上层" },
        { "dayOfWeek": 4, "time": "20:00", "name": "ICC25H下层" }
      ],
      "playerLimits": [
        { "type": "total", "limit": 25 },
        { "type": "role", "role": "tank", "limit": 2 }
      ],
      "serverRegion": { "server": "克尔苏加德", "region": "CN", "faction": 1 },
      "visibility": { "isPublic": true, "notShare": false }
    }
  ]
}
```

### 2.5.2 创建模板

- `POST /api/v1/games/{gameKey}/templates`

Body（示例）：
```json
{
  "title": "ICC25H周六团",
  "description": "准备合剂...",
  "size": 25,
  "lineupList": [
    { "dayOfWeek": 4, "time": "20:00", "name": "ICC25H上层" }
  ],
  "playerLimits": [...],
  "serverRegion": { "server": "克尔苏加德", "region": "CN", "faction": 1 },
  "visibility": { "isPublic": false, "notShare": true }
}
```

### 2.5.3 更新模板

- `PATCH /api/v1/games/{gameKey}/templates/{id}`

### 2.5.4 删除模板

- `DELETE /api/v1/games/{gameKey}/templates/{id}`

### 2.5.5 从模板创建活动

- `POST /api/v1/games/{gameKey}/activities/from-template`

Body：
```json
{
  "templateId": 1,
  "deadline": "2026-01-20T20:00:00Z",
  "overrides": {
    "title": "01月20日 周六 ICC25H"
  }
}
```

## 3. 活动详情（含访问控制）

### 3.1 获取活动详情

- `GET /api/v1/games/{gameKey}/activities/{id}`

说明：
- 若 `accessMode=code|invite` 且用户未授权，返回 403，并给出 `requiresAccess: true`

### 3.2 输入访问码（获取访问授权）

- `POST /api/v1/games/{gameKey}/activities/{id}/access/code`

Body：
```json
{ "code": "123456" }
```

返回：
- `accessGranted: true`（并建立会话授权，具体由实现决定）

### 3.3 使用邀请链接 token 授权

- `POST /api/v1/games/{gameKey}/activities/{id}/access/invite`

Body：
```json
{ "token": "invite_xxx" }
```

返回同上。

### 3.4 邀请链接管理（Owner/Assistant）

- `POST /api/v1/games/{gameKey}/activities/{id}/invites`（创建）
- `GET /api/v1/games/{gameKey}/activities/{id}/invites`（列表）
- `DELETE /api/v1/games/{gameKey}/activities/{id}/invites/{inviteId}`（撤销）

## 4. 报名（Signup Details）

### 4.1 单人报名

- `POST /api/v1/games/{gameKey}/activities/{id}/signup`

Body（示例）：
```json
{
  "registerType": "character|custom",
  "character": { "characterId": "xxx", "name": "Bob" },
  "job": "Paladin",
  "spec": "Protection",
  "subSpec": "Holy",
  "score": { "value": 98.7, "source": "wcl" },
  "itemLevel": 512,
  "desc": "时间OK",
  "showContactToLeaders": true,
  "lineupAllowList": [1,2]
}
```

返回：
- 报名结果与当前限制校验信息（可选）

### 4.2 报名状态变更（请假/取消）

- `PATCH /api/v1/games/{gameKey}/details/{detailId}`

Body（示例）：
```json
{ "status": "signed|leave|cancel" }
```

权限：
- 当事人可改自己
- Owner/Assistant 可移除他人（cancel）

### 4.3 报名组（Signup Groups）

- `GET /api/v1/games/{gameKey}/signup-groups`
- `POST /api/v1/games/{gameKey}/signup-groups`
- `PATCH /api/v1/games/{gameKey}/signup-groups/{groupId}`
- `DELETE /api/v1/games/{gameKey}/signup-groups/{groupId}`

### 4.4 报名组批量报名

- `POST /api/v1/games/{gameKey}/activities/{id}/signup-with-group`

Body：
```json
{
  "groupId": 123,
  "desc": "统一备注",
  "showContactToLeaders": true,
  "lineupAllowList": [1]
}
```

## 5. 阵容（Lineups）

### 5.1 获取阵容

- `GET /api/v1/games/{gameKey}/activities/{id}/lineups`

返回：
- lineup 列表 + slot 明细（可按需裁剪）

### 5.2 保存单个阵容

- `PUT /api/v1/games/{gameKey}/lineups/{lineupId}`

Body：
```json
{
  "name": "第1车",
  "slots": [
    { "index": 1, "detailId": 10 },
    { "index": 2, "detailId": 0 }
  ]
}
```

### 5.3 批量保存阵容

- `PUT /api/v1/games/{gameKey}/activities/{id}/lineups:batch`

### 5.4 AI 排班

- `POST /api/v1/games/{gameKey}/activities/{id}/auto-lineup`

Body（网格约束示例）：
```json
{
  "strategy": "grid",
  "grid": [
    { "index": 1, "constraint": { "baseType": 1 } },
    { "index": 2, "constraint": { "baseType": 1 } },
    { "index": 3, "constraint": { "baseType": 2 } },
    { "index": 4, "constraint": { "baseType": 2 } },
    { "index": 5, "constraint": { "baseType": 2 } },
    { "index": 6, "constraint": { "baseType": 4, "jobKey": "Shaman", "specKey": "Elemental" } },
    { "index": 7, "constraint": { "baseType": 4, "jobKey": "Priest", "specKey": "Shadow" } },
    { "index": 8, "constraint": { "baseType": 0 } },
    { "index": 9, "constraint": { "baseType": 5 } }
  ]
}
```

约束类型（baseType）：
- `0`：任意角色
- `1`：坦克
- `2`：治疗
- `3`：近战DPS
- `4`：远程DPS
- `5`：空位（不填人）

可选细化：
- `jobKey`：指定职业（如 Paladin）
- `specKey`：指定专长（如 Protection）

返回：
```json
{
  "code": 200,
  "data": [
    {
      "lineupId": 1,
      "lineupName": "第1车",
      "slots": "10,25,33,48,52,61,78,0,0,..."
    },
    {
      "lineupId": 2,
      "lineupName": "第2车",
      "slots": "12,26,35,49,55,..."
    }
  ]
}
```

说明：
- `slots` 是 detailId 列表，按 index 顺序排列，0 表示空位
- AI 排班结果**不自动保存**，前端需展示预览，用户确认后调用保存接口
- AI 会避免同一用户在同一阵容中出现多次（优先级较低时会放宽）
- AI 会优先匹配"指定本车"的角色

### 5.5 保存 AI 排班模板

用户可以保存 AI 排班的网格约束模板，下次复用。

- `PUT /api/v1/games/{gameKey}/auto-lineup-template`

Body：
```json
{
  "size": 25,
  "grid": [
    { "index": 1, "constraint": { "baseType": 1 } },
    { "index": 2, "constraint": { "baseType": 1 } }
  ]
}
```

说明：
- 按 `gameKey` + `userId` 存储
- 前端在打开 AI 排班弹窗时加载用户保存的模板

## 6. 日志与导出

### 6.1 获取活动日志

- `GET /api/v1/games/{gameKey}/activities/{id}/logs?page=1&pageSize=50`

返回（示例）：
```json
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "activityId": 100,
      "action": "signup",
      "actor": { "id": "uuid", "name": "Alice" },
      "target": { "characterName": "Arthas", "job": "DeathKnight", "spec": "Blood" },
      "detail": "报名成功",
      "createdAt": "2026-01-13T10:00:00Z"
    },
    {
      "id": 2,
      "action": "status_change",
      "actor": { "id": "uuid", "name": "Alice" },
      "target": { "characterName": "Arthas" },
      "detail": "请假",
      "createdAt": "2026-01-13T12:00:00Z"
    },
    {
      "id": 3,
      "action": "lineup_save",
      "actor": { "id": "uuid", "name": "Bob" },
      "target": { "lineupName": "第1车" },
      "detail": "保存阵容",
      "createdAt": "2026-01-13T14:00:00Z"
    }
  ],
  "hasMore": false
}
```

日志 action 类型：
- `signup`：报名
- `status_change`：状态变更（请假/取消/被移除）
- `lineup_save`：保存阵容
- `activity_update`：活动信息变更
- `assistant_add`：添加管理员
- `assistant_remove`：移除管理员

### 6.2 导出报名列表

- `GET /api/v1/games/{gameKey}/activities/{id}/export?format=csv|excel`

Query 参数：
- `format`：导出格式（csv/excel）
- `status`：筛选状态（signed/leave/all，默认 signed）
- `fields`：导出字段（逗号分隔，如 name,job,spec,score,itemLevel）

返回：
- `Content-Type: text/csv` 或 `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`
- 文件下载

### 6.3 导出 MRT 名单

- `GET /api/v1/games/{gameKey}/lineups/{lineupId}/mrt`

返回：
```json
{
  "code": 200,
  "data": {
    "text": "角色名1\n角色名2\n角色名3\n..."
  }
}
```

说明：
- 按 MRT 插件格式导出（适用于 WoW）
- 25人阵容按小队顺序排列（列优先）

## 7. 社交与协作

账号层（不随 game 拆分）：
- `POST /api/v1/accounts/follow`
- `POST /api/v1/accounts/unfollow`
- `POST /api/v1/accounts/addBlacklist`
- `POST /api/v1/accounts/removeBlacklist`
- `POST /api/v1/accounts/addMilestoneAssistant`（建议改名 assistants）

## 8. 通知设置（MVP=Email）

### 8.1 查询/更新通知设置

- `GET /api/v1/notifications/settings`
- `PATCH /api/v1/notifications/settings`

示例：
```json
{
  "emailEnabled": true,
  "events": {
    "deadlineReminder": true,
    "assistantInvite": true,
    "statusChanged": true,
    "activityUpdated": true
  },
  "digestWindowMinutes": 5,
  "timezone": "America/Los_Angeles",
  "language": "en"
}
```

### 8.2 测试发送（可选）

- `POST /api/v1/notifications/test-email`

