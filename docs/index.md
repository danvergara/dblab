# Welcome to **dblab**

![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/gopher-dblab.png){ width="300" : .center }

Cross-platform, zero dependencies, terminal based UI application for your Data Bases.  

![integration tests](https://github.com/danvergara/dblab/actions/workflows/ci.yaml/badge.svg)  ![unit tests](https://github.com/danvergara/dblab/actions/workflows/test.yaml/badge.svg)  [![Release](https://img.shields.io/github/release/danvergara/dblab.svg?label=Release)](https://github.com/danvergara/dblab/releases)

---

**Documentation**: <a href="https://dblab.app" target="_blank">https://dblab.app</a>

**Source Code**: <a href="https://github.com/danvergara/dblab" target="_blank">https://github.com/danvergara/dblab</a>

---

## Overview

dblab is a fast and lightweight interactive terminal based UI application for PostgreSQL, MySQL and SQLite3,
written in Go and works on OSX and Linux machines. Main idea behind using Go for backend development
is to utilize ability of the compiler to produce zero-dependency binaries for
multiple platforms. dblab was created as an attempt to build very simple and portable
application to work with local or remote PostgreSQL/MySQL/SQLite3/Oracle/SQL Server databases.
  

## Features

  * Cross-platform support OSX/Linux 32/64-bit  
  * Simple installation (distributed as a single binary)  
  * Zero dependencies.  

## Installation

{--if you need to work with SQLite3, install the CGO enabled binary using the proper bash script listed below.--}

{==

The above comment is deprecated and CGO is not needed anymore.   
There will be a single binary capable to deal with all supported clients.
 
==}

### Homebrew installation

It works with Linux too.

```{ .sh .copy }
brew install --cask danvergara/tools/dblab
```

Or

```{ .sh .copy }
brew tap danvergara/tools
brew install --cask dblab
```

### Manual Binary Installation
The binaries are compatible with Linux and OSX.  
You can manually download and install the binary release from [the release page](https://github.com/danvergara/dblab/releases).

### Automated installation/update
> Don't forget to always verify what you're piping into bash

Install the binary using our bash script:

```{ .sh .copy }
curl https://raw.githubusercontent.com/danvergara/dblab/master/scripts/install_update_linux.sh | bash
```

