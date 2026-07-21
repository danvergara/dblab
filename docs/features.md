The key features are:

  * Cross-platform support for macOS/Linux/Windows (32/64-bit)  
  * Simple installation (distributed as a single binary)  
  * Zero dependencies.  
  * Vim-style query editor (normal and insert modes, line-oriented editing commands).  
  * Multi-query execution: write multiple SQL statements separated by `;` and run them concurrently with results in separate tabs.  
  * Single-query execution: execute only the query on the current cursor line with `ctrl+r`, without running other statements in the editor.  
  * Connection profiles with secure credential storage in the OS keyring.  
  * Query history: executed queries are automatically saved and can be browsed or re-used from a searchable list.  
  * Read-only mode: use `--readonly` to prevent accidental writes by forcing the database session into read-only mode (PostgreSQL, MySQL, SQLite, Oracle, and SQL Server).