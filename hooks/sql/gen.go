//go:build ignore

package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"sort"
	"strings"
)

var connIfaces = []string{
	"driver.ConnBeginTx",
	"driver.ConnPrepareContext",
	"driver.Execer",
	"driver.ExecerContext",
	"driver.NamedValueChecker",
	"driver.Pinger",
	"driver.Queryer",
	"driver.QueryerContext",
	"driver.SessionResetter",
	"driver.Validator",
}

var stmtIfaces = []string{
	"driver.StmtExecContext",
	"driver.StmtQueryContext",
	"driver.ColumnConverter",
	"driver.NamedValueChecker",
}

func getHash(s []string) string {
	h := md5.New()
	io.WriteString(h, strings.Join(s, "|"))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func main() {
	comboConn := all(connIfaces)

	sort.Slice(comboConn, func(i, j int) bool {
		return len(comboConn[i]) < len(comboConn[j])
	})

	comboStmt := all(stmtIfaces)

	sort.Slice(comboStmt, func(i, j int) bool {
		return len(comboStmt[i]) < len(comboStmt[j])
	})

	b := bytes.NewBuffer(nil)
	b.WriteString("// Code generated. DO NOT EDIT.\n\n")
	b.WriteString("package sql\n\n")
	b.WriteString(`import "database/sql/driver"`)
	b.WriteString("\n\n")

	b.WriteString("func wrapConn(dc driver.Conn, opts Options) driver.Conn {\n")
	b.WriteString("\tc := &wrapperConn{conn: dc, opts: opts}\n")

	for idx := len(comboConn) - 1; idx >= 0; idx-- {
		ifaces := comboConn[idx]
		n := len(ifaces)
		if n == 0 {
			continue
		}
		h := getHash(ifaces)
		b.WriteString(fmt.Sprintf("\tif _, ok := dc.(wrapConn%04d_%s); ok {\n", n, h))
		b.WriteString("\treturn struct {\n")
		b.WriteString("\t\tdriver.Conn\n")
		b.WriteString(fmt.Sprintf("\t\t\t%s", strings.Join(ifaces, "\n\t\t\t")))
		b.WriteString("\t\t\n}{")
		for idx := range ifaces {
			if idx > 0 {
				b.WriteString(", ")
				b.WriteString("c")
			} else if idx == 0 {
				b.WriteString("c")
			} else {
				b.WriteString("c")
			}
		}
		b.WriteString(", c}\n")
		b.WriteString("}\n\n")
	}
	b.WriteString("return c\n")
	b.WriteString("}\n")

	for idx := len(comboConn) - 1; idx >= 0; idx-- {
		ifaces := comboConn[idx]
		n := len(ifaces)
		if n == 0 {
			continue
		}
		h := getHash(ifaces)
		b.WriteString(fmt.Sprintf("// %s\n", strings.Join(ifaces, "|")))
		b.WriteString(fmt.Sprintf("type wrapConn%04d_%s interface {\n", n, h))
		for _, iface := range ifaces {
			b.WriteString(fmt.Sprintf("\t%s\n", iface))
		}
		b.WriteString("}\n\n")
	}

	b.WriteString("func wrapStmt(stmt driver.Stmt, query string, opts Options) driver.Stmt {\n")
	b.WriteString("\tc := &wrapperStmt{stmt: stmt, query: query, opts: opts}\n")

	for idx := len(comboStmt) - 1; idx >= 0; idx-- {
		ifaces := comboStmt[idx]
		n := len(ifaces)
		if n == 0 {
			continue
		}
		h := getHash(ifaces)
		b.WriteString(fmt.Sprintf("\tif _, ok := stmt.(wrapStmt%04d_%s); ok {\n", n, h))
		b.WriteString("\treturn struct {\n")
		b.WriteString("\t\tdriver.Stmt\n")
		b.WriteString(fmt.Sprintf("\t\t\t%s", strings.Join(ifaces, "\n\t\t\t")))
		b.WriteString("\t\t\n}{")
		for idx := range ifaces {
			if idx > 0 {
				b.WriteString(", ")
				b.WriteString("c")
			} else if idx == 0 {
				b.WriteString("c")
			} else {
				b.WriteString("c")
			}
		}
		b.WriteString(", c}\n")
		b.WriteString("}\n\n")
	}
	b.WriteString("return c\n")
	b.WriteString("}\n")

	for idx := len(comboStmt) - 1; idx >= 0; idx-- {
		ifaces := comboStmt[idx]
		n := len(ifaces)
		if n == 0 {
			continue
		}
		h := getHash(ifaces)
		b.WriteString(fmt.Sprintf("// %s\n", strings.Join(ifaces, "|")))
		b.WriteString(fmt.Sprintf("type wrapStmt%04d_%s interface {\n", n, h))
		for _, iface := range ifaces {
			b.WriteString(fmt.Sprintf("\t%s\n", iface))
		}
		b.WriteString("}\n\n")
	}

	fmt.Printf("%s\n", b.String())
}

// all returns all combinations for a given string array.
func all[T any](set []T) (subsets [][]T) {
	length := uint(len(set))
	for subsetBits := 1; subsetBits < (1 << length); subsetBits++ {
		var subset []T
		for object := uint(0); object < length; object++ {
			if (subsetBits>>object)&1 == 1 {
				subset = append(subset, set[object])
			}
		}
		subsets = append(subsets, subset)
	}
	return subsets
}
