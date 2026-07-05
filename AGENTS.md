# AGENTS.md — go-stdx

## 项目定位与边界

Go stdlib 扩展库，长期对标 Java 里 Guava 的位置——项目手写 helper 之前先来这里找。与现有生态分工：samber/lo 占泛型集合、gods 占数据结构、lancet 是 kitchen-sink；go-stdx 拼的是纪律而非广度。**名字即收录纪律**：stdlib 没有的才进（手搓 max、复刻 slices.Clone、包一层 strconv 都不属于这里）。永远零依赖；stdlib 补齐等价物后这里的条目废弃删除。符合定位即可收录，不设重复次数门槛。

## 代码地图

子包镜像 stdlib 命名（`slicesx` / `osx` / `filepathx` / `tarx` / `shellx` / `randx` / `uuid`），调用点读起来像它扩展的那个标准库。一包一职责，包内保持小。

## 关键约定

1. 收录评审问三件事：stdlib 真没有吗（含最新版本）？是真通用还是某项目的业务形状？零依赖能实现吗？
2. 语义约定：Uniq 系保留**首次出现**且保序；所有函数 nil-in-nil-out。
3. 消费方（ccr / hostel / …）内部不许再就地手写本库已有的操作——消费仓的 AGENTS.md 应有对应约定。

## References

- 首个消费方与孵化史：[case-code-review](https://github.com/qiankunli/case-code-review) `pkg/stdx`（已迁出）
- 第二个消费方：[hostel](https://github.com/qiankunli/hostel)——`osx`/`randx`/`shellx`/`tarx`/`filepathx` 均由其内部手写实现迁入
