# Welcome to **dblab**

![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/gopher-dblab.png){ width="300" : .center }

Cross-platform, zero dependencies, terminal-based UI application for your databases.  

![integration tests](https://github.com/danvergara/dblab/actions/workflows/ci.yaml/badge.svg)  ![unit tests](https://github.com/danvergara/dblab/actions/workflows/test.yaml/badge.svg)  [![Release](https://img.shields.io/github/release/danvergara/dblab.svg?label=Release)](https://github.com/danvergara/dblab/releases)

---

**Documentation**: <a href="https://dblab.app" target="_blank">https://dblab.app</a>

**Source Code**: <a href="https://github.com/danvergara/dblab" target="_blank">https://github.com/danvergara/dblab</a>

---

## Overview

dblab is a fast and lightweight interactive terminal-based UI application for PostgreSQL, MySQL, and SQLite3,
written in Go and works on macOS, Linux, and Windows machines. The main idea behind using Go for backend development
is to utilize the ability of the compiler to produce zero-dependency binaries for
multiple platforms. dblab was created as an attempt to build a very simple and portable
application to work with local or remote PostgreSQL/MySQL/SQLite3/Oracle/SQL Server databases.
  

<div style="position:relative;padding-bottom:56.25%;height:0;overflow:hidden;max-width:100%;margin:1.5rem 0;border-radius:10px;border:1px solid rgba(128,128,128,0.2);box-shadow:0 4px 20px rgba(0,0,0,0.15);">
  <iframe src="https://bisque.cloud/p/github/danvergara-dblab" title="dblab — narrated walkthrough" loading="lazy"
    style="position:absolute;top:0;left:0;width:100%;height:100%;border:0;"
    allow="autoplay; fullscreen; encrypted-media" allowfullscreen></iframe>
</div>

## Features

  * Cross-platform support for macOS/Linux/Windows (32/64-bit)  
  * Simple installation (distributed as a single binary)  
  * Zero dependencies.  
  * Vim-style query editor (normal and insert modes, line-oriented editing commands).  
  * Multi-query execution: write multiple SQL statements separated by `;` and run them concurrently with results in separate tabs.  
  * Connection profiles with secure credential storage in the OS keyring.  
  * Query history: executed queries are automatically saved and can be browsed or re-used from a searchable list.  
  * Read-only mode: use `--readonly` to prevent accidental writes by forcing the database session into read-only mode (PostgreSQL, MySQL, SQLite, Oracle, and SQL Server).  

## Installation

{--if you need to work with SQLite3, install the CGO enabled binary using the proper bash script listed below.--}

{==

The above comment is deprecated and CGO is not needed anymore.   
There will be a single binary capable of dealing with all supported clients.
 
==}

### Homebrew installation

It works with Linux, too.

```{ .sh .copy }
brew install --cask danvergara/tools/dblab
```

Or

```{ .sh .copy }
brew tap danvergara/tools
brew install --cask dblab
```

### Manual Binary Installation
The binaries are compatible with Linux, macOS, and Windows.  
You can manually download and install the binary release from [the release page](https://github.com/danvergara/dblab/releases).

### Automated installation/update
> Don't forget to always verify what you're piping into bash

Install the binary using our bash script:

```{ .sh .copy }
curl https://raw.githubusercontent.com/danvergara/dblab/master/scripts/install_update_linux.sh | bash
```

