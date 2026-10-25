package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/scanner"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type edit struct {
	start int
	end   int
	text  []byte
}

type fileChange struct {
	path string
	data []byte
	mode os.FileMode
}

func main() {
	root := "tools/formatgo/testdata"
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	root = filepath.Clean(root)
	rootSlash := filepath.ToSlash(root)
	testDataRoot := strings.HasPrefix(rootSlash, "tools/formatgo/testdata")

	var paths []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == ".git" || info.Name() == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		if !testDataRoot &&
			(strings.HasPrefix(filepath.ToSlash(path), "tools/formatgo/")) {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	if err != nil {
		fail(err)
	}

	sort.Strings(paths)
	fmt.Printf("Found %d Go files.\n", len(paths))

	// Format and validate every file before writing any changes.
	var changes []fileChange
	for i, path := range paths {
		fmt.Printf("[%d/%d] Checking %s\n", i+1, len(paths), path)

		info, err := os.Stat(path)
		if err != nil {
			fail(err)
		}
		source, err := os.ReadFile(path)
		if err != nil {
			fail(err)
		}

		result, err := formatSource(path, source)
		if err != nil {
			fail(err)
		}

		if !bytes.Equal(source, result) {
			changes = append(changes, fileChange{
				path: path,
				data: result,
				mode: info.Mode().Perm(),
			})
		}
	}

	// Stage all results first so processing errors never leave partial edits.
	var staged []string
	cleanup := func() {
		for _, path := range staged {
			_ = os.Remove(path)
		}
	}

	for _, change := range changes {
		tmp, err := os.CreateTemp(filepath.Dir(change.path), ".formatgo-*.tmp")
		if err != nil {
			cleanup()
			fail(err)
		}

		tmpPath := tmp.Name()
		if _, err := tmp.Write(change.data); err != nil {
			_ = tmp.Close()
			_ = os.Remove(tmpPath)
			cleanup()
			fail(err)
		}
		if err := tmp.Chmod(change.mode); err != nil {
			_ = tmp.Close()
			_ = os.Remove(tmpPath)
			cleanup()
			fail(err)
		}
		if err := tmp.Close(); err != nil {
			_ = os.Remove(tmpPath)
			cleanup()
			fail(err)
		}
		staged = append(staged, tmpPath)
	}

	for i, change := range changes {
		if err := os.Rename(staged[i], change.path); err != nil {
			cleanup()
			fail(fmt.Errorf("write %s: %w", change.path, err))
		}
		fmt.Println("Formatted:", change.path)
	}

	fmt.Printf("Formatting complete. Files checked: %d, files changed: %d\n",
		len(paths), len(changes))
}

func formatSource(path string, source []byte) ([]byte, error) {
	result := append([]byte(nil), source...)

	// Multiple passes allow nested calls to be compacted without reparsing
	// after every individual edit.
	for pass := 0; pass < 100; pass++ {
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, result, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}
		tf := fset.File(file.Pos())

		var edits []edit

		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || call.Ellipsis.IsValid() || containsFuncLiteral(call) {
				return true
			}

			start := tf.Offset(call.Lparen)
			end := tf.Offset(call.Rparen) + 1
			if !validRange(start, end, len(result)) ||
				!bytes.ContainsAny(result[start:end], "\r\n") {
				return true
			}

			compacted, err := compactDelimited(result[start:end])
			if err != nil {
				return true
			}
			edits = append(edits, edit{start: start, end: end, text: compacted})
			return true
		})

		// Add multiline function declaration parameter lists.
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			start := tf.Offset(fn.Type.Func)
			end := tf.Offset(fn.Type.Params.Closing) + 1
			if !validRange(start, end, len(result)) ||
				!bytes.ContainsAny(result[start:end], "\r\n") {
				continue
			}
			compacted, err := compactDelimited(result[start:end])
			if err != nil {
				return nil, fmt.Errorf("declaration in %s: %w", path, err)
			}
			edits = append(edits, edit{start: start, end: end, text: compacted})
		}

		if len(edits) == 0 {
			break
		}

		// Apply edits right to left so earlier offsets remain valid.
		sort.Slice(edits, func(i, j int) bool {
			if edits[i].start != edits[j].start {
				return edits[i].start > edits[j].start
			}
			return edits[i].end > edits[j].end
		})

		var lastStart int = len(result) + 1
		changed := false
		for _, e := range edits {
			// Skip overlapping edits; a later pass handles outer expressions.
			if e.end > lastStart {
				continue
			}
			if !validRange(e.start, e.end, len(result)) {
				continue
			}
			if bytes.Equal(result[e.start:e.end], e.text) {
				lastStart = e.start
				continue
			}
			result = replaceRange(result, e.start, e.end, e.text)
			lastStart = e.start
			changed = true
		}

		if !changed {
			break
		}
	}

	if _, err := parser.ParseFile(token.NewFileSet(), path, result, parser.ParseComments); err != nil {
		return nil, fmt.Errorf("invalid output for %s: %w", path, err)
	}

	result, err := format.Source(result)
	if err != nil {
		return nil, fmt.Errorf("gofmt %s: %w", path, err)
	}
	if _, err := parser.ParseFile(token.NewFileSet(), path, result, parser.ParseComments); err != nil {
		return nil, fmt.Errorf("invalid formatted output for %s: %w", path, err)
	}

	return result, nil
}

func replaceRange(src []byte, start, end int, replacement []byte) []byte {
	out := make([]byte, 0, len(src)-(end-start)+len(replacement))
	out = append(out, src[:start]...)
	out = append(out, replacement...)
	out = append(out, src[end:]...)
	return out
}

func validRange(start, end, length int) bool {
	return start >= 0 && start < end && end <= length
}

func compactDelimited(src []byte) ([]byte, error) {
	open := bytes.IndexByte(src, '(')
	if open < 0 {
		return nil, fmt.Errorf("opening parenthesis not found")
	}

	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	var s scanner.Scanner
	var scanErrors []string
	s.Init(file, src, func(_ token.Position, msg string) {
		scanErrors = append(scanErrors, msg)
	}, scanner.ScanComments)

	depth, closeOffset := 0, -1
	for {
		pos, tok, _ := s.Scan()
		if tok == token.EOF {
			break
		}
		offset := file.Offset(pos)
		if offset < open {
			continue
		}
		switch tok {
		case token.LPAREN:
			depth++
		case token.RPAREN:
			depth--
			if depth == 0 {
				closeOffset = offset
			}
		}
		if closeOffset >= 0 {
			break
		}
	}
	if closeOffset < 0 {
		return nil, fmt.Errorf("closing parenthesis not found")
	}

	inner, err := compactTokens(src[open+1 : closeOffset])
	if err != nil {
		return nil, err
	}
	inner = bytes.TrimSpace(inner)
	inner = bytes.TrimSuffix(inner, []byte(","))
	inner = bytes.TrimSpace(inner)

	var out bytes.Buffer
	out.Write(src[:open+1])
	if len(inner) > 0 {
		out.WriteByte(' ')
		out.Write(inner)
		out.WriteByte(' ')
	}
	out.WriteByte(')')
	out.Write(src[closeOffset+1:])
	return out.Bytes(), nil
}

func compactTokens(src []byte) ([]byte, error) {
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	var s scanner.Scanner
	var scanErrors []string
	s.Init(file, src, func(_ token.Position, msg string) {
		scanErrors = append(scanErrors, msg)
	}, scanner.ScanComments)

	var out bytes.Buffer
	var previous token.Token
	var previousLit string

	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}

		if tok == token.COMMENT {
			if strings.HasPrefix(lit, "//") {
				return nil, fmt.Errorf("line comment prevents safe compaction")
			}
			if out.Len() > 0 {
				out.WriteByte(' ')
			}
			out.WriteString(lit)
			out.WriteByte(' ')
			previous = token.COMMENT
			previousLit = lit
			continue
		}

		if tok == token.SEMICOLON {
			// Scanner-generated semicolons use "\n"; discard only those.
			if lit == "\n" {
				continue
			}
			return nil, fmt.Errorf("explicit semicolon prevents safe compaction")
		}

		if out.Len() > 0 && needsTokenSpace(previous, previousLit, tok, lit) {
			out.WriteByte(' ')
		}
		if lit != "" {
			out.WriteString(lit)
		} else {
			out.WriteString(tok.String())
		}
		previous = tok
		previousLit = lit
	}

	if len(scanErrors) > 0 {
		return nil, fmt.Errorf("cannot compact safely: %s", strings.Join(scanErrors, "; "))
	}
	return bytes.TrimSpace(out.Bytes()), nil
}

func needsTokenSpace(prev token.Token, prevLit string, current token.Token, currentLit string) bool {
	if prev == token.COMMENT {
		return true
	}

	isWordToken := func(tok token.Token) bool {
		switch tok {
		case token.IDENT, token.INT, token.FLOAT, token.IMAG,
			token.CHAR, token.STRING:
			return true
		default:
			return false
		}
	}
	if isWordToken(prev) && isWordToken(current) {
		return true
	}

	// Prevent neighboring operators from merging into a different token.
	if prevLit != "" && currentLit != "" {
		combined := prevLit + currentLit
		fset := token.NewFileSet()
		file := fset.AddFile("", fset.Base(), len(combined))
		var s scanner.Scanner
		s.Init(file, []byte(combined), nil, 0)
		_, first, _ := s.Scan()
		_, second, _ := s.Scan()
		if first != prev || second != current {
			return true
		}
	}
	return false
}

func containsFuncLiteral(node ast.Node) bool {
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		if _, ok := n.(*ast.FuncLit); ok {
			found = true
			return false
		}
		return !found
	})
	return found
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
