# 06 迁移与里程碑计划

## 1. 迁移目标

- 将现有“集合石（milestone）”能力从旧系统形态，迁移为可扩展的 Web 平台（支持多游戏）。
- 尽可能复用现有数据（milestone_*、关注/黑名单/管理员关系），同时引入 Web 侧新增能力（访问码/邀请链接、通知 outbox、OAuth 身份）。

## 2. 两种迁移路线

### 路线 A：沿用同一数据库（低成本）

适用：你希望快速上线，且现有 DB 可持续演进。

- 直接复用 milestone 表：
  - `milestone_activity`
  - `milestone_detail`
  - `milestone_lineup`
  - `milestone_lineup_slot`
  - `milestone_signup_group`
  - `milestone_signup_group_char`
  - `milestone_report`
- 复用社交/协作表：
  - `account_follows`
  - `account_blacklists`
  - `milestone_assistants`

新增表（Web 扩展）：
- `activity_invites`（邀请链接）
- `activity_access_grants`（可选：访问授权记录）
- `notifications` / `notification_deliveries` / `user_notification_settings`
- `account_identities`（OAuth 绑定：google/discord/bnet）

优点：
- 不需要搬数据
- 可渐进增加字段（如 access_mode/rules_json/payload_json）

风险：
- 旧系统字段命名与新平台契约不完全一致，需要在 API 层做适配/兼容

### 路线 B：新库/新表（更干净）

适用：你更重视长期整洁，愿意一次性做迁移脚本。

做法：
- 新建一套 Core 表结构（见 `04_DATA_MODEL.md`）
- 写迁移脚本把 milestone_* 与社交关系导入新表

优点：
- 数据结构与 API/业务对齐
- 多游戏扩展更自然

风险：
- 迁移成本更高
- 需要明确回滚策略与灰度策略

## 3. MVP 建议选择

建议先用 **路线 A（同库演进）** 快速跑通 MVP，然后再评估是否需要迁移到路线 B。

理由：
- 你现有 milestone 逻辑已在生产验证
- Web MVP 的主要不确定性在“产品形态与全球化”而不是 DB 结构

## 4. 关键兼容点（从旧到新）

### 4.1 gameKey / X-Game

现状：
- 现有后端通过 header `X-Game` 区分 game（在 milestone 数据里也有 `game_type` 字段）

新平台建议：
- API 路由以 `gameKey` path 参数为准（例如 `/games/wow/...`）
- DB 仍保留 `game_type/game_key` 字段作为隔离索引
- 兼容期：允许 `X-Game` 与 path 任选其一，但最终收敛到 path

### 4.2 not_share → 访问码/邀请链接

现状：
- 小程序利用生态限制“转发”

Web：
- 引入 `access_mode` 与 `access_code_hash`
- 邀请链接 token 可撤销与过期

迁移策略：
- 旧数据中 `not_share=1` 的活动可默认映射为 `access_mode=code`（并生成随机访问码），由团长自行重置/配置（或保持 private，仅通过链接访问）

### 4.3 报名字段 payload 扩展

现状：
- milestone_detail 的字段较偏 WoW（class/spec/wcl_score/item_level）

新平台：
- 逐步增加 `payload_json`（或通过 JSON 字段映射），把游戏特有字段放入 payload

## 5. 灰度发布与回滚

### 5.1 灰度策略（推荐）

- Phase 1：只读 Web（大厅 + 活动详情）→ 验证海外访问与 i18n/时区
- Phase 2：开放报名（单人报名）→ 验证限制与风控
- Phase 3：开放阵容编排 + AI 排班 → 验证高频交互
- Phase 4：开放报名组与批量报名、通知 Email

### 5.2 回滚策略

若采用路线 A（同库）：
- 回滚仅需停止 Web 新服务，不会破坏旧服务数据（但要避免新增字段导致旧代码异常）
- 新增表（invites/notifications）不影响旧服务

若采用路线 B（新库）：
- 回滚可切回旧服务与旧库
- 需要明确数据双写或导回策略（一般不建议双写，先只读/灰度写）

## 5.5 现有表结构映射

以下是现有 `milestone_*` 表与新平台模型的字段映射关系：

### milestone_activity → activities

| 旧字段 | 新字段 | 说明 |
|--------|--------|------|
| ID | id | 主键 |
| game_type | game_key | 游戏标识（可能需要映射：classic → wow-classic） |
| uid | owner_id | 团长 ID |
| title | title | 标题 |
| description | description | 描述 |
| deadline | deadline_at | 截止时间 |
| size | size | 团队规模 |
| is_public | is_public | 是否公开 |
| not_share | access_mode | `not_share=1` → `access_mode=private/code` |
| server | server_region_json.server | 服务器 |
| region | server_region_json.region | 地区 |
| faction | server_region_json.faction | 阵营 |
| player_limit_json | player_limit_json | 人员限制（格式兼容） |
| time_list_json | lineup_list_json | 子活动列表 |
| show_contact | visibility_json.showContact | 是否公开联系方式 |

### milestone_detail → signup_details

| 旧字段 | 新字段 | 说明 |
|--------|--------|------|
| ID | id | 主键 |
| activity_id | activity_id | 活动 ID |
| uid | user_id | 用户 ID |
| character_name | character_name | 角色名 |
| class | job_key | 职业 |
| spec | spec_key | 专长 |
| sub_spec | sub_spec_key | 副专长 |
| wcl_score | score_value + score_source="wcl" | WCL 分数 |
| item_level | item_level | 装等 |
| status | status | 状态（1=signed, 2=leave, 5=cancel） |
| desc | desc | 备注 |
| lineup_list_str | lineup_allow_list | 指定车队 |
| show_contact | show_contact_to_leaders | 是否公开联系方式 |
| register_type | payload_json.registerType | 报名类型（1=绑定角色, 2=自定义） |
| wcl_id | payload_json.characterRef | 角色引用 |

### milestone_lineup → lineups

| 旧字段 | 新字段 | 说明 |
|--------|--------|------|
| ID | id | 主键 |
| activity_id | activity_id | 活动 ID |
| name | name | 阵容名 |

### milestone_lineup_slot → lineup_slots

| 旧字段 | 新字段 | 说明 |
|--------|--------|------|
| ID | id | 主键 |
| lineup_id | lineup_id | 阵容 ID |
| index | slot_index | 槽位索引 |
| detail_id | detail_id | 报名记录 ID |

### milestone_signup_group → signup_groups

| 旧字段 | 新字段 | 说明 |
|--------|--------|------|
| ID | id | 主键 |
| uid | user_id | 用户 ID |
| name | name | 组名 |
| desc | description | 描述 |

### milestone_signup_group_char → signup_group_members

| 旧字段 | 新字段 | 说明 |
|--------|--------|------|
| ID | id | 主键 |
| group_id | group_id | 报名组 ID |
| character_name | character_name | 角色名 |
| class | job_key | 职业 |
| spec | spec_key | 专长 |

### milestone_assistants → activity_assistants

| 旧字段 | 新字段 | 说明 |
|--------|--------|------|
| ID | id | 主键 |
| uid | owner_id | 团长 ID |
| assistant_uid | assistant_id | 管理员 ID |
| status | status | 状态（invited/accepted/dismissed） |

## 6. 里程碑（与 PRD 对齐）

- **M1**：账号体系（Email+OAuth）+ 大厅搜索 + 活动详情只读
- **M2**：活动创建/模板 + 单人报名 + 日志
- **M3**：报名组 + 批量报名
- **M4**：阵容编排 + 保存/批量保存 + 导出
- **M5**：AI 排班（grid/constraint）
- **M6**：关注/黑名单/管理员协作 + Email 通知完善

## 7. 里程碑详细任务分解

### M1：账号体系 + 大厅 + 活动详情只读

**后端任务**：
- [ ] OAuth 集成（Battle.net / Discord / Google）
- [ ] Email 注册/登录
- [ ] JWT + Refresh Token 实现
- [ ] `GET /games` + `GET /games/{gameKey}/manifest`
- [ ] `GET /games/{gameKey}/lobby/activities`（大厅搜索）
- [ ] `GET /games/{gameKey}/activities/{id}`（活动详情只读）
- [ ] 访问码校验 `POST /activities/{id}/access/code`

**前端任务**：
- [ ] 登录/注册页面
- [ ] 首页（大厅搜索 + 结果列表）
- [ ] 活动详情页（只读报名列表）
- [ ] 访问码输入页

### M2：活动创建/模板 + 单人报名 + 日志

**后端任务**：
- [ ] `POST/PATCH/DELETE /activities`
- [ ] `GET/POST/PATCH/DELETE /templates`
- [ ] `POST /activities/from-template`
- [ ] `POST /activities/{id}/signup`
- [ ] `GET /activities/{id}/logs`

**前端任务**：
- [ ] 创建/编辑活动表单
- [ ] 模板管理页面
- [ ] 报名弹窗（单人报名）
- [ ] 活动日志页面

### M3：报名组 + 批量报名

**后端任务**：
- [ ] `GET/POST/PATCH/DELETE /signup-groups`
- [ ] `POST /activities/{id}/signup-with-group`

**前端任务**：
- [ ] 报名组管理页面
- [ ] 批量报名 Tab

### M4：阵容编排 + 保存/批量保存 + 导出

**后端任务**：
- [ ] `GET /activities/{id}/lineups`
- [ ] `PUT /lineups/{lineupId}`
- [ ] `PUT /activities/{id}/lineups:batch`
- [ ] `GET /activities/{id}/export`
- [ ] `GET /lineups/{lineupId}/mrt`

**前端任务**：
- [ ] 阵容编排页面（拖拽排位）
- [ ] 角色筛选面板
- [ ] 团队Buff统计面板
- [ ] 导出功能

### M5：AI 排班

**后端任务**：
- [ ] `POST /activities/{id}/auto-lineup`
- [ ] `PUT /auto-lineup-template`

**前端任务**：
- [ ] AI 排班弹窗
- [ ] 阵容模板编辑器

### M6：社交 + 通知

**后端任务**：
- [ ] 关注/黑名单/管理员 API
- [ ] Notification Outbox + Worker
- [ ] Email 模板（i18n）

**前端任务**：
- [ ] 联系人页面
- [ ] 通知设置页面

