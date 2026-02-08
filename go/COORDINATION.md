# 协调和合并总体计划

## 项目概览
- **项目**：OpenCode Golang 重构
- **主分支**：`dev`（默认）
- **工作分支**：
  - `golang-refactor` (Agent 1) - TUI & 后端
  - `golang-refactor-v2` (Agent 2) - 架构、性能、文档
- **目标完成日期**：2026-02-15

---

## 分支结构和责任

### 主分支 `dev`
- 生产就绪的代码
- 所有工作基于此分支从 origin fork

### 工作分支

#### `golang-refactor` (Agent 1)
**职责**：核心功能和 TUI 开发

**已完成**：
- Phase 1: 脚手架、后端服务、API、基础 TUI（5 commits）
  - 文件管理（FileService）
  - 项目管理（ProjectService）
  - HTTP API 接口（8 个端点）
  - TUI 框架（文件树、编辑器、状态栏）
  - 4 个单元测试（全部通过）
  - 编译成功（8.3MB + 11MB 二进制）

**进行中**：
- Phase 2: TUI 完整实现（预定 Feb 10-12）
  - 语法高亮
  - 富文本编辑
  - 命令面板
  - 快捷键系统
  - 文件监听
  - TUI-Backend 通信

**计划**：
- Phase 3: 后端扩展（Feb 12-14）
  - 搜索和索引
  - Git 集成
  - 缓存和优化
- Phase 4: CI/CD（Feb 14-15）
  - 跨平台编译
  - 性能基准
  - 发布流程

**GitHub**：[golang-refactor 分支](https://github.com/oliveagle/opencode/tree/golang-refactor)

---

#### `golang-refactor-v2` (Agent 2)
**职责**：架构验证、性能优化、文档、GUI 规划

**即将开始**：
- Phase 1: 架构评审（Feb 8-10）
  - API 接口评审
  - 核心包评审
  - 类型系统评估
  - TUI 架构评审
  - 输出：详细评审报告

- Phase 2: 性能分析（Feb 10-12）
  - 基准测试实现
  - 内存分析
  - 优化建议
  - 输出：性能报告 + 优化方案

- Phase 3: 配置和文档（Feb 12-14）
  - XDG 标准支持
  - 配置格式设计
  - 跨平台兼容性
  - 日志策略
  - 输出：配置规范 + 文档库

- Phase 4: GUI 基础（Feb 14-15，取决于 TUI 完成）
  - 框架技术选型（Fyne vs Gio）
  - 设计系统
  - 架构原型
  - 输出：技术选型报告 + GUI 原型

**GitHub**：[golang-refactor-v2 分支](https://github.com/oliveagle/opencode/tree/golang-refactor-v2)

---

## 协作协议

### 代码共享和所有权

```
go/
├── pkg/
│   ├── types/        ← 共享（Agent 1 最终决定）
│   ├── core/         ← 共享（Agent 1 最终决定）
│   ├── api/          ← 共享（Agent 1 最终决定）
│   ├── tui/          ← Agent 1 独占
│   └── storage/      ← Agent 1 独占
├── cmd/
│   ├── tui/          ← Agent 1 独占
│   ├── backend/      ← Agent 1 独占
│   └── gui/          ← Agent 2（后续）
├── test/
│   ├── core/         ← Agent 1 独占
│   └── benchmarks/   ← Agent 2 独占
├── doc/              ← Agent 2 独占
├── config/           ← Agent 2 独占
└── COLLABORATION.md  ← 共享维护
```

### 冲突解决规则
1. **共享代码冲突**：Agent 1 优先决定（TUI 优先级高）
2. **文档冲突**：Agent 2 最终决定
3. **架构决策**：在线讨论，不能达成共识时 Agent 1 决定
4. **性能优化**：合并前征求 Agent 2 意见

### PR 和代码审查
```
Agent 1 流程：
  1. 在 golang-refactor 分支开发
  2. 完成 phase 后在 GitHub 创建 PR 到 dev
  3. 等审核后合并到 dev

Agent 2 流程：
  1. 在 golang-refactor-v2 分支开发
  2. 完成 phase 后在 GitHub 创建 PR 到 golang-refactor
  3. Agent 1 审核并决定是否合并
  4. 如需合并到 dev，转给 Agent 1 处理
```

### 日常同步
```bash
# Agent 1 同步 Agent 2 的建议
git fetch origin golang-refactor-v2
git log origin/golang-refactor-v2 --oneline

# Agent 2 同步 Agent 1 的最新代码
git fetch origin golang-refactor
git rebase origin/golang-refactor
```

---

## 合并时间表和策略

| 日期 | Agent | Phase | 输出 | 目标分支 |
|------|-------|-------|------|---------|
| Feb 10 | Agent 1 | Phase 2 开始 | TUI Alpha | - |
| Feb 10 | Agent 2 | Phase 1 完成 | 架构报告 | golang-refactor |
| Feb 12 | Agent 1 | Phase 2 完成 | TUI Beta | - |
| Feb 12 | Agent 2 | Phase 2 完成 | 性能报告 | golang-refactor |
| Feb 13 | Agent 2 | Phase 3 完成 | 配置规范 | golang-refactor |
| Feb 14 | Agent 1 | Phase 3 完成 | 后端扩展 | - |
| Feb 14 | Agent 2 | Phase 4 完成 | GUI 设计 | - (待 TUI 完成) |
| Feb 15 | Agent 1 | 合并到 dev | v1.0 Ready | **dev** |
| Feb 15 | Agent 2 | 文档归档 | 完整文档库 | **dev** |

### 最终合并流程（Feb 15）
```
golang-refactor:
  git checkout dev
  git pull origin dev
  git merge --no-ff golang-refactor -m "Merge: golang-refactor (v1.0)"
  git push origin dev

golang-refactor-v2:
  git checkout dev  
  git pull origin dev
  git merge --no-ff golang-refactor-v2 -m "Merge: golang-refactor-v2 (docs & design)"
  git push origin dev

Tag release:
  git tag -a v1.0-golang -m "OpenCode Golang Refactor v1.0"
  git push origin v1.0-golang
```

---

## 关键里程碑检查清单

### Agent 1 Checklist
- [ ] Phase 2: TUI 功能完善（Feb 12）
  - [ ] 语法高亮或文本装饰
  - [ ] 撤销/重做功能
  - [ ] 搜索和替换
  - [ ] 分割窗口
  - [ ] File watcher
  - [ ] 快捷键配置
  - [ ] ≥80% 代码覆盖率测试

- [ ] Phase 3: 后端扩展（Feb 14）
  - [ ] 文件搜索
  - [ ] Git 集成（基础）
  - [ ] 缓存系统
  - [ ] 并发处理增强

- [ ] Phase 4: CI/CD（Feb 15）
  - [ ] GitHub Action 工作流
  - [ ] Linux/macOS/Windows 编译
  - [ ] ARM64 支持
  - [ ] 发布脚本

- [ ] 最终检查
  - [ ] 所有测试通过
  - [ ] 代码覆盖率 ≥80%
  - [ ] README.md 更新
  - [ ] 版本号 bump

### Agent 2 Checklist
- [ ] Phase 1: 架构评审（Feb 10）
  - [ ] API 设计评估
  - [ ] 核心库审查
  - [ ] 扩展性分析
  - [ ] 输出 3+ 评审文档

- [ ] Phase 2: 性能分析（Feb 12）
  - [ ] ≥10 个基准测试
  - [ ] 内存使用报告
  - [ ] 优化建议清单
  - [ ] 瓶颈识别

- [ ] Phase 3: 配置和文档（Feb 14）
  - [ ] XDG 实现规划
  - [ ] 配置格式设计
  - [ ] 跨平台检查
  - [ ] ≥20 页文档

- [ ] Phase 4: GUI 基础（Feb 15）
  - [ ] 框架对比报告
  - [ ] 设计系统定义
  - [ ] 可行性原型
  - [ ] 技术选型建议

- [ ] 最终检查
  - [ ] 所有设计文档完整
  - [ ] API 文档生成
  - [ ] 迁移指南（JS → Go）
  - [ ] 版本历史记录

---

## 通信和问题跟踪

### GitHub Issues
使用 label：
- `agent1` - Agent 1 负责
- `agent2` - Agent 2 负责
- `golang-refactor` - TUI 相关
- `golang-refactor-v2` - 文档/性能相关
- `bug` - 问题追踪
- `architecture` - 架构决策

### 周期同步
- **每日**：git push，保持分支最新
- **每周**：PR review 和合并
- **推荐**：在 GitHub PR 上讨论架构决策

### 快速参考
```bash
# Agent 1 开发流程
git checkout golang-refactor
git pull origin golang-refactor
# ... 编辑代码
go test -v ./...
git commit -m "feat: ..."
git push origin golang-refactor

# Agent 2 开发流程
git checkout golang-refactor-v2
git fetch origin golang-refactor  # 查看 Agent 1 进展
# ... 编辑文档/性能分析
git commit -m "docs: ... \n\nCo-authored-by: Agent 1"
git push origin golang-refactor-v2
```

---

## 预期产物

### Agent 1 产物
```
go/
├── cmd/
│   ├── tui/          (功能完整的 TUI)
│   └── backend/      (扩展的后端服务)
├── pkg/
│   ├── core/         (≥95% 测试覆盖)
│   ├── api/          (完整 API)
│   └── tui/          (完整 UI 组件)
├── test/             (集成测试)
└── Makefile          (构建脚本)
```

### Agent 2 产物
```
go/
├── doc/
│   ├── architecture-review.md
│   ├── performance-analysis.md
│   ├── optimization-roadmap.md
│   ├── xdg-support.md
│   ├── gui-framework-comparison.md
│   └── design-system.md
├── config/
│   └── config.example.yaml
├── test/
│   └── benchmarks/
└── AGENT2_PLAN.md
```

### Merged to `dev`
- 完整的 Golang OpenCode v1.0
- 完整的文档库
- 性能基准数据
- GUI 技术选型报告

---

## 版本历史

| 版本 | 日期 | 内容 |
|------|------|------|
| v0.1 | 2026-02-08 | 脚手架 + 基础服务 |
| v0.2 | 2026-02-10 | TUI Alpha + 架构评审 |
| v0.3 | 2026-02-12 | TUI Beta + 性能分析 |
| v0.4 | 2026-02-14 | TUI Complete + 后端扩展 |
| **v1.0** | **2026-02-15** | **完整发布** |

---

**最后更新**：2026-02-08  
**维护人**：Agent 1（作为协调者）
