# 基于 docs/architecture.md 的完整开发与落地计划

## 前置分析与设计思想锚定

### 1. 开源整合策略决策

#### 架构核心哲学：以“单点 Fork + 异源参考”打破耦合陷阱
> **一句话原则**：只承载一个传输实现（go-qemu），避免将多个项目的类型系统与底层逻辑硬性碰撞或合并，其他项目仅作为行为基准、测试夹具或生成思路的参考源。

#### 模块分工与能力参考矩阵

为避免“类型系统硬焊接”导致的依赖膨胀与维护灾难，采用**单主干、多参考**的解耦策略：

```
                    ┌─────────────────────────┐
                    │   核心产品 (Production) │
                    └────────────┬────────────┘
                                │
                    [唯一的物理 Fork / 传输层]
                                │
                ┌─────────────┴─────────────┐
                │    go-qemu / qmp 包       │
                └─────────────┬─────────────┘
                                │
    ┌──────────────────────────┼──────────────────────────┐
    ▼                          ▼                          ▼
┌──────────────┐           ┌──────────────┐           ┌──────────────────────────────────────┐
│  qemu.qmp    │           │   qapi-rs    │           │ Proxmox VTP                          │
│ (Python 官方)│           │ (Rust 项目)  │           │ (qemu-server, using qemu 5.0.0)      │
└──────┬───────┘           └──────┬───────┘           └──────┬───────────────────────────────┘
       │                          │                          │
[行为对照标尺]            [生成思路参照]            [夹具与规则生成，参照 Proxmox PVE qemu-server]
```

1. **传输基座（唯一的物理 Fork）：`go-qemu/qmp`**
   - **选用理由**：Go 语言生态中经过实际生产验证、质量最高的 QMP 通信与握手实现。
   - **边界约束**：仅保留其 Socket 通信、JSON 解析与最基础的 QMP 握手协议，剥离上层业务逻辑。

2. **行为与协议标尺：`qemu.qmp` (Python 官方)**
   - **角色定义**：**黑盒行为对照**（绝对不链接/嵌入到生产包中）。
   - **借鉴点**：参照官方客户端的问候顺序（Greeting）、Capabilities 协商逻辑以及 Error Object 的标准异常结构。
   - **交付形式**：作为集成测试与 CI 的行为基准标尺。

3. **类型绑定生成思路：`qapi-rs` (Rust)**
   - **角色定义**：**代码生成（Codegen）设计模式参照**。
   - **借鉴点**：参考其“模式进（Schema-in）、索引出（Index-out）”的类型绑定范式。
   - **解耦精髓**：**轻量生成（Index-Only Codegen）** — 本仓库的 Generator 只产出字段索引和少量元数据映射（Meta-index），绝不完整复制或硬克隆对方整套臃肿的类型系统。

4. **驱动编排与测试夹具：Proxmox `qemu-server`**
   - **角色定义**：**规则配置引擎与测试夹具参照**。
   - **借鉴点**：借鉴 `config_to_command` 的链式/规则驱动构建思想，以及 `cfg2cmd` 的测试夹具结构。
   - **隔离原则**：测试数据与 Config 规则映射统一存放在 `harness/env/fixtures` 下，严格拦截在生产包外。

#### 关键链路执行规范

1. **严格的通信握手与防御式降级（Handshake Pipeline）**
   - 在启动任何虚拟机编排（Orchestration）之前，必须建立硬性的握手检验隔离带：
   ```
   [连接 Socket] ──> [QMP 握手 (Go-QMP)] ──> [基于 QEMU 5.0.0 与目标二进制校验]
                                                       │
                                       ┌──────────────┴──────────────┐
                                       ▼                             ▼
                                   [ 成功 ]                      [ 失败 ]
                                       │                             │
                               [ 允许进入编排阶段 ]           [ 立即中断！拒绝启动编排 ]
                                                                   │
                                                           [ 优先修补传输层适配 ]
   ```
   - **校验机制**：合入任何传输代码前，必须基于 QEMU 5.0.0 及目标 QEMU 二进制完成双向握手测试。
   - **熔断规则**：若握手或 Capabilities 协商失败，**禁止进入编排层**；研发优先级强制收敛至“修改传输层”，确保基座完好后再暴露上层 API。

2. **“轻索引”代码生成策略（Index-Only Codegen）**
   - **输入**：QEMU 标准 QAPI JSON Schema。
   - **输出**：轻量级 Go 结构体 + 运行时字段 Index（`schema_meta_gen.go` 中的 `CommandSet` 等）。
   - **优势**：既保障了 Go 语言调用时的类型安全（Type Safety），又不会因 QEMU 版本升级导致大面积类型断言崩溃。

3. **测试与环境夹具隔离（Test Fixtures Policy）**
   - 所有的 Proxmox `cfg2cmd` 测试案例、QEMU 参数构建期望（Expected Command Lines）均收容在 `harness/env/fixtures`。
   - 生产代码仅包含通用渲染逻辑，不允许出现任何特定硬件配置或硬编码的 Vendor 规则。

---

### 2. 第一原则：围绕 QEMU 命令行参数的精准构建
根据 `docs/architecture.md` 的要求，整套系统的第一原则是：**在任何 QEMU 版本组合下，都能精准构建出正确的命令行参数，并且能在问题发生之前（而非之后）发现参数错误或版本不兼容**。
`VMConfig` 中间表示（IR）与具体的 QEMU 版本完全解耦，命令行构建器根据目标 QEMU 版本和场景提示（`SceneHints`），通过版本感知的规则引擎生成最终的命令行参数。

### 3. 约束清单与硬性边界
- **C01（硬约束）**：无 libvirt，所有操作直接通过 QMP Domain Socket / TCP 与 QEMU 进程直连。
- **C02（硬约束）**：支持 QEMU 5.0.0 到 10.2.2 的直接跨版本热迁移，虚拟机业务零中断。
- **C03（硬约束）**：跨版本热迁移能力必须协商（`MigrationNegotiator`），不能假设双端能力相同。
- **C04（功能约束）**：三层 Dry-Run 校验（静态校验、进程冒烟、QMP 运行时语义）缺一不可。
- **C05（功能约束）**：`cfg2cmd` 与历史 Proxmox PVE Perl 实现的等效性必须可验证（通过 `cfg2cmd-equivalents.yaml`）。
- **C06（功能约束）**：HCI 企业定制补丁与原生 QEMU 能力分层管理（通过 `capability.yaml` 元数据）。
- **C07（功能约束）**：多架构支持（x86_64、aarch64、c86_64 兆芯）。
- **C08（功能约束）**：场景分支引擎（VDI、AMD SEV/Intel TDX 机密计算、等保合规、CXL/NUMA 内存分层）。
- **C09（设计约束）**：AI-friendly 架构（语义清晰、配置驱动、Schema 驱动）。
- **C10（设计约束）**：引入新 QEMU 版本的边际成本极低（仅需添加 Schema 和运行 `make gen`）。

---

## 阶段规划与深度落地任务

```
阶段一  Fork 整理、传输层现代化与基础设施构建   （传输基座与 Handshake Pipeline）
阶段二  Schema 解析器与代码生成器（qapi-gen）  （Index-Only Codegen 双版本覆盖）
阶段三  版本适配层与企业补丁管理                （Feature Gates & Capability Meta）
阶段四  命令行构建器、场景引擎与 cfg2cmd 兼容  （第一原则与 Proxmox Fixtures）
阶段五  Dry-Run 三层校验体系                    （全三层预防网）
阶段六  跨版本热迁移编排、能力协商与回滚机制    （5.0→10.2.2 无缝升级）
阶段七  Harness SOP 完整验收与版本矩阵          （七大场景 SOP 验收）
```

---

### 阶段一：Fork 整理、传输层现代化与基础设施构建

#### 目标
整理 `go-qemu` 仓库，存档旧版 2.11 代码；补强传输层 `qmp/`，支持 Context 取消、连接就绪等待、Handshake 熔断与结构化错误分类；建立版本类型 `QEMUVersion` 和 `VersionRange`。

#### 具体任务

1.1 **仓库清理与架构隔离**
- 保持 `qmp/` 核心传输层代码，清理/存档旧 `qemu/` 包（移入 `_archive/qemu-2.11-legacy/`）。
- 升级 `go.mod` 至 Go 1.22+。

1.2 **传输层补强（qmp/ 包扩展）**
- `qmp/context.go`：提供 Context 感知的 QMP 命令发送和超时处理。
- `qmp/ready.go`：实现 `WaitForSocket(ctx, socketPath)` 轮询机制，等待 QEMU 启动后 QMP socket 就绪。
- `qmp/errors.go`：定义结构化错误类型 `CommandError`（对照 `qemu.qmp` 解析 QMP ErrorClass）与 `TransportError`。

1.3 **版本管理核心（internal/version）**
- `internal/version/version.go`：定义 `QEMUVersion`，包含 `Major`, `Minor`, `Micro`, `HCITag`（如 `"hci-6.8.0"`）。支持版本号解析（`Parse`）与版本比较（`GTE`, `LT`）。
- `internal/version/range.go`：定义 `VersionRange` 闭开区间 `[From, To)`，实现 `Contains(v QEMUVersion)` 方法。

#### 验收标准
- `go test ./qmp/... ./internal/version/...` 全部通过。
- 能够通过 `qmp` 包连接到运行中的 QEMU 5.0.0 或 10.2.2 进程，成功完成 Handshake、发送 `query-version` 并接收事件。

---

### 阶段二：Schema 解析器与代码生成器（qapi-gen）

#### 目标
重建 QAPI Schema 解析器与代码生成器 `qapi-gen`（借鉴 `qapi-rs` 的 Index-Only Codegen 策略），支持 QEMU 5.0+ 多文件 JSON Schema（对照 Python `qemu.qmp` 标准）、`features` 字段以及条件编译指令 (`if: condition`)，为 QEMU 5.0.0 和 10.2.2 生成强类型 Go 绑定与索引映射。

#### 新增与强化目录
```
qapi/                           -- Schema 解析与代码生成
├── schema.go                   -- QAPISchema 结构体与多文件 ParseDir/ParseFile 递归加载
├── schema_types.go             -- SchemaEntry, CommandDef, StructDef, EnumDef, UnionDef
├── schema_features.go          -- features 字段提取与匹配
├── schema_conditions.go        -- 条件编译表达式 ConditionExpr 解析
├── gen.go                      -- Generator 主体
├── gen_types.go                -- 生成 types_gen.go
├── gen_commands.go             -- 生成 commands_gen.go
└── gen_meta.go                 -- 生成 schema_meta_gen.go (包含 CommandSet 静态命令集合)

cmd/qapi-gen/main.go            -- qapi-gen CLI 工具入口
schemas/
├── qemu-5.0.0/qapi/            -- 5.0.0 Schema 提取文件
└── qemu-10.2.2/qapi/           -- 10.2.2 Schema 提取文件
internal/gen/
├── v5/                         -- 生成代码（gen_v5）
└── v10/                        -- 生成代码（gen_v10）
```

#### 核心实现要点
1. **递归处理 include**：`parseWithIncludes` 递归解析 `qapi-schema.json` 中的 `include` 指令（如 `block.json`, `migration.json`, `cxl.json` 等），防止循环包含。
2. **Schema 差异识别**：能够识别 5.0.0 与 10.2.2 之间的 Schema 差异（例如 10.2.2 中新增的 `postcopy-preempt` 枚举值和 `cxl` 命令）。
3. **Index-Only 静态命令集合生成**：生成 `gen_v5.CommandSet` 与 `gen_v10.CommandSet` (`map[string]struct{}`)，用于无实例时的静态 Feature Gate 校验。

#### 验收标准
- 执行 `make gen` 无报错，在 `internal/gen/v5` 和 `internal/gen/v10` 中正确产出编译通过的 Go 代码。
- 单元测试验证 Schema 解析器能够识别出 10.2.2 中特有的字段和命令。

---

### 阶段三：版本适配层与企业补丁管理（Feature Gates & Capability Meta）

#### 目标
实现 `FeatureGateRegistry`，通过**三重验证**（版本区间 + 架构兼容 + Schema 运行时/静态确认）查询特性可用性；引入 `capability.yaml` 元数据机制，实现 HCI 企业定制补丁与社区原生 QEMU 能力的分离管理。

#### 新增结构与目录
```
internal/capability/
├── capability.go               -- CapabilityMeta 数据结构 (BackportedPatches, PrivatePatches, DisabledFeatures)
└── loader.go                   -- capability.yaml 加载解析与 Schema 交叉验证

internal/feature/
├── gate.go                     -- FeatureGate 定义与 Reason 结构体
├── registry.go                 -- FeatureGateRegistry 管理类（三重验证）
├── gates_core.go               -- 核心 QMP 命令门控
├── gates_migrate.go            -- 迁移门控（postcopy-ram, postcopy-preempt, multifd 等）
├── gates_device.go             -- 设备与磁盘门控 (blockdev-preferred, explicit-backend)
├── gates_arch_x86.go           -- x86_64 架构门控
├── gates_arch_arm.go           -- aarch64 架构门控
└── gates_arch_c86.go           -- c86_64 兆芯架构门控
```

#### 三重验证逻辑
1. **版本区间检查**：`gate.Available.Contains(r.currentVer)`。
2. **架构兼容检查**：`r.currentArch` 在 `gate.Arches` 许可范围内。
3. **Schema 确认 / 企业补丁校验**：
   - 若连接运行实例，用 `liveSchema.HasKey(gate.SchemaKey)` 精确查找。
   - 若加载了 `CapabilityMeta`，检查该特性是否被 `disabled_upstream_features` 禁用，或在 `private_patches` / `backported_patches` 中被显式激活。
   - 若无运行实例，用 `staticCommandSet` 兜底检查。

#### 验收标准
- 在 QEMU 5.0.0 环境下，`migrate.postcopy-preempt` 门控查询返回 `false`。
- 在 QEMU 10.2.2 环境下，`migrate.postcopy-preempt` 门控查询返回 `true`。
- 若企业 `capability.yaml` 声明禁用了某特性，即使版本号支持，`IsEnabled` 依然准确返回 `false`。

---

### 阶段四：命令行构建器、场景引擎与 cfg2cmd 兼容

#### 目标
实现 `VMConfig` 中间表示（IR）与声明式规则引擎 `BuildRule`；借鉴 Proxmox `qemu-server` 的 `config_to_command` 构建模式；引入 `SceneHints` 场景分支引擎（VDI/SEV/TDX/等保合规/内存分层）；产出包含完整决策链的 `BuildReport` 和 `CmdlineDiff` 报告；实现 Proxmox `cfg2cmd` 对比验证器与等效映射规则 (`cfg2cmd-equivalents.yaml`)，隔离测试夹具在 `harness/env/fixtures`。

#### 核心模块
```
internal/vmconfig/
├── config.go                   -- VMConfig IR (CPU, Memory, Machine, Disks, Networks, Display, VMX)
├── loader_hci.go               -- HCI YAML 格式配置加载
├── loader_pve.go               -- Proxmox PVE 历史格式配置加载
└── validator.go                -- 配置语义校验

internal/cmdline/
├── builder.go                  -- CmdlineBuilder 构建器入口
├── option.go                   -- WithQEMUVersion, WithArch, WithScene 等选项
├── report.go                   -- BuildReport, ArgEntry, RiskLevel (RiskNone, RiskDeprecated, RiskBehaviorChange)
├── diff.go                     -- CmdlineDiff 版本间差异计算与格式化输出
├── conflict.go                 -- 参数冲突检测器
└── rules/                      -- 构建规则库
    ├── rule.go                 -- BuildRule 接口
    ├── registry.go             -- DefaultRegistry
    ├── cpu.go                  -- CPU & SMP 规则
    ├── memory.go               -- Memory & memory-backend-ram 规则
    ├── machine.go              -- Machine type (跟随 QEMU Major.Minor，如 pc-i440fx-5.0 / 10.2)
    ├── disk_selector.go        -- 根据 blockdev-preferred 选择 -drive 或 -blockdev
    ├── network.go              -- -netdev 与 -device virtio-net-pci
    └── qmp_socket.go           -- -chardev 与 -mon 参数
```

#### 关键实现细节
1. **磁盘参数切换**：
   - QEMU 5.0.0：生成旧式 `-drive file=...,format=raw,if=none,id=drive-scsi0`。
   - QEMU 10.2.2：生成新式 `-blockdev driver=raw,node-name=drive-scsi0,...` 及配套 `-device scsi-hd`。
2. **机器类型自动跟随**：
   - 5.0.0 → `-machine pc-i440fx-5.0,accel=kvm`。
   - 10.2.2 → `-machine pc-i440fx-10.2,accel=kvm`。
3. **`CmdlineDiff` 风险分级**：
   - 标注 `RiskBehaviorChange`（例如机器类型大版本变化）。
   - 标注 `RiskDeprecated`（例如使用了废弃的 `-drive` 写法）。
4. **`cfg2cmd` 等效校验与夹具隔离**：
   - 对比 Proxmox PVE Perl `qm` 输出与 Go 实现输出，解析 `cfg2cmd-equivalents.yaml` 规则进行匹配，夹具统一置于 `harness/env/fixtures`。

#### 验收标准
- 相同的 `VMConfig`，在 5.0.0 下生成 `-drive`，在 10.2.2 下生成 `-blockdev`。
- `DiffFrom` 生成的报告能准确识别并打印新增、删除和变更的命令行参数。

---

### 阶段五：Dry-Run 三层校验体系

#### 目标
实现三层 Dry-Run 校验体系，在任何修改应用到虚拟机之前提前捕获语法、进程和运行时错误。

```
Layer 1: Static Validation  (毫秒级，无副作用) -- 参数格式/必填项/内存4MB对齐/冲突检测
Layer 2: Process Smoke     (秒级，启动 -S 进程) -- 测试参数启动/QMP Socket 就绪/关键设备存在性
Layer 3: Runtime Semantic  (运行时，QMP 交互)  -- Schema 检查/迁移能力交集计算/热插拔可行性
```

#### 核心代码实现（internal/dryrun）
- `runner.go`：管理 `Run(ctx, cfg)` 入口与 LayerMask 位掩码（`Layer1Static`, `Layer2Smoke`, `Layer3Runtime`）。
- `layer1_static.go`：静态检查必填参数（`-machine`, `-m`, `-smp`）、内存对齐、`-chardev`/`-mon` 配置、冲突参数。
- `layer2_smoke.go`：管理带 `-S` 参数的冻结态临时 QEMU 进程，测试 QMP socket 能否正常 Dial 并响应 `query-version` 和 `query-qmp-schema`。
- `layer3_runtime.go`：在运行实例上校验 QMP 命令字段完整性与迁移能力支持。

#### 验收标准
- 合法配置通过全部三层校验。
- 缺少 `-machine` 在 Layer 1 立即报错，不触发 Layer 2。
- 错误的 QEMU 二进制路径在 Layer 2 准确报错。

---

### 阶段六：跨版本热迁移编排、能力协商与回滚机制

#### 目标
实现 `MigrationNegotiator`（跨版本能力协商），参照 `qemu.qmp` 规范自动计算 5.0 源端与 10.2.2 目标端的能力交集，强制禁用目标端独有能力（如 `postcopy-preempt`）；实现包含明确回滚策略的热升级编排器（`LocalUpgrader`）与跨机迁移器。

#### 核心架构
```go
// MigrationNegotiator 核心算法：
// 1. 查询源端 (5.0) 迁移能力集合
// 2. 查询目标端 (10.2.2) 迁移能力集合
// 3. 取两端已知能力的交集（以源端状态为基准）
// 4. 将目标端独有能力（如 postcopy-preempt）显式设置为 State=false 强制禁用
// 5. 应用协商后的能力列表到目标端
```

#### 迁移状态机与回滚策略
```
IDLE -> PREPARING -> NEGOTIATING -> MIGRATING -> COMPLETING -> DONE
               \             \            \             \
                -------> FAILED -------> ROLLING_BACK -> ROLLBACK_DONE
```
- 在 `COMPLETING`（switchover）之前，源端 VM 保持完整可用状态，任何异常触发目标端进程终止与源端继续运行。

#### 验收标准
- 在 Mock 及真实测试中，从 5.0 迁移至 10.2.2 时，`postcopy-preempt` 被自动添加到禁用列表并应用于目标端。
- 升级过程中断或失败时，回滚机制能正确保持源端 VM 稳定运行。

---

### 阶段七：Harness SOP 完整验收与版本矩阵

#### 目标
实现 `docs/architecture.md` 第六章定义的 Harness SOP **七大验收场景**，配置自动化版本矩阵测试，生成结构化的 JSON/HTML 验收报告。

#### 场景实现清单（harness/scenarios/）
1. **Scenario 1: `local-upgrade-5.0-to-10.2.2`**（同机 QEMU 热升级，业务零中断，停机时间 < 500ms）。
2. **Scenario 2: `cross-host-migrate`**（跨宿主机 5.0 到 10.2.2 热迁移）。
3. **Scenario 3: `cross-version-migration-negotiation`**（专项验证迁移能力协商与交集计算）。
4. **Scenario 4: `cmdline-diff-5.0-vs-10.2.2`**（Cmdline Full Dump 与废弃/行为变更 Warning 报告）。
5. **Scenario 5: `patch-verify`**（热补丁更新前后 Schema Diff 与特性门控稳定性验证）。
6. **Scenario 6: `cfg2cmd-compat-pve`**（Proxmox 历史 Perl `qm` cfg2cmd 输出等效性对比）。
7. **Scenario 7: `api-stability`**（上层业务 API 在 QEMU 版本更新前后无感知稳定性验证）。

#### Harness Runner & Matrix
- `harness/matrix/matrix.yaml`：定义 QEMU 版本 (5.0.0 x 10.2.2) x 架构 (x86_64) x 场景矩阵。
- `harness/cmd/run/main.go`：CLI 驱动入口，支持一键运行 SOP 并导出结果。

---

## 每阶段完成标准（Checklist）

| 阶段 | 交付标志与验收标准 |
|---|---|
| 阶段一 | `go test ./qmp/... ./internal/version/...` 全部通过，传输层支持 Handshake 熔断、Context 与 Ready 等待 |
| 阶段二 | `make gen` 为 5.0.0 和 10.2.2 生成代码无报错，代码编译通过，Schema 解析器能够区分版本差异 |
| 阶段三 | `FeatureGateRegistry` 单元测试通过，`capability.yaml` 元数据正确加载并影响门控判断 |
| 阶段四 | 相同 `VMConfig` 在 5.0 下生成 `-drive`，在 10.2.2 下生成 `-blockdev`，`CmdlineDiff` 输出清晰 |
| 阶段五 | 三层 Dry-Run 校验测试通过，错误配置在对应层准确拦截 |
| 阶段六 | `MigrationNegotiator` 测试通过（`postcopy-preempt` 自动禁用），`LocalUpgrader` 端到端执行无报错 |
| 阶段七 | Harness SOP 矩阵全部 `required=true` 场景 PASSED，结构化报告成功生成 |

---

## 参考资料

[1] go-qemu 项目，digitalocean/go-qemu，GitHub 仓库: https://github.com/digitalocean/go-qemu
[2] qemu.qmp，QEMU 官方 Python QMP 客户端实现: https://gitlab.com/qemu-project/qemu/-/tree/master/python/qemu/qmp
[3] qmp-rs，Rust 实现的 QMP 动态客户端: https://github.com/oxidecomputer/qmp-rs
[4] qapi-rs，Rust 实现的类型化 QAPI 绑定: https://github.com/arcnmx/qapi-rs
[5] proxmox-ve 源码树，PVE::QemuServer Perl 模块: https://git.proxmox.com/?p=qemu-server.git
[6] libvirt-go，Go 语言的 libvirt 绑定: https://github.com/libvirt/libvirt-go
[7] QEMU QAPI Schema 文档: https://www.qemu.org/docs/master/devel/qapi-code-gen.html
[8] QEMU 5.0 发行说明: https://wiki.qemu.org/ChangeLog/5.0
[9] QEMU 迁移协议文档: https://www.qemu.org/docs/master/devel/migration/main.html
[10] KVM API 文档: https://www.kernel.org/doc/html/latest/virt/kvm/api.html
[11] QEMU QMP 参考手册: https://www.qemu.org/docs/master/interop/qemu-qmp-ref.html
[12] Proxmox VE 开发者文档: https://pve.proxmox.com/wiki/Developer_Documentation
