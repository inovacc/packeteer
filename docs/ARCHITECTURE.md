# Architecture

## System Overview

```mermaid
flowchart TB
    subgraph CLI["CLI Layer (cmd/)"]
        ROOT[packeteer root]
        SERVE[serve command]
        VERSION[version command]
        AICTX[aicontext command]
    end

    subgraph MCP["MCP Server (internal/server/)"]
        SERVER[MCP Server<br/>stdio transport]
    end

    subgraph TOOLS["Tool Handlers (internal/tools/)"]
        T1[list_interfaces]
        T2[capture_packets]
        T3[read_pcap]
        T4[extract_fields]
        T5[get_statistics]
        T6[get_capture_info]
        T7[filter_pcap]
        T8[merge_pcaps]
        T9[list_protocols]
        T10[decode_packet]
    end

    subgraph SAFETY["Safety Layer (internal/safety/)"]
        VAL[Path Validation]
        SAN[Filter Sanitization]
        CLAMP[Timeout/Count Clamping]
    end

    subgraph EXEC["Executor (internal/executor/)"]
        IFACE[CommandExecutor<br/>interface]
        REAL[RealExecutor]
        MOCK[MockExecutor]
    end

    subgraph OUTPUT["Output (internal/output/)"]
        TRUNC[Truncation]
        FMT[Formatting]
    end

    subgraph EXT["External Tools"]
        TSHARK[tshark]
        CAPINFOS[capinfos]
        EDITCAP[editcap]
        MERGECAP[mergecap]
    end

    ROOT --> SERVE
    ROOT --> VERSION
    ROOT --> AICTX
    SERVE --> SERVER

    SERVER --> T1 & T2 & T3 & T4 & T5
    SERVER --> T6 & T7 & T8 & T9 & T10

    T1 & T2 & T3 & T4 & T5 --> SAFETY
    T6 & T7 & T8 & T9 & T10 --> SAFETY
    T1 & T2 & T3 & T4 & T5 --> IFACE
    T6 & T7 & T8 & T9 & T10 --> IFACE
    T1 & T2 & T3 & T4 & T5 --> TRUNC
    T6 & T7 & T8 & T9 & T10 --> TRUNC

    IFACE --> REAL
    IFACE --> MOCK
    REAL --> TSHARK & CAPINFOS & EDITCAP & MERGECAP
```

## MCP Request Flow

```mermaid
sequenceDiagram
    participant AI as AI Assistant
    participant MCP as MCP Server<br/>(stdio)
    participant TH as Tool Handler
    participant SAF as Safety Layer
    participant EXEC as CommandExecutor
    participant TS as tshark/capinfos

    AI->>MCP: CallTool (JSON-RPC)
    MCP->>TH: Dispatch to handler<br/>(typed input)

    TH->>SAF: Validate inputs
    alt Validation fails
        SAF-->>TH: Error
        TH-->>MCP: IsError=true
        MCP-->>AI: Error response
    else Validation passes
        SAF-->>TH: OK
        TH->>EXEC: Execute(binary, args)
        EXEC->>TS: exec.CommandContext
        TS-->>EXEC: stdout, stderr
        EXEC-->>TH: output bytes
        TH->>TH: Truncate if needed
        TH-->>MCP: CallToolResult
        MCP-->>AI: JSON response
    end
```

## Server Lifecycle

```mermaid
sequenceDiagram
    participant CLI as Cobra CLI
    participant SRV as server.New()
    participant MCP as mcp.Server
    participant STDIO as StdioTransport

    CLI->>SRV: Create executor + logger
    SRV->>MCP: mcp.NewServer()
    SRV->>MCP: AddTool() x10
    SRV-->>CLI: *mcp.Server

    CLI->>MCP: server.Run(ctx, stdio)
    MCP->>STDIO: Listen stdin

    loop Process requests
        STDIO->>MCP: JSON-RPC request
        MCP->>MCP: Route to tool handler
        MCP->>STDIO: JSON-RPC response
    end

    Note over CLI,STDIO: Context cancellation or EOF stops server
```

## Key Design Decisions

- **Hexagonal Architecture:** `CommandExecutor` interface decouples tool handlers from real CLI execution, enabling full unit testing with `MockExecutor`
- **Safety First:** All inputs validated before any command execution — path traversal blocked, filters sanitized, timeouts clamped
- **Stdio Transport:** MCP JSON-RPC over stdin/stdout; all logging to stderr to avoid protocol corruption
- **Output Truncation:** Large captures truncated to 512KB with metadata about omitted data
