package main

/*
   Convey's Game of Life
   @author Vadzim Shandrykau
   ahamemnonvs@gmail.com
*/

func main() {
	universe := newUniverse(128, 128)
	universe.ToGif("gol.gif")
}
