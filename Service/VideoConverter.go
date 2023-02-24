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
	args := []string{"-i", inputFile, "-vf", "scale=1280:720", "-c:v", "libx264", "-preset", "medium", "-crf", "23", "-c:a", "copy", outpath}
	// 创建 *exec.Cmd
	cmd := exec.Command("ffmpeg/bin/ffmpeg", args...)
	// 运行 ffmpeg 命令
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("转码成功")
	return filename
}
