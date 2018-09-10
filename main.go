package main

import (
	"flag"
	"fmt"
    "os"
	"path/filepath"
	"bufio"
	"bytes"
	"strconv"
)

func readRegions(path string) [][2]int {
	var res [][2]int

    fh, err := os.Open(path)

    if err != nil {
		fmt.Printf("Can't open region list file %s\n", path)
		os.Exit(2)
    }

	defer fh.Close()

    fr := bufio.NewReader(fh)
	buf := make([]byte, 1024)
    for {
		buf, _ , err = fr.ReadLine()
		if err != nil {
			return res
		}
		fields := bytes.Split(buf, []byte(","))
		a, _ := strconv.Atoi(string(fields[0]))
		b, _ := strconv.Atoi(string(fields[1]))
		res = append(res, [2]int{a, b})
	}
}

func multisect(arguments []string) {
	subCommand := flag.NewFlagSet("musect", flag.ExitOnError)

	regionListPath := subCommand.String("regions", "", `Path of file to specify regions. ex: input_regions.txt`)
	inpath := subCommand.String("input", "", `Path of file to read. ex: ../data/grants2012/ipg120110.xml`)
	outprefix := subCommand.String("outprefix", "", `Prefix of path to output result. ex: ../data/tmp/tempvar`)

	subCommand.Parse(arguments)

	
	err := os.MkdirAll(filepath.Dir(*outprefix), 0744)
    if err != nil {
		fmt.Printf("Can't create dir  %s\n", filepath.Dir(*outprefix))
		os.Exit(2)
    }

	regions := readRegions(*regionListPath)
    fh, err := os.Open(*inpath)

    if err != nil {
		fmt.Printf("Can't open input file %s\n", *inpath)
		os.Exit(2)
    }

	defer fh.Close()
    fr := bufio.NewReader(fh)

	endnum := 0
	for i, v := range(regions) {
		fpath := fmt.Sprintf("%s_%08d.txt", *outprefix, i)
		endnum = subsectOne(fr, endnum, fpath, v[0], v[1])
		// fmt.Printf("%s: %d: %d\n", fpath, v[0], v[1])
	}
}

func subsectOne(reader *bufio.Reader, current int, outpath string, start int, end int) int {
	if start < current {
		fmt.Printf("Interleave regions, not supported: %d, %d\n", current, start)
		os.Exit(2)
	}

	fout, err := os.Create(outpath)
    if err != nil {
		fmt.Printf("Can't create output file: %s\n", outpath)
		os.Exit(2)
    }
	defer fout.Close()


	fw := bufio.NewWriter(fout)

	buf := make([]byte, 1024)
	count := current
    for {
		buf, _ , err = reader.ReadLine()
		if err != nil {
			fw.Flush()
			return -1
		}
		count++
		if count >= start {
			fmt.Fprintf(fw, "%s\n", buf)
		}

		// do not need to read one line more.
		if count >= end {
			fw.Flush()
			return count
		}
	}

}

func subsectStdout(arguments []string)  {
	subCommand := flag.NewFlagSet("one", flag.ExitOnError)

	start := subCommand.Int("start", -1, "Specifiy start of region. ex: 100")
	end := subCommand.Int("end", -1, "Specify end of region. ex: 1000")
	inpath := subCommand.String("input", "", `Path of file to output result. ex: ../data/grants2012/ipg120110.xml`)

	subCommand.Parse(arguments)


	if *start == -1 || *end == -1 || *inpath == "" {
		fmt.Println(`Invalid argument for command "sub"`)
		subCommand.PrintDefaults()
		os.Exit(2)
	}

    fh, err := os.Open(*inpath)

    if err != nil {
		fmt.Printf("Can't open input file %s\n", *inpath)
		os.Exit(2)
    }

	defer fh.Close()



    fr := bufio.NewReader(fh)

	buf := make([]byte, 1024)
	count := 0
    for {
		buf, _ , err = fr.ReadLine()
		if err != nil {
			return
		}
		count++
		if count >= *start {
			fmt.Printf("%s\n", buf)
		}

		// do not need to read one line more.
		if count >= *end {
			return			
		}
	}

}

func main() {

	if len(os.Args) == 1 {
		fmt.Println("usage: musect <command> [<args>]")
		fmt.Println("Commands are: ")
		fmt.Println(" one   Print one region of a large file.")
		fmt.Println(" list  Write regions of a large file to files.")
		return
	}



	switch os.Args[1] {
	case "list":
		multisect(os.Args[2:])
	case "one":
		subsectStdout(os.Args[2:])
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}


}