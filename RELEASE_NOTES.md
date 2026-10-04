# Flex-WavelogGate v0.1.0

The first public release of **Flex-WavelogGate**, created by **Oliver (OM0RX)**.

## Highlights

- Dedicated xCAT RigCtlD connection for FlexRadio stations.
- Live frequency and mode reporting to Wavelog.
- Correct voice, CW, RTTY, AM and FM mode normalization.
- FT8 and FT4 recognition on their conventional dial frequencies when xCAT
  reports a generic data mode such as `PKTUSB`.
- Existing WaveLogGate pairing and station profiles migrate into a separate
  Flex-WavelogGate configuration.
- Wavelog QSY support retained; Flex-WavelogGate never keys PTT.
- Rotor support remains separate and unchanged.

## Downloads

| Platform | Choose |
|---|---|
| macOS Apple Silicon | `Flex-WavelogGate-darwin-arm64.dmg` |
| macOS Intel | `Flex-WavelogGate-darwin-amd64.dmg` |
| Windows 64-bit | `Flex-WavelogGate-windows-amd64-x64.exe` |
| Windows ARM64 | `Flex-WavelogGate-windows-arm64.exe` |
| Windows 32-bit | `Flex-WavelogGate-windows-i386-x86.exe` |
| Ubuntu/Debian | matching `.deb` architecture and WebKit version |
| Fedora/RHEL | matching `.rpm` architecture and WebKit version |

Download `SHA256SUMS.txt` to verify your selected artifact.

## Before starting

Read the [xCAT setup guide](docs/XCAT_SETUP.md). Use the **RigCtlD port**, which
is normally `4532`; do not use xCAT's CAT port `5002`.

The public macOS and Windows builds are community builds and may not be signed
with a commercial platform certificate. On macOS, Control-click the app and
choose **Open**. On Windows, review the SmartScreen details before choosing to
run it.

## Proven configuration

The release was validated with DL3LSM xCAT 2.0 on macOS, a live FlexRadio slice,
and a paired Wavelog instance. Live checks covered LSB voice operation and FT8
at 7.074 MHz without sending PTT or rotor commands.
