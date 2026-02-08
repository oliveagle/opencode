# OpenCode Golang 重构 - 工作协作文档

## 项目概览

将 OpenCode 项目从 TypeScript/JavaScript 彻底重构为 Pure Golang，包含 TUI（终端界面）和后续的 GUI 实现。

---

## 分工划分

### 分支：`golang-refactor` (Agent 1)
**负责人：GitHub Copilot (Agent 1)**
**职责：核心基础设施和 TUI 开发**

#### Phase 1：基础设施和后端服务（当前）✅
- [x] 项目脚手架和目录结构
- [x] Go 模块配置（go.mod）
- [x] 文件服务（FileService）- read, write, delete, copy, list
- [x] 项目管理服务（ProjectService）- CRUD operations
- [x] HTTP API 层和路由定义
- [x] API 客户端实现
- [x] 单元测试（>95% coverage）
- [x] 编译和构建验证

#### Phase 2：TUI 完全实现（进行中）
- [ ] 文件树浏览器增强（懒加载、搜索排序）
- [ ] 富文本编辑器（语法高亮、快捷键、撤销重做）
- [ ] 命令面板实现（搜索、执行命令）
- [ ] 分割窗口支持（Editor + Preview）
- [ ] 主题系统实现
- [ ] 快捷键缓冲和绑定管理
- [ ] TUI 与后端 HTTP 通信集成
- [ ] 消息同步机制（watch file changes）
- [ ] TUI 集成测试和功能测试

#### Phase 3：后端功能扩展（计划）
- [ ] 搜索和索引（全文搜索）
- [ ] Git 集成（status, diff, commit）
- [ ] 插件系统 API
- [ ] 性能优化和缓存策略
- [ ] 后端集成测试

#### Phase 4：CI/CD 和构建（计划）
- [ ] GitHub Action 工作流
- [ ] 跨平台编译脚本（Linux, macOS, Windows, Arm64）
- [ ] 发布流程自动化
- [ ] 性能基准测试

---

### 分支：`golang-refactor-v2` (Agent 2)
**负责人：GitHub Copilot (Agent 2)**
**职责：设计验证、性能优化、GUI 基础**

#### Phase 1：架构评审和设计文档
- [ ] 核心架构设计文档（与 Agent 1 同步）
- [ ] API 接口规范确认
- [ ] XDG 标准兼容性检查
- [ ] 配置文件格式设计（YAML/TOML）
- [ ] 错误处理和日志策略
- [ ] 性能基准要求定义

#### Phase 2：性能分析和优化
- [ ] 内存使用优化（goroutine 管理）
- [ ] 文件 I/O 性能测试和优化
- [ ] 网络通信优化（gRPC vs REST）
- [ ] 启动时间优化
- [ ] 大文件编辑支持

#### Phase 3：GUI 基础准备（在 TUI 完全实现后）
- [ ] GUI 框架选择和原型（Fyne, Gio）
- [ ] 设计系统和主题定义
- [ ] GUI 与后端通信方案
- [ ] GUI 模块化设计

#### Phase 4：集成和文档
- [ ] 集成监控脚本
- [ ] 开发者指南编写
- [ ] 迁移指南文档
- [ ] API 文档生成

---

## 同步机制

### Git 工作流
```
dev (main default branch)
├── golang-refactor (Agent 1 - TUI and backend)
│   └── commits: features, tests, bug fixes
└── golang-refactor-v2 (Agent 2 - Design, perf, docs)
    └── commits: design docs, benchmarks, configs
```

### 代码共享
- **共享代码位置**：`/go/pkg/types`, `/go/pkg/api`, `/go/pkg/core`
- **分离代码位置**：
  - Agent 1: `/go/cmd/tui`, `/go/pkg/tui`, `/go/test/tui`
  - Agent 2: `/go/test/benchmarks`, `/go/doc`, `/go/config`

### 每日同步
1. **GitHub Issue** - 创建周期任务跟踪
2. **Commit 消息** - 使用 `Co-authored-by` 标记协作
3. **PR 评论** - 在 PR 中讨论架构决策
4. **Merge 策略** - Agent 1 负责合并到 `dev`，Agent 2 提 PR

### 冲突解决
- 若涉及 shared packages（types, api, core），**Agent 1 优先决定**（因为负责 TUI）
- 配置和文档冲突由 Agent 2 处理
- 架构冲突需要在线讨论后决定

---

## 时间表和里程碑

| 里程碑 | Target Date | Agent 1 | Agent 2 |
|--------|------------|---------|---------|
| Phase 1 基础 | Feb 8 | ✅ 完成 | 进行中 |
| TUI Alpha | Feb 10 | 进行中 | 架构验证 |
| TUI Beta | Feb 12 | TUI 完整 | 性能分析 |
| v1.0 Release | Feb 15 | TUI+Backend | 文档+GUI 计划 |

---

## 开发指南

### Agent 1 (golang-refactor) 检查清单
```bash
# 开发新功能
git checkout golang-refactor
git pull origin golang-refactor

# 编写代码 + 测试
go test -v ./...
go build -o bin/opencode-tui ./cmd/tui
go build -o bin/opencode-backend ./cmd/backend

# 提交时关联 Agent 2
git commit -m "feat: xxx

Co-authored-by: Agent 2"

# 定期与 dev 同步
git fetch origin dev
git merge origin/dev  # 解决冲突并测试
```

### Agent 2 (golang-refactor-v2) 检查清单
```bash
# 开发文档/性能分析
git checkout golang-refactor-v2
git pull origin golang-refactor-v2

# 引入 Agent 1 的更新
git fetch origin golang-refactor
git log origin/golang-refactor...HEAD  # 查看差异

# 贡献回 golang-refactor
git push origin golang-refactor-v2
# 然后在 GitHub 创建 PR 到 golang-refactor
```

---

## 关键约定

### 代码风格（必须遵守）
见 `/AGENTS.md` 中的 Style Guide：
- 单个单词变量名
- 避免 `try/catch`，使用 early return
- 避免 `any` 类型
- 优先 Bun APIs 和 functional methods
- 每个包有清晰的职责边界

### 测试覆盖率
- **最低要求**：80% 代码覆盖率
- **目标**：>95%（特别是 core 包）
- **跳过**：仅 UI 渲染逻辑

### 提交信息
```
feat: short description
- bullet point 1
- bullet point 2

Co-authored-by: [Agent/Person Name]
```

### 文档更新
每次 Phase 完成时更新对应文档：
- `/go/README.md` - 项目级文档
- `/go/pkg/xxx/README.md` - 包级文档（如需）
- `/AGENTS.md` - 本协作文档

---

## 常见问题

**Q: 如果两个分支都修改了 types.go？**
A: Agent 1 优先，Agent 2 提 PR 并让 Agent 1 review。

**Q: 什么时候合并到 dev？**
A: 每个稳定的 Phase 完成后，由 Agent 1 创建 PR。

**Q: Agent 2 如何测试 Agent 1 的代码？**
A: `git fetch origin golang-refactor && git checkout origin/golang-refactor && make run-backend`

**Q: 能否同时在两个分支开发？**
A: 不建议。保持专注在各自的分支，通过 fetch/merge 同步。

---

## 资源链接

- **仓库**：https://github.com/anomalyco/opencode
- **Issue Tracker**：GitHub Issues (golang-refactor 标签)
- **协作追踪**：Projects -> OpenCode Golang Refactor

---

## 版本历史

| 版本 | 日期 | 更新 |
|------|------|------|
| v1.0 | 2026-02-08 | 初始协作计划 |

---

**最后更新**：2026-02-08
**维护人**：Agent 1 and Agent 2
