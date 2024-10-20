# strux
ECS implementation in Go


# Strux: A High-Performance Entity Component System in Go

[![MIT License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/AMMPTT/strux)](https://goreportcard.com/report/github.com/AMMPTT/strux)
[![GoDoc](https://godoc.org/github.com/AMMPTT/strux?status.svg)](https://godoc.org/github.com/AMMPTT/strux)

## Table of Contents

1. [Introduction](#introduction)
2. [Features](#features)
3. [Installation](#installation)
4. [Usage](#usage)
5. [Architecture](#architecture)
6. [Performance Considerations](#performance-considerations)
7. [Contributing](#contributing)
8. [License](#license)
9. [Citation](#citation)

## Introduction

Strux is a high-performance Entity Component System (ECS) implementation in Go, designed for building efficient and scalable simulation systems. This library provides a flexible and powerful framework for managing entities, components, and systems in a cache-friendly and concurrent manner.

The ECS architectural pattern is widely used in game development and complex simulations due to its ability to efficiently manage large numbers of objects with varying behaviors. Strux aims to bring these benefits to Go developers working on performance-critical applications.

## Features

- **Efficient Component Storage**: Utilizes contiguous memory storage for components, improving cache locality and iteration performance.
- **Flexible Entity Management**: Lightweight entity representation with support for a large number of entities.
- **Archetype-based Querying**: Fast querying of entities based on their component compositions.
- **Concurrent System Execution**: Support for multi-threaded system updates to leverage multi-core processors.
- **Event System**: Decoupled communication between systems and components through a flexible event system.
- **Type-safe Components**: Compile-time type checking for components using Go's interface system.
- **Serialization Support**: Built-in methods for saving and loading the world state.

## Installation

To install Strux, use the following command:

```bash
git clone https://github.com/AMMPTT/strux.git
```

## Usage

Here's a quick example of how to use Strux:

```go
package main

import (
    "fmt"
    "github.com/AMMPTT/strux/ecs"
    "github.com/AMMPTT/strux/components"
)

type PositionComponent struct {
    X, Y float32
}

func (p *PositionComponent) IsComponentData() {}

type VelocityComponent struct {
    VX, VY float32
}

func (v *VelocityComponent) IsComponentData() {}

type MovementSystem struct{}

func (m *MovementSystem) Update(dt float32) {
    // Implementation of movement system
}

func main() {
    world := ecs.NewWorld()
    
    // Create an entity
    entity := world.CreateEntity()
    
    // Add components to the entity
    world.AddComponent(entity, &PositionComponent{X: 0, Y: 0})
    world.AddComponent(entity, &VelocityComponent{VX: 1, VY: 1})
    
    // Add a system to the world
    world.AddSystem(&MovementSystem{})
    
    // Run the simulation
    for i := 0; i < 100; i++ {
        world.Update(0.016) // 60 FPS
    }
}
```

For more detailed usage examples, please refer to the [examples](./examples) directory.

## Architecture

Strux follows a classic Entity Component System architecture with some optimizations for performance:

1. **Entities**: Simple integer identifiers representing game objects or elements in the simulation.
2. **Components**: Pure data structures that can be attached to entities.
3. **Systems**: Logic that operates on entities with specific component combinations.
4. **World**: The container that manages entities, components, and systems.

Key design decisions include:

- Use of contiguous arrays for component storage to improve cache locality.
- Archetype-based entity management for efficient querying and iteration.
- Concurrent system updates with mutex-based thread safety.
- Event system for decoupled communication between systems and components.

For a more detailed architectural overview, please refer to the [ARCHITECTURE.md](./ARCHITECTURE.md) file.

## Performance Considerations

Strux is designed with performance in mind, making several tradeoffs to optimize for common use cases in simulation and game development:

1. **Data Locality**: Components of the same type are stored together in memory, improving cache utilization.
2. **Minimal Indirection**: The design minimizes pointer usage to reduce cache misses.
3. **Concurrent System Execution**: Systems can be updated concurrently, leveraging multi-core processors.
4. **Efficient Querying**: The archetype-based approach allows for fast querying of entities with specific component combinations.

However, users should be aware of potential performance implications:

- The use of Go's garbage collector may introduce occasional pauses in large-scale simulations.
- Mutex-based concurrency may lead to contention in highly parallel scenarios.
- Dynamic component addition/removal can be slower compared to static component compositions.

For performance-critical applications, users may need to implement custom memory management or explore alternative concurrency models.

## Contributing

We welcome contributions to Strux! Please see our [CONTRIBUTING.md](./CONTRIBUTING.md) file for details on how to contribute, our code of conduct, and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.

## Citation

If you use Strux in your research or wish to refer to it in your publications, please use the following BibTeX entry:

```bibtex
@software{strux,
  author = {{Strux Contributors}},
  title = {Strux: A High-Performance Entity Component System in Go},
  year = {2023},
  url = {https://github.com/AMMPTT/strux},
  version = {1.0.0}
}
```

## Self Review
Below is a lightweight self review going over some of the "theoretical" implications of tradeoffs involved should help provide insight as to potential tradeoffs during the design decisions. One particular thing of note, is this authoring approach kept close to mind, that archetypes that contain components, are not very commonly going to be constantly adding and removing components, but rather any changes in archetype state will be due to data changes within that archetype. That isn't to say there wont be cases where adding or removing a component makes sense, but we have found that establishing the data models prior to implementing the simulation aspects leads to an easier to manage simulation model, both from a project standpoint and also from a resource managament standpoint as its easier to track responsiblity for data that doesnt fundamentally change in size. 


---

# Entity Component System (ECS) Implementation Review

## Overview

This ECS implementation in Go is designed to support simulation systems by providing a flexible and efficient framework for managing entities, components, and systems. The design aims to balance performance, flexibility, and ease of use.

## Key Components

1. World
2. Entity
3. Component
4. System
5. EntityManager
6. EventManager

## Design Implications and Tradeoffs

### 1. Component Storage

The implementation uses a component array approach, where components of the same type are stored together in contiguous memory.

**Pros:**
- Improves cache locality for systems that process components of the same type
- Allows for efficient iteration over components of a specific type

**Cons:**
- Removing components can be slower due to the need to maintain contiguous memory
- May lead to fragmentation if entities frequently add/remove components

**Tradeoff:** The choice of using component arrays prioritizes performance for systems that process many entities with the same component type, at the cost of some flexibility in dynamic component management.

### 2. Entity Representation

Entities are represented as simple uint32 values.

**Pros:**
- Lightweight and easy to pass around
- Allows for a large number of entities 

**Cons:**
- No built-in mechanism for entity versioning or recycling

**Tradeoff:** This simple entity representation keeps the system lightweight but may require additional management for entity recycling in long-running simulations.

### 3. Component Interface

Components are defined using an interface with a marker method (`IsComponentData()`).

**Pros:**
- Allows for type safety and compiler checks
- Provides flexibility in component implementation

**Cons:**
- Introduces a small runtime overhead due to interface method calls
- May lead to some use of reflection, which can impact performance

**Tradeoff:** The use of an interface for components provides type safety and flexibility at the cost of some runtime performance.

### 4. Archetype-based Entity Management

The EntityManager uses an archetype-based approach for managing entities and their components.

**Pros:**
- Efficient querying of entities with specific component combinations
- Supports cache-friendly iteration over entities with the same component types

**Cons:**
- Increased memory usage due to storing multiple copies of entity-component associations
- More complex implementation compared to simpler approaches

**Tradeoff:** The archetype-based approach improves query and iteration performance at the cost of increased memory usage and implementation complexity.

### 5. Concurrency Model

The implementation uses mutex locks for thread safety and allows for concurrent system updates.

**Pros:**
- Supports multi-threaded simulations
- Prevents data races when accessing shared data

**Cons:**
- Potential for lock contention in highly concurrent scenarios
- May limit scalability for systems with frequent component access

**Tradeoff:** The use of mutex locks provides thread safety but may impact performance in highly concurrent scenarios. Alternative concurrency models (e.g., lock-free data structures) could be explored for better scalability.

### 6. Memory Management

The implementation uses Go's built-in garbage collector and does not implement custom memory management.

**Pros:**
- Simplifies memory management
- Leverages Go's efficient garbage collector

**Cons:**
- Less control over memory allocation and deallocation
- May lead to occasional GC pauses in large simulations

**Tradeoff:** Relying on Go's garbage collector simplifies the implementation but may introduce some performance variability in large-scale simulations.

### 7. Event System

The implementation includes an event system for communication between systems and components.

**Pros:**
- Allows for decoupled communication between parts of the simulation
- Supports reactive programming patterns

**Cons:**
- May introduce overhead for event dispatch and handling
- Can make the flow of the simulation harder to reason about if overused

**Tradeoff:** The event system provides flexibility in system communication at the cost of potential performance overhead and increased complexity in understanding the simulation flow.

## Performance Considerations

1. **Data Locality:** The use of component arrays improves data locality, which can lead to better cache utilization and improved performance for systems processing many entities with the same component type.

2. **Pointer Usage:** The implementation generally avoids excessive use of pointers, which helps reduce indirection and improves cache efficiency. However, some uses of pointers (e.g., in the ComponentArray) may introduce some level of indirection.

3. **Allocation Patterns:** The implementation relies on Go's built-in allocation mechanisms. For performance-critical simulations, custom allocation strategies (e.g., object pools, custom allocators) could be considered to reduce GC pressure.

4. **Polymorphism:** The design generally avoids heavy use of polymorphism, favoring a more data-oriented approach with component arrays. This aligns well with the goals of an ECS and can lead to better performance in many scenarios.

## Suggestions for Improvement

1. **Entity Recycling:** Implement an entity recycling mechanism to avoid potential issues with entity ID exhaustion in long-running simulations.

2. **Custom Allocators:** For performance-critical simulations, consider implementing custom memory allocators or object pools for frequently created/destroyed components.

3. **Lock-Free Data Structures:** Explore the use of lock-free data structures for highly concurrent scenarios to reduce lock contention.

4. **Component Storage Optimization:** Consider implementing a hybrid approach that combines the benefits of sparse sets and dense arrays for component storage, potentially improving both iteration performance and memory usage.

5. **Code Generation:** Implement a code generation tool to create type-safe, performance-optimized systems and queries, reducing the reliance on reflection and interface method calls.

## Conclusion

This ECS implementation in Go provides a solid foundation for building simulation systems. It makes reasonable tradeoffs between performance, flexibility, and ease of use. The design leans towards a data-oriented approach, which aligns well with the principles of ECS and can lead to good performance in many scenarios.

However, for extremely performance-critical or large-scale simulations, further optimizations and custom implementations of certain subsystems (e.g., memory management, concurrency control) may be necessary. The current implementation provides a good balance for a wide range of simulation scenarios while leaving room for domain-specific optimizations.

---

For any questions, issues, or suggestions, please open an issue on the GitHub repository or contact the maintainers directly.
