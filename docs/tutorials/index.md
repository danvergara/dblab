Once the installation is done you can start by typing dblab 

```bash
$ dblab
Select the database driver:

[x] postgres
[ ] mysql
[ ] sqlite3
```
When you don't use any flag with `dblab`, a prompt will ask for the driver to be used.  
You can navigate through the options with the <kbd>Arrow Up</kbd>  and <kbd>Arrow Down</kbd> keys or with <kbd>j</kbd> an <kbd>k</kbd> keys.
When the right driver is selected you can press <kbd>Enter</kbd> to apply the selection.

For this example we are going to use a sample SQlite database file with a few tables from [here](https://raw.githubusercontent.com/danvergara/dblab/master/docs/tutorials/resources/EssentialSQL.db), but you can use your own sqlite file.

In this case we are going to choose the sqlite3 driver, so the prompt will ask for the path of the db file and the size limit of the result from the queries

```bash
Introduce the connection params:

> File Path
> Limit
```
Then you will be asked to select the ssl mode for the connection with your database, in the case of sqlite3 you can just press <kbd>Enter</kbd>.

```bash
Select the ssl mode (just press enter if you selected sqlite3):
```

If everything went well you should see the UI  
![dblab](https://raw.githubusercontent.com/danvergara/dblab/main/docs/tutorials/resources/images/full-ui.png){ width="500" : .center }

For further knowledge on the navigation of the UI you can check this [first steps in navigation](https://dblab.danvergara.com/tutorials/navigation/)


