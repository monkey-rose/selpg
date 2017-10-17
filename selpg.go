package main

import (
	"fmt"
	"os"
	"strconv"
	"bufio"
)

type selpg_args struct {
	start_page int
	end_page int
	in_filename string
	page_len int
	page_type int
	print_dest string
}

var progname string
const INT32MAX = 1<<31 - 1

func process_args(ac int, av []string, psa *selpg_args) {
	//参数个数
	if ac < 3 {
		fmt.Fprintf(os.Stderr, "%s: not enough arguments\n", progname)
		usage()
		os.Exit(1)
	}

	//-s
	s1 := av[1];
	if s1[:2] != "-s" {
		fmt.Fprintf(os.Stderr, "%s: 1st arg should be -sstart_page\n", progname)
		fmt.Fprintf(os.Stderr, "%s", s1[:2])
		usage()
		os.Exit(2)
	}
	i,_ := strconv.Atoi(s1[2:])
	if i < 1 || i > (INT32MAX - 1) {
		fmt.Fprintf(os.Stderr, "%s: invalid start page %s\n", progname, s1[2:])
		usage()
		os.Exit(3)
	}
	(*psa).start_page = i

	//-e
	s1 = av[2];
	if s1[:2] != "-e" {
		fmt.Fprintf(os.Stderr, "%s: 1st arg should be -sstart_page\n", progname)
		usage()
		os.Exit(4)
	}
	i,_ = strconv.Atoi(s1[2:])
	if i < 1 || i > (INT32MAX - 1) || i < (*psa).start_page {
		fmt.Fprintf(os.Stderr, "%s: invalid start page %s\n", progname, s1[2:]);
		usage();
		os.Exit(5);
	}
	(*psa).end_page = i

	//other arguments
	argno := 3
	for argno <= (ac - 1) && av[argno][0] == '-' {
		s1 := av[argno]
		switch s1[1] {
		case 'l':
			i,_ := strconv.Atoi(s1[2:])
			if i < 1 || i > (INT32MAX - 1) {
				fmt.Fprintf(os.Stderr, "%s: invalid page length %s\n", progname, s1[2:])
				usage()
				os.Exit(6)
			}
			(*psa).page_len = i
			argno = argno+1
			continue
			break
		case 'f':
			if s1 != "-f" {
				fmt.Fprintf(os.Stderr, "%s: option should be \"-f\"\n", progname)
				usage()
				os.Exit(7)
			}
			(*psa).page_type = 'f'
			argno = argno+1
			continue
			break
		case 'd':
			if len([]rune(s1[2:])) < 1 {
				fmt.Fprintf(os.Stderr,
					"%s: -d option requires a printer destination\n", progname)
				usage()
				os.Exit(8)
			}
			(*psa).print_dest = s1[2:]
			argno = argno+1
			continue
			break
		default:
			fmt.Fprintf(os.Stderr, "%s: unknown option %s\n", progname, s1)
			usage()
			os.Exit(9)
			break
		}
	}

	//filename
	if argno <= ac-1 {
		(*psa).in_filename = av[argno]
		_, err := os.Open(av[argno])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: input file \"%s\" does not exist\n", progname, av[argno])
			os.Exit(10)
		}
	}
}

func process_input(psa *selpg_args) {
	var fin = os.Stdin
	var fout = os.Stdout
	var err error
	//input
	if (*psa).in_filename == "" {
		fin = os.Stdin
	} else {
		fin, err = os.Open((*psa).in_filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: could not open input file \"%s\"\n", progname, (*psa).in_filename)
			os.Exit(11)
		}
	}

	//output
	if (*psa).print_dest == "" {
		fout = os.Stdout
	} else {
		fmt.Fprintf(os.Stderr, "%s: could not open pipe to \"%s\"\n", progname, (*psa).print_dest)
		os.Exit(13)
	}

	inputReader := bufio.NewReader(fin)
	var page_ctr int
	if (*psa).page_type == 'l' {
		line_ctr := 0
		page_ctr = 1
		var line string
		for true {
			line, err = inputReader.ReadString('\n')
			if err != nil {
				break
			}
			line_ctr++
			if line_ctr > (*psa).page_len {
				page_ctr++
				line_ctr = 1
			}
			if page_ctr >= (*psa).start_page && page_ctr <= (*psa).end_page {
				fmt.Fprintf(fout, "%s", line)
			}
		}
	} else {
		var c = 'a'
		page_ctr = 1
		for true {
			c, _, err = inputReader.ReadRune()
			if err != nil {
				break
			}
			if c == '\f' {
				page_ctr++
			}
			if page_ctr >= (*psa).start_page && page_ctr <= (*psa).end_page {
				fmt.Fprintf(fout, "%c", c)
			}
		}
	}
	if page_ctr < (*psa).start_page {
		fmt.Fprintf(os.Stderr,
		"%s: start_page (%d) greater than total pages (%d), no output written\n", progname, (*psa).start_page, page_ctr)
	} else if page_ctr < (*psa).end_page {
		fmt.Fprintf(os.Stderr,"%s: end_page (%d) greater than total pages (%d), less output than expected\n", progname, (*psa).end_page, page_ctr);
	}
	fin.Close()
	fout.Close()
	fmt.Fprintf(os.Stderr, "%s: done\n", progname)
}

func main() {
	sa := selpg_args{}
	progname = os.Args[0]
	sa.start_page = -1
	sa.end_page = -1
	sa.in_filename = ""
	sa.page_len = 32
	sa.page_type = 'l'
	sa.print_dest = ""

	process_args(len(os.Args), os.Args, &sa)
	process_input(&sa)
}

func usage() {
	fmt.Fprintf(os.Stderr,
	"\nUSAGE: %s -sstart_page -eend_page [ -f | -llines_per_page ] [ -ddest ] [ in_filename ]\n", progname)
}
