# go-qemu (现代化与多版本支持 Fork)

`go-qemu` 是一个基于 Go 语言的 QEMU 管理与自动化控制框架。本仓库在原 go-qemu 传输层基础上进行了彻底重建，支持现代 QEMU 版本（重点覆盖 **QEMU 5.0.0** 与 **QEMU 10.2.2**），提供强大的 QAPI Schema 代码生成、Feature Gates 特性门控、版本感知命令行构建、三层 Dry-Run 校验以及跨版本热迁移编排能力。

---

## 核心特性

1. **传输层现代化 (`qmp/`)**
   - Context 感知的命令发送与超时控制 (`ContextMonitor`)。
   - 结构化错误分类 (`CommandError`, `TransportError`, `ErrorClass`)。
   - Socket 自动等待与重试机制 (`WaitForSocket`)。

2. **多文件 QAPI Schema 解析器与代码生成器 (`qapi/`, `cmd/qapi-gen/`)**
   - 支持 QEMU 5.0+ 多文件嵌套包含 (`include`)、条件编译 (`if`) 与 `features` 声明解析。
   - 为 QEMU 5.0.0 与 QEMU 10.2.2 自动生成强类型 Go 绑定 (`internal/gen/v5/`, `internal/gen/v10/`) 及静态 `CommandSet` 元数据。

3. **版本适配层与 Feature Gates (`internal/version/`, `internal/feature/`)**
   - 精确的版本区间匹配与解析 (`QEMUVersion`, `VersionRange`)。
   - `FeatureGateRegistry` 实现 **版本 + 架构 + Schema 运行时** 三重校验，自动识别版本功能差异（例如 `postcopy-preempt`、`disk.blockdev-preferred`）。

4. **版本感知命令行构建器 (`internal/vmconfig/`, `internal/cmdline/`)**
   - VMConfig 通用配置中间表示，支持 YAML/HCI 与 Proxmox PVE 格式加载。
   - 模块化 `BuildRule` 规则引擎，自动切换 `-drive` (5.0) 与 `-blockdev` (10.2.2)，机器类型版本跟随 (`pc-i440fx-5.0` vs `pc-i440fx-10.2`)。
   - `BuildReport` 包含风险等级提示与 `CmdlineDiff` 跨版本参数对比报告。

5. **三层 Dry-Run 校验体系 (`internal/dryrun/`, `internal/process/`)**
   - **Layer 1 静态校验**：参数完整性、内存对齐与冲突检测（毫秒级，无副作用）。
   - **Layer 2 进程冒烟**：启动冻结态 (`-S`) 临时 QEMU 进程校验参数有效性。
   - **Layer 3 运行时语义**：通过 QMP 校验运行时 Schema 与命令语义。

6. **跨版本热迁移与升级编排 (`internal/migrate/`, `internal/upgrade/`)**
   - `MigrationNegotiator` 自动协商两端迁移能力，自动禁用目标端独有能力（例如 5.0 到 10.2.2 迁移时自动禁用 `postcopy-preempt`）。
   - `LocalUpgrader` 同机热升级编排，提供完整状态机与回滚机制。

7. **Harness SOP 完整验收套件 (`harness/`, `cmd/qvm/`)**
   - 提供标准验收场景（同机热升级、跨版本协商、Cmdline Diff、PVE 格式兼容、Feature Matrix）。
   - 提供 CLI 工具 `qvm` 与 `harness` 命令行运行器。

---

## 目录结构

```
go-qemu/
├── qmp/               # 传输层（Context 支持、Socket 重试、结构化错误）
├── qapi/              # QAPI Schema 解析器与代码生成器引擎
├── cmd/
│   ├── qapi-gen/      # Schema 代码生成器 CLI
│   └── qvm/           # QEMU VM 命令行与工具 CLI
├── internal/
│   ├── version/       # 版本解析与区间运算
│   ├── feature/       # Feature Gates 特性注册表与三重校验
│   ├── vmconfig/      # VM 配置 IR 与 YAML/PVE 加载器
│   ├── cmdline/       # 命令行构建器、BuildReport 与 CmdlineDiff
│   │   └── rules/     # 模块化构建规则 (-drive, -blockdev, -machine 等)
│   ├── process/       # QEMU 进程管理与 socket 监听
│   ├── dryrun/        # 三层 Dry-Run 校验体系
│   ├── migrate/       # 跨版本迁移能力协商器与监控器
│   ├── upgrade/       # 同机热升级编排器与回滚处理
│   └── gen/           # 自动生成的 Go 绑定 (v5 / v10)
├── harness/           # Harness SOP 验收测试套件与矩阵 runner
├── schemas/           # QEMU 5.0.0 与 10.2.2 QAPI Schema 定义文件
├── _archive/          # 存档的旧版本代码 (QEMU 2.11 legacy)
├── Makefile           # 构建与测试入口
└── README.md
```

---

## 快速开始

### 1. 代码生成

根据 `schemas/` 下的 JSON 文件生成 Go 类型定义与命令集合：

```bash
make gen
```

### 2. 构建 CLI 工具

构建 `qapi-gen` 与 `qvm` 命令行工具：

```bash
make build
```

### 3. 运行单元测试与集成测试

运行全量测试套件：

```bash
make test
```

### 4. 运行 Harness SOP 验收矩阵

执行完整 Harness SOP 场景验证：

```bash
make harness
```

---

## 命令行工具 Usage (`qvm`)

### 构建命令行参数

```bash
./bin/qvm build-cmdline --config harness/env/fixtures/vm_basic_x86.yaml --version 10.2.2
```

### 执行 Dry-Run 三层校验

```bash
./bin/qvm dry-run --config harness/env/fixtures/vm_basic_x86.yaml
```

### 执行同机热升级 (5.0 -> 10.2.2)

```bash
./bin/qvm upgrade --config harness/env/fixtures/vm_basic_x86.yaml --vmid 100
```

---

## 许可证

基于 Apache 2.0 许可证开源。
