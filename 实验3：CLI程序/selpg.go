package main

import (
    "fmt"
	"os"
	"bufio"
	"github.com/spf13/pflag"
	"os/exec"
	"io"
)

type selpg_args struct {
	spage  int   		//start page
	epage    int    	//end page
	infile string 		//input file name
	plen    int    		//page length
	ptype   bool 		//page type -f or -l
	printdest  string	//print destination, ouput file name
	
}

var programName string 		//program name is the first argument

func main() {
	selArgs := selpg_args{}
	programName = os.Args[0]

	initArgs(&selArgs)				//init the flag arguments, setting default values, and get arguments from the command line
	handleArgs(len(os.Args), &selArgs) //handle the arguments
	process(&selArgs)				//run the CLI command
}

func initArgs(args *selpg_args) {
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"USAGE: \n%s -s start_page -e end_page [ -f | -l lines_per_page ]" + 
			" [ -d dest ] [ in_filename ]\n", )
		pflag.PrintDefaults()
	}
	pflag.IntVarP(&args.spage,"start", "s", 0, "start page")
	pflag.IntVarP(&args.epage,"end","e", 0, "emd page")
	//这里跟官方文档不同，使用10行作为默认值
	pflag.IntVarP(&args.plen,"linenum", "l", 10, "page length (lines)")
	pflag.BoolVarP(&args.ptype,"printdes","f", false, "'l' for lines-delimited, 'f' for form-feed-delimited. default is 'l'")
	pflag.StringVarP(&args.printdest, "othertype","d", "", "print destination")
	pflag.Parse() //解析
}

func handleArgs(argNum int, args *selpg_args) {
	/* 检查参数合不合法 */
	if argNum < 3 {
		fmt.Fprintf(os.Stderr, "%s: not enough arguments\n", programName)
		pflag.Usage()
		os.Exit(1)
	}

	/* 第一个参数，spage*/
	if os.Args[1][0] != '-' || os.Args[1][1] != 's' {
		fmt.Fprintf(os.Stderr, "%s: 1st arg should be -s=spage\n", programName)
		pflag.Usage()
		os.Exit(2)
	}
	if args.spage < 1  {
		fmt.Fprintf(os.Stderr, "%s: invalid start page %s\n", programName, args.spage)
		pflag.Usage()
		os.Exit(3)
	}

	/* 第二个参数，epage*/
	if os.Args[3][0] != '-' || os.Args[3][1] != 'e' {
		fmt.Fprintf(os.Stderr, "%s: 2nd arg should be -e=epage\n", programName)
		pflag.Usage()
		os.Exit(4)
	}
	if args.epage < 1  || args.epage < args.spage  {
		fmt.Fprintf(os.Stderr, "%s: invalid end page %s\n", programName, args.epage)
		pflag.Usage()
		os.Exit(5)
	}
    
	/* 处理可选择的参数 */
	if args.plen != 5 {
		if args.plen < 1  {
			fmt.Fprintf(os.Stderr, "%s: invalid page length %s\n", programName, args.plen)
			pflag.Usage()
			os.Exit(6)
		}
	}


	/* 第三个参数，infile*/
	if pflag.NArg() > 0 {
		args.infile = pflag.Arg(0)
		/*检查文件是否存在 */
		file, err := os.Open(args.infile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: input file \"%s\" does not exist\n", programName, args.infile)
			os.Exit(7)
		}
		/* 是否可读 */
		file, err = os.OpenFile(args.infile, os.O_RDONLY, 0666)
		if err != nil {
			if os.IsPermission(err) {
				fmt.Fprintf(os.Stderr, "%s: input file \"%s\" exists but cannot be read\n", programName, args.infile)
				os.Exit(8)
			}
		}
		file.Close()
	}
}

func process(args *selpg_args) {
	fin := os.Stdin
	fout := os.Stdout
	var (
		 pageCnt int
		 lineCnt int
		 err error
		 err1 error
		 err2 error
		 line string
		 cmd *exec.Cmd
		 stdin io.WriteCloser
	)
	/* 处理fin输入 */
	if args.infile != "" {
		fin, err1 = os.Open(args.infile)
		if err1 != nil {
			fmt.Fprintf(os.Stderr, "%s: could not open input file \"%s\"\n", programName, args.infile)
			os.Exit(11)
		}
	}

	if args.printdest != "" {
		//使用exec库的Command的函数调用，command返回cmd结构来执行带有相关参数的命令，它仅仅设定cmd结构中的Path和Args参数
		cmd = exec.Command("cat", "-n")
		stdin, err = cmd.StdinPipe()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		stdin = nil
	}

/* begin one of two main loops based on page type */
//NewReader 相当于 NewReaderSize(rd, 4096)
	rd := bufio.NewReader(fin)
	if args.ptype == false {
		lineCnt = 0
		pageCnt = 1
		for true {
			line, err2 = rd.ReadString('\n')
			if err2 != nil { /* error or EOF */
				break
			}
			lineCnt++
			if lineCnt > args.plen {
				pageCnt++
				lineCnt = 1
			}
			if pageCnt >= args.spage && pageCnt <= args.epage {
				fmt.Fprintf(fout, "%s", line)
			}
		}
	} else {
		pageCnt = 1
		for true {
			c, err3 := rd.ReadByte()
			if err3 != nil { /* error or EOF */
				break
			}
			if c == '\f' {
				pageCnt++
			}
			if pageCnt >= args.spage && pageCnt <= args.epage {
				fmt.Fprintf(fout, "%c", c)
			}
		}
		fmt.Print("\n")
	}

	/* end main loop */
	if pageCnt < args.spage {
		fmt.Fprintf(os.Stderr, "%s: spage (%d) greater than total pages (%d), no output written\n", programName, args.spage, pageCnt)
	} else if pageCnt < args.epage {
			fmt.Fprintf(os.Stderr, "%s: epage (%d) greater than total pages (%d), less output than expected\n", programName, args.epage, pageCnt)
	}
	
	if args.printdest != "" {
		stdin.Close()
		cmd.Stdout = fout
		cmd.Run()
	}
	fmt.Fprintf(os.Stderr,"\n---------------\nProcess end\n")
	fin.Close()
	fout.Close()
}

