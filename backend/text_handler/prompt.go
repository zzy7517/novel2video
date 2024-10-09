package text_handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"novel2video/backend"
	"novel2video/backend/llm"
	"novel2video/backend/util"
)

var sys = `
#Task: #
从输入中提取画面信息

#Rules:#
1. 提取画面感强烈的元素，输出的内容是一组词语
2. 优先输出画面感强烈的名词，比如“谁知道呢，或许做了什么亏心事，惹得神灵降怒了吧…”这句输入中，神灵的画面感最强，则输出神灵
3. 输出画面感强的动词，比如开车，不输出画面感不强的动词，比如嘲讽
4. 不要输出形容词
5. 关于文本中描述的场景，可以适当结合上下文发散
6. 不要输出心理描写，不要输出情绪
7. 避免使用模糊或不明确的描述
8. 每个输入一定要有一个输出
9. 如果出现了人，需要具体到名字
10. 如果没有具体的名字，先结合上下文推断，如果还是没有，则不输出名字
11. 如果没有出现人，可以只描述下场景
12. 如果句子中不是很具体的内容，可以挑选文本中某个词语发散一下
13. 如果某一行输入无法提取出内容，输出“无”，但不要不输出

#example#
input
1.炎炎八月。
2.滴滴滴——！
3.刺耳的蝉鸣混杂着此起彼伏的鸣笛声，回荡在人流湍急的街道上，灼热的阳光炙烤着灰褐色的沥青路面，热量涌动，整个街道仿佛都扭曲了起来。
4.路边为数不多的几团树荫下，几个小年轻正簇在一起，叼着烟等待着红绿灯。
5.突然，一个正在吞云吐雾的小年轻似乎是发现了什么，轻咦了一声，目光落在了街角某处。
6.“阿诺，你在看什么？”他身旁的同伴问道。
7.那个名为阿诺的年轻人呆呆的望着街角，半晌才开口，“你说……盲人怎么过马路？”
8.同伴一愣，迟疑了片刻之后，缓缓开口：“一般来说，盲人出门都有人照看，或者导盲犬引导，要是在现代点的城市的话，马路边上也有红绿灯的语音播报，实在不行的话，或许能靠着声音和导盲杖一点点挪过去？”
9.阿诺摇了摇头，“那如果即没人照看，又没导盲犬，也没有语音播报，甚至连导盲杖都用来拎花生油了呢？”
10. 中年男子话刚刚脱口，便是不出意外的在人头汹涌的广场上带起了一阵嘲讽的骚动。
11. 当初的少年，自信而且潜力无可估量，不知让得多少少女对其春心荡漾，当然，这也包括以前的萧媚。
12. “唉…”莫名的轻叹了一口气，萧媚脑中忽然浮现出三年前那意气风发的少年，四岁练气，十岁拥有九段斗之气，十一岁突破十段斗之气，成功凝聚斗之气旋，一跃成为家族百年之内最年轻的斗者！

output
1. 夏天,炎热的街道
2. 街道,很多汽车
3. 街道,很多汽车
4. 红绿灯,马路,几个年轻人
5. 阿诺,向远处看
6. 阿诺,街角
6. 阿诺,盲人过马路
7. 导盲犬,盲人
8. 导盲犬,红绿灯,盲人
9. 阿诺
10. 广场,人群
11. 少女,喜欢
13. 萧媚,叹气

#Input Format:#
0.输入0
1.输入1
2.输入2
每个数字开头的行代表一个输入，每行输入必须对应一行输出

#Output Format:#
每行输入需要对应一行输出，每个输出用空行隔开
输入和输出的序号需要对应
不要输出多余内容
这是一个输出示例
0. 输出0
1. 输出1
2. 输出2

# 检查 #
输出的行数是否和输入一致，如果不一致，则重新生成输出内容
`

var fragmentsLen = 30

func ExtractSceneFromTexts(c *gin.Context) {
	lines, err := util.ReadLinesFromDirectory(util.NovelFragmentsDir)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to read fragments", err)
		return
	}
	err = os.RemoveAll(util.PromptsDir)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to remove directory", err)
		return
	}
	err = os.MkdirAll(util.PromptsDir, os.ModePerm)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to create directory", err)
		return
	}
	var offset int
	promptsMid := generateInputPrompts(lines, fragmentsLen)
	for _, p := range promptsMid {
		res, err := llm.QueryLLM(c.Request.Context(), p, sys, "doubao", 1, 8192)
		if err != nil {
			logrus.Errorf("query gemini failed, err is %v", err)
			continue
		}
		lines := strings.Split(res, "\n")
		// 编译正则表达式，用于匹配行首的序号和点
		re := regexp.MustCompile(`^\d+\.\s*`)
		// 遍历每一行，去掉序号和点
		var t2iPrompts []string
		for _, line := range lines {
			line = re.ReplaceAllString(line, "")
			if len(strings.TrimSpace(line)) > 0 {
				t2iPrompts = append(t2iPrompts, line)
				offset++
			}
		}

		err = saveListToFiles(t2iPrompts, util.PromptsDir+"/", offset-len(t2iPrompts))
		if err != nil {
			backend.HandleError(c, http.StatusInternalServerError, "save list to file failed", err)
			return
		}
	}
	// 从目录中读取所有文件并返回内容
	lines, err = util.ReadLinesFromDirectory(util.PromptsDir)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to read fragments", err)
		return
	}
	logrus.Infof("extract prompts from novel fragments finished")
	c.JSON(http.StatusOK, lines)
}

// 每100个生成为一组，发给ai，提取提示词
func generateInputPrompts(list []string, step int) []string {
	var prompts []string
	for i := 0; i < len(list); i += step {
		end := i + step
		if end > len(list) {
			end = len(list)
		}
		var builder strings.Builder
		for j, v := range list[i:end] {
			builder.WriteString(strconv.Itoa(j) + ". ")
			builder.WriteString(fmt.Sprintf("%v", v))
			if j != len(list[i:end])-1 {
				builder.WriteString("\n")
			}
		}
		prompt := builder.String()
		logrus.Infof("prompt is %v", prompt)
		prompts = append(prompts, prompt)
	}
	return prompts
}

func GetPromptsEn(c *gin.Context) {
	err := os.RemoveAll(util.PromptsEnDir)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to remove directory", err)
		return
	}
	err = os.MkdirAll(util.PromptsEnDir, os.ModePerm)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to create directory", err)
		return
	}
	lines, err := util.ReadLinesFromDirectory(util.PromptsDir)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to read fragments", err)
		return
	}
	characterMap, err := getLocalCharactersMap(c.Request.Context())
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to get local characters", err)
		return
	}
	for i, _ := range lines {
		for key, value := range characterMap {
			if strings.Contains(lines[i], strings.TrimSpace(key)) {
				lines[i] = strings.ReplaceAll(lines[i], strings.TrimSpace(key), value)
			}
		}
	}
	lines, err = translatePrompts(lines)
	if err != nil {
		logrus.Errorf("translate prompts failed, err %v", err)
		backend.HandleError(c, http.StatusInternalServerError, "translate failed", err)
		return
	}
	err = saveListToFiles(lines, util.PromptsEnDir+"/", 0)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to save promptsEn", err)
		return
	}
	logrus.Infof("translate prompts to English finished")
	c.JSON(http.StatusOK, lines)
}

func translatePrompts(lines []string) (res []string, err error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := range lines {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			newValue, err := llm.LLMTranslate(context.Background(), lines[i])
			var cnt int
			for cnt < 3 && err != nil {
				logrus.Errorf("translated faild, err %v", err)
				newValue, err = llm.LLMTranslate(context.Background(), lines[i])
				cnt++
			}
			if err != nil {
				logrus.Errorf("translated faild, err %v", err)
				return
			}
			mu.Lock()
			lines[i] = newValue
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	return lines, nil
}

func SavePromptEn(c *gin.Context) {
	type SaveRequest struct {
		Index   int    `json:"index"`
		Content string `json:"content"`
	}
	var req SaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		backend.HandleError(c, http.StatusBadRequest, "parse request body failed", err)
		return
	}

	filePath := filepath.Join(util.PromptsEnDir, fmt.Sprintf("%d.txt", req.Index))
	if err := os.MkdirAll(util.PromptsEnDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}
	if err := os.WriteFile(filePath, []byte(req.Content), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Attachment saved successfully"})
}
