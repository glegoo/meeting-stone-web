# 05 通知系统设计（MVP = Email）

## 1. 目标与原则

### 1.1 目标

- Web MVP 先提供 **Email 通知**，覆盖你确认的核心事件：
  - 截止提醒
  - 管理员邀请/接受/解除
  - 报名状态变更（请假/取消/被移除）
  - 活动变更（标题/描述/时间/公开性/访问控制变化）
- 设计可扩展到 Discord/Web Push（后续），且不侵入核心业务链路。

### 1.2 原则

- **主链路不直接发送**：业务只写 Outbox/Notification 记录，异步 Worker 负责投递。
- **可控与可追踪**：每次投递有状态、有重试、有失败原因。
- **幂等与去重**：避免重复发送（尤其是活动频繁更新时）。
- **尊重用户设置**：事件级开关 + 频控 + 退订。
- **多语言与时区**：邮件模板 i18n（en/zh），时间按用户时区渲染。

## 2. MVP 渠道

- Email（唯一必做）

后续：
- Discord Webhook / Discord DM（OAuth 绑定）
- Web Push（PWA）

## 3. MVP 事件定义

### 3.1 deadline_reminder（截止提醒）

触发：
- 定时任务扫描：`deadline_at - now` 命中提醒窗口（例如：24h、3h、1h、30m）

接收者：
- 报名状态 = signed 且活动未过期的用户
- 可选：Owner（可配置）

频控：
- 每个活动每个用户每个窗口最多 1 封

### 3.2 assistant_invite（管理员邀请/接受/解除）

触发：
- Owner 发起邀请
- Assistant 接受
- Owner/Assistant 解除关系

接收者：
- 邀请：被邀请人
- 接受/解除：双方

### 3.3 status_changed（报名状态变更）

触发：
- 用户自行请假/取消
- Owner/Assistant 移除报名

接收者：
- 当事人（必发）
- 若被移除：可选通知 Owner（可配置）

### 3.4 activity_updated（活动变更）

触发字段（MVP 建议只对“重要变更”发）：
- `deadline_at`、`title`、`description`、`is_public`、`access_mode`、访问码更新

接收者：
- 已报名用户 + Assistant
- 可按变更类型细分：仅时间/截止类变更才触发（减少噪音）

去抖/合并：
- 同一活动在 `digest_window_minutes` 内多次变更 → 合并成摘要邮件（只发一次）

## 4. 用户设置与退订

### 4.1 设置项（建议）

- `email_enabled`：总开关
- `events.deadline_reminder`
- `events.assistant_invite`
- `events.status_changed`
- `events.activity_updated`
- `digest_window_minutes`：合并窗口（默认 5 分钟）
- `timezone`：默认用户时区
- `locale`：默认语言
- 频控策略：每小时/每天上限

### 4.2 退订链接（Email 必备）

每封邮件都包含：
- 一键退订（关闭全部邮件）
- 通知设置入口（关闭某类事件）

实现：
- 邮件中的 URL 带签名 token（短有效期），可免登录完成退订动作（安全要求：token 一次性或短期有效）。

## 5. 投递架构（Outbox + Worker）

### 5.1 数据流

1. 业务事件发生（例如活动更新）→ 写入 `notifications`（或 `outbox`）
2. Worker 拉取 `pending` 记录
3. 渲染模板（根据 locale/timezone + payload）
4. 发送 Email（Provider）
5. 写入 `notification_deliveries`，更新状态

### 5.2 幂等与去重

- `dedupe_key` 建议结构：
  - `deadline:{activityId}:{userId}:{window}`
  - `activity_updated:{activityId}:{userId}:{digestWindowStart}`
  - `assistant_invite:{activityId}:{inviteId}:{userId}`
  - `status_changed:{detailId}:{status}:{updatedAtBucket}`

### 5.3 重试与死信

- `attempt < maxAttempts` 时指数退避重试（例如 1m/5m/30m）
- 超过重试次数 → `failed`，进入“死信”可人工排查

### 5.4 限流

- Provider 维度限流（QPS）
- 用户维度限流（hour/day）

## 6. 模板与 i18n

### 6.1 模板组织

- `templates/email/{locale}/{event}.mjml` 或 `html`
- 统一布局（header/footer）
- 事件模板只负责内容区域

### 6.2 payload 规范（建议）

所有事件 payload 至少包含：
- `activityId`
- `activityTitle`
- `activityUrl`（注意：若活动需要访问控制，URL 可能指向“输入访问码页”或携带 invite token）
- `deadlineAtUtc`
- `userTimezone`
- `gameKey`

## 7. 安全与隐私

- 不在邮件中暴露敏感联系方式（除非是 owner/assistant 且用户勾选公开给团长；MVP 建议邮件不包含联系方式）
- 访问控制：不要在邮件中包含明文访问码；邀请链接 token 需可撤销且可过期

## 8. 邮件模板示例

### 8.1 截止提醒模板

```html
<!-- templates/email/zh/deadline_reminder.mjml -->
<mjml>
  <mj-body>
    <mj-section>
      <mj-column>
        <mj-text font-size="20px" color="#333">
          ⏰ 活动即将截止
        </mj-text>
        <mj-text>
          您报名的活动 <strong>{{activityTitle}}</strong> 将于 
          <strong>{{deadlineFormatted}}</strong> 截止。
        </mj-text>
        <mj-button href="{{activityUrl}}" background-color="#46c3fc">
          查看活动详情
        </mj-button>
      </mj-column>
    </mj-section>
    <mj-section>
      <mj-column>
        <mj-text font-size="12px" color="#999">
          不想收到此类邮件？<a href="{{unsubscribeUrl}}">退订</a> 或 
          <a href="{{settingsUrl}}">管理通知设置</a>
        </mj-text>
      </mj-column>
    </mj-section>
  </mj-body>
</mjml>
```

### 8.2 模板变量规范

| 变量 | 类型 | 说明 |
|------|------|------|
| `activityTitle` | string | 活动标题 |
| `activityUrl` | string | 活动详情链接 |
| `deadlineFormatted` | string | 截止时间（按用户时区格式化） |
| `ownerName` | string | 团长昵称 |
| `characterName` | string | 角色名（状态变更通知） |
| `changesSummary` | string | 变更摘要（活动变更通知） |
| `unsubscribeUrl` | string | 一键退订链接 |
| `settingsUrl` | string | 通知设置链接 |

## 9. 定时任务配置

### 9.1 截止提醒任务

```go
// 建议使用 cron 表达式或定时器
// 每 5 分钟检查一次即将截止的活动

func (s *NotificationService) ScanDeadlineReminders() {
    windows := []time.Duration{
        24 * time.Hour,  // 提前 24 小时
        3 * time.Hour,   // 提前 3 小时
        1 * time.Hour,   // 提前 1 小时
    }
    
    for _, window := range windows {
        activities := s.repo.FindActivitiesDeadlineWithin(window, window-5*time.Minute)
        for _, activity := range activities {
            s.CreateDeadlineReminderNotifications(activity, window)
        }
    }
}
```

### 9.2 通知合并任务

```go
// 每分钟检查一次待合并的通知

func (s *NotificationService) ProcessDigestWindow() {
    // 查找 digest_window 已过期的待发送通知
    pendingNotifications := s.repo.FindPendingDigestNotifications()
    
    // 按 (user_id, activity_id, event_type) 分组
    groups := groupNotifications(pendingNotifications)
    
    for _, group := range groups {
        if len(group) > 1 {
            // 合并为摘要
            s.MergeToDigest(group)
        } else {
            // 直接发送
            s.MarkReadyToSend(group[0])
        }
    }
}
```

## 10. MVP 验收清单

- [ ] 能配置用户通知设置（开关、时区、语言、合并窗口）
- [ ] 4 类事件都能生成通知记录并成功发送 Email
- [ ] 具备去抖/合并机制（activity_updated）
- [ ] 具备重试与投递日志（delivery）
- [ ] 每封邮件都含退订入口
- [ ] 截止提醒定时任务正常运行（24h/3h/1h）
- [ ] 邮件模板支持中英文
- [ ] 时间按用户时区正确显示
- [ ] 退订链接有效（一键退订/管理设置）

