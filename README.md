# go-stdx

Go stdlib extensions — **admitted only when the standard library doesn't offer it.**

The name is the admission rule: a hand-rolled `max`, a `slices.Clone` re-implementation, or a `strconv` wrapper does not belong here — use stdlib. What earns a slot is the loop real projects keep re-writing because `slices`/`maps` deliberately omit it, plus the tiny primitives not worth a heavyweight dependency.

Subpackages mirror stdlib naming so call sites read like the standard library they extend:

| package | what | why not stdlib / a lib |
|---|---|---|
| `slicesx` | `Uniq`, `UniqBy` — first-occurrence, order-preserving dedup | the "seen map" loop everyone re-writes; stdlib's `slices` has no transform family |
| `uuid` | `V4()` — random RFC-4122 v4 id | one function; not worth a full UUID dependency |

Rules of the house:

- **Zero dependencies**, forever. Everything here leans only on the standard library.
- **stdlib-first**: when Go's standard library grows an equivalent, the entry here is deprecated and removed.
- **Rule of three**: nothing enters speculatively — an entry needs repeated hand-written copies across real projects before it earns a slot.

```go
import (
    "github.com/qiankunli/go-stdx/slicesx"
    "github.com/qiankunli/go-stdx/uuid"
)

ids := slicesx.Uniq(rawIDs)
ref := uuid.V4()
```

Used by [case-code-review](https://github.com/qiankunli/case-code-review) and other Go projects under this account.

## License

MIT
