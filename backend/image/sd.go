package image

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"novel2video/backend/util"
)

var url = "http://10.193.239.248:7860"

func GenerateImage(prompt string, seed int, width int, height int, order int) error {
	payload := map[string]interface{}{
		"prompt":          "anime " + prompt + " <lora:The Garden of Words_20230619154444:1>",
		"negative_prompt": "(painting by bad-artist-anime:0.9), (painting by bad-artist:0.9), watermark, text, error, blurry, jpeg artifacts, cropped, worst quality, low quality, normal quality, jpeg artifacts, signature, watermark, username, artist name,deformed,distorted,disfigured,doll,poorly drawn,bad anatomy,wrong anatomy,bad hand,bad fingers,NSFW",
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

	outputFilename := fmt.Sprintf("%v/%d.png", util.ImageDir, order)

	if err := os.WriteFile(outputFilename, imageData, 0644); err != nil {
		return fmt.Errorf("failed to write image file: %v", err)
	}

	logrus.Infof("Image saved to %v", outputFilename)
	return nil
}
