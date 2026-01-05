# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

**Communication Language**: Please communicate in Chinese (中文) when working in this repository.

## Overview

This is a Go generic pool implementation (`github.com/ndsky1003/pool`) providing an adaptive ring buffer pool for object reuse and memory optimization. The pool automatically scales up/down based on usage patterns.

## Architecture

The codebase consists of three main files:

- **pool.go** - Core `AdaptiveRingPool[T]` implementation:
  - Ring buffer with head/tail pointers for object storage
  - Lock-based Get/Put operations with atomic counters for statistics
  - Auto-scaling logic that adjusts capacity based on hit rate
  - Uses padding to reduce false sharing on cache lines (64-byte alignment)

- **option.go** - Configuration system using functional options pattern:
  - `DefaultOptions()` provides sensible defaults (Min: 32, Max: 512, scale factors: 1.2/0.8)
  - `With*()` functions for custom configuration

- **api.go** - Global registry API for type-based pool management:
  - `Regist[T]()` - Register a pool for type T
  - `Get[T]()` - Get object from registered pool
  - `Unregist[T]()` - Unregister a pool
  - Uses `sync.Map` for concurrent-safe type-keyed storage

## Key Design Decisions

1. **Auto-scaling triggers**: Scaling only happens in `Put()` to minimize performance impact on `Get()`
2. **Hit rate thresholds**: >0.8 triggers scale-up, <0.2 triggers scale-down
3. **Ring buffer resize**: Handles wrapped data correctly when copying to new buffer
4. **Memory safety**: Maximum capacity prevents unbounded growth, minimum capacity preserves baseline

## Development

### Building and Testing

```bash
# Run tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with race detector
go test -race ./...

# Build the module
go build ./...
```

### Module Information

- **Module**: `github.com/ndsky1003/pool`
- **Go version**: 1.23.0
- **No external dependencies**

### Code Style

- The codebase uses Chinese comments for explaining complex logic
- Field padding is used explicitly for cache line optimization (64-byte alignment)
- Atomic operations are used for high-frequency counters (hitCount, getCount)
