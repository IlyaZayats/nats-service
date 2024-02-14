package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/IlyaZayats/servord/internal/interfaces"
	"github.com/IlyaZayats/servord/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/stan.go"
	"net/http"
)

type NatsHandler struct {
	s  *services.OrderService
	mc interfaces.Cache
	ns *stan.Conn
}

type OrderOutput struct {
	data string
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
		//logrus.Println(string(cacheList[i]))
		err := n.mc.Set(keys[i], cacheList[i])
		if err != nil {
			return err
		}
	}
	var sub stan.Subscription
	go func() {
		sub, err = (*n.ns).Subscribe("orders", func(m *stan.Msg) {
			buff := new(bytes.Buffer)
			if err := json.Compact(buff, m.Data); err != nil {
				fmt.Println(err)
			}
			key, err := n.s.InsertOrder(buff.Bytes())
			if err != nil {
				fmt.Println(err.Error())
			} else {
				err := n.mc.Set(key, m.Data)
				if err != nil {
					doneC <- err
					return
				}
			}

		})
		for {
			if err != nil {
				doneC <- err
			}
		}
		sub.Unsubscribe()
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order, err := h.cache.Get(request.OrderUid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var output map[string]interface{}
	err = json.Unmarshal(order, &output)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//logrus.Println(dude)
	//var output OrderOutput
	//err = json.Unmarshal([]byte(dude), &output)
	//logrus.Println(output.data)
	c.JSON(http.StatusOK, output)
}
