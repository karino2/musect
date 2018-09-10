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
```

