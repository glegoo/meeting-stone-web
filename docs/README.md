# Meeting Stone Web 文档索引

本目录包含 **产品设计（PRD）**、**技术方案**、**API 设计**、**数据模型**、**通知系统**、**迁移计划**、**新增游戏适配指南** 以及 **UX/UI 设计指南**。

## 文档列表

| 文档 | 说明 | 主要读者 |
|------|------|----------|
| [01_PRD.md](01_PRD.md) | 产品需求文档（多游戏、访问控制、权限、通知、时区/i18n） | PM、全员 |
| [02_ARCHITECTURE.md](02_ARCHITECTURE.md) | 总体架构（Core + Game Adapter + Game Manifest + 前端组件） | 后端、前端、架构 |
| [03_API.md](03_API.md) | API 设计（活动/模板/报名/阵容/日志/导出） | 前后端 |
| [04_DATA_MODEL.md](04_DATA_MODEL.md) | 数据模型/表结构草案（Core 域 + 扩展字段 + Outbox + 日志） | 后端、DBA |
| [05_NOTIFICATIONS.md](05_NOTIFICATIONS.md) | 通知系统设计（MVP=Email） | 后端 |
| [06_MIGRATION_PLAN.md](06_MIGRATION_PLAN.md) | 迁移与里程碑（表结构映射、任务分解、灰度/回滚） | 后端、运维 |
| [07_ADD_NEW_GAME.md](07_ADD_NEW_GAME.md) | 新增游戏指南（Adapter/Manifest/ScoreProvider/团队Buff） | 后端 |
| [08_UX_UI_GUIDE.md](08_UX_UI_GUIDE.md) | UX/UI 设计指南（页面结构、交互流程、组件规范） | **UI 设计师**、前端 |

## 全局约定（重要）

- **多游戏策略**：每个活动都属于一个 `gameKey`（如 `wow-classic`、`ff14`），所有核心数据（活动/模板/报名/阵容/报名组）均以 `gameKey` 作为隔离维度。
- **核心架构**：平台 Core 不写死任何游戏规则，游戏差异通过 **Game Adapter** 与 **Game Manifest** 注入。
- **访问控制**：Web 无法限制转发，`not_share` 语义在 Web 侧升级为"访问码/邀请链接"等可撤销的访问控制。
- **通知**：MVP 仅做 Email；事件覆盖截止提醒/管理员邀请/状态变更/活动变更；支持 i18n 与时区。

## 快速导航

### 产品层面
- 了解产品功能范围 → [01_PRD.md](01_PRD.md)
- 了解页面设计规范 → [08_UX_UI_GUIDE.md](08_UX_UI_GUIDE.md)

### 技术层面
- 了解系统架构 → [02_ARCHITECTURE.md](02_ARCHITECTURE.md)
- 了解 API 接口 → [03_API.md](03_API.md)
- 了解数据模型 → [04_DATA_MODEL.md](04_DATA_MODEL.md)

### 实施层面
- 了解通知系统 → [05_NOTIFICATIONS.md](05_NOTIFICATIONS.md)
- 了解迁移计划 → [06_MIGRATION_PLAN.md](06_MIGRATION_PLAN.md)
- 新增游戏支持 → [07_ADD_NEW_GAME.md](07_ADD_NEW_GAME.md)

## 文档更新日志

| 日期 | 版本 | 更新内容 |
|------|------|----------|
| 2026-01-13 | v1.1 | 完善所有文档：增加模板功能、子活动规则、阵容编排细节、团队Buff配置、表结构映射、里程碑任务分解、WoW适配器示例等 |
| 2026-01-13 | v1.0 | 初始版本 |
