# cassette

This is a personal fork of [TemirlanN-dev/cassette](https://github.com/TemirlanN-dev/cassette) (MIT license) for my own experiments, on the branch `mac-hwcnts`. Cassette is a proxy, that is, a program between an app and an LLM API. It records the API replies to a tape (a JSONL file) one time. After that, it plays the tape back as the server.

In this README, "the app" is the command that you give to cassette after `--`. Cassette starts the app and points its API address at cassette.

This branch adds code that runs inside cassette at each request and at the end of each reply. At those points, cassette reads the CPU counters of the app. This works on macOS and does not need root. The branch also adds controls for the replay timing.

## What this branch adds

| Addition | Used in | What it does |
|---|---|---|
| `--marks <file>` | `replay`, `auto` | Writes one JSON line when a request that the tape answers gets to cassette (`req`), and one line when the reply is fully sent (`end`). Each line holds the counters of the app. In `auto` mode, requests that go to the real API get no marks. |
| `--latency-scale x` | `replay`, `auto`, with `--realtime` | Multiplies all recorded delays by `x`. The value 0.5 makes the server two times faster. The value 2 makes it two times slower. |
| `header_delay_ms` | tape field | Holds the time from the request to the reply headers, that is, the time until the reply starts (the work of the server plus the network). `--realtime` plays this time back. Older tapes do not have this field, so replay uses 0. |
| Signposts | macOS | Each mark also sends a "Points of Interest" signpost (a time marker that Xcode Instruments records). The marks then show in Instruments traces on the clock of the trace. |
| `internal/counters` | package | On macOS it reads the counters with `proc_pid_rusage`. On other systems it is a stub that returns no counters. |

## Run it

1. Build the program. You need Go. On macOS you also need the Xcode command line tools, because cgo (the Go feature that calls C code) uses them.

   ```sh
   go build -o cassette .
   ```

2. Record the replies of the app. Cassette sets `ANTHROPIC_BASE_URL` and `OPENAI_BASE_URL` for the app to its own address.

   ```sh
   ./cassette record --target https://api.anthropic.com --cassette tape.jsonl -- <the app>
   ```

3. Play the tape back. Without `--realtime`, cassette sends each saved reply at once.

   ```sh
   ./cassette replay --cassette tape.jsonl -- <the app>
   ```

4. To play the tape back with the recorded timing and read the counters at each request and reply end, add these flags:

   ```sh
   ./cassette replay --realtime --latency-scale 1 --cassette tape.jsonl --marks marks.jsonl -- <the app>
   ```

   Marks need the app after `--`, because cassette reads only a process that it started. Without the app, every line has `ok: false`.

If a request in the replay is not in the tape, cassette answers with status 404 and writes `MISS` to its log. Some request fields change between runs. For these fields, put `ignore` and `replace` rules in a `cassette.yaml` file, and give the file with `--config`.

## Replay under Xcode Instruments

You can record a trace of all processes in Instruments while cassette plays a tape back. The signposts then show each `req` and `end` on the clock of the trace. This step does not need root. Wait a few seconds after you start `xctrace`, because it starts recording only after a short delay. Stop it with an interrupt signal, so that it saves the trace.

```sh
xcrun xctrace record --template 'CPU Counters' --all-processes --output run.trace --time-limit 120s &
sleep 3
./cassette replay --realtime --cassette tape.jsonl --marks marks.jsonl -- <the app>
kill -INT %1; wait
```

## The marks file

Each line is one reading. All counter values are totals since the app started. To get the value for one step, subtract one line from another. For example, `end` of request k minus `req` of request k is the wait for that reply. `req` of request k+1 minus `end` of request k is the work that the app does between two replies.

| Field | Meaning |
|---|---|
| `seq`, `ev` | The number of the request in this run, and `req` or `end`. |
| `t_ns`, `t_wall_ns` | CLOCK_MONOTONIC (a clock that only counts up) and Unix time, read at the same moment. Use `t_wall_ns` to line up logs that use the wall clock. |
| `pid`, `ok` | The process that cassette read, and whether the read worked. |
| `instr`, `cycles`, `p_instr`, `p_cycles` | Instructions and cycles, and the part of each that ran on P-cores (the fast cores). |
| `energy_nj` | The energy estimate of macOS for the process, in nanojoules. It is not a measurement. |
| `user_ns`, `sys_ns`, `child_user_ns`, `child_sys_ns`, `footprint_bytes` | CPU time, the CPU time of child processes that already stopped, and memory. |

Pair the lines by `seq`, not by time order, because two requests can be open at the same time.

## Limits

- Cassette reads only the app process. It does not count the child processes of the app, except their CPU time after they stop.
- If the app is a wrapper (a shell or a script that starts the real program), cassette reads the wrapper. Use `exec` in the wrapper, or give the real program after `--`.
- Cassette sees only requests and replies. It cannot see what the app does between them.
- The counters work only on macOS. On Linux and Windows, the CI builds the stub.

## Test

Run the tests and the benchmark for one mark. One mark takes about 3 µs on an M4.

```sh
go vet ./... && go test ./...
go test ./internal/proxy -run '^$' -bench Mark
```

## Upstream

The proxy, the matching rules, and the `record`, `replay`, `show` and `diff` commands come from upstream cassette. The original README is in the git history. It explains the other features: the rules in `cassette.yaml`, the `auto` mode, and the `show` and `diff` commands. The license is MIT (see `LICENSE`).
