
# DeepDiveContainers

**DeepDiveContainers** is a project written in Go that delves into the inner workings of containers, showcasing how they utilize Linux kernel features such as namespaces, cgroups, and syscalls. This repository serves as an educational tool for understanding containerization from scratch.

---

## Table of Contents

- [Overview](#overview)
- [Key Learnings](#key-learnings)
- [Syscalls Used](#syscalls-used)
- [tmp-rootfs Directory](#tmp-rootfs-directory)
- [TUN/TAP Network Interfaces](#tuntap-network-interfaces)
- [Features Implemented](#features-implemented)
- [Technologies Used](#technologies-used)
- [Getting Started](#getting-started)
- [Project Structure](#project-structure)
- [License](#license)

---

## Overview

This project explores how containers are built under the hood by leveraging Linux kernel features. It demonstrates how to isolate processes, manage resource limits, and create minimal runtime environments by employing syscalls such as `pivot_root`, `mount`, and `clone`. The implementation is written entirely in Go, providing a practical way to learn container fundamentals.

---

## Key Learnings

Through this project, I explored:

1. **Namespaces**:
   - **UTS Namespace**: Allows setting a hostname for the container using `syscall.Sethostname`.
   - **PID Namespace**: Isolates process IDs within the container.
   - **Network Namespace**: Configures container-specific networking.

2. **Syscalls**:
   - `pivot_root`: Changes the root filesystem of the container to isolate it from the host.
   - `mount`: Mounts filesystems like `/proc` inside the container for process visibility.
   - `clone`: Creates new processes with specific namespace flags.

3. **Root Filesystem**:
   - Learned how to initialize a root filesystem for containers and use `pivot_root` to isolate it.

4. **Process Isolation**:
   - Used the `exec.Command` and `syscall.SysProcAttr` to run commands in an isolated environment.

5. **Networking**:
   - Explored basic networking setups using tools like `slirp4netns` to enable container-specific networking.

---

## Syscalls Used

1. **`syscall.Clone`**:
   - **Description**: Creates a new process with specific namespaces for isolation. The `Cloneflags` parameter enables the creation of new namespaces (e.g., UTS, PID, IPC, NET, and USER namespaces).
   - **Purpose in the Project**:
     - Isolates the container's hostname (UTS namespace).
     - Creates a separate process tree (PID namespace).
     - Sets up isolated IPC and networking (IPC and NET namespaces).
     - Enables rootless operation by mapping UIDs and GIDs (USER namespace).

2. **`syscall.PivotRoot`**:
   - **Description**: Changes the root filesystem of the current process to a new directory. This is crucial for isolating the container's filesystem from the host.
   - **Purpose in the Project**:
     - Ensures that the container only sees the files in the `tmp-rootfs` directory (a BusyBox image).
     - Provides filesystem isolation by pivoting to a minimal root filesystem.

3. **`syscall.Mount`**:
   - **Description**: Mounts a filesystem to a specified target. This syscall is used to bind mount directories or mount virtual filesystems like `proc`.
   - **Purpose in the Project**:
     - Mounts the `proc` filesystem inside the container to provide process-related information.
     - Uses bind mounts to set up directories for the container.

4. **`syscall.Sethostname`**:
   - **Description**: Sets the hostname of the current UTS namespace.
   - **Purpose in the Project**:
     - Assigns a unique hostname (`"container"`) to the container, separating it from the host.

5. **`syscall.Unmount`**:
   - **Description**: Detaches a mounted filesystem.
   - **Purpose in the Project**:
     - Unmounts the old root directory after the `pivot_root` operation to complete the filesystem isolation.

6. **`syscall.UidMappings` and `syscall.GidMappings`**:
   - **Description**: Maps container user IDs (UIDs) and group IDs (GIDs) to host IDs for rootless operation.
   - **Purpose in the Project**:
     - Allows the container to run as a non-root user while maintaining isolation.

---

## tmp-rootfs Directory

- **What is `tmp-rootfs`?**
  - The `tmp-rootfs` directory contains a minimal BusyBox image that is used as the root filesystem for the container.
  - BusyBox is a compact set of Unix utilities that provides the basic tools needed to run a minimal environment.

- **Role in the Project**:
  - Serves as the isolated root filesystem for the container.
  - Ensures that the container only has access to the files and tools within the `tmp-rootfs` directory, providing filesystem isolation from the host.

---

## TUN/TAP Network Interfaces

- **What Are TUN/TAP Interfaces?**
  - TUN (network TUNnel) and TAP (network TAP) are virtual network kernel interfaces.
  - TUN interfaces operate at the IP level (layer 3), while TAP interfaces work at the Ethernet frame level (layer 2).

- **Role in Rootless Containers**:
  - TUN/TAP interfaces are used to create virtual network interfaces for containers, enabling network isolation and connectivity.
  - In the project, `slirp4netns` is used to configure networking for rootless containers. It leverages TUN interfaces to forward network packets between the container and the host.

- **How It Works**:
  1. A TAP interface is created for the container.
  2. The container's network namespace is connected to the host's network through the TAP interface.
  3. `slirp4netns` or similar tools emulate a network stack for the container, allowing it to access external networks without requiring root privileges.

- **Benefits for Rootless Containers**:
  - Provides a secure and isolated networking environment.
  - Allows containers to communicate with external networks without requiring elevated permissions on the host.

---

## Features Implemented

- **Namespace Isolation**:
  - Isolates UTS, PID, IPC, Network, and Mount namespaces using `syscall.CLONE_*` flags.

- **Filesystem Isolation**:
  - Implements `pivot_root` to sandbox the container's filesystem.
  - Mounts `/proc` inside the container for process visibility.

- **Custom Hostname**:
  - Sets a custom hostname for the container using syscalls.

- **Process Management**:
  - Runs commands inside the container using an isolated process.

- **Networking**:
  - Configures container-specific networking using external tools.

---

## Technologies Used

- **Language**:
  - Go (100%)

- **Linux Kernel Features**:
  - Namespaces
  - Syscalls (`pivot_root`, `mount`, `clone`, etc.)

---

## Getting Started

### Prerequisites

- Linux-based system with kernel version 4.8 or higher.
- Go (1.17 or higher).

### Running the Project

1. Clone the repository:
   ```bash
   git clone https://github.com/Ali-Afifi/DeepDiveContainers.git
   cd DeepDiveContainers
   ```

2. Build the project:
   ```bash
   go build -o container
   ```

3. Run the container:
   ```bash
   ./container run <command> <parameters>
   ```

Example:
   ```bash
   ./container run /bin/bash
   ```

---

## Project Structure

```
DeepDiveContainers/
├── main.go       # Entry point for the container implementation.
├── rootfs.go     # Utilities for managing the container's root filesystem.
├── README.md     # Project documentation.
└── utils/
    ├── syscall_helpers.go  # Helper functions for syscalls.
    └── ...
```

---

## License

This project is licensed under the MIT License. See the LICENSE file for details.


Here’s the formatted "References" section to be added to the README file:

---

## References

1. [Linux Clone System Call](https://medium.com/brundas-tech-notes/linux-clone-system-call-a45f717fa697) - A detailed explanation of the Linux `clone` system call and its use in creating isolated processes.
2. [Linux Namespaces](https://medium.com/@teddyking/linux-namespaces-850489d3ccf) - An overview of Linux namespaces and how they provide isolation for containers.
3. [Pivot Root vs Chroot for Containers](https://tbhaxor.com/pivot-root-vs-chroot-for-containers/) - A comparison of `pivot_root` and `chroot` in the context of containerization.
4. [Introduction to Linux Interfaces for Virtual Networking](https://developers.redhat.com/blog/2018/10/22/introduction-to-linux-interfaces-for-virtual-networking#) - A guide to TUN/TAP interfaces and other Linux networking tools for virtualized environments.
5. [Building a Container From Scratch in Go (YouTube)](https://www.youtube.com/watch?v=8fi7uSYlOdc) - A video tutorial on building a container from scratch using Go.

---
