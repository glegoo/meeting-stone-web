# 07 新增游戏指南（Adapter / Manifest / ScoreProvider）

> 目标：让你未来接入 FF14 或其他组队游戏时，不需要重写 Core；新增游戏主要是“新增配置 + 实现适配器/评分源”。

## 1. 新增游戏需要准备什么？

至少需要定义三件事：

1) **Game Manifest**（前端渲染契约）
2) **Game Adapter**（规则/校验/阵容能力）
3) **ScoreProvider**（评分来源与归一化）

建议你把每个游戏的实现放在：

- `api/internal/games/{gameKey}/...`
- `web/src/games/{gameKey}/...`（如果少量前端差异需要定制；优先用 manifest 驱动）

## 1.1 gameKey 的命名与粒度（重要）

`gameKey` 的本质是 **规则集（ruleset）**。只要职业/服务器/评分口径/规模等任一关键规则发生显著变化，就应该使用不同的 `gameKey`，以便：

- 使用不同 Adapter/Manifest
- 数据自然隔离（活动/报名/阵容互不影响）

以魔兽世界为例，建议使用不同 `gameKey`：

- `wow-retail`
- `wow-classic`（或更细：`wow-era` / `wow-sod` / `wow-cata` …）

同时增加一个聚合字段（用于 UI/统计，不参与业务隔离）：

- `gameFamily` / `parentGameKey`：统一为 `wow`

## 1.2 同一 gameKey 内的小版本差异：用 Manifest 版本化/phase 管理

当差异属于“同一规则集内的小幅演进”（例如赛季阶段增加职业/专长展示、字典增量），优先用：

- `manifestVersion`（递增）
- 或 `phase` / `effectiveFrom`（配置化）

只有当差异大到需要不同 Adapter（规则解释都不同）时，才拆新的 `gameKey`。

## 2. Game Manifest 设计要点

Manifest 的目标是“让前端不用写死规则”。

### 2.1 必备字段（建议）

- `gameKey`：唯一标识
- `roles`：职责体系（例如 tank/heal/dps 或 melee/ranged）
- `jobs/specs`：职业与专长（Job/Class + Spec）
- `partySizes`：允许规模（例如 8/24/40）
- `features`：
  - `serverRegion`：是否有服务器/大区
  - `faction`：是否有阵营
  - `buffComposition`：是否启用团队增益展示
  - `scoreProvider`：评分源类型（wcl/fflogs/custom）
- `rosterShape`：阵容布局（grid/list，列数等）
- `constraints`：slot 约束支持哪些类型（any/role/job/spec/empty）

### 2.2 版本化

建议加入 `manifestVersion`，便于前端缓存与兼容：

- `manifestVersion: 1`

## 3. Game Adapter 接口（建议能力）

### 3.1 规则/校验

- `ValidateActivityCreate(input)`：活动创建校验
  - 是否允许 server/region/faction
  - size 是否允许
  - player limits 是否符合该游戏规则
- `ValidateSignup(input)`：报名校验
  - 角色名规则（不同游戏不同）
  - job/spec 是否存在
  - payload_json 是否符合 schema

### 3.2 阵容与 AI 排班能力

- `NormalizeConstraints(gridConstraints)`：把前端的 slot 约束规范化
- `AutoLineup(activity, lineups, details, constraints)`：输出推荐排班（不保存）
- `CheckUserDuplicationPolicy(activity)`：同账号重复入队策略（阻止/提示/允许）

### 3.3 可选：团队增益

- `ComputeBuffComposition(lineupDetails)`：输出 buff 覆盖情况（适用于 WoW/部分游戏）

## 4. ScoreProvider（评分源）设计要点

不同游戏评分差异极大，建议统一输出：

- `score_value`：0..100（或 0..1000），并定义全平台一致的区间
- `score_tier`：S/A/B/C（可选）
- `score_source`：wcl/fflogs/custom
- `score_meta`：原始数据（百分位、难度、赛季等）

### 4.1 WoW（WCL）

现有实现参考：
- `wow-classic-recruit/pkg/services/mileStoneService.go` 中 `rpglogs.GetCharacterListWclInfo(...)`

### 4.2 FF14（FFLogs）

建议接入方式：
- 以 FFLogs character/encounter API 拉取最近表现或 best percentile
- 将 percentile 映射到 `score_value`（例如 percentile=95 → score_value=95）

注意：
- 需要处理不同难度/版本/职业的对比口径

## 5. 新增游戏的实施步骤（Checklist）

1. **确定 gameKey**：例如 `ff14`
2. **补齐字典**：
   - roles/jobs/specs/icons/colors
   - partySizes/rosterShape
3. **实现 Manifest Endpoint**：
   - `GET /api/v1/games/{gameKey}/manifest`
4. **实现 Adapter**：
   - 活动校验/报名校验/限制解释
   - AI 排班（如果该游戏需要）
5. **接入 ScoreProvider**：
   - FFLogs 或自定义分数
6. **前端适配**：
   - 使用 manifest 驱动表单与选择器
   - 检查特殊字段：是否有服务器/阵营、是否有多队形
7. **回归测试**：
   - 创建活动/搜索/报名/阵容/AI/通知（至少走通一次）

## 6. 团队Buff（Raid Composition）配置

对于需要团队Buff展示的游戏（如 WoW CTM/MOP/正式服），需要在 Manifest 中配置：

### 6.1 Buff 分类定义

```json
{
  "buffComposition": {
    "enabled": true,
    "categories": [
      {
        "id": "attack_power",
        "name": "攻击强度",
        "compact": "攻强",
        "icon": "ability_warrior_battleshout",
        "providedBy": [
          { "job": "Warrior", "spec": "*" },
          { "job": "DeathKnight", "spec": "Blood" },
          { "job": "Hunter", "spec": "*", "note": "通过宠物" }
        ]
      },
      {
        "id": "spell_power",
        "name": "法术强度",
        "compact": "法强",
        "icon": "spell_arcane_arcanepotency",
        "providedBy": [
          { "job": "Mage", "spec": "Arcane" },
          { "job": "Shaman", "spec": "Elemental" }
        ]
      },
      {
        "id": "crit",
        "name": "暴击",
        "compact": "暴击",
        "icon": "spell_nature_lightning",
        "providedBy": [
          { "job": "Mage", "spec": "*" },
          { "job": "Druid", "spec": "Balance" }
        ]
      },
      {
        "id": "haste",
        "name": "急速",
        "compact": "急速",
        "icon": "spell_nature_bloodlust",
        "providedBy": [
          { "job": "Shaman", "spec": "*" },
          { "job": "Mage", "spec": "*" }
        ]
      }
    ]
  }
}
```

### 6.2 前端计算逻辑

```typescript
function computeBuffComposition(lineup: LineupSlot[], manifest: GameManifest) {
  const categories = manifest.buffComposition.categories;
  
  return categories.map(category => {
    const activeSpecs = lineup
      .filter(slot => slot.detail)
      .map(slot => ({ job: slot.detail.jobKey, spec: slot.detail.specKey }));
    
    const isActive = category.providedBy.some(provider => 
      activeSpecs.some(active => 
        active.job === provider.job && 
        (provider.spec === '*' || active.spec === provider.spec)
      )
    );
    
    return {
      ...category,
      active: isActive
    };
  });
}
```

## 7. WoW Adapter 示例结构

以下是 WoW Classic 适配器的示例代码结构：

### 7.1 目录结构

```
api/internal/games/
├── registry.go           # 游戏注册表
├── interface.go          # Adapter 接口定义
└── wow-classic/
    ├── adapter.go        # 核心适配器
    ├── manifest.go       # Manifest 数据
    ├── validator.go      # 校验逻辑
    ├── lineup.go         # 阵容/AI 排班逻辑
    ├── score.go          # 评分 Provider
    └── data/
        ├── jobs.json     # 职业定义
        ├── specs.json    # 专长定义
        ├── servers.json  # 服务器列表
        └── buffs.json    # 团队Buff定义
```

### 7.2 Adapter 接口

```go
type GameAdapter interface {
    // 基本信息
    GetGameKey() string
    GetManifest() *Manifest
    
    // 校验
    ValidateActivityCreate(input *CreateActivityInput) error
    ValidateSignup(input *SignupInput) error
    
    // 阵容
    AutoLineup(activity *Activity, details []*SignupDetail, constraints []SlotConstraint) ([][]int64, error)
    CheckDuplicationPolicy() DuplicationPolicy
    
    // 评分
    GetScoreProvider() ScoreProvider
}
```

### 7.3 Manifest 实现

```go
func (a *WowClassicAdapter) GetManifest() *Manifest {
    return &Manifest{
        GameKey:         "wow-classic",
        GameFamily:      "wow",
        ManifestVersion: 1,
        Roles: []Role{
            {Key: "tank", Label: "坦克", Icon: "Tank.png"},
            {Key: "healer", Label: "治疗", Icon: "Healer.png"},
            {Key: "dps", Label: "输出", Icon: "Melee.png"},
        },
        Jobs: a.loadJobs(),      // 从 data/jobs.json 加载
        Specs: a.loadSpecs(),    // 从 data/specs.json 加载
        Features: Features{
            ServerRegion:    true,
            Faction:         true,
            BuffComposition: true,  // CTM/MOP 开启
            ScoreProvider:   "wcl",
        },
        PartySizes:  []int{10, 25, 40},
        RosterShape: RosterShape{Type: "grid", Cols: 5},
        BuffComposition: a.loadBuffs(),
    }
}
```

## 8. 常见坑

- **时区**：FF14/WoW 用户分布广，活动时间展示必须按用户时区，且邮件中明确标注
- **角色名规则**：不同游戏字符集/长度/空格规则不同，校验必须放到 adapter
- **评分口径**：不同赛季/版本/难度口径不同，必须在 score_meta 里保留原始上下文
- **访问控制**：Web 链接可扩散，建议默认支持访问码/邀请链接（尤其是非公开活动）
- **团队Buff版本差异**：WoW 不同版本的 Buff 系统差异巨大（经典旧世 vs CTM vs 正式服），需要按 gameKey 独立配置
- **服务器/阵营合并**：WoW 服务器合并/连接频繁，服务器列表需要定期更新
- **双天赋/多天赋**：部分游戏/版本支持双天赋或多天赋，报名时需要处理主副天赋逻辑

## 9. 版本差异处理示例

### 9.1 WoW 不同版本的差异

| 特性 | 经典旧世 | CTM/MOP | 正式服 |
|------|----------|---------|--------|
| 团队规模 | 20/40 | 10/25 | 10/20/25/40 |
| 双天赋 | 无 | 有 | 有 |
| 团队Buff | 简化 | 完整 | 完整 |
| 服务器类型 | 普通/PvP | 普通/PvP | 已合并 |
| 评分来源 | WCL Classic | WCL | WCL |

### 9.2 gameKey 建议

```
wow-era        # 经典旧世 60 级
wow-sod        # 发现赛季
wow-hc         # 硬核
wow-cata       # CTM 怀旧服
wow-mop        # MOP 怀旧服
wow-retail     # 正式服（当前版本）
```

每个 gameKey 对应独立的 Manifest 和 Adapter 配置。

