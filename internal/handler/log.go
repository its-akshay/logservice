package handler

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/logservice/internal/model"
	"github.com/logservice/internal/repo"
)

type LogHandler struct {
	repo repo.LogRepository
}

func NewLogHandler(r repo.LogRepository) *LogHandler {
	return &LogHandler{repo: r}
}

// CreateLogs godoc
// @Summary Create a new log
// @Description Insert a log entry
// @Tags logs
// @Accept json
// @Produce json
// @Param log body model.CreateLogRequest true "Log input"
// @Success 201 {object} model.Log
// @Failure 400 {object} map[string]string
// @Router /logs [post]
func (h *LogHandler) CreateLogs(c *gin.Context) {
	var req model.CreateLogRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request"})
		return
	}

	if req.Level == "" || req.Service == "" || req.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "level, service, message are required"})
		return
	}
	log := &model.Log{
		Level:   req.Level,
		Service: req.Service,
		Message: req.Message,
	}
	if err := h.repo.Insert(c.Request.Context(), log); err != nil {
		c.JSON(500, gin.H{"error": "failed to insert log"})
		return
	}
	c.JSON(http.StatusCreated, log)
}


// CreateLogsBatch godoc
// @Summary Create logs in batch
// @Description Insert multiple log entries (max 500)
// @Tags logs
// @Accept json
// @Produce json
// @Param logs body []model.CreateLogRequest true "List of logs"
// @Success 201 {array} model.Log
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /logs/batch [post]
func (h *LogHandler) CreateLogsBatch(c *gin.Context) {
	var reqs []model.CreateLogRequest

	if err := c.ShouldBindJSON(&reqs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// limit check
	if len(reqs) == 0 || len(reqs) > 500 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "batch size must be between 1 and 500"})
		return
	}

	var logs []model.Log

	for _, r := range reqs {
		// validation
		if r.Level == "" || r.Service == "" || r.Message == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "all fields required in each entry"})
			return
		}

		logs = append(logs, model.Log{
			Level:   r.Level,
			Service: r.Service,
			Message: r.Message,
		})
	}

	// insert one by one (simple approach)
	for i := range logs {
		if err := h.repo.Insert(c.Request.Context(), &logs[i]); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert batch"})
			return
		}
	}

	c.JSON(201, logs)
}

// SearchLogs godoc
// @Summary Search logs
// @Description Filter logs by level/service/regex
// @Tags logs
// @Produce json
// @Param level query string false "Log level"
// @Param service query string false "Service name"
// @Param regex query string false "Regex filter"
// @Param limit query int false "Limit"
// @Success 200 {array} model.Log
// @Router /logs/search [get]
func (h *LogHandler) SearchLogs(c *gin.Context) {
	regexStr := c.Query("regex")
	level := c.Query("level")
	service := c.Query("service")

	limit := 100
	if l := c.Query("limit"); l != "" {
		val, err := strconv.Atoi(l)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid limit"})
			return
		}
		limit = val
	}
	var re *regexp.Regexp
	if regexStr != "" {
		var err error
		re, err = regexp.Compile(regexStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid regex"})
			return
		}
	}

	logs, err := h.repo.List(c.Request.Context(), repo.Filter{
		Level:   level,
		Service: service,
		Limit:   limit,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "db error"})
	}

	

	if re != nil {
		var result []model.Log
		for _, l := range logs {
			if re.MatchString(l.Message) {
				result = append(result, l)
			}
			
		}
		logs = result
	}
	c.JSON(http.StatusOK, logs)
}
