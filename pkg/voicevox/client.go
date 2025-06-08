package voicevox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Client はVOICEVOX APIとの通信を担当する構造体です
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient は新しいVOICEVOX APIクライアントを作成します
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

// Speaker はVOICEVOXの話者情報を表す構造体です
type Speaker struct {
	Name      string `json:"name"`
	SpeakerID int    `json:"speaker_id"`
	StyleID   int    `json:"style_id"`
	StyleName string `json:"style_name"`
}

// GetSpeakers は利用可能な話者の一覧を取得します
func (c *Client) GetSpeakers() ([]Speaker, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/speakers")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	var speakers []Speaker
	if err := json.NewDecoder(resp.Body).Decode(&speakers); err != nil {
		return nil, err
	}

	return speakers, nil
}

// AudioQueryRequest は音声合成のためのクエリリクエストを表す構造体です
type AudioQueryRequest struct {
	Text      string `json:"text"`
	SpeakerID int    `json:"speaker"`
}

// AudioQuery は音声合成のためのクエリ結果を表す構造体です
type AudioQuery struct {
	AccentPhrases      []interface{} `json:"accent_phrases"`
	SpeedScale         float64       `json:"speedScale"`
	PitchScale         float64       `json:"pitchScale"`
	IntonationScale    float64       `json:"intonationScale"`
	VolumeScale        float64       `json:"volumeScale"`
	PrePhonemeLength   float64       `json:"prePhonemeLength"`
	PostPhonemeLength  float64       `json:"postPhonemeLength"`
	OutputSamplingRate int           `json:"outputSamplingRate"`
	OutputStereo       bool          `json:"outputStereo"`
	Kana               string        `json:"kana"`
}

// CreateAudioQuery はテキストから音声合成のためのクエリを作成します
func (c *Client) CreateAudioQuery(text string, speakerID int) (*AudioQuery, error) {
	params := url.Values{}
	params.Add("text", text)
	params.Add("speaker", fmt.Sprintf("%d", speakerID))

	req, err := http.NewRequest("POST", c.BaseURL+"/audio_query?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s, body: %s", resp.Status, string(body))
	}

	var query AudioQuery
	if err := json.NewDecoder(resp.Body).Decode(&query); err != nil {
		return nil, err
	}

	return &query, nil
}

// SynthesizeVoice は音声合成を実行し、音声データを返します
func (c *Client) SynthesizeVoice(query *AudioQuery, speakerID int) ([]byte, error) {
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("speaker", fmt.Sprintf("%d", speakerID))

	req, err := http.NewRequest("POST", c.BaseURL+"/synthesis?"+params.Encode(), bytes.NewBuffer(queryJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "audio/wav")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}
