package Service

import (
	"fmt"
	"github.com/segmentio/ksuid"
	"os"
	"os/exec"
	"path"
)

func VideoConverter(filepath string) (out string) {
	// 设置视频源文件路径
	inputFile := filepath
	// 设置转码后文件路径
	filename := ksuid.New().String()
	outfile := path.Join("out", filename+".mp4")
	outpath := outfile
	dir := path.Dir(outfile)
	os.MkdirAll(dir, os.FileMode(0755))
	// 设置 ffmpeg 命令行参数
	cmd := exec.Command("ffmpeg/bin/ffmpeg", "-i", inputFile, "-profile:v", "main", "-movflags", "+faststart", "-crf", "26", "-y", outpath)

	err := cmd.Run()
	if err != nil {
		fmt.Println("转码失败")
		return ""
	}
	fmt.Println("转码成功")
	return filename
}
