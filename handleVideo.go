package main

import (
	"bytes"
	"fmt"

	"log"
	"os/exec"
)

func main() {
	fmt.Println(GetFrame(1))
}

func GetFrame(index int) string {
	filename := "http://a.vskitcdn.com/V/4a7b36bc-bb67-49a8-aed8-9a0347392798.mp4"
	//outPut := "/Users/huangcheng/Downloads/123456.jpg"
	outPut := "videoCoverImage/123456.jpg"

	// cmd := exec.Command("ffmpeg", "-i", filename, "-vframes", strconv.Itoa(index), "-s", fmt.Sprintf("%dx%d", width, height), "-f", "singlejpeg", "-")
	//_ = exec.Command("mkdir cover_image")
	cmd := exec.Command("ffmpeg", "-i", filename, "-r", "5", "-vframes", "1", "-f", "image2", outPut)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return "bad"
	}
	fmt.Printf("command output: %q", out.String())
	return fmt.Sprintf("command output: %q", out.String())

	//buf := new(bytes.Buffer)
	//
	//cmd.Stdout = buf
	//
	//if cmd.Run() != nil {
	//	panic("could not generate frame")
	//}
	//
	//return buf
}
