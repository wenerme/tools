# Stardict

## Reference
* [StarDictFileFormat](http://www.huzheng.org/stardict/StarDictFileFormat)

## Files
Every dictionary consists of these files:

0. somedict.ifo
0. somedict.idx or somedict.idx.gz
0. somedict.dict or somedict.dict.dz
0. somedict.syn (optional)

StarDict search for dictionaries in the following predefined directories:

0. gStarDictDataDir + "/dic",
0. "/usr/share/stardict/dic",
0. g_get_home_dir() + "/.stardict/dic".
