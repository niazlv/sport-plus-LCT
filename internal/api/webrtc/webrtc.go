package webrtc

import (
	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/pion/webrtc/v3"
	"github.com/wI2L/fizz"
)

var (
	PeerConnection *webrtc.PeerConnection
	DataChannel    *webrtc.DataChannel
)

func Setup(rg *fizz.RouterGroup) {
	api := rg.Group("webrtc", "WebRTC", "WebRTC related endpoints")
	InitWebRTC()

	api.POST("/offer", []fizz.OperationOption{fizz.Summary("Create an SDP offer")}, tonic.Handler(CreateOffer, 200))
	api.POST("/answer", []fizz.OperationOption{fizz.Summary("Set an SDP answer")}, tonic.Handler(SetAnswer, 200))
}

func InitWebRTC() {
	// Создаем новое соединение PeerConnection
	var err error
	PeerConnection, err = webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		panic(err)
	}

	// Создаем новый DataChannel
	DataChannel, err = PeerConnection.CreateDataChannel("data", nil)
	if err != nil {
		panic(err)
	}
}

func CreateOffer(c *gin.Context) (*webrtc.SessionDescription, error) {
	// Создаем SDP предложение
	offer, err := PeerConnection.CreateOffer(nil)
	if err != nil {
		return nil, err
	}

	// Устанавливаем локальное SDP
	err = PeerConnection.SetLocalDescription(offer)
	if err != nil {
		return nil, err
	}

	// Возвращаем SDP предложение клиенту
	return &offer, nil
}

func SetAnswer(c *gin.Context) error {
	var answer webrtc.SessionDescription
	if err := c.ShouldBindJSON(&answer); err != nil {
		return err
	}

	// Устанавливаем удаленное SDP
	err := PeerConnection.SetRemoteDescription(answer)
	if err != nil {
		return err
	}

	return nil
}
