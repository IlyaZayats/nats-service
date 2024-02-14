package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/IlyaZayats/servord/internal/interfaces"
	"github.com/IlyaZayats/servord/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"net/http"
)

type NatsHandler struct {
	s  *services.OrderService
	mc interfaces.Cache
	ns *stan.Conn
}

func NewNatsHandler(s *services.OrderService, mc interfaces.Cache, ns *stan.Conn) (*NatsHandler, error) {
	return &NatsHandler{s: s, mc: mc, ns: ns}, nil
}

func (n *NatsHandler) Run() error {
	doneC := make(chan error)
	keys, cacheList, err := n.s.InitCache()
	if err != nil {
		return err
	}
	for i := range cacheList {
		err := n.mc.Set(keys[i], cacheList[i])
		if err != nil {
			return err
		}
	}
	var sub stan.Subscription
	go func() {
		sub, err = (*n.ns).Subscribe("orders", func(m *stan.Msg) {
			logrus.Println("Got new nats-streaming msg!")
			buff := new(bytes.Buffer)
			if err := json.Compact(buff, m.Data); err != nil {
				logrus.Error(err.Error())
				return
			}
			key, err := n.s.InsertOrder(buff.Bytes())
			if err != nil {
				logrus.Error(err.Error())
			} else if err := n.mc.Set(key, m.Data); err != nil {
				logrus.Error(err.Error())
				doneC <- err
				return
			}
		})
		for {
			if err != nil {
				if errSub := sub.Unsubscribe(); errSub != nil {
					logrus.Error(err.Error())
				}
				doneC <- err
			}
		}
	}()
	return <-doneC
}

type OrderHandler struct {
	engine *gin.Engine
	cache  interfaces.Cache
}

func NewOrderHandler(engine *gin.Engine, cache interfaces.Cache) (*OrderHandler, error) {
	h := &OrderHandler{
		engine: engine,
		cache:  cache,
	}
	h.InitRoute()
	return h, nil
}

func (h *OrderHandler) InitRoute() {
	h.engine.GET("/order/:order_uid", h.GetOrder)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	var request GetOrderRequest
	if err := c.BindUri(&request); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Println(fmt.Sprintf("Got new order request! ID: %s", request.OrderUid))
	order, err := h.cache.Get(request.OrderUid)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var output map[string]interface{}
	err = json.Unmarshal(order, &output)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, output)
}
