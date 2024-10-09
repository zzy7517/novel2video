package tts

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/wychl/edgetts"

	"novel2video/backend/util"
)

func byEdgeTTS() {
	err := os.RemoveAll(util.AudioDir)
	if err != nil {
		logrus.Errorf("remove dir %v failed, %v", util.AudioDir, err)
		return
	}
	err = os.MkdirAll(util.AudioDir, os.ModePerm)
	if err != nil {
		logrus.Errorf("create dir %v failed, %v", util.AudioDir, err)
		return
	}
	lines, err := util.ReadLinesFromDirectory(util.NovelFragmentsDir)
	if err != nil {
		logrus.Errorf("read novel fragments failed, %v", err)
		return
	}
	cli := edgetts.New(nil)
	for i, line := range lines {
		data, err := cli.TTS(line, "en-US-EmmaMultilingualNeural", edgetts.WithRate("+25%"))
		if err != nil {
			logrus.Errorf("tts convert failed, err: %v", err)
			continue
		}
		fullPath := filepath.Join(util.AudioDir, strconv.Itoa(i)+".mp3")
		err = WriteToFile(fullPath, data)
		if err != nil {
			logrus.Errorf("tts convert failed, err: %v", err)
		}
	}
}

func WriteToFile(file string, data []byte) error {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}
