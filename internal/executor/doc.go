// Package executor provides the CommandExecutor interface and implementations
// for running Wireshark CLI tools (tshark, capinfos, editcap, mergecap).
// RealExecutor calls actual binaries; MockExecutor enables testing without Wireshark installed.
package executor
