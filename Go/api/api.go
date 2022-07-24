package api

import (
	"example.com/kafka-serializer-publisher/marshaller"
	"example.com/kafka-serializer-publisher/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	publisher model.Publisher
)

// PostJson godoc
// @Summary      Sends message to topic
// @Description  Sends JSON formatted payload to specified topic
// @Tags         Producer
// @Accept       json
// @Produce      plain
// @Param   Payload body model.JsonRequest{payload=object} true "request payload"
// @Success      204
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      500  {object}  Data  "Error sending json payload"
// @Router       /json [post]
func PostJson(c *gin.Context) {
	var req model.JsonRequest
	err := c.BindJSON(&req)
	if err != nil {
		ReplyWithError(c, http.StatusBadRequest, err)
		return
	}
	doIt(c, marshaller.Json, &req)
}

// PostAvro godoc
// @Summary      Sends message to topic
// @Description  Sends AVRO formatted payload to specified topic
// @Tags         Producer
// @Accept       json
// @Produce      plain
// @Param   Payload body model.AvroRequest true "request payload"
// @Success      204
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      500  {object}  Data  "Error sending avro payload"
// @Router       /avro [post]
func PostAvro(c *gin.Context) {
	var req model.AvroRequest
	err := c.BindJSON(&req)
	if err != nil {
		ReplyWithError(c, http.StatusBadRequest, err)
		return
	}
	doIt(c, marshaller.Avro, &req)
}

func doIt(c *gin.Context, contentType model.EventContentType, req model.Publishable) {
	err := publisher.Send(c.Request.Context(), contentType, req)
	if err != nil {
		ReplyWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// CreateRoutes - Creates the needed routes and returns and function to shutdown as needed
func CreateRoutes(router *gin.Engine, sndr model.Publisher) {
	publisher = sndr
	pv1 := router.Group("/")
	{
		pv1.POST("json", PostJson)
		pv1.POST("avro", PostAvro)
	}
}

type Data struct {
	Status int
	Meta   interface{}
	Data   interface{}
}

func ReplyWithError(c *gin.Context, code int, err error) {
	var res Data
	res.Status = code
	res.Data = err.Error()
	c.JSON(code, res)
}
