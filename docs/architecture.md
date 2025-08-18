# Architecture & Package Structure

This document describes the clean, modular architecture of RxGo, including the package structure and API design choices.

## Overview

RxGo provides a modular, clean API that combines the best of both worlds:

- **Simple Observable API** - Clean, intuitive reactive programming
- **Full Reactive Streams 1.0.4 Compliance** - Built-in backpressure support
- **Modular Package Structure** - Organized into logical subpackages
- **Flexible Imports** - Use what you need, when you need it

## Package Structure

The library is organized into clean, focused subpackages:

### Core Packages

#### `github.com/droxer/RxGo/pkg/observable`
- **Purpose**: Core Observable API and basic reactive programming
- **Key Features**:
  - `Observable[T]` - Main observable type with generics support
  - `Subscriber[T]` - Functional subscriber interface
  - `Just()`, `Range()`, `Create()` - Observable creation functions
  - Operators: `Map()`, `Filter()`, `ObserveOn()`, `Merge()`, `Concat()`, `Take()`, `Skip()`, `Distinct()` - Data transformation functions
  - Context-based cancellation support
  - Type-safe throughout with Go generics

#### `github.com/droxer/RxGo/pkg/scheduler`
- **Purpose**: Advanced threading and execution control
- **Key Features**:
  - `Scheduler` interface - Unified scheduling abstraction
  - **Trampoline** - Immediate execution for lightweight operations
  - **NewThread** - New goroutine for each task
  - **SingleThread** - Sequential processing on dedicated goroutine
  - **Computation** - Fixed thread pool for CPU-bound work
  - **IO** - Cached thread pool for I/O-bound work

#### `github.com/droxer/RxGo/pkg/streams`
- **Purpose**: Full Reactive Streams 1.0.4 compliance
- **Key Features**:
  - `Publisher[T]` - Type-safe data source with demand control
  - `Subscriber[T]` - Complete subscriber interface with lifecycle
  - `Subscription` - Request/cancel control with backpressure
  - `Processor[T,R]` - Transforming publisher that implements both Publisher and Subscriber
  - Processors: `MapProcessor`, `FilterProcessor`, `FlatMapProcessor`, `MergeProcessor`, `ConcatProcessor`, `TakeProcessor`, `SkipProcessor`, `DistinctProcessor` - Data transformation types
  - Full backpressure support with multiple strategies

#### `adapters.go` (in `github.com/droxer/RxGo/pkg/`)
- **Purpose**: Interoperability between Observable and Reactive Streams APIs
- **Key Features**:
  - `ObservablePublisherAdapter[T]` - Convert Observable to Publisher
  - `PublisherToObservableAdapter[T]` - Convert Publisher to Observable

## Design Principles

### 1. Modularity
Each package has a single, well-defined responsibility:
- `observable/` - Core Observable API and basic reactive programming
- `scheduler/` - Execution control and threading models
- `streams/` - Full Reactive Streams 1.0.4 compliance
- `adapters.go` - Interoperability between Observable and Reactive Streams APIs

### 2. Push vs Pull Models
RxGo implements two distinct reactive models to serve different use cases:

**Push Model (`pkg/observable`)**:
- Observable API follows a push-based model
- Producer controls emission rate
- Data is pushed to subscribers as soon as it's available
- No built-in backpressure mechanism
- Simple and intuitive for basic use cases
- Best for scenarios where producers and consumers operate at similar speeds

**Pull Model (`pkg/streams`)**:
- Reactive Streams specification follows a pull-based model with backpressure
- Subscriber controls consumption rate through demand requests
- Publishers must respect subscriber demand (`Request(n)` calls)
- Full backpressure support with multiple strategies
- Production-ready with Reactive Streams 1.0.4 compliance
- Best for scenarios where producers may outpace consumers

### 3. Operators vs Processors

RxGo provides two distinct approaches to data transformation:

**Operators (`pkg/observable/operators.go`)**:
- Function-based transformations for the Observable API
- Composed using function calls: `observable.Map(observable, transformFunc)`
- Work with push-based Observable model
- No built-in backpressure support
- Simple syntax for basic transformations
- Examples: `Map`, `Filter`, `ObserveOn`, `Merge`, `Concat`, `Take`, `Skip`, `Distinct`

**Processors (`pkg/streams/processors.go`)**:
- Type-based transformations for the Reactive Streams API
- Composed using Subscribe method calls: `publisher.Subscribe(ctx, processor)`
- Work with pull-based Publisher/Subscriber model
- Full Reactive Streams 1.0.4 compliance with backpressure
- Implement both Publisher and Subscriber interfaces
- Examples: `MapProcessor`, `FilterProcessor`, `FlatMapProcessor`, `MergeProcessor`, `ConcatProcessor`, `TakeProcessor`, `SkipProcessor`, `DistinctProcessor`

For a detailed explanation of push vs pull models and backpressure strategies, see the [Push vs Pull Models documentation](./push-pull-models.md).

### When to Use Each Model

**Use the Push Model (`pkg/observable`) when**:
- Building simple reactive applications
- Working with predictable data rates
- Prototyping or learning reactive programming
- Data sources that naturally emit at controlled rates

**Use the Pull Model (`pkg/streams`) when**:
- Building production systems with potentially unbounded data streams
- Need to handle producer/consumer speed mismatches
- Working with external data sources (network, file I/O, etc.)
- Requiring Reactive Streams 1.0.4 compliance
- Need for backpressure control to prevent memory issues

### 4. Type Safety
Full generics support throughout the API:
- Type-safe Observable creation and processing
- Compile-time type checking
- No interface{} or reflection usage

### 5. Context Integration
Native Go context support:
- Graceful cancellation using `context.Context`
- Goroutine leak prevention
- Timeout and deadline support

### 6. Performance
Optimized for Go's concurrency model:
- Zero-allocation signal delivery where possible
- Lock-free data structures
- Efficient goroutine management
- Bounded buffers to prevent memory issues
