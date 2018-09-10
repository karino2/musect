# musect
Large file splitter written in go.

Read multiple region specified by regionlist file from large inputfile, 

## Usage example:


### Print one region (for debug)

```
musect one -start 1 -end 5 -input "../data/grants2012/ipg120110.xml"
```


### Multi sect file

```
musect list -regions regionlist.txt -outprefix temp/result  -input ../data/grants2012/ipg120103.xml
```

The format of regionlist.txt is like

```
100,230
1000,1020
23523,23600
```

Result path is like

```
temp/result_00000000.txt
temp/result_00000001.txt
temp/result_00000002.txt
```

### Python sample

```
import subprocess

def subfile(fpath, start, end):
    return subprocess.check_output(["musect", "one", "-start", str(start), "-end", str(end), "-input", fpath],universal_newlines=True).split("\n")


TEMP_SUBFILE_DIR="temp_subfile"

import os
import glob
def _start_end_list_to_file(start_end_list, fpath):
    with open(fpath, "w") as f:
        for start, end in start_end_list:
            f.write("{},{}\n".format(start, end))

REGIONLIST_PATH="{}/regions.txt".format(TEMP_SUBFILE_DIR)
TEMP_OUT_PREFIX="{}/result".format(TEMP_SUBFILE_DIR)
            
def subfiles(fpath, start_end_list):
    os.makedirs(TEMP_SUBFILE_DIR, exist_ok =True)
    _start_end_list_to_file(start_end_list, REGIONLIST_PATH)
    subprocess.call(["musect", "list", "-regions", REGIONLIST_PATH, "-outprefix", TEMP_OUT_PREFIX, "-input", fpath])
    result_files = sorted(glob.glob("{}_*.txt".format(TEMP_OUT_PREFIX)))
    res = []
    for fpath in result_files:
        with open(fpath) as f:
            lines = [line.rstrip("\n") for line in f]
            res.append(lines)
    os.remove(REGIONLIST_PATH)
    [os.remove(path) for path in result_files]
    os.rmdir(TEMP_SUBFILE_DIR)
    return res
```

