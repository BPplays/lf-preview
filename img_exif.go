package main

import "sync"


func image_exif(image_file string, width, height int, file string, tags [][]string, thumbnail_type string) (string) {
	output := ""

	ch := make(chan order_string)
	// ch2 := make(chan string)

	// var gr_array [2]string


	var wg sync.WaitGroup

	wg.Add(2)
	go image_gr(image_file, width, height, ch, 0, &wg, thumbnail_type)
	go exif_fmt_gr(file, tags, ch, 1, &wg)


	go func() {
		wg.Wait()
		close(ch)
		// close(ch2)
	}()

	// gr_array[0] = "test0"
	// gr_array[1] = "test1"
	// output = output + fmt.Sprintln(gr_array[0])
	var temp_slice [20]string

	lines := 0
	for result := range ch {
		lines += countRune(result.content, '\n')
		temp_slice[result.order] = result.content
	}



	if lines > height {




		new_height := height-countRune(temp_slice[1], '\n')

		if new_height < 2 {
			temp_slice[0] = ""
		} else {
			ch2 := make(chan order_string)

			wg.Add(1)
			go image_gr(image_file, width, new_height, ch2, 0, &wg, thumbnail_type)
			go func() {
				wg.Wait()
				close(ch2)
			}()

			for result := range ch2 {
				temp_slice[result.order] = result.content
			}
		}
		// fmt.Println(lines)
		// fmt.Printf("h: %v, nh: %v\n", height, new_height)


	}



	for _, val := range temp_slice {
		output = output + val
	}

	// output = output + fmt.Sprintln(sep1)
	// output = output + fmt.Sprintln(sep1)
	// output = output + fmt.Sprintln(gr_array[1])

	// close(ch)
	// close(ch2)

	return output
}