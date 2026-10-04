# Flex-WavelogGate

Flex-WavelogGate is a customised WaveLogGate desktop application for stations
where a FlexRadio is controlled by [xCAT](https://dl3lsm.blogspot.com/) and radio
state is reported to a paired [Wavelog](https://github.com/wavelog/wavelog)
instance.

This fork keeps WaveLogGate's Wavelog pairing, station profiles, QSO forwarding,
QSY endpoint, WebSocket status, and rotator implementation. Its radio path is
intentionally xCAT-only.

## What changed

- The desktop product and bundle are named **Flex-WavelogGate**.
- The only selectable radio backend is xCAT's RigCtlD-compatible TCP service.
- The default xCAT endpoint is `127.0.0.1:4532`.
- Polling only uses xCAT's supported `f` (frequency) and `m` (mode/passband)
  commands. It does not continuously send unsupported Hamlib SAT, split, or
  power probes.
- Wavelog QSY requests use xCAT's `F` and `M` commands.
- Existing WaveLogGate profiles are copied on first launch so their Wavelog URL,
  API key, station ID, and radio name remain paired.
- Rotor control is unchanged and remains a separate `rotctld` integration.

## Critical xCAT port choice

xCAT can expose two different services:

| xCAT field | Typical port | Use here? |
|---|---:|---|
| CAT Port | `5002` | No — this is the radio-specific CAT protocol |
| RigCtlD Port | `4532` | **Yes** — Flex-WavelogGate uses this protocol |

Pointing generic WaveLogGate Hamlib settings at xCAT's CAT port `5002` does not
work. Flex-WavelogGate connects to xCAT's RigCtlD port `4532` instead.

## macOS setup

1. Start SmartSDR/FlexRadio and xCAT.
2. In xCAT, select the station/slice and ensure the row reports `CAT and RigCtlD`.
3. Confirm the RigCtlD port is `4532` and the status says `Slice present`.
4. Start Flex-WavelogGate.
5. Under **Configuration → Wavelog**, enter the Wavelog URL, API key, station,
   and desired radio name.
6. Under **xCAT Radio Control**, enable the backend and use host `127.0.0.1`,
   RigCtlD port `4532`.
7. Save. The Status page should show the active slice frequency and mode.

Frequency and mode changes are posted to Wavelog's `/api/radio` endpoint (or
`/api/v2/radio` for a v2 token) whenever the observed state changes, with a
30-minute refresh even when it remains stable.

## Existing WaveLogGate settings

On first launch, if Flex-WavelogGate has no configuration of its own, it reads
the existing WaveLogGate configuration and writes a separate migrated copy:

```text
~/Library/Application Support/WavelogGate/config.json
  → ~/Library/Application Support/Flex-WavelogGate/config.json
```

The original file is not edited or removed. Wavelog pairing fields are
preserved; legacy FLRig/Hamlib radio selections are disabled and xCAT is enabled
at `127.0.0.1:4532`.

API keys are never compiled into the application or committed to this
repository. They remain in the per-user configuration file.

## Radio behavior and limits

- Read: active slice frequency and mode.
- Wavelog reporting: radio name, frequency, and mode.
- Mode compatibility: xCAT voice/CW/RTTY/AM/FM modes are normalized to values
  Wavelog can select. On conventional FT8 and FT4 dial frequencies, xCAT's
  generic data modes (`PKTUSB`, `DIGU`, and equivalents) are reported as `FT8`
  or `FT4`; away from those dial frequencies they safely fall back to `USB` or
  `LSB` because CAT cannot identify the application protocol by itself.
- QSY from Wavelog: frequency and, when enabled, mode.
- PTT: never asserted by Flex-WavelogGate.
- Split-VFO and RF-power reporting: not queried because xCAT 2.0 does not expose
  the generic Hamlib commands required by the upstream client.
- Rotor control: unchanged; this fork does not route rotor commands through
  xCAT.

## Other retained WaveLogGate services

| Port | Protocol | Purpose |
|---:|---|---|
| `2333` | UDP inbound | QSO packets from WSJT-X/FLDigi |
| `54321` | HTTP/HTTPS inbound | QSY request from Wavelog |
| `54322` | WebSocket | Local radio status and Wavelog events |
| `54323` | Secure WebSocket | TLS WebSocket endpoint |

For WSJT-X, use the **Secondary UDP Server** with `localhost:2333`.

## Development

Requirements:

- Go 1.25+
- Wails CLI 2.13+
- Bun

Generate Wails bindings and build:

```bash
wails generate module
cd frontend && bun install --frozen-lockfile && bun run build && cd ..
wails build -clean -platform darwin/arm64
```

Run the automated tests:

```bash
go test ./...
```

With xCAT running locally, the read-only live integration test is:

```bash
go test -tags integration -run TestLiveXCAT -v ./internal/radio
```

That test only requests frequency and mode. It does not tune the radio or key
PTT.

## Upstream and licence

Flex-WavelogGate is based on WaveLogGate v2.1.1 at upstream commit
`9344a4852ad461b60a50fbb384e2afdf17f6e414`. See [LICENSE](LICENSE) for the
upstream licence terms.
