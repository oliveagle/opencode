# OpenCode - Go Refactor

Golang 重构版本，包含 TUI（终端界面）和后端服务。

## 项目结构

```
go/
├── cmd/
│   ├── tui/           # TUI 应用入口
│   └── backend/       # 后端服务入口
├── pkg/
│   ├── core/          # 核心服务逻辑
│   ├── api/           # API 类型定义
│   ├── storage/       # 存储层（数据库、缓存）
│   └── tui/           # TUI 组件库
├── test/              # 测试用例
├── go.mod             # Go 模块定义
├── Makefile           # 构建脚本
└── README.md          # 本文件
```

## 依赖

- Go 1.23+
- tview (TUI framework)
- tcell (terminal handling)
- cobra (CLI)
- viper (config management)

## 构建

```bash
# 下载依赖
make mod-download

# 构建所有
make build

# 仅构建 TUI
make build-tui

# 仅构建后端
make build-backend
```

## 运行

```bash
# 运行 TUI
make run-tui

# 运行后端
make run-backend
```

## 开发阶段

当前状态：**脚手架初始化完成**

### Phase 1: TUI 基础（当前）
- [ ] 文件浏览器
- [ ] 文本编辑器
- [ ] 命令面板
- [ ] 快捷键处理
- [ ] 主题支持

### Phase 2: 后端服务
- [ ] 文件 I/O 服务
- [ ] 项目管理
- [ ] API 层

### Phase 3: GUI（TUI 完成后）
- [ ] 使用 Fyne 或 Gio 实现

## 测试

```bash
make test
```

## 清理

```bash
make clean
```
