package aigchq

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestCurrentModelConstants(t *testing.T) {
	cases := map[string]string{
		ModelGeminiAuto:            "gemini-auto",
		ModelGeminiAutoExtended:    "gemini-auto-extended",
		ModelGemini35FlashLite:     "gemini-3.5-flash-lite",
		ModelGemini38Flash:         "gemini-3.8-flash",
		ModelGemini38FlashExtended: "gemini-3.8-flash-extended",
		ModelGemini31Pro:           "gemini-3.1-pro",
		ModelGemini4Argon:          "gemini-4-argon",
		ModelGemini4ArgonExtended:  "gemini-4-argon-extended",
		ModelGemini36Flash:         "gemini-3.6-flash",
		ModelGemini36FlashExtended: "gemini-3.6-flash-extended",
		ModelGemini3ProImage:       "gemini-3-pro-image",
		ModelGemini3FlashImage:     "gemini-3-flash-image",
		ModelChatGPTAuto:           "auto",
		ModelChatGPTGPT6:           "gpt-6",
		ModelChatGPTImage2:         "gpt-image-2",
	}
	for got, want := range cases {
		if got != want {
			t.Fatalf("constant %q, want %q", got, want)
		}
	}
}

func TestChatCompletionRequestMarshalsChatGPTTools(t *testing.T) {
	req := ChatCompletionRequest{
		Model: ModelChatGPTAuto,
		Messages: []Message{{
			Role: "user",
			Content: []ContentPart{{
				Type: "input_file",
				File: &InputFile{
					Data:     "data:text/plain;base64,aGVsbG8=",
					Filename: "note.txt",
					MimeType: "text/plain",
				},
			}},
		}},
		WebSearch:       true,
		ReasoningEffort: ReasoningEffortMedium,
		ImageGeneration: true,
		DeepResearch:    true,
	}
	raw, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	text := string(raw)
	for _, field := range []string{
		`"model":"auto"`,
		`"web_search":true`,
		`"reasoning_effort":"medium"`,
		`"image_generation":true`,
		`"deep_research":true`,
		`"type":"input_file"`,
		`"filename":"note.txt"`,
	} {
		if !strings.Contains(text, field) {
			t.Fatalf("marshaled request missing %s: %s", field, text)
		}
	}
}
