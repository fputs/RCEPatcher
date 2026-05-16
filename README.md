# RCEPatcher

Patches the WoW WOTLK 3.3.5 client to remove the remote code execution exploit. Any server you connect to with an
unpatched client can remotely execute code on your computer.

This is a cross-platform Go port of the original C# [RCEPatcher](https://github.com/stoneharry/RCEPatcher) that works
on Linux, macOS, or Windows.

## Usage

RCEPatcher will inform you if `Wow.exe` is already patched. If not, it creates a backup called `Wow.exe.backup` and then
patches `Wow.exe` in place.

Either download the executable for your OS and architecture from [the release page](https://git.burkey.co/eburk/RCEPatcher/releases/latest)
or build yourself, then run the program. If you download the binary from the releases page, be aware that none of the
binaries are codesigned or notarized and you may get security popups on your platform.
```shell
# OPTIONAL: build yourself
go build

# Run RCEPatcher
./RCEPatcher /path/to/Wow.exe
```