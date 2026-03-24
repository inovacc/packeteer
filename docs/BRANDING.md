# TShark AI Branding Names

## Project Identity

- **Current Name:** tshark_ai
- **Purpose:** Go MCP server that wraps Wireshark CLI tools (tshark, capinfos, editcap, mergecap) for AI-driven packet capture and network analysis
- **Target Audience:** Security engineers, network analysts, DevOps, and AI/LLM tool users
- **Domain:** Network forensics, packet analysis, MCP protocol

---

## Project Name Candidates

| # | Name | Rationale |
|---|------|-----------|
| 1 | **tshark_ai** | Current name — descriptive but underscore is non-idiomatic for Go CLIs |
| 2 | **Sharkline** | Shark (Wireshark) + line (command-line) — clean, memorable |
| 3 | **Packeteer** | Evokes packet mastery; professional, action-oriented |
| 4 | **Wiretap** | Classic network analysis metaphor — short, punchy, evocative |
| 5 | **Finsniff** | Fin (shark fin) + sniff (packet sniffing) — playful, unique |
| 6 | **Netjaw** | Net (network) + jaw (shark anatomy) — aggressive, memorable |
| 7 | **Dissectr** | From "dissector" (Wireshark term) — technical, brandable |
| 8 | **SharkCast** | Shark + cast (to capture/throw a net) — action metaphor |
| 9 | **Pcapture** | Pcap + capture — direct domain reference, developer-friendly |
| 10 | **Dorsal** | Shark's dorsal fin — elegant, abstract, brandable |
| 11 | **Carcharias** | Genus of great white shark — scientific, distinctive |
| 12 | **Jawline** | Jaw (shark) + line (CLI) — dual meaning with sleek connotation |

**Recommended:** **Packeteer** — professional, action-oriented, instantly communicates packet mastery without being tied to "tshark" or "Wireshark" trademarks.

---

## Feature Names

| Feature | Current Name | Branded Name Options |
|---------|-------------|---------------------|
| Live packet capture | `capture_packets` | `snare`, `trapline`, `livecatch` |
| Pcap file analysis | `read_pcap` | `dissect`, `unravel`, `inspect` |
| Field extraction | `extract_fields` | `pluck`, `harvest`, `distill` |
| Protocol statistics | `get_statistics` | `census`, `tally`, `survey` |
| Interface listing | `list_interfaces` | `scan`, `probe`, `enumerate` |
| Capture file info | `get_capture_info` | `manifest`, `profile`, `dossier` |
| Pcap filtering | `filter_pcap` | `sieve`, `refine`, `winnow` |
| Pcap merging | `merge_pcaps` | `fuse`, `splice`, `weave` |
| Protocol listing | `list_protocols` | `codex`, `registry`, `catalog` |
| Packet decode | `decode_packet` | `reveal`, `unfold`, `expose` |

---

## Component Names

| Component | Branded Name Options |
|-----------|---------------------|
| Command executor (CLI wrapper) | `helm`, `rigging`, `harness` |
| Safety/validation layer | `guardrail`, `bulkhead`, `reef` |
| Output truncation | `trimmer`, `breaker`, `spillway` |
| MCP server core | `bridge`, `conning`, `deck` |
| Filter sanitizer | `strainer`, `sluice`, `grate` |

---

## Taglines

| # | Tagline | Style |
|---|---------|-------|
| 1 | **Packet intelligence, on demand.** | Short & punchy |
| 2 | **Wireshark for your AI.** | Descriptive |
| 3 | **Let your LLM read the wire.** | Aspirational |
| 4 | **MCP-native packet capture and analysis.** | Technical |
| 5 | **Sniff, dissect, understand.** | Action-driven |
| 6 | **Network forensics meets model context.** | Domain bridge |
| 7 | **Every packet tells a story.** | Narrative |
| 8 | **Deep packet intelligence for AI assistants.** | Professional |

---

## CLI Branding Themes

### Theme 1: Naval / Maritime
```
capture_packets  → intercept
read_pcap        → chart
extract_fields   → salvage
get_statistics   → survey
list_interfaces  → fleet
merge_pcaps      → convoy
filter_pcap      → trawl
decode_packet    → fathom
```

### Theme 2: Forensics / Investigation
```
capture_packets  → wiretap
read_pcap        → examine
extract_fields   → extract
get_statistics   → profile
list_interfaces  → canvas
merge_pcaps      → consolidate
filter_pcap      → isolate
decode_packet    → reconstruct
```

### Theme 3: Minimal / Verb-only
```
capture_packets  → capture
read_pcap        → read
extract_fields   → extract
get_statistics   → stats
list_interfaces  → interfaces
merge_pcaps      → merge
filter_pcap      → filter
decode_packet    → decode
```

---

## Color Palette Suggestions

| Role | Color Name | Hex Code | Rationale |
|------|-----------|----------|-----------|
| Primary | **Deep Ocean** | `#0F4C75` | Dark blue — network, depth, trust |
| Secondary | **Shark Grey** | `#3C4F65` | Steel grey — technical, professional |
| Accent | **Signal Cyan** | `#00D4AA` | Bright teal — data flow, packet highlights |
| Warning | **Alert Coral** | `#E74C3C` | Red-coral — errors, dropped packets |
| Muted | **Wire Slate** | `#8B9DAF` | Soft blue-grey — secondary text, borders |

---

## Logo Concepts

1. **Shark Fin + Packet Wave** — A minimalist shark dorsal fin emerging from a stylized sine wave representing network traffic; conveys packet capture from the data stream.

2. **Hexagonal Shark Eye** — A shark's eye rendered inside a hexagon (echoing hexagonal architecture), with concentric rings suggesting protocol layers and deep packet inspection.

3. **Wire Mesh Jaw** — Abstract shark jaw outline composed of interconnected nodes and edges (network topology), representing the tool's ability to bite into and dissect network data.

4. **Terminal Shark** — A shark silhouette formed entirely from ASCII/terminal characters (`>`, `|`, `/`), referencing the CLI-first nature of the tool — works great as a monochrome favicon.

---

## Icon Generation

```bash
iconforge forge --generate \
  --name packeteer \
  --primary "#0F4C75" \
  --secondary "#3C4F65" \
  --accent "#00D4AA" \
  --output build/icons
```
