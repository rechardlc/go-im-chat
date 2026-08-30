# Make 使用教程

Make 是独立安装的命令行工具，不绑任何语言或仓库。你在目录里放一个名为 `Makefile` 的文件，把终端里反复敲的命令写成 **target**，之后只输入 `make 名字`。

前端类比：根目录 `package.json` 的 `scripts`。`make build` ≈ `npm run build`。差别是 Make 默认按**文件是否过期**决定要不要重跑，而不只是「按名字执行一段脚本」。

下面整篇用一个**独立空目录**练习，和任何业务仓库无关。

---

## 1. 安装与确认

```bash
make --version    # 应看到 GNU Make x.x
```

- macOS：`xcode-select --install` 或 `brew install make`
- Linux：`sudo apt install make` / `sudo yum install make`
- Windows：不要用纯 PowerShell 硬跑 POSIX 配方。用 **Git Bash 或 WSL**，再装 GNU Make（`choco install make` / `scoop install make` 也可以，但配方里的 `rm`、`test` 仍需要 Unix shell）

约定：命令都在**放着 Makefile 的那个目录**执行。`make` 默认找当前目录的 `Makefile`（或 `makefile`）。

---

## 2. 最小例子（先当脚本用）

新建空目录，写入 `Makefile`：

```makefile
.PHONY: hello
hello:
	echo hello, make
```

```bash
make hello
# hello, make
```

发生了什么：

| 名词 | 含义 | 前端类比 |
|------|------|----------|
| `hello` | target 名，`make hello` 的那个词 | `"scripts"` 里的 key |
| 下面缩进的行 | recipe（配方），真正丢给 shell 的命令 | script 的 value |
| `.PHONY: hello` | 声明这是动作，不是要生成的文件 | script 本来就不会和磁盘文件搞混，Make 需要显式说 |

三条硬规则：

1. 配方必须用 **Tab** 缩进，空格会报 `missing separator`。
2. 每行配方是一条独立 shell。`cd foo` 只对**这一行**有效，下一行回到 Makefile 所在目录。要连着做：`cd foo && npm test`。
3. 裸跑 `make`（不写名字）会执行**文件里第一个 target**。所以第一个 target 要么是你最常用的，要么专门写一个 `help`。

到这里，Make 已经能当「命令别名」用。下一节才是它和 npm scripts 真正不同的地方。

---

## 3. Make 真正在优化什么：文件过期

Make 的原始模型是：**target 是一个要产出的文件**；依赖是产出它需要的原料；原料比产物新，才重跑配方。

独立小项目，目录如下：

```
hello-make/
  Makefile
  src/
    greet.c
```

`src/greet.c`：

```c
#include <stdio.h>
int main(void) { printf("hi\n"); return 0; }
```

`Makefile`：

```makefile
out/greet: src/greet.c
	mkdir -p out
	cc -o out/greet src/greet.c

.PHONY: clean
clean:
	rm -rf out
```

```bash
make out/greet    # 第一次：编译
make out/greet    # 第二次：Nothing to be done（源文件没变）
touch src/greet.c
make out/greet    # 源文件更新了：再编译
make clean        # 删产物；clean 是动作，所以要 .PHONY
```

读法：

```
target: 依赖1 依赖2
	配方
```

- target `out/greet` 是磁盘上的产物。
- 依赖 `src/greet.c` 是原料。
- Make 比较时间戳：产物不存在，或依赖比产物新 → 跑配方；否则跳过。

这就是增量构建。前端类比：不是每次都 `npm run build`，而是「源码比 `dist/` 新才重建」——Make 把这套做成了默认语义。

`.PHONY` 的意义现在清楚了：`clean`、`hello` 这类名字**不是文件**。若目录里碰巧有个叫 `clean` 的文件，没有 `.PHONY` 时 Make 会认为产物已存在、跳过删除。动作类 target 一律声明 `.PHONY`。

---

## 4. 语法速查（够用的那一部分）

### 4.1 依赖链

```makefile
all: out/greet

out/greet: out/greet.o
	cc -o out/greet out/greet.o

out/greet.o: src/greet.c
	mkdir -p out
	cc -c -o out/greet.o src/greet.c
```

`make all` → 发现要 `out/greet` → 发现要 `out/greet.o` → 先编译 `.o`，再链接。依赖是图，不是简单从上到下执行。

常驻进程（`go run`、`npm run dev`）不要做成依赖链：依赖必须**结束**才会跑下一个。两个 dev server 要用并行，见 4.5。

### 4.2 变量

```makefile
CC      := cc
CFLAGS  := -Wall -O2
SRC     := src/greet.c
BIN     := out/greet

$(BIN): $(SRC)
	mkdir -p $(dir $(BIN))
	$(CC) $(CFLAGS) -o $(BIN) $(SRC)
```

| 写法 | 含义 |
|------|------|
| `:=` | 立刻求值，几乎总是用这个 |
| `=` | 延迟求值，引用其他变量时才展开，容易绕晕 |
| `$(CC)` | 展开变量 |
| `$(dir $(BIN))` | 内置函数：取目录部分 `out/` |

覆盖：`make CC=clang out/greet` 会用命令行传入的 `CC`。≈ `npm run build -- --flag`，但是 Make 变量是一等公民。

### 4.3 自动变量（配方里常用）

配方执行时 Make 填好的临时变量：

| 变量 | 含义 |
|------|------|
| `$@` | 当前 target 名（要生成的那个文件） |
| `$<` | 第一个依赖 |
| `$^` | 全部依赖，空格分隔 |

```makefile
out/greet: src/greet.c
	mkdir -p out
	cc -o $@ $<
```

`$@` → `out/greet`，`$<` → `src/greet.c`。模式规则里几乎必用。

### 4.4 模式规则（一类文件一种配方）

```makefile
out/%.o: src/%.c
	mkdir -p out
	cc -c -o $@ $<

out/app: out/a.o out/b.o
	cc -o $@ $^
```

`%` 是通配。`out/a.o` 会去找 `src/a.c`。不必为每个 `.c` 手写一条。前端类比：webpack 的 `test: /\.ts$/` 一条规则覆盖一类文件，而不是每个文件一条 loader。

### 4.5 `@`、`-`、`-j`

```makefile
.PHONY: silent risky parallel

silent:
	@echo 这行命令本身不会被打印，只看输出

risky:
	-false          # 这行失败也继续
	echo still here

# 两个独立 target 并行（都不是常驻也可以加速）
# make -j2 a b
```

| 前缀 / 开关 | 作用 |
|-------------|------|
| `@` | 不 echo 命令本身 |
| `-` | 这行失败（非 0）不让 Make 停下来 |
| `-jN` | 最多 N 路并行跑没有互相依赖的 target |
| `-n` | dry-run，只打印将执行的命令 |

并行例子：`make -j2 out/a.o out/b.o`。有依赖的仍会按图等待。

两个**常驻**进程：

```makefile
.PHONY: api web local
api:
	cd api && npm run dev
web:
	cd web && npm run dev
local:
	$(MAKE) -j2 api web
```

不要写成 `local: api web`：Make 会等 `api` 退出才启动 `web`。`$(MAKE)` 是再调一次 Make，能正确继承 `-j`；不要在配方里写死 `make`。

### 4.6 条件与「命令存在才拷文件」

```makefile
.PHONY: env
env:
	@test -f .env || cp .env.example .env
```

`test -f` 是 POSIX：文件存在则为真。`||` 只在不存在时拷贝。这是脚本习惯，不是 Make 独有。

---

## 5. 独立小项目：从零写到能用

目标：一个迷你 C 程序，`make` 默认构建，`make run` 运行，`make clean` 清理。把下面文件拷进空目录即可。

**src/main.c**

```c
#include <stdio.h>
int main(void) {
    printf("ok\n");
    return 0;
}
```

**Makefile**

```makefile
CC      := cc
CFLAGS  := -Wall
SRC     := $(wildcard src/*.c)
OBJ     := $(patsubst src/%.c,out/%.o,$(SRC))
BIN     := out/app

.PHONY: all run clean

all: $(BIN)

$(BIN): $(OBJ)
	mkdir -p $(dir $@)
	$(CC) -o $@ $^

out/%.o: src/%.c
	mkdir -p out
	$(CC) $(CFLAGS) -c -o $@ $<

run: $(BIN)
	./$(BIN)

clean:
	rm -rf out
```

```bash
make          # 默认 all → 编译 out/app
make run      # 依赖 all 的产物，然后执行
make clean
```

多写几个 `src/*.c` 再 `make`：只重编改过的 `.c` 对应的 `.o`，再重新链接。这就是 Make 比「一个 bash 脚本从头编」值钱的地方。

`$(wildcard src/*.c)` 列出匹配文件；`$(patsubst src/%.c,out/%.o,$(SRC))` 把路径换成 `.o`。遇到再查，不必一次记完所有函数。

---

## 6. 常见坑

1. **Tab 不是空格。** 编辑器开「以空格代替 Tab」会写废 Makefile。
2. **默认按文件想问题。** 动作类 target（`dev`、`test`、`clean`）必须 `.PHONY`，否则同名文件会让命令被跳过。
3. **依赖是「等结束」。** 长驻进程用 `-j`，不要用依赖串起来。
4. **每一行配方是新 shell。** `cd`、`export` 不跨行；用 `&&` 或 `\` 续行。
5. **第一个 target 是默认。** 习惯把 `all` 放文件最上面。
6. **配方失败默认停。** 需要「停不掉也继续」才加行首 `-`。
7. **Windows 路径与 shell。** 配方按 Unix 写（`rm -rf`、`mkdir -p`）。Git Bash / WSL 最省事。

---

## 7. 和 npm scripts 怎么选

| | npm scripts | Make |
|--|-------------|------|
| 默认语义 | 按名字跑一段命令 | 按文件时间戳增量构建 |
| 依赖 | `prebuild` 约定，或自己 `&&` | 一等：`target: deps` |
| 并行 | `concurrently` 等 | `-j` |
| 适合 | Node 项目日常脚本 | 多文件编译、codegen、混合语言、与语言无关的任务图 |

只是少敲几下终端：Make 或 npm scripts 都行。一旦有「很多中间产物、只想重做脏的那部分」，Make 的模型更合适。现代也有 Ninja、Bazel、Taskfile；概念都是「图 + 过期检测」，先把 Make 吃透即可。

---

