package main

import (
	"context"
	"encoding/json"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type TwitchWS struct {
	con    *websocket.Conn
	ctx    context.Context
	server *Server
}

func InintTwitchWS(server *Server) *TwitchWS {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "wss://ws.mjrlegends.com:2096", nil)
	if err != nil {
		log.Printf("error: %v", err)
	}

	s := &TwitchWS{
		con:    c,
		ctx:    ctx,
		server: server,
	}

	s.connect()

	go s.readMessages()

	return s
}

func (c *TwitchWS) connect() {
	err := wsjson.Write(c.ctx, c.con, MJRRequest{
		Type:      "LISTEN",
		Nonce:     "4jgUaUv0zdxBMe2tN6YSZaCROCwkO92baSaFzgT50sWFySI15ErkVpoIqfqLwoZ6",
		ChannelId: 32907202,
		Topics:    []string{"channel_points_reward_redeem"},
		Token:     MJR_TOKEN,
	})
	if err != nil {
		log.Printf("error: %v\n", err)
	}
}

func (c *TwitchWS) getColorFromMsg(msg string) *Color {
	// if msg in color_map.keys():
	// 	return color_map[msg]
	// else:
	r, _ := regexp.Compile("#[a-f0-9]{6}")
	if r.MatchString(msg) {
		usr_input := msg[1:]
		return &Color{
			r: getUInt8(usr_input[0:2]),
			g: getUInt8(usr_input[2:4]),
			b: getUInt8(usr_input[4:6]),
		}
	}
	return nil
}

func getUInt8(str string) uint8 {
	i, err := strconv.ParseInt(str, 16, 8)
	if err != nil {
		panic(err)
	}
	return uint8(i)
}

func (c *TwitchWS) readMessages() {
	defer c.con.Close(websocket.StatusInternalError, "the sky is falling")
	for {
		// Read in a new message as JSON and map it to a Message object
		_, bytes, err := c.con.Read(context.Background())
		if err != nil {
			log.Printf("error: %v\n", err)
			break
		}

		var message MJRResponse
		if err := json.Unmarshal(bytes, &message); err != nil {
			log.Println("error:", err)
		}

		if message.Type == "MESSAGE" && message.Topic == "channel_points_reward_redeem" {
			var data ChannePointRedeemResponse
			if err := json.Unmarshal(bytes, &data); err != nil {
				log.Println("error:", err)
			}

			redemption := data.Message.Redemption
			if redemption.Reward.Id == "c63fb418-8463-4a95-8fb5-04ffac7b964e" {
				usrInput := strings.TrimSpace(strings.ToLower(redemption.UserInput))
				color := c.getColorFromMsg(usrInput)
				if color != nil {
					c.server.leds.display = SOLID
					c.server.leds.colors = []Color{*color}
				} else if strings.HasPrefix(usrInput, "rainbow") {
					c.server.leds.display = RAINBOW
				} else if strings.HasPrefix(usrInput, "colorblocks") {
					c.server.leds.display = BLOCK_COLOR
					c.server.leds.colors = []Color{}
					// 	pot_colors = usr_input.split(' ')
					// 	for col in pot_colors:
					// 		color = self.get_color_from_msg(col)
					// 		if color is not None:
					// 			data.colors.append(color)
				} else if strings.HasPrefix(usrInput, "coloralternate") {
					c.server.leds.display = ALTERNATE_COLOR
					// 	data.colors = []
					// 	pot_colors = usr_input.split(' ')
					// 	for col in pot_colors:
					// 		color = self.get_color_from_msg(col)
					// 		if color is not None:
					// 			data.colors.append(color)
				} else if strings.HasPrefix(usrInput, "police") {
					c.server.leds.display = POLICE
				}
			}
		}
	}
}
