# 04 数据模型（表结构草案）

> 目标：支持多游戏（gameKey 隔离）、访问控制（访问码/邀请链接）、可扩展报名字段（payload_json）、可扩展规则（rules_json）、通知 Outbox（异步投递）。

以下为“建议的逻辑模型”，并不强制一次性全部实现；MVP 可先实现必要字段，后续再扩展。

## 1. 多游戏隔离

所有核心业务表都应包含：

- `game_key`：`varchar(32)`，索引（与业务查询组合索引）

原因：
- 同一用户可跨游戏创建活动/报名
- Manifest、字典与规则由 `game_key` 决定

### 1.1 gameKey 的粒度与聚合

- **建议**：`game_key` 表达“规则集（ruleset）”。例如魔兽世界的不同版本（职业/服务器/评分口径不同）应使用不同 `game_key`：`wow-retail`、`wow-classic`、`wow-sod` …\n
- 为了在产品层做聚合展示，建议增加一个“家族/父级”字段（不参与业务隔离，仅用于 UI/统计）：\n
  - `game_family`（或 `parent_game_key`），例如统一为 `wow`\n

可以落地为独立字典表（推荐）：

### 1.2 games（游戏/规则集字典表，推荐）

| 字段 | 类型 | 说明 |
|---|---|---|
| game_key | varchar | 主键：规则集标识（wow-retail） |
| game_family | varchar | 聚合标识（wow） |
| name | varchar | 展示名 |
| status | enum | active/disabled |
| manifest_version | int | 当前 manifest 版本（便于缓存与升级） |
| created_at/updated_at | datetime | |

## 2. 核心域模型（Core）

### 2.1 activities（活动）

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint | 主键 |
| game_key | varchar | 游戏标识 |
| owner_id | varchar | 账号 ID（uuid） |
| title | varchar | 标题 |
| description | text | 描述 |
| deadline_at | datetime | 截止时间（UTC） |
| size | int | 团队规模 |
| is_public | bool | 是否公开（大厅可见） |
| access_mode | enum | `public/private/code/invite` |
| access_code_hash | varchar | 访问码 hash（bcrypt/argon2） |
| access_code_updated_at | datetime | 用于失效/审计 |
| server_region_json | json | 服务器/大区/阵营（不同游戏差异大） |
| player_limit_json | json | 人员限制（总量/职责/职业/专长） |
| rules_json | json | 活动扩展规则（Adapter 解释） |
| created_at/updated_at/deleted_at | datetime | |

索引建议：
- `(game_key, is_public, deadline_at)`（大厅）
- `(game_key, owner_id, deadline_at)`（我的活动）

### 2.2 lineups（子活动/车队）

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint | 主键 |
| game_key | varchar | 冗余隔离字段（方便查询/校验） |
| activity_id | bigint | 外键 |
| name | varchar | 第1车/第2车… |
| created_at/updated_at/deleted_at | datetime | |

### 2.3 lineup_slots（阵容槽位）

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint | 主键 |
| lineup_id | bigint | 外键 |
| slot_index | int | 1..size |
| detail_id | bigint | 报名记录 ID（为空=空位） |
| created_at/updated_at/deleted_at | datetime | |

约束建议：
- `(lineup_id, slot_index)` 唯一

### 2.4 signup_details（报名记录）

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint | 主键 |
| game_key | varchar | 游戏标识 |
| activity_id | bigint | 外键 |
| user_id | varchar | 账号 ID（uuid） |
| character_name | varchar | 角色名（展示用） |
| job_key | varchar | 职业（游戏内 Job/Class） |
| spec_key | varchar | 专长/天赋 |
| sub_spec_key | varchar | 副专长/副天赋 |
| score_value | decimal | 归一化分数（可空） |
| score_source | varchar | `wcl/fflogs/custom/...` |
| item_level | int | 装等/物品等级（可空） |
| status | enum | `signed/leave/cancel` |
| desc | varchar/text | 备注 |
| lineup_allow_list | varchar | 允许的阵容列表（可先用字符串/JSON） |
| show_contact_to_leaders | bool | 是否向团长公开联系方式 |
| payload_json | json | 报名扩展字段（Adapter 解释） |
| created_at/updated_at/deleted_at | datetime | |

索引建议：
- `(activity_id, status)`（详情页列表）
- `(user_id, activity_id)`（我的报名）

### 2.5 signup_groups（报名组）

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint | 主键 |
| game_key | varchar | 游戏标识 |
| user_id | varchar | 账号 ID（uuid） |
| name | varchar | 组名 |
| description | varchar/text | 描述 |
| created_at/updated_at/deleted_at | datetime | |

### 2.6 signup_group_members（报名组角色）

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint | 主键 |
| group_id | bigint | 外键 |
| character_name | varchar | 角色名 |
| job_key/spec_key/sub_spec_key | varchar | 职业/专长 |
| character_ref_json | json | 绑定角色引用（不同游戏不同） |
| payload_json | json | 扩展字段 |
| created_at/updated_at/deleted_at | datetime | |

### 2.7 templates（活动模板）

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint | 主键 |
| game_key | varchar | 游戏标识 |
| owner_id | varchar | 账号 ID（uuid） |
| title | varchar | 模板标题 |
| description | text | 描述 |
| size | int | 团队规模 |
| lineup_list_json | json | 子活动列表 `[{dayOfWeek, time, name}]` |
| player_limit_json | json | 人员限制 |
| server_region_json | json | 服务器/大区/阵营 |
| visibility_json | json | 可见性设置 `{isPublic, notShare}` |
| created_at/updated_at/deleted_at | datetime | |

索引建议：
- `(game_key, owner_id)`

### 2.8 activity_logs（活动日志）

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint | 主键 |
| game_key | varchar | 游戏标识（冗余，便于查询） |
| activity_id | bigint | 外键 |
| action | varchar | 操作类型：`signup/status_change/lineup_save/activity_update/assistant_add/assistant_remove` |
| actor_id | varchar | 操作者 ID |
| target_json | json | 操作对象（角色名/阵容名/变更字段等） |
| detail | varchar | 操作描述 |
| created_at | datetime | |

索引建议：
- `(activity_id, created_at DESC)`

### 2.9 auto_lineup_templates（AI 排班模板）

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint | 主键 |
| game_key | varchar | 游戏标识 |
| user_id | varchar | 用户 ID |
| size | int | 团队规模（10/25/40） |
| grid_json | json | 网格约束 `[{index, constraint: {baseType, jobKey?, specKey?}}]` |
| created_at/updated_at | datetime | |

约束：
- `(game_key, user_id, size)` 唯一

## 3. 访问控制（Web 特有）

### 3.1 activity_invites（邀请链接 token）

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint | 主键 |
| activity_id | bigint | 外键 |
| created_by | varchar | owner/assistant |
| name | varchar | 备注（可选） |
| token_hash | varchar | token hash（不存明文） |
| expires_at | datetime | 过期时间 |
| revoked_at | datetime | 撤销时间 |
| created_at | datetime | |

### 3.2 activity_access_grants（可选：访问授权记录）

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint | 主键 |
| activity_id | bigint | 外键 |
| user_id | varchar | 授权给谁（登录用户） |
| method | enum | `code/invite` |
| granted_at | datetime | 授权时间 |
| expires_at | datetime | 授权过期（可选） |

## 4. 社交与协作

建议沿用现有概念，但命名更通用：

- `account_follows`：关注
- `account_blacklists`：黑名单
- `activity_assistants`（或沿用 milestone_assistants）：活动管理员关系

## 5. 通知系统（Outbox + Delivery）

### 5.1 notifications（通知事件）

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint | 主键 |
| user_id | varchar | 收件人 |
| channel | enum | `email`（MVP） |
| event_type | enum | `deadline_reminder/assistant_invite/status_changed/activity_updated` |
| locale | varchar | `en/zh` |
| timezone | varchar | IANA（America/Los_Angeles） |
| payload_json | json | 渲染模板所需数据（活动标题/时间/链接等） |
| dedupe_key | varchar | 幂等键（避免重复发） |
| scheduled_at | datetime | 计划发送时间（截止提醒） |
| created_at | datetime | |

### 5.2 notification_deliveries（投递记录）

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint | 主键 |
| notification_id | bigint | 外键 |
| provider | varchar | `ses/sendgrid/smtp` |
| status | enum | `pending/sent/failed` |
| attempt | int | 重试次数 |
| last_error | text | 错误信息 |
| sent_at | datetime | |

### 5.3 user_notification_settings（用户通知设置）

| 字段 | 类型 | 说明 |
|---|---|---|
| user_id | varchar | 主键 |
| email_enabled | bool | 总开关 |
| events_json | json | 每类事件开关 |
| digest_window_minutes | int | 合并窗口 |
| rate_limit_json | json | 限速策略 |
| locale | varchar | 默认语言 |
| timezone | varchar | 默认时区 |
| updated_at | datetime | |

## 6. 多游戏扩展字段的约束（建议）

- Core 只读写 `rules_json/payload_json`，解释与校验交给 Adapter。
- API 层可以做“结构校验（JSON schema）”，具体语义校验交给 Adapter。

