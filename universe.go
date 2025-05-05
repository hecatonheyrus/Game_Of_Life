package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"os"
)

/*
   Flat two-dimensional universe where Life plays its Game...
*/

type Universe struct {
	CellsUpdated int
	Width        int
	Height       int
	Cells        [](*Cell)
}

/*
  Cell can be dead (0) or alive (1)
*/

type Cell struct {
	X              int
	Y              int
	State          int
	LiveNeighbours int
}

/*
Rules of Life:
1 - if the cell is alive and live neighbours count < 2 => underpopulation => death
2-  if the cell is alive and live neighbours count > 3 => overpopulation => death
3 - if the cell is dead and live neigbours count = 3 => reproduction
*/
func (c *Cell) UpdateCellStatus() int {
	updated := 0

	if c.State == 1 {
		if c.LiveNeighbours < 2 ||
			c.LiveNeighbours > 3 {
			c.State = 0
			updated = 1
		}
	} else if c.State == 0 {
		if c.LiveNeighbours == 3 {
			c.State = 1
			updated = 1
		}
	}

	return updated
}

func (u *Universe) GetCellStateByIndex(idx int) int {

	state := u.Cells[idx].State
	return state
}

/*
   Calculates 8 directions encoded as -1, 0, +1
*/

func (u *Universe) Directions(dirs []int) [][]int {

	var result [][]int

	for j := 0; j < len(dirs); j++ {
		for i := 0; i < len(dirs); i++ {
			if dirs[i] == 0 && dirs[j] == 0 {
				continue
			}

			dir := []int{dirs[j], dirs[i]}
			result = append(result, dir)
		}
	}

	return result
}

/*
   For a given cell, calculates number of its live neighbours
*/

func (c *Cell) GetLiveNeighbours(u *Universe) {

	live_neigbours := 0

	d := []int{0, 1, -1}

	directions := u.Directions(d)

	for _, d := range directions {
		index, error := u.GetIndex(c.X+d[0], c.Y+d[1])
		if error == nil {
			live_neigbours += u.GetCellStateByIndex(index)
		}
	}

	c.LiveNeighbours = live_neigbours

}

/*
   Creates a new universe given width(w) and height(h)
*/

func newUniverse(w, h int) *Universe {

	u := Universe{Width: w, Height: w, CellsUpdated: 0}

	for i := 0; i < u.Height; i++ {
		for j := 0; j < u.Width; j++ {
			c := &Cell{X: i, Y: j, State: 0, LiveNeighbours: 0}
			if j > u.Width/2 && i < u.Height/2 {
				if i%2 == 0 || i%7 == 0 {
					c.State = 1
				}
			}
			u.Cells = append(u.Cells, c)
		}
	}

	return &u

}

func (u *Universe) GetIndex(row, column int) (int, error) {

	if row < 0 || row >= u.Width {
		return -1, fmt.Errorf("Row index %d out of range", row)
	}

	if column < 0 || column >= u.Height {
		return -1, fmt.Errorf("Column index %d out of range", row)
	}

	return row*u.Width + column, nil
}

/*
One tick of the universe
*/
func (u *Universe) Tick() int {

	for i := 0; i < u.Height; i++ {
		for j := 0; j < u.Width; j++ {

			idx, error := u.GetIndex(i, j)
			if error != nil {
				continue
			}
			cell := u.Cells[idx]
			cell.GetLiveNeighbours(u)
			cells_updated := cell.UpdateCellStatus()
			u.CellsUpdated += cells_updated

		}
	}
	return u.CellsUpdated

}

/*
Prints out current state of universe
For debugging purposes only
*/
func (u *Universe) PrintUniverse() {

	for i, cell := range u.Cells {
		if i%u.Width == 0 {
			fmt.Println("")
		}
		fmt.Printf("%d ", cell.State)

	}
	fmt.Println("-----------------------------------------")
}

/*
   Saves Life to GIF
*/

func (u *Universe) ToGif(GifFile string) {

	var palette = []color.Color{color.White, color.Black}
	const (
		nframes = 512
		delay   = 20
		scale   = 2
	)

	anim := gif.GIF{LoopCount: nframes}

	for i := 0; i < nframes; i++ {

		u.Tick()

		rect := image.Rect(0, 0, u.Width*scale, u.Height*scale)
		img := image.NewPaletted(rect, palette)

		for _, cell := range u.Cells {

			for i := 0; i < scale; i++ {
				img.SetColorIndex(cell.X*scale+i, cell.Y*scale+i, uint8(cell.State))
			}
		}

		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)

	}

	f, err := os.OpenFile(GifFile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()

	gif.EncodeAll(f, &gif.GIF{
		Image: anim.Image,
		Delay: anim.Delay,
	})

}
