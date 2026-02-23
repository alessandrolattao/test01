package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type AIService struct {
	DB          *gorm.DB
	WhisperURL  string
	OllamaURL   string
	OllamaModel string
}

type whisperResponse struct {
	Text string `json:"text"`
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
}

type AIScores struct {
	DescriptiveSkills int `json:"descriptive_skills"`
	CriticalThinking  int `json:"critical_thinking"`
	Engagement        int `json:"engagement"`
	Structure         int `json:"structure"`
	PracticalInfo     int `json:"practical_info"`
}

type AIAnalysisResult struct {
	Summary        string   `json:"summary"`
	Strengths      []string `json:"strengths"`
	Weaknesses     []string `json:"weaknesses"`
	Recommendation string   `json:"recommendation"`
	Scores         AIScores `json:"scores"`
	Score          int      `json:"score"`
}

// TranscribeAndAnalyze runs the full pipeline in a goroutine.
func (s *AIService) TranscribeAndAnalyze(candidateID string, audioPath string) {
	go func() {
		log.Printf("[AI] Starting pipeline for candidate %s", candidateID)

		// Step 1: Transcribe
		s.updateStatus(candidateID, "transcribing")

		transcript, err := s.transcribe(audioPath)
		if err != nil {
			log.Printf("[AI] Transcription failed for %s: %v", candidateID, err)
			s.updateStatus(candidateID, "failed")
			return
		}

		s.DB.Model(&struct{ ID string }{}).
			Table("candidates").
			Where("id = ?", candidateID).
			Update("transcript", transcript)

		log.Printf("[AI] Transcription done for %s (%d chars)", candidateID, len(transcript))

		// Step 2: Analyze
		s.updateStatus(candidateID, "analyzing")

		analysis, err := s.analyze(transcript)
		if err != nil {
			log.Printf("[AI] Analysis failed for %s: %v", candidateID, err)
			s.updateStatus(candidateID, "failed")
			return
		}

		analysisJSON, _ := json.Marshal(analysis)

		s.DB.Table("candidates").
			Where("id = ?", candidateID).
			Updates(map[string]interface{}{
				"ai_analysis":     string(analysisJSON),
				"ai_score":        analysis.Score,
				"analysis_status": "completed",
			})

		log.Printf("[AI] Pipeline completed for %s (score: %d)", candidateID, analysis.Score)
	}()
}

func (s *AIService) updateStatus(candidateID string, status string) {
	s.DB.Table("candidates").
		Where("id = ?", candidateID).
		Update("analysis_status", status)
}

func (s *AIService) transcribe(audioPath string) (string, error) {
	file, err := os.Open(audioPath)
	if err != nil {
		return "", fmt.Errorf("opening audio file: %w", err)
	}
	defer func() { _ = file.Close() }()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("audio_file", "audio.webm")
	if err != nil {
		return "", fmt.Errorf("creating form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return "", fmt.Errorf("copying audio data: %w", err)
	}
	_ = writer.Close()

	url := s.WhisperURL + "/asr?language=it&output=json"
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("whisper request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("whisper returned %d: %s", resp.StatusCode, string(body))
	}

	var result whisperResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decoding whisper response: %w", err)
	}

	return strings.TrimSpace(result.Text), nil
}

func (s *AIService) analyze(transcript string) (*AIAnalysisResult, error) {
	prompt := fmt.Sprintf(`/no_think
Sei un recruiter esperto che seleziona reviewer di location turistiche (ristoranti, attrazioni, hotel, esperienze di viaggio).

Al candidato e' stato chiesto di: presentarsi brevemente e raccontare un luogo turistico visitato, descrivendolo come una recensione (cosa lo ha colpito, punti di forza, cosa migliorebbe, se lo consiglierebbe).

REGOLA IMPORTANTE: Se il candidato non parla di un luogo turistico o non fa nessuna descrizione/recensione, tutti i punteggi devono essere 0. Sii severo su questo.

Valuta la trascrizione su queste 5 categorie:
- descriptive_skills (0-25): Capacita' descrittiva. Usa aggettivi, crea immagini vivide, fa venire voglia di andarci? "Era bello" = punteggio basso. Descrizioni sensoriali e dettagliate = punteggio alto.
- critical_thinking (0-25): Occhio critico. Sa bilanciare pro e contro? Nota dettagli rilevanti per un turista? Solo lodi senza critiche = punteggio medio.
- engagement (0-20): Coinvolgimento. Il tono e' appassionato e naturale o monotono e meccanico? Si percepisce entusiasmo per il settore?
- structure (0-15): Struttura del discorso. Presentazione organizzata logicamente? C'e' un filo conduttore o salta da un concetto all'altro?
- practical_info (0-15): Info pratiche utili. Menziona dettagli concreti (prezzi, orari, come arrivare, per chi e' adatto, consigli)?

Trascrizione:
"""%s"""

Rispondi SOLO con un JSON valido, nient'altro:
{
  "summary": "riassunto della presentazione in 2-3 frasi",
  "strengths": ["punto di forza 1", "punto di forza 2"],
  "weaknesses": ["punto debole 1", "punto debole 2"],
  "recommendation": "raccomandazione finale per il recruiter in 1-2 frasi",
  "scores": {
    "descriptive_skills": 0,
    "critical_thinking": 0,
    "engagement": 0,
    "structure": 0,
    "practical_info": 0
  },
  "score": 0
}

Il campo "score" e' la somma dei 5 sub-score (max 100). Rispondi SOLO con il JSON.`, transcript)

	ollamaReq := ollamaRequest{
		Model:  s.OllamaModel,
		Prompt: prompt,
		Stream: false,
	}

	body, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, fmt.Errorf("marshaling ollama request: %w", err)
	}

	url := s.OllamaURL + "/api/generate"
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("ollama request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama returned %d: %s", resp.StatusCode, string(respBody))
	}

	var ollamaResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("decoding ollama response: %w", err)
	}

	// Extract JSON from the response (might have extra text around it)
	jsonStr := extractJSON(ollamaResp.Response)

	var result AIAnalysisResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("parsing analysis JSON: %w (raw: %s)", err, ollamaResp.Response)
	}

	// Clamp score to 0-100
	if result.Score < 0 {
		result.Score = 0
	}
	if result.Score > 100 {
		result.Score = 100
	}

	return &result, nil
}

// extractJSON tries to pull a JSON object from a string that might have extra text.
func extractJSON(s string) string {
	// Try to find JSON between first { and last }
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start >= 0 && end > start {
		return s[start : end+1]
	}
	return s
}

// ReanalyzeCandidate re-runs the AI pipeline for a candidate that already has audio.
func (s *AIService) ReanalyzeCandidate(candidateID string) error {
	var audioPath string
	err := s.DB.Table("candidates").
		Where("id = ?", candidateID).
		Pluck("audio_path", &audioPath).Error
	if err != nil {
		return fmt.Errorf("candidate not found: %w", err)
	}
	if audioPath == "" {
		return fmt.Errorf("no audio file for candidate %s", candidateID)
	}

	// Reset fields
	s.DB.Table("candidates").
		Where("id = ?", candidateID).
		Updates(map[string]interface{}{
			"transcript":      nil,
			"ai_analysis":     nil,
			"ai_score":        nil,
			"analysis_status": "pending",
		})

	s.TranscribeAndAnalyze(candidateID, audioPath)
	return nil
}

// ParseAIScore extracts the score from the response text as a fallback.
func ParseAIScore(text string) int {
	re := regexp.MustCompile(`"score"\s*:\s*(\d+)`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		score, err := strconv.Atoi(matches[1])
		if err == nil {
			return score
		}
	}
	return 0
}
