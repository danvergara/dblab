One of the main features of dblab is to be a simple but very useful UI to interact with your database.  
![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/full-ui.png){ width="700" : .center }  

If the query panel is active, type the desired query and press <kbd>ctrl+e</kbd> to see the results on the rows panel below.
Otherwise, you might be located at the tables panel, then you can navigate by using the arrows <kbd>Up</kbd> and <kbd>Down</kbd> (or the keys <kbd>k</kbd> and <kbd>j</kbd> respectively). If you want to see the rows of a table, press <kbd>Enter</kbd>. To see the schema of a table, locate yourself on the `tables` panel and press <kbd>tab</kbd> to switch to the `columns` panel, then switch <kbd>shift+tab</kbd> to switch back.

Now, there's a menu to navigate between hidden views by just clicking on the desired options:

- Data: Will show the result of the executed query. Press <kbd>ctrl+e</kbd> to execute the query.
    ![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/rows-result.png){ width="600" : .center }
- Columns: Will show the schema of the table selected  
    ![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/structure-result.png){ width="400" : .center }
- Indexes: Will show the indexes of the table selected  
    ![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/indexes-result.png){ width="400" : .center }
- Constraints: Will show the constraints of the table selected  
    ![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/constraints-result.png){ width="400" : .center }

In order to be able to see the information of `Columns`, `Indexes` or `Constraints` first you need to select a table from the left menu.  

To navigate there you can use:

- <kbd>tab</kbd>: If the resultset panel is active, press tab to navigate to the next metadata tab
- <kbd>shift+tab</kbd>: If the resultset panel is active, press shift+tab to navigate to the previous metadata tab
 
Once you are placed above the correct name in the left menu press <kbd>Enter</kbd> to select the table.
Now you can navigate to the different panels to see the information related to it.

![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/left-menu.png){ width="400" : .center }

The navigation buttons were removed since they are too slow to really navigate the content of a table. The user is better off typing a `SELECT` statement with proper `OFFSET` and `LIMIT`.

The `--db` flag is now optional (except for Oracle), meaning that the user will be able to see the list of databases they have access to. The regular list of tables will be replaced with a tree structure showing a list of databases and their respective list of tables, branching off each database. Due to the nature of the vast majority of DBMSs that don't allow cross-database queries, dblab has to open an independent connection for each database. The side effect of this decision, is that the user has to press `Enter` on the specific database of interest. An indicator showing the current active database will appear at the bottom-right of the screen. To change the focus, just hit enter on another database. Once a database is selected, the usual behavior of inspecting tables remains the same.

![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/screenshots/tree-view.png){ width="400" : .center }

At the time to navigate queries' resultsets, the cell will be highlighted, so the user can know which table cell is on. This is important because you can press the `Enter` key on a cell of interest to copy its content.

### Key Bindings
| Key                                    | Description                           |
|----------------------------------------|----------------------------------------|
|<kbd>ctrl+e</kbd>                       | If the query editor is active, execute the query |
|<kbd>Ctrl+D</kbd>                       | Cleans the whole text from the query editor, when the editor is selected |
|<kbd>Enter</kbd>                        | If the tables panel is active, list all the rows as a result set on the rows panel and display the structure of the table on the structure panel |
|<kbd>tab</kbd>                          | If the resultset panel is active, press tab to navigate to the next metadata tab |
|<kbd>shift+tab</kbd>                    | If the resultset panel is active, press shift+tab to navigate to the previous metadata tab |
|<kbd>Ctrl+H</kbd>                       | Toggle to the panel on the left |
|<kbd>Ctrl+J</kbd>                       | Toggle to the panel below |
|<kbd>Ctrl+K</kbd>                       | Toggle to the panel above |
|<kbd>Ctrl+L</kbd>                       | Toggle to the panel on the right |
|<kbd>Arrow Up</kbd>                     | Vertical scrolling on the panel. Views: rows, table, constraints, structure and indexes |
|<kbd>k</kbd>                            | Vertical scrolling on the panel. Views: rows, table, constraints, structure and indexes |
|<kbd>Arrow Down</kbd>                   | Vertical on the panel. Views: rows, table, constraints, structure and indexes |
|<kbd>j</kbd>                            | Vertical on the panel. Views: rows, table, constraints, structure and indexes |
|<kbd>Arrow Right</kbd>                  | Horizontal scrolling on the panel. Views: rows, constraints, structure and indexes |
|<kbd>l</kbd>                            | Horizontal scrolling on the panel. Views: rows, constraints, structure and indexes |
|<kbd>Arrow Left</kbd>                   | Horizontal scrolling on the panel. Views: rows, constraints, structure and indexes |
|<kbd>h</kbd>                            | Horizontal scrolling on the panel. Views: rows, constraints, structure and indexes |
|<kbd>g</kbd>                            | Move cursor to the top of the panel's dataset. Views: rows, constraints, structure and indexes |
|<kbd>G</kbd>                            | Move cursor to the bottom of the panel's dataset. Views: rows, constraints, structure and indexes |
|<kbd>Ctrl+c</kbd>                       | Quit |
