# moo

A small, friendly code generator that adds Moo-style ergonomics to Go structs:
functional options, defaults, validation, `Clone()`, and `String()`.

**No reflection. No runtime magic.** `moo` reads struct tags at build time
and writes plain, idiomatic Go that GoLand (and every other Go tool) can
autocomplete, refactor, and navigate normally.

---

## Table of contents

- [Installation](#installation)
- [Quick start](#quick-start)
- [Example 1 — A user signup form](#example-1--a-user-signup-form)
- [Example 2 — A server config with defaults](#example-2--a-server-config-with-defaults)
- [The complete tag reference](#the-complete-tag-reference)
- [Semantics — the rules, spelled out](#semantics--the-rules-spelled-out)
- [What actually gets generated](#what-actually-gets-generated)
- [Project layout](#project-layout)
- [Daily workflow](#daily-workflow)
- [Troubleshooting](#troubleshooting)
- [Limitations](#limitations)
- [License](#license)

---

## Installation

```bash
go install github.com/thomas666-beast/moo/cmd/moo@latest
```

If the command `moo` isn't found afterwards, add Go's bin directory to your
`PATH`:

```bash
export PATH="$HOME/go/bin:$PATH"
```

Verify:

```bash
moo -h
```

Should print:

```
usage: moo [-debug] [-stdout] <file.go>
  -debug
        print parsed IR instead of generating code
  -stdout
        write generated code to stdout instead of a file
```

---

## Quick start

1. **Write** a Go struct and add `moo:"..."` tags to the fields you care about.
2. **Run** `go generate ./...` (or `moo yourfile.go`).
3. **Use** the generated `NewX(...)`, `WithField(...)`, `Clone()`, and `String()` helpers.

That's it.

---

## Example 1 — A user signup form

The simplest useful case: a struct that must have valid data before you use it.

### Step 1 — `user.go`

```go
package main

//go:generate moo user.go

type User struct {
    Name  string   `moo:"required;min=1;max=64"`
    Email string   `moo:"required;match=email"`
    Age   int      `moo:"min=0;max=150"`
    Role  string   `moo:"default=\"user\";enum=user|admin"`
    Tags  []string `moo:""`
}
```

Read the tags out loud:

- `Name` — must be set, between 1 and 64 characters.
- `Email` — must be set, must look like an email.
- `Age` — if set, must be 0–150.
- `Role` — defaults to `"user"`; if set, must be `"user"` or `"admin"`.
- `Tags` — no rules, but still gets a `WithTags` option and is deep-copied
  by `Clone`.

### Step 2 — Generate

```bash
go generate ./...
```

This creates `zz_moo_user.go` next to `user.go`. Commit it — the generated
code is part of your package.

### Step 3 — `main.go`

```go
package main

import "fmt"

func main() {
    // Happy path
    u, err := NewUser(
        WithName("Alice"),
        WithEmail("alice@example.com"),
        WithAge(30),
        WithTags([]string{"beta", "early"}),
    )
    if err != nil {
        panic(err)
    }
    fmt.Println(u)
    // User{Name: "Alice", Email: "alice@example.com", Age: 30, Role: "user", Tags: [beta early]}

    // Validation catches problems at construction time
    _, err = NewUser(WithName("Bob"))
    fmt.Println("missing email:", err)
    // missing email: Email is required

    _, err = NewUser(WithName("Bob"), WithEmail("not-an-email"))
    fmt.Println("bad email:", err)
    // bad email: Email must be a valid email: ...

    // Clone is independent
    v := u.Clone()
    v.Tags[0] = "changed"
    fmt.Println("original:", u.Tags) // [beta early]
    fmt.Println("clone:   ", v.Tags) // [changed early]
}
```

### What you got for free

| Function        | What it does                                                 |
|-----------------|--------------------------------------------------------------|
| `NewUser(...)`  | Builds a user, applies defaults, validates, returns `error`. |
| `WithName(...)` | One functional option per writable field.                    |
| `u.Validate()`  | Re-check any time (e.g., after mutation).                    |
| `u.Clone()`     | Deep-copies slices/maps so the copy is independent.          |
| `u.String()`    | Readable output for logs — `User{Name: "Alice", ...}`.       |

---

## Example 2 — A server config with defaults

The same idea, but for settings you load at startup. Configs usually have
defaults and every field is optional, so `required` disappears.

### Step 1 — `config.go`

```go
package main

//go:generate moo config.go

type Config struct {
    Host     string `moo:"default=\"localhost\""`
    Port     int    `moo:"default=8080;min=1;max=65535"`
    Timeout  int    `moo:"default=30;min=1"`
    LogLevel string `moo:"default=\"info\";enum=debug|info|warn|error"`
}
```

### Step 2 — Generate

```bash
go generate ./...
```

### Step 3 — Use it

```go
package main

import "fmt"

func main() {
    // All defaults
    c, err := NewConfig()
    if err != nil {
        panic(err)
    }
    fmt.Println(c)
    // Config{Host: "localhost", Port: 8080, Timeout: 30, LogLevel: "info"}

    // Override just what you need
    c, err = NewConfig(
        WithPort(9090),
        WithLogLevel("debug"),
    )
    if err != nil {
        panic(err)
    }
    fmt.Println(c)
    // Config{Host: "localhost", Port: 9090, Timeout: 30, LogLevel: "debug"}

    // Bad values are rejected immediately
    _, err = NewConfig(WithPort(99999))
    fmt.Println("bad port:", err)
    // bad port: Port must be <= 65535

    _, err = NewConfig(WithLogLevel("verbose"))
    fmt.Println("bad level:", err)
    // bad level: LogLevel must be one of: debug|info|warn|error
}
```

### Why this pattern is nice

- You never write a 40-line `NewConfig` again.
- Every default lives next to the field it belongs to — no scattered
  `if c.Port == 0 { c.Port = 8080 }`.
- Bounds (`min`, `max`) and allowed values (`enum`) are checked the moment
  the config is built, not at first use.

---

## The complete tag reference

Tags are a `;`-separated list. Each item is either a bare word or
`key=value`.

| Tag        | Example                  | Meaning                                                             |
|------------|--------------------------|---------------------------------------------------------------------|
| `-`        | `moo:"-"`                | Ignore this field completely. No option, no validation, no `Clone`. |
| `readonly` | `moo:"readonly"`         | No `WithX` option. Still validated, cloned, printed.                |
| `required` | `moo:"required"`         | Must be non-zero / non-nil / non-empty.                             |
| `default`  | `moo:"default=42"`       | Go expression used by `NewX`. Written verbatim into the constructor.|
| `min`      | `moo:"min=1"`            | Numeric: `>= 1`. String/slice/map: length `>= 1`.                   |
| `max`      | `moo:"max=100"`          | Numeric: `<= 100`. String/slice/map: length `<= 100`.               |
| `enum`     | ``moo:"enum=a\|b\|c"``   | Value must be one of these (checked only when non-empty).           |
| `match`    | `moo:"match=email"`      | Value must satisfy a named validator (checked only when non-empty). |

Combine any of them:

```go
Name     string `moo:"required;min=1;max=64"`
Email    string `moo:"required;match=email"`
Role     string `moo:"default=\"user\";enum=user|admin"`
ID       int64  `moo:"readonly"`
Internal string `moo:"-"`
```

---

## Semantics — the rules, spelled out

Understanding this table is 90% of understanding `moo`.

| Tag           | When is it checked?               | Why?                                            |
|---------------|-----------------------------------|-------------------------------------------------|
| `required`    | Always.                           | "Must be set" is unconditional.                 |
| `min` / `max` | Always.                           | `0` is a legitimate numeric value.              |
| `enum`        | Only when the field is non-empty. | Empty usually means "use default / not set".    |
| `match=email` | Only when the field is non-empty. | Same reasoning.                                 |

**Rule of thumb:** if you want `enum` or `match` to *also* require a value,
add `required`:

```go
Email string `moo:"required;match=email"`
```

For field kinds, the checks adapt automatically:

| Field kind     | `required` fails when…       | `min`/`max` applies to… |
|----------------|------------------------------|-------------------------|
| `string`       | `== ""`                      | length                  |
| numbers / bool | `== 0` / `== false`          | value                   |
| pointer        | `== nil`                     | (skipped)               |
| slice / map    | `len(...) == 0`              | length                  |

---

## What actually gets generated

For each struct with at least one `moo:` tag, you get exactly six things:

1. **`type UserOption func(*User)`** — the option type.
2. **`WithName(v T) UserOption`** — one per writable field.
3. **`NewUser(opts ...UserOption) (*User, error)`** — applies defaults, runs
   options, calls `Validate`.
4. **`(u *User) Validate() error`** — re-runs every check at any time.
5. **`(u *User) Clone() *User`** — safe copy; slices and maps are duplicated.
6. **`(u *User) String() string`** — readable text, great for logs.

Nothing else. The generated file starts with:

```go
// Code generated by moo. DO NOT EDIT.
```

so GoLand greys it out and skips it for most inspections.

---

## Project layout

```
cmd/moo/             CLI entry
internal/model/      IR types (StructSpec, FieldSpec)
internal/tag/        tag grammar parser
internal/parser/     Go AST → IR
internal/gen/        IR → Go source, with golden tests
e2e/                 end-to-end package that uses generated code
```

---

## Daily workflow

```bash
# One-time setup
go install github.com/thomas666-beast/moo/cmd/moo@latest

# After editing a struct in your project
go generate ./...
go test ./...
```

If `go generate ./...` leaves `git status` dirty, your checked-in generated
files are stale. Regenerate and commit.

### Drift detection (CI)

Add this to CI to fail when someone forgets to regenerate:

```bash
go generate ./...
git diff --exit-code -- '*.go'
```

---

## Troubleshooting

### `moo: executable file not found in $PATH`

`go generate` looks for `moo` as a command. Install it:

```bash
go install github.com/thomas666-beast/moo/cmd/moo@latest
```

Then check:

```bash
which moo
```

If empty, `$GOPATH/bin` (usually `~/go/bin`) isn't on your `PATH`:

```bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

### You're developing `moo` locally

Install from your checkout instead of from `@latest`:

```bash
cd /path/to/moo
go install ./cmd/moo
```

Now `moo` is on your `PATH` for any project on your machine.

### You don't want to install anything

Skip `go generate` and run the tool directly:

```bash
go run github.com/thomas666-beast/moo/cmd/moo user.go
```

Works, but you'll type that each time. Installing is better.

### `moo:` tag is ignored

- The struct must have **at least one** field with a `moo:` tag to be
  considered. Structs with no tagged fields are skipped entirely.
- Unexported fields are skipped for options (you can't set them from outside
  the package) but still appear in `Clone` and `String`.
- Embedded (anonymous) fields are not supported yet.

### Generated file looks stale

Delete `zz_moo_*.go` and rerun `go generate ./...`. If that doesn't help,
your `moo` binary is stale — `go install ./cmd/moo` from the `moo` repo.

---

## Limitations

- `enum` currently only meaningful for `string` fields.
- `Clone` deep-copies slices and maps, but not pointers or nested slices.
- `String` uses `%q` for strings and `%v` otherwise; no JSON output.
- No support for embedded structs yet.
- No roles/interfaces beyond plain Go interfaces.
- One file per invocation; package-wide generation is on the roadmap.

---

## License

MIT
