One of the main features of dblab is its simple but very useful UI for interacting with your database.  
![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/full-ui.png){ width="700" : .center }

This tutorial walks through the UI step by step. For the complete list of key bindings, see the [key bindings reference](../usage.md#key-bindings).

### The three panels

The UI is split into three panels:

- the **sidebar tree** on the left, listing the database, its schemas and its tables
- the **query editor** on the top right, where you write SQL
- the **result set** panel below the editor, where results and table metadata are displayed

Move focus between them with <kbd>Ctrl+H</kbd>, <kbd>Ctrl+J</kbd>, <kbd>Ctrl+K</kbd> and <kbd>Ctrl+L</kbd> — left, down, up and right respectively. The focused panel is highlighted with a brighter border.

### Selecting a table

Focus the sidebar and move through the tree with the <kbd>Arrow Up</kbd> and <kbd>Arrow Down</kbd> keys, or with <kbd>k</kbd> and <kbd>j</kbd>.

![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/left-menu.png){ width="400" : .center }

On a long catalog you don't have to walk the whole tree: <kbd>alt+k</kbd> and <kbd>alt+j</kbd> jump to the first and last visible node, and <kbd>/</kbd> starts a search by name — type part of the table you're after and press <kbd>Esc</kbd> when you're done searching. If a name is wider than the panel, <kbd>h</kbd> and <kbd>l</kbd> scroll the tree sideways.

Once the table you want is highlighted, press <kbd>Enter</kbd> to select it. dblab loads its rows into the result set panel and fills in the metadata tabs described below.

### Inspecting a table

The result set panel has one tab per view of the selected table. Press <kbd>tab</kbd> to move to the next tab and <kbd>shift+tab</kbd> to move back.

- **Data**: the rows of the table, or the result of the query you executed
    ![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/rows-result.png){ width="600" : .center }
- **Columns**: the schema of the table selected  
    ![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/structure-result.png){ width="400" : .center }
- **Indexes**: the indexes of the table selected  
    ![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/indexes-result.png){ width="400" : .center }
- **Constraints**: the constraints of the table selected  
    ![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/constraints-result.png){ width="400" : .center }

!!! note

    In order to see anything under `Columns`, `Indexes` or `Constraints`, you first need to select a table from the sidebar.

While moving through a result set, the selected cell is highlighted so you can see where you are. Press <kbd>Enter</kbd> on a cell to copy its content.

### Writing your first query

Focus the query editor with the panel navigation keys. The editor works like Vim: it starts in **normal** mode, where keystrokes are commands rather than text.

Press <kbd>i</kbd> to switch to **insert** mode, then type your query:

```{ .sql .copy }
SELECT * FROM customers LIMIT 10;
```

Press <kbd>Escape</kbd> to go back to normal mode, then press <kbd>ctrl+e</kbd> to execute. The result appears in the **Data** tab.

Normal mode also gives you the line-oriented editing commands you'd expect — <kbd>dd</kbd> to delete a line, <kbd>yy</kbd> and <kbd>p</kbd> to copy and paste one, <kbd>Ctrl+D</kbd> to clear the editor. The [key bindings reference](../usage.md#query-editor-normal-mode) has the full list.

### Running one query out of several

You don't have to clear the editor between queries. Keep several statements around, put the cursor on the one you care about, and press <kbd>ctrl+r</kbd> to execute only that line.

If you'd rather run all of them at once, separate them with semicolons and press <kbd>ctrl+e</kbd>:

```{ .sql .copy }
SELECT * FROM users; SELECT * FROM orders; SELECT count(*) FROM products;
```

![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/dblab-multi-query.png){ width="700" : .center }

Each statement gets its own result tab — "query #1", "query #2" and so on, three of them for the example above. Up to 5 statements can be run per batch. If one of them fails, its tab shows the error and the others still show their results.

While a batch is running, press <kbd>Ctrl+c</kbd> to cancel it; press <kbd>Ctrl+c</kbd> again to quit dblab.

### Reusing a query you ran before

![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/assets/tutorials/images/query-history.png){ width="400" : .center }

Every query you execute is saved to a local history file, so it survives across sessions. Press <kbd>F8</kbd> to open the history view, which lists your past queries newest-first. Type to filter the list, press <kbd>Enter</kbd> to load the highlighted query back into the editor, or press <kbd>Esc</kbd> to go back without picking anything.

### When you forget a key binding

![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/screenshots/dblab-help-modal.png){ width="700" : .center }

Press <kbd>?</kbd> at any point to bring up the help modal, which lists every key binding in a centered overlay. Press <kbd>Esc</kbd> to dismiss it; focus returns to the query editor.

!!! tip

    There are no pagination controls in the result set panel — they proved too slow to page through a table effectively. To work through a large table, write a `SELECT` with an explicit `OFFSET` and `LIMIT` instead.

### Next steps

All the key bindings shown here are defaults, and every one of them can be changed through the `.dblab.yaml` file. See [key bindings configuration](../usage.md#key-bindings-configuration) for how to do that, and the [key bindings reference](../usage.md#key-bindings) for the complete list.
