# xCAT setup for Flex-WavelogGate

This guide describes the tested xCAT integration for a FlexRadio 6000/8000
station. The proven setup uses **DL3LSM xCAT 2.0 on macOS**.

> There is unrelated cluster-management software also named xCAT. Do not
> install that product. The radio application is published by Mario, DL3LSM.

## Download xCAT

Use the author's official pages:

- [DL3LSM xCAT/xDAX downloads](https://dl3lsm.blogspot.com/p/downloads.html)
- [DL3LSM radio software blog](https://dl3lsm.blogspot.com/)

Read the **Very important notice** on the download page before installing.
DL3LSM documents xCAT 2.0 for SmartSDR 2.5.1 and newer, and keeps xCAT 1.0 for
SmartSDR 2.4.9 and older.

Flex-WavelogGate itself is released for macOS, Windows and Linux. DL3LSM's xCAT
2.0 build used for this integration is macOS software. The Windows/Linux builds
of Flex-WavelogGate may connect to an xCAT RigCtlD endpoint on another machine
only when that endpoint is reachable over the LAN. Running both applications on
the same Mac with `127.0.0.1` is the recommended and tested arrangement.

## Configure the xCAT slice

1. Start the FlexRadio client or Maestro connection you intend to use.
2. Open xCAT and select the correct station at the top of its window.
3. Click **Add…** or edit the existing slice entry.
4. Select the required slice, normally **A**.
5. Select **CAT and RigCtlD** as the protocol.
6. Keep xCAT's normal CAT port if another application needs it. A typical CAT
   port is `5002`, but Flex-WavelogGate does not use it.
7. Set or note the **RigCtlD Port**. The recommended value is `4532`.
8. Confirm the xCAT row reports **Slice present** and the connection indicator
   is green.

The important distinction is:

| xCAT field | Typical value | Flex-WavelogGate |
|---|---:|---|
| CAT Port | `5002` | Do not use |
| RigCtlD Port | `4532` | Use this |

The **Switch TX** and **Split Mode** choices are not required for status
reporting. Flex-WavelogGate does not assert PTT, and its first release does not
control split VFOs through xCAT.

## Configure Flex-WavelogGate

1. Open **Configuration → Wavelog**.
2. Enter the full Wavelog base URL, API key, station and radio name.
3. Test the Wavelog connection and save it.
4. Open **xCAT Radio Control** and enable xCAT.
5. Use host `127.0.0.1` and port `4532` when xCAT is on the same Mac.
6. Save and return to **Status**.

The status card should show the selected slice frequency and mode. In Wavelog,
open **Hardware Interfaces**, set the new radio as default, then open **Add QSO**.
Wavelog's frequency, band and mode should follow the radio after its polling
interval.

## Mode behavior

| xCAT reports | Wavelog receives |
|---|---|
| `USB`, `LSB` | the corresponding voice sideband |
| `CW`, `CWR` | `CW` |
| `RTTY`, `RTTYR` | `RTTY` |
| `AM`, `FM`, `PKTFM` | `AM` or `FM` |
| data mode on an FT8 dial frequency | `FT8` |
| data mode on an FT4 dial frequency | `FT4` |
| data mode elsewhere | underlying `USB` or `LSB` |

CAT identifies the radio's demodulation mode, not the application creating the
audio. Therefore FT8/FT4 identification uses their conventional dial
frequencies. Other digital applications continue to log their exact mode from
their normal QSO/ADIF message path.

## Read-only connection test

On macOS or Linux, this requests only frequency and mode:

```bash
printf 'f\nm\n' | nc 127.0.0.1 4532
```

A healthy response contains a frequency in hertz, followed by the mode and
passband. The test does not tune the radio or key PTT.

## Troubleshooting

- **Connection refused:** xCAT is not running, RigCtlD is not enabled, or the
  configured host/port is wrong.
- **Connected but no data:** verify xCAT says **Slice present** and is attached
  to the intended station.
- **CAT port works elsewhere but Flex-WavelogGate fails:** check that you used
  RigCtlD `4532`, not CAT `5002`.
- **Wavelog shows a stale radio:** open **Hardware Interfaces** and select the
  Flex-WavelogGate radio as the default.
- **Mode shows USB instead of a digital protocol:** the radio is in a generic
  data mode away from a recognized FT8/FT4 dial frequency. This is deliberate;
  CAT cannot safely guess the application protocol.
- **Remote host fails:** allow the RigCtlD port through the host firewall and
  confirm xCAT is listening on an address reachable from the other computer.

Never expose RigCtlD directly to the public internet. Keep it on a trusted LAN
or use a private VPN.
