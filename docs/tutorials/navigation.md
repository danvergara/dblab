One of the main features of dblab is to be a simple but very useful UI to interact with your database.  
![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/full-ui.png){ width="700" : .center }  

- <kbd>Ctrl+I</kbd> to move from Rows to Indexes and backwards.
- <kbd>Ctrl+T</kbd> to move from Rows to Constraints and backwards.
- <kbd>Ctrl+S</kbd> to move from Rows to Structure and backwards.
  
When selected each panel will show different information in the bottom box:

- Rows: Will show the result of the executed query. Press <kbd>Ctrl+Space</kbd> to execute the query.
    ![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/rows-result.png){ width="600" : .center }
- Structure: Will show the schema of the table selected  
    ![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/structure-result.png){ width="400" : .center }
- Constraints: Will show the constraints of the table selected  
    ![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/constraints-result.png){ width="400" : .center }
- Indexes: Will show the indexes of the table selected  
    ![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/indexes-result.png){ width="400" : .center }

In order to be able to see the information of `Structure`, `Constraints` or `Indexes` first you need to select a table from the left menu.  

To navigate there you can use:

- <kbd>Ctrl+H</kbd>: Toggle to the panel on the left
- <kbd>Ctrl+J</kbd>: Toggle to the panel below
- <kbd>Ctrl+K</kbd>: Toggle to the panel above
- <kbd>Ctrl+L</kbd>: Toggle to the panel on the right
 
Once you are placed above the correct name in the menu press <kbd>Enter</kbd> to select the table.
Now you can navigate to the different panels to see the information related to it.
![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/left-menu.png){ width="400" : .center }

### Key Bindings
Key                                     | Description
----------------------------------------|---------------------------------------
<kbd>Ctrl+Space</kbd>                   | If the query panel is active, execute the query
<kbd>Ctrl+D</kbd>                       | Cleans the whole text from the query editor, when the editor is selected
<kbd>Enter</kbd>                        | If the tables panel is active, list all the rows as a result set on the rows panel and display the structure of the table on the structure panel
<kbd>Ctrl+S</kbd>                       | If the rows panel is active, switch to the schema panel. The opposite is true
<kbd>Ctrl+T</kbd>                       | If the rows panel is active, switch to the constraints view. The opposite is true
<kbd>Ctrl+I</kbd>                       | If the rows panel is active, switch to the indexes view. The opposite is true
<kbd>Ctrl+H</kbd>                       | Toggle to the panel on the left
<kbd>Ctrl+J</kbd>                       | Toggle to the panel below
<kbd>Ctrl+K</kbd>                       | Toggle to the panel above
<kbd>Ctrl+L</kbd>                       | Toggle to the panel on the right
<kbd>Arrow Up</kbd>                     | Next row of the result set on the panel. Views: rows, table, constraints, structure and indexes
<kbd>k</kbd>                            | Next row of the result set on the panel. Views: rows, table, constraints, structure and indexes
<kbd>Arrow Down</kbd>                   | Previous row of the result set on the panel. Views: rows, table, constraints, structure and indexes
<kbd>j</kbd>                            | Previous row of the result set on the panel. Views: rows, table, constraints, structure and indexes
<kbd>Arrow Right</kbd>                  | Horizontal scrolling on the panel. Views: rows, constraints, structure and indexes
<kbd>l</kbd>                            | Horizontal scrolling on the panel. Views: rows, constraints, structure and indexes
<kbd>Arrow Left</kbd>                   | Horizontal scrolling on the panel. Views: rows, constraints, structure and indexes
<kbd>h</kbd>                            | Horizontal scrolling on the panel. Views: rows, constraints, structure and indexes
<kbd>g</kbd>                            | Move cursor to the top of the panel's dataset. Views: rows, constraints, structure and indexes
<kbd>G</kbd>                            | Move cursor to the bottom of the panel's dataset. Views: rows, constraints, structure and indexes
<kbd>Ctrl-F</kbd>                       | Move down by one page. Views: rows, constraints, structure and indexes
<kbd>Ctrl-B</kbd>                       | Move up by one page. Views: rows, constraints, structure and indexes
<kbd>Ctrl+c</kbd>                       | Quit 
