# AGENTS.md — go-stdx

## 项目定位与边界

Go stdlib 扩展库，**名字即收录纪律**：stdlib 没有的才进（手搓 max、复刻 slices.Clone、包一层 strconv 都不属于这里）。永远零依赖；stdlib 补齐等价物后这里的条目废弃删除；**三次法则**——真实项目里重复手写 ≥3 次才配一个槽位，不预铸。

## 代码地图

子包镜像 stdlib 命名（`slicesx` / `uuid` / …），调用点读起来像它扩展的那个标准库。一包一职责，包内保持小。

## 关键约定

1. 收录评审问三件事：stdlib 真没有吗（含最新版本）？有没有 ≥3 处真实重复？零依赖能实现吗？
2. 语义约定：Uniq 系保留**首次出现**且保序；所有函数 nil-in-nil-out。
3. 消费方（ccr / hostel / …）内部不许再就地手写本库已有的操作——消费仓的 AGENTS.md 应有对应约定。

## References

- 首个消费方与孵化史：[case-code-review](https://github.com/qiankunli/case-code-review) `pkg/stdx`（已迁出）
