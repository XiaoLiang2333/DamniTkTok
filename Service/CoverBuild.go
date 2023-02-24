package Service

import (
	"fmt"
	"github.com/segmentio/ksuid"
	"os"
	"os/exec"
	"path"
)

func Cover(filepath string) (covername string) {
	inputFile := filepath
	filename := ksuid.New().String()
	outfile := path.Join("cover", filename+".jpg")
	dir := path.Dir(outfile)
	os.MkdirAll(dir, os.FileMode(0755))
	// 执行 FFmpeg 命令进行截取第一帧图片
	cmd := exec.Command("ffmpeg/bin/ffmpeg", "-ss", "00:00:01.000", "-i", inputFile, "-vframes", "1", outfile)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("截取图片失败：%v\n", err)
		return ""
	}
	fmt.Println("图片截取完成")
	return filename
}
