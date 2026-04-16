Output screenshots/dblab-demo.gif

Set FontSize 14
Set Width 1920
Set Height 1080
Set Theme "Catppuccin Mocha"

# Give the terminal a moment to breathe before acting
Sleep 1s

# Type the dblab command and press Enter
Type "dblab"
Enter

# Wait for the app to initialize and render
Sleep 2s

# Press enter to select the Postgres driver
Enter

# Wait for the next form
Sleep 1s

# Fill out the connection params form
Type "localhost"
Enter
Sleep 500ms
Type "5432"
Enter
Sleep 500ms
Type "postgres"
Enter
Sleep 500ms
Type "password"
Enter
Sleep 500ms
Type "users"
Enter
Sleep 500ms
Type "50"
Enter

# Wait for the last form
Sleep 1s

# Press enter to select to disable the ssl mode
Enter

# Wait a second for the app to initialize and render
Sleep 1s

# Pres i to enter into INSERT mode
Type "i"
Sleep 500ms

# Simulate the dblab interaction
Type "SELECT * FROM users LIMIT 30;"
Ctrl+E
Sleep 1s
Ctrl+J
Sleep 500ms
Type "G"
Sleep 500ms
Type "g"

# Admire the output for a bit.
Sleep 3s
