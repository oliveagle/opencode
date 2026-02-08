# Agent 2 工作计划 - golang-refactor-v2

## 分支信息
- **分支名**：`golang-refactor-v2`
- **基准**：`dev` (包含 Agent 1 的所有工作)
- **职责**：架构验证、性能优化、GUI 规划、文档

---

## Phase 1：架构评审（当前）

### 目标
验证 Agent 1 的架构设计，确保可扩展性和性能。

### 任务清单
- [ ] **1.1** API 接口评审（pkg/api/）
  - 评估 REST vs gRPC 的性能差异
  - 创建基准测试计划
  - 文档：`/go/doc/api-design.md`

- [ ] **1.2** 核心包评审（pkg/core/）
  - FileService: 检查大文件处理、并发模型
  - ProjectService: 配置持久化策略
  - 文档：`/go/doc/core-design.md`

- [ ] **1.3** 类型系统评审（pkg/types/）
  - 确认数据结构定义完整性
  - 性能指标：序列化/反序列化开销
  - 文档：`/go/doc/types-spec.md`

- [ ] **1.4** TUI 架构评审（pkg/tui/）
  - 组件化设计评估
  - 状态管理模式
  - 文档：`/go/doc/tui-architecture.md`

### 输出物
- [ ] `go/doc/architecture-review.md` - 完整设计评审报告
- [ ] `go/doc/performance-targets.md` - 性能指标定义
- [ ] `go/doc/scalability-plan.md` - 可扩展性规划

---

## Phase 2：性能分析（Feb 10-12）

### 目标
建立性能基准，识别瓶颈，制定优化方案。

### 任务清单
- [ ] **2.1** 基准测试实现
  - 文件 I/O 性能测试
  - API 通信开销
  - TUI 渲染性能
  - 位置：`go/test/benchmarks/`

- [ ] **2.2** 内存管理分析
  - Goroutine 泄漏检测
  - 堆内存排查
  - 工具：pprof, graphviz
  - 文档：`go/doc/memory-analysis.md`

- [ ] **2.3** 架构优化建议
  - 缓存策略（文件、索引）
  - 并发模型改进
  - 通信协议评估（gRPC vs REST vs IPC）
  - 文档：`go/doc/optimization-roadmap.md`

### 输出物
- [ ] 基准测试报告（CSV + 图表）
- [ ] 优化建议清单（优先级排序）
- [ ] 内存池、缓存设计方案

---

## Phase 3：配置和文档（Feb 12-14）

### 目标
为项目标准化配置、XDG 支持、跨平台兼容性。

### 任务清单
- [ ] **3.1** XDG 标准实现规划
  - XDG_CONFIG_HOME: `~/.config/opencode/`
  - XDG_DATA_HOME: `~/.local/share/opencode/`
  - XDG_CACHE_HOME: `~/.cache/opencode/`
  - 文档：`go/doc/xdg-support.md`

- [ ] **3.2** 配置文件格式设计
  - 主配置：YAML or TOML？
  - 主题系统设计
  - 快捷键配置格式
  - 文档：`go/config/schema.md`

- [ ] **3.3** 跨平台支持检查
  - 路径处理（Windows vs Unix）
  - 终端兼容性
  - 构建脚本检查
  - 文档：`go/doc/cross-platform.md`

- [ ] **3.4** 日志和监控策略
  - 日志级别定义
  - 结构化日志设计
  - 性能指标导出
  - 文档：`go/doc/logging-strategy.md`

### 输出物
- [ ] `go/config/config.example.yaml`
- [ ] XDG 兼容性 checklist
- [ ] 跨平台测试计划

---

## Phase 4：GUI 基础设计（Feb 14-15）

### 目标
为后续 GUI 实现打下基础（在 TUI 完成后启动）。

### 任务清单
- [ ] **4.1** GUI 框架技术选型
  - Fyne: 轻量级、易学，但功能有限
  - Gio: 高性能、现代 UI，但学习陡峭
  - 对比分析和推荐
  - 文档：`go/doc/gui-framework-comparison.md`

- [ ] **4.2** GUI 设计系统
  - 色彩方案（深色/浅色）
  - 字体和布局规范
  - 组件库设计
  - 文档：`go/doc/design-system.md`

- [ ] **4.3** 架构设计（GUI）
  - 与后端通信方案（共用 HTTP API？）
  - 状态同步机制
  - 多窗口支持
  - 文档：`go/doc/gui-architecture.md`

- [ ] **4.4** 原型实现
  - 创建最小化功能原型
  - 验证框架可行性
  - 性能评估
  - 代码：`go/cmd/gui-prototype/`

### 输出物
- [ ] 技术选型报告
- [ ] 设计系统 Figma/原型
- [ ] GUI 可行性研究

---

## 工作指南

### 日常工作流
```bash
# 切换到 golang-refactor-v2
git checkout golang-refactor-v2

# 定期检查 Agent 1 的最新进展
git fetch origin golang-refactor
git log origin/golang-refactor...HEAD

# 如果需要 merge Agent 1 的改动
git rebase origin/golang-refactor

# 编写文档和分析
echo "analysis" > go/doc/your-analysis.md
git add go/doc/your-analysis.md
git commit -m "docs: your analysis

Co-authored-by: Agent 1"
```

### 文档位置
- 架构设计：`/go/doc/architecture-*.md`
- 性能分析：`/go/doc/performance-*.md`
- 技术方案：`/go/doc/technical-*.md`

### 代码位置
- 基准测试：`/go/test/benchmarks/`
- 配置示例：`/go/config/`
- GUI 原型：`/go/cmd/gui-prototype/`（后续）

### 与 Agent 1 交互
1. **提交 PR**：完成 phase 后创建 PR 到 golang-refactor
2. **Code Review**：Agent 1 review 后 merge
3. **实时讨论**：关键决策在 PR comments 中讨论

---

## 关键里程碑

| 日期 | 内容 | 输出 |
|------|------|------|
| Feb 8 | Phase 1 启动 | 任务分解 |
| Feb 10 | Phase 1 完成 | 架构评审报告 |
| Feb 12 | Phase 2 完成 | 性能基准 + 优化建议 |
| Feb 13 | Phase 3 完成 | XDG + 配置规划 |
| Feb 14 | Phase 4 完成 | GUI 技术选型 |
| Feb 15 | 总结报告 | 完整文档库 |

---

## 关键指标

- **目标覆盖率**：>95% 代码库评审
- **基准数量**：≥10 个核心性能指标
- **文档页数**：≥20 页设计文档
- **优化建议**：≥5 个主要优化方向

---

## 备注

- 这是初始计划，可根据进度调整
- Agent 1 优先级更高，Agent 2 协力支持
- 每周五同步一次进度
- 所有决策默认 Agent 1 有最终发言权
