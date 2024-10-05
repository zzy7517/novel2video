package image

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

var url = "http://10.193.239.248:7860/"

func generateImage(prompt string, seed int, width int, height int, order int) error {
	payload := map[string]interface{}{
		"prompt":          "anime" + prompt,
		"negative_prompt": "booty, boob, (nsfw), (painting by bad-artist-anime:0.9), (painting by bad-artist:0.9), watermark, text, error, blurry, jpeg artifacts, cropped, worst quality, low quality, normal quality, jpeg artifacts, signature, watermark, username, artist name, (worst quality, low quality:1.4), bad anatomy",
		"cfg_scale":       7,
		"steps":           30,
		"width":           width,
		"height":          height,
		"override_settings": map[string]string{
			"sd_vae": "Automatic",
		},
		"enable_hr":            true,
		"denoising_strength":   0.7,
		"hr_upscaler":          "Latent",
		"hr_resize_x":          1024,
		"hr_resize_y":          1024,
		"hr_sampler_name":      "Euler",
		"hr_second_pass_steps": 28,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	resp, err := http.Post(url+"/sdapi/v1/txt2img", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}

	images, ok := response["images"].([]interface{})
	if !ok || len(images) == 0 {
		return fmt.Errorf("no images found in response")
	}

	imageData, err := base64.StdEncoding.DecodeString(images[0].(string))
	if err != nil {
		return fmt.Errorf("failed to decode image: %v", err)
	}

	outputFilename := fmt.Sprintf("temp/images/%d.png", order)

	if err := ioutil.WriteFile(outputFilename, imageData, 0644); err != nil {
		return fmt.Errorf("failed to write image file: %v", err)
	}

	logrus.Infof("Image saved to", outputFilename)
	return nil
}

func main() {
	err := generateImage("(Best Quality), a boy, Anime, sitting, eating, ((masterpiece)) <lora:ChosenChineseStyleNsfw_v20:1>", 114514191981, 540, 960, 3)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
