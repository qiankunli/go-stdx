# go-stdx

Go stdlib extensions — **admitted only when the standard library doesn't offer it.**

The long-term ambition is the slot Guava fills in Java: the utility layer a project reaches for before hand-writing a helper. Go's ecosystem splits that slot today — [samber/lo](https://github.com/samber/lo) owns generic collection transforms, [gods](https://github.com/emirpasic/gods) owns data structures, [lancet](https://github.com/duke-git/lancet) goes kitchen-sink — and go-stdx competes on discipline, not breadth: zero dependencies, stdlib-mirror naming, and only what fits that positioning.

The name is the admission rule: a hand-rolled `max`, a `slices.Clone` re-implementation, or a `strconv` wrapper does not belong here — use stdlib. What earns a slot is the loop real projects keep re-writing because the stdlib deliberately omits it, plus the tiny primitives not worth a heavyweight dependency.

Subpackages mirror stdlib naming so call sites read like the standard library they extend:

| package | what | why not stdlib / a lib |
|---|---|---|
| `slicesx` | `Uniq`, `UniqBy` — first-occurrence, order-preserving dedup | the "seen map" loop everyone re-writes; stdlib's `slices` has no transform family |
| `stringsx` | `Truncate` / `TruncateEllipsis` — byte-budget cuts; `FirstNonBlank`; `SplitAndTrim` | the log-truncation and human-written-list helpers every service re-writes (`cmp.Or` covers non-empty, not non-blank) |
| `osx` | `EnvStr` / `EnvBool` / `EnvInt` / `EnvInt64` / `EnvDuration` — typed env lookups with defaults; `WriteFileAtomic` — temp+rename write | config-from-env boilerplate in every service; `os.WriteFile` can leave readers a torn file |
| `ptrx` | `To` / `Value` / `FormatOr` — optional-field pointer helpers | no stdlib answer to `&literal` (k8s.io/utils/ptr exists because of the gap) |
| `filepathx` | `DirBytes` — recursive regular-file byte count | the quota/GC accounting walk everyone re-writes |
| `tarx` | `PackDir` / `UnpackDir` — directory ⇄ tar.gz with zip-slip defense | `archive/tar` leaves both loops to the caller, and the extraction loop is famously easy to get wrong |
| `shellx` | `Quote` — POSIX single-quoting | Go has no `shlex`; unquoted interpolation into a shell line is an injection |
| `randx` | `Hex(n)` — n random bytes as lowercase hex | the "short random id" helper every daemon re-writes |
| `uuid` | `V4`, `V7`, `V7Hex` — random / time-ordered ids | the id shapes services need without a full UUID dependency |

Rules of the house:

- **Zero dependencies**, forever. Everything here leans only on the standard library.
- **stdlib-first**: when Go's standard library grows an equivalent, the entry here is deprecated and removed.
- **Positioning is the bar**: anything genuinely generic that projects would otherwise hand-write belongs here — no waiting for N copies to accumulate first.

```go
import (
    "github.com/qiankunli/go-stdx/osx"
    "github.com/qiankunli/go-stdx/randx"
    "github.com/qiankunli/go-stdx/slicesx"
)

port := osx.EnvInt("APP_PORT", 8080)
id := "job-" + randx.Hex(6)
ids := slicesx.Uniq(rawIDs)
```

Used by [case-code-review](https://github.com/qiankunli/case-code-review), [hostel](https://github.com/qiankunli/hostel), and other Go projects under this account.

## License

MIT
