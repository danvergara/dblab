One of the main features of dblab is its simple but very useful UI for interacting with your database.  
![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/full-ui.png){ width="700" : .center }  

### Query editor

The query editor uses **normal** and **insert** modes (similar to Vim). When you focus the query editor, it starts in **normal** mode. Press <kbd>i</kbd> to enter insert mode and type or edit SQL; press <kbd>Escape</kbd> to return to normal mode (the cursor moves one character to the left, as in Vim). In insert mode, use the arrow keys to move the cursor; in normal mode, use <kbd>h</kbd>, <kbd>j</kbd>, <kbd>k</kbd>, and <kbd>l</kbd> instead (configurable in `.dblab.yaml` with `--keybindings` or `-k`; see [Key bindings configuration](../usage.md#key-bindings-configuration)). In normal mode, <kbd>dd</kbd> deletes the current line, <kbd>yy</kbd> yanks the current line into an internal register, <kbd>p</kbd> pastes that line after the current line, and <kbd>x</kbd> deletes the character under the cursor. <kbd>0</kbd> and <kbd>$</kbd> move to the beginning or end of the current line in the query buffer. Press <kbd>ctrl+e</kbd> to execute the query (this uses the `keybindings.editor.execute-query` binding); whitespace-only queries are ignored.

Otherwise, you might be located at the tables panel, where you can navigate using the arrows <kbd>Up</kbd> and <kbd>Down</kbd> (or the keys <kbd>k</kbd> and <kbd>j</kbd> respectively). If you want to see the rows of a table, press <kbd>Enter</kbd>. To see the schema of a table, locate yourself on the `tables` panel and press <kbd>tab</kbd> to switch to the `columns` panel, then use <kbd>shift+tab</kbd> to switch back.

Now, there's a menu to navigate between hidden views by just clicking on the desired options:

- Data: Will show the result of the executed query. Press <kbd>ctrl+e</kbd> to execute the query.
    ![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/rows-result.png){ width="600" : .center }
- Columns: Will show the schema of the table selected  
    ![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/structure-result.png){ width="400" : .center }
- Indexes: Will show the indexes of the table selected  
    ![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/indexes-result.png){ width="400" : .center }
- Constraints: Will show the constraints of the table selected  
    ![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/constraints-result.png){ width="400" : .center }

In order to be able to see the information for `Columns`, `Indexes`, or `Constraints`, first you need to select a table from the left menu.  

To navigate there you can use:

- <kbd>tab</kbd>: If the result set panel is focused, press tab to navigate to the next metadata tab.
- <kbd>shift+tab</kbd>: If the result set panel is focused, press shift+tab to navigate to the previous metadata tab.
 
Once the correct name is highlighted in the left menu, press <kbd>Enter</kbd> to select the table.
Now you can navigate to the different panels to see the related information.

![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/left-menu.png){ width="400" : .center }

The navigation buttons were removed since they are too slow to navigate the content of a table effectively. The user is better off typing a `SELECT` statement with proper `OFFSET` and `LIMIT`.

The `--db` flag is now optional (except for Oracle), meaning that the user will be able to see the list of databases they have access to. The regular list of tables will be replaced with a tree structure showing a list of databases and their respective list of tables, branching off each database. Due to the nature of the vast majority of DBMSs that don't allow cross-database queries, dblab has to open an independent connection for each database. The side effect of this decision is that the user has to press `Enter` on the specific database of interest. An indicator showing the current active database will appear at the bottom-right of the screen. To change the focus, just hit enter on another database. Once a database is selected, the usual behavior of inspecting tables remains the same.

![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/screenshots/tree-view.png){ width="400" : .center }

When navigating query result sets, the cell will be highlighted so the user can see which table cell is selected. This is important because you can press the `Enter` key on a cell of interest to copy its content.

### Key Bindings
| Key                                    | Description                           |
|----------------------------------------|----------------------------------------|
|<kbd>ctrl+e</kbd>                       | If the query editor is focused, execute the query (also works in insert and normal mode) |
|<kbd>i</kbd>                            | If the query editor is focused in normal mode, enter insert mode |
|<kbd>Escape</kbd>                       | If the query editor is focused in insert mode, return to normal mode |
|<kbd>dd</kbd>                           | If the query editor is focused in normal mode, delete the current line |
|<kbd>yy</kbd>                           | If the query editor is focused in normal mode, yank the current line |
|<kbd>p</kbd>                            | If the query editor is focused in normal mode, paste the yanked or deleted line after the current line |
|<kbd>x</kbd>                            | If the query editor is focused in normal mode, delete the character under the cursor |
|<kbd>Enter</kbd>                        | If the tables panel is focused, list all rows as a result set on the rows panel and display the structure of the table on the structure panel |
|<kbd>tab</kbd>                          | If the result set panel is focused, press tab to navigate to the next metadata tab |
|<kbd>shift+tab</kbd>                    | If the result set panel is focused, press shift+tab to navigate to the previous metadata tab |
|<kbd>Ctrl+H</kbd>                       | Toggle to the panel on the left |
|<kbd>Ctrl+J</kbd>                       | Toggle to the panel below |
|<kbd>Ctrl+K</kbd>                       | Toggle to the panel above |
|<kbd>Ctrl+L</kbd>                       | Toggle to the panel on the right |
|<kbd>Arrow Up</kbd>                     | If the query editor is focused in insert mode, move the cursor up. If the results panel is focused, navigate the table upward (all tabs on the results panel). |
|<kbd>k</kbd>                            | If the query editor is focused in normal mode, move the cursor up. If the results panel is focused, navigate the table upward (all tabs on the results panel). |
|<kbd>Arrow Down</kbd>                   | If the query editor is focused in insert mode, move the cursor down. If the results panel is focused, navigate the table downward (all tabs on the results panel). |
|<kbd>j</kbd>                            | If the query editor is focused in normal mode, move the cursor down. If the results panel is focused, navigate the table downward (all tabs on the results panel). |
|<kbd>Arrow Right</kbd>                  | If the query editor is focused in insert mode, move the cursor right. If the results panel is focused, navigate the table to the right (all tabs on the results panel). |
|<kbd>l</kbd>                            | If the query editor is focused in normal mode, move the cursor right. If the results panel is focused, navigate the table to the right (all tabs on the results panel). |
|<kbd>Arrow Left</kbd>                   | If the query editor is focused in insert mode, move the cursor left. If the results panel is focused, navigate the table to the left (all tabs on the results panel). |
|<kbd>h</kbd>                            | If the query editor is focused in normal mode, move the cursor left. If the results panel is focused, navigate the table to the left (all tabs on the results panel). |
|<kbd>g</kbd>                            | If the results panel is focused, move to the top of the dataset (all tabs on the results panel). |
|<kbd>G</kbd>                            | If the results panel is focused, move to the bottom of the dataset (all tabs on the results panel). |
|<kbd>0</kbd>                            | If the query editor is focused in normal mode, move to the start of the current line. If the results panel is focused, move to the left edge of the row (all tabs on the results panel). |
|<kbd>$</kbd>                            | If the query editor is focused in normal mode, move to the end of the current line. If the results panel is focused, move to the right edge of the row (all tabs on the results panel). |
|<kbd>Ctrl+c</kbd>                       | Quit |
