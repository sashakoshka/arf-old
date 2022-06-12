package lineFile

import (
        "os"
        "fmt"
        "bufio"
        "strings"
        "strconv"
)

type LineFile struct {
        file   *os.File
        path   string
        module string
        lines  []string
}

func Open (path string, moduleName string) (lineFile *LineFile, err error) {
        lineFile = &LineFile {
                module: moduleName,
                path:   path,
        }
        
        lineFile.file, err = os.Open(path)
        defer lineFile.file.Close()

        scanner := bufio.NewScanner(lineFile.file)
        for scanner.Scan() {
                lineFile.lines = append(lineFile.lines, scanner.Text())
        }

        err = scanner.Err()
        return
}

func (lineFile *LineFile) GetLine (row int) (line string) {
        return lineFile.lines[row]
}

func (lineFile *LineFile) GetLength () (length int) {
        return len(lineFile.lines)
}

func (lineFile *LineFile) PrintWarning (
        column int,
        row int,
        cause ...interface {},
) {
        lineFile.printMistake("\033[33m!!!\033[0m", column, row, cause...)
}

func (lineFile *LineFile) PrintError (
        column int,
        row int,
        cause ...interface {},
) {
        lineFile.printMistake("\033[31mERR\033[0m", column, row, cause...)
}

func (lineFile *LineFile) PrintFatal (
        cause ...interface {},
) {
        fmt.Println ("\033[31mXXX\033[0m", "\033[90min\033[0m", lineFile.path,
                "\033[90mof\033[0m", lineFile.module)
        fmt.Print("    ")
        fmt.Println(cause...)
}

func (lineFile *LineFile) printMistake (
        kind string,
        column int,
        row int,
        cause ...interface{},
) {
        fmt.Println (
                kind, "\033[90min\033[0m", lineFile.path,
                "\033[34m" + strconv.Itoa(row) + ":" +
                strconv.Itoa(column),
                "\033[90mof\033[0m", lineFile.module)
        fmt.Println("   ", strings.TrimLeft(lineFile.lines[row], " "))

        fmt.Print("    ")
        for column > 0 {
                fmt.Print("-")
                column --
        }
        fmt.Println("^")
        
        fmt.Print("    ")
        fmt.Println(cause...)
}
