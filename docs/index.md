# Welcome to **dblab**

![dblab](https://github.com/danvergara/dblab/blob/main/assets/gopher-dblab.png){ width="300" }



## Overview

**dblab** is a fast and lightweight interactive terminal based UI application for PostgreSQL, MySQL and SQLite3, written in Go and works on OSX, Linux and Windows machines.   

The main idea behind using Go for backend development is to utilize the ability of the compiler to produce zero-dependency binaries for multiple platforms. 
  
dblab was created as an attempt to build a very simple and portable application to interact with local or remote PostgreSQL/MySQL/SQLite3 databases.  
  
The key features are:

  * Cross-platform support OSX/Linux/Windows 32/64-bit  
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

```
$ brew install danvergara/tools/dblab
```

Or

```
$ brew tap danvergara/tools
$ brew install dblab
```

### Manual Binary Installation
The binaries are compatible with Linux, OSX and Windows.  
You can manually download and install the binary release from [the release page](https://github.com/danvergara/dblab/releases).

### Automated installation/update
> Don't forget to always verify what you're piping into bash

Install the binary using our bash script:

```sh
curl https://raw.githubusercontent.com/danvergara/dblab/master/scripts/install_update_linux.sh | bash
```

