package ascii

import (
	"fmt"
	"strings"
)

// NOTE: here's my own repo for reference -> https://github.com/Mr-Robot-err-404/tui-spinner

const (
	Base rune = 0x2800
)

type Grid struct {
	cells  [][]bool
	offset int
}

type Coord struct {
	row int
	col int
}

type Rotation struct {
	head map[Coord]Coord
	tail map[Coord]Coord
}
type SpinnerTheme struct {
	Inner string
	Outer string
}
type Centre struct {
	m      map[Coord]bool
	bounds [2]Coord
}

type Spinner struct {
	size     int
	grid     Grid
	centre   Centre
	head     []Coord
	tail     []Coord
	rotation Rotation
	theme    SpinnerTheme
	offset   int
}

var braiile_map = map[Coord]byte{
	{col: 0, row: 0}: 0,
	{col: 0, row: 1}: 1,
	{col: 0, row: 2}: 2,
	{col: 0, row: 3}: 6,
	{col: 1, row: 0}: 3,
	{col: 1, row: 1}: 4,
	{col: 1, row: 2}: 5,
	{col: 1, row: 3}: 7,
}

func (s Spinner) Walk() Spinner {
	s.head = step(s.head, s.rotation.head)

	for _, pos := range s.head {
		s.grid.set(pos.row, pos.col, true)
	}
	for _, pos := range s.tail {
		s.grid.set(pos.row, pos.col, false)
	}
	s.tail = step(s.tail, s.rotation.tail)
	return s
}

func (sp Spinner) Render_frame() string {
	height := (len(sp.grid.cells) + 3) / 4
	width := (len(sp.grid.cells[0]) + 1) / 2

	screen := make([][]byte, height)

	for i := range screen {
		screen[i] = make([]byte, width)
	}
	active := sp.theme.Outer
	s := strings.Builder{}
	s.WriteString(active)

	for row, current := range sp.grid.cells {
		for col := range current {
			if !sp.grid.cells[row][col] {
				continue
			}
			i := row / 4
			j := col / 2
			n := braiile_map[Coord{row: row % 4, col: col % 2}]
			screen[i][j] |= 1 << n
		}
	}
	for i, line := range screen {
		for j, b := range line {
			r := Base + rune(b)
			s.WriteRune(r)

			if should_switch(sp.centre.bounds, i, j) {
				active = switch_theme(sp.theme, active)
				s.WriteString(Reset)
				s.WriteString(active)
			}
		}
		if i != len(screen)-1 {
			s.WriteString("\n")
		}
	}
	s.WriteString(Reset)
	return s.String()
}

func MakeSpinner(size int, theme SpinnerTheme) Spinner {
	dm := calc_dimension(size)
	offset := vertical_offset(size)

	cells := make([][]bool, dm+offset)
	for i := range cells {
		cells[i] = make([]bool, dm)
	}
	grid := Grid{cells: cells, offset: offset}

	centre_map, start, end := make_centre(size, dm)
	bounds := [2]Coord{
		{row: (start.row + grid.offset) / 4, col: (start.col / 2) - 1},
		{row: (end.row + grid.offset) / 4, col: end.col / 2},
	}
	head, tail := key_nodes(dm, size)

	for i := range head {
		grid.fill(tail[i], head[i])
	}
	for c := range centre_map {
		grid.set(c.row, c.col, true)
	}
	width, height := len(cells)-offset, len(cells[0])

	return Spinner{
		size:   size,
		grid:   Grid{cells: cells, offset: offset},
		head:   head,
		tail:   tail,
		centre: Centre{m: centre_map, bounds: bounds},
		rotation: Rotation{
			tail: make_tail_map(width, height, size),
			head: make_head_map(width, height, size),
		},
		theme:  normalize_theme(theme.Inner, theme.Outer),
		offset: offset,
	}
}

func calc_dimension(size int) int {
	low := 8
	x := 5 * (size - 2)
	return low + x
}

func vertical_offset(size int) int {
	if size == 2 {
		return 2
	}
	return 0
}

func normalize_theme(inner string, outer string) SpinnerTheme {
	return SpinnerTheme{
		Inner: normalize(inner),
		Outer: normalize(outer),
	}
}

func normalize(s string) string {
	if strings.HasPrefix(s, OpenSequence) {
		s = fmt.Sprintf("%s%s", OpenSequence, s)
	}
	idx := strings.Index(s, Reset)

	if idx == -1 {
		return s
	}
	return s[0:idx]
}

func make_centre(size int, width int) (map[Coord]bool, Coord, Coord) {
	centre := make(map[Coord]bool, size*size)

	mid := (width / 2)
	offset := size / 2

	start := Coord{
		row: mid - offset,
		col: mid - offset,
	}
	for i := range size {
		for j := range size {
			centre[Coord{row: start.row + i, col: start.col + j}] = true
		}
	}
	return centre, start, Coord{row: start.row + size - 1, col: start.col + size - 1}
}

func make_head_map(width, height, size int) map[Coord]Coord {
	m := make(map[Coord]Coord)

	end_col := width - 1
	end_row := height - 1

	for n := range size {
		m[Coord{row: n, col: end_col}] = Coord{row: size, col: end_col - n}
	}

	for n := range size {
		m[Coord{row: end_row, col: end_col - n}] = Coord{row: end_row - n, col: end_col - size}
	}

	for n := range size {
		m[Coord{row: end_row - n, col: 0}] = Coord{row: end_col - size, col: n}
	}

	for n := range size {
		m[Coord{row: 0, col: n}] = Coord{row: n, col: size}
	}
	return m
}

func make_tail_map(width, height, size int) map[Coord]Coord {
	m := make(map[Coord]Coord)

	end_col := width - 1
	end_row := height - 1

	for n := range size {
		m[Coord{row: size, col: n}] = Coord{row: n, col: 0}
	}

	for n := range size {
		m[Coord{row: n, col: end_col - size}] = Coord{row: 0, col: end_col - n}
	}

	for n := range size {
		m[Coord{row: end_row - size, col: end_col - n}] = Coord{row: end_row - n, col: end_col}
	}

	for n := range size {
		m[Coord{row: end_row - n, col: size}] = Coord{row: end_row, col: n}
	}
	return m
}

func rotate_nodes(nodes []Coord, rotation map[Coord]Coord) ([]Coord, bool) {
	transform := []Coord{}

	for _, pos := range nodes {
		next, ok := rotation[pos]
		if !ok {
			return transform, false
		}
		transform = append(transform, next)
	}
	return transform, true
}

func x_dir(nodes []Coord) int {
	for _, pos := range nodes {
		if pos.row == 0 {
			return 1
		}
	}
	return -1
}

func y_dir(nodes []Coord) int {
	for _, pos := range nodes {
		if pos.col == 0 {
			return -1
		}
	}
	return 1
}

func traversing_x(nodes []Coord) bool {
	prev := nodes[0]

	for i := 1; i < len(nodes); i++ {
		curr := nodes[i]

		if curr.col != prev.col {
			return false
		}
		prev = curr
	}
	return true
}
func traversing_y(nodes []Coord) bool {
	prev := nodes[0]

	for i := 1; i < len(nodes); i++ {
		curr := nodes[i]

		if curr.row != prev.row {
			return false
		}
		prev = curr
	}
	return true
}

func step(nodes []Coord, rotate map[Coord]Coord) []Coord {
	next, ok := rotate_nodes(nodes, rotate)
	if ok {
		return next
	}
	if traversing_x(nodes) {
		dir := x_dir(nodes)

		for i := range nodes {
			nodes[i].col += dir
		}
	}
	if traversing_y(nodes) {
		dir := y_dir(nodes)

		for i := range nodes {
			nodes[i].row += dir
		}
	}
	return nodes
}

func switch_theme(theme SpinnerTheme, active string) string {
	if theme.Outer == active {
		return theme.Inner
	}
	return theme.Outer
}

func should_switch(bounds [2]Coord, row int, col int) bool {
	if row >= bounds[0].row && row <= bounds[1].row {
		return col == bounds[0].col || col == bounds[1].col
	}
	return false
}

func full_square(grid Grid) Grid {
	for i := range len(grid.cells) {
		for j := range len(grid.cells[0]) {
			grid.set(i, j, true)
		}
	}
	return grid
}

func (grid *Grid) get(row, col int) bool {
	return grid.cells[row+grid.offset][col]
}
func (grid *Grid) set(row, col int, value bool) {
	grid.cells[row+grid.offset][col] = value
}

func (grid *Grid) fill(start, end Coord) {
	x, y := 1, 1
	row, col := start.row, start.col
	grid.set(row, col, true)

	if end.row < start.row {
		y = -1
	}
	if end.col < start.col {
		x = -1
	}
	for row != end.row {
		row += y
		grid.set(row, col, true)
	}
	for col != end.col {
		col += x
		grid.set(row, col, true)
	}
}

func key_nodes(width int, size int) ([]Coord, []Coord) {
	head := []Coord{}
	tail := []Coord{}
	rem := (width % 2) + ((size - 2) / 2)
	mid := (width / 2) + rem

	for n := range size {
		head = append(head, Coord{row: n, col: mid})
		tail = append(tail, Coord{row: mid, col: n})
	}
	return head, tail
}
