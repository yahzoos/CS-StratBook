# Tags Module

This function allows for the automated creation of json metadata for the annotation files. This data can then be joined into a single json file, SQLite db, or MongoDB to be querried against.

# Json Structure
Here is an example of an empty json file.
```json
{
  "file_name": "",
  "file_path": "",
  "ImagePath": "",
  "NadeName": "",
  "Description": "",
  "map_name": "",
  "side": "",
  "nade_type": "",
  "site_location": ""
}
```

Here is the example within the `Sample` folder
```json
{
  "file_name": "T2Camera.txt",
  "file_path": "Sample\\T2Camera\\T2Camera.txt",
  "ImagePath": "Sample\\T2Camera\\T2Camera.png",
  "NadeName": "T2Camera",
  "Description": "Smokes Camera from T Spawn",
  "map_name": "de_train",
  "side": "T",
  "nade_type": "smoke",
  "site_location": "A"
}
```
Where: \
 `file_name` is the name of the annotation txt file.\
 `file_path` is the full path to the annotation txt file.\
 `image_path` is the full path to the annotation png file.\
 `nade_name` is the name of the parent folder - ideally matches the file names.\
 `description` required user input. Describes the purpose of the grenade.\
 `map_name` is the name of the map. Pulled fromt he annotation txt file.\
 `side` optional user input. Can be T/CT or empty.\
 `nade_type` is the type of grenade. smoke/flash/molotov/he_grenade.\
 `site_location` optional user input. Can be A/B/Mid or empty.

 



