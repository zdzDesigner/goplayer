package main

// func ctrl(end chan struct{}, ch chan struct{}, currIndex chan int, names []string) {
// 	var (
// 		order string
// 		cmds  = []string{"next", "pre", "del!", "jump", "exit"}
// 	)
// 	for {
// 		fmt.Scanf("%s", &order)
// 		fmt.Println("order::", order, "expect::", strings.Join(cmds, " | "))
// 		if util.Contains(cmds, order) {
// 			ch <- struct{}{}
// 			// time.Sleep(time.Second)
// 			i := <-currIndex
// 			if order == "del!" {
// 				delSong(names[i])
// 			}
// 			switch order {
// 			case "next":
// 			case "del!":
// 				i++
// 			case "pre":
// 				i--
// 			case "jump":
// 				i += 10
// 			case "exit":
// 				end <- struct{}{}
// 			default:
//
// 			}
//
// 			order = ""
// 			ch = make(chan struct{}, 1)
//
// 			// go play(names[1], ch)
// 			go func() {
// 				currIndex <- control(ch, i, names)
// 				fmt.Println("currIndex::", i)
// 			}()
// 		}
//
// 	}
// }
//
// // 删除资源
// func delSong(name string) {
// 	var err error
// 	defer func() {
// 		if err != nil {
// 			fmt.Println("del song err::", err)
// 		}
// 	}()
//
// 	f, list, err := getIgnoreDetail()
// 	list = append(list, name)
// 	bts, err := json.Marshal(list)
// 	if err != nil {
// 		return
// 	}
// 	fmt.Println("bts::", string(bts))
// 	_, err = f.Write(bts)
// 	if err != nil {
// 		return
// 	}
// }
