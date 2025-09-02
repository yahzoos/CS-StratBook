# CS-StratBook
Yahzoos' annotated nade guides for CS2

The Executables can be found in the CS-Stratbook Windows or Linux folder.

## Home Tab
When the application first loads it will create a settings.json file. The paths can be edited in the UI. 

Tags.json is the 'main' metadata file. This contains all of the annotation metadata in a single source.

The annotations folder is where the individual annotations are stored. If you change this path, be sure to find and replace the path to match in tags.json. >[!NOTE] On linux you will need to replace the `~` with the actual hard path>

Best Practice would be to store all of the annotation files in a different folder and only move the ones you would want to use into the default csgo path.

The generate new tags can be used when new (single) nade annotations are placed in the Annotation Folder Path. It will bring up a new window where a description, side and site can be added.

## Metadata Explorer Tab
 Click the refresh button if new annotations were added.

 Select a map and any filters, if none are selected it will show everything for the map. Click apply filters to have everything show in the right hand grid.

 Select a nade from the grid to have the details shown.

 Add/Remove will add the nade to the File Generator tab.

 The edit button has no functionality right now. To edit the metadata, manually edit the tags.json file. (click refresh button to reload the file)

 ## File Generator Tab

 This tab has all the nades selected from the previous tab.

 Write a name for the new annotation file (make sure to end with .txt)
 

# Using the annotation files
- In windows, place the contents of the \local folder into "C:\Program Files (x86)\Steam\steamapps\common\Counter-Strike Global Offensive\game\csgo\annotations\local"

If an annotation file was generated, create a folder in the above path with the same name as the txt file generated, and place the txt file in that new folder.

Example: If you created a new file called Top_Bannana_Control.txt create a folder and place it at this path C:\Program Files (x86)\Steam\steamapps\common\Counter-Strike Global Offensive\game\csgo\annotations\local\Top_Bannana_Control\Top_Bannana_Control.txt

To load the annotation file, load up a game and execute the following:

```
sv_cheats 1
sv_allow_annotations
annotations_load **FileName**
```

# Using the prac.cfg
The prac.cfg will create a practice server with the ability to buy full nades, enable the nade preview, and enable cheats and annotations.

Place the prac.cfg into ..\Steam\steamapps\common\Counter-Strike Global Offensive\game\csgo\cfg\

To load the prc.cfg file, load up a game and execute the following:

```
exec prac.cfg
```

## Using the Annotation Commands
Use this for documenation for the commands in CS2
https://steamcommunity.com/sharedfiles/filedetails/?id=3367125162
